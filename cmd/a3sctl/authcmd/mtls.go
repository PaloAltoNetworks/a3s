package authcmd

import (
	"context"
	"time"

	"github.com/spf13/cobra"
	"github.com/spf13/viper"
	"go.aporeto.io/a3s/cmd/a3sctl/helpers"
	"go.aporeto.io/a3s/pkgs/authlib"
	"go.aporeto.io/a3s/pkgs/permissions"
	"go.aporeto.io/manipulate/manipcli"
	"go.aporeto.io/manipulate/maniphttp"
	"go.aporeto.io/tg/tglib"
)

func makeMTLSCmd(mmaker manipcli.ManipulatorMaker, restrictions *permissions.Restrictions) *cobra.Command {

	cmd := &cobra.Command{
		Use:              "mtls",
		Short:            "Use a configured MTLS authentication source.",
		TraverseChildren: true,
		RunE: func(cmd *cobra.Command, args []string) error {

			fCert := viper.GetString("cert")
			fKey := viper.GetString("key")
			fPass := helpers.ReadFlag("passphrase: ", "pass", true)
			fSourceName := viper.GetString("source-name")
			fSourceNamespace := viper.GetString("source-namespace")
			fAudience := viper.GetStringSlice("audience")
			fCloak := viper.GetStringSlice("cloak")
			fQRCode := viper.GetBool("qrcode")
			fValidity := viper.GetDuration("validity")

			if fSourceNamespace == "" {
				fSourceNamespace = viper.GetString("namespace")
			}

			t, err := GetMTLSToken(
				mmaker,
				fCert,
				fKey,
				fPass,
				fSourceNamespace,
				fSourceName,
				fAudience,
				fCloak,
				fValidity,
				restrictions,
			)
			if err != nil {
				return err
			}

			printToken(t, fQRCode)

			return nil
		},
	}

	cmd.Flags().String("cert", "", "Path to the certificate in PEM format.")
	cmd.Flags().String("key", "", "Path to the certificate key in PEM format.")
	cmd.Flags().String("pass", "", "Passphrase for the certificate key.")

	return cmd
}

// GetMTLSToken retrieves a token using
// the provided MTLS source information.
func GetMTLSToken(
	mmaker manipcli.ManipulatorMaker,
	certPath string,
	keyPath string,
	keyPass string,
	sourceNamespace string,
	sourceName string,
	audience []string,
	cloak []string,
	validity time.Duration,
	restrictions *permissions.Restrictions,
) (string, error) {

	cert, key, err := tglib.ReadCertificatePEM(certPath, keyPath, keyPass)
	if err != nil {
		return "", err
	}

	clientCert, err := tglib.ToTLSCertificate(cert, key)
	if err != nil {
		return "", err
	}

	m, err := mmaker(maniphttp.OptionTLSClientCertificates(clientCert))
	if err != nil {
		return "", err
	}

	opts := []authlib.Option{
		authlib.OptAudience(audience...),
		authlib.OptCloak(cloak...),
		authlib.OptValidity(validity),
	}

	if restrictions != nil {
		opts = append(opts,
			authlib.OptRestrictions(*restrictions),
		)
	}

	client := authlib.NewClient(m)
	return client.AuthFromCertificate(
		context.Background(),
		sourceNamespace,
		sourceName,
		opts...,
	)
}
