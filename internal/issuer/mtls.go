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
	token  *token.IdentityToken
	source *api.MTLSSource
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
func (c *MTLSIssuer) FromCertificate(cert *x509.Certificate) error {

	pool := x509.NewCertPool()
	pool.AppendCertsFromPEM([]byte(c.source.CertificateAuthority))

	chains, err := cert.Verify(x509.VerifyOptions{
		Roots:     pool,
		KeyUsages: []x509.ExtKeyUsage{x509.ExtKeyUsageClientAuth},
	})
	if err != nil {
		return fmt.Errorf("unable Verify certificate: %w", err)
	}

	var fingerprints []string
	for _, chain := range chains {
		for _, cert := range chain {
			fingerprints = append(
				fingerprints,
				fmt.Sprintf("%02X", sha1.Sum(cert.Raw)),
			)
		}
	}

	if v := cert.Subject.CommonName; v != "" {
		c.token.Identity = append(c.token.Identity, fmt.Sprintf("commonname=%s", v))
	}

	if v := cert.SerialNumber; v != nil {
		c.token.Identity = append(c.token.Identity, fmt.Sprintf("serialnumber=%s", v))
	}

	if vs := cert.Subject.Country; len(vs) != 0 {
		for _, v := range vs {
			c.token.Identity = append(c.token.Identity, fmt.Sprintf("country=%s", v))
		}
	}

	if vs := cert.Subject.Locality; len(vs) != 0 {
		for _, v := range vs {
			c.token.Identity = append(c.token.Identity, fmt.Sprintf("locality=%s", v))
		}
	}

	if vs := cert.Subject.Organization; len(vs) != 0 {
		for _, v := range vs {
			c.token.Identity = append(c.token.Identity, fmt.Sprintf("organization=%s", v))
		}
	}

	if vs := cert.Subject.OrganizationalUnit; len(vs) != 0 {
		for _, v := range vs {
			c.token.Identity = append(c.token.Identity, fmt.Sprintf("organizationalunit=%s", v))
		}
	}

	if vs := cert.Subject.PostalCode; len(vs) != 0 {
		for _, v := range vs {
			c.token.Identity = append(c.token.Identity, fmt.Sprintf("postalcode=%s", v))
		}
	}

	if vs := cert.Subject.Province; len(vs) != 0 {
		for _, v := range vs {
			c.token.Identity = append(c.token.Identity, fmt.Sprintf("province=%s", v))
		}
	}

	if vs := cert.Subject.StreetAddress; len(vs) != 0 {
		for _, v := range vs {
			c.token.Identity = append(c.token.Identity, fmt.Sprintf("streetaddress=%s", v))
		}
	}

	if vs := cert.EmailAddresses; len(vs) != 0 {
		for _, v := range vs {
			c.token.Identity = append(c.token.Identity, fmt.Sprintf("email=%s", v))
		}
	}

	if vs := cert.DNSNames; len(vs) != 0 {
		for _, v := range vs {
			c.token.Identity = append(c.token.Identity, fmt.Sprintf("dnsname=%s", v))
		}
	}

	if len(fingerprints) > 0 {
		// if > 0 it is guaranteed to have at least 2 items.
		c.token.Identity = append(c.token.Identity, fmt.Sprintf("fingerprint=%s", fingerprints[0]))
		c.token.Identity = append(c.token.Identity, fmt.Sprintf("issuerchain=%s", strings.Join(fingerprints[1:], ",")))
	}

	return nil
}

// Issue issues the token.IdentityToken derived from the the user certificate.
func (c *MTLSIssuer) Issue() *token.IdentityToken {

	return c.token
}
