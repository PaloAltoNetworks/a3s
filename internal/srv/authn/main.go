package authn

import (
	"go.aporeto.io/a3s/internal/authorizer"
	"go.aporeto.io/a3s/internal/srv/authn/internal/processors"
	"go.aporeto.io/a3s/pkgs/api"
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
	auth authorizer.Authorizer,
	pubsub bahamut.PubSubClient,
) error {

	jwtCert, jwtKey, err := cfg.JWTCertificate()
	if err != nil {
		return err
	}

	bahamut.RegisterProcessorOrDie(
		server,
		processors.NewIssueProcessor(m,
			jwtCert,
			jwtKey,
			cfg.JWTMaxValidity),
		api.IssueIdentity,
	)

	return nil
}
