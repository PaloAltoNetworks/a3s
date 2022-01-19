package processors

import (
	"fmt"
	"net/http"

	"go.aporeto.io/a3s/pkgs/api"
	"go.aporeto.io/a3s/pkgs/authorizer"
	"go.aporeto.io/a3s/pkgs/bearermanip"
	"go.aporeto.io/a3s/pkgs/importing"
	"go.aporeto.io/a3s/pkgs/permissions"
	"go.aporeto.io/a3s/pkgs/token"
	"go.aporeto.io/bahamut"
	"go.aporeto.io/elemental"
)

// A ImportProcessor is a bahamut processor for Import.
type ImportProcessor struct {
	bmanipMaker bearermanip.MakerFunc
	authz       authorizer.Authorizer
}

// NewImportProcessor returns a new ImportProcessor .
func NewImportProcessor(bmanipMaker bearermanip.MakerFunc, authz authorizer.Authorizer) *ImportProcessor {
	return &ImportProcessor{
		bmanipMaker: bmanipMaker,
		authz:       authz,
	}
}

// ProcessCreate handles the creates requests for HTTPSource.
func (p *ImportProcessor) ProcessCreate(bctx bahamut.Context) error {

	req := bctx.InputData().(*api.Import)
	ns := bctx.Request().Namespace

	values := []elemental.Identifiables{
		req.LDAPSources,
		req.OIDCSources,
		req.A3SSources,
		req.MTLSSources,
		req.HTTPSources,
		req.Authorizations,
	}

	restrictions, err := permissions.GetRestrictions(token.FromRequest(bctx.Request()))
	if err != nil {
		return err
	}

	for _, lst := range values {

		if len(lst.List()) == 0 {
			continue
		}

		for _, perm := range []string{"retrieve-many", "create", "delete"} {
			ok, err := p.authz.CheckAuthorization(
				bctx.Context(),
				bctx.Claims(),
				perm,
				ns,
				lst.Identity().Category,
				authorizer.OptionCheckRestrictions(restrictions),
				authorizer.OptionCheckSourceIP(bctx.Request().ClientIP),
			)
			if err != nil {
				return err
			}
			if !ok {
				return elemental.NewError(
					"Permission Denied",
					fmt.Sprintf("You don't have the permission to '%s' on '%s'", perm, lst.Identity().Category),
					"a3s:import",
					http.StatusForbidden,
				)
			}
		}
	}

	for _, lst := range values {

		if len(lst.List()) == 0 {
			continue
		}

		if err := importing.Import(
			bctx.Context(),
			api.Manager(),
			p.bmanipMaker(bctx),
			ns,
			req.Label,
			lst,
			bctx.Request().Parameters.Get("delete").BoolValue(),
		); err != nil {
			return err
		}
	}

	return nil
}
