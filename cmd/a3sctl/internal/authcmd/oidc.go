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
	"go.aporeto.io/a3s/pkgs/token"
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
			fCheck := viper.GetBool("check")
			fRefresh := viper.GetBool("refresh")

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
				fSourceNamespace,
				fSourceName,
				authD.code,
				authD.state,
				authlib.OptAudience(fAudience...),
				authlib.OptCloak(fCloak...),
				authlib.OptRestrictions(*restrictions),
				authlib.OptRefresh(fRefresh),
			)
			if err != nil {
				return err
			}

			return token.Fprint(
				os.Stdout,
				t,
				token.PrintOptionDecoded(fCheck),
				token.PrintOptionQRCode(fQRCode),
				token.PrintOptionRaw(true),
			)
		},
	}
	cmd.Flags().String("source-name", "default", "The name of the auth source.")
	cmd.Flags().String("source-namespace", "", "The namespace of the auth source. If omitted, uses --namespace.")

	cmd.SetHelpFunc(func(cmd *cobra.Command, args []string) {
		_ = cmd.Flags().MarkHidden("token")
		cmd.Parent().HelpFunc()(cmd, args)
	})

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
		_, _ = w.Write([]byte(`Authenticated. You can close this window`))
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

	_ = server.Shutdown(shutdownCtx) // nolint
}
