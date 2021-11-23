package processors

import (
	"go.aporeto.io/a3s/pkgs/api"
	"go.aporeto.io/a3s/pkgs/crud"
	"go.aporeto.io/bahamut"
	"go.aporeto.io/manipulate"
)

// A A3SSourcesProcessor is a bahamut processor for A3SSource.
type A3SSourcesProcessor struct {
	manipulator manipulate.Manipulator
}

// NewA3SSourcesProcessor returns a new A3SSourcesProcessor.
func NewA3SSourcesProcessor(manipulator manipulate.Manipulator) *A3SSourcesProcessor {
	return &A3SSourcesProcessor{
		manipulator: manipulator,
	}
}

// ProcessCreate handles the creates requests for A3SSource.
func (p *A3SSourcesProcessor) ProcessCreate(bctx bahamut.Context) error {
	return crud.Create(bctx, p.manipulator, bctx.InputData().(*api.A3SSource))
}

// ProcessRetrieveMany handles the retrieve many requests for A3SSource.
func (p *A3SSourcesProcessor) ProcessRetrieveMany(bctx bahamut.Context) error {
	return crud.RetrieveMany(bctx, p.manipulator, &api.A3SSourcesList{})
}

// ProcessRetrieve handles the retrieve requests for A3SSource.
func (p *A3SSourcesProcessor) ProcessRetrieve(bctx bahamut.Context) error {
	return crud.Retrieve(bctx, p.manipulator, api.NewA3SSource())
}

// ProcessUpdate handles the update requests for A3SSource.
func (p *A3SSourcesProcessor) ProcessUpdate(bctx bahamut.Context) error {
	return crud.Update(bctx, p.manipulator, bctx.InputData().(*api.A3SSource))
}

// ProcessDelete handles the delete requests for A3SSource.
func (p *A3SSourcesProcessor) ProcessDelete(bctx bahamut.Context) error {
	return crud.Delete(bctx, p.manipulator, api.NewA3SSource())
}

// ProcessInfo handles the info request for A3SSource.
func (p *A3SSourcesProcessor) ProcessInfo(bctx bahamut.Context) error {
	return crud.Info(bctx, p.manipulator, api.A3SSourceIdentity)
}
