package processors

import (
	"go.aporeto.io/a3s/pkgs/api"
	"go.aporeto.io/a3s/pkgs/crud"
	"go.aporeto.io/bahamut"
	"go.aporeto.io/manipulate"
)

// A NamespacesProcessor is a bahamut processor for Namespaces.
type NamespacesProcessor struct {
	manipulator manipulate.Manipulator
}

// NewNamespacesProcessor returns a new NamespacesProcessor.
func NewNamespacesProcessor(manipulator manipulate.Manipulator) *NamespacesProcessor {
	return &NamespacesProcessor{
		manipulator: manipulator,
	}
}

// ProcessCreate handles the creates requests for Namespaces.
func (p *NamespacesProcessor) ProcessCreate(bctx bahamut.Context) error {
	return crud.Create(bctx, p.manipulator, bctx.InputData().(*api.Namespace))
}

// ProcessRetrieveMany handles the retrieve many requests for Namespaces.
func (p *NamespacesProcessor) ProcessRetrieveMany(bctx bahamut.Context) error {
	return crud.RetrieveMany(bctx, p.manipulator, &api.NamespacesList{})
}

// ProcessRetrieve handles the retrieve requests for Namespaces.
func (p *NamespacesProcessor) ProcessRetrieve(bctx bahamut.Context) error {
	return crud.Retrieve(bctx, p.manipulator, api.NewNamespace())
}

// ProcessUpdate handles the update requests for Namespaces.
func (p *NamespacesProcessor) ProcessUpdate(bctx bahamut.Context) error {
	return crud.Update(bctx, p.manipulator, bctx.InputData().(*api.Namespace))
}

// ProcessDelete handles the delete requests for Namespaces.
func (p *NamespacesProcessor) ProcessDelete(bctx bahamut.Context) error {
	return crud.Delete(bctx, p.manipulator, api.NewNamespace())
}

// ProcessInfo handles the info request for Namespaces.
func (p *NamespacesProcessor) ProcessInfo(bctx bahamut.Context) error {
	return crud.Info(bctx, p.manipulator, api.NamespaceIdentity)
}
