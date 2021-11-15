package main

import (
	"fmt"
	"os"
	"path"
	"strings"

	"github.com/mitchellh/go-homedir"
	"github.com/spf13/cobra"
	"github.com/spf13/viper"
	"go.aporeto.io/a3s/cmd/a3sctl/authcmd"
	"go.aporeto.io/a3s/cmd/a3sctl/compcmd"
	"go.aporeto.io/a3s/pkgs/api"
	"go.aporeto.io/manipulate/manipcli"
)

var cfgFile string

func main() {

	cobra.OnInitialize(initCobra)

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
	rootCmd.PersistentFlags().StringVar(&cfgFile, "config", "", "config file (default: $HOME/.config/a3sctl/default.yaml)")

	mflags := manipcli.ManipulatorFlagSet()
	mmaker := manipcli.ManipulatorMakerFromFlags()

	apiCmd := manipcli.New(api.Manager(), mmaker)
	apiCmd.PersistentFlags().AddFlagSet(mflags)

	authCmd := authcmd.New(mmaker)
	authCmd.PersistentFlags().AddFlagSet(mflags)

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

func initCobra() {

	viper.SetEnvPrefix("a3sctl")
	viper.AutomaticEnv()
	viper.SetEnvKeyReplacer(strings.NewReplacer("-", "_"))

	home, err := homedir.Dir()
	if err != nil {
		fmt.Println("error: unable to find home dir:", err)
		return
	}

	hpath := path.Join(home, ".config", "a3sctl")
	if _, err := os.Stat(hpath); os.IsNotExist(err) {
		if err := os.Mkdir(hpath, os.ModePerm); err != nil {
			fmt.Printf("error: failed to create %s: %s\n", hpath, err)
			return
		}
	}

	if cfgFile == "" {
		cfgFile = os.Getenv("A3SCTL_CONFIG")
	}

	if cfgFile != "" {

		if _, err := os.Stat(cfgFile); os.IsNotExist(err) {
			fmt.Println("error: config file does not exist:", err)
			os.Exit(1)
		}

		viper.SetConfigType("yaml")
		viper.SetConfigFile(cfgFile)
		_ = viper.ReadInConfig()

		return
	}

	viper.AddConfigPath(hpath)
	viper.AddConfigPath("/usr/local/etc/a3sctl")
	viper.AddConfigPath("/etc/a3sctl")

	if cfgName := os.Getenv("A3SCTL_CONFIG_NAME"); cfgName != "" {
		viper.SetConfigName(cfgName)
	} else {
		viper.SetConfigName("default")
	}

	_ = viper.ReadInConfig()
}
