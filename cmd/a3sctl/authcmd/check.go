package authcmd

import (
	"fmt"

	"github.com/golang-jwt/jwt/v4"
	"github.com/hokaccha/go-prettyjson"
	"github.com/spf13/cobra"
	"github.com/spf13/viper"
)

func makeCheckCmd() *cobra.Command {

	cmd := &cobra.Command{
		Use:              "check",
		Short:            "Check the token",
		TraverseChildren: true,
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
	_ = cobra.MarkFlagRequired(cmd.Flags(), "token")

	return cmd
}
