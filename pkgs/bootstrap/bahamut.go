package bootstrap

import (
	"context"
	"crypto/tls"
	"errors"
	"net/http"

	"github.com/fatih/structs"
	"github.com/opentracing/opentracing-go"
	"go.aporeto.io/a3s/pkgs/conf"
	"go.aporeto.io/bahamut"
	"go.aporeto.io/bahamut/authorizer/simple"
	"go.aporeto.io/bahamut/gateway/upstreamer/push"
	"go.aporeto.io/elemental"
	"go.aporeto.io/manipulate"
	"go.uber.org/zap"
	"go.uber.org/zap/zapcore"
)

// ConfigureBahamut returns a list of bahamut.Option based on provided configuration.
func ConfigureBahamut(
	ctx context.Context,
	cfg interface{},
	pubsub bahamut.PubSubClient,
	apiManager elemental.ModelManager,
	healthHandler bahamut.HealthServerFunc,
	requestAuthenticators []bahamut.RequestAuthenticator,
	sessionAuthenticators []bahamut.SessionAuthenticator,
	authorizers []bahamut.Authorizer,
) (opts []bahamut.Option) {

	modelManagers := map[int]elemental.ModelManager{0: apiManager, 1: apiManager}

	l, err := zap.NewStdLogAt(zap.L(), zapcore.DebugLevel)
	if err != nil {
		zap.L().Fatal("Unable to create logger for bahamut HTTP server", zap.Error(err))
	}

	// Default options.
	opts = []bahamut.Option{
		// bahamut.OptServiceInfo(serviceName, serviceVersion, subversions),
		bahamut.OptModel(modelManagers),
		bahamut.OptAuthenticators(requestAuthenticators, sessionAuthenticators),
		bahamut.OptAuthorizers(authorizers),
		bahamut.OptOpentracingTracer(opentracing.GlobalTracer()),
		bahamut.OptDisableCompression(),
		bahamut.OptHTTPLogger(l),
		bahamut.OptErrorTransformer(ErrorTransformer),
	}

	cs := structs.New(cfg)

	if f, ok := cs.FieldOk(structs.Name(conf.APIServerConf{})); ok {
		c := f.Value().(conf.APIServerConf)
		opts = append(
			opts,
			bahamut.OptRestServer(c.ListenAddress),
			bahamut.OptMaxConnection(c.MaxConnections),
		)

		zap.L().Info("Max TCP connections", zap.Int("max", c.MaxConnections))

		tlscfg, err := c.TLSConfig()
		if err != nil {
			zap.L().Fatal("Unable to configure tls", zap.Error(err))
		}

		if tlscfg != nil {

			opts = append(opts,
				bahamut.OptTLS(tlscfg.Certificates, nil),
				bahamut.OptTLSNextProtos([]string{"h2"}), // enable http2 support.
			)

			if clientCA := tlscfg.ClientCAs; clientCA != nil {
				opts = append(opts, bahamut.OptMTLS(clientCA, tls.RequireAndVerifyClientCert))
			}
		}

		if c.CORSDefaultOrigin != "" || len(c.CORSAdditionalOrigins) > 0 {
			opts = append(
				opts,
				bahamut.OptCORSAccessControl(
					bahamut.NewDefaultCORSController(
						c.CORSDefaultOrigin,
						c.CORSAdditionalOrigins,
					),
				),
			)
			zap.L().Info("CORS origin configured",
				zap.String("default", c.CORSDefaultOrigin),
				zap.Strings("additional", c.CORSAdditionalOrigins),
			)
		}
	}

	if f, ok := cs.FieldOk(structs.Name(conf.HealthConfiguration{})); ok {
		c := f.Value().(conf.HealthConfiguration)
		if c.EnableHealth {
			opts = append(
				opts,
				bahamut.OptHealthServer(c.HealthListenAddress, healthHandler),
				bahamut.OptHealthServerMetricsManager(bahamut.NewPrometheusMetricsManager()),
			)
		}
	}

	if f, ok := cs.FieldOk(structs.Name(conf.ProfilingConf{})); ok {
		c := f.Value().(conf.ProfilingConf)
		if c.ProfilingEnabled {
			opts = append(opts, bahamut.OptProfilingLocal(c.ProfilingListenAddress))
		}
	}

	if f, ok := cs.FieldOk(structs.Name(conf.RateLimitingConf{})); ok {
		c := f.Value().(conf.RateLimitingConf)
		if c.RateLimitingEnabled {
			opts = append(opts, bahamut.OptRateLimiting(float64(c.RateLimitingRPS), c.RateLimitingBurst))
			zap.L().Info("Rate limit configured",
				zap.Int("rps", c.RateLimitingRPS),
				zap.Int("burst", c.RateLimitingBurst),
			)
		}
	}

	if f, ok := cs.FieldOk(structs.Name(conf.HTTPTimeoutsConf{})); ok {
		c := f.Value().(conf.HTTPTimeoutsConf)
		opts = append(opts, bahamut.OptTimeouts(c.TimeoutRead, c.TimeoutWrite, c.TimeoutIdle))

		zap.L().Debug("Timeouts configured",
			zap.Duration("read", c.TimeoutRead),
			zap.Duration("write", c.TimeoutWrite),
			zap.Duration("idle", c.TimeoutIdle),
		)
	}

	if f, ok := cs.FieldOk(structs.Name(conf.NATSPublisherConf{})); ok {
		c := f.Value().(conf.NATSPublisherConf)
		opts = append(opts,
			bahamut.OptPushServer(pubsub, c.NATSPublishTopic),
			bahamut.OptPushServerEnableSubjectHierarchies(),
		)
	}

	return opts
}

// MakeBahamutGatewayNotifier returns the bahamut options needed
// to make A3S announce itself to a bahamut gateway.
func MakeBahamutGatewayNotifier(
	ctx context.Context,
	pubsub bahamut.PubSubClient,
	serviceName string,
	gatewayTopic string,
	anouncedAddress string,
	nopts ...push.NotifierOption,
) []bahamut.Option {

	opts := []bahamut.Option{}

	if gatewayTopic == "" {
		return nil
	}

	nw := push.NewNotifier(
		pubsub,
		gatewayTopic,
		serviceName,
		anouncedAddress,
		nopts...,
	)

	opts = append(opts,
		bahamut.OptPostStartHook(nw.MakeStartHook(ctx)),
		bahamut.OptPreStopHook(nw.MakeStopHook()),
	)

	zap.L().Info(
		"Gateway topic set",
		zap.String("topic", gatewayTopic),
		zap.String("service", serviceName),
	)

	return opts
}

// ErrorTransformer transforms a disconnected error into an not acceptable.
// This avoid 500 errors due to clients being disconnected.
func ErrorTransformer(err error) error {

	switch {

	case errors.As(err, &manipulate.ErrDisconnected{}),
		errors.As(err, &manipulate.ErrDisconnected{}),
		errors.Is(err, context.Canceled):

		return elemental.NewError(
			"Client Disconnected",
			err.Error(),
			"a3s",
			http.StatusNotAcceptable,
		)

	case manipulate.IsObjectNotFoundError(err):

		return elemental.NewError(
			"Not Found",
			err.Error(),
			"a3s",
			http.StatusNotFound,
		)

	default:
		return nil
	}
}

// MakeIdentifiableRetriever returns a bahamut.IdentifiableRetriever to handle patches as classic update.
func MakeIdentifiableRetriever(
	manipulator manipulate.Manipulator,
	apiManager elemental.ModelManager,
) bahamut.IdentifiableRetriever {

	return func(req *elemental.Request) (elemental.Identifiable, error) {

		identity := req.Identity

		obj := apiManager.Identifiable(identity)
		obj.SetIdentifier(req.ObjectID)

		if err := manipulator.Retrieve(nil, obj); err != nil {
			return nil, err
		}

		return obj, nil
	}
}

// MakePublishHandler returns a bahamut.PushPublishHandler that publishes all events but the
// ones related to the given identities.
func MakePublishHandler(excludedIdentities []elemental.Identity) bahamut.PushPublishHandler {

	return simple.NewPublishHandler(func(event *elemental.Event) (bool, error) {
		for _, i := range excludedIdentities {
			if event.Identity == i.Name {
				return false, nil
			}
		}
		return true, nil
	})
}
