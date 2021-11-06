package permissions

import (
	"fmt"
	"net"

	"github.com/apparentlymart/go-cidr/cidr"
	"github.com/golang-jwt/jwt/v4"
	"go.aporeto.io/elemental"
)

// ErrRestrictionsViolation represents an error
// during restrictions computations.
type ErrRestrictionsViolation struct {
	Err error
}

func (e ErrRestrictionsViolation) Error() string {
	return fmt.Sprintf("restriction violation: %s", e.Err)
}

func (e ErrRestrictionsViolation) Unwrap() error {
	return e.Err
}

// Restrictions are a collection of restrictions
// that the policy engine should apply for authz based
// on the token
type Restrictions struct {
	Namespace   string   `json:"namespace"`
	Permissions []string `json:"perms"`
	Networks    []string `json:"networks"`
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

// RestrictNamespace returns the namespace to use based on the
// receiver and the new requested one.
func (r Restrictions) RestrictNamespace(requested string) (string, error) {

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
		return "", ErrRestrictionsViolation{
			Err: fmt.Errorf("restricted namespace must be empty, '%s' or one of its children", r.Namespace),
		}
	}
}

// RestrictNetworks returns the networks to use based on the
// receiver and the new requested ones.
func (r Restrictions) RestrictNetworks(requested []string) ([]string, error) {

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
				return nil, ErrRestrictionsViolation{
					Err: fmt.Errorf("restricted networks must not overlap the current ones"),
				}
			}
		}

		return requested, nil
	}
}

// RestrictPermissions returns the permissions to use based on the
// receiver and the new requested ones.
func (r Restrictions) RestrictPermissions(requested []string) ([]string, error) {

	if len(requested) == 0 {
		return r.Permissions, nil
	}

	if len(r.Permissions) == 0 {
		return requested, nil
	}

	if !Parse(r.Permissions, "").Contains(Parse(requested, "")) {
		return nil, ErrRestrictionsViolation{
			Err: fmt.Errorf("restricted permissions must not be more permissive than the current ones"),
		}
	}

	return requested, nil
}
