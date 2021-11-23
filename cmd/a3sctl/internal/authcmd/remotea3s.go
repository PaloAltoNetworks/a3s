package authcmd

import (
	"context"

	"github.com/spf13/cobra"
	"github.com/spf13/viper"
	"go.aporeto.io/a3s/pkgs/authlib"
	"go.aporeto.io/a3s/pkgs/permissions"
	"go.aporeto.io/manipulate/manipcli"
)

func makeRemoteA3SCmd(mmaker manipcli.ManipulatorMaker, restrictions *permissions.Restrictions) *cobra.Command {

	cmd := &cobra.Command{
		Use:              "remote-a3s",
		Short:            "Use a remote A3S identity token.",
		TraverseChildren: true,
		RunE: func(cmd *cobra.Command, args []string) error {

			fToken := viper.GetString("access-token")
			fSourceName := viper.GetString("source-name")
			fSourceNamespace := viper.GetString("source-namespace")
			fAudience := viper.GetStringSlice("audience")
			fCloak := viper.GetStringSlice("cloak")
			fQRCode := viper.GetBool("qrcode")
			fValidity := viper.GetDuration("validity")

			if fToken == "" {
				fToken = viper.GetString("token")
			}

			m, err := mmaker()
			if err != nil {
				return err
			}

			client := authlib.NewClient(m)
			t, err := client.AuthFromRemoteA3S(
				context.Background(),
				fToken,
				fSourceNamespace,
				fSourceName,
				authlib.OptAudience(fAudience...),
				authlib.OptCloak(fCloak...),
				authlib.OptRestrictions(*restrictions),
				authlib.OptValidity(fValidity),
			)
			if err != nil {
				return err
			}

			printToken(t, fQRCode)

			return nil
		},
	}

	cmd.Flags().String("access-token", "", "Valid remote a3s token.")

	return cmd
}
