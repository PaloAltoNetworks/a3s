package permissions

import (
	"fmt"
	"net"

	"github.com/apparentlymart/go-cidr/cidr"
	"github.com/golang-jwt/jwt"
	"go.aporeto.io/elemental"
)

// Restrictions are a collection of restrictions
// that the policy engine should apply for authz based
// on the token
type Restrictions struct {
	Namespace   string   `json:"namespace"`
	Permissions []string `json:"perms"`
	Networks    []string `json:"networks"`
}

// ComputeNamespaceRestriction will return the namespace to use based on the
// receiver and the new requested one.
func (r Restrictions) ComputeNamespaceRestriction(requested string) (string, error) {

	switch {

	case r.Namespace == "":
		return requested, nil

	case requested == "":
		return r.Namespace, nil

	case r.Namespace == requested:
		return r.Namespace, nil

	case elemental.IsNamespaceChildrenOfNamespace(requested, r.Namespace):
		return requested, nil

	default:
		return "", fmt.Errorf("the new namespace restriction must be empty, '%s' or one of its children", r.Namespace)
	}
}

// ComputeNetworkRestrictions will return the networks to use based on the
// receiver and the new requested ones.
func (r Restrictions) ComputeNetworkRestrictions(requested []string) ([]string, error) {

	switch {

	case len(requested) == 0:
		return r.Networks, nil

	case len(r.Networks) == 0:
		return requested, nil

	default:

		for _, substr := range requested {

			_, sub, err := net.ParseCIDR(substr)
			if err != nil {
				return nil, err
			}

			valid := false

			for _, osubstr := range r.Networks {

				_, osub, err := net.ParseCIDR(osubstr)
				if err != nil {
					return nil, err
				}

				valid = valid || cidr.VerifyNoOverlap([]*net.IPNet{sub}, osub) == nil
			}

			if !valid {
				return nil, fmt.Errorf("the new network restrictions must not overlap any of the original ones")
			}
		}

		return requested, nil
	}
}

// ComputePermissionsRestrictions will return the networks to use based on the
// receiver and the new requested ones.
func (r Restrictions) ComputePermissionsRestrictions(requested []string) ([]string, error) {

	if len(requested) == 0 {
		return r.Permissions, nil
	}

	if len(r.Permissions) == 0 {
		return requested, nil
	}

	if !Contains(Parse(r.Permissions, ""), Parse(requested, "")) {
		return nil, fmt.Errorf("the new permissions restrictions must not be broader than the existing ones")
	}

	return requested, nil
}

// GetRestrictions returns the eventual Restrictions
// embedded in the given token.
func GetRestrictions(tokenString string) (Restrictions, error) {

	s := struct {
		R Restrictions `json:"restrictions"`
		jwt.Claims
	}{}

	parser := jwt.Parser{}
	if _, _, err := parser.ParseUnverified(tokenString, &s); err != nil {
		return Restrictions{}, fmt.Errorf("unable to compute authz restrictions from token: %w", err)
	}

	return s.R, nil
}
