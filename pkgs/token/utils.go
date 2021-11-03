package token

import (
	"go.aporeto.io/bahamut"
	"go.aporeto.io/elemental"
)

func FromRequest(req *elemental.Request) string {
	if hreq := req.HTTPRequest(); hreq != nil {
		if cookie, err := hreq.Cookie("x-a3s-token"); err == nil {
			return cookie.Value
		}
	}
	return req.Password
}

func FromSession(session bahamut.Session) string {
	if cookie, err := session.Cookie("x-a3s-token"); err == nil {
		return cookie.Value
	}
	return session.Token()
}
