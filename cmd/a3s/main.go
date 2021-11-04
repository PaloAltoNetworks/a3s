package main

import (
	"context"
	"crypto/sha1"
	"crypto/tls"
	"errors"
	"fmt"
	"net/http"
	"os"
	"strings"

	"go.aporeto.io/a3s/internal/srv/authn"
	"go.aporeto.io/a3s/internal/srv/authz"
	"go.aporeto.io/a3s/internal/srv/policy"
	"go.aporeto.io/a3s/pkgs/api"
	"go.aporeto.io/a3s/pkgs/authenticator"
	"go.aporeto.io/a3s/pkgs/authorizer"
	"go.aporeto.io/a3s/pkgs/bootstrap"
	"go.aporeto.io/a3s/pkgs/indexes"
	"go.aporeto.io/a3s/pkgs/permissions"
	"go.aporeto.io/bahamut"
	"go.aporeto.io/elemental"
	"go.aporeto.io/manipulate"
	"go.aporeto.io/tg/tglib"
	"go.uber.org/zap"
)

var (
	publicResources = []string{
		api.IssueIdentity.Category,
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

	manipulator := bootstrap.MakeMongoManipulator(cfg.MongoConf)
	if err := indexes.Ensure(manipulator, api.Manager(), "a3s"); err != nil {
		zap.L().Fatal("Unable to ensure indexes", zap.Error(err))
	}

	if err := createRootNamespaceIfNeeded(manipulator); err != nil {
		zap.L().Fatal("Unable to handle root namespace", zap.Error(err))
	}

	if cfg.Init {
		if err := initRootPermissions(ctx, manipulator, cfg.InitRootUserCAPath); err != nil {
			zap.L().Fatal("unable to initialize root permissions", zap.Error(err))
			return
		}

		zap.L().Info("Root permissions initialized")
		return
	}

	jwtCert, _, err := cfg.AuthNConf.JWTCertificate()
	if err != nil {
		zap.L().Fatal("Unable to get JWT certificate", zap.Error(err))
	}

	pubsub := bootstrap.MakeNATSClient(cfg.NATSConf)
	defer pubsub.Disconnect() // nolint: errcheck

	pauthn := authenticator.NewPrivate(jwtCert)
	retriever := permissions.NewRetriever(manipulator)
	pauthz := authorizer.New(
		ctx,
		retriever,
		pubsub,
		authorizer.OptionIgnoredResources(publicResources...),
	)

	server := bahamut.New(
		append(
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
			bahamut.OptMTLS(nil, tls.RequestClientCert),
			bahamut.OptErrorTransformer(errorTransformer),
			bahamut.OptIdentifiableRetriever(bootstrap.MakeIdentifiableRetriever(manipulator)),
		)...,
	)

	if err := authn.Init(ctx, cfg.AuthNConf, server, manipulator, pubsub); err != nil {
		zap.L().Fatal("Unable to initialize authn module", zap.Error(err))
	}

	if err := policy.Init(ctx, cfg.PolicyConf, server, manipulator, retriever, pubsub); err != nil {
		zap.L().Fatal("Unable to initialize policy module", zap.Error(err))
	}

	if err := authz.Init(ctx, server, manipulator, retriever); err != nil {
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

func initRootPermissions(ctx context.Context, m manipulate.Manipulator, caPath string) error {

	caData, err := os.ReadFile(caPath)
	if err != nil {
		return fmt.Errorf("unable to read root user ca: %w", err)
	}

	caCerts, err := tglib.ParseCertificates(caData)
	if err != nil {
		return fmt.Errorf("unable to parse root user ca: %w", err)
	}

	chain := make([]string, len(caCerts))
	for i, cert := range caCerts {
		chain[i] = fmt.Sprintf("%02X", sha1.Sum(cert.Raw))
	}

	source := api.NewMTLSSource()
	source.Namespace = "/"
	source.Name = "root"
	source.Description = "Root auth source used to bootstrap permissions."
	source.CertificateAuthority = string(caData)
	if err := m.Create(manipulate.NewContext(ctx), source); err != nil {
		return fmt.Errorf("unable to create root mtls auth source: %w", err)
	}

	auth := api.NewAuthorization()
	auth.Namespace = "/"
	auth.Name = "root-mtls-authorization"
	auth.Description = "Root authorization for certificates issued from the CA declared  in the root auth mtls source."
	auth.Subject = [][]string{
		{
			"@sourcetype=mtls",
			"@sourcename=root",
			"@sourcenamespace=/",
			fmt.Sprintf("issuerchain=%s", strings.Join(chain, ",")),
		},
	}
	auth.FlattenedSubject = auth.Subject[0]
	auth.Permissions = []string{"*,*"}
	auth.TargetNamespace = "/"
	auth.Hidden = true

	if err := m.Create(manipulate.NewContext(ctx), auth); err != nil {
		return fmt.Errorf("unable to create root auth: %w", err)
	}

	return nil
}
