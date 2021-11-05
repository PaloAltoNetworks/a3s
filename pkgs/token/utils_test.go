package token

import (
	"crypto/ecdsa"
	"crypto/sha1"
	"fmt"
	"net/http"
	"testing"

	"github.com/golang-jwt/jwt/v4"
	. "github.com/smartystreets/goconvey/convey"
	"go.aporeto.io/a3s/pkgs/api"
	"go.aporeto.io/bahamut"
	"go.aporeto.io/elemental"
)

func TestFromRequest(t *testing.T) {

	Convey("Calling FromRequest with a bearer token", t, func() {
		httpreq, _ := http.NewRequest(http.MethodGet, "https://toto.com/namespaces", nil)
		httpreq.Header.Set("Authorization", "Bearer token")
		req, _ := elemental.NewRequestFromHTTPRequest(httpreq, api.Manager())

		t := FromRequest(req)
		So(t, ShouldEqual, "token")
	})

	Convey("Calling FromRequest with a cookie token", t, func() {
		httpreq, _ := http.NewRequest(http.MethodGet, "https://toto.com/namespaces", nil)
		httpreq.AddCookie(&http.Cookie{Name: "x-a3s-token", Value: "token"})
		req, _ := elemental.NewRequestFromHTTPRequest(httpreq, api.Manager())
		t := FromRequest(req)
		So(t, ShouldEqual, "token")
	})

	Convey("Calling FromRequest with both bearer and  cookie token", t, func() {
		httpreq, _ := http.NewRequest(http.MethodGet, "https://toto.com/namespaces", nil)
		httpreq.Header.Set("Authorization", "Bearer token1")
		httpreq.AddCookie(&http.Cookie{Name: "x-a3s-token", Value: "token2"})
		req, _ := elemental.NewRequestFromHTTPRequest(httpreq, api.Manager())
		t := FromRequest(req)
		So(t, ShouldEqual, "token2")
	})
}

func TestFromSession(t *testing.T) {

	Convey("Calling FromRequest with a bearer token", t, func() {
		s := bahamut.NewMockSession()
		s.MockToken = "token"
		t := FromSession(s)
		So(t, ShouldEqual, "token")
	})

	Convey("Calling FromRequest with a cookie token", t, func() {
		s := bahamut.NewMockSession()
		s.MockCookies = map[string]*http.Cookie{"x-a3s-token": {Name: "x-a3s-token", Value: "token"}}
		t := FromSession(s)
		So(t, ShouldEqual, "token")
	})

	Convey("Calling FromRequest with both bearer and  cookie token", t, func() {
		s := bahamut.NewMockSession()
		s.MockCookies = map[string]*http.Cookie{"x-a3s-token": {Name: "x-a3s-token", Value: "token2"}}
		s.MockToken = "token"
		t := FromSession(s)
		So(t, ShouldEqual, "token2")
	})
}

func TestMakeKeyFunc(t *testing.T) {

	Convey("Given I make a key func", t, func() {

		cert1, _ := getECCert()

		keychain := NewJWKS()
		_ = keychain.Append(cert1)

		f := makeKeyFunc(keychain)

		Convey("Calling on token with invalid alg should fail", func() {
			t := &jwt.Token{
				Method: jwt.SigningMethodHS512,
				Header: map[string]interface{}{"alg": "nope"},
			}
			_, err := f(t)
			So(err, ShouldNotBeNil)
			So(err.Error(), ShouldEqual, "unexpected signing method: nope")
		})

		Convey("Calling on token with missing kid header should fail", func() {
			t := &jwt.Token{
				Method: jwt.SigningMethodES256,
				Header: map[string]interface{}{},
			}
			_, err := f(t)
			So(err, ShouldNotBeNil)
			So(err.Error(), ShouldEqual, "token has no KID in its header")
		})

		Convey("Calling on token with empty kid header should fail", func() {
			t := &jwt.Token{
				Method: jwt.SigningMethodES256,
				Header: map[string]interface{}{"kid": ""},
			}
			_, err := f(t)
			So(err, ShouldNotBeNil)
			So(err.Error(), ShouldEqual, "token has no KID in its header")
		})

		Convey("Calling on token with kid not found", func() {
			t := &jwt.Token{
				Method: jwt.SigningMethodES256,
				Header: map[string]interface{}{"kid": "kidding"},
			}
			_, err := f(t)
			So(err, ShouldNotBeNil)
			So(err.Error(), ShouldEqual, "unable to find kid 'kidding': kid not found in JWKS")
		})

		Convey("Calling on token with valid kid", func() {
			t := &jwt.Token{
				Method: jwt.SigningMethodES256,
				Header: map[string]interface{}{"kid": fmt.Sprintf("%02X", sha1.Sum(cert1.Raw))},
			}
			k, err := f(t)
			So(err, ShouldBeNil)
			So(k, ShouldNotBeNil)
			So(k.(*ecdsa.PublicKey).Equal(cert1.PublicKey), ShouldBeTrue)

		})
	})
}
