package token

import (
	"fmt"

	"github.com/golang-jwt/jwt/v4"
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

func makeKeyFunc(keychain *JWKS) jwt.Keyfunc {

	return func(token *jwt.Token) (interface{}, error) {

		if _, ok := token.Method.(*jwt.SigningMethodECDSA); !ok {
			return nil, fmt.Errorf("unexpected signing method: %s", token.Header["alg"])
		}

		kid, ok := token.Header["kid"].(string)
		if !ok || kid == "" {
			return nil, fmt.Errorf("token has no KID in its header")
		}

		k, err := keychain.Get(kid)
		if err != nil {
			return nil, fmt.Errorf("unable to find kid '%s': %w", kid, err)
		}

		return k.PublicKey(), nil
	}
}
