package main

import (
	"crypto/sha256"
	"fmt"
	"os"
	"path"
	"strings"
	"time"

	"github.com/golang-jwt/jwt/v4"
	"github.com/mitchellh/go-homedir"
	"github.com/spf13/cobra"
	"github.com/spf13/viper"
	"go.aporeto.io/a3s/cmd/a3sctl/authcmd"
	"go.aporeto.io/a3s/cmd/a3sctl/compcmd"
	"go.aporeto.io/a3s/internal/conf"
	"go.aporeto.io/a3s/pkgs/api"
	"go.aporeto.io/a3s/pkgs/bootstrap"
	"go.aporeto.io/a3s/pkgs/token"
	"go.aporeto.io/manipulate/manipcli"
)

var (
	cfgFile  string
	logLevel string
)

func main() {

	cobra.OnInitialize(initCobra)

	mflags := manipcli.ManipulatorFlagSet()
	mmaker := manipcli.ManipulatorMakerFromFlags()

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
			if err := viper.BindPFlags(cmd.Flags()); err != nil {
				return err
			}
			if err := handleAutoAuth(mmaker); err != nil {
				return fmt.Errorf("unable to handle autoauth: %w", err)
			}
			return nil
		},
	}
	rootCmd.PersistentFlags().StringVar(&cfgFile, "config", "", "config file (default: $HOME/.config/a3sctl/default.yaml)")
	rootCmd.PersistentFlags().StringVar(&logLevel, "log-level", "warn", "Log level. Can be debug, info, warn or error")

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
		fmt.Fprintf(os.Stderr, "error: %s\n", err)
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

func handleAutoAuth(mmaker manipcli.ManipulatorMaker) error {

	if viper.GetString("token") != "" {
		return nil
	}

	home, err := homedir.Dir()
	if err != nil {
		return fmt.Errorf("unable to find home dir: %w", err)
	}

	cache, ok := os.LookupEnv("XDG_HOME_CACHE")
	if ok {
		cache = path.Join(home, cache, "a3scl")
	} else {
		cache = path.Join(home, ".cache", "a3sctl")
	}

	if _, err := os.Stat(cache); os.IsNotExist(err) {
		if err := os.Mkdir(cache, 0700); err != nil {
			return fmt.Errorf("failed to create cache %s: %w", cache, err)
		}
	}

	tokenCache := path.Join(cache, fmt.Sprintf("token-%x", sha256.Sum256([]byte(viper.GetString("api")))))
	data, err := os.ReadFile(tokenCache)
	if err != nil {
		if !os.IsNotExist(err) {
			return err
		}

		t, err := authcmd.GetMTLSToken(
			mmaker,
			viper.GetString("autoauth.cert"),
			viper.GetString("autoauth.key"),
			viper.GetString("autoauth.keyPass"),
			viper.GetString("autoauth.source.namespace"),
			viper.GetString("autoauth.source.name"),
			viper.GetStringSlice("autoauth.audience"),
			viper.GetStringSlice("autoauth.cloak"),
			24*time.Hour,
			nil,
		)
		if err != nil {
			return fmt.Errorf("unable to retrieve token from autoauth info: %w", err)
		}

		data = []byte(t)

		if err := os.WriteFile(tokenCache, data, 0600); err != nil {
			return fmt.Errorf("unable to write token cache: %w", err)
		}
	}

	idt := &token.IdentityToken{}
	p := &jwt.Parser{}
	if _, _, err := p.ParseUnverified(string(data), idt); err != nil {
		return fmt.Errorf("unable to parse cached token: %w", err)
	}

	if time.Until(idt.ExpiresAt.Time) <= time.Duration(idt.ExpiresAt.Unix()/2) {
		if err := os.Remove(tokenCache); err != nil {
			return fmt.Errorf("unable to clean currently cached token: %w", err)
		}
		return handleAutoAuth(mmaker)
	}

	viper.Set("token", string(data))

	return nil
}
