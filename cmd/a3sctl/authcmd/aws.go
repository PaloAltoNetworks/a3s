package authcmd

import (
	"context"

	"github.com/spf13/cobra"
	"github.com/spf13/viper"
	"go.aporeto.io/a3s/pkgs/authlib"
	"go.aporeto.io/manipulate/manipcli"
)

func makeAWSCmd(mmaker manipcli.ManipulatorMaker) *cobra.Command {

	cmd := &cobra.Command{
		Use:              "aws",
		Short:            "Use an AWS identity token.",
		TraverseChildren: true,
		RunE: func(cmd *cobra.Command, args []string) error {

			fToken := viper.GetString("token")
			fAccessKeyID := viper.GetString("access-key-id")
			fSecretAccessKey := viper.GetString("access-key-secret")
			fAudience := viper.GetStringSlice("audience")
			fCloak := viper.GetStringSlice("cloak")
			fQRCode := viper.GetBool("qrcode")

			m, err := mmaker()
			if err != nil {
				return err
			}

			client := authlib.NewClient(m)
			t, err := client.AuthFromAWS(
				context.Background(),
				fAccessKeyID,
				fSecretAccessKey,
				fToken,
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

	cmd.Flags().String("token", "", "Valid Azure token.")
	cmd.Flags().String("access-key-id", "", "Access key ID for the token.")
	cmd.Flags().String("access-key-secret", "", "Secret for access key.")

	return cmd
}
