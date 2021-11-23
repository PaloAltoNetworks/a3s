package remotea3sissuer

import (
	"context"
	"crypto/tls"
	"crypto/x509"
	"fmt"
	"net/http"
	"strings"

	"go.aporeto.io/a3s/pkgs/api"
	"go.aporeto.io/a3s/pkgs/token"
)

const wellKnownSuffix = ".well-known/jwks.json"

// New retrurns new Remote A3S issuer.
func New(
	ctx context.Context,
	source *api.A3SSource,
	tokenString string,
) (token.Issuer, error) {

	c := newRemoteA3SIssuer(source)
	if err := c.fromToken(ctx, tokenString); err != nil {
		return nil, err
	}

	return c, nil
}

type remoteA3SIssuer struct {
	token  *token.IdentityToken
	source *api.A3SSource
}

func newRemoteA3SIssuer(source *api.A3SSource) *remoteA3SIssuer {
	return &remoteA3SIssuer{
		source: source,
		token: token.NewIdentityToken(token.Source{
			Type:      "remotea3s",
			Namespace: source.Namespace,
			Name:      source.Name,
		}),
	}
}

func (c *remoteA3SIssuer) fromToken(ctx context.Context, tokenString string) error {

	endpoint := c.source.Endpoint
	if endpoint == "" {
		endpoint = c.source.Issuer
	}
	if !strings.HasSuffix(endpoint, wellKnownSuffix) {
		endpoint = fmt.Sprintf("%s/%s", strings.TrimRight(endpoint, "/"), wellKnownSuffix)
	}

	root := x509.NewCertPool()
	root.AppendCertsFromPEM([]byte(c.source.CertificateAuthority))
	client := &http.Client{
		Transport: &http.Transport{
			TLSClientConfig: &tls.Config{
				RootCAs: root,
			},
		},
	}

	jwks, err := token.NewRemoteJWKS(ctx, client, endpoint)
	if err != nil {
		return ErrRemoteA3S{Err: fmt.Errorf("unable to retrieve remote jwks: %w", err)}
	}

	idt, err := token.Parse(tokenString, jwks, c.source.Issuer, c.source.Audience)
	if err != nil {
		return ErrRemoteA3S{Err: fmt.Errorf("unable to parse remote a3s token: %w", err)}

	}

	c.token.Identity = make([]string, len(idt.Identity))
	var i int
	for _, claim := range idt.Identity {
		if strings.HasPrefix(claim, "@") {
			continue
		}
		c.token.Identity[i] = claim
		i++
	}
	c.token.Identity = c.token.Identity[:i]

	return nil
}

// Issue issues a token.IdentityToken derived from the initial token.
func (c *remoteA3SIssuer) Issue() *token.IdentityToken {

	return c.token
}
