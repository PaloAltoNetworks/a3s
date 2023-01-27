package token

import (
	"context"
	"crypto/sha256"
	"crypto/tls"
	"crypto/x509"
	"fmt"
	"net/http"
	"strings"

	"github.com/golang-jwt/jwt/v4"
	"go.aporeto.io/bahamut"
	"go.aporeto.io/elemental"
)

// FromRequest retrieves the token from the given elemental.Request
// first looking at the cookie x-a3s-token, then the request.Password.
func FromRequest(req *elemental.Request) string {
	if hreq := req.HTTPRequest(); hreq != nil {
		if cookie, err := hreq.Cookie("x-a3s-token"); err == nil {
			return cookie.Value
		}
	}
	return req.Password
}

// FromSession retrieves the token from the given bahamut.Session
// first looking at the cookie x-a3s-token, then the session.Token(.
func FromSession(session bahamut.Session) string {
	if cookie, err := session.Cookie("x-a3s-token"); err == nil {
		return cookie.Value
	}
	return session.Token()
}

// Fingerprint returns the fingerprint of the given certificate.
func Fingerprint(cert *x509.Certificate) string {

	return fmt.Sprintf("%02X", sha256.Sum256(cert.Raw)) // #nosec
}

// JWKSFromTokenIssuer will retrieve a remote JWKS from the issuer field
// in the given idt, using the eventually given tlsConfig to retrieve the JWKS..
// You usually want to pass a non verified IdentityToken here (from ParseUnverified for instance)
// so you can then correctly verify it using Parse().
func JWKSFromTokenIssuer(ctx context.Context, idt *IdentityToken, tlsConfig *tls.Config) (*JWKS, error) {

	wellKnownSuffix := ".well-known/jwks.json"

	endpoint := idt.Issuer
	if !strings.HasSuffix(endpoint, wellKnownSuffix) {
		endpoint = fmt.Sprintf("%s/%s", strings.TrimRight(endpoint, "/"), wellKnownSuffix)
	}

	client := &http.Client{
		Transport: &http.Transport{
			TLSClientConfig: tlsConfig,
		},
	}

	jwks, err := NewRemoteJWKS(ctx, client, endpoint)
	if err != nil {
		return nil, fmt.Errorf("unable to retrieve remote jwks: %w", err)
	}

	return jwks, nil
}

func makeKeyFunc(keychain *JWKS) jwt.Keyfunc {

	return func(token *jwt.Token) (any, error) {

		if _, ok := token.Method.(*jwt.SigningMethodECDSA); !ok {
			return nil, fmt.Errorf("unexpected signing method: %s", token.Header["alg"])
		}

		kid, ok := token.Header["kid"].(string)
		if !ok || kid == "" {
			return nil, fmt.Errorf("token has no KID in its header")
		}

		k, err := keychain.Get(kid)
		if err != nil {
			return nil, fmt.Errorf("unable to find kid '%s': %w", kid, err)
		}

		return k.PublicKey(), nil
	}
}
