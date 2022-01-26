package authcmd

import (
	"fmt"
	"time"

	"github.com/golang-jwt/jwt/v4"
	"github.com/hokaccha/go-prettyjson"
	"github.com/spf13/cobra"
	"github.com/spf13/viper"
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
			fPrint := viper.GetBool("print")
			fQRCode := viper.GetBool("qrcode")

			return DisplayToken(fToken, fPrint, fQRCode)
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

// DisplayToken displays information about the token.
// It does not verify its signature or validity.
func DisplayToken(token string, printRaw bool, qrcode bool) error {
	claims := jwt.MapClaims{}
	p := jwt.Parser{}

	t, _, err := p.ParseUnverified(token, &claims)
	if err != nil {
		return err
	}

	data, err := prettyjson.Marshal(claims)
	if err != nil {
		return err
	}

	fmt.Println("alg:", t.Method.Alg())
	fmt.Println("kid:", t.Header["kid"])
	if exp, ok := claims["exp"].(float64); ok {

		remaining := time.Until(time.Unix(int64(exp), 0))
		if remaining <= 0 {
			fmt.Println("exp: the token has expired", -remaining.Truncate(time.Second), "ago")
		} else {
			fmt.Println("exp:", remaining.Truncate(time.Second))
		}
	}
	fmt.Println()

	fmt.Println(string(data))

	if printRaw {
		fmt.Println()
		printToken(token, qrcode)
	}

	return nil
}
