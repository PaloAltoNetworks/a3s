package api

import "go.aporeto.io/elemental"

var (
	identityNamesMap = map[string]elemental.Identity{
		"authorization": AuthorizationIdentity,
		"issue":         IssueIdentity,
		"namespace":     NamespaceIdentity,
		"role":          RoleIdentity,
		"root":          RootIdentity,
	}

	identitycategoriesMap = map[string]elemental.Identity{
		"authorizations": AuthorizationIdentity,
		"issue":          IssueIdentity,
		"namespaces":     NamespaceIdentity,
		"roles":          RoleIdentity,
		"root":           RootIdentity,
	}

	aliasesMap = map[string]elemental.Identity{}

	indexesMap = map[string][][]string{
		"authorization": {
			{"namespace", "flattenedSubject", "disabled"},
			{"namespace", "flattenedSubject", "propagate"},
			{"namespace"},
			{"namespace", "ID"},
		},
		"issue": nil,
		"namespace": {
			{"namespace", "name"},
			{"namespace"},
			{"namespace", "ID"},
			{"name"},
		},
		"role": nil,
		"root": nil,
	}
)

// ModelVersion returns the current version of the model.
func ModelVersion() float64 { return 1 }

type modelManager struct{}

func (f modelManager) IdentityFromName(name string) elemental.Identity {

	return identityNamesMap[name]
}

func (f modelManager) IdentityFromCategory(category string) elemental.Identity {

	return identitycategoriesMap[category]
}

func (f modelManager) IdentityFromAlias(alias string) elemental.Identity {

	return aliasesMap[alias]
}

func (f modelManager) IdentityFromAny(any string) (i elemental.Identity) {

	if i = f.IdentityFromName(any); !i.IsEmpty() {
		return i
	}

	if i = f.IdentityFromCategory(any); !i.IsEmpty() {
		return i
	}

	return f.IdentityFromAlias(any)
}

func (f modelManager) Identifiable(identity elemental.Identity) elemental.Identifiable {

	switch identity {

	case AuthorizationIdentity:
		return NewAuthorization()
	case IssueIdentity:
		return NewIssue()
	case NamespaceIdentity:
		return NewNamespace()
	case RoleIdentity:
		return NewRole()
	case RootIdentity:
		return NewRoot()
	default:
		return nil
	}
}

func (f modelManager) SparseIdentifiable(identity elemental.Identity) elemental.SparseIdentifiable {

	switch identity {

	case AuthorizationIdentity:
		return NewSparseAuthorization()
	case IssueIdentity:
		return NewSparseIssue()
	case NamespaceIdentity:
		return NewSparseNamespace()
	case RoleIdentity:
		return NewSparseRole()
	default:
		return nil
	}
}

func (f modelManager) Indexes(identity elemental.Identity) [][]string {

	return indexesMap[identity.Name]
}

func (f modelManager) IdentifiableFromString(any string) elemental.Identifiable {

	return f.Identifiable(f.IdentityFromAny(any))
}

func (f modelManager) Identifiables(identity elemental.Identity) elemental.Identifiables {

	switch identity {

	case AuthorizationIdentity:
		return &AuthorizationsList{}
	case IssueIdentity:
		return &IssuesList{}
	case NamespaceIdentity:
		return &NamespacesList{}
	case RoleIdentity:
		return &RolesList{}
	default:
		return nil
	}
}

func (f modelManager) SparseIdentifiables(identity elemental.Identity) elemental.SparseIdentifiables {

	switch identity {

	case AuthorizationIdentity:
		return &SparseAuthorizationsList{}
	case IssueIdentity:
		return &SparseIssuesList{}
	case NamespaceIdentity:
		return &SparseNamespacesList{}
	case RoleIdentity:
		return &SparseRolesList{}
	default:
		return nil
	}
}

func (f modelManager) IdentifiablesFromString(any string) elemental.Identifiables {

	return f.Identifiables(f.IdentityFromAny(any))
}

func (f modelManager) Relationships() elemental.RelationshipsRegistry {

	return relationshipsRegistry
}

func (f modelManager) AllIdentities() []elemental.Identity {
	return AllIdentities()
}

var manager = modelManager{}

// Manager returns the model elemental.ModelManager.
func Manager() elemental.ModelManager { return manager }

// AllIdentities returns all existing identities.
func AllIdentities() []elemental.Identity {

	return []elemental.Identity{
		AuthorizationIdentity,
		IssueIdentity,
		NamespaceIdentity,
		RoleIdentity,
		RootIdentity,
	}
}

// AliasesForIdentity returns all the aliases for the given identity.
func AliasesForIdentity(identity elemental.Identity) []string {

	switch identity {
	case AuthorizationIdentity:
		return []string{}
	case IssueIdentity:
		return []string{}
	case NamespaceIdentity:
		return []string{}
	case RoleIdentity:
		return []string{}
	case RootIdentity:
		return []string{}
	}

	return nil
}
