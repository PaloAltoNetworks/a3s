package authcmd

import (
	"context"
	"time"

	"github.com/spf13/cobra"
	"github.com/spf13/viper"
	"go.aporeto.io/a3s/cmd/a3sctl/internal/helpers"
	"go.aporeto.io/a3s/pkgs/authlib"
	"go.aporeto.io/a3s/pkgs/permissions"
	"go.aporeto.io/manipulate/manipcli"
)

func makeLDAPCmd(mmaker manipcli.ManipulatorMaker, restrictions *permissions.Restrictions) *cobra.Command {

	cmd := &cobra.Command{
		Use:              "ldap",
		Short:            "Use a configured LDAP authentication source.",
		TraverseChildren: true,
		RunE: func(cmd *cobra.Command, args []string) error {

			flags := cmd.Flags()
			fSourceName, _ := flags.GetString("source-name")
			fSourceNamespace, _ := flags.GetString("source-namespace")
			fAudience := viper.GetStringSlice("audience")
			fUser := helpers.ReadFlag("username: ", "user", false)
			fPass := helpers.ReadFlag("password: ", "pass", true)
			fCloak := viper.GetStringSlice("cloak")
			fQRCode := viper.GetBool("qrcode")
			fValidity := viper.GetDuration("validity")

			if fSourceNamespace == "" {
				fSourceNamespace = viper.GetString("namespace")
			}

			t, err := GetLDAPToken(
				mmaker,
				fUser,
				fPass,
				fSourceNamespace,
				fSourceName,
				fAudience,
				fCloak,
				fValidity,
				restrictions,
			)
			if err != nil {
				return err
			}

			printToken(t, fQRCode)

			return nil
		},
	}

	cmd.Flags().String("user", "", "The LDAP username to use. Use '-' to prompt.")
	cmd.Flags().String("pass", "", "The password associateds to the user. Use '-' to prompt.")
	cmd.Flags().String("source-name", "default", "The name of the auth source.")
	cmd.Flags().String("source-namespace", "", "The namespace of the auth source. If omitted, uses --namespace.")

	cmd.SetHelpFunc(func(cmd *cobra.Command, args []string) {
		_ = cmd.Flags().MarkHidden("token")
		cmd.Parent().HelpFunc()(cmd, args)
	})

	return cmd
}

// GetLDAPToken retrieves a token using the
// provided LDAP source.
func GetLDAPToken(
	mmaker manipcli.ManipulatorMaker,
	user string,
	pas string,
	sourceNamespace string,
	sourceName string,
	audience []string,
	cloak []string,
	validity time.Duration,
	restrictions *permissions.Restrictions,
) (string, error) {

	m, err := mmaker()
	if err != nil {
		return "", err
	}

	opts := []authlib.Option{
		authlib.OptAudience(audience...),
		authlib.OptCloak(cloak...),
		authlib.OptValidity(validity),
	}

	if restrictions != nil {
		opts = append(opts,
			authlib.OptRestrictions(*restrictions),
		)
	}

	client := authlib.NewClient(m)
	return client.AuthFromLDAP(
		context.Background(),
		user,
		pas,
		sourceNamespace,
		sourceName,
		opts...,
	)
}
