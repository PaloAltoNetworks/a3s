package processors

import (
	"go.aporeto.io/a3s/pkgs/api"
	"go.aporeto.io/a3s/pkgs/authorizer"
	"go.aporeto.io/a3s/pkgs/bearermanip"
	"go.aporeto.io/a3s/pkgs/importing"
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

	for _, lst := range values {
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
