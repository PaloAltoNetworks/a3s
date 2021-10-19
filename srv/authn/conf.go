package authn

import (
	"time"
)

// Conf holds the service configuration.
type Conf struct {
	CookieDomain         string        `mapstructure:"jwt-cookie-domain"    desc:"Defines the domain for the token cookies"              required:"true"`
	CookieSameSitePolicy string        `mapstructure:"jwt-cookie-same-site" desc:"Define SameSite policy applied to token cookies"       default:"strict" allowed:"strict,lax,none"`
	JWTCertPath          string        `mapstructure:"jwt-cert"             desc:"Secret to use to sign the JWT"                         secret:"true" file:"true"`
	JWTKeyPass           string        `mapstructure:"jwt-key-pass"         desc:"JWT certificate key password"                          secret:"true" file:"true"`
	JWTKeyPath           string        `mapstructure:"jwt-key"              desc:"Path to the JWT certificate key pem file"              secret:"true" file:"true"`
	JWTMaxValidity       time.Duration `mapstructure:"jwt-max-validity"     desc:"Maximum duration of the validity of the issued tokens" default:"720h"`
}
