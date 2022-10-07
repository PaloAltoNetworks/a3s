package main

import (
	"context"
	"crypto/tls"
	"crypto/x509"
	"errors"
	"fmt"
	"net"
	"net/http"
	"net/url"
	"os"
	"strings"
	"time"

	"github.com/ghodss/yaml"
	"github.com/globalsign/mgo"
	"go.aporeto.io/a3s/internal/hasher"
	"go.aporeto.io/a3s/internal/processors"
	"go.aporeto.io/a3s/internal/ui"
	"go.aporeto.io/a3s/pkgs/api"
	"go.aporeto.io/a3s/pkgs/authenticator"
	"go.aporeto.io/a3s/pkgs/authorizer"
	"go.aporeto.io/a3s/pkgs/bearermanip"
	"go.aporeto.io/a3s/pkgs/bootstrap"
	"go.aporeto.io/a3s/pkgs/conf"
	"go.aporeto.io/a3s/pkgs/importing"
	"go.aporeto.io/a3s/pkgs/indexes"
	"go.aporeto.io/a3s/pkgs/notification"
	"go.aporeto.io/a3s/pkgs/nscache"
	"go.aporeto.io/a3s/pkgs/permissions"
	"go.aporeto.io/a3s/pkgs/push"
	"go.aporeto.io/a3s/pkgs/token"
	"go.aporeto.io/bahamut"
	"go.aporeto.io/elemental"
	"go.aporeto.io/manipulate"
	"go.aporeto.io/manipulate/manipmongo"
	"go.aporeto.io/tg/tglib"
	"go.uber.org/zap"

	gwpush "go.aporeto.io/bahamut/gateway/upstreamer/push"
)

var (
	publicResources = []string{
		api.IssueIdentity.Category,
		api.PermissionsIdentity.Category,
		api.AuthzIdentity.Category,
	}
	pushExcludedResources = []elemental.Identity{
		api.PermissionsIdentity,

		// safety: these ones are not an identifiable, so it would not be pushed anyway.
		api.IssueIdentity,
		api.AuthzIdentity,
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

	if cfg.InitDB {
		if err := createMongoDBAccount(cfg.MongoConf, cfg.InitDBUsername); err != nil {
			zap.L().Fatal("Unable to create mongodb account", zap.Error(err))
		}

		if !cfg.InitContinue {
			return
		}
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
		if cfg.InitRootUserCAPath != "" {
			initialized, err := initRootPermissions(ctx, m, cfg.InitRootUserCAPath, cfg.JWT.JWTIssuer, cfg.InitContinue)
			if err != nil {
				zap.L().Fatal("Unable to initialize root permissions", zap.Error(err))
				return
			}

			if initialized {
				zap.L().Info("Root auth initialized")
			}
		}

		if cfg.InitPlatformCAPath != "" {
			initialized, err := initPlatformPermissions(ctx, m, cfg.InitPlatformCAPath, cfg.JWT.JWTIssuer, cfg.InitContinue)
			if err != nil {
				zap.L().Fatal("Unable to initialize platform permissions", zap.Error(err))
				return
			}

			if initialized {
				zap.L().Info("Platform auth initialized")
			}
		}

		if cfg.InitData != "" {
			initialized, err := initData(ctx, m, cfg.InitData)
			if err != nil {
				zap.L().Fatal("Unable to init provisionning data", zap.Error(err))
				return
			}

			if initialized {
				zap.L().Info("Initial provisionning initialized")
			}
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

	if cfg.MTLSHeader.Enabled {
		if cfg.MTLSHeader.Passphrase == "" {
			zap.L().Fatal("You must pass --mtls-header-passphrase when --mtls-header-enabled is set")
		}
		var cipher string
		switch len(cfg.MTLSHeader.Passphrase) {
		case 16:
			cipher = "AES-128"
		case 24:
			cipher = "AES-192"
		case 32:
			cipher = "AES-256"
		default:
			zap.L().Fatal("The value for --mtls-header-passphrase must be 16, 24 or 32 bytes long to select AES-128, AES-192 or AES-256")
		}
		zap.L().Info("MTLS header trust set", zap.String("header", cfg.MTLSHeader.HeaderKey), zap.String("cipher", cipher))
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

	trustedIssuers, err := cfg.JWT.TrustedIssuers()
	if err != nil {
		zap.L().Fatal("Unable to build trusted issuers list", zap.Error(err))
	}
	if len(trustedIssuers) > 0 {
		zap.L().Info("Trusted issuers set",
			zap.Strings(
				"issuers",
				func() []string {
					out := make([]string, len(trustedIssuers))
					for i, o := range trustedIssuers {
						out[i] = o.URL
					}
					return out
				}(),
			),
		)
	}

	pubsub := bootstrap.MakeNATSClient(cfg.NATSConf)
	defer pubsub.Disconnect() // nolint: errcheck

	pauthn := authenticator.New(
		jwks,
		cfg.JWT.JWTIssuer,
		cfg.JWT.JWTAudience,
		authenticator.OptionIgnoredResources(publicResources...),
		authenticator.OptionExternalTrustedIssuers(trustedIssuers...),
	)
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
			api.Manager(),
			nil,
			[]bahamut.RequestAuthenticator{pauthn},
			[]bahamut.SessionAuthenticator{pauthn},
			[]bahamut.Authorizer{pauthz},
		),
		bahamut.OptPushDispatchHandler(push.NewDispatcher(pauthz)),
		bahamut.OptPushPublishHandler(bootstrap.MakePublishHandler(pushExcludedResources)),
		bahamut.OptMTLS(nil, tls.RequestClientCert),
		bahamut.OptErrorTransformer(errorTransformer),
		bahamut.OptIdentifiableRetriever(bootstrap.MakeIdentifiableRetriever(m, api.Manager())),
	)

	if cfg.GWTopic != "" {

		gwAnnouncedAddress := cfg.GWAnnouncedAddress
		if gwAnnouncedAddress == "" {
			gwAnnouncedAddress = getNotifierEndpoint(cfg.ListenAddress)
		}

		opts = append(
			opts,
			bootstrap.MakeBahamutGatewayNotifier(
				ctx,
				pubsub,
				"a3s",
				cfg.GWTopic,
				gwAnnouncedAddress,
				gwpush.OptionNotifierPrefix(cfg.GWAnnouncePrefix),
				gwpush.OptionNotifierPrivateAPIOverrides(cfg.GWPrivateOverrides()),
			)...,
		)

		zap.L().Info(
			"Gateway announcement configured",
			zap.String("address", gwAnnouncedAddress),
			zap.String("topic", cfg.GWTopic),
			zap.String("prefix", cfg.GWAnnouncePrefix),
			zap.Any("overrides", cfg.GWOverridePrivate),
		)
	}

	certData, err := os.ReadFile(cfg.APIServerConf.TLSCertificate)
	bmanipPool := x509.NewCertPool()
	bmanipPool.AppendCertsFromPEM(certData)
	if err != nil {
		zap.L().Fatal("Unable to read server TLS certificate", zap.Error(err))
	}
	bmanipMaker := bearermanip.Configure(
		ctx,
		publicAPIURL,
		&tls.Config{
			RootCAs: bmanipPool,
		},
	)

	server := bahamut.New(opts...)

	if err := server.RegisterCustomRouteHandler("/.well-known/jwks.json", makeJWKSHandler(jwks)); err != nil {
		zap.L().Fatal("Unable to install jwks handler", zap.Error(err))
	}

	if err := server.RegisterCustomRouteHandler("/ui/login.html", makeUILoginHandler(publicAPIURL)); err != nil {
		zap.L().Fatal("Unable to install UI login handler", zap.Error(err))
	}

	// Reusing `makeUILoginHandler` since we are serving the same html file. The UI will render the content based on the URL.
	if err := server.RegisterCustomRouteHandler("/ui/request.html", makeUILoginHandler(publicAPIURL)); err != nil {
		zap.L().Fatal("Unable to install UI request handler", zap.Error(err))
	}

	bahamut.RegisterProcessorOrDie(server,
		processors.NewIssueProcessor(
			m,
			jwks,
			cfg.JWT.JWTDefaultValidity,
			cfg.JWT.JWTMaxValidity,
			cfg.JWT.JWTIssuer,
			cfg.JWT.JWTAudience,
			cookiePolicy,
			cookieDomain,
			cfg.MTLSHeader.Enabled,
			cfg.MTLSHeader.HeaderKey,
			cfg.MTLSHeader.Passphrase,
		),
		api.IssueIdentity,
	)
	bahamut.RegisterProcessorOrDie(server, processors.NewMTLSSourcesProcessor(m), api.MTLSSourceIdentity)
	bahamut.RegisterProcessorOrDie(server, processors.NewLDAPSourcesProcessor(m), api.LDAPSourceIdentity)
	bahamut.RegisterProcessorOrDie(server, processors.NewOIDCSourcesProcessor(m), api.OIDCSourceIdentity)
	bahamut.RegisterProcessorOrDie(server, processors.NewHTTPSourcesProcessor(m), api.HTTPSourceIdentity)
	bahamut.RegisterProcessorOrDie(server, processors.NewA3SSourcesProcessor(m), api.A3SSourceIdentity)
	bahamut.RegisterProcessorOrDie(server, processors.NewPermissionsProcessor(retriever), api.PermissionsIdentity)
	bahamut.RegisterProcessorOrDie(server, processors.NewAuthzProcessor(pauthz, jwks, cfg.JWT.JWTIssuer, cfg.JWT.JWTAudience), api.AuthzIdentity)
	bahamut.RegisterProcessorOrDie(server, processors.NewNamespacesProcessor(m, pubsub), api.NamespaceIdentity)
	bahamut.RegisterProcessorOrDie(server, processors.NewAuthorizationProcessor(m, pubsub, retriever, cfg.JWT.JWTIssuer), api.AuthorizationIdentity)
	bahamut.RegisterProcessorOrDie(server, processors.NewImportProcessor(bmanipMaker, pauthz), api.ImportIdentity)

	notification.Subscribe(ctx, pubsub, nscache.NotificationNamespaceChanges, makeNamespaceCleaner(ctx, m))

	server.Run(ctx)
}

func createMongoDBAccount(cfg conf.MongoConf, username string) error {

	m := bootstrap.MakeMongoManipulator(cfg, &hasher.Hasher{})

	db, close, _ := manipmongo.GetDatabase(m)
	defer close()

	user := mgo.User{
		Username: username,
		OtherDBRoles: map[string][]mgo.Role{
			"a3s": {mgo.RoleReadWrite, mgo.RoleDBAdmin},
		},
	}

	if err := db.UpsertUser(&user); err != nil {
		return fmt.Errorf("unable to upsert the user: %w", err)
	}

	zap.L().Info("Successfully created mongodb account", zap.String("user", username))

	return nil
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

func initRootPermissions(ctx context.Context, m manipulate.Manipulator, caPath string, issuer string, ifNeeded bool) (bool, error) {

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
	source.Description = "Auth source to authenticate root users"
	source.CA = string(caData)
	certs, err := tglib.ParseCertificates([]byte(source.CA))
	if err != nil {
		return false, err
	}
	source.Fingerprints = make([]string, len(certs))
	source.SubjectKeyIDs = make([]string, len(certs))
	for i, cert := range certs {
		source.Fingerprints[i] = token.Fingerprint(cert)
		source.SubjectKeyIDs[i] = fmt.Sprintf("%02X", cert.SubjectKeyId)
	}
	if err := m.Create(manipulate.NewContext(ctx), source); err != nil {
		if errors.As(err, &manipulate.ErrConstraintViolation{}) && ifNeeded {
			return false, nil
		}
		return false, fmt.Errorf("unable to create root mtls auth source: %w", err)
	}

	auth := api.NewAuthorization()
	auth.Namespace = "/"
	auth.Name = "root-mtls-authorization"
	auth.Description = "Authorization to allow root users"
	auth.TrustedIssuers = []string{issuer}
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

func initPlatformPermissions(ctx context.Context, m manipulate.Manipulator, caPath string, issuer string, ifNeeded bool) (bool, error) {

	caData, err := os.ReadFile(caPath)
	if err != nil {
		return false, fmt.Errorf("unable to read platform ca: %w", err)
	}

	caCerts, err := tglib.ParseCertificates(caData)
	if err != nil {
		return false, fmt.Errorf("unable to parse platform ca: %w", err)
	}

	chain := make([]string, len(caCerts))
	for i, cert := range caCerts {
		chain[i] = token.Fingerprint(cert)
	}

	source := api.NewMTLSSource()
	source.Namespace = "/"
	source.Name = "platform"
	source.Description = "Auth source used to authenticate internal platform services"
	source.CA = string(caData)
	certs, err := tglib.ParseCertificates([]byte(source.CA))
	if err != nil {
		return false, err
	}
	source.Fingerprints = make([]string, len(certs))
	source.SubjectKeyIDs = make([]string, len(certs))
	for i, cert := range certs {
		source.Fingerprints[i] = token.Fingerprint(cert)
		source.SubjectKeyIDs[i] = fmt.Sprintf("%02X", cert.SubjectKeyId)
	}

	if err := m.Create(manipulate.NewContext(ctx), source); err != nil {
		if errors.As(err, &manipulate.ErrConstraintViolation{}) && ifNeeded {
			return false, nil
		}
		return false, fmt.Errorf("unable to create platform mtls auth source: %w", err)
	}

	auth := api.NewAuthorization()
	auth.Namespace = "/"
	auth.Name = "platform-mtls-authorization"
	auth.Description = "Authorization to allow internal services"
	auth.TrustedIssuers = []string{issuer}
	auth.Subject = [][]string{
		{
			"@source:type=mtls",
			"@source:name=platform",
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

func initData(ctx context.Context, m manipulate.Manipulator, dataPath string) (bool, error) {

	data, err := os.ReadFile(dataPath)
	if err != nil {
		return false, fmt.Errorf("unable to read init import file: %w", err)
	}

	importFile := api.NewImport()
	if err := yaml.Unmarshal(data, importFile); err != nil {
		return false, fmt.Errorf("unable to unmarshal import file: %w", err)
	}

	values := []elemental.Identifiables{
		importFile.LDAPSources,
		importFile.OIDCSources,
		importFile.A3SSources,
		importFile.MTLSSources,
		importFile.HTTPSources,
		importFile.Authorizations,
	}

	for _, lst := range values {
		for i, o := range lst.List() {
			if o.(elemental.Namespaceable).GetNamespace() == "" {
				return false, fmt.Errorf(
					"missing namespace property for object '%s' at index %d ",
					lst.Identity().Name,
					i,
				)
			}
		}
	}

	for _, lst := range values {

		if len(lst.List()) == 0 {
			continue
		}

		if err := importing.Import(
			ctx,
			api.Manager(),
			m,
			"/",
			"a3s:init:data",
			lst,
			false,
		); err != nil {
			return false, fmt.Errorf("unable to import '%s': %w", lst.Identity().Name, err)
		}
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

		q := req.URL.Query()

		redirect := q.Get("redirect")
		audience := q.Get("audience")

		if proxy := q.Get("proxy"); proxy != "" {
			api = proxy
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
