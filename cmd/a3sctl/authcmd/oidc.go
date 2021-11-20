package authcmd

import (
	"context"
	"fmt"
	"net/http"
	"os"
	"time"

	"github.com/spf13/cobra"
	"github.com/spf13/viper"
	"go.aporeto.io/a3s/pkgs/authlib"
	"go.aporeto.io/a3s/pkgs/permissions"
	"go.aporeto.io/manipulate/manipcli"
	"go.uber.org/zap"
)

type oidcAuthData struct {
	code  string
	state string
}

func makeOIDCCmd(mmaker manipcli.ManipulatorMaker, restrictions *permissions.Restrictions) *cobra.Command {

	cmd := &cobra.Command{
		Use:   "oidc",
		Short: "Authenticate using OIDC source.",
		RunE: func(cmd *cobra.Command, args []string) error {

			fSourceName := viper.GetString("source-name")
			fSourceNamespace := viper.GetString("source-namespace")
			fAudience := viper.GetStringSlice("audience")
			fCloak := viper.GetStringSlice("cloak")
			fQRCode := viper.GetBool("qrcode")

			if fSourceNamespace == "" {
				fSourceNamespace = viper.GetString("namespace")
			}

			srvCtx, srvCancel := context.WithCancel(context.Background())
			defer srvCancel()

			authDataCh := make(chan oidcAuthData)

			go startOIDCCallbackServer(srvCtx, authDataCh)

			ctx, cancel := context.WithTimeout(context.Background(), 30*time.Second)
			defer cancel()

			m, err := mmaker()
			if err != nil {
				return err
			}

			client := authlib.NewClient(m)
			url, err := client.AuthFromOIDCStep1(
				ctx,
				fSourceNamespace,
				fSourceName,
				"http://localhost:65333",
			)
			if err != nil {
				return err
			}

			fmt.Fprintln(os.Stderr, "Open this URL in your browser:", url)

			authD := <-authDataCh
			srvCancel()

			t, err := client.AuthFromOIDCStep2(
				ctx,
				authD.code,
				authD.state,
				authlib.OptAudience(fAudience...),
				authlib.OptCloak(fCloak...),
				authlib.OptRestrictions(*restrictions),
			)
			if err != nil {
				return err
			}

			printToken(t, fQRCode)

			return nil
		},
	}

	return cmd
}

func startOIDCCallbackServer(srvCtx context.Context, out chan oidcAuthData) {

	mux := http.NewServeMux()
	mux.HandleFunc("/", http.HandlerFunc(func(w http.ResponseWriter, req *http.Request) {

		out <- oidcAuthData{
			code:  req.URL.Query().Get("code"),
			state: req.URL.Query().Get("state"),
		}

		w.WriteHeader(http.StatusOK)
		w.Write([]byte(`Authenticated. You can close this window`)) // nolint
	}))

	server := &http.Server{
		Addr:    "localhost:65333",
		Handler: mux,
	}

	go func() {
		if err := server.ListenAndServe(); err != nil {
			if err != http.ErrServerClosed {
				zap.L().Fatal("Unable to start temporary webserver", zap.Error(err))
			}
		}
	}()

	<-srvCtx.Done()

	shutdownCtx, cancel := context.WithTimeout(context.Background(), 1*time.Second)
	defer cancel()

	server.Shutdown(shutdownCtx) // nolint
}
