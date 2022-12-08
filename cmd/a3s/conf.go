package main

import (
	"crypto"
	"crypto/x509"
	"fmt"
	"strings"
	"time"

	"go.aporeto.io/a3s/pkgs/authenticator"
	"go.aporeto.io/a3s/pkgs/conf"
	"go.aporeto.io/a3s/pkgs/lombric"
	"go.aporeto.io/tg/tglib"
)

var (
	version = "v0.0.0"
	commit  = "dev"
)

// Conf holds the main configuration flags.
type Conf struct {
	Init               bool   `mapstructure:"init"              desc:"If set, initialize the root permissions using the CAs passed in --init-root-ca and --init-platform-ca"`
	InitContinue       bool   `mapstructure:"init-continue"     desc:"Continues normal boot after init."`
	InitDB             bool   `mapstructure:"init-db"           desc:"If set, initialize the database using the mongo config passed in and init-db-username"`
	InitDBUsername     string `mapstructure:"init-db-username"  desc:"If init-db is set, this will define the username to use on db initialization"           default:"CN=a3s,OU=root,O=system"`
	InitPlatformCAPath string `mapstructure:"init-platform-ca"  desc:"Path to the platform CA to use to initialize platform permissions"`
	InitRootUserCAPath string `mapstructure:"init-root-ca"      desc:"Path to the root CA to use to initialize root permissions"`
	InitData           string `mapstructure:"init-data"         desc:"Path to an import file containing initial provisionning data"`

	JWT        JWTConf        `mapstructure:",squash"`
	MTLSHeader MTLSHeaderConf `mapstructure:",squash"`

	conf.APIServerConf       `mapstructure:",squash"`
	conf.GatewayConf         `mapstructure:",squash"`
	conf.HTTPTimeoutsConf    `mapstructure:",squash"`
	conf.HealthConfiguration `mapstructure:",squash"`
	conf.LoggingConf         `mapstructure:",squash"`
	conf.MongoConf           `mapstructure:",squash" override:"mongo-db=a3s"`
	conf.NATSPublisherConf   `mapstructure:",squash"`
	conf.ProfilingConf       `mapstructure:",squash"`
	conf.RateLimitingConf    `mapstructure:",squash"`
}

// Prefix returns the configuration prefix.
func (c *Conf) Prefix() string { return "a3s" }

// PrintVersion prints the current version.
func (c *Conf) PrintVersion() {
	fmt.Printf("a3s %s (%s)", version, commit)
}

func newConf() Conf {
	c := Conf{}
	lombric.Initialize(&c)
	return c
}

// JWTConf holds the configuration related to jwt management.
type JWTConf struct {
	JWTAudience        string        `mapstructure:"jwt-audience"         desc:"Default audience for delivered jwt"`
	JWTCertPath        string        `mapstructure:"jwt-cert"             desc:"Secret to use to sign the JWT"                         secret:"true" file:"true"`
	JWTCookieDomain    string        `mapstructure:"jwt-cookie-domain"    desc:"Defines the domain for the cookie"`
	JWTCookiePolicy    string        `mapstructure:"jwt-cookie-policy"    desc:"Define same site policy applied to token cookies"      default:"strict" allowed:"strict,lax,none"`
	JWTDefaultValidity time.Duration `mapstructure:"jwt-default-validity" desc:"Default duration of the validity of the issued tokens" default:"24h"`
	JWTIssuer          string        `mapstructure:"jwt-issuer"           desc:"Value used for issuer jwt field"`
	JWTKeyPass         string        `mapstructure:"jwt-key-pass"         desc:"JWT certificate key password"                          secret:"true" file:"true"`
	JWTKeyPath         string        `mapstructure:"jwt-key"              desc:"Path to the JWT certificate key pem file"              secret:"true" file:"true"`
	JWTMaxValidity     time.Duration `mapstructure:"jwt-max-validity"     desc:"Maximum duration of the validity of the issued tokens" default:"720h"`
	JWTTrustedIssuers  []string      `mapstructure:"jwt-trusted-issuer"   desc:"List of externally trusted issuers"`

	jwtCert *x509.Certificate
	jwtKey  crypto.PrivateKey
}

// JWTCertificate returns the certificate used to verify JWTs.
func (c *JWTConf) JWTCertificate() (*x509.Certificate, crypto.PrivateKey, error) {

	if c.jwtCert != nil {
		return c.jwtCert, c.jwtKey, nil
	}

	jwtCert, jwtKey, err := tglib.ReadCertificatePEM(c.JWTCertPath, c.JWTKeyPath, c.JWTKeyPass)
	if err != nil {
		return nil, nil, fmt.Errorf("unable to read jwt certificate: %w", err)
	}

	c.jwtCert = jwtCert
	c.jwtKey = jwtKey

	return jwtCert, jwtKey, nil
}

// TrustedIssuers parses --jwt-trusted-issuers and returns a list
// of prepopulated authenticator.RemoteIssuer.
func (c *JWTConf) TrustedIssuers() ([]authenticator.RemoteIssuer, error) {

	if len(c.JWTTrustedIssuers) == 0 {
		return nil, nil
	}

	out := make([]authenticator.RemoteIssuer, len(c.JWTTrustedIssuers))

	for i, r := range c.JWTTrustedIssuers {

		ri := authenticator.RemoteIssuer{}
		parts := strings.SplitN(r, "@", 2)

		ri.URL = parts[0]

		if len(parts) == 2 {
			cert, err := tglib.ParseCertificatePEM(parts[1])
			if err != nil {
				return nil, fmt.Errorf("unable to parse cert at '%s': %w", parts[1], err)
			}

			ri.Pool = x509.NewCertPool()
			ri.Pool.AddCert(cert)
		}

		out[i] = ri
	}

	return out, nil
}

// MTLSHeaderConf holds the configuration for trusted certificate header.
type MTLSHeaderConf struct {
	Enabled    bool   `mapstructure:"mtls-header-enabled"    desc:"Trust the value of the defined header containing a user certificate. This is insecure if there is no proper tls verification happening upstream"`
	HeaderKey  string `mapstructure:"mtls-header-key"        desc:"The header to check for user certificates" default:"x-tls-certificate"`
	Passphrase string `mapstructure:"mtls-header-passphrase" desc:"The passphrase to decrypt the AES encrypted header content. It is mandatory if --mtls-header-enabled is set."`
}
