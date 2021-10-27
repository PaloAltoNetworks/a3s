package processors

import (
	"net/http"
	"strings"

	"go.aporeto.io/a3s/pkgs/api"
	"go.aporeto.io/a3s/pkgs/crud"
	"go.aporeto.io/bahamut"
	"go.aporeto.io/elemental"
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
func (p *NamespacesProcessor) ProcessCreate(ctx bahamut.Context) error {

	ns := ctx.InputData().(*api.Namespace)

	if strings.Contains(ns.Name, "/") {
		return elemental.NewError("Validation Error", "Namespace name must not contain any / when created", "a3s:policy", http.StatusUnprocessableEntity)
	}

	ns.Namespace = ctx.Request().Namespace
	ns.Name = strings.Join([]string{ns.Namespace, ns.Name}, "/")

	return p.manipulator.Create(manipulate.NewContext(ctx.Context()), ns)
}

// ProcessRetrieveMany handles the retrieve many requests for Namespaces.
func (p *NamespacesProcessor) ProcessRetrieveMany(ctx bahamut.Context) error {

	mctx, err := crud.TranslateContext(ctx)
	if err != nil {
		return err
	}

	nss := api.NamespacesList{}
	if err := p.manipulator.RetrieveMany(mctx, &nss); err != nil {
		return err
	}

	ctx.SetOutputData(nss)

	return nil
}

// ProcessRetrieve handles the retrieve requests for Namespaces.
func (p *NamespacesProcessor) ProcessRetrieve(ctx bahamut.Context) error {

	mctx, err := crud.TranslateContext(ctx)
	if err != nil {
		return err
	}

	ns := api.NewNamespace()
	ns.ID = ctx.Request().ObjectID
	if err := p.manipulator.Retrieve(mctx, ns); err != nil {
		return err
	}

	ctx.SetOutputData(ns)

	return nil
}

// ProcessUpdate handles the update requests for Namespaces.
func (p *NamespacesProcessor) ProcessUpdate(ctx bahamut.Context) error {

	ns := ctx.InputData().(*api.Namespace)
	ns.ID = ctx.Request().ObjectID

	mctx, err := crud.TranslateContext(ctx)
	if err != nil {
		return err
	}

	ens := api.NewNamespace()
	ens.ID = ns.ID
	if err := p.manipulator.Retrieve(mctx, ens); err != nil {
		return elemental.NewError(
			"Not Found",
			"Object not found",
			"a3s:policy",
			http.StatusNotFound,
		)
	}

	elemental.BackportUnexposedFields(ens, ns)

	if err := p.manipulator.Update(mctx, ns); err != nil {
		return err
	}

	ctx.SetOutputData(ns)

	return nil
}

// ProcessDelete handles the delete requests for Namespaces.
func (p *NamespacesProcessor) ProcessDelete(ctx bahamut.Context) error {

	mctx, err := crud.TranslateContext(ctx)
	if err != nil {
		return err
	}

	ns := api.NewNamespace()
	ns.ID = ctx.Request().ObjectID
	if err := p.manipulator.Retrieve(mctx, ns); err != nil {
		return err
	}

	ctx.SetOutputData(ns)

	return p.manipulator.Delete(mctx, ns)
}

// ProcessInfo handles the info request for Namespaces.
func (p *NamespacesProcessor) ProcessInfo(ctx bahamut.Context) error {

	mctx, err := crud.TranslateContext(ctx)
	if err != nil {
		return err
	}

	c, err := p.manipulator.Count(mctx, api.NamespaceIdentity)
	if err != nil {
		return err
	}

	ctx.SetCount(c)

	return nil
}
