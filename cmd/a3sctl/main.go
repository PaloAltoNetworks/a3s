package main

import (
	"fmt"
	"strings"

	"github.com/spf13/cobra"
	"github.com/spf13/viper"
	"go.aporeto.io/a3s/cmd/a3sctl/authcmd"
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
		Use:           "a3sctl",
		Short:         "Controls a3s APIs",
		SilenceUsage:  true,
		SilenceErrors: true,
		PersistentPreRunE: func(cmd *cobra.Command, args []string) error {
			return viper.BindPFlags(cmd.PersistentFlags())
		},
	}

	rootCmd.PersistentFlags().AddFlagSet(manipcli.ManipulatorFlagSet())
	m := manipcli.ManipulatorMakerFromFlags()
	apiCmd := manipcli.GenerateCobraCommand(api.Manager(), m)

	rootCmd.AddCommand(
		apiCmd,
		authcmd.AuthCmd,
		completionCmd,
	)

	if err := rootCmd.Execute(); err != nil {
		fmt.Printf("error: %s", err)
	}

}
