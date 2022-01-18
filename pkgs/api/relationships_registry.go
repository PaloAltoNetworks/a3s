package api

import "go.aporeto.io/elemental"

var relationshipsRegistry elemental.RelationshipsRegistry

func init() {

	relationshipsRegistry = elemental.RelationshipsRegistry{}

	relationshipsRegistry[A3SSourceIdentity] = &elemental.Relationship{
		Create: map[string]*elemental.RelationshipInfo{
			"root": {},
		},
		Update: map[string]*elemental.RelationshipInfo{
			"root": {},
		},
		Patch: map[string]*elemental.RelationshipInfo{
			"root": {},
		},
		Delete: map[string]*elemental.RelationshipInfo{
			"root": {},
		},
		Retrieve: map[string]*elemental.RelationshipInfo{
			"root": {},
		},
		RetrieveMany: map[string]*elemental.RelationshipInfo{
			"root": {
				Parameters: []elemental.ParameterDefinition{
					{
						Name: "q",
						Type: "string",
					},
				},
			},
		},
		Info: map[string]*elemental.RelationshipInfo{
			"root": {
				Parameters: []elemental.ParameterDefinition{
					{
						Name: "q",
						Type: "string",
					},
				},
			},
		},
	}

	relationshipsRegistry[AuthorizationIdentity] = &elemental.Relationship{
		Create: map[string]*elemental.RelationshipInfo{
			"root": {},
		},
		Update: map[string]*elemental.RelationshipInfo{
			"root": {},
		},
		Patch: map[string]*elemental.RelationshipInfo{
			"root": {},
		},
		Delete: map[string]*elemental.RelationshipInfo{
			"root": {
				Parameters: []elemental.ParameterDefinition{
					{
						Name: "q",
						Type: "string",
					},
				},
			},
		},
		Retrieve: map[string]*elemental.RelationshipInfo{
			"root": {},
		},
		RetrieveMany: map[string]*elemental.RelationshipInfo{
			"root": {
				Parameters: []elemental.ParameterDefinition{
					{
						Name: "q",
						Type: "string",
					},
				},
			},
		},
		Info: map[string]*elemental.RelationshipInfo{
			"root": {
				Parameters: []elemental.ParameterDefinition{
					{
						Name: "q",
						Type: "string",
					},
				},
			},
		},
	}

	relationshipsRegistry[AuthzIdentity] = &elemental.Relationship{
		Create: map[string]*elemental.RelationshipInfo{
			"root": {},
		},
	}

	relationshipsRegistry[HTTPSourceIdentity] = &elemental.Relationship{
		Create: map[string]*elemental.RelationshipInfo{
			"root": {},
		},
		Update: map[string]*elemental.RelationshipInfo{
			"root": {},
		},
		Patch: map[string]*elemental.RelationshipInfo{
			"root": {},
		},
		Delete: map[string]*elemental.RelationshipInfo{
			"root": {},
		},
		Retrieve: map[string]*elemental.RelationshipInfo{
			"root": {},
		},
		RetrieveMany: map[string]*elemental.RelationshipInfo{
			"root": {
				Parameters: []elemental.ParameterDefinition{
					{
						Name: "q",
						Type: "string",
					},
				},
			},
		},
		Info: map[string]*elemental.RelationshipInfo{
			"root": {
				Parameters: []elemental.ParameterDefinition{
					{
						Name: "q",
						Type: "string",
					},
				},
			},
		},
	}

	relationshipsRegistry[IdentityModifierIdentity] = &elemental.Relationship{}

	relationshipsRegistry[ImportIdentity] = &elemental.Relationship{
		Create: map[string]*elemental.RelationshipInfo{
			"root": {},
		},
	}

	relationshipsRegistry[IssueIdentity] = &elemental.Relationship{
		Create: map[string]*elemental.RelationshipInfo{
			"root": {},
		},
	}

	relationshipsRegistry[LDAPSourceIdentity] = &elemental.Relationship{
		Create: map[string]*elemental.RelationshipInfo{
			"root": {},
		},
		Update: map[string]*elemental.RelationshipInfo{
			"root": {},
		},
		Patch: map[string]*elemental.RelationshipInfo{
			"root": {},
		},
		Delete: map[string]*elemental.RelationshipInfo{
			"root": {},
		},
		Retrieve: map[string]*elemental.RelationshipInfo{
			"root": {},
		},
		RetrieveMany: map[string]*elemental.RelationshipInfo{
			"root": {
				Parameters: []elemental.ParameterDefinition{
					{
						Name: "q",
						Type: "string",
					},
				},
			},
		},
		Info: map[string]*elemental.RelationshipInfo{
			"root": {
				Parameters: []elemental.ParameterDefinition{
					{
						Name: "q",
						Type: "string",
					},
				},
			},
		},
	}

	relationshipsRegistry[MTLSSourceIdentity] = &elemental.Relationship{
		Create: map[string]*elemental.RelationshipInfo{
			"root": {},
		},
		Update: map[string]*elemental.RelationshipInfo{
			"root": {},
		},
		Patch: map[string]*elemental.RelationshipInfo{
			"root": {},
		},
		Delete: map[string]*elemental.RelationshipInfo{
			"root": {},
		},
		Retrieve: map[string]*elemental.RelationshipInfo{
			"root": {},
		},
		RetrieveMany: map[string]*elemental.RelationshipInfo{
			"root": {
				Parameters: []elemental.ParameterDefinition{
					{
						Name: "q",
						Type: "string",
					},
				},
			},
		},
		Info: map[string]*elemental.RelationshipInfo{
			"root": {
				Parameters: []elemental.ParameterDefinition{
					{
						Name: "q",
						Type: "string",
					},
				},
			},
		},
	}

	relationshipsRegistry[NamespaceIdentity] = &elemental.Relationship{
		Create: map[string]*elemental.RelationshipInfo{
			"root": {},
		},
		Update: map[string]*elemental.RelationshipInfo{
			"root": {},
		},
		Patch: map[string]*elemental.RelationshipInfo{
			"root": {},
		},
		Delete: map[string]*elemental.RelationshipInfo{
			"root": {},
		},
		Retrieve: map[string]*elemental.RelationshipInfo{
			"root": {},
		},
		RetrieveMany: map[string]*elemental.RelationshipInfo{
			"root": {
				Parameters: []elemental.ParameterDefinition{
					{
						Name: "q",
						Type: "string",
					},
				},
			},
		},
		Info: map[string]*elemental.RelationshipInfo{
			"root": {
				Parameters: []elemental.ParameterDefinition{
					{
						Name: "q",
						Type: "string",
					},
				},
			},
		},
	}

	relationshipsRegistry[OIDCSourceIdentity] = &elemental.Relationship{
		Create: map[string]*elemental.RelationshipInfo{
			"root": {},
		},
		Update: map[string]*elemental.RelationshipInfo{
			"root": {},
		},
		Patch: map[string]*elemental.RelationshipInfo{
			"root": {},
		},
		Delete: map[string]*elemental.RelationshipInfo{
			"root": {},
		},
		Retrieve: map[string]*elemental.RelationshipInfo{
			"root": {},
		},
		RetrieveMany: map[string]*elemental.RelationshipInfo{
			"root": {
				Parameters: []elemental.ParameterDefinition{
					{
						Name: "q",
						Type: "string",
					},
				},
			},
		},
		Info: map[string]*elemental.RelationshipInfo{
			"root": {
				Parameters: []elemental.ParameterDefinition{
					{
						Name: "q",
						Type: "string",
					},
				},
			},
		},
	}

	relationshipsRegistry[PermissionsIdentity] = &elemental.Relationship{
		Create: map[string]*elemental.RelationshipInfo{
			"root": {},
		},
	}

	relationshipsRegistry[RootIdentity] = &elemental.Relationship{}

}
