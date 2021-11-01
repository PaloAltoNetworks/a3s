package claims

import (
	"crypto/x509"
	"fmt"
	"time"

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
	cert *x509.Certificate,
	issuer string,
	audience string,
	validityStr string,
	restrictions permissions.Restrictions,
) error {

	orest, err := permissions.GetRestrictions(tokenString)
	if err != nil {
		return ErrComputeRestrictions{Err: err}
	}

	restrictedNamespace, err := orest.ComputeNamespaceRestriction(restrictions.Namespace)
	if err != nil {
		return ErrComputeRestrictions{Err: err}
	}

	restrictedNetworks, err := orest.ComputeNetworkRestrictions(restrictions.Networks)
	if err != nil {
		return ErrComputeRestrictions{Err: err}
	}

	restrictedPermissions, err := orest.ComputePermissionsRestrictions(restrictions.Permissions)
	if err != nil {
		return ErrComputeRestrictions{Err: err}
	}

	c.token = token.NewIdentityToken(token.Source{})
	if err := c.token.Parse(tokenString, cert, issuer, audience); err != nil {
		return ErrInputToken{Err: err}
	}

	c.token.ExpiresAt, err = computeNewValidity(c.token.ExpiresAt, validityStr)
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

func computeNewValidity(originalExpUNIX int64, requestedValidityStr string) (int64, error) {

	if originalExpUNIX == 0 {
		return 0, fmt.Errorf("unable to compute new validity: original expiration is zero")
	}

	if requestedValidityStr == "" {
		return originalExpUNIX, nil
	}

	now := time.Now()

	requestedValidity, err := time.ParseDuration(requestedValidityStr)
	if err != nil {
		return 0, err
	}

	originalExp := time.Unix(originalExpUNIX, 0)
	if now.Add(requestedValidity).After(originalExp) {
		return originalExp.Unix(), nil
	}

	return now.Add(requestedValidity).Unix(), nil
}
