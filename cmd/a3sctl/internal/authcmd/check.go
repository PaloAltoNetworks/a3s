package authcmd

import (
	"fmt"
	"os"

	"github.com/spf13/cobra"
	"github.com/spf13/viper"
	"go.aporeto.io/a3s/pkgs/token"
	"go.aporeto.io/manipulate/manipcli"
)

func makeCheckCmd(mmaker manipcli.ManipulatorMaker) *cobra.Command {

	cmd := &cobra.Command{
		Use:              "check",
		Aliases:          []string{"verify"},
		Short:            "Check the token",
		TraverseChildren: true,
		PersistentPreRunE: func(cmd *cobra.Command, args []string) error {
			if err := cmd.Root().PersistentPreRunE(cmd, args); err != nil {
				return err
			}
			if err := HandleAutoAuth(
				mmaker,
				"",
				nil,
				nil,
				false,
				false,
			); err != nil {
				return fmt.Errorf("auto auth error: %w", err)
			}

			return nil
		},
		RunE: func(cmd *cobra.Command, args []string) error {
			fToken := viper.GetString("token")
			fCheck := viper.GetBool("check")
			fQRCode := viper.GetBool("qrcode")

			return token.Fprint(
				os.Stdout,
				fToken,
				token.PrintOptionDecoded(fCheck),
				token.PrintOptionQRCode(fQRCode),
				token.PrintOptionRaw(true),
			)
		},
	}

	cmd.Flags().Bool("print", false, "Print the token string.")

	cmd.SetHelpFunc(func(cmd *cobra.Command, args []string) {
		_ = cmd.Flags().MarkHidden("namespace")
		_ = cmd.Flags().MarkHidden("audience")
		_ = cmd.Flags().MarkHidden("cloak")
		_ = cmd.Flags().MarkHidden("validity")
		_ = cmd.Flags().MarkHidden("encoding")
		_ = cmd.Flags().MarkHidden("restrict-namespace")
		_ = cmd.Flags().MarkHidden("restrict-permissions")
		_ = cmd.Flags().MarkHidden("restrict-network")
		_ = cmd.Flags().MarkHidden("refresh")
		cmd.Parent().HelpFunc()(cmd, args)
	})

	return cmd
}
