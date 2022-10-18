package permissions

// A Transformer is an object that can manipulate a
// permissions map based on a provided mapping.
type Transformer interface {
	Transform(permissions PermissionMap) PermissionMap
}

type transformer struct {
	roleExpander map[string]map[string][]string
}

// NewTransformer returns a new Transformer.
func NewTransformer(roleExpander map[string]map[string][]string) Transformer {
	return &transformer{
		roleExpander: roleExpander,
	}
}

// Transform modifies the provided permissions map based
// on a given mapping. This can be used to expand roles
// into their lower-level permissions.
func (t *transformer) Transform(permissions PermissionMap) PermissionMap {

	if t.roleExpander == nil {
		return permissions
	}

	for name, auths := range t.roleExpander {
		if _, ok := permissions[name]; !ok {
			continue
		}

		delete(permissions, name)

		for identity, verbs := range auths {
			for _, verb := range verbs {
				permissions[identity][verb] = true
			}
		}
	}

	return permissions
}
