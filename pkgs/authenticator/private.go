package authenticator

import (
	"crypto/x509"
	"fmt"
	"net/http"

	"go.aporeto.io/a3s/pkgs/token"
	"go.aporeto.io/bahamut"
	"go.aporeto.io/elemental"
)

// A PrivateAuthenticator is a bahamut.PrivateAuthenticator compliant structure to authentify
// requests using a a3s token.
type PrivateAuthenticator struct {
	jwtCert *x509.Certificate
}

// NewPrivateAuthenticator returns a new *Authenticator that will make a call
func NewPrivateAuthenticator(cert *x509.Certificate) *PrivateAuthenticator {
	return &PrivateAuthenticator{
		jwtCert: cert,
	}
}

// AuthenticateSession authenticates the given session.
func (a *PrivateAuthenticator) AuthenticateSession(session bahamut.Session) (bahamut.AuthAction, error) {

	action, claims, err := a.commonAuth(getSessionToken(session))
	if err != nil {
		return bahamut.AuthActionKO, err
	}

	session.SetClaims(claims)

	return action, nil
}

// AuthenticateRequest authenticates the request from the given bahamut.Context.
func (a *PrivateAuthenticator) AuthenticateRequest(bctx bahamut.Context) (bahamut.AuthAction, error) {

	token := getToken(bctx.Request())

	action, claims, err := a.commonAuth(token)
	if err != nil {
		return bahamut.AuthActionKO, err
	}

	bctx.SetClaims(claims)

	return action, nil
}

func (a *PrivateAuthenticator) commonAuth(tokenString string) (bahamut.AuthAction, []string, error) {

	if tokenString == "" {
		return bahamut.AuthActionKO, nil, elemental.NewError(
			"Unauthorized",
			"Missing token in Authorization header",
			"a3s:authn",
			http.StatusUnauthorized,
		)
	}

	mc, err := token.Verify(tokenString, a.jwtCert)
	if err != nil {
		return bahamut.AuthActionKO, nil, elemental.NewError(
			"Unauthorized",
			fmt.Sprintf("Authentication rejected with error: %s", err),
			"a3s:authn",
			http.StatusUnauthorized,
		)
	}

	return bahamut.AuthActionContinue, mc.Identity, nil
}

func getToken(req *elemental.Request) string {
	if hreq := req.HTTPRequest(); hreq != nil {
		if cookie, err := hreq.Cookie("x-a3s-token"); err == nil {
			return cookie.Value
		}
	}
	return req.Password
}

func getSessionToken(session bahamut.Session) string {
	if cookie, err := session.Cookie("x-a3s-token"); err == nil {
		return cookie.Value
	}
	return session.Token()
}
