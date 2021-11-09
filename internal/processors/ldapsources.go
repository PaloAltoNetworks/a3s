package processors

import (
	"go.aporeto.io/a3s/pkgs/api"
	"go.aporeto.io/a3s/pkgs/crud"
	"go.aporeto.io/bahamut"
	"go.aporeto.io/manipulate"
)

// A LDAPSourcesProcessor is a bahamut processor for LDAPSource.
type LDAPSourcesProcessor struct {
	manipulator manipulate.Manipulator
}

// NewLDAPSourcesProcessor returns a new LDAPSourcesProcessor.
func NewLDAPSourcesProcessor(manipulator manipulate.Manipulator) *LDAPSourcesProcessor {
	return &LDAPSourcesProcessor{
		manipulator: manipulator,
	}
}

// ProcessCreate handles the creates requests for LDAPSource.
func (p *LDAPSourcesProcessor) ProcessCreate(bctx bahamut.Context) error {
	return crud.Create(bctx, p.manipulator, bctx.InputData().(*api.LDAPSource))
}

// ProcessRetrieveMany handles the retrieve many requests for LDAPSource.
func (p *LDAPSourcesProcessor) ProcessRetrieveMany(bctx bahamut.Context) error {
	return crud.RetrieveMany(bctx, p.manipulator, &api.LDAPSourcesList{})
}

// ProcessRetrieve handles the retrieve requests for LDAPSource.
func (p *LDAPSourcesProcessor) ProcessRetrieve(bctx bahamut.Context) error {
	return crud.Retrieve(bctx, p.manipulator, api.NewLDAPSource())
}

// ProcessUpdate handles the update requests for LDAPSource.
func (p *LDAPSourcesProcessor) ProcessUpdate(bctx bahamut.Context) error {
	return crud.Update(bctx, p.manipulator, bctx.InputData().(*api.LDAPSource))
}

// ProcessDelete handles the delete requests for LDAPSource.
func (p *LDAPSourcesProcessor) ProcessDelete(bctx bahamut.Context) error {
	return crud.Delete(bctx, p.manipulator, api.NewLDAPSource())
}

// ProcessInfo handles the info request for LDAPSource.
func (p *LDAPSourcesProcessor) ProcessInfo(bctx bahamut.Context) error {
	return crud.Info(bctx, p.manipulator, api.LDAPSourceIdentity)
}
