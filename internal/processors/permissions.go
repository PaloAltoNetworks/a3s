package processors

import (
	"go.aporeto.io/a3s/pkgs/api"
	"go.aporeto.io/a3s/pkgs/permissions"
	"go.aporeto.io/bahamut"
)

// A PermissionsProcessor is a bahamut processor for Permissionss.
type PermissionsProcessor struct {
	retriever permissions.Retriever
}

// NewPermissionsProcessor returns a new PermissionsProcessor.
func NewPermissionsProcessor(retriever permissions.Retriever) *PermissionsProcessor {
	return &PermissionsProcessor{
		retriever: retriever,
	}
}

// ProcessCreate handles the creates requests for Permissionss.
func (p *PermissionsProcessor) ProcessCreate(bctx bahamut.Context) error {

	req := bctx.InputData().(*api.Permissions)

	restrictions := permissions.Restrictions{
		Namespace:   req.RestrictedNamespace,
		Networks:    req.RestrictedNetworks,
		Permissions: req.RestrictedPermissions,
	}

	perms, err := p.retriever.Permissions(
		bctx.Context(),
		req.Claims,
		req.Namespace,
		permissions.OptionRetrieverID(req.ID),
		permissions.OptionRetrieverSourceIP(req.IP),
		permissions.OptionRetrieverRestrictions(restrictions),
		permissions.OptionOffloadRestrictions(req.OffloadRestrictions),
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
