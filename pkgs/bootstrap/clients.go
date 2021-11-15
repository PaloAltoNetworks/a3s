package bootstrap

import (
	"context"
	"fmt"
	"time"

	"go.aporeto.io/a3s/internal/conf"
	"go.aporeto.io/a3s/pkgs/api"
	"go.aporeto.io/a3s/pkgs/sharder"
	"go.aporeto.io/bahamut"
	"go.aporeto.io/elemental"
	"go.aporeto.io/manipulate"
	"go.aporeto.io/manipulate/manipmongo"
	"go.uber.org/zap"
)

// MakeNATSClient returns a connected pubsub server client.
// This function is not meant to be used outside of the platform. It will fatal
// anytime it feels like it.
func MakeNATSClient(cfg conf.NATSConf) bahamut.PubSubClient {

	opts := []bahamut.NATSOption{
		bahamut.NATSOptClientID(cfg.NATSClientID),
		bahamut.NATSOptClusterID(cfg.NATSClusterID),
		bahamut.NATSOptCredentials(cfg.NATSUser, cfg.NATSPassword),
	}

	tlscfg, err := cfg.TLSConfig()
	if err != nil {
		zap.L().Fatal("Unable to prepare TLS config for nats", zap.Error(err))
	}

	if tlscfg != nil {
		opts = append(opts, bahamut.NATSOptTLS(tlscfg))
	}

	pubsub := bahamut.NewNATSPubSubClient(cfg.NATSURL, opts...)

	ctx, cancel := context.WithTimeout(context.Background(), 30*time.Second)
	defer cancel()

	if err := pubsub.Connect(ctx); err != nil {
		zap.L().Fatal("Could not connect to nats", zap.Error(err))
	}

	zap.L().Info("Connected to nats", zap.String("server", cfg.NATSURL))

	return pubsub
}

// MakeMongoManipulator returns a configured mongo manipulator.
// This function is not meant to be used outside of the platform. It will fatal
// anytime it feels like it.
func MakeMongoManipulator(cfg conf.MongoConf, hasher sharder.Hasher, additionalOptions ...manipmongo.Option) manipulate.TransactionalManipulator {

	var consistency manipulate.ReadConsistency
	switch cfg.MongoConsistency {
	case "strong":
		consistency = manipulate.ReadConsistencyStrong
	case "monotonic":
		consistency = manipulate.ReadConsistencyMonotonic
	case "eventual":
		consistency = manipulate.ReadConsistencyEventual
	case "nearest":
		consistency = manipulate.ReadConsistencyNearest
	case "weakest":
		consistency = manipulate.ReadConsistencyWeakest
	default:
		panic(fmt.Sprintf("unknown consistency '%s'", cfg.MongoConsistency))
	}

	opts := append(
		[]manipmongo.Option{
			manipmongo.OptionCredentials(cfg.MongoUser, cfg.MongoPassword, cfg.MongoAuthDB),
			manipmongo.OptionConnectionPoolLimit(cfg.MongoPoolSize),
			manipmongo.OptionDefaultReadConsistencyMode(consistency),
			manipmongo.OptionTranslateKeysFromModelManager(api.Manager()),
			manipmongo.OptionSharder(sharder.New(hasher)),
			manipmongo.OptionDefaultRetryFunc(func(i manipulate.RetryInfo) error {
				info := i.(manipmongo.RetryInfo)
				zap.L().Debug("mongo manipulator retry",
					zap.Int("try", info.Try()),
					zap.String("operation", string(info.Operation)),
					zap.String("identity", info.Identity.Name),
					zap.Error(info.Err()),
				)
				return nil
			}),
		},
		additionalOptions...,
	)

	tlscfg, err := cfg.TLSConfig()
	if err != nil {
		zap.L().Fatal("Unable to prepare TLS config for nats", zap.Error(err))
	}

	if tlscfg != nil {
		opts = append(opts, manipmongo.OptionTLS(tlscfg))
	}

	if cfg.MongoAttrEncryptKey != "" {
		encrypter, err := elemental.NewAESAttributeEncrypter(cfg.MongoAttrEncryptKey)
		if err != nil {
			zap.L().Fatal("Unable to create mongodb attribute encrypter", zap.Error(err))
		}
		opts = append(opts, manipmongo.OptionAttributeEncrypter(encrypter))
		zap.L().Info("Attribute encryption", zap.String("status", "enabled"))
	} else {
		zap.L().Warn("Attribute encryption", zap.String("status", "disabled"))
	}

	m, err := manipmongo.New(cfg.MongoURL, cfg.MongoDBName, opts...)
	if err != nil {
		zap.L().Fatal("Unable to connect to mongo", zap.Error(err))
	}

	zap.L().Info("Connected to mongodb", zap.String("url", cfg.MongoURL), zap.String("db", cfg.MongoDBName))

	return m
}
