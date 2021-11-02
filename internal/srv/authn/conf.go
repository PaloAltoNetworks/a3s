package authn

import (
	"crypto"
	"crypto/x509"
	"fmt"
	"time"

	"go.aporeto.io/tg/tglib"
)

// Conf holds the service configuration.
type Conf struct {
	JWTCertPath    string        `mapstructure:"jwt-cert"             desc:"Secret to use to sign the JWT"                         secret:"true" file:"true"`
	JWTKeyPass     string        `mapstructure:"jwt-key-pass"         desc:"JWT certificate key password"                          secret:"true" file:"true"`
	JWTKeyPath     string        `mapstructure:"jwt-key"              desc:"Path to the JWT certificate key pem file"              secret:"true" file:"true"`
	JWTMaxValidity time.Duration `mapstructure:"jwt-max-validity"     desc:"Maximum duration of the validity of the issued tokens" default:"720h"`

	jwtCert *x509.Certificate
	jwtKey  crypto.PrivateKey
}

// JWTCertificate returns the certificate used to verify JWTs.
func (c *Conf) JWTCertificate() (*x509.Certificate, crypto.PrivateKey, error) {

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
