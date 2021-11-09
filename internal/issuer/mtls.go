package issuer

import (
	"crypto/sha1"
	"crypto/x509"
	"fmt"
	"strings"

	"go.aporeto.io/a3s/pkgs/api"
	"go.aporeto.io/a3s/pkgs/token"
)

// MTLSIssuer issues IdentityToken from a TLS certificate.
type MTLSIssuer struct {
	signerFingerprints []string
	token              *token.IdentityToken
	certificate        *x509.Certificate
	source             *api.MTLSSource
}

// NewMTLSIssuer returns a new MTLSIssuer.
func NewMTLSIssuer(source *api.MTLSSource) *MTLSIssuer {

	return &MTLSIssuer{
		source: source,
		token: token.NewIdentityToken(token.Source{
			Type:      "mtls",
			Namespace: source.Namespace,
			Name:      source.Name,
		}),
	}
}

// FromCertificate prepares the issuer according to the provided x509.Certificate.
func (c *MTLSIssuer) FromCertificate(certificate *x509.Certificate) error {

	pool := x509.NewCertPool()
	pool.AppendCertsFromPEM([]byte(c.source.CertificateAuthority))

	chains, err := certificate.Verify(x509.VerifyOptions{
		Roots:     pool,
		KeyUsages: []x509.ExtKeyUsage{x509.ExtKeyUsageClientAuth},
	})
	if err != nil {
		return fmt.Errorf("unable Verify certificate: %w", err)
	}

	for _, chain := range chains {
		for _, cert := range chain {
			c.signerFingerprints = append(
				c.signerFingerprints,
				fmt.Sprintf("%02X", sha1.Sum(cert.Raw)),
			)
		}
	}

	c.certificate = certificate

	return nil
}

// Issue issues the token.IdentityToken derived from the the user certificate.
func (c *MTLSIssuer) Issue() *token.IdentityToken {

	if v := c.certificate.Subject.CommonName; v != "" {
		c.token.Identity = append(c.token.Identity, fmt.Sprintf("commonname=%s", v))
	}

	if v := c.certificate.SerialNumber; v != nil {
		c.token.Identity = append(c.token.Identity, fmt.Sprintf("serialnumber=%s", v))
	}

	if vs := c.certificate.Subject.Country; len(vs) != 0 {
		for _, v := range vs {
			c.token.Identity = append(c.token.Identity, fmt.Sprintf("country=%s", v))
		}
	}

	if vs := c.certificate.Subject.Locality; len(vs) != 0 {
		for _, v := range vs {
			c.token.Identity = append(c.token.Identity, fmt.Sprintf("locality=%s", v))
		}
	}

	if vs := c.certificate.Subject.Organization; len(vs) != 0 {
		for _, v := range vs {
			c.token.Identity = append(c.token.Identity, fmt.Sprintf("organization=%s", v))
		}
	}

	if vs := c.certificate.Subject.OrganizationalUnit; len(vs) != 0 {
		for _, v := range vs {
			c.token.Identity = append(c.token.Identity, fmt.Sprintf("organizationalunit=%s", v))
		}
	}

	if vs := c.certificate.Subject.PostalCode; len(vs) != 0 {
		for _, v := range vs {
			c.token.Identity = append(c.token.Identity, fmt.Sprintf("postalcode=%s", v))
		}
	}

	if vs := c.certificate.Subject.Province; len(vs) != 0 {
		for _, v := range vs {
			c.token.Identity = append(c.token.Identity, fmt.Sprintf("province=%s", v))
		}
	}

	if vs := c.certificate.Subject.StreetAddress; len(vs) != 0 {
		for _, v := range vs {
			c.token.Identity = append(c.token.Identity, fmt.Sprintf("streetaddress=%s", v))
		}
	}

	if vs := c.certificate.EmailAddresses; len(vs) != 0 {
		for _, v := range vs {
			c.token.Identity = append(c.token.Identity, fmt.Sprintf("email=%s", v))
		}
	}

	if vs := c.certificate.DNSNames; len(vs) != 0 {
		for _, v := range vs {
			c.token.Identity = append(c.token.Identity, fmt.Sprintf("dnsname=%s", v))
		}
	}

	if len(c.signerFingerprints) > 0 {
		// if > 0 it is guaranteed to have at least 2 items.
		c.token.Identity = append(c.token.Identity, fmt.Sprintf("fingerprint=%s", c.signerFingerprints[0]))
		c.token.Identity = append(c.token.Identity, fmt.Sprintf("issuerchain=%s", strings.Join(c.signerFingerprints[1:], ",")))
	}

	return c.token
}
