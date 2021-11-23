package authlib

import (
	"context"
	"encoding/json"

	"go.aporeto.io/a3s/pkgs/api"
	"go.aporeto.io/a3s/pkgs/authlib/internal/providers"
	"go.aporeto.io/manipulate"
)

// A Client allows to interract with a midgard server.
type Client struct {
	manipulator manipulate.Manipulator
}

// NewClient returns a new Client.
func NewClient(m manipulate.Manipulator) *Client {

	return &Client{
		manipulator: m,
	}
}

// AuthFromCertificate requests an identity token from the currently configured Certificate in the tls config of the client.
func (a *Client) AuthFromCertificate(ctx context.Context, sourceNamespace string, sourceName string, options ...Option) (string, error) {

	cfg := newConfig()
	for _, opt := range options {
		opt(&cfg)
	}

	req := api.NewIssue()
	req.SourceType = api.IssueSourceTypeMTLS
	req.SourceName = sourceName
	req.SourceNamespace = sourceNamespace

	applyOptions(req, cfg)

	return a.sendRequest(ctx, req)
}

// AuthFromLDAP requests a token using the provided credentials using LDAP auth source with the given name and namespace.
func (a *Client) AuthFromLDAP(ctx context.Context, username string, password string, sourceNamespace string, sourceName string, options ...Option) (string, error) {

	cfg := newConfig()
	for _, opt := range options {
		opt(&cfg)
	}

	req := api.NewIssue()
	req.SourceType = api.IssueSourceTypeLDAP
	req.SourceName = sourceName
	req.SourceNamespace = sourceNamespace
	req.InputLDAP = &api.IssueLDAP{
		Username: username,
		Password: password,
	}

	applyOptions(req, cfg)

	return a.sendRequest(ctx, req)
}

// AuthFromA3S requests a token using the provided local a3s token.
func (a *Client) AuthFromA3S(ctx context.Context, token string, options ...Option) (string, error) {

	cfg := newConfig()
	for _, opt := range options {
		opt(&cfg)
	}

	req := api.NewIssue()
	req.SourceType = api.IssueSourceTypeA3S
	req.InputA3S = &api.IssueA3S{
		Token: token,
	}

	applyOptions(req, cfg)

	return a.sendRequest(ctx, req)
}

// AuthFromRemoteA3S requests a token using the provided remote a3s token.
func (a *Client) AuthFromRemoteA3S(ctx context.Context, token string, sourceNamespace string, sourceName string, options ...Option) (string, error) {

	cfg := newConfig()
	for _, opt := range options {
		opt(&cfg)
	}

	req := api.NewIssue()
	req.SourceType = api.IssueSourceTypeRemoteA3S
	req.SourceNamespace = sourceNamespace
	req.SourceName = sourceName

	req.InputRemoteA3S = &api.IssueRemoteA3S{
		Token: token,
	}

	applyOptions(req, cfg)

	return a.sendRequest(ctx, req)
}

// AuthFromAWS requests a token using the provided AWS sts information. If accessKeyID, secretAccessKey and token are empty,
// the function will assume it is running on an AWS instance and will try to retrieve them using the magic IP.
func (a *Client) AuthFromAWS(ctx context.Context, accessKeyID, secretAccessKey, token string, options ...Option) (string, error) {

	cfg := newConfig()
	for _, opt := range options {
		opt(&cfg)
	}

	s := &struct {
		AccessKeyID     string `json:"AccessKeyId"`
		SecretAccessKey string
		Token           string
	}{}

	if accessKeyID == "" && secretAccessKey == "" && token == "" {
		awsToken, err := providers.AWSServiceRoleToken()
		if err != nil {
			return "", err
		}
		if err := json.Unmarshal([]byte(awsToken), &s); err != nil {
			return "", err
		}
	} else {
		s.AccessKeyID = accessKeyID
		s.SecretAccessKey = secretAccessKey
		s.Token = token
	}

	req := api.NewIssue()
	req.SourceType = api.IssueSourceTypeAWS
	req.InputAWS = &api.IssueAWS{
		ID:     s.AccessKeyID,
		Secret: s.SecretAccessKey,
		Token:  s.Token,
	}

	applyOptions(req, cfg)

	return a.sendRequest(ctx, req)
}

// AuthFromGCP requests a token using the provided GCP token. If token is empty,
// the function will assume it is running on an GCP instance and will try to retrieve it using the magic IP.
func (a *Client) AuthFromGCP(ctx context.Context, token string, audience string, options ...Option) (string, error) {

	cfg := newConfig()
	for _, opt := range options {
		opt(&cfg)
	}

	var err error

	if token == "" {
		token, err = providers.GCPServiceAccountToken(ctx, audience)
		if err != nil {
			return "", err
		}
	}

	req := api.NewIssue()
	req.SourceType = api.IssueSourceTypeGCP
	req.InputGCP = &api.IssueGCP{
		Token:    token,
		Audience: audience,
	}

	applyOptions(req, cfg)

	return a.sendRequest(ctx, req)
}

// AuthFromAzure requests a token using the provided Azure token. If token is empty,
// the function will assume it is running on an Azure instance and will try to retrieve it using the magic IP.
func (a *Client) AuthFromAzure(ctx context.Context, token string, options ...Option) (string, error) {

	var err error

	if token == "" {
		token, err = providers.AzureServiceIdentityToken()
		if err != nil {
			return "", err
		}
	}

	cfg := newConfig()
	for _, opt := range options {
		opt(&cfg)
	}

	req := api.NewIssue()
	req.SourceType = api.IssueSourceTypeAzure
	req.InputAzure = &api.IssueAzure{
		Token: token,
	}

	applyOptions(req, cfg)

	return a.sendRequest(ctx, req)
}

// AuthFromOIDCStep1 performs the first step of the OIDC ceremony using the configured OIDC auth source identified by
// its name and namespace. The functiion will return the provider URL to use to autenticate.
func (a *Client) AuthFromOIDCStep1(ctx context.Context, sourceNamespace string, sourceName string, redirectURL string) (string, error) {

	req := api.NewIssue()
	req.SourceType = api.IssueSourceTypeOIDC
	req.SourceNamespace = sourceNamespace
	req.SourceName = sourceName
	req.InputOIDC = &api.IssueOIDC{
		RedirectURL:    redirectURL,
		NoAuthRedirect: true,
	}

	_, err := a.sendRequest(ctx, req)
	if err != nil {
		return "", err
	}

	return req.InputOIDC.AuthURL, nil
}

// AuthFromOIDCStep2 finishes the OIDC ceremony using the code and state you obtained after performing
// the authentication against the OIDC provider.
func (a *Client) AuthFromOIDCStep2(ctx context.Context, code string, state string, options ...Option) (string, error) {

	cfg := newConfig()
	for _, opt := range options {
		opt(&cfg)
	}

	req := api.NewIssue()
	req.SourceType = api.IssueSourceTypeOIDC
	req.InputOIDC = &api.IssueOIDC{
		Code:  code,
		State: state,
	}

	applyOptions(req, cfg)

	return a.sendRequest(ctx, req)
}

func (a *Client) sendRequest(ctx context.Context, req *api.Issue) (string, error) {

	if err := a.manipulator.Create(manipulate.NewContext(ctx), req); err != nil {
		return "", err
	}

	return req.Token, nil
}

func applyOptions(req *api.Issue, cfg config) {

	if cfg.validity != 0 {
		req.Validity = cfg.validity.String()
	}
	req.Cloak = cfg.cloak
	req.Opaque = cfg.opaque
	req.Audience = cfg.audience
	req.RestrictedPermissions = cfg.restrictions.Permissions
	req.RestrictedNamespace = cfg.restrictions.Namespace
	req.RestrictedNetworks = cfg.restrictions.Networks
}
