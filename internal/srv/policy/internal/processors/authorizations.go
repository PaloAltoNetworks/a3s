package processors

import (
	"go.aporeto.io/a3s/pkgs/api"
	"go.aporeto.io/a3s/pkgs/crud"
	"go.aporeto.io/a3s/pkgs/notification"
	"go.aporeto.io/a3s/pkgs/nscache"
	"go.aporeto.io/bahamut"
	"go.aporeto.io/elemental"
	"go.aporeto.io/manipulate"
)

// A AuthorizationsProcessor is a bahamut processor for Authorizations.
type AuthorizationsProcessor struct {
	manipulator manipulate.Manipulator
	pubsub      bahamut.PubSubClient
}

// NewAuthorizationProcessor returns a new AuthorizationsProcessor.
func NewAuthorizationProcessor(manipulator manipulate.Manipulator, pubsub bahamut.PubSubClient) *AuthorizationsProcessor {
	return &AuthorizationsProcessor{
		manipulator: manipulator,
		pubsub:      pubsub,
	}
}

// ProcessCreate handles the creates requests for Authorizations.
func (p *AuthorizationsProcessor) ProcessCreate(bctx bahamut.Context) error {
	return crud.Create(bctx, p.manipulator, bctx.InputData().(*api.Authorization),
		crud.OptionPostWriteHook(p.makeNotify()),
	)
}

// ProcessRetrieveMany handles the retrieve many requests for Authorizations.
func (p *AuthorizationsProcessor) ProcessRetrieveMany(bctx bahamut.Context) error {
	return crud.RetrieveMany(bctx, p.manipulator, &api.AuthorizationsList{})
}

// ProcessRetrieve handles the retrieve requests for Authorizations.
func (p *AuthorizationsProcessor) ProcessRetrieve(bctx bahamut.Context) error {
	return crud.Retrieve(bctx, p.manipulator, api.NewAuthorization())
}

// ProcessUpdate handles the update requests for Authorizations.
func (p *AuthorizationsProcessor) ProcessUpdate(bctx bahamut.Context) error {
	return crud.Update(bctx, p.manipulator, bctx.InputData().(*api.Authorization),
		crud.OptionPostWriteHook(p.makeNotify()),
	)
}

// ProcessDelete handles the delete requests for Authorizations.
func (p *AuthorizationsProcessor) ProcessDelete(bctx bahamut.Context) error {
	return crud.Delete(bctx, p.manipulator, api.NewAuthorization(),
		crud.OptionPostWriteHook(p.makeNotify()),
	)
}

// ProcessInfo handles the info request for Authorizations.
func (p *AuthorizationsProcessor) ProcessInfo(bctx bahamut.Context) error {
	return crud.Info(bctx, p.manipulator, api.AuthorizationIdentity)
}

func (p *AuthorizationsProcessor) makeNotify() crud.PostWriteHook {
	return func(obj elemental.Identifiable) {
		_ = notification.Publish(
			p.pubsub,
			nscache.NotificationNamespaceChanges,
			&notification.Message{
				Data: obj.(*api.Authorization).Namespace,
			},
		)
	}
}
