package processors

import (
	"go.aporeto.io/a3s/pkgs/api"
	"go.aporeto.io/a3s/pkgs/bearermanip"
	"go.aporeto.io/a3s/pkgs/importing"
	"go.aporeto.io/bahamut"
	"go.aporeto.io/elemental"
)

// A HTTPSourcesProcessor is a bahamut processor for HTTPSource.
type ImportProcessor struct {
	bmanipMaker bearermanip.MakerFunc
}

// NewImportProcessor returns a new ImportProcessor .
func NewImportProcessor(bmanipMaker bearermanip.MakerFunc) *ImportProcessor {
	return &ImportProcessor{
		bmanipMaker: bmanipMaker,
	}
}

// ProcessCreate handles the creates requests for HTTPSource.
func (p *ImportProcessor) ProcessCreate(bctx bahamut.Context) error {

	req := bctx.InputData().(*api.Import)
	ns := bctx.Request().Namespace

	for _, lst := range []elemental.Identifiables{
		req.LDAPSources,
		req.OIDCSources,
		req.A3SSources,
		req.MTLSSources,
		req.HTTPSources,
		req.Authorizations,
	} {
		if err := importing.Import(
			bctx.Context(),
			api.Manager(),
			p.bmanipMaker(bctx),
			ns,
			req.Label,
			lst,
			req.Mode == api.ImportModeRemove,
		); err != nil {
			return err
		}
	}

	return nil
}
