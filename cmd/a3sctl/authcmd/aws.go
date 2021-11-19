package authcmd

import (
	"context"

	"github.com/spf13/cobra"
	"github.com/spf13/viper"
	"go.aporeto.io/a3s/pkgs/authlib"
	"go.aporeto.io/a3s/pkgs/permissions"
	"go.aporeto.io/manipulate/manipcli"
)

func makeAWSCmd(mmaker manipcli.ManipulatorMaker, restrictions *permissions.Restrictions) *cobra.Command {

	cmd := &cobra.Command{
		Use:              "aws",
		Short:            "Use an AWS identity token.",
		TraverseChildren: true,
		RunE: func(cmd *cobra.Command, args []string) error {

			fToken := viper.GetString("access-token")
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
				authlib.OptRestrictions(*restrictions),
			)
			if err != nil {
				return err
			}

			printToken(t, fQRCode)

			return nil
		},
	}

	cmd.Flags().String("access-token", "", "Valid AWS token.")
	cmd.Flags().String("access-key-id", "", "Access key ID for the token.")
	cmd.Flags().String("access-key-secret", "", "Secret for access key.")

	return cmd
}
