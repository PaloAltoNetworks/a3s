package authcmd

import (
	"crypto/sha256"
	"fmt"
	"os"
	"path"
	"time"

	"github.com/golang-jwt/jwt/v4"
	"github.com/mitchellh/go-homedir"
	"github.com/spf13/viper"
	"go.aporeto.io/a3s/cmd/a3sctl/helpers"
	"go.aporeto.io/a3s/pkgs/token"
	"go.aporeto.io/manipulate/manipcli"
)

// HandleAutoAuth handles automatic retrieval of tokens based on
// the current config file.
// If will check for `autoauth.enable` to retrieve desired auto auth
// method. Setting it to empty will disable auto auth.
// Support:
//
// autoauth.enable: mtls
//      autoauth.mtls.cert: path to the client certificate
//      autoauth.mtls.key: path to the client certificate key
//      autoauth.mtls.keyPass: optional passphrase to the certificate.
//      autoauth.mtls.source.name: the name of the MTLS source to use.
//      autoauth.mtls.source.namespace: the namespace of the MTLS source to use.
func HandleAutoAuth(mmaker manipcli.ManipulatorMaker) error {

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

	method := viper.GetString("autoauth.enable")

	tokenCache := path.Join(cache, fmt.Sprintf("token-%s-%x", method, sha256.Sum256([]byte(viper.GetString("api")))))
	data, err := os.ReadFile(tokenCache)

	if err != nil {
		if !os.IsNotExist(err) {
			return err
		}

		switch method {
		case "mtls", "MTLS":
			t, err := GetMTLSToken(
				mmaker,
				viper.GetString("autoauth.mtls.cert"),
				viper.GetString("autoauth.mtls.key"),
				helpers.ReadFlag("passphrase: ", "autoauth.mtls.pass", true),
				viper.GetString("autoauth.mtls.source.namespace"),
				viper.GetString("autoauth.mtls.source.name"),
				viper.GetStringSlice("autoauth.mtls.audience"),
				viper.GetStringSlice("autoauth.mtls.cloak"),
				24*time.Hour,
				nil,
			)
			if err != nil {
				return fmt.Errorf("unable to retrieve token from autoauth info: %w", err)
			}
			data = []byte(t)

		case "ldap", "LDAP":
			t, err := GetLDAPToken(
				mmaker,
				helpers.ReadFlag("username: ", "autoauth.ldap.user", false),
				helpers.ReadFlag("password: ", "autoauth.ldap.pass", true),
				viper.GetString("autoauth.ldap.source.namespace"),
				viper.GetString("autoauth.ldap.source.name"),
				viper.GetStringSlice("autoauth.ldap.audience"),
				viper.GetStringSlice("autoauth.ldap.cloak"),
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
		return HandleAutoAuth(mmaker)
	}

	viper.Set("token", string(data))

	return nil
}
