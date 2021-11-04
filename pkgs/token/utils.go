package token

import (
	"go.aporeto.io/bahamut"
	"go.aporeto.io/elemental"
)

// FromRequest retrieves the token from the given elemental.Request
// first looking at the cookie x-a3s-token, then the request.Password.
func FromRequest(req *elemental.Request) string {
	if hreq := req.HTTPRequest(); hreq != nil {
		if cookie, err := hreq.Cookie("x-a3s-token"); err == nil {
			return cookie.Value
		}
	}
	return req.Password
}

// FromSession retrieves the token from the given bahamut.Session
// first looking at the cookie x-a3s-token, then the session.Token(.
func FromSession(session bahamut.Session) string {
	if cookie, err := session.Cookie("x-a3s-token"); err == nil {
		return cookie.Value
	}
	return session.Token()
}
