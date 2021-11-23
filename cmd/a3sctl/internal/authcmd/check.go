package authcmd

import (
	"fmt"

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
				viper.GetString("auto-auth-method"),
				viper.GetBool("refresh-cached-token"),
			); err != nil {
				return fmt.Errorf("auto auth error: %w", err)
			}
			return nil
		},
		RunE: func(cmd *cobra.Command, args []string) error {
			fToken := viper.GetString("token")

			claims := jwt.MapClaims{}
			p := jwt.Parser{}

			t, _, err := p.ParseUnverified(fToken, &claims)
			if err != nil {
				return err
			}

			data, err := prettyjson.Marshal(claims)
			if err != nil {
				return err
			}

			fmt.Println("alg:", t.Method.Alg())
			fmt.Println("kid:", t.Header["kid"])
			fmt.Println()

			fmt.Println(string(data))

			return nil
		},
	}

	cmd.Flags().String("token", "", "The token to verify.")

	return cmd
}
