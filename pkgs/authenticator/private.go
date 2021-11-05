package authenticator

import (
	"fmt"
	"net/http"

	"go.aporeto.io/a3s/pkgs/token"
	"go.aporeto.io/bahamut"
	"go.aporeto.io/elemental"
)

// A Private is a bahamut.Authenticator compliant structure to authentify
// requests using a a3s token.
type Private struct {
	jwks *token.JWKS
}

// NewPrivate returns a new *Authenticator that will make a call
func NewPrivate(jwks *token.JWKS) *Private {
	return &Private{
		jwks: jwks,
	}
}

// AuthenticateSession authenticates the given session.
func (a *Private) AuthenticateSession(session bahamut.Session) (bahamut.AuthAction, error) {

	action, claims, err := a.commonAuth(token.FromSession(session))
	if err != nil {
		return bahamut.AuthActionKO, err
	}

	session.SetClaims(claims)

	return action, nil
}

// AuthenticateRequest authenticates the request from the given bahamut.Context.
func (a *Private) AuthenticateRequest(bctx bahamut.Context) (bahamut.AuthAction, error) {

	token := token.FromRequest(bctx.Request())

	action, claims, err := a.commonAuth(token)
	if err != nil {
		return bahamut.AuthActionKO, err
	}

	bctx.SetClaims(claims)

	return action, nil
}

func (a *Private) commonAuth(tokenString string) (bahamut.AuthAction, []string, error) {

	if tokenString == "" {
		return bahamut.AuthActionKO, nil, elemental.NewError(
			"Unauthorized",
			"Missing token in Authorization header",
			"a3s:authn",
			http.StatusUnauthorized,
		)
	}

	idt := &token.IdentityToken{}
	if err := idt.Parse(tokenString, a.jwks, "", ""); err != nil {
		return bahamut.AuthActionKO, nil, elemental.NewError(
			"Unauthorized",
			fmt.Sprintf("Authentication rejected with error: %s", err),
			"a3s:authn",
			http.StatusUnauthorized,
		)
	}

	return bahamut.AuthActionContinue, idt.Identity, nil
}
