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
	requiredIssuer string,
	audience jwt.ClaimStrings,
	validity time.Duration,
) (token.Issuer, error) {

	c := newA3SIssuer()
	if err := c.fromToken(
		tokenString,
		keychain,
		requiredIssuer,
		audience,
		validity,
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
	audience jwt.ClaimStrings,
	validity time.Duration,
) error {

	orest, err := permissions.GetRestrictions(tokenString)
	if err != nil {
		return ErrComputeRestrictions{Err: err}
	}

	if c.token, err = token.Parse(tokenString, keychain, issuer, ""); err != nil {
		return ErrInputToken{Err: err}
	}

	if len(audience) == 0 && len(c.token.Audience) != 0 {
		return ErrInputToken{Err: fmt.Errorf("you cannot request a token with no audience from a token that has one")}
	}

	for _, aud := range audience {
		if !c.token.VerifyAudience(aud, true) {
			return ErrInputToken{Err: fmt.Errorf("requested audience '%s' is not declared in initial token", aud)}
		}
	}

	if !orest.Zero() {
		c.token.Restrictions = &orest
	}

	c.token.ExpiresAt, err = computeNewValidity(c.token.ExpiresAt, validity, c.token.Refresh)
	if err != nil {
		return ErrComputeRestrictions{Err: err}
	}

	return nil
}

// Issue issues a token.IdentityToken derived from the initial token.
func (c *a3sIssuer) Issue() *token.IdentityToken {

	return c.token
}

func computeNewValidity(originalExpUNIX *jwt.NumericDate, requestedValidity time.Duration, isRefresh bool) (*jwt.NumericDate, error) {

	if originalExpUNIX == nil || originalExpUNIX.Unix() == 0 {
		return nil, fmt.Errorf("unable to compute new validity: original expiration is zero")
	}

	if requestedValidity == 0 {
		return originalExpUNIX, nil
	}

	now := time.Now()

	originalExp := originalExpUNIX.Local()
	if now.Add(requestedValidity).After(originalExp) && !isRefresh {
		return nil, fmt.Errorf("the request validity is greater than the original non refresh token")
	}

	return jwt.NewNumericDate(now.Add(requestedValidity)), nil
}
