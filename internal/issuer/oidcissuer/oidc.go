package oidcissuer

import (
	"fmt"
	"sort"
	"strings"

	"go.aporeto.io/a3s/pkgs/token"
)

// New returns a new Azure issuer.
func New(claims map[string]interface{}) token.Issuer {

	c := newOIDCIssuer()
	c.fromClaims(claims)
	return c
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

func (c *oidcIssuer) fromClaims(claims map[string]interface{}) {

	c.token.Identity = computeOIDClaims(claims)
}

func computeOIDClaims(claims map[string]interface{}) []string {

	out := []string{}

	for k, v := range claims {
		k = strings.TrimLeft(k, "@")
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
		case float64:
			out = append(out, fmt.Sprintf("%s=%f", k, claim))
		case []float64:
			for _, item := range claim {
				out = append(out, fmt.Sprintf("%s=%f", k, item))
			}
		case bool:
			out = append(out, fmt.Sprintf("%s=%t", k, claim))
		case []interface{}:
			for _, item := range claim {
				if claimValue, ok := item.(string); ok {
					out = append(out, fmt.Sprintf("%s=%s", strings.TrimLeft(k, "@"), claimValue))
				}
			}
		default:
			out = append(out, fmt.Sprintf("%s=%s", k, v))
		}
	}

	sort.Strings(out)

	return out
}
