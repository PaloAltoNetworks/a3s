package identitymodifier

import (
	"bytes"
	"context"
	"crypto/tls"
	"crypto/x509"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"strings"

	"go.aporeto.io/a3s/pkgs/api"
	"go.aporeto.io/a3s/pkgs/token"
	"go.aporeto.io/tg/tglib"
)

// An IdentityModifier can modify a set of claims that are about
// to be delivered.
type IdentityModifier interface {

	// Modify returns the list of given claims, after eventually applying modifications.
	// It must complete before the given context expires.
	// It will return an error if the returned slice contains claims prefixed by '@'.
	Modify(ctx context.Context, in []string) (out []string, err error)
}

type identityModifier struct {
	caPool     *x509.CertPool
	url        string
	clientCert tls.Certificate
	method     string
	src        token.Source
}

// NewRemote returns a new HTTP based IdentityModifier.
// The remote server will receive the given method and the claims, either as
// a json array in the request body (for POST/PUT/PATCH) or in the query
// parameter `claim` (for GET). Any other http method will make the function
// to return an error.
// The server must return 200 if it modified the list, 204 if it did not.
// Anything else is considered as an error.
func NewRemote(m *api.IdentityModifier, src token.Source) (IdentityModifier, error) {

	switch strings.ToUpper(string(m.Method)) {
	case http.MethodGet, http.MethodPost, http.MethodPut, http.MethodPatch:
	default:
		return nil, fmt.Errorf("invalid http method: %s", m.Method)
	}

	var caPool *x509.CertPool
	var err error
	if len(m.CA) != 0 {
		caPool = x509.NewCertPool()
		caPool.AppendCertsFromPEM([]byte(m.CA))
	} else {
		caPool, err = x509.SystemCertPool()
		if err != nil {
			return nil, fmt.Errorf("unable to prepare system cert pool: %w", err)
		}
	}

	xc, xk, err := tglib.ReadCertificate([]byte(m.Certificate), []byte(m.Key), "")
	if err != nil {
		return nil, fmt.Errorf("unable to create certificate: %w", err)
	}

	clientCert, err := tglib.ToTLSCertificate(xc, xk)
	if err != nil {
		return nil, fmt.Errorf("unable to convert to tls.Certificate: %w", err)
	}

	return &identityModifier{
		url:        m.URL,
		caPool:     caPool,
		clientCert: clientCert,
		method:     string(m.Method),
		src:        src,
	}, nil
}

// Modify calls the remove service for modification.
func (m *identityModifier) Modify(ctx context.Context, in []string) (out []string, err error) {

	client := &http.Client{
		Transport: &http.Transport{
			TLSClientConfig: &tls.Config{
				RootCAs:      m.caPool,
				Certificates: []tls.Certificate{m.clientCert},
			},
		},
	}

	var buffer io.Reader

	switch m.method {
	case http.MethodGet:
	case http.MethodPost, http.MethodPut, http.MethodPatch:
		data, err := json.Marshal(in)
		if err != nil {
			return nil, fmt.Errorf("unable to encode body: %w", err)
		}
		buffer = bytes.NewBuffer(data)
	}

	req, err := http.NewRequestWithContext(ctx, m.method, m.url, buffer)
	if err != nil {
		return nil, fmt.Errorf("unable to build http request: %w", err)
	}

	if m.method == http.MethodGet {
		values := req.URL.Query()
		for _, c := range in {
			values.Add("claim", c)
		}
		req.URL.RawQuery = values.Encode()
	}

	req.Header.Set("x-a3s-source-type", m.src.Type)
	if v := m.src.Namespace; v != "" {
		req.Header.Set("x-a3s-source-namespace", v)
	}
	if v := m.src.Name; v != "" {
		req.Header.Set("x-a3s-source-name", v)
	}

	resp, err := client.Do(req)
	if err != nil {
		return nil, fmt.Errorf("unable to perform request: %w", err)
	}
	defer resp.Body.Close() // nolint: errcheck

	switch resp.StatusCode {
	case http.StatusOK:
	case http.StatusNoContent:
		return in, nil
	default:
		return nil, fmt.Errorf("service returned an error: %s", resp.Status)
	}

	if err := json.NewDecoder(resp.Body).Decode(&out); err != nil {
		return nil, fmt.Errorf("unable to decode response: %w", err)
	}

	for _, o := range out {
		if strings.HasPrefix(o, "@") {
			return nil, fmt.Errorf("invalid returned claim '%s': must not be prefixed by @", o)
		}
	}

	return out, nil
}
