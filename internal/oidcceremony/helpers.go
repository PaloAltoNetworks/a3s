package oidcceremony

import (
	"crypto/rand"
	"crypto/sha256"
	"crypto/tls"
	"crypto/x509"
	"encoding/base64"
	"encoding/json"
	"fmt"
	"net/http"
	"time"

	"go.aporeto.io/bahamut"
)

// GenerateNonce generate a nonce.
func GenerateNonce(nonceSourceSize int) (string, error) {

	nonceSource := make([]byte, nonceSourceSize)
	_, err := rand.Read(nonceSource)
	if err != nil {
		return "", err
	}
	sha := sha256.Sum256(nonceSource) // #nosec

	return base64.RawStdEncoding.EncodeToString(sha[:]), nil
}

// MakeOIDCProviderClient returns a OIDC client using the given CA.
func MakeOIDCProviderClient(ca string) (*http.Client, error) {

	var pool *x509.CertPool
	var err error

	if ca != "" {
		pool = x509.NewCertPool()
		if !pool.AppendCertsFromPEM([]byte(ca)) {
			return nil, fmt.Errorf("unable to append given ca to ca pool")
		}
	} else {
		pool, err = x509.SystemCertPool()
		if err != nil {
			return nil, fmt.Errorf("unable to initialize system root ca pool: %s", err)
		}
	}

	return &http.Client{
		Timeout: 10 * time.Second,
		Transport: &http.Transport{
			TLSClientConfig: &tls.Config{
				RootCAs: pool,
			},
			Proxy: http.ProxyFromEnvironment,
		},
	}, nil
}

// RedirectErrorEventually will configure the redirect url if given for the
// given bahamut.Context
func RedirectErrorEventually(ctx bahamut.Context, url string, err error) error {

	if url == "" {
		return err
	}

	d, e := json.Marshal(err)
	if e != nil {
		return err
	}

	ctx.SetRedirect(fmt.Sprintf("%s?error=%s", url, string(d)))

	return nil
}
