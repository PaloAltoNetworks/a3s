package authcmd

import (
	"fmt"
	"os"

	"github.com/mdp/qrterminal"
	"github.com/spf13/cobra"
	"go.aporeto.io/a3s/cmd/a3sctl/internal/help"
	"go.aporeto.io/a3s/pkgs/permissions"
	"go.aporeto.io/manipulate/manipcli"
)

// New returns a new auth command.
func New(mmaker manipcli.ManipulatorMaker) *cobra.Command {

	// rootCmd is the root cmd for authentication subcommand.
	rootCmd := &cobra.Command{
		Use:              "auth",
		Short:            "Authenticate and retrieve a token.",
		Long:             help.Load("auth"),
		TraverseChildren: true,
	}
	rootCmd.PersistentFlags().Duration("validity", 0, "The validity of the requested token.")
	rootCmd.PersistentFlags().StringSlice("audience", nil, "Requested audience for the token.")
	rootCmd.PersistentFlags().String("source-name", "default", "The name of the auth source.")
	rootCmd.PersistentFlags().String("source-namespace", "", "The namespace of the auth source. If omitted, --namespace will be used.")
	rootCmd.PersistentFlags().StringSlice("cloak", nil, "Cloak identity claims. Only claims with a prefix matching of of the given string will be used in the token.")
	rootCmd.PersistentFlags().Bool("qrcode", false, "If passed, display the token as a QR code.")

	// Freaking pglags and its non configurable split char
	// and missing GetStringArray...
	restrictions := &permissions.Restrictions{}
	rootCmd.PersistentFlags().StringArrayVar(&restrictions.Permissions, "restrict-permissions", nil, "Restrict the permissions to what is given.")
	rootCmd.PersistentFlags().StringArrayVar(&restrictions.Networks, "restrict-network", nil, "Rrestrict the origin networks from which the token can be used.")
	rootCmd.PersistentFlags().StringVar(&restrictions.Namespace, "restrict-namespace", "", "Rrestrict the namespace from which the token can be used.")

	rootCmd.AddCommand(
		makeCheckCmd(mmaker),
		makePermsCmd(mmaker),
		makeAutoCmd(mmaker),
		makeMTLSCmd(mmaker, restrictions),
		makeLDAPCmd(mmaker, restrictions),
		makeAzureCmd(mmaker, restrictions),
		makeGCPCmd(mmaker, restrictions),
		makeAWSCmd(mmaker, restrictions),
		makeOIDCCmd(mmaker, restrictions),
		makeRemoteA3SCmd(mmaker, restrictions),
		makeA3SCmd(mmaker, restrictions),
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
