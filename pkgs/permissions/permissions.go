package permissions

import (
	"strings"
)

// Permissions represents a parsed permission string.
type Permissions map[string]bool

// A PermissionMap represents a map of resource to Permissions
type PermissionMap map[string]Permissions

// Parse parses the given list of permission string in the form
// resource,action1,action2:targetID1,targetID2 and returns the
// PermissionMap.
func Parse(authStrings []string, targetID string) PermissionMap {

	auths := PermissionMap{}

	for _, item := range authStrings {

		segments := strings.SplitN(item, ":", 3)

		var resource, actions, ids string
		switch len(segments) {
		case 0, 1:
			continue
		case 2:
			resource = segments[0]
			actions = segments[1]
		case 3:
			resource = segments[0]
			actions = segments[1]
			ids = segments[2]

		}

		if len(ids) > 0 {

			// We did not receive any targetID, so this rule does not apply.
			if targetID == "" {
				continue
			}

			accept := false
			for _, tid := range strings.Split(ids, ",") {
				if tid == targetID {
					accept = true
					break
				}
			}

			if !accept {
				continue
			}
		}

		parts := strings.Split(actions, ",")

		if _, ok := auths[resource]; !ok {
			auths[resource] = map[string]bool{}
		}

		for _, action := range parts {
			auths[resource][action] = true
		}
	}

	return auths
}

// Copy returns a copy of the perms
func (p PermissionMap) Copy() PermissionMap {

	var copy = make(PermissionMap, len(p))

	for i, m := range p {
		copy[i] = make(Permissions, len(m))
		for k, v := range m {
			copy[i][k] = v
		}
	}
	return copy
}

// Contains returns true if the given Authorization is equal or lesser
// than the receiver.
func (p PermissionMap) Contains(other PermissionMap) bool {

	if len(p) == 0 {
		return false
	}

	star := p["*"]

	for identity, decorators := range other {

		if _, ok := p[identity]; !ok && len(star) == 0 {
			return false
		}

		for decorator := range decorators {
			if !p[identity][decorator] && !star[decorator] {
				ok1 := p[identity]["*"]
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
func (p PermissionMap) Intersect(other PermissionMap) PermissionMap {

	// If one or the other are empty, the intersection is nil.
	if len(p) == 0 || len(other) == 0 {
		return PermissionMap{}
	}

	// first we copy the base, since we are going to
	// modify it.
	candidate := PermissionMap{}
	for k, v := range p {
		candidate[k] = Permissions{}
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
				candidate[k] = Permissions{}
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
	for resource, perms := range candidate {

		// Otherwise we check check if the other
		// has the identity declared.
		rperms, ok := other[resource]

		// If it does not, and we have no * declared
		// we remove the identity from the candidate
		// and continue
		if !ok && !rstartok {
			delete(candidate, resource)
			continue
		}

		// We may have nil perms here in case
		// of no identity, but global permissions
		// so we eventually initialize the map.
		if rperms == nil {
			rperms = Permissions{}
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
			candidate[resource] = rperms
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

// Allows returns true if the given operation on the given identity is allowed in the
// given perms.
func (p PermissionMap) Allows(operation string, resource string) bool {

	allowed := func(p Permissions, m string) bool {
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

	if p, ok := p["*"]; ok {
		if allowed(p, operation) {
			return true
		}
	}

	for i, p := range p {
		if resource != i {
			continue
		}
		if allowed(p, operation) {
			return true
		}
	}

	return false
}
