package authcmd

import (
	"fmt"

	"github.com/ghodss/yaml"
	"github.com/golang-jwt/jwt/v4"
	"github.com/spf13/cobra"
	"github.com/spf13/viper"
	"go.aporeto.io/a3s/pkgs/api"
	"go.aporeto.io/a3s/pkgs/token"
	"go.aporeto.io/manipulate/manipcli"
)

func makePermsCmd(mmaker manipcli.ManipulatorMaker) *cobra.Command {

	cmd := &cobra.Command{
		Use:              "permissions",
		Aliases:          []string{"perms", "perm"},
		Short:            "Check the permissions",
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

			fToken := viper.GetString("access-token")
			fNamespace := viper.GetString("namespace")

			if fToken == "" {
				fToken = viper.GetString("token")
			}

			idt := &token.IdentityToken{}
			p := jwt.Parser{}
			if _, _, err := p.ParseUnverified(fToken, idt); err != nil {
				return err
			}

			m, err := mmaker()
			if err != nil {
				return err
			}

			perms := api.NewPermissions()
			perms.Claims = idt.Identity
			perms.Namespace = fNamespace

			if r := idt.Restrictions; r != nil {
				perms.RestrictedNamespace = r.Namespace
				perms.RestrictedNetworks = r.Networks
				perms.RestrictedPermissions = r.Permissions
			}

			if err := m.Create(nil, perms); err != nil {
				return err
			}

			data, err := yaml.Marshal(perms.Permissions)
			if err != nil {
				return err
			}

			fmt.Println(string(data))

			return nil
		},
	}

	cmd.Flags().String("access-token", "", "Valid remote a3s token. If omitted, uses --token.")

	cmd.SetHelpFunc(func(cmd *cobra.Command, args []string) {
		_ = cmd.Flags().MarkHidden("audience")
		_ = cmd.Flags().MarkHidden("cloak")
		_ = cmd.Flags().MarkHidden("validity")
		_ = cmd.Flags().MarkHidden("restrict-namespace")
		_ = cmd.Flags().MarkHidden("restrict-permissions")
		_ = cmd.Flags().MarkHidden("restrict-network")
		cmd.Parent().HelpFunc()(cmd, args)
	})

	return cmd
}
