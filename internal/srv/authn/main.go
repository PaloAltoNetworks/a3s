package authn

import (
	"fmt"
	"io/ioutil"
	"net/http"

	"go.aporeto.io/bahamut"
	"go.aporeto.io/manipulate"
	"go.aporeto.io/tg/tglib"
	"go.uber.org/zap"
	"golang.org/x/net/context"
)

// Init initializes the module.
func Init(
	ctx context.Context,
	cfg Conf,
	server bahamut.Server,
	m manipulate.Manipulator,
	pubsub bahamut.PubSubClient,
) error {

	// JWT certificate
	certData, err := ioutil.ReadFile(cfg.JWTCertPath)
	if err != nil {
		return fmt.Errorf("unable to read jwt cert: %w", err)
	}

	keyData, err := ioutil.ReadFile(cfg.JWTKeyPath)
	if err != nil {
		return fmt.Errorf("unable to read jwt key: %w", err)
	}

	jwtCert, jwtKey, err := tglib.ReadCertificate(certData, keyData, cfg.JWTKeyPass)
	if err != nil {
		return fmt.Errorf("unable to parse jwt certificate: %w", err)
	}

	// Cookie
	cookieDomain := cfg.CookieDomain
	zap.L().Info("Cookie domain set", zap.String("domain", cookieDomain))

	var cookieSameSitePolicy http.SameSite
	switch cfg.CookieSameSitePolicy {
	case "strict":
		cookieSameSitePolicy = http.SameSiteStrictMode
	case "lax":
		cookieSameSitePolicy = http.SameSiteLaxMode
	case "none":
		cookieSameSitePolicy = http.SameSiteNoneMode
	case "default":
		cookieSameSitePolicy = http.SameSiteDefaultMode
	default:
		return fmt.Errorf("unknown SameSite policy '%s'", cfg.CookieSameSitePolicy)
	}

	zap.L().Info("Cookie policy set", zap.String("mode", cfg.CookieSameSitePolicy))

	// Let the police do the job
	_, _, _, _ = cookieSameSitePolicy, cookieDomain, jwtCert, jwtKey

	return nil
}
