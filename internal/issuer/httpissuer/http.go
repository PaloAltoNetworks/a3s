package httpissuer

import (
	"bytes"
	"context"
	"crypto/tls"
	"crypto/x509"
	"encoding/json"
	"fmt"
	"net/http"
	"strings"

	"go.aporeto.io/a3s/internal/identitymodifier"
	"go.aporeto.io/a3s/pkgs/api"
	"go.aporeto.io/a3s/pkgs/token"
	"go.aporeto.io/tg/tglib"
)

// A Credentials represents user provided Credentials.
type Credentials struct {
	Username string `json:"username,omitempty"`
	Password string `json:"password,omitempty"`
	TOTP     string `json:"totp,omitempty"`
}

// New retrurns new Remote http issuer.
func New(
	ctx context.Context,
	source *api.HTTPSource,
	creds Credentials,
) (token.Issuer, error) {

	c := newHTTPIssuer(source)
	if err := c.fromCredentials(ctx, creds); err != nil {
		return nil, err
	}

	return c, nil
}

type httpIssuer struct {
	token  *token.IdentityToken
	source *api.HTTPSource
}

func newHTTPIssuer(source *api.HTTPSource) *httpIssuer {
	return &httpIssuer{
		source: source,
		token: token.NewIdentityToken(token.Source{
			Type:      "http",
			Namespace: source.Namespace,
			Name:      source.Name,
		}),
	}
}

func (c *httpIssuer) fromCredentials(ctx context.Context, creds Credentials) error {

	root := x509.NewCertPool()
	root.AppendCertsFromPEM([]byte(c.source.CA))

	clientCert, clientKey, err := tglib.ReadCertificate([]byte(c.source.Certificate), []byte(c.source.Key), "")
	if err != nil {
		return ErrHTTP{Err: fmt.Errorf("unable to read certificate: %w", err)}
	}
	cert, err := tglib.ToTLSCertificate(clientCert, clientKey)
	if err != nil {
		return ErrHTTP{Err: fmt.Errorf("unable to convert to tls certificate: %s", err)}
	}

	client := &http.Client{
		Transport: &http.Transport{
			TLSClientConfig: &tls.Config{
				RootCAs:      root,
				Certificates: []tls.Certificate{cert},
			},
		},
	}

	buf := &bytes.Buffer{}
	if err := json.NewEncoder(buf).Encode(creds); err != nil {
		return ErrHTTP{Err: fmt.Errorf("unable to encode body: %w", err)}
	}

	req, err := http.NewRequestWithContext(ctx, http.MethodPost, c.source.Endpoint, buf)
	if err != nil {
		return ErrHTTP{Err: fmt.Errorf("unable to build request: %w", err)}
	}
	req.Header.Add("Content-Type", "application/json")

	resp, err := client.Do(req)
	if err != nil {
		return ErrHTTP{Err: fmt.Errorf("unable to send request: %w", err)}
	}

	if resp.StatusCode != http.StatusOK {
		return ErrHTTPResponse{Err: fmt.Errorf("server responded with '%s'", resp.Status)}
	}

	var rawClaims []string
	if err := json.NewDecoder(resp.Body).Decode(&rawClaims); err != nil {
		return ErrHTTPResponse{Err: fmt.Errorf("unable to decode response body: %w", err)}
	}

	var i int
	claims := make([]string, len(rawClaims))
	for _, v := range rawClaims {
		if !strings.HasPrefix(v, "@") {
			claims[i] = v
			i++
		}
	}
	c.token.Identity = claims[:i]

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

// Issue issues a token.IdentityToken derived from the initial token.
func (c *httpIssuer) Issue() *token.IdentityToken {

	return c.token
}
