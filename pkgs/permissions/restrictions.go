package permissions

import (
	"encoding/json"
	"errors"
	"fmt"
	"net"
	"strings"

	"github.com/apparentlymart/go-cidr/cidr"
	"github.com/golang-jwt/jwt"
	"go.aporeto.io/elemental"
)

// Restrictions are a collection of restrictions
// that the policy engine should apply for authz based
// on the token
type Restrictions struct {
	Namespace   string
	Permissions []string
	Networks    []string
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

	if !Contains(
		ResolveRestrictions(Restrictions{Permissions: r.Permissions}),
		ResolveRestrictions(Restrictions{Permissions: requested}),
	) {
		return nil, fmt.Errorf("the new permissions restrictions must not be broader than the existing ones")
	}

	return requested, nil
}

// GetRestrictions returns the eventual Restrictions
// embedded in the given token.
func GetRestrictions(token string) (Restrictions, error) {

	ns, perms, networks, err := ExtractRestrictions(token)
	if err != nil {
		return Restrictions{}, fmt.Errorf("unable to compute authz restrictions from token: %w", err)
	}

	return Restrictions{
		Namespace:   ns,
		Permissions: perms,
		Networks:    networks,
	}, nil
}

// ResolveRestrictions resolves the given restrictions into a standard permission map.
func ResolveRestrictions(restrictions Restrictions) PermissionMap {

	resolved := PermissionMap{}

	for _, perm := range restrictions.Permissions {

		parts := strings.Split(perm, ",")

		if _, ok := resolved[parts[0]]; !ok {
			resolved[parts[0]] = Permissions{}
		}

		for _, r := range parts[1:] {
			resolved[parts[0]][r] = true
		}
	}

	return resolved
}

// ExtractRestrictions extracts the eventual authz restrictions embded in the token.
func ExtractRestrictions(token string) (ns string, perms []string, networks []string, err error) {

	claims, err := UnsecureClaimsMap(token)
	if err != nil {
		return "", nil, nil, err
	}

	restrictions, ok := claims["restrictions"].(map[string]interface{})
	if !ok {
		return "", nil, nil, nil
	}

	lns, ok := restrictions["namespace"]
	if ok {
		ns, ok = lns.(string)
		if !ok {
			return "", nil, nil, fmt.Errorf("invalid restrictions.namespace claim type")
		}
	}

	lai, ok := restrictions["perms"]
	if ok {
		permsIface, ok := lai.([]interface{})
		if !ok {
			return "", nil, nil, fmt.Errorf("invalid restrictions.permissions claim type")
		}

		for _, perm := range permsIface {
			pstr, ok := perm.(string)
			if !ok {
				return "", nil, nil, fmt.Errorf("invalid restrictions.permissions claim item type")
			}
			perms = append(perms, pstr)
		}
	}

	lnet, ok := restrictions["networks"]
	if ok {
		lnetIface, ok := lnet.([]interface{})
		if !ok {
			return "", nil, nil, fmt.Errorf("invalid restrictions.networks claim type")
		}

		for _, net := range lnetIface {
			nstr, ok := net.(string)
			if !ok {
				return "", nil, nil, fmt.Errorf("invalid restrictions.networks claim item type")
			}
			networks = append(networks, nstr)
		}
	}

	return ns, perms, networks, nil
}

// UnsecureClaimsMap decodes the claims in the given JWT token without
// verifying its validity. Only use or trust this after proper validation.
func UnsecureClaimsMap(token string) (claims map[string]interface{}, err error) {

	if token == "" {
		return nil, errors.New("invalid jwt: empty")
	}

	parts := strings.Split(token, ".")
	if len(parts) != 3 {
		return nil, errors.New("invalid jwt: not enough segments")
	}

	data, err := jwt.DecodeSegment(parts[1])
	if err != nil {
		return nil, fmt.Errorf("invalid jwt: %s", err)
	}

	claims = map[string]interface{}{}
	if err := json.Unmarshal(data, &claims); err != nil {
		return nil, fmt.Errorf("invalid jwt: %s", err)
	}

	return claims, nil
}
