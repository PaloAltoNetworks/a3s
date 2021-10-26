package policy

import (
	"go.aporeto.io/a3s/pkgs/api"
	"go.aporeto.io/a3s/srv/policy/internal/processors"
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
	pubsub bahamut.PubSubClient,
) error {

	bahamut.RegisterProcessorOrDie(server, processors.NewNamespacesProcessor(m), api.NamespaceIdentity)

	return nil
}
