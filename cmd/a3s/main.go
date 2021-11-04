package main

import (
	"context"
	"crypto/tls"
	"errors"
	"fmt"
	"net/http"

	"go.aporeto.io/a3s/internal/authorizer"
	"go.aporeto.io/a3s/internal/srv/authn"
	"go.aporeto.io/a3s/internal/srv/policy"
	"go.aporeto.io/a3s/pkgs/api"
	"go.aporeto.io/a3s/pkgs/authenticator"
	"go.aporeto.io/a3s/pkgs/bootstrap"
	"go.aporeto.io/a3s/pkgs/indexes"
	"go.aporeto.io/bahamut"
	"go.aporeto.io/elemental"
	"go.aporeto.io/manipulate"
	"go.uber.org/zap"
)

func main() {

	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	bahamut.InstallSIGINTHandler(cancel)

	cfg := newConf()

	if close := bootstrap.ConfigureLogger("a3s", cfg.LoggingConf); close != nil {
		defer close()
	}

	manipulator := bootstrap.MakeMongoManipulator(cfg.MongoConf)
	if err := indexes.Ensure(manipulator, api.Manager(), "a3s"); err != nil {
		zap.L().Fatal("Unable to ensure indexes", zap.Error(err))
	}
	if err := createRootNamespaceIfNeeded(manipulator); err != nil {
		zap.L().Fatal("Unable to handle root namespace", zap.Error(err))
	}

	jwtCert, _, err := cfg.AuthNConf.JWTCertificate()
	if err != nil {
		zap.L().Fatal("Unable to get JWT certificate", zap.Error(err))
	}

	pubsub := bootstrap.MakeNATSClient(cfg.NATSConf)
	defer pubsub.Disconnect() // nolint: errcheck

	authz := authorizer.NewLocalAuthorizer(manipulator)

	server := bahamut.New(
		append(
			bootstrap.ConfigureBahamut(
				ctx,
				cfg,
				pubsub,
				nil,
				[]bahamut.RequestAuthenticator{
					authenticator.NewPublic(api.IssueIdentity.Name),
					authenticator.NewPrivate(jwtCert),
				},
				[]bahamut.SessionAuthenticator{
					authenticator.NewPrivate(jwtCert),
				},
				nil,
				// []bahamut.Authorizer{
				// 	authz,
				// },
			),
			bahamut.OptMTLS(nil, tls.RequestClientCert),
			bahamut.OptErrorTransformer(errorTransformer),
			bahamut.OptIdentifiableRetriever(bootstrap.MakeIdentifiableRetriever(manipulator)),
		)...,
	)

	if err := authn.Init(ctx, cfg.AuthNConf, server, manipulator, pubsub); err != nil {
		zap.L().Fatal("Unable to initialize authn module", zap.Error(err))
	}

	if err := policy.Init(ctx, cfg.PolicyConf, server, manipulator, authz, pubsub); err != nil {
		zap.L().Fatal("Unable to initialize policy module", zap.Error(err))
	}

	server.Run(ctx)
}

func errorTransformer(err error) error {

	if errors.As(err, &manipulate.ErrObjectNotFound{}) {
		return elemental.NewError("Not Found", err.Error(), "a3s", http.StatusNotFound)
	}

	if errors.As(err, &manipulate.ErrConstraintViolation{}) {
		return elemental.NewError("Constraint Violation", err.Error(), "a3s", http.StatusUnprocessableEntity)
	}

	if errors.As(err, &manipulate.ErrCannotCommunicate{}) {
		return elemental.NewError("Communication Error", err.Error(), "a3s", http.StatusServiceUnavailable)
	}

	return err
}

func createRootNamespaceIfNeeded(m manipulate.Manipulator) error {

	mctx := manipulate.NewContext(context.Background(),
		manipulate.ContextOptionFilter(
			elemental.NewFilterComposer().
				WithKey("name").Equals("/").
				Done(),
		),
	)

	c, err := m.Count(mctx, api.NamespaceIdentity)
	if err != nil {
		return fmt.Errorf("unable to check if root namespace exists: %w", err)
	}

	if c == 1 {
		return nil
	}

	if c > 1 {
		panic("more than one namespace / found")
	}

	ns := api.NewNamespace()
	ns.Name = "/"
	ns.Namespace = "root"

	if err := m.Create(nil, ns); err != nil {
		return fmt.Errorf("unable to create root namespace: %w", err)
	}

	return nil
}
