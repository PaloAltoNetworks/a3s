package main

import (
	"fmt"
	"strings"

	"github.com/spf13/cobra"
	"github.com/spf13/viper"
	"go.aporeto.io/a3s/cmd/a3sctl/authcmd"
	"go.aporeto.io/a3s/cmd/a3sctl/compcmd"
	"go.aporeto.io/a3s/pkgs/api"
	"go.aporeto.io/manipulate/manipcli"
)

func main() {

	cobra.OnInitialize(func() {
		viper.SetEnvPrefix("a3sctl")
		viper.AutomaticEnv()
		viper.SetEnvKeyReplacer(strings.NewReplacer("-", "_"))
	})

	rootCmd := &cobra.Command{
		Use:              "a3sctl",
		Short:            "Controls a3s APIs",
		SilenceUsage:     true,
		SilenceErrors:    true,
		TraverseChildren: true,
		PersistentPreRunE: func(cmd *cobra.Command, args []string) error {
			if err := viper.BindPFlags(cmd.PersistentFlags()); err != nil {
				return err
			}
			return viper.BindPFlags(cmd.Flags())
		},
	}

	mflags := manipcli.ManipulatorFlagSet()
	mmaker := manipcli.ManipulatorMakerFromFlags()

	apiCmd := manipcli.New(api.Manager(), mmaker)
	apiCmd.PersistentFlags().AddFlagSet(mflags)
	// cobra.MarkFlagRequired(apiCmd.PersistentFlags(), "api")
	// cobra.MarkFlagRequired(apiCmd.PersistentFlags(), "namespace")

	authCmd := authcmd.New(mmaker)
	authCmd.PersistentFlags().AddFlagSet(mflags)
	// cobra.MarkFlagRequired(authCmd.PersistentFlags(), "api")

	compCmd := compcmd.New()

	rootCmd.AddCommand(
		apiCmd,
		authCmd,
		compCmd,
	)

	if err := rootCmd.Execute(); err != nil {
		fmt.Printf("error: %s\n", err)
	}

}
