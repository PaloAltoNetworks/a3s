package api

import "go.aporeto.io/elemental"

var (
	identityNamesMap = map[string]elemental.Identity{
		"a3ssource":     A3SSourceIdentity,
		"authorization": AuthorizationIdentity,
		"authz":         AuthzIdentity,

		"issue": IssueIdentity,

		"ldapsource":  LDAPSourceIdentity,
		"mtlssource":  MTLSSourceIdentity,
		"namespace":   NamespaceIdentity,
		"oidcsource":  OIDCSourceIdentity,
		"permissions": PermissionsIdentity,
		"root":        RootIdentity,
	}

	identitycategoriesMap = map[string]elemental.Identity{
		"a3ssources":     A3SSourceIdentity,
		"authorizations": AuthorizationIdentity,
		"authz":          AuthzIdentity,

		"issue": IssueIdentity,

		"ldapsources": LDAPSourceIdentity,
		"mtlssources": MTLSSourceIdentity,
		"namespaces":  NamespaceIdentity,
		"oidcsources": OIDCSourceIdentity,
		"permissions": PermissionsIdentity,
		"root":        RootIdentity,
	}

	aliasesMap = map[string]elemental.Identity{}

	indexesMap = map[string][][]string{
		"a3ssource": {
			{"namespace", "name"},
			{":shard", ":unique", "zone", "zHash"},
			{"namespace"},
			{"namespace", "ID"},
		},
		"authorization": {
			{"namespace", "flattenedSubject", "disabled"},
			{"namespace", "flattenedSubject", "propagate"},
			{":shard", ":unique", "zone", "zHash"},
			{"namespace"},
			{"namespace", "ID"},
		},
		"authz": nil,
		"issue": nil,
		"ldapsource": {
			{"namespace", "name"},
			{":shard", ":unique", "zone", "zHash"},
			{"namespace"},
			{"namespace", "ID"},
		},
		"mtlssource": {
			{"namespace", "name"},
			{":shard", ":unique", "zone", "zHash"},
			{"namespace"},
			{"namespace", "ID"},
		},
		"namespace": {
			{"namespace", "name"},
			{":shard", ":unique", "zone", "zHash"},
			{"namespace"},
			{"namespace", "ID"},
			{"name"},
		},
		"oidcsource": {
			{"namespace", "name"},
			{":shard", ":unique", "zone", "zHash"},
			{"namespace"},
			{"namespace", "ID"},
		},
		"permissions": nil,
		"root":        nil,
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

	case A3SSourceIdentity:
		return NewA3SSource()
	case AuthorizationIdentity:
		return NewAuthorization()
	case AuthzIdentity:
		return NewAuthz()
	case IssueIdentity:
		return NewIssue()
	case LDAPSourceIdentity:
		return NewLDAPSource()
	case MTLSSourceIdentity:
		return NewMTLSSource()
	case NamespaceIdentity:
		return NewNamespace()
	case OIDCSourceIdentity:
		return NewOIDCSource()
	case PermissionsIdentity:
		return NewPermissions()
	case RootIdentity:
		return NewRoot()
	default:
		return nil
	}
}

func (f modelManager) SparseIdentifiable(identity elemental.Identity) elemental.SparseIdentifiable {

	switch identity {

	case A3SSourceIdentity:
		return NewSparseA3SSource()
	case AuthorizationIdentity:
		return NewSparseAuthorization()
	case AuthzIdentity:
		return NewSparseAuthz()
	case IssueIdentity:
		return NewSparseIssue()
	case LDAPSourceIdentity:
		return NewSparseLDAPSource()
	case MTLSSourceIdentity:
		return NewSparseMTLSSource()
	case NamespaceIdentity:
		return NewSparseNamespace()
	case OIDCSourceIdentity:
		return NewSparseOIDCSource()
	case PermissionsIdentity:
		return NewSparsePermissions()
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

	case A3SSourceIdentity:
		return &A3SSourcesList{}
	case AuthorizationIdentity:
		return &AuthorizationsList{}
	case AuthzIdentity:
		return &AuthzsList{}
	case IssueIdentity:
		return &IssuesList{}
	case LDAPSourceIdentity:
		return &LDAPSourcesList{}
	case MTLSSourceIdentity:
		return &MTLSSourcesList{}
	case NamespaceIdentity:
		return &NamespacesList{}
	case OIDCSourceIdentity:
		return &OIDCSourcesList{}
	case PermissionsIdentity:
		return &PermissionsList{}
	default:
		return nil
	}
}

func (f modelManager) SparseIdentifiables(identity elemental.Identity) elemental.SparseIdentifiables {

	switch identity {

	case A3SSourceIdentity:
		return &SparseA3SSourcesList{}
	case AuthorizationIdentity:
		return &SparseAuthorizationsList{}
	case AuthzIdentity:
		return &SparseAuthzsList{}
	case IssueIdentity:
		return &SparseIssuesList{}
	case LDAPSourceIdentity:
		return &SparseLDAPSourcesList{}
	case MTLSSourceIdentity:
		return &SparseMTLSSourcesList{}
	case NamespaceIdentity:
		return &SparseNamespacesList{}
	case OIDCSourceIdentity:
		return &SparseOIDCSourcesList{}
	case PermissionsIdentity:
		return &SparsePermissionsList{}
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
		A3SSourceIdentity,
		AuthorizationIdentity,
		AuthzIdentity,
		IssueIdentity,
		LDAPSourceIdentity,
		MTLSSourceIdentity,
		NamespaceIdentity,
		OIDCSourceIdentity,
		PermissionsIdentity,
		RootIdentity,
	}
}

// AliasesForIdentity returns all the aliases for the given identity.
func AliasesForIdentity(identity elemental.Identity) []string {

	switch identity {
	case A3SSourceIdentity:
		return []string{}
	case AuthorizationIdentity:
		return []string{}
	case AuthzIdentity:
		return []string{}
	case IssueIdentity:
		return []string{}
	case LDAPSourceIdentity:
		return []string{}
	case MTLSSourceIdentity:
		return []string{}
	case NamespaceIdentity:
		return []string{}
	case OIDCSourceIdentity:
		return []string{}
	case PermissionsIdentity:
		return []string{}
	case RootIdentity:
		return []string{}
	}

	return nil
}
