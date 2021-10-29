package permissions

import (
	"fmt"
	"net/http"
	"strings"

	"go.aporeto.io/elemental"
)

// Copy returns a copy of the perms
func Copy(perms map[string]map[string]bool) map[string]map[string]bool {

	var copy = make(map[string]map[string]bool, len(perms))

	for i, m := range perms {
		copy[i] = make(map[string]bool, len(m))
		for k, v := range m {
			copy[i][k] = v
		}
	}
	return copy
}

// Contains returns true if the given Authorization is equal or lesser
// than the receiver.
func Contains(perms map[string]map[string]bool, other map[string]map[string]bool) bool {

	if len(perms) == 0 {
		return false
	}

	star := perms["*"]

	for identity, decorators := range other {

		if _, ok := perms[identity]; !ok && len(star) == 0 {
			return false
		}

		for decorator := range decorators {
			if !perms[identity][decorator] && !star[decorator] {
				ok1 := perms[identity]["*"]
				ok2 := star["*"]
				if !ok1 && !ok2 {
					return false
				}
			}
		}
	}

	return true
}

// Intersect returns the intersection between first set and second set.
func Intersect(base map[string]map[string]bool, other map[string]map[string]bool) map[string]map[string]bool {

	// If one or the other are empty, the intersection is nil.
	if len(base) == 0 || len(other) == 0 {
		return map[string]map[string]bool{}
	}

	// first we copy the base, since we are going to
	// modify it.
	candidate := map[string]map[string]bool{}
	for k, v := range base {
		candidate[k] = map[string]bool{}
		for kk, vv := range v {
			candidate[k][kk] = vv
		}
	}

	// If the candidate has a * in it,
	// we copy all the other's key in the base map
	// that are not already there
	if _, ok := candidate["*"]; ok {
		delete(candidate, "*")
		for k, v := range other {
			if _, ok := candidate[k]; !ok {
				candidate[k] = map[string]bool{}
				for kk, vv := range v {
					candidate[k][kk] = vv
				}
			}
		}
	}

	// If the other as a star, we keep track of
	// the general permissions
	rstartperms, rstartok := other["*"]

	// now we loop on all the permission of the out candidate
	for identity, perms := range candidate {

		// Otherwise we check check if the other
		// has the identity declared.
		rperms, ok := other[identity]

		// If it does not, and we have no * declared
		// we remove the identity from the candidate
		// and continue
		if !ok && !rstartok {
			delete(candidate, identity)
			continue
		}

		// We may have nil perms here in case
		// of no identity, but global permissions
		// so we eventually initialize the map.
		if rperms == nil {
			rperms = map[string]bool{}
		}

		// If we have some global perms we backport them
		// to the current set of perms.
		if rstartok {
			for k, v := range rstartperms {
				rperms[k] = v
			}
		}

		// We now check if the candidate permissions of the
		// current identity is *. If it is,
		// then we simply apply the other permissions.
		// and we continue
		if allowed, ok := perms["*"]; ok && allowed {
			candidate[identity] = rperms
			continue
		}

		// Otherwise we loop of the candidate perms.
		for perm := range perms {

			// If the restricted permissions are not here and there is
			// no * declared, we remove the permission from the candidate.
			allowed, ok := rperms[perm]
			allowedAny, okAny := rperms["*"]
			if (!ok || !allowed) && (!okAny || !allowedAny) {
				delete(perms, perm)
			}
		}
	}

	return candidate
}

// IsAllowed returns true if the given operation on the given identity is allowed in the
// given perms.
func IsAllowed(perms map[string]map[string]bool, operation elemental.Operation, identity elemental.Identity) bool {

	method, err := OperationToMethod(operation)
	if err != nil {
		panic(err)
	}

	allowed := func(p map[string]bool, m string) bool {
		if authorized := p["*"]; authorized {
			if authorized {
				return true
			}
		}

		if authorized := p[m]; authorized {
			return true
		}

		return false
	}

	if p, ok := perms["*"]; ok {
		if allowed(p, method) {
			return true
		}
	}

	for i, p := range perms {
		if identity.Name != i {
			continue
		}
		if allowed(p, method) {
			return true
		}
	}

	return false
}

// OperationToMethod is a helper that returns the HTTP method associated with a given elemental.Operation
func OperationToMethod(op elemental.Operation) (string, error) {
	var method string

	switch op {
	case elemental.OperationCreate:
		method = http.MethodPost
	case elemental.OperationDelete:
		method = http.MethodDelete
	case elemental.OperationUpdate, elemental.OperationPatch:
		method = http.MethodPut
	case elemental.OperationRetrieve, elemental.OperationRetrieveMany, elemental.OperationInfo:
		method = http.MethodGet
	default:
		return "", fmt.Errorf("unsupported operation: %s", op)
	}

	return strings.ToLower(method), nil
}
