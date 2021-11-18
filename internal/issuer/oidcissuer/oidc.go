package oidcissuer

import (
	"fmt"
	"strings"

	"go.aporeto.io/a3s/pkgs/token"
)

// New returns a new Azure issuer.
func New(claims map[string]interface{}) (token.Issuer, error) {

	c := newOIDCIssuer()
	if err := c.fromClaims(claims); err != nil {
		return nil, err
	}

	return c, nil
}

type oidcIssuer struct {
	token *token.IdentityToken
}

func newOIDCIssuer() *oidcIssuer {
	return &oidcIssuer{
		token: token.NewIdentityToken(token.Source{
			Type: "oidc",
		}),
	}
}

// Issue returns the IdentityToken.
func (c *oidcIssuer) Issue() *token.IdentityToken {

	return c.token
}

func (c *oidcIssuer) fromClaims(claims map[string]interface{}) (err error) {

	for k := range claims {
		if strings.HasPrefix(k, "@source") {
			return ErrOIDC{Err: fmt.Errorf("cannot handle claims starting with '@' as this is reserved")}
		}
	}

	c.token.Identity = computeOIDClaims(claims)

	return nil
}

func computeOIDClaims(claims map[string]interface{}) []string {

	out := []string{}

	for k, v := range claims {
		switch claim := v.(type) {
		case string:
			out = append(out, fmt.Sprintf("%s=%s", k, claim))
		case []string:
			for _, item := range claim {
				out = append(out, fmt.Sprintf("%s=%s", k, item))
			}
		case int:
			out = append(out, fmt.Sprintf("%s=%d", k, claim))
		case []int:
			for _, item := range claim {
				out = append(out, fmt.Sprintf("%s=%d", k, item))
			}
		case bool:
			out = append(out, fmt.Sprintf("%s=%t", k, claim))
		case []interface{}:
			for _, item := range claim {
				if claimValue, ok := item.(string); ok {
					out = append(out, fmt.Sprintf("%s=%s", k, claimValue))
				}
			}
		}
	}

	return out
}
