package policy

import (
	"go.aporeto.io/a3s/internal/srv/policy/internal/processors"
	"go.aporeto.io/a3s/pkgs/api"
	"go.aporeto.io/a3s/pkgs/permissions"
	"go.aporeto.io/bahamut"
	"go.aporeto.io/manipulate"
	"golang.org/x/net/context"
)

// Init initializes the module.
func Init(
	ctx context.Context,
	cfg Conf,
	server bahamut.Server,
	m manipulate.Manipulator,
	retriever permissions.Retriever,
	pubsub bahamut.PubSubClient,
) error {

	bahamut.RegisterProcessorOrDie(server, processors.NewNamespacesProcessor(m, pubsub), api.NamespaceIdentity)
	bahamut.RegisterProcessorOrDie(server, processors.NewAuthorizationProcessor(m, pubsub, retriever), api.AuthorizationIdentity)

	return nil
}
