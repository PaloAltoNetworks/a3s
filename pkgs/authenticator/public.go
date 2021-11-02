package authenticator

import "go.aporeto.io/bahamut"

// A PublicAuthenticator can be used to define
// public resources. If the authenticator does not
// know the resource it will return bahamut.AuthActionContinue
// to let the rest of authenticator chain decide.
type PublicAuthenticator struct {
	publicResources map[string]struct{}
}

// NewPublicAuthenticator returns a new PublicAuthenticator.
func NewPublicAuthenticator(resources ...string) *PublicAuthenticator {

	m := make(map[string]struct{}, len(resources))
	for _, r := range resources {
		m[r] = struct{}{}
	}

	return &PublicAuthenticator{
		publicResources: m,
	}
}

// AuthenticateSession authenticates the given session.
func (a *PublicAuthenticator) AuthenticateSession(session bahamut.Session) (bahamut.AuthAction, error) {

	return bahamut.AuthActionContinue, nil
}

// AuthenticateRequest authenticates the request from the given bahamut.Context.
func (a *PublicAuthenticator) AuthenticateRequest(bctx bahamut.Context) (bahamut.AuthAction, error) {

	if _, ok := a.publicResources[bctx.Request().Identity.Name]; ok {
		return bahamut.AuthActionOK, nil
	}

	return bahamut.AuthActionContinue, nil
}
