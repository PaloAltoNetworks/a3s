package token

import (
	"crypto"
	"crypto/ecdsa"
	"crypto/x509"
	"fmt"
	"time"

	"github.com/gofrs/uuid"
	"github.com/golang-jwt/jwt/v4"
	"go.aporeto.io/a3s/pkgs/permissions"
)

// An Issuer is an object that can issue IdentityToken.
type Issuer interface {
	Issue() *IdentityToken
}

// A Source represents the authentication source info
// used to derive an IdentityToken.
type Source struct {
	Type      string `json:"type"`
	Namespace string `json:"namespace,omitempty"`
	Name      string `json:"name,omitempty"`
}

// An IdentityToken represents a normalized identity token.
type IdentityToken struct {
	Identity     []string                  `json:"identity"`
	Opaque       map[string]string         `json:"opaque,omitempty"`
	Restrictions *permissions.Restrictions `json:"restrictions,omitempty"`
	Source       Source                    `json:"source"`

	jwt.StandardClaims
}

// NewIdentityToken returns a new IdentityToken with the
// given Source.
func NewIdentityToken(source Source) *IdentityToken {
	return &IdentityToken{
		Source: source,
	}
}

// Parse returns a validated IdentityToken from the given token string.
func (t *IdentityToken) Parse(tokenString string, cert *x509.Certificate, issuer string, audience string) error {

	token, err := jwt.ParseWithClaims(tokenString, t, func(token *jwt.Token) (interface{}, error) {
		if _, ok := token.Method.(*jwt.SigningMethodECDSA); ok {
			return cert.PublicKey.(*ecdsa.PublicKey), nil
		}
		return nil, fmt.Errorf("unexpected signing method: %s", token.Header["alg"])
	})

	if err != nil {
		return fmt.Errorf("unable to parse jwt: %w", err)
	}

	claims := token.Claims.(*IdentityToken)

	if !claims.VerifyIssuer(issuer, true) {
		return fmt.Errorf("issuer '%s' is not acceptable. want '%s'", claims.Issuer, issuer)
	}

	if !claims.VerifyAudience(audience, false) {
		return fmt.Errorf("audience '%s' is not acceptable. want '%s'", claims.Audience, audience)
	}

	return nil
}

// JWT returns the signed JWT string.
func (t *IdentityToken) JWT(key crypto.PrivateKey, kid string, exp time.Time) (string, error) {

	t.Id = uuid.Must(uuid.NewV4()).String()
	t.IssuedAt = time.Now().Unix()

	if !exp.IsZero() {
		t.ExpiresAt = exp.Unix()
	}

	if t.Source.Type == "" {
		return "", fmt.Errorf("invalid identity token: missing source type")
	}

	t.Identity = append(t.Identity, fmt.Sprintf("@sourcetype=%s", t.Source.Type))

	if t.Source.Namespace != "" {
		t.Identity = append(t.Identity, fmt.Sprintf("@sourcenamespace=%s", t.Source.Namespace))
	}

	if t.Source.Name != "" {
		t.Identity = append(t.Identity, fmt.Sprintf("@sourcename=%s", t.Source.Name))
	}

	j := jwt.NewWithClaims(jwt.SigningMethodES256, t)

	if kid != "" {
		j.Header["kid"] = kid
	}

	return j.SignedString(key)
}
