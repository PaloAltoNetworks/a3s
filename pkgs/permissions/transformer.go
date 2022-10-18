package permissions

// A Transformer is an object that can manipulate a
// permissions map based on a provided mapping.
type Transformer interface {
	Transform(permissions []string) []string
}

type transformer struct {
	roleExpander map[string][]string
}

// NewTransformer returns a new Transformer.
func NewTransformer(roleExpander map[string][]string) Transformer {
	return &transformer{
		roleExpander: roleExpander,
	}
}

// Transform modifies the provided permissions map based
// on a given mapping. This can be used to expand roles
// into their lower-level permissions.
func (t *transformer) Transform(permissions []string) []string {

	if t.roleExpander == nil || len(permissions) == 0 {
		return permissions
	}

	perms := map[string]struct{}{}
	for _, p := range permissions {
		perms[p] = struct{}{}
	}

	for name, auths := range t.roleExpander {
		if _, ok := perms[name]; !ok {
			continue
		}

		delete(perms, name)

		for _, auth := range auths {
			perms[auth] = struct{}{}
		}
	}

	newPermissions := make([]string, 0, len(perms))
	for perm := range perms {
		newPermissions = append(newPermissions, perm)
	}

	return newPermissions
}
