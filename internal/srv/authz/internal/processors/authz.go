package processors

import (
	"go.aporeto.io/a3s/pkgs/api"
	"go.aporeto.io/a3s/pkgs/permissions"
	"go.aporeto.io/bahamut"
)

// A AuthzProcessor is a bahamut processor for Authzs.
type AuthzProcessor struct {
	retriever permissions.Retriever
}

// NewAuthzProcessor returns a new AuthzProcessor.
func NewAuthzProcessor(retriever permissions.Retriever) *AuthzProcessor {
	return &AuthzProcessor{
		retriever: retriever,
	}
}

// ProcessCreate handles the creates requests for Authzs.
func (p *AuthzProcessor) ProcessCreate(bctx bahamut.Context) error {

	req := bctx.InputData().(*api.Authz)

	restrictions := permissions.Restrictions{
		Namespace:   req.RestrictedNamespace,
		Networks:    req.RestrictedNetworks,
		Permissions: req.RestrictedPermissions,
	}

	perms, err := p.retriever.Permissions(
		bctx.Context(),
		req.Claims,
		req.TargetNamespace,
		permissions.OptionRetrieverID(req.TargetID),
		permissions.OptionRetrieverSourceIP(req.ClientIP),
		permissions.OptionRetrieverRestrictions(restrictions),
	)

	switch err {
	case nil:
		req.Permissions = permsToMap(perms)
	default:
		req.Error = err.Error()
	}

	bctx.SetOutputData(req)

	return nil
}

func permsToMap(p permissions.PermissionMap) map[string]map[string]bool {

	out := make(map[string]map[string]bool, len(p))

	for resource, perms := range p {
		out[resource] = make(map[string]bool, len(perms))
		for action, allowed := range perms {
			out[resource][action] = allowed
		}
	}

	return out
}
