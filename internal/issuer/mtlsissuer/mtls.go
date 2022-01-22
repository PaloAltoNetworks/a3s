package mtlsissuer

import (
	"context"
	"crypto/x509"
	"fmt"
	"strings"

	"go.aporeto.io/a3s/internal/identitymodifier"
	"go.aporeto.io/a3s/pkgs/api"
	"go.aporeto.io/a3s/pkgs/token"
)

// New returns a new MTLS issuer.
func New(ctx context.Context, source *api.MTLSSource, cert *x509.Certificate) (token.Issuer, error) {

	c := newMTLSIssuer(source)
	if err := c.fromCertificate(ctx, cert); err != nil {
		return nil, err
	}

	return c, nil
}

type mtlsIssuer struct {
	token  *token.IdentityToken
	source *api.MTLSSource
}

func newMTLSIssuer(source *api.MTLSSource) *mtlsIssuer {

	return &mtlsIssuer{
		source: source,
		token: token.NewIdentityToken(token.Source{
			Type:      "mtls",
			Namespace: source.Namespace,
			Name:      source.Name,
		}),
	}
}

func (c *mtlsIssuer) fromCertificate(ctx context.Context, cert *x509.Certificate) error {

	pool := x509.NewCertPool()
	if !pool.AppendCertsFromPEM([]byte(c.source.CA)) {
		return fmt.Errorf("unable to prepare x509 verifier: could not append cert from source.CA")
	}

	chains, err := cert.Verify(
		x509.VerifyOptions{
			Roots: pool,
			KeyUsages: []x509.ExtKeyUsage{
				x509.ExtKeyUsageClientAuth,
			},
		},
	)
	if err != nil {
		return fmt.Errorf("unable to verify certificate: %w", err)
	}

	var fingerprints []string
	for _, chain := range chains {
		for _, cert := range chain {
			fingerprints = append(
				fingerprints,
				token.Fingerprint(cert),
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

	if len(cert.AuthorityKeyId) > 0 {
		c.token.Identity = append(c.token.Identity, fmt.Sprintf("akid=%02X", cert.AuthorityKeyId))
	}

	if srcmod := c.source.Modifier; srcmod != nil {

		m, err := identitymodifier.NewRemote(srcmod, c.token.Source)
		if err != nil {
			return fmt.Errorf("unable to prepare source modifier: %w", err)
		}

		if c.token.Identity, err = m.Modify(ctx, c.token.Identity); err != nil {
			return fmt.Errorf("unable to call modifier: %w", err)
		}
	}

	return nil
}

// Issue issues the token.IdentityToken derived from the the user certificate.
func (c *mtlsIssuer) Issue() *token.IdentityToken {

	return c.token
}
