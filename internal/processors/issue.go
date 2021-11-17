package processors

import (
	"context"
	"crypto/tls"
	"fmt"
	"net/http"
	"time"

	"github.com/golang-jwt/jwt/v4"
	"go.aporeto.io/a3s/internal/issuer/a3sissuer"
	"go.aporeto.io/a3s/internal/issuer/awsissuer"
	"go.aporeto.io/a3s/internal/issuer/azureissuer"
	"go.aporeto.io/a3s/internal/issuer/gcpissuer"
	"go.aporeto.io/a3s/internal/issuer/ldapissuer"
	"go.aporeto.io/a3s/internal/issuer/mtlsissuer"
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
	exp := time.Now().Add(validity)

	audience := req.Audience
	if len(audience) == 0 {
		audience = jwt.ClaimStrings{p.audience}
	}

	var issuer token.Issuer

	switch req.SourceType {

	case api.IssueSourceTypeMTLS:
		issuer, err = p.handleCertificateIssue(bctx.Context(), req, bctx.Request().TLSConnectionState)

	case api.IssueSourceTypeLDAP:
		issuer, err = p.handleLDAPIssue(bctx.Context(), req)

	case api.IssueSourceTypeAWSSecurityToken:
		issuer, err = p.handleAWSIssue(bctx.Context(), req)

	case api.IssueSourceTypeAzureIdentityToken:
		issuer, err = p.handleAzureIssue(bctx.Context(), req)

	case api.IssueSourceTypeGCPIdentityToken:
		issuer, err = p.handleGCPIssue(bctx.Context(), req)

	case api.IssueSourceTypeA3SIdentityToken:
		issuer, err = p.handleTokenIssue(bctx.Context(), req, validity)
		// we reset to 0 to skip setting exp during issuing of the token
		// as the token issers already caps it.
		exp = time.Time{}
	}

	if err != nil {
		return elemental.NewError("Unauthorized", err.Error(), "a3s:authn", http.StatusUnauthorized)
	}

	idt := issuer.Issue()
	k := p.jwks.GetLast()

	if req.Token, err = idt.JWT(k.PrivateKey(), k.KID, p.issuer, audience, exp, req.Cloak); err != nil {
		return err
	}

	req.InputLDAP = nil
	req.InputAWSSTS = nil
	req.InputToken = nil
	req.Validity = time.Until(idt.ExpiresAt.Time).Round(time.Second).String()

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

	userCert := tlsState.PeerCertificates[0]
	iss, err := mtlsissuer.New(src, userCert)
	if err != nil {
		return nil, err
	}

	return iss, nil
}

func (p *IssueProcessor) handleLDAPIssue(ctx context.Context, req *api.Issue) (token.Issuer, error) {

	out, err := retrieveSource(ctx, p.manipulator, req.SourceNamespace, req.SourceName, api.LDAPSourceIdentity)
	if err != nil {
		return nil, err
	}

	src := out.(*api.LDAPSource)
	iss, err := ldapissuer.New(src, req.InputLDAP.Username, req.InputLDAP.Password)
	if err != nil {
		return nil, err
	}

	return iss, nil
}

func (p *IssueProcessor) handleAWSIssue(ctx context.Context, req *api.Issue) (token.Issuer, error) {

	iss, err := awsissuer.New(req.InputAWSSTS.ID, req.InputAWSSTS.Secret, req.InputAWSSTS.Token)
	if err != nil {
		return nil, err
	}

	return iss, nil
}

func (p *IssueProcessor) handleAzureIssue(ctx context.Context, req *api.Issue) (token.Issuer, error) {

	iss, err := azureissuer.New(ctx, req.InputAzure.Token)
	if err != nil {
		return nil, err
	}

	return iss, nil
}

func (p *IssueProcessor) handleGCPIssue(ctx context.Context, req *api.Issue) (token.Issuer, error) {

	iss, err := gcpissuer.New(req.InputGCP.Token, req.InputGCP.Audience)
	if err != nil {
		return nil, err
	}

	return iss, nil
}
func (p *IssueProcessor) handleTokenIssue(ctx context.Context, req *api.Issue, validity time.Duration) (token.Issuer, error) {

	iss, err := a3sissuer.New(
		req.InputToken.Token,
		p.jwks,
		p.issuer,
		p.audience,
		validity,
		permissions.Restrictions{
			Namespace:   req.RestrictedNamespace,
			Networks:    req.RestrictedNetworks,
			Permissions: req.RestrictedPermissions,
		},
	)
	if err != nil {
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
