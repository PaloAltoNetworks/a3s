package processors

import (
	"context"
	"crypto/tls"
	"fmt"
	"net/http"
	"time"

	"github.com/coreos/go-oidc/v3/oidc"
	"github.com/golang-jwt/jwt/v4"
	"go.aporeto.io/a3s/internal/issuer/a3sissuer"
	"go.aporeto.io/a3s/internal/issuer/awsissuer"
	"go.aporeto.io/a3s/internal/issuer/azureissuer"
	"go.aporeto.io/a3s/internal/issuer/gcpissuer"
	"go.aporeto.io/a3s/internal/issuer/httpissuer"
	"go.aporeto.io/a3s/internal/issuer/ldapissuer"
	"go.aporeto.io/a3s/internal/issuer/mtlsissuer"
	"go.aporeto.io/a3s/internal/issuer/oidcissuer"
	"go.aporeto.io/a3s/internal/issuer/remotea3sissuer"
	"go.aporeto.io/a3s/internal/oidcceremony"
	"go.aporeto.io/a3s/pkgs/api"
	"go.aporeto.io/a3s/pkgs/permissions"
	"go.aporeto.io/a3s/pkgs/token"
	"go.aporeto.io/bahamut"
	"go.aporeto.io/bahamut/authorizer/mtls"
	"go.aporeto.io/elemental"
	"go.aporeto.io/manipulate"
	"golang.org/x/oauth2"
)

// A IssueProcessor is a bahamut processor for Issue.
type IssueProcessor struct {
	manipulator          manipulate.Manipulator
	jwks                 *token.JWKS
	maxValidity          time.Duration
	audience             string
	cookieSameSitePolicy http.SameSite
	cookieDomain         string
	issuer               string
	mtlsHeaderEnabled    bool
	mtlsHeaderKey        string
	mtlsHeaderPass       string
}

// NewIssueProcessor returns a new IssueProcessor.
func NewIssueProcessor(
	manipulator manipulate.Manipulator,
	jwks *token.JWKS,
	maxValidity time.Duration,
	issuer string,
	audience string,
	cookieSameSitePolicy http.SameSite,
	cookieDomain string,
	mtlsHeaderEnabled bool,
	mtlsHeaderKey string,
	mtlsHeaderPass string,
) *IssueProcessor {

	return &IssueProcessor{
		manipulator:          manipulator,
		jwks:                 jwks,
		maxValidity:          maxValidity,
		issuer:               issuer,
		audience:             audience,
		cookieSameSitePolicy: cookieSameSitePolicy,
		cookieDomain:         cookieDomain,
		mtlsHeaderEnabled:    mtlsHeaderEnabled,
		mtlsHeaderKey:        mtlsHeaderKey,
		mtlsHeaderPass:       mtlsHeaderPass,
	}
}

// ProcessCreate handles the creates requests for Issue.
func (p *IssueProcessor) ProcessCreate(bctx bahamut.Context) (err error) {

	req := bctx.InputData().(*api.Issue)

	validity, _ := time.ParseDuration(req.Validity) // elemental already validated this
	if validity > p.maxValidity {
		return elemental.NewError(
			"Invalid validity",
			fmt.Sprintf(
				"The requested validity '%s' is greater than the maximum allowed ('%s')",
				req.Validity,
				p.maxValidity,
			),
			"a3s:authn",
			http.StatusBadRequest,
		)
	}

	if validity == 0 {
		validity = p.maxValidity
	}

	exp := time.Now().Add(validity)

	audience := req.Audience
	if len(audience) == 0 {
		audience = jwt.ClaimStrings{p.audience}
	}

	var issuer token.Issuer

	switch req.SourceType {

	case api.IssueSourceTypeMTLS:
		issuer, err = p.handleCertificateIssue(
			bctx.Context(),
			req,
			bctx.Request().TLSConnectionState,
			bctx.Request().Headers.Get(p.mtlsHeaderKey),
		)

	case api.IssueSourceTypeLDAP:
		issuer, err = p.handleLDAPIssue(bctx.Context(), req)

	case api.IssueSourceTypeHTTP:
		issuer, err = p.handleHTTPIssue(bctx.Context(), req)

	case api.IssueSourceTypeAWS:
		issuer, err = p.handleAWSIssue(bctx.Context(), req)

	case api.IssueSourceTypeAzure:
		issuer, err = p.handleAzureIssue(bctx.Context(), req)

	case api.IssueSourceTypeGCP:
		issuer, err = p.handleGCPIssue(bctx.Context(), req)

	case api.IssueSourceTypeRemoteA3S:
		issuer, err = p.handleRemoteA3SIssue(bctx.Context(), req)

	case api.IssueSourceTypeOIDC:
		issuer, err = p.handleOIDCIssue(bctx, req)
		if issuer == nil && err == nil {
			return nil
		}

	case api.IssueSourceTypeA3S:
		if req.Validity == "" {
			validity = 0
		}
		issuer, err = p.handleTokenIssue(bctx.Context(), req, validity, audience)
		// we reset to 0 to skip setting exp during issuing of the token
		// as the token issers already caps it.
		exp = time.Time{}
	}

	if err != nil {
		return elemental.NewError("Unauthorized", err.Error(), "a3s:authn", http.StatusUnauthorized)
	}

	idt := issuer.Issue()

	if err := idt.Restrict(permissions.Restrictions{
		Namespace:   req.RestrictedNamespace,
		Networks:    req.RestrictedNetworks,
		Permissions: req.RestrictedPermissions,
	}); err != nil {
		return elemental.NewError(
			"Restrictions Error",
			err.Error(),
			"a3s:authn",
			http.StatusBadRequest,
		)
	}

	if req.TokenType == api.IssueTokenTypeRefresh {
		idt.Refresh = true
	}

	k := p.jwks.GetLast()
	tkn, err := idt.JWT(
		k.PrivateKey(),
		k.KID,
		p.issuer,
		audience,
		exp,
		req.Cloak,
	)
	if err != nil {
		return err
	}

	req.Validity = time.Until(idt.ExpiresAt.Time).Round(time.Second).String()
	req.InputLDAP = nil
	req.InputAWS = nil
	req.InputAzure = nil
	req.InputGCP = nil
	req.InputOIDC = nil
	req.InputA3S = nil
	req.InputRemoteA3S = nil

	if req.Cookie {
		domain := req.CookieDomain
		if domain == "" {
			domain = p.cookieDomain
		}
		bctx.AddOutputCookies(
			&http.Cookie{
				Name:     "x-a3s-token",
				Value:    tkn,
				HttpOnly: true,
				Secure:   true,
				Expires:  idt.ExpiresAt.Time,
				SameSite: p.cookieSameSitePolicy,
				Domain:   domain,
			},
		)
	} else {
		req.Token = tkn
	}

	bctx.SetOutputData(req)

	return nil
}

func (p *IssueProcessor) handleCertificateIssue(ctx context.Context, req *api.Issue, tlsState *tls.ConnectionState, tlsHeader string) (token.Issuer, error) {

	// We get the peer certificate.
	cert := tlsState.PeerCertificates

	// If mtls header is enabled, and the header is not empty
	// we will use it instead of the cert from the tls state.
	if p.mtlsHeaderEnabled && tlsHeader != "" {

		// First we create an elemental.AESAttributeEncrypter
		cipher, err := elemental.NewAESAttributeEncrypter(p.mtlsHeaderPass)
		if err != nil {
			return nil, fmt.Errorf("unable to build AES encrypter: %w", err)
		}

		// Then we decrypt the content of the header.
		header, err := cipher.DecryptString(tlsHeader)
		if err != nil {
			return nil, fmt.Errorf("unable to decrypt header: %w", err)
		}

		// Then we try to extract a certificate out of the decrypted blob.
		cert, err = mtls.CertificatesFromHeader(header)
		if err != nil {
			return nil, fmt.Errorf("unable to retrieve certificate from mtls header: %w", err)
		}

		// If we reach here, the decoded certificate from the header will be used to
		// to match against the source.
	}

	if len(cert) == 0 {
		return nil, elemental.NewError("Bad Request", "No user certificates", "a3s:authn", http.StatusBadRequest)
	}

	out, err := retrieveSource(ctx, p.manipulator, req.SourceNamespace, req.SourceName, api.MTLSSourceIdentity)
	if err != nil {
		return nil, err
	}
	src := out.(*api.MTLSSource)

	iss, err := mtlsissuer.New(ctx, src, cert[0])
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
	iss, err := ldapissuer.New(ctx, src, req.InputLDAP.Username, req.InputLDAP.Password)
	if err != nil {
		return nil, err
	}

	return iss, nil
}

func (p *IssueProcessor) handleHTTPIssue(ctx context.Context, req *api.Issue) (token.Issuer, error) {

	out, err := retrieveSource(ctx, p.manipulator, req.SourceNamespace, req.SourceName, api.HTTPSourceIdentity)
	if err != nil {
		return nil, err
	}

	src := out.(*api.HTTPSource)
	iss, err := httpissuer.New(ctx, src, httpissuer.Credentials{
		Username: req.InputHTTP.Username,
		Password: req.InputHTTP.Password,
		TOTP:     req.InputHTTP.TOTP,
	})
	if err != nil {
		return nil, err
	}

	return iss, nil
}

func (p *IssueProcessor) handleAWSIssue(ctx context.Context, req *api.Issue) (token.Issuer, error) {

	iss, err := awsissuer.New(req.InputAWS.ID, req.InputAWS.Secret, req.InputAWS.Token)
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

func (p *IssueProcessor) handleTokenIssue(ctx context.Context, req *api.Issue, validity time.Duration, audience []string) (token.Issuer, error) {

	iss, err := a3sissuer.New(
		req.InputA3S.Token,
		p.jwks,
		p.issuer,
		audience,
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

func (p *IssueProcessor) handleRemoteA3SIssue(ctx context.Context, req *api.Issue) (token.Issuer, error) {

	out, err := retrieveSource(ctx, p.manipulator, req.SourceNamespace, req.SourceName, api.A3SSourceIdentity)
	if err != nil {
		return nil, err
	}

	src := out.(*api.A3SSource)
	iss, err := remotea3sissuer.New(ctx, src, req.InputRemoteA3S.Token)
	if err != nil {
		return nil, err
	}

	return iss, nil
}

func (p *IssueProcessor) handleOIDCIssue(bctx bahamut.Context, req *api.Issue) (token.Issuer, error) {

	state := req.InputOIDC.State
	code := req.InputOIDC.Code

	out, err := retrieveSource(
		bctx.Context(),
		p.manipulator,
		req.SourceNamespace,
		req.SourceName,
		api.OIDCSourceIdentity,
	)
	if err != nil {
		return nil, err
	}
	src := out.(*api.OIDCSource)

	if code == "" && state == "" {

		client, err := oidcceremony.MakeOIDCProviderClient(src.CA)
		if err != nil {
			return nil, oidcceremony.RedirectErrorEventually(
				bctx,
				req.InputOIDC.RedirectErrorURL,
				elemental.NewError(
					"Bad Request",
					err.Error(),
					"a3s:authn",
					http.StatusBadRequest,
				),
			)
		}

		oidcCtx := oidc.ClientContext(bctx.Context(), client)
		provider, err := oidc.NewProvider(oidcCtx, src.Endpoint)
		if err != nil {
			return nil, oidcceremony.RedirectErrorEventually(
				bctx,
				req.InputOIDC.RedirectErrorURL,
				elemental.NewError(
					"Bad Request",
					err.Error(),
					"a3s:authn",
					http.StatusBadRequest,
				),
			)
		}

		oauth2Config := oauth2.Config{
			ClientID:     src.ClientID,
			ClientSecret: src.ClientSecret,
			RedirectURL:  req.InputOIDC.RedirectURL,
			Endpoint:     provider.Endpoint(),
			Scopes:       append([]string{oidc.ScopeOpenID}, src.Scopes...),
		}

		state, err = oidcceremony.GenerateNonce(12)
		if err != nil {
			return nil, oidcceremony.RedirectErrorEventually(
				bctx,
				req.InputOIDC.RedirectErrorURL,
				err,
			)
		}

		cacheItem := &oidcceremony.CacheItem{
			State:            state,
			ProviderEndpoint: src.Endpoint,
			CA:               src.CA,
			ClientID:         src.ClientID,
			OAuth2Config:     oauth2Config,
		}

		if err := oidcceremony.Set(p.manipulator, cacheItem); err != nil {
			return nil, oidcceremony.RedirectErrorEventually(
				bctx,
				req.InputOIDC.RedirectErrorURL,
				err,
			)
		}

		authURL := oauth2Config.AuthCodeURL(state)

		if req.InputOIDC.NoAuthRedirect {
			req.InputOIDC.AuthURL = authURL
			bctx.SetOutputData(req)
		} else {
			bctx.SetRedirect(authURL)
		}

		return nil, nil
	}

	oidcReq, err := oidcceremony.Get(p.manipulator, state)
	if err != nil {
		return nil, err
	}

	if err := oidcceremony.Delete(p.manipulator, state); err != nil {
		return nil, err
	}

	client, err := oidcceremony.MakeOIDCProviderClient(oidcReq.CA)
	if err != nil {
		return nil, fmt.Errorf("unable to create oidc http client: %s", err)
	}

	oidcctx := oidc.ClientContext(bctx.Context(), client)

	oauth2Token, err := oidcReq.OAuth2Config.Exchange(oidcctx, code)
	if err != nil {
		return nil, elemental.NewError(
			"OAuth2 Error",
			err.Error(),
			"a3s:authn",
			http.StatusNotAcceptable,
		)
	}

	rawIDToken, ok := oauth2Token.Extra("id_token").(string)
	if !ok {
		return nil, fmt.Errorf("missing ID token")
	}

	provider, err := oidc.NewProvider(oidcctx, oidcReq.ProviderEndpoint)
	if err != nil {
		return nil, elemental.NewError(
			"OIDC Error",
			err.Error(),
			"a3s:authn",
			http.StatusUnauthorized,
		)
	}

	verifier := provider.Verifier(
		&oidc.Config{
			ClientID: oidcReq.ClientID,
		},
	)

	idToken, err := verifier.Verify(oidcctx, rawIDToken)
	if err != nil {
		return nil, elemental.NewError(
			"OAuth2 Verification Error",
			err.Error(),
			"a3s:authn",
			http.StatusNotAcceptable,
		)
	}

	claims := map[string]interface{}{}
	if err := idToken.Claims(&claims); err != nil {
		return nil, elemental.NewError(
			"Claims Decoding Error",
			err.Error(),
			"a3s:authn",
			http.StatusNotAcceptable,
		)
	}

	return oidcissuer.New(bctx.Context(), src, claims)
}

func retrieveSource(
	ctx context.Context,
	m manipulate.Manipulator,
	namespace string,
	name string,
	identity elemental.Identity,
) (elemental.Identifiable, error) {

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
