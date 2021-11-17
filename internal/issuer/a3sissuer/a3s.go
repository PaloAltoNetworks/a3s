package a3sissuer

import (
	"fmt"
	"time"

	"github.com/golang-jwt/jwt/v4"
	"go.aporeto.io/a3s/pkgs/permissions"
	"go.aporeto.io/a3s/pkgs/token"
)

// New retrurns new A3S issuer.
func New(
	tokenString string,
	keychain *token.JWKS,
	issuer string,
	audience string,
	validity time.Duration,
	restrictions permissions.Restrictions,
) (token.Issuer, error) {

	c := newA3SIssuer()
	if err := c.fromToken(
		tokenString,
		keychain,
		issuer,
		audience,
		validity,
		restrictions,
	); err != nil {
		return nil, err
	}

	return c, nil
}

type a3sIssuer struct {
	token *token.IdentityToken
}

func newA3SIssuer() *a3sIssuer {
	return &a3sIssuer{}
}

func (c *a3sIssuer) fromToken(
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

	if c.token, err = token.Parse(tokenString, keychain, issuer, audience); err != nil {
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
func (c *a3sIssuer) Issue() *token.IdentityToken {

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
