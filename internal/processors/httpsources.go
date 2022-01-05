package processors

import (
	"go.aporeto.io/a3s/pkgs/api"
	"go.aporeto.io/a3s/pkgs/crud"
	"go.aporeto.io/bahamut"
	"go.aporeto.io/manipulate"
)

// A HTTPSourcesProcessor is a bahamut processor for HTTPSource.
type HTTPSourcesProcessor struct {
	manipulator manipulate.Manipulator
}

// NewHTTPSourcesProcessor returns a new HTTPSourcesProcessor.
func NewHTTPSourcesProcessor(manipulator manipulate.Manipulator) *HTTPSourcesProcessor {
	return &HTTPSourcesProcessor{
		manipulator: manipulator,
	}
}

// ProcessCreate handles the creates requests for HTTPSource.
func (p *HTTPSourcesProcessor) ProcessCreate(bctx bahamut.Context) error {
	return crud.Create(bctx, p.manipulator, bctx.InputData().(*api.HTTPSource))
}

// ProcessRetrieveMany handles the retrieve many requests for HTTPSource.
func (p *HTTPSourcesProcessor) ProcessRetrieveMany(bctx bahamut.Context) error {
	return crud.RetrieveMany(bctx, p.manipulator, &api.HTTPSourcesList{})
}

// ProcessRetrieve handles the retrieve requests for HTTPSource.
func (p *HTTPSourcesProcessor) ProcessRetrieve(bctx bahamut.Context) error {
	return crud.Retrieve(bctx, p.manipulator, api.NewHTTPSource())
}

// ProcessUpdate handles the update requests for HTTPSource.
func (p *HTTPSourcesProcessor) ProcessUpdate(bctx bahamut.Context) error {
	return crud.Update(bctx, p.manipulator, bctx.InputData().(*api.HTTPSource))
}

// ProcessDelete handles the delete requests for HTTPSource.
func (p *HTTPSourcesProcessor) ProcessDelete(bctx bahamut.Context) error {
	return crud.Delete(bctx, p.manipulator, api.NewHTTPSource())
}

// ProcessInfo handles the info request for HTTPSource.
func (p *HTTPSourcesProcessor) ProcessInfo(bctx bahamut.Context) error {
	return crud.Info(bctx, p.manipulator, api.HTTPSourceIdentity)
}
