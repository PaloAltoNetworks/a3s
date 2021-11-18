package authcmd

import (
	"context"

	"github.com/spf13/cobra"
	"github.com/spf13/viper"
	"go.aporeto.io/a3s/pkgs/authlib"
	"go.aporeto.io/manipulate/manipcli"
	"go.aporeto.io/manipulate/maniphttp"
	"go.aporeto.io/tg/tglib"
)

func makeMTLSCmd(mmaker manipcli.ManipulatorMaker) *cobra.Command {

	cmd := &cobra.Command{
		Use:              "mtls",
		Short:            "Use a configured MTLS authentication source.",
		TraverseChildren: true,
		RunE: func(cmd *cobra.Command, args []string) error {

			fCert := viper.GetString("cert")
			fKey := viper.GetString("key")
			fPass := viper.GetString("pass")
			fSourceName := viper.GetString("source-name")
			fSourceNamespace := viper.GetString("source-namespace")
			fAudience := viper.GetStringSlice("audience")
			fCloak := viper.GetStringSlice("cloak")
			fQRCode := viper.GetBool("qrcode")

			cert, key, err := tglib.ReadCertificatePEM(fCert, fKey, fPass)
			if err != nil {
				return err
			}

			clientCert, err := tglib.ToTLSCertificate(cert, key)
			if err != nil {
				return err
			}

			m, err := mmaker(maniphttp.OptionTLSClientCertificates(clientCert))
			if err != nil {
				return err
			}

			client := authlib.NewClient(m)
			t, err := client.AuthFromCertificate(
				context.Background(),
				fSourceNamespace,
				fSourceName,
				authlib.OptAudience(fAudience...),
				authlib.OptCloak(fCloak...),
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
	_ = cobra.MarkFlagRequired(cmd.Flags(), "cert")
	_ = cobra.MarkFlagRequired(cmd.Flags(), "key")

	return cmd
}
