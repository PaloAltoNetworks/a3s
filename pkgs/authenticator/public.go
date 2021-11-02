package authenticator

import "go.aporeto.io/bahamut"

// A Public can be used to define
// public resources. If the authenticator does not
// know the resource it will return bahamut.AuthActionContinue
// to let the rest of authenticator chain decide.
type Public struct {
	publicResources map[string]struct{}
}

// NewPublic returns a new PublicAuthenticator.
func NewPublic(resources ...string) *Public {

	m := make(map[string]struct{}, len(resources))
	for _, r := range resources {
		m[r] = struct{}{}
	}

	return &Public{
		publicResources: m,
	}
}

// AuthenticateSession authenticates the given session.
func (a *Public) AuthenticateSession(session bahamut.Session) (bahamut.AuthAction, error) {

	return bahamut.AuthActionContinue, nil
}

// AuthenticateRequest authenticates the request from the given bahamut.Context.
func (a *Public) AuthenticateRequest(bctx bahamut.Context) (bahamut.AuthAction, error) {

	if _, ok := a.publicResources[bctx.Request().Identity.Name]; ok {
		return bahamut.AuthActionOK, nil
	}

	return bahamut.AuthActionContinue, nil
}
