package authenticator

import (
	"context"
	"crypto"
	"crypto/sha1"
	"crypto/x509"
	"crypto/x509/pkix"
	"encoding/pem"
	"fmt"
	"net/http"
	"testing"

	"github.com/golang-jwt/jwt/v4"
	. "github.com/smartystreets/goconvey/convey"
	"go.aporeto.io/a3s/pkgs/api"
	"go.aporeto.io/a3s/pkgs/token"
	"go.aporeto.io/bahamut"
	"go.aporeto.io/elemental"
	"go.aporeto.io/tg/tglib"
)

func getECCert() (*x509.Certificate, crypto.PrivateKey) {

	certBlock, keyBlock, err := tglib.Issue(pkix.Name{})
	if err != nil {
		panic(err)
	}

	cert, err := tglib.ParseCertificate(pem.EncodeToMemory(certBlock))
	if err != nil {
		panic(err)
	}

	key, err := tglib.PEMToKey(keyBlock)
	if err != nil {
		panic(err)
	}

	return cert, key
}

func makeToken(claims *token.IdentityToken, signMethod jwt.SigningMethod, key crypto.PrivateKey, kid string) string {
	token := jwt.NewWithClaims(signMethod, claims)
	token.Header["kid"] = kid
	t, err := token.SignedString(key)
	if err != nil {
		panic(err)
	}
	return t
}

func TestNewAuthenticator(t *testing.T) {

	Convey("Given I create a new Authenticator", t, func() {

		c, _ := getECCert()
		jwks := token.NewJWKS()
		_ = jwks.Append(c)
		a := NewPrivate(jwks)

		Convey("Then a should be correct", func() {
			So(a.jwks, ShouldEqual, jwks)
		})
	})
}

func TestCommonAuth(t *testing.T) {

	Convey("Given I have an authenticator and a token", t, func() {

		c, k := getECCert()
		_, k2 := getECCert()
		jwks := token.NewJWKS()
		_ = jwks.Append(c)
		a := NewPrivate(jwks)

		kid1 := fmt.Sprintf("%02X", sha1.Sum(c.Raw))

		Convey("Calling commonAuth on a token signed by the signer should work", func() {

			token := makeToken(
				&token.IdentityToken{Identity: []string{"color=blue"}},
				jwt.SigningMethodES256,
				k,
				kid1,
			)

			action, claims, err := a.commonAuth(token)

			So(err, ShouldBeNil)
			So(action, ShouldEqual, bahamut.AuthActionContinue)
			So(claims, ShouldResemble, []string{"color=blue"})
		})

		Convey("Calling commonAuth on a token signed by the wrong signer should fail", func() {

			token := makeToken(
				&token.IdentityToken{Identity: []string{"color=blue"}},
				jwt.SigningMethodES256,
				k2,
				kid1,
			)

			action, claims, err := a.commonAuth(token)

			So(err, ShouldNotBeNil)
			So(err.Error(), ShouldEqual, `error 401 (a3s:authn): Unauthorized: Authentication rejected with error: unable to parse jwt: crypto/ecdsa: verification error`)
			So(action, ShouldEqual, bahamut.AuthActionKO)
			So(claims, ShouldBeNil)
		})

		Convey("Calling commonAuth on a bad token should fail", func() {

			token := "that's not good"

			action, claims, err := a.commonAuth(token)

			So(err, ShouldNotBeNil)
			So(err.Error(), ShouldEqual, `error 401 (a3s:authn): Unauthorized: Authentication rejected with error: unable to parse jwt: token contains an invalid number of segments`)
			So(action, ShouldEqual, bahamut.AuthActionKO)
			So(claims, ShouldBeNil)
		})

		Convey("Calling commonAuth on an empty token should fail", func() {

			token := ""

			action, claims, err := a.commonAuth(token)

			So(err, ShouldNotBeNil)
			So(err.Error(), ShouldEqual, `error 401 (a3s:authn): Unauthorized: Missing token in Authorization header`)
			So(action, ShouldEqual, bahamut.AuthActionKO)
			So(claims, ShouldBeNil)
		})
	})
}

func TestAuthenticateSession(t *testing.T) {

	Convey("Given I have an authenticator", t, func() {

		c, k := getECCert()
		_, k2 := getECCert()
		jwks := token.NewJWKS()
		_ = jwks.Append(c)
		a := NewPrivate(jwks)

		kid1 := fmt.Sprintf("%02X", sha1.Sum(c.Raw))

		Convey("Calling AuthenticateSession on a session that has a valid token should work", func() {

			session := bahamut.NewMockSession()
			session.MockToken = makeToken(
				&token.IdentityToken{Identity: []string{"color=blue"}},
				jwt.SigningMethodES256,
				k,
				kid1,
			)

			action, err := a.AuthenticateSession(session)

			So(err, ShouldBeNil)
			So(action, ShouldEqual, bahamut.AuthActionContinue)
		})

		Convey("Call AuthenticateSession with a valid token in cookies should work", func() {

			session := bahamut.NewMockSession()
			session.MockCookies = map[string]*http.Cookie{
				"x-a3s-token": {
					Name: "x-a3s-token",
					Value: makeToken(
						&token.IdentityToken{Identity: []string{"color=blue"}},
						jwt.SigningMethodES256,
						k,
						kid1,
					),
				},
			}

			action, err := a.AuthenticateSession(session)

			So(err, ShouldBeNil)
			So(action, ShouldEqual, bahamut.AuthActionContinue)
		})

		Convey("Calling AuthenticateSession on a session that has a invalid token should fail", func() {

			session := bahamut.NewMockSession()
			session.MockToken = makeToken(
				&token.IdentityToken{Identity: []string{"color=blue"}},
				jwt.SigningMethodES256,
				k2,
				kid1,
			)

			action, err := a.AuthenticateSession(session)

			So(err, ShouldNotBeNil)
			So(err.Error(), ShouldEqual, "error 401 (a3s:authn): Unauthorized: Authentication rejected with error: unable to parse jwt: crypto/ecdsa: verification error")
			So(action, ShouldEqual, bahamut.AuthActionKO)
		})
	})
}

func TestAuthenticateRequest(t *testing.T) {

	Convey("Given I have an authenticator", t, func() {

		c, k := getECCert()
		_, k2 := getECCert()
		jwks := token.NewJWKS()
		_ = jwks.Append(c)
		a := NewPrivate(jwks)

		kid1 := fmt.Sprintf("%02X", sha1.Sum(c.Raw))

		Convey("Call AuthenticateSession with a valid token should work", func() {

			ctx := bahamut.NewMockContext(context.Background())
			ctx.MockRequest = elemental.NewRequest()
			ctx.Request().Username = "Bearer"
			ctx.Request().Password = makeToken(
				&token.IdentityToken{Identity: []string{"color=blue"}},
				jwt.SigningMethodES256,
				k,
				kid1,
			)

			action, err := a.AuthenticateRequest(ctx)

			So(err, ShouldBeNil)
			So(action, ShouldEqual, bahamut.AuthActionContinue)
		})

		Convey("Call AuthenticateSession with a valid token in cookies should work", func() {

			req, _ := http.NewRequest(http.MethodGet, "http://url.com/hello", nil)
			req.AddCookie(&http.Cookie{
				Name: "x-a3s-token",
				Value: makeToken(
					&token.IdentityToken{Identity: []string{"color=blue"}},
					jwt.SigningMethodES256,
					k,
					kid1,
				),
			})
			ctx := bahamut.NewMockContext(context.Background())
			ctx.MockRequest, _ = elemental.NewRequestFromHTTPRequest(req, api.Manager())
			action, err := a.AuthenticateRequest(ctx)

			So(err, ShouldBeNil)
			So(action, ShouldEqual, bahamut.AuthActionContinue)
		})

		Convey("Calling AuthenticateSession with an invalid token should work", func() {

			ctx := bahamut.NewMockContext(context.Background())
			ctx.MockRequest = elemental.NewRequest()
			ctx.Request().Username = "Bearer"
			ctx.Request().Password = makeToken(
				&token.IdentityToken{Identity: []string{"color=blue"}},
				jwt.SigningMethodES256,
				k2,
				kid1,
			)

			action, err := a.AuthenticateRequest(ctx)

			So(err, ShouldNotBeNil)
			So(err.Error(), ShouldEqual, "error 401 (a3s:authn): Unauthorized: Authentication rejected with error: unable to parse jwt: crypto/ecdsa: verification error")
			So(action, ShouldEqual, bahamut.AuthActionKO)
		})
	})
}
