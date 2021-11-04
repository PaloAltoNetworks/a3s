package api

import "go.aporeto.io/elemental"

var relationshipsRegistry elemental.RelationshipsRegistry

func init() {

	relationshipsRegistry = elemental.RelationshipsRegistry{}

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
			"root": {},
		},
		Info: map[string]*elemental.RelationshipInfo{
			"root": {},
		},
	}

	relationshipsRegistry[IssueIdentity] = &elemental.Relationship{
		Create: map[string]*elemental.RelationshipInfo{
			"root": {
				Parameters: []elemental.ParameterDefinition{
					{
						Name: "asCookie",
						Type: "boolean",
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
			"root": {},
		},
		Info: map[string]*elemental.RelationshipInfo{
			"root": {},
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
			"root": {},
		},
		Info: map[string]*elemental.RelationshipInfo{
			"root": {},
		},
	}

	relationshipsRegistry[PermissionsIdentity] = &elemental.Relationship{
		Create: map[string]*elemental.RelationshipInfo{
			"root": {},
		},
	}

	relationshipsRegistry[RootIdentity] = &elemental.Relationship{}

}
