package authcmd

import (
	"fmt"

	"github.com/golang-jwt/jwt/v4"
	"github.com/hokaccha/go-prettyjson"
	"github.com/spf13/cobra"
	"github.com/spf13/viper"
	"go.aporeto.io/a3s/pkgs/api"
	"go.aporeto.io/manipulate/manipcli"
	"go.aporeto.io/manipulate/maniphttp"
	"go.aporeto.io/tg/tglib"
)

func init() {

	authMTLSCmd.Flags().String("cert", "", "Path to the certificate in PEM format.")
	authMTLSCmd.Flags().String("key", "", "Path to the certificate key in PEM format.")
	authMTLSCmd.Flags().String("pass", "", "Passphrase for the certificate key.")

	AuthCmd.PersistentFlags().AddFlagSet(manipcli.ManipulatorFlagSet())
	AuthCmd.PersistentFlags().StringSlice("audience", nil, "Requested audience for the token")
	AuthCmd.PersistentFlags().String("source-name", "default", "The name of the auth source")
	AuthCmd.PersistentFlags().String("source-namespace", "/", "The namespace of the auth source")

	authCheckCmd.Flags().String("token", "", "The token to verify")

	AuthCmd.AddCommand(
		authCheckCmd,
		authMTLSCmd,
	)
}

// AuthCmd is the root cmd for authentication subcommand.
var AuthCmd = &cobra.Command{
	Use:   "auth",
	Short: "Authenticate and retrieve a token.",
}

var authCheckCmd = &cobra.Command{
	Use:   "check",
	Short: "Check the token",
	PersistentPreRunE: func(cmd *cobra.Command, args []string) error {
		if err := viper.BindPFlags(cmd.PersistentFlags()); err != nil {
			return err
		}
		return viper.BindPFlags(cmd.Flags())
	},
	RunE: func(cmd *cobra.Command, args []string) error {

		claims := jwt.MapClaims{}
		p := jwt.Parser{}

		t, _, err := p.ParseUnverified(viper.GetString("token"), &claims)
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

var authMTLSCmd = &cobra.Command{
	Use:   "mtls",
	Short: "Use a client certificate to obtain a token.",
	PersistentPreRunE: func(cmd *cobra.Command, args []string) error {
		if err := viper.BindPFlags(cmd.PersistentFlags()); err != nil {
			return err
		}
		return viper.BindPFlags(cmd.Flags())
	},
	RunE: func(cmd *cobra.Command, args []string) error {

		cert, key, err := tglib.ReadCertificatePEM(
			viper.GetString("cert"),
			viper.GetString("key"),
			viper.GetString("pass"),
		)
		if err != nil {
			return err
		}

		clientCert, err := tglib.ToTLSCertificate(cert, key)
		if err != nil {
			return err
		}

		m, err := manipcli.ManipulatorMakerFromFlags(
			maniphttp.OptionTLSClientCertificates(clientCert),
		)()
		if err != nil {
			return err
		}

		iss := api.NewIssue()
		iss.SourceName = viper.GetString("source-name")
		iss.SourceNamespace = viper.GetString("source-namespace")
		iss.SourceType = api.IssueSourceTypeMTLS
		iss.Audience = viper.GetStringSlice("audience")

		if err := m.Create(nil, iss); err != nil {
			return err
		}

		fmt.Println(iss.Token)

		return nil
	},
}
