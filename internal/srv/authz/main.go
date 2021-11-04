package authz

import (
	"context"

	"go.aporeto.io/a3s/internal/srv/authz/internal/processors"
	"go.aporeto.io/a3s/pkgs/api"
	"go.aporeto.io/a3s/pkgs/permissions"
	"go.aporeto.io/bahamut"
	"go.aporeto.io/manipulate"
)

// Init initializes the module.
func Init(
	ctx context.Context,
	server bahamut.Server,
	m manipulate.Manipulator,
	retriever permissions.Retriever,
) error {

	// jwtCert, jwtKey, err := cfg.JWTCertificate()
	// if err != nil {
	// 	return err
	// }

	bahamut.RegisterProcessorOrDie(server, processors.NewAuthzProcessor(retriever), api.AuthzIdentity)

	return nil
}
