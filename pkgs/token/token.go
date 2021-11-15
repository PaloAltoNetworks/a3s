package token

import (
	"crypto"
	"fmt"
	"strings"
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
	Source       Source                    `json:"-"`

	jwt.RegisteredClaims
}

// NewIdentityToken returns a new IdentityToken with the
// given Source.
func NewIdentityToken(source Source) *IdentityToken {
	return &IdentityToken{
		Source: source,
	}
}

// Parse returns a validated IdentityToken from the given token string.
func Parse(tokenString string, keychain *JWKS, issuer string, audience string) (*IdentityToken, error) {

	t := &IdentityToken{}
	token, err := jwt.ParseWithClaims(tokenString, t, makeKeyFunc(keychain))
	if err != nil {
		return nil, fmt.Errorf("unable to parse jwt: %w", err)
	}

	claims := token.Claims.(*IdentityToken)

	for _, c := range claims.Identity {
		switch {
		case strings.HasPrefix(c, "@sourcename="):
			claims.Source.Name = strings.TrimPrefix(c, "@sourcename=")
		case strings.HasPrefix(c, "@sourcenamespace="):
			claims.Source.Namespace = strings.TrimPrefix(c, "@sourcenamespace=")
		case strings.HasPrefix(c, "@sourcetype="):
			claims.Source.Type = strings.TrimPrefix(c, "@sourcetype=")
		}
	}

	if claims.Issuer != issuer {
		return nil, fmt.Errorf("issuer '%s' is not acceptable. want '%s'", claims.Issuer, issuer)
	}

	if !claims.VerifyAudience(audience, false) {
		return nil, fmt.Errorf("audience '%s' is not acceptable. want '%s'", claims.Audience, audience)
	}

	return t, nil
}

// JWT returns the signed JWT string.
func (t *IdentityToken) JWT(key crypto.PrivateKey, kid string, iss string, aud jwt.ClaimStrings, exp time.Time, cloak []string) (string, error) {

	t.ID = uuid.Must(uuid.NewV4()).String()
	t.IssuedAt = jwt.NewNumericDate(time.Now())
	t.Issuer = iss
	t.Audience = aud

	if !exp.IsZero() {
		t.ExpiresAt = jwt.NewNumericDate(exp)
	}

	if len(cloak) > 0 {
		var cloakedIdentity []string
		for _, claim := range t.Identity {
			for _, c := range cloak {
				if strings.HasPrefix(claim, c) {
					cloakedIdentity = append(cloakedIdentity, claim)
				}
			}
		}
		t.Identity = cloakedIdentity
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
