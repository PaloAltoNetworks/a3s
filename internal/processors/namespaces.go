package processors

import (
	"net/http"
	"strings"
	"time"

	"go.aporeto.io/a3s/pkgs/api"
	"go.aporeto.io/a3s/pkgs/crud"
	"go.aporeto.io/a3s/pkgs/notification"
	"go.aporeto.io/a3s/pkgs/nscache"
	"go.aporeto.io/bahamut"
	"go.aporeto.io/elemental"
	"go.aporeto.io/manipulate"
	"go.uber.org/zap"
)

// A NamespacesProcessor is a bahamut processor for Namespaces.
type NamespacesProcessor struct {
	manipulator manipulate.Manipulator
	pubsub      bahamut.PubSubClient
}

// NewNamespacesProcessor returns a new NamespacesProcessor.
func NewNamespacesProcessor(manipulator manipulate.Manipulator, pubsub bahamut.PubSubClient) *NamespacesProcessor {
	return &NamespacesProcessor{
		manipulator: manipulator,
		pubsub:      pubsub,
	}
}

// ProcessCreate handles the creates requests for Namespaces.
func (p *NamespacesProcessor) ProcessCreate(bctx bahamut.Context) error {

	ns := bctx.InputData().(*api.Namespace)

	if strings.Contains(ns.Name, "/") {
		return elemental.NewError(
			"Validation Error",
			"Name must not contain any '/' during creation",
			"a3s",
			http.StatusUnprocessableEntity,
		)
	}

	rns := bctx.Request().Namespace
	if rns == "/" {
		rns = ""
	}

	ns.Name = strings.Join([]string{rns, ns.Name}, "/")

	return crud.Create(bctx, p.manipulator, ns, crud.OptionPostWriteHook(p.makeNotify(bctx.Request().Operation)))
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
	return crud.Update(bctx, p.manipulator, bctx.InputData().(*api.Namespace),
		crud.OptionPostWriteHook(p.makeNotify(bctx.Request().Operation)),
	)
}

// ProcessDelete handles the delete requests for Namespaces.
func (p *NamespacesProcessor) ProcessDelete(bctx bahamut.Context) error {
	return crud.Delete(bctx, p.manipulator, api.NewNamespace(),
		crud.OptionPostWriteHook(func(obj elemental.Identifiable) {

			ndr := api.NewNamespaceDeletionRecord()
			ndr.Namespace = obj.(*api.Namespace).Name
			ndr.DeleteTime = time.Now()

			if err := p.manipulator.Create(manipulate.NewContext(bctx.Context()), ndr); err != nil {
				zap.L().Error("Unable to create namespace deletion record",
					zap.String("namespace", ndr.Namespace),
					zap.Error(err),
				)
			}

			p.makeNotify(bctx.Request().Operation)(obj)
		}),
	)
}

// ProcessInfo handles the info request for Namespaces.
func (p *NamespacesProcessor) ProcessInfo(bctx bahamut.Context) error {
	return crud.Info(bctx, p.manipulator, api.NamespaceIdentity)
}

func (p *NamespacesProcessor) makeNotify(op elemental.Operation) crud.PostWriteHook {
	return func(obj elemental.Identifiable) {
		_ = notification.Publish(
			p.pubsub,
			nscache.NotificationNamespaceChanges,
			&notification.Message{
				Type: string(op),
				Data: obj.(*api.Namespace).Name,
			},
		)
	}
}
