package main

import (
	"crypto"
	"crypto/x509"
	"fmt"
	"time"

	"go.aporeto.io/a3s/pkgs/conf"
	"go.aporeto.io/addedeffect/lombric"
	"go.aporeto.io/tg/tglib"
)

// Conf holds the main configuration flags.
type Conf struct {
	Init               bool   `mapstructure:"init"          desc:"If set, initialize the root permissions usingf the CA passed in --init-root-ca and exit"`
	InitContinue       bool   `mapstructure:"init-continue" desc:"Continues normal boot after init."`
	InitRootUserCAPath string `mapstructure:"init-root-ca"  desc:"Path to the root CA to use to initialize root permissions"`

	JWT JWTConf `mapstructure:",squash"`

	conf.APIServerConf       `mapstructure:",squash"`
	conf.HealthConfiguration `mapstructure:",squash"`
	conf.HTTPTimeoutsConf    `mapstructure:",squash"`
	conf.LoggingConf         `mapstructure:",squash"`
	conf.NATSPublisherConf   `mapstructure:",squash"`
	conf.ProfilingConf       `mapstructure:",squash"`
	conf.RateLimitingConf    `mapstructure:",squash"`
	conf.MongoConf           `mapstructure:",squash" override:"mongo-db=a3s"`
}

// Prefix returns the configuration prefix.
func (c *Conf) Prefix() string { return "a3s" }

// PrintVersion prints the current version.
func (c *Conf) PrintVersion() {
	fmt.Printf("authn 0.0.1")
}

func newConf() Conf {
	c := Conf{}
	lombric.Initialize(&c)
	return c
}

// JWTConf holds the configuration related to jwt management.
type JWTConf struct {
	JWTAudience    string        `mapstructure:"jwt-audience"         desc:"Default audience for delivered jwt"`
	JWTCertPath    string        `mapstructure:"jwt-cert"             desc:"Secret to use to sign the JWT"                         secret:"true" file:"true"`
	JWTIssuer      string        `mapstructure:"jwt-issuer"           desc:"Value used for issuer jwt field"`
	JWTKeyPass     string        `mapstructure:"jwt-key-pass"         desc:"JWT certificate key password"                          secret:"true" file:"true"`
	JWTKeyPath     string        `mapstructure:"jwt-key"              desc:"Path to the JWT certificate key pem file"              secret:"true" file:"true"`
	JWTMaxValidity time.Duration `mapstructure:"jwt-max-validity"     desc:"Maximum duration of the validity of the issued tokens" default:"720h"`

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
