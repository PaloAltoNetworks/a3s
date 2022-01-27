package main

import (
	"errors"
	"fmt"
	"os"
	"path"
	"strings"

	"github.com/mitchellh/go-homedir"
	"github.com/spf13/cobra"
	"github.com/spf13/viper"
	"go.aporeto.io/a3s/cmd/a3sctl/internal/authcmd"
	"go.aporeto.io/a3s/cmd/a3sctl/internal/compcmd"
	"go.aporeto.io/a3s/cmd/a3sctl/internal/flagsets"
	"go.aporeto.io/a3s/cmd/a3sctl/internal/help"
	"go.aporeto.io/a3s/pkgs/api"
	"go.aporeto.io/a3s/pkgs/bootstrap"
	"go.aporeto.io/a3s/pkgs/conf"
	"go.aporeto.io/manipulate/manipcli"
	"go.uber.org/zap"
)

var (
	cfgFile  string
	cfgName  string
	logLevel string
)

func main() {

	cobra.OnInitialize(initCobra)

	rootCmd := &cobra.Command{
		Use:              "a3sctl",
		Short:            "Controls a3s APIs",
		Long:             help.Load("a3sctl"),
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

	rootCmd.PersistentFlags().StringVar(&cfgFile, "config", "", "config file (default: $HOME/.config/a3sctl/default.yaml)")
	rootCmd.PersistentFlags().StringVar(&cfgName, "config-name", "", "default config name (default: default)")
	rootCmd.PersistentFlags().StringVar(&logLevel, "log-level", "warn", "Log level. Can be debug, info, warn or error")

	apiCmd := manipcli.New(api.Manager(), mmaker, manipcli.OptionArgumentsPrefix("with"))
	apiCmd.PersistentFlags().AddFlagSet(mflags)
	apiCmd.PersistentFlags().AddFlagSet(flagsets.MakeAutoAuthFlags())
	apiCmd.PersistentPreRunE = func(cmd *cobra.Command, args []string) error {
		if err := rootCmd.PersistentPreRunE(cmd, args); err != nil {
			return err
		}
		if err := authcmd.HandleAutoAuth(
			mmaker,
			viper.GetString("auto-auth-method"),
			nil,
			nil,
			viper.GetBool("refresh"),
			viper.GetBool("renew-cached-token"),
		); err != nil {
			return fmt.Errorf("auto auth error: %w", err)
		}
		return nil
	}

	authCmd := authcmd.New(mmaker)
	authCmd.PersistentFlags().AddFlagSet(mflags)

	compCmd := compcmd.New()

	rootCmd.AddCommand(
		apiCmd,
		authCmd,
		compCmd,
	)

	if err := rootCmd.Execute(); err != nil {
		fmt.Fprintf(os.Stderr, "error: %s\n", err)
		os.Exit(1)
	}
}

func initCobra() {

	viper.SetEnvPrefix("a3sctl")
	viper.AutomaticEnv()
	viper.SetEnvKeyReplacer(strings.NewReplacer("-", "_"))

	bootstrap.ConfigureLogger("a3sctl", conf.LoggingConf{
		LogLevel:  logLevel,
		LogFormat: "console",
	})

	home, err := homedir.Dir()
	if err != nil {
		zap.L().Fatal("unable to find home dir", zap.Error(err))
	}

	configPath := path.Join(home, ".config", "a3sctl")
	if _, err := os.Stat(configPath); os.IsNotExist(err) {
		if err := os.Mkdir(configPath, os.ModePerm); err != nil {
			zap.L().Fatal("unable to create config folder",
				zap.String("path", configPath),
				zap.Error(err),
			)
		}
	}

	if cfgFile == "" {
		cfgFile = os.Getenv("A3SCTL_CONFIG")
	}

	if cfgFile != "" {
		if _, err := os.Stat(cfgFile); os.IsNotExist(err) {
			zap.L().Fatal("config file does not exist", zap.Error(err))
		}

		viper.SetConfigType("yaml")
		viper.SetConfigFile(cfgFile)

		if err = viper.ReadInConfig(); err != nil {
			zap.L().Fatal("unable to read config",
				zap.String("path", cfgFile),
				zap.Error(err),
			)
		}

		zap.L().Debug("using config file", zap.String("path", cfgFile))
		return
	}

	viper.AddConfigPath(configPath)
	viper.AddConfigPath("/usr/local/etc/a3sctl")
	viper.AddConfigPath("/etc/a3sctl")

	if cfgName == "" {
		cfgName = os.Getenv("A3SCTL_CONFIG_NAME")
	}

	if cfgName == "" {
		cfgName = "default"
	}

	viper.SetConfigName(cfgName)

	if err = viper.ReadInConfig(); err != nil {
		if !errors.As(err, &viper.ConfigFileNotFoundError{}) {
			zap.L().Fatal("unable to read config",
				zap.Error(err),
			)
		}
	}

	zap.L().Debug("using config name", zap.String("name", cfgName))
}
