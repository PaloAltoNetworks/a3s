package processors

import (
	"context"
	"crypto/tls"
	"crypto/x509"
	"fmt"
	"net/http"
	"time"

	"github.com/golang-jwt/jwt/v4"
	"go.aporeto.io/a3s/internal/issuer"
	"go.aporeto.io/a3s/pkgs/api"
	"go.aporeto.io/a3s/pkgs/permissions"
	"go.aporeto.io/a3s/pkgs/token"
	"go.aporeto.io/bahamut"
	"go.aporeto.io/elemental"
	"go.aporeto.io/manipulate"
)

// A IssueProcessor is a bahamut processor for Issue.
type IssueProcessor struct {
	manipulator manipulate.Manipulator
	jwks        *token.JWKS
	maxValidity time.Duration
	audience    string
	issuer      string
}

// NewIssueProcessor returns a new IssueProcessor.
func NewIssueProcessor(manipulator manipulate.Manipulator, jwks *token.JWKS, maxValidity time.Duration, issuer string, audience string) *IssueProcessor {

	return &IssueProcessor{
		manipulator: manipulator,
		jwks:        jwks,
		maxValidity: maxValidity,
		issuer:      issuer,
		audience:    audience,
	}
}

// ProcessCreate handles the creates requests for Issue.
func (p *IssueProcessor) ProcessCreate(bctx bahamut.Context) (err error) {

	req := bctx.InputData().(*api.Issue)
	validity, _ := time.ParseDuration(req.Validity) // elemental already validated this

	var issuer token.Issuer

	switch req.SourceType {

	case api.IssueSourceTypeMTLS:
		if issuer, err = p.handleCertificateIssue(bctx.Context(), req, bctx.Request().TLSConnectionState); err != nil {
			return err
		}

	case api.IssueSourceTypeA3SIdentityToken:
		if issuer, err = p.handleTokenIssue(bctx.Context(), req); err != nil {
			return err
		}
	}

	idt := issuer.Issue()
	k := p.jwks.GetLast()

	if req.Token, err = idt.JWT(
		k.PrivateKey(),
		k.KID,
		p.issuer,
		append(jwt.ClaimStrings{p.audience}, req.Audience...),
		time.Now().Add(validity),
	); err != nil {
		return err
	}

	bctx.SetOutputData(req)

	return nil
}

func (p *IssueProcessor) handleCertificateIssue(ctx context.Context, req *api.Issue, tlsState *tls.ConnectionState) (token.Issuer, error) {

	if tlsState == nil || len(tlsState.PeerCertificates) == 0 {
		return nil, elemental.NewError("Bad Request", "No client certificates", "a3s", http.StatusBadRequest)
	}

	out, err := retrieveSource(ctx, p.manipulator, req.SourceNamespace, req.SourceName, api.MTLSSourceIdentity)
	if err != nil {
		return nil, err
	}
	src := out.(*api.MTLSSource)

	pool := x509.NewCertPool()
	pool.AppendCertsFromPEM([]byte(src.CertificateAuthority))

	userCert := tlsState.PeerCertificates[0]
	iss := issuer.NewMTLSIssuer(pool, req.SourceNamespace, req.SourceName)
	if err := iss.FromCertificate(userCert); err != nil {
		return nil, err
	}

	return iss, nil
}

func (p *IssueProcessor) handleTokenIssue(ctx context.Context, req *api.Issue) (token.Issuer, error) {

	iss := issuer.NewTokenIssuer()

	token, ok := req.Metadata["token"].(string)
	if !ok || token == "" {
		return nil, elemental.NewError(
			"Bad Request",
			"This source needs the token to be passed as metadata key 'token'",
			"a3s:authn",
			http.StatusBadRequest,
		)
	}

	if err := iss.FromToken(
		token,
		p.jwks,
		p.issuer,
		p.audience,
		req.Validity,
		permissions.Restrictions{
			Namespace:   req.RestrictedNamespace,
			Networks:    req.RestrictedNetworks,
			Permissions: req.RestrictedPermissions,
		},
	); err != nil {
		return nil, err
	}

	return iss, nil
}

func retrieveSource(ctx context.Context, m manipulate.Manipulator, namespace string, name string, identity elemental.Identity) (elemental.Identifiable, error) {

	if namespace == "" {
		return nil, elemental.NewError(
			"Bad Request",
			"You must set sourceNamespace and sourceName",
			"a3s:auth",
			http.StatusBadRequest,
		)
	}

	if name == "" {
		return nil, elemental.NewError(
			"Bad Request",
			"You must set sourceNamespace and sourceName",
			"a3s:auth",
			http.StatusBadRequest,
		)
	}

	mctx := manipulate.NewContext(ctx,
		manipulate.ContextOptionNamespace(namespace),
		manipulate.ContextOptionFilter(
			elemental.NewFilterComposer().WithKey("name").Equals(name).
				Done(),
		),
	)

	identifiables := api.Manager().IdentifiablesFromString(identity.Name)
	if err := m.RetrieveMany(mctx, identifiables); err != nil {
		return nil, err
	}

	lst := identifiables.List()
	switch len(lst) {
	case 0:
		return nil, elemental.NewError(
			"Not Found",
			"Unable to find the request auth source",
			"a3s:authn",
			http.StatusNotFound,
		)
	case 1:
	default:
		return nil, fmt.Errorf("more than one auth source found")
	}

	return lst[0], nil
}
