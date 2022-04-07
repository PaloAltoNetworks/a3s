package authenticator

import (
	"context"
	"crypto"
	"crypto/x509"
	"crypto/x509/pkix"
	"encoding/json"
	"encoding/pem"
	"net/http"
	"net/http/httptest"
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
	if claims.Issuer == "" {
		claims.Issuer = "iss"
	}
	claims.Audience = jwt.ClaimStrings{"aud"}
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
		a := New(
			jwks,
			"iss",
			"aud",
			OptionExternalTrustedIssuers(RemoteIssuer{URL: "a"}),
			OptionIgnoredResources("a", "b"),
		)

		Convey("Then a should be correct", func() {
			So(a.jwks, ShouldEqual, jwks)
			So(a.externalTrustedIssuers, ShouldResemble, map[string]RemoteIssuer{"a": {URL: "a"}})
			So(a.ignoredResources, ShouldResemble, map[string]struct{}{"a": {}, "b": {}})
		})
	})
}

func TestCommonAuth(t *testing.T) {

	Convey("Given I have an authenticator and a token", t, func() {

		c, k := getECCert()
		_, k2 := getECCert()
		jwks := token.NewJWKS()
		_ = jwks.Append(c)
		a := New(jwks, "iss", "aud")

		kid1 := token.Fingerprint(c)

		Convey("Calling commonAuth on a token signed by the signer should work", func() {

			token := makeToken(
				&token.IdentityToken{Identity: []string{"color=blue", "@source:type=test"}},
				jwt.SigningMethodES256,
				k,
				kid1,
			)

			action, claims, err := a.commonAuth(context.Background(), token)

			So(err, ShouldBeNil)
			So(action, ShouldEqual, bahamut.AuthActionContinue)
			So(claims, ShouldResemble, []string{"color=blue", "@source:type=test"})
		})

		Convey("Calling commonAuth on a refresh token should fail", func() {

			token := makeToken(
				&token.IdentityToken{Identity: []string{"color=blue", "@source:type=test"}, Refresh: true},
				jwt.SigningMethodES256,
				k,
				kid1,
			)

			action, claims, err := a.commonAuth(context.Background(), token)

			So(err, ShouldNotBeNil)
			So(err.Error(), ShouldEqual, "error 401 (a3s:authn): Unauthorized: Authentication impossible from a refresh token")
			So(action, ShouldEqual, bahamut.AuthActionKO)
			So(claims, ShouldBeNil)
		})

		Convey("Calling commonAuth on a token signed by the wrong signer should fail", func() {

			token := makeToken(
				&token.IdentityToken{Identity: []string{"color=blue", "@source:type=test"}},
				jwt.SigningMethodES256,
				k2,
				kid1,
			)

			action, claims, err := a.commonAuth(context.Background(), token)

			So(err, ShouldNotBeNil)
			So(err.Error(), ShouldEqual, `error 401 (a3s:authn): Unauthorized: Authentication rejected with error: unable to parse jwt: crypto/ecdsa: verification error`)
			So(action, ShouldEqual, bahamut.AuthActionKO)
			So(claims, ShouldBeNil)
		})

		Convey("Calling commonAuth on a bad token should fail", func() {

			token := "that's not good"

			action, claims, err := a.commonAuth(context.Background(), token)

			So(err, ShouldNotBeNil)
			So(err.Error(), ShouldEqual, `error 401 (a3s:authn): Unauthorized: Authentication rejected with error: unable to parse jwt: token contains an invalid number of segments`)
			So(action, ShouldEqual, bahamut.AuthActionKO)
			So(claims, ShouldBeNil)
		})

		Convey("Calling commonAuth on an empty token should fail", func() {

			token := ""

			action, claims, err := a.commonAuth(context.Background(), token)

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
		a := New(jwks, "iss", "aud")

		kid1 := token.Fingerprint(c)

		Convey("Calling AuthenticateSession on a session that has a valid token should work", func() {

			session := bahamut.NewMockSession()
			session.MockToken = makeToken(
				&token.IdentityToken{Identity: []string{"color=blue", "@source:type=test"}},
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
						&token.IdentityToken{Identity: []string{"color=blue", "@source:type=test"}},
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
				&token.IdentityToken{Identity: []string{"color=blue", "@source:type=test"}},
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

		Convey("Call AuthenticateRequest on a public resource source work", func() {

			a := New(jwks, "iss", "aud", OptionIgnoredResources("hello"))

			ctx := bahamut.NewMockContext(context.Background())
			req := elemental.NewRequest()
			req.Identity = elemental.MakeIdentity("hello", "hello")
			ctx.MockRequest = req
			action, err := a.AuthenticateRequest(ctx)

			So(err, ShouldBeNil)
			So(action, ShouldEqual, bahamut.AuthActionOK)
		})

		a := New(jwks, "iss", "aud")

		kid1 := token.Fingerprint(c)

		Convey("Call AuthenticateSession with a valid token should work", func() {

			ctx := bahamut.NewMockContext(context.Background())
			ctx.MockRequest = elemental.NewRequest()
			ctx.Request().Username = "Bearer"
			ctx.Request().Password = makeToken(
				&token.IdentityToken{Identity: []string{"color=blue", "@source:type=test"}},
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
					&token.IdentityToken{Identity: []string{"color=blue", "@source:type=test"}},
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
				&token.IdentityToken{Identity: []string{"color=blue", "@source:type=test"}},
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

func TestHandleFederatedToken(t *testing.T) {

	Convey("Given I have federated token", t, func() {

		c, k := getECCert()
		c2, k2 := getECCert()
		jwks := token.NewJWKS()
		_ = jwks.Append(c)

		kid1 := token.Fingerprint(c)

		_ = k
		Convey("I call the function with no externalTrustedIssuers", func() {
			a := New(jwks, "iss", "aud")
			rjwks, rissuer, err := a.handleFederatedToken(context.Background(), "not a token but it does not matter")
			So(err, ShouldBeNil)
			So(rjwks, ShouldBeNil)
			So(rissuer, ShouldBeEmpty)
		})

		Convey("When I call the function with an invalid token", func() {
			a := New(jwks, "iss", "aud", OptionExternalTrustedIssuers(RemoteIssuer{URL: "a"}))
			rjwks, rissuer, err := a.handleFederatedToken(context.Background(), "not a token and it will make it fail")
			So(err, ShouldNotBeNil)
			So(err.Error(), ShouldEqual, "Unable to parse input token: unable to parse unverified jwt: token contains an invalid number of segments")
			So(rjwks, ShouldBeNil)
			So(rissuer, ShouldBeEmpty)
		})

		Convey("When I call the function with a token using the local issuer", func() {
			a := New(jwks, "iss", "aud", OptionExternalTrustedIssuers(RemoteIssuer{URL: "a"}))
			token := makeToken(
				&token.IdentityToken{
					RegisteredClaims: jwt.RegisteredClaims{Issuer: "iss"},
					Identity:         []string{"@source:type=test"},
				},
				jwt.SigningMethodES256,
				k2,
				kid1,
			)
			rjwks, rissuer, err := a.handleFederatedToken(context.Background(), token)
			So(err, ShouldBeNil)
			So(rjwks, ShouldBeNil)
			So(rissuer, ShouldBeEmpty)
		})

		Convey("When I call the function with a token using * as issuer", func() {
			a := New(jwks, "iss", "aud", OptionExternalTrustedIssuers(RemoteIssuer{URL: "a"}))
			token := makeToken(
				&token.IdentityToken{
					RegisteredClaims: jwt.RegisteredClaims{Issuer: "*"},
					Identity:         []string{"@source:type=test"},
				},
				jwt.SigningMethodES256,
				k2,
				kid1,
			)
			rjwks, rissuer, err := a.handleFederatedToken(context.Background(), token)
			So(err, ShouldNotBeNil)
			So(err.Error(), ShouldEqual, "invalid iss field in token: what are you trying to do here?")
			So(rjwks, ShouldBeNil)
			So(rissuer, ShouldBeEmpty)
		})

		Convey("When I call the function with a token issuer not in trusted list", func() {
			a := New(jwks, "iss", "aud", OptionExternalTrustedIssuers(RemoteIssuer{URL: "a"}))
			token := makeToken(
				&token.IdentityToken{
					RegisteredClaims: jwt.RegisteredClaims{Issuer: "b"},
					Identity:         []string{"@source:type=test"},
				},
				jwt.SigningMethodES256,
				k2,
				kid1,
			)
			rjwks, rissuer, err := a.handleFederatedToken(context.Background(), token)
			So(err, ShouldBeNil)
			So(rjwks, ShouldBeNil)
			So(rissuer, ShouldBeEmpty)
		})

		Convey("When I call the function with an external token issuer in the trusted list", func() {

			var called int
			ts := httptest.NewTLSServer(http.HandlerFunc(func(w http.ResponseWriter, req *http.Request) {
				called++
				j := token.NewJWKS()
				j.Append(c2) // nolint
				d, _ := json.Marshal(j)
				w.Write(d) // nolint
			}))

			pool := x509.NewCertPool()
			pool.AddCert(ts.Certificate())

			a := New(jwks, "iss", "aud", OptionExternalTrustedIssuers(
				RemoteIssuer{
					URL:  ts.URL,
					Pool: pool,
				},
			))
			token := makeToken(
				&token.IdentityToken{
					RegisteredClaims: jwt.RegisteredClaims{Issuer: ts.URL},
					Identity:         []string{"@source:type=test"},
				},
				jwt.SigningMethodES256,
				k,
				kid1,
			)
			rjwks, rissuer, err := a.handleFederatedToken(context.Background(), token)
			So(err, ShouldBeNil)
			So(rjwks, ShouldNotBeNil)
			So(rissuer, ShouldEqual, ts.URL)
			So(called, ShouldEqual, 1)

			Convey("I do it again, it should be cached", func() {
				rjwks, rissuer, err := a.handleFederatedToken(context.Background(), token)
				So(err, ShouldBeNil)
				So(rjwks, ShouldNotBeNil)
				So(rissuer, ShouldEqual, ts.URL)
				So(called, ShouldEqual, 1)
			})
		})

		Convey("When I call the function with an external token issuer in the trusted list but with no pool and untrusted cert", func() {

			ts := httptest.NewTLSServer(http.HandlerFunc(func(w http.ResponseWriter, req *http.Request) {
			}))

			a := New(jwks, "iss", "aud", OptionExternalTrustedIssuers(
				RemoteIssuer{
					URL: ts.URL,
				},
			))
			token := makeToken(
				&token.IdentityToken{
					RegisteredClaims: jwt.RegisteredClaims{Issuer: ts.URL},
					Identity:         []string{"@source:type=test"},
				},
				jwt.SigningMethodES256,
				k,
				kid1,
			)
			rjwks, rissuer, err := a.handleFederatedToken(context.Background(), token)
			So(err, ShouldNotBeNil)
			So(err.Error(), ShouldStartWith, "Unable to retrieve remote jwks: unable to retrieve remote jwks:")
			So(err.Error(), ShouldEndWith, "x509: certificate signed by unknown authority")
			So(rjwks, ShouldBeNil)
			So(rissuer, ShouldBeEmpty)
		})
	})
}
