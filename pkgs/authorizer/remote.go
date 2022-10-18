package authorizer

import (
	"context"

	"go.aporeto.io/a3s/pkgs/api"
	"go.aporeto.io/a3s/pkgs/permissions"
	"go.aporeto.io/elemental"
	"go.aporeto.io/manipulate"
	"go.aporeto.io/manipulate/maniphttp"
)

type remoteAuthorizer struct {
	Authorizer
}

// NewRemote returns a ready to use bahamut.Authorizer that can be used over the API.
// This is meant to be use by external bahamut service.
// Updates of the namespace/authorization state comes from the websocket.
func NewRemote(ctx context.Context, m manipulate.Manipulator, options ...Option) Authorizer {

	subscriber := maniphttp.NewSubscriber(
		m,
		maniphttp.SubscriberOptionRecursive(true),
		maniphttp.SubscriberOptionNamespace(maniphttp.ExtractNamespace(m)),
	)

	pcfg := elemental.NewPushConfig()
	pcfg.FilterIdentity(api.NamespaceIdentity.Name)
	pcfg.FilterIdentity(api.AuthorizationIdentity.Name)

	subscriber.Start(ctx, pcfg)

	return &remoteAuthorizer{
		Authorizer: New(
			ctx,
			permissions.NewRemoteRetrieverWithTransformer(m, nil),
			&webSocketPubSub{subscriber: subscriber},
			options...,
		),
	}
}
