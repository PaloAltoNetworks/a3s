package issuer

import (
	"fmt"
	"time"

	"github.com/golang-jwt/jwt/v4"
	"go.aporeto.io/a3s/pkgs/permissions"
	"go.aporeto.io/a3s/pkgs/token"
)

// An ErrComputeRestrictions represents a generic
// restrictions computation error.
type ErrComputeRestrictions struct {
	Err error
}

func (e ErrComputeRestrictions) Error() string {
	return fmt.Sprintf("unable to compute restrictions: %s", e.Err)
}

func (e ErrComputeRestrictions) Unwrap() error {
	return e.Err
}

// An ErrInputToken represents a generic input token error.
type ErrInputToken struct {
	Err error
}

func (e ErrInputToken) Error() string {
	return fmt.Sprintf("unable to parse input token: %s", e.Err)
}

func (e ErrInputToken) Unwrap() error {
	return e.Err
}

// TokenIssuer represents the claims from a PCC token.
type TokenIssuer struct {
	token *token.IdentityToken
}

// NewTokenIssuer returns a new TokenIssuer.
func NewTokenIssuer() *TokenIssuer {
	return &TokenIssuer{}
}

// FromToken reads the claims from original token.
func (c *TokenIssuer) FromToken(
	tokenString string,
	keychain *token.JWKS,
	issuer string,
	audience string,
	validity time.Duration,
	restrictions permissions.Restrictions,
) error {

	orest, err := permissions.GetRestrictions(tokenString)
	if err != nil {
		return ErrComputeRestrictions{Err: err}
	}

	restrictedNamespace, err := orest.RestrictNamespace(restrictions.Namespace)
	if err != nil {
		return ErrComputeRestrictions{Err: err}
	}

	restrictedNetworks, err := orest.RestrictNetworks(restrictions.Networks)
	if err != nil {
		return ErrComputeRestrictions{Err: err}
	}

	restrictedPermissions, err := orest.RestrictPermissions(restrictions.Permissions)
	if err != nil {
		return ErrComputeRestrictions{Err: err}
	}

	c.token = &token.IdentityToken{
		RegisteredClaims: jwt.RegisteredClaims{
			Issuer:   issuer,
			Audience: jwt.ClaimStrings{audience},
		},
	}
	if err := c.token.Parse(tokenString, keychain, issuer, audience); err != nil {
		return ErrInputToken{Err: err}
	}

	c.token.ExpiresAt, err = computeNewValidity(c.token.ExpiresAt, validity)
	if err != nil {
		return ErrComputeRestrictions{Err: err}
	}

	if restrictedNamespace != "" || len(restrictedPermissions) > 0 || len(restrictedNetworks) > 0 {
		c.token.Restrictions = &permissions.Restrictions{
			Namespace:   restrictedNamespace,
			Networks:    restrictedNetworks,
			Permissions: restrictedPermissions,
		}
	}

	return nil
}

// Issue issues a token.IdentityToken derived from the initial token.
func (c *TokenIssuer) Issue() *token.IdentityToken {

	return c.token
}

func computeNewValidity(originalExpUNIX *jwt.NumericDate, requestedValidity time.Duration) (*jwt.NumericDate, error) {

	if originalExpUNIX == nil || originalExpUNIX.Unix() == 0 {
		return nil, fmt.Errorf("unable to compute new validity: original expiration is zero")
	}

	if requestedValidity == 0 {
		return originalExpUNIX, nil
	}

	now := time.Now()

	originalExp := originalExpUNIX.Local()
	if now.Add(requestedValidity).After(originalExp) {
		return jwt.NewNumericDate(originalExp), nil
	}

	return jwt.NewNumericDate(now.Add(requestedValidity)), nil
}
