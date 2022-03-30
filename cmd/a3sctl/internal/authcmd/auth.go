package authcmd

import (
	"github.com/spf13/cobra"
	"go.aporeto.io/a3s/cmd/a3sctl/internal/help"
	"go.aporeto.io/a3s/pkgs/permissions"
	"go.aporeto.io/manipulate/manipcli"
)

// New returns a new auth command.
func New(mmaker manipcli.ManipulatorMaker) *cobra.Command {

	// cmd is the root cmd for authentication subcommand.
	cmd := &cobra.Command{
		Use:              "auth",
		Short:            "Authenticate and retrieve a token.",
		Long:             help.Load("auth"),
		TraverseChildren: true,
	}
	cmd.PersistentFlags().Duration("validity", 0, "The validity of the requested token.")
	cmd.PersistentFlags().StringSlice("audience", nil, "Requested audience for the token.")
	cmd.PersistentFlags().StringSlice("cloak", nil, "Cloak identity claims. Only claims with a prefix matching of of the given string will be used in the token.")
	cmd.PersistentFlags().Bool("qrcode", false, "If passed, display the token as a QR code.")
	cmd.PersistentFlags().Bool("check", false, "If passed, display the decoded token")
	cmd.PersistentFlags().Bool("refresh", false, "If set, ask for a refresh token.")

	// Freaking pglags and its non configurable split char
	// and missing GetStringArray...
	restrictions := &permissions.Restrictions{}
	cmd.PersistentFlags().StringArrayVar(&restrictions.Permissions, "restrict-permissions", nil, "Restrict the permissions to what is given.")
	cmd.PersistentFlags().StringArrayVar(&restrictions.Networks, "restrict-network", nil, "Rrestrict the origin networks from which the token can be used.")
	cmd.PersistentFlags().StringVar(&restrictions.Namespace, "restrict-namespace", "", "Rrestrict the namespace from which the token can be used.")

	cmd.AddCommand(
		makeCheckCmd(mmaker),
		makePermsCmd(mmaker),
		makeAutoCmd(mmaker),
		makeMTLSCmd(mmaker, restrictions),
		makeLDAPCmd(mmaker, restrictions),
		makeHTTPCmd(mmaker, restrictions),
		makeAzureCmd(mmaker, restrictions),
		makeGCPCmd(mmaker, restrictions),
		makeAWSCmd(mmaker, restrictions),
		makeOIDCCmd(mmaker, restrictions),
		makeRemoteA3SCmd(mmaker, restrictions),
		makeA3SCmd(mmaker, restrictions),
	)

	return cmd
}
