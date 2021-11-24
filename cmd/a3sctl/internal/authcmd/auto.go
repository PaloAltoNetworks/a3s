package authcmd

import (
	"crypto/sha256"
	"fmt"
	"os"
	"path"
	"time"

	"github.com/golang-jwt/jwt/v4"
	"github.com/mitchellh/go-homedir"
	"github.com/spf13/cobra"
	"github.com/spf13/viper"
	"go.aporeto.io/a3s/cmd/a3sctl/internal/helpers"
	"go.aporeto.io/a3s/pkgs/token"
	"go.aporeto.io/manipulate/manipcli"
	"go.uber.org/zap"
)

func makeAutoCmd(mmaker manipcli.ManipulatorMaker) *cobra.Command {

	cmd := &cobra.Command{
		Use:              "auto",
		Short:            "Use config file auto auth.",
		TraverseChildren: true,
		RunE: func(cmd *cobra.Command, args []string) error {

			if viper.ConfigFileUsed() == "" {
				return fmt.Errorf("auto subcommand is available only when using a config file")
			}

			return HandleAutoAuth(
				mmaker,
				viper.GetString("auto-auth-method"),
				viper.GetStringSlice("override-audience"),
				viper.GetStringSlice("override-cloak"),
				true,
			)
		},
	}

	cmd.Flags().StringSlice("override-audience", nil, "Override audience")
	cmd.Flags().StringSlice("override-cloak", nil, "Override cloak")

	return cmd
}

// HandleAutoAuth handles automatic retrieval of tokens based on
// the current config file.
// If will check for `autoauth.enable` to retrieve desired auto auth
// method. Setting it to empty will disable auto auth.
// Support:
//
// autoauth.enable: mtls
//      autoauth.mtls.cert: path to the client certificate
//      autoauth.mtls.key: path to the client certificate key
//      autoauth.mtls.pass: optional passphrase to the certificate.
//      autoauth.mtls.source.name: the name of the MTLS source to use.
//      autoauth.mtls.source.namespace: the namespace of the MTLS source to use.
// autoauth.enable: ldap
//      autoauth.ldap.user: the username.
//      autoauth.ldap.pass: the password.
//      autoauth.ldap.source.name: the name of the LDAP source to use.
//      autoauth.ldap.source.namespace: the namespace of the LDAP source to use.
func HandleAutoAuth(mmaker manipcli.ManipulatorMaker, method string, overrideAudience []string, overrideCloak []string, refresh bool) error {

	if viper.GetString("token") != "" {
		zap.L().Debug("autoauth: using --token")
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

	if method == "" {
		method = viper.GetString("autoauth.enable")
	}

	tokenCache := path.Join(cache, fmt.Sprintf("token-%s-%x", method, sha256.Sum256([]byte(viper.GetString("api")))))

	if refresh {
		if err := os.Remove(tokenCache); err != nil && !os.IsNotExist(err) {
			return err
		}
	}

	overrideIfNeeded := func(key string, override []string) []string {
		if override != nil {
			return override
		}
		return viper.GetStringSlice(key)
	}

	data, err := os.ReadFile(tokenCache)
	if err != nil {
		if !os.IsNotExist(err) {
			return err
		}

		switch method {
		case "mtls", "MTLS":
			zap.L().Debug("autoauth: retrieving token using autoauth.mtls")
			t, err := GetMTLSToken(
				mmaker,
				viper.GetString("autoauth.mtls.cert"),
				viper.GetString("autoauth.mtls.key"),
				helpers.ReadFlag("passphrase: ", "autoauth.mtls.pass", true),
				viper.GetString("autoauth.mtls.source.namespace"),
				viper.GetString("autoauth.mtls.source.name"),
				overrideIfNeeded("autoauth.mtls.audience", overrideAudience),
				overrideIfNeeded("autoauth.mtls.cloak", overrideCloak),
				24*time.Hour,
				nil,
			)
			if err != nil {
				return fmt.Errorf("unable to retrieve token from autoauth info: %w", err)
			}
			data = []byte(t)

		case "ldap", "LDAP":
			zap.L().Debug("autoauth: retrieving token using autoauth.ldap")
			t, err := GetLDAPToken(
				mmaker,
				helpers.ReadFlag("username: ", "autoauth.ldap.user", false),
				helpers.ReadFlag("password: ", "autoauth.ldap.pass", true),
				viper.GetString("autoauth.ldap.source.namespace"),
				viper.GetString("autoauth.ldap.source.name"),
				overrideIfNeeded("autoauth.ldap.audience", overrideAudience),
				overrideIfNeeded("autoauth.ldap.cloak", overrideCloak),
				24*time.Hour,
				nil,
			)
			if err != nil {
				return fmt.Errorf("unable to retrieve token from autoauth info: %w", err)
			}
			data = []byte(t)

		case "":
			return nil

		default:
			return fmt.Errorf("unsupported auto auth method: %s", method)
		}

		if err := os.WriteFile(tokenCache, data, 0600); err != nil {
			return fmt.Errorf("unable to write token cache: %w", err)
		}
		zap.L().Debug("autoauth: token cached", zap.String("path", tokenCache))
	}

	idt := &token.IdentityToken{}
	p := &jwt.Parser{}
	if _, _, err := p.ParseUnverified(string(data), idt); err != nil {
		return fmt.Errorf("unable to parse cached token: %w", err)
	}

	if time.Until(idt.ExpiresAt.Time) <= time.Duration(idt.ExpiresAt.Unix()/2) {
		zap.L().Debug("autoauth: token about to expire. removing", zap.String("path", tokenCache))
		if err := os.Remove(tokenCache); err != nil {
			return fmt.Errorf("unable to clean currently cached token: %w", err)
		}
		return HandleAutoAuth(mmaker, method, overrideAudience, overrideCloak, false)
	}

	zap.L().Debug("autoauth: token set from cache", zap.String("path", tokenCache))
	viper.Set("token", string(data))

	return nil
}
