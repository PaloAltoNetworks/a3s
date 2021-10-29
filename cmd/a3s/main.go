package main

import (
	"context"
	"errors"
	"net/http"

	"go.aporeto.io/a3s/internal/srv/policy"
	"go.aporeto.io/a3s/pkgs/bootstrap"
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
	// db.Bootstrap(manipulator, serviceName)

	pubsub := bootstrap.MakeNATSClient(cfg.NATSConf)
	defer pubsub.Disconnect() // nolint: errcheck

	server := bahamut.New(
		append(
			bootstrap.ConfigureBahamut(
				ctx,
				cfg,
				pubsub,
				nil,
				nil,
				nil,
				nil,
			),
			bahamut.OptErrorTransformer(errorTransformer),
			bahamut.OptIdentifiableRetriever(bootstrap.MakeIdentifiableRetriever(manipulator)),
		)...,
	)

	// if err := authn.Init(ctx, cfg.AuthNConf, server, manipulator, pubsub); err != nil {
	// 	zap.L().Fatal("Unable to initialize authn module", zap.Error(err))
	// }

	if err := policy.Init(ctx, cfg.PolicyConf, server, manipulator, pubsub); err != nil {
		zap.L().Fatal("Unable to initialize policy module", zap.Error(err))
	}

	server.Run(ctx)
}

func errorTransformer(err error) error {

	if errors.As(err, &manipulate.ErrObjectNotFound{}) {
		return elemental.NewError("Not Found", err.Error(), "a3s", http.StatusNotFound)
	}

	if errors.As(err, &manipulate.ErrCannotCommunicate{}) {
		return elemental.NewError("Communication Error", err.Error(), "a3s", http.StatusServiceUnavailable)
	}

	return err
}
