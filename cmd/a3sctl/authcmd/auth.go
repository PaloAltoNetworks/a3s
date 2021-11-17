package authcmd

import (
	"fmt"
	"os"

	"github.com/golang-jwt/jwt/v4"
	"github.com/hokaccha/go-prettyjson"
	"github.com/mdp/qrterminal"
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
	rootCmd.PersistentFlags().StringSlice("cloak", nil, "Cloak identity claims. Only claims with a prefix matching of of the given string will be used in the token.")
	rootCmd.PersistentFlags().Bool("qrcode", false, "If passed, display the token as a QR code.")

	checkCmd := &cobra.Command{
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
	checkCmd.Flags().String("token", "", "The token to verify.")
	_ = cobra.MarkFlagRequired(checkCmd.Flags(), "token")

	mtlsCmd := &cobra.Command{
		Use:              "mtls",
		Short:            "Use a configured MTLS authentication source.",
		TraverseChildren: true,
		RunE: func(cmd *cobra.Command, args []string) error {

			fCert := viper.GetString("cert")
			fKey := viper.GetString("key")
			fPass := viper.GetString("pass")
			fSourceName := viper.GetString("source-name")
			fSourceNamespace := viper.GetString("source-namespace")
			fAudience := viper.GetStringSlice("audience")
			fCloak := viper.GetStringSlice("cloak")
			fQRCode := viper.GetBool("qrcode")

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
			iss.Cloak = fCloak

			if err := m.Create(nil, iss); err != nil {
				return err
			}

			printToken(iss.Token, fQRCode)

			return nil
		},
	}
	mtlsCmd.Flags().String("cert", "", "Path to the certificate in PEM format.")
	mtlsCmd.Flags().String("key", "", "Path to the certificate key in PEM format.")
	mtlsCmd.Flags().String("pass", "", "Passphrase for the certificate key.")
	_ = cobra.MarkFlagRequired(mtlsCmd.Flags(), "cert")
	_ = cobra.MarkFlagRequired(mtlsCmd.Flags(), "key")

	azureCmd := &cobra.Command{
		Use:              "azure",
		Short:            "Use an Azure identity token.",
		TraverseChildren: true,
		RunE: func(cmd *cobra.Command, args []string) error {

			fToken := viper.GetString("token")
			fAudience := viper.GetStringSlice("audience")
			fCloak := viper.GetStringSlice("cloak")
			fQRCode := viper.GetBool("qrcode")

			iss := api.NewIssue()
			iss.SourceType = api.IssueSourceTypeAzure
			iss.Audience = fAudience
			iss.Cloak = fCloak
			iss.InputAzure = &api.IssueAzure{
				Token: fToken,
			}

			m, err := mmaker()
			if err != nil {
				return err
			}

			if err := m.Create(nil, iss); err != nil {
				return err
			}

			printToken(iss.Token, fQRCode)

			return nil
		},
	}
	azureCmd.Flags().String("token", "", "Valid Azure token.")

	gcpCmd := &cobra.Command{
		Use:              "gcp",
		Short:            "Use an GCP identity token.",
		TraverseChildren: true,
		RunE: func(cmd *cobra.Command, args []string) error {

			fToken := viper.GetString("token")
			fTokenAudience := viper.GetString("token-audience")
			fAudience := viper.GetStringSlice("audience")
			fCloak := viper.GetStringSlice("cloak")
			fQRCode := viper.GetBool("qrcode")

			iss := api.NewIssue()
			iss.SourceType = api.IssueSourceTypeGCP
			iss.Audience = fAudience
			iss.Cloak = fCloak
			iss.InputGCP = &api.IssueGCP{
				Token:    fToken,
				Audience: fTokenAudience,
			}

			m, err := mmaker()
			if err != nil {
				return err
			}

			if err := m.Create(nil, iss); err != nil {
				return err
			}

			printToken(iss.Token, fQRCode)

			return nil
		},
	}
	gcpCmd.Flags().String("token", "", "Valid Azure token.")
	gcpCmd.Flags().String("token-audience", "", "Required GCP token audience.")

	ldapCmd := &cobra.Command{
		Use:              "ldap",
		Short:            "Use a configured LDAP authentication source.",
		TraverseChildren: true,
		RunE: func(cmd *cobra.Command, args []string) error {

			flags := cmd.Flags()
			fSourceName, _ := flags.GetString("source-name")
			fSourceNamespace, _ := flags.GetString("source-namespace")
			fAudience := viper.GetStringSlice("audience")
			fUser := viper.GetString("user")
			fPass := viper.GetString("pass")
			fCloak := viper.GetStringSlice("cloak")
			fQRCode := viper.GetBool("qrcode")

			iss := api.NewIssue()
			iss.SourceType = api.IssueSourceTypeLDAP
			iss.SourceName = fSourceName
			iss.SourceNamespace = fSourceNamespace
			iss.Audience = fAudience
			iss.Cloak = fCloak
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

			printToken(iss.Token, fQRCode)

			return nil
		},
	}
	ldapCmd.Flags().String("user", "", "The LDAP username to use.")
	ldapCmd.Flags().String("pass", "", "The password associateds to the user.")
	_ = cobra.MarkFlagRequired(ldapCmd.Flags(), "user")
	_ = cobra.MarkFlagRequired(ldapCmd.Flags(), "pass")

	rootCmd.AddCommand(
		checkCmd,
		mtlsCmd,
		ldapCmd,
		azureCmd,
		gcpCmd,
	)

	return rootCmd
}

func printToken(token string, qrCode bool) {

	if !qrCode {
		fmt.Println(token)
		return
	}

	qrterminal.GenerateWithConfig(
		token,
		qrterminal.Config{
			Writer:         os.Stdout,
			Level:          qrterminal.M,
			HalfBlocks:     true,
			QuietZone:      1,
			BlackChar:      qrterminal.BLACK_BLACK,
			WhiteBlackChar: qrterminal.WHITE_BLACK,
			WhiteChar:      qrterminal.WHITE_WHITE,
			BlackWhiteChar: qrterminal.BLACK_WHITE,
		},
	)
}
