package token

import (
	"context"
	"crypto/ecdsa"
	"crypto/sha256"
	"crypto/tls"
	"crypto/x509"
	"encoding/json"
	"fmt"
	"math/big"
	"net/http"
	"net/http/httptest"
	"reflect"
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
				Header: map[string]interface{}{"kid": Fingerprint(cert1)},
			}
			k, err := f(t)
			So(err, ShouldBeNil)
			So(k, ShouldNotBeNil)
			So(k.(*ecdsa.PublicKey).Equal(cert1.PublicKey), ShouldBeTrue)

		})
	})
}

func TestFingerprint(t *testing.T) {
	cert, _ := getECCert()
	type args struct {
		cert *x509.Certificate
	}
	tests := []struct {
		name string
		args func(t *testing.T) args

		want1 string
	}{
		{
			"standard",
			func(*testing.T) args {
				return args{
					cert: cert,
				}
			},
			fmt.Sprintf("%02X", sha256.Sum256(cert.Raw)),
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			tArgs := tt.args(t)

			got1 := Fingerprint(tArgs.cert)

			if !reflect.DeepEqual(got1, tt.want1) {
				t.Errorf("Fingerprint got1 = %v, want1: %v", got1, tt.want1)
			}
		})
	}
}

func TestJWKSFromTokenIssuer(t *testing.T) {

	Convey("given a working remote jwks", t, func() {

		ts := httptest.NewTLSServer(http.HandlerFunc(func(w http.ResponseWriter, req *http.Request) {
			j := NewJWKS()
			j.Keys = []*JWKSKey{
				{
					KID: "kid",
					X:   "vs1LjF38O2OOSmy0Nbo45zfroQ1ME7GLuZ69uuiCOkk",
					Y:   "vs1LjF38O2OOSmy0Nbo45zfroQ1ME7GLuZ69uuiCOkk",
				},
			}

			d, _ := json.Marshal(j)
			w.Write(d) // nolint
		}))

		idt := &IdentityToken{
			RegisteredClaims: jwt.RegisteredClaims{
				Issuer: ts.URL,
			},
		}

		Convey("Calling JWKSFromTokenIssuer with certificates issue should fail", func() {

			jwks, err := JWKSFromTokenIssuer(
				context.Background(),
				idt,
				nil,
			)

			So(err, ShouldNotBeNil)
			So(err.Error(), ShouldStartWith, "unable to retrieve remote jwks: remote jwks error: unable to send request:")
			So(jwks, ShouldBeNil)
		})

		Convey("Calling JWKSFromTokenIssuer without problem should work", func() {

			jwks, err := JWKSFromTokenIssuer(
				context.Background(),
				idt,
				&tls.Config{
					InsecureSkipVerify: true, // nolint
				},
			)

			So(err, ShouldBeNil)
			So(len(jwks.Keys), ShouldEqual, 1)
			So(len(jwks.keyMap), ShouldEqual, 1)
			So(jwks.keyMap, ShouldContainKey, "kid")
			So(jwks.Keys[0].x, ShouldHaveSameTypeAs, big.NewInt(42))
			So(jwks.Keys[0].y, ShouldHaveSameTypeAs, big.NewInt(42))
		})
	})
}
