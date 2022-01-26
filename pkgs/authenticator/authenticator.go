package authenticator

import (
	"fmt"
	"net/http"

	"go.aporeto.io/a3s/pkgs/token"
	"go.aporeto.io/bahamut"
	"go.aporeto.io/elemental"
)

// A Authenticator is a bahamut.Authenticator compliant structure to authentify
// requests using an a3s token.
type Authenticator struct {
	jwks             *token.JWKS
	issuer           string
	audience         string
	ignoredResources map[string]struct{}
}

// New returns a new Authenticator that will use the provided JWKS
// to cryptographically verify a request or session token.
// It will validate the token comes from the given issuer and has the
// correct audience.
func New(jwks *token.JWKS, issuer string, audience string, options ...Option) *Authenticator {

	cfg := config{}
	for _, o := range options {
		o(&cfg)
	}

	m := make(map[string]struct{}, len(cfg.ignoredResources))
	for _, r := range cfg.ignoredResources {
		m[r] = struct{}{}
	}

	return &Authenticator{
		jwks:             jwks,
		issuer:           issuer,
		audience:         audience,
		ignoredResources: m,
	}
}

// AuthenticateSession authenticates the given session.
func (a *Authenticator) AuthenticateSession(session bahamut.Session) (bahamut.AuthAction, error) {

	action, claims, err := a.commonAuth(token.FromSession(session))
	if err != nil {
		return bahamut.AuthActionKO, err
	}

	session.SetClaims(claims)

	return action, nil
}

// AuthenticateRequest authenticates the request from the given bahamut.Context.
func (a *Authenticator) AuthenticateRequest(bctx bahamut.Context) (bahamut.AuthAction, error) {

	if _, ok := a.ignoredResources[bctx.Request().Identity.Category]; ok {
		return bahamut.AuthActionOK, nil
	}

	token := token.FromRequest(bctx.Request())

	action, claims, err := a.commonAuth(token)
	if err != nil {
		return bahamut.AuthActionKO, err
	}

	bctx.SetClaims(claims)

	return action, nil
}

func (a *Authenticator) commonAuth(tokenString string) (bahamut.AuthAction, []string, error) {

	if tokenString == "" {
		return bahamut.AuthActionKO, nil, elemental.NewError(
			"Unauthorized",
			"Missing token in Authorization header",
			"a3s:authn",
			http.StatusUnauthorized,
		)
	}

	idt, err := token.Parse(tokenString, a.jwks, a.issuer, a.audience)
	if err != nil {
		return bahamut.AuthActionKO, nil, elemental.NewError(
			"Unauthorized",
			fmt.Sprintf("Authentication rejected with error: %s", err),
			"a3s:authn",
			http.StatusUnauthorized,
		)
	}

	if idt.Refresh {
		return bahamut.AuthActionKO, nil, elemental.NewError(
			"Unauthorized",
			fmt.Sprintf("Authentication impossible from a refresh token"),
			"a3s:authn",
			http.StatusUnauthorized,
		)
	}

	return bahamut.AuthActionContinue, idt.Identity, nil
}
