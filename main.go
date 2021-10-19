package main

import (
	"context"

	"go.aporeto.io/a3s/pkgs/bootstrap"
	"go.aporeto.io/a3s/srv/authn"
	"go.aporeto.io/bahamut"
	"go.aporeto.io/elemental"
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
			bahamut.OptUnmarshallers(map[elemental.Identity]bahamut.CustomUmarshaller{
				// 	gaia.IssueIdentity: unmarshallers.FormData,
			}),
			// bahamut.OptTraceCleaner(tracecleaner.Clean),
		)...,
	)

	if err := authn.Init(ctx, cfg.AuthNConf, server, manipulator, pubsub); err != nil {
		zap.L().Fatal("Unable to initialize authn module", zap.Error(err))
	}

	server.Run(ctx)
}
