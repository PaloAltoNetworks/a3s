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

// New returns a new auth command.
func New(mmaker manipcli.ManipulatorMaker) *cobra.Command {

	// rootCmd is the root cmd for authentication subcommand.
	rootCmd := &cobra.Command{
		Use:              "auth",
		Short:            "Authenticate and retrieve a token.",
		TraverseChildren: true,
	}
	rootCmd.PersistentFlags().StringSlice("audience", nil, "Requested audience for the token.")
	rootCmd.PersistentFlags().String("source-name", "default", "The name of the auth source.")
	rootCmd.PersistentFlags().String("source-namespace", "/", "The namespace of the auth source.")
	viper.BindPFlags(rootCmd.Flags())

	checkCmd := &cobra.Command{
		Use:              "check",
		Short:            "Check the token",
		TraverseChildren: true,
		RunE: func(cmd *cobra.Command, args []string) error {

			flags := cmd.Flags()
			fToken, _ := flags.GetString("token")

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
	checkCmd.Flags().String("token", "", "The token to verify.")
	cobra.MarkFlagRequired(checkCmd.Flags(), "token")
	viper.BindPFlags(checkCmd.Flags())

	mtlsCmd := &cobra.Command{
		Use:              "mtls",
		Short:            "Use a configured MTLS authentication source.",
		TraverseChildren: true,
		RunE: func(cmd *cobra.Command, args []string) error {

			flags := cmd.Flags()
			fCert, _ := flags.GetString("cert")
			fKey, _ := flags.GetString("key")
			fPass, _ := flags.GetString("pass")
			fSourceName, _ := flags.GetString("source-name")
			fSourceNamespace, _ := flags.GetString("source-namespace")
			fAudience, _ := flags.GetStringSlice("audience")

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

			iss := api.NewIssue()
			iss.SourceType = api.IssueSourceTypeMTLS
			iss.SourceName = fSourceName
			iss.SourceNamespace = fSourceNamespace
			iss.Audience = fAudience

			if err := m.Create(nil, iss); err != nil {
				return err
			}

			fmt.Println(iss.Token)

			return nil
		},
	}
	mtlsCmd.Flags().String("cert", "", "Path to the certificate in PEM format.")
	mtlsCmd.Flags().String("key", "", "Path to the certificate key in PEM format.")
	mtlsCmd.Flags().String("pass", "", "Passphrase for the certificate key.")
	cobra.MarkFlagRequired(mtlsCmd.Flags(), "cert")
	cobra.MarkFlagRequired(mtlsCmd.Flags(), "key")
	viper.BindPFlags(mtlsCmd.Flags())

	ldapCmd := &cobra.Command{
		Use:              "ldap",
		Short:            "Use a configured LDAP authentication source.",
		TraverseChildren: true,
		RunE: func(cmd *cobra.Command, args []string) error {

			flags := cmd.Flags()
			fSourceName, _ := flags.GetString("source-name")
			fSourceNamespace, _ := flags.GetString("source-namespace")
			fAudience, _ := flags.GetStringSlice("audience")
			fUser, _ := flags.GetString("user")
			fPass, _ := flags.GetString("pass")

			iss := api.NewIssue()
			iss.SourceType = api.IssueSourceTypeLDAP
			iss.SourceName = fSourceName
			iss.SourceNamespace = fSourceNamespace
			iss.Audience = fAudience
			iss.InputLDAP = &api.IssueLDAP{
				Username: fUser,
				Password: fPass,
			}

			m, err := mmaker()
			if err != nil {
				return err
			}

			if err := m.Create(nil, iss); err != nil {
				return err
			}

			fmt.Println(iss.Token)

			return nil
		},
	}
	ldapCmd.Flags().String("user", "", "The LDAP username to use.")
	ldapCmd.Flags().String("pass", "", "The password associateds to the user.")
	cobra.MarkFlagRequired(ldapCmd.Flags(), "user")
	cobra.MarkFlagRequired(ldapCmd.Flags(), "pass")
	viper.BindPFlags(ldapCmd.Flags())

	rootCmd.AddCommand(
		checkCmd,
		mtlsCmd,
		ldapCmd,
	)

	return rootCmd
}
