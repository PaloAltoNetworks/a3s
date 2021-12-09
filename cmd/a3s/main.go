package main

import (
	"context"
	"crypto/tls"
	"errors"
	"fmt"
	"net"
	"net/http"
	"net/url"
	"os"
	"strings"
	"time"

	"github.com/globalsign/mgo"
	"go.aporeto.io/a3s/internal/hasher"
	"go.aporeto.io/a3s/internal/processors"
	"go.aporeto.io/a3s/internal/ui"
	"go.aporeto.io/a3s/pkgs/api"
	"go.aporeto.io/a3s/pkgs/authenticator"
	"go.aporeto.io/a3s/pkgs/authorizer"
	"go.aporeto.io/a3s/pkgs/bootstrap"
	"go.aporeto.io/a3s/pkgs/indexes"
	"go.aporeto.io/a3s/pkgs/notification"
	"go.aporeto.io/a3s/pkgs/nscache"
	"go.aporeto.io/a3s/pkgs/permissions"
	"go.aporeto.io/a3s/pkgs/token"
	"go.aporeto.io/bahamut"
	"go.aporeto.io/elemental"
	"go.aporeto.io/manipulate"
	"go.aporeto.io/manipulate/manipmongo"
	"go.aporeto.io/tg/tglib"
	"go.uber.org/zap"
)

var (
	publicResources = []string{
		api.IssueIdentity.Category,
		api.PermissionsIdentity.Category,
		api.AuthzIdentity.Category,
	}
)

func main() {

	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	bahamut.InstallSIGINTHandler(cancel)

	cfg := newConf()

	if close := bootstrap.ConfigureLogger("a3s", cfg.LoggingConf); close != nil {
		defer close()
	}

	m := bootstrap.MakeMongoManipulator(cfg.MongoConf, &hasher.Hasher{})
	if err := indexes.Ensure(m, api.Manager(), "a3s"); err != nil {
		zap.L().Fatal("Unable to ensure indexes", zap.Error(err))
	}

	if err := manipmongo.EnsureIndex(m, elemental.MakeIdentity("oidccache", "oidccache"), mgo.Index{
		Key:         []string{"time"},
		ExpireAfter: 1 * time.Minute,
		Name:        "index_expiration_exp",
	}); err != nil {
		zap.L().Fatal("Unable to create exp expiration index for oidccache", zap.Error(err))
	}

	if err := createRootNamespaceIfNeeded(m); err != nil {
		zap.L().Fatal("Unable to handle root namespace", zap.Error(err))
	}

	if cfg.Init {
		initialized, err := initRootPermissions(ctx, m, cfg.InitRootUserCAPath, cfg.InitContinue)
		if err != nil {
			zap.L().Fatal("unable to initialize root permissions", zap.Error(err))
			return
		}

		if initialized {
			zap.L().Info("Root auth initialized")
		}

		if !cfg.InitContinue {
			return
		}
	}

	jwtCert, jwtKey, err := cfg.JWT.JWTCertificate()
	if err != nil {
		zap.L().Fatal("Unable to get JWT certificate", zap.Error(err))
	}

	zap.L().Info("JWT info configured",
		zap.String("iss", cfg.JWT.JWTIssuer),
		zap.String("aud", cfg.JWT.JWTAudience),
	)

	jwks := token.NewJWKS()
	if err := jwks.AppendWithPrivate(jwtCert, jwtKey); err != nil {
		zap.L().Fatal("unable to build JWKS", zap.Error(err))
	}

	publicAPIURL := cfg.PublicAPIURL
	if publicAPIURL == "" {
		publicAPIURL = fmt.Sprintf("https://%s", getNotifierEndpoint(cfg.ListenAddress))
	}

	zap.L().Info("Announced public API", zap.String("url", publicAPIURL))
	cookiePolicy := http.SameSiteDefaultMode
	switch cfg.JWT.JWTCookiePolicy {
	case "strict":
		cookiePolicy = http.SameSiteStrictMode
	case "lax":
		cookiePolicy = http.SameSiteLaxMode
	case "none":
		cookiePolicy = http.SameSiteNoneMode
	}
	zap.L().Info("Cookie policy set", zap.String("policy", cfg.JWT.JWTCookiePolicy))

	cookieDomain := cfg.JWT.JWTCookieDomain
	if cookieDomain == "" {
		u, err := url.Parse(publicAPIURL)
		if err != nil {
			zap.L().Fatal("Unable to parse publicAPIURL", zap.Error(err))
		}
		cookieDomain = u.Hostname()
	}
	zap.L().Info("Cookie domain set", zap.String("domain", cookieDomain))

	pubsub := bootstrap.MakeNATSClient(cfg.NATSConf)
	defer pubsub.Disconnect() // nolint: errcheck

	pauthn := authenticator.NewPrivate(jwks, cfg.JWT.JWTIssuer, cfg.JWT.JWTAudience)
	retriever := permissions.NewRetriever(m)
	pauthz := authorizer.New(
		ctx,
		retriever,
		pubsub,
		authorizer.OptionIgnoredResources(publicResources...),
	)

	opts := append(
		bootstrap.ConfigureBahamut(
			ctx,
			cfg,
			pubsub,
			nil,
			[]bahamut.RequestAuthenticator{
				authenticator.NewPublic(publicResources...),
				pauthn,
			},
			[]bahamut.SessionAuthenticator{
				pauthn,
			},
			[]bahamut.Authorizer{
				pauthz,
			},
		),
		bahamut.OptPushDispatchHandler(authorizer.NewPushDispatchHandler(m, pauthz)),
		bahamut.OptPushPublishHandler(bootstrap.MakePublishHandler(nil)),
		bahamut.OptMTLS(nil, tls.RequestClientCert),
		bahamut.OptErrorTransformer(errorTransformer),
		bahamut.OptIdentifiableRetriever(bootstrap.MakeIdentifiableRetriever(m)),
	)

	if cfg.NATSGWTopic != "" {
		opts = append(
			opts,
			bootstrap.MakeBahamutGatewayNotifier(
				ctx,
				pubsub,
				cfg.NATSGWTopic,
				getNotifierEndpoint(cfg.ListenAddress),
			)...,
		)
	}

	server := bahamut.New(opts...)

	if err := server.RegisterCustomRouteHandler("/.well-known/jwks.json", makeJWKSHandler(jwks)); err != nil {
		zap.L().Fatal("Unable to install jwks handler", zap.Error(err))
	}

	if err := server.RegisterCustomRouteHandler("/ui/login.html", makeUILoginHandler(publicAPIURL)); err != nil {
		zap.L().Fatal("Unable to install UI login handler", zap.Error(err))
	}

	bahamut.RegisterProcessorOrDie(server, processors.NewIssueProcessor(m, jwks, cfg.JWT.JWTMaxValidity, cfg.JWT.JWTIssuer, cfg.JWT.JWTAudience, cookiePolicy, cookieDomain), api.IssueIdentity)
	bahamut.RegisterProcessorOrDie(server, processors.NewMTLSSourcesProcessor(m), api.MTLSSourceIdentity)
	bahamut.RegisterProcessorOrDie(server, processors.NewLDAPSourcesProcessor(m), api.LDAPSourceIdentity)
	bahamut.RegisterProcessorOrDie(server, processors.NewOIDCSourcesProcessor(m), api.OIDCSourceIdentity)
	bahamut.RegisterProcessorOrDie(server, processors.NewA3SSourcesProcessor(m), api.A3SSourceIdentity)
	bahamut.RegisterProcessorOrDie(server, processors.NewPermissionsProcessor(retriever), api.PermissionsIdentity)
	bahamut.RegisterProcessorOrDie(server, processors.NewAuthzProcessor(pauthz, jwks, cfg.JWT.JWTIssuer, cfg.JWT.JWTAudience), api.AuthzIdentity)
	bahamut.RegisterProcessorOrDie(server, processors.NewNamespacesProcessor(m, pubsub), api.NamespaceIdentity)
	bahamut.RegisterProcessorOrDie(server, processors.NewAuthorizationProcessor(m, pubsub, retriever), api.AuthorizationIdentity)

	notification.Subscribe(ctx, pubsub, nscache.NotificationNamespaceChanges, makeNamespaceCleaner(ctx, m))

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

func initRootPermissions(ctx context.Context, m manipulate.Manipulator, caPath string, ifNeeded bool) (bool, error) {

	caData, err := os.ReadFile(caPath)
	if err != nil {
		return false, fmt.Errorf("unable to read root user ca: %w", err)
	}

	caCerts, err := tglib.ParseCertificates(caData)
	if err != nil {
		return false, fmt.Errorf("unable to parse root user ca: %w", err)
	}

	chain := make([]string, len(caCerts))
	for i, cert := range caCerts {
		chain[i] = token.Fingerprint(cert)
	}

	source := api.NewMTLSSource()
	source.Namespace = "/"
	source.Name = "root"
	source.Description = "Root auth source used to bootstrap permissions."
	source.CA = string(caData)
	if err := m.Create(manipulate.NewContext(ctx), source); err != nil {
		if errors.As(err, &manipulate.ErrConstraintViolation{}) && ifNeeded {
			return false, nil
		}
		return false, fmt.Errorf("unable to create root mtls auth source: %w", err)
	}

	auth := api.NewAuthorization()
	auth.Namespace = "/"
	auth.Name = "root-mtls-authorization"
	auth.Description = "Root authorization for certificates issued from the CA declared  in the root auth mtls source."
	auth.Subject = [][]string{
		{
			"@source:type=mtls",
			"@source:name=root",
			"@source:namespace=/",
			fmt.Sprintf("issuerchain=%s", strings.Join(chain, ",")),
		},
	}
	auth.FlattenedSubject = auth.Subject[0]
	auth.Permissions = []string{"*:*"}
	auth.TargetNamespaces = []string{"/"}
	auth.Hidden = true

	if err := m.Create(manipulate.NewContext(ctx), auth); err != nil {
		return false, fmt.Errorf("unable to create root auth: %w", err)
	}

	return true, nil
}

func makeJWKSHandler(jwks *token.JWKS) http.HandlerFunc {

	return func(w http.ResponseWriter, req *http.Request) {

		jwks.RLock()
		defer jwks.RUnlock()

		data, err := elemental.Encode(elemental.EncodingTypeJSON, jwks)
		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}

		w.Header().Add("Content-Type", "application/json")
		_, _ = w.Write(data)
	}
}

func makeUILoginHandler(api string) http.HandlerFunc {

	return func(w http.ResponseWriter, req *http.Request) {

		redirect := req.URL.Query().Get("redirect")
		if redirect == "" {
			http.Error(w, "Missing redirect query parameter", http.StatusBadRequest)
			return
		}

		audience := req.URL.Query().Get("audience")
		if audience == "" {
			audience = redirect
		}

		data, err := ui.GetLogin(api, redirect, audience)
		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}

		w.Header().Add("Content-Type", "text/html")
		_, _ = w.Write(data)
	}
}

func makeNamespaceCleaner(ctx context.Context, m manipulate.Manipulator) notification.Handler {

	return func(msg *notification.Message) {

		if msg.Type != string(elemental.OperationDelete) {
			return
		}

		ns := msg.Data.(string)

		for _, i := range api.Manager().AllIdentities() {
			mctx := manipulate.NewContext(
				ctx,
				manipulate.ContextOptionNamespace(ns),
				manipulate.ContextOptionRecursive(true),
			)
			if err := m.DeleteMany(mctx, i); err != nil {
				zap.L().Error("Unable to clean namespace", zap.String("ns", ns), zap.Error(err))
			}
		}
	}
}

func getNotifierEndpoint(listenAddress string) string {

	_, port, err := net.SplitHostPort(listenAddress)
	if err != nil {
		zap.L().Fatal("Unable to parse listen address", zap.Error(err))
	}

	host, err := os.Hostname()
	if err != nil {
		zap.L().Fatal("Unable to retrieve hostname", zap.Error(err))
	}

	addrs, err := net.LookupHost(host)
	if err != nil {
		zap.L().Fatal("Unable to resolve hostname", zap.Error(err))
	}

	if len(addrs) == 0 {
		zap.L().Fatal("Unable to find any IP in resolved hostname")
	}

	var endpoint string
	for _, addr := range addrs {
		ip := net.ParseIP(addr)
		if len(ip.To4()) == net.IPv4len {
			endpoint = addr
			break
		}
	}

	if endpoint == "" {
		endpoint = addrs[0]
	}

	return fmt.Sprintf("%s:%s", endpoint, port)
}
