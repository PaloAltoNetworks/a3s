package authcmd

import (
	"context"

	"github.com/spf13/cobra"
	"github.com/spf13/viper"
	"go.aporeto.io/a3s/pkgs/authlib"
	"go.aporeto.io/a3s/pkgs/permissions"
	"go.aporeto.io/manipulate/manipcli"
)

func makeGCPCmd(mmaker manipcli.ManipulatorMaker, restrictions *permissions.Restrictions) *cobra.Command {

	cmd := &cobra.Command{
		Use:              "gcp",
		Short:            "Use an GCP identity token.",
		TraverseChildren: true,
		RunE: func(cmd *cobra.Command, args []string) error {

			fToken := viper.GetString("access-token")
			fTokenAudience := viper.GetString("token-audience")
			fAudience := viper.GetStringSlice("audience")
			fCloak := viper.GetStringSlice("cloak")
			fQRCode := viper.GetBool("qrcode")

			m, err := mmaker()
			if err != nil {
				return err
			}

			client := authlib.NewClient(m)
			t, err := client.AuthFromGCP(
				context.Background(),
				fToken,
				fTokenAudience,
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

	cmd.Flags().String("access-token", "", "Valid GCP token.")
	cmd.Flags().String("token-audience", "a3s", "Required GCP token audience.")

	return cmd
}
