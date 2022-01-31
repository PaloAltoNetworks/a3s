package token

import (
	"context"
	"crypto/ecdsa"
	"crypto/elliptic"
	"encoding/base64"
	"encoding/json"
	"errors"
	"math/big"
	"net/http"
	"net/http/httptest"
	"testing"
	"time"

	. "github.com/smartystreets/goconvey/convey"
)

func TestNewJWKS(t *testing.T) {
	Convey("Calling TestNewJWKS should work", t, func() {

		k := NewJWKS()
		So(k.keyMap, ShouldNotBeNil)
	})
}

func TestJWKSCrud(t *testing.T) {

	Convey("Given I have a new JWKS", t, func() {

		cert1, key1 := getECCert()
		kid1 := Fingerprint(cert1)
		k := NewJWKS()
		So(k.GetLast(), ShouldBeNil)

		Convey("Appending a cert should work", func() {

			err := k.AppendWithPrivate(cert1, key1)
			So(err, ShouldBeNil)
			So(k.keyMap, ShouldContainKey, kid1)
			So(len(k.keyMap), ShouldEqual, 1)
			So(len(k.Keys), ShouldEqual, 1)
			So(k.GetLast(), ShouldNotBeNil)

			Convey("Appending it again should fail", func() {
				err := k.Append(cert1)
				So(err, ShouldNotBeNil)
				So(err, ShouldEqual, ErrJWKSKeyExists)
				So(k.keyMap, ShouldContainKey, kid1)
				So(len(k.keyMap), ShouldEqual, 1)
				So(len(k.Keys), ShouldEqual, 1)
			})

			Convey("Getting the key should work", func() {
				kk, err := k.Get(kid1)
				So(err, ShouldBeNil)
				So(kk, ShouldNotBeNil)
				So(kk.KID, ShouldEqual, kid1)
				So(kk.PrivateKey(), ShouldEqual, key1)
				So(k.keyMap, ShouldContainKey, kid1)
				So(len(k.Keys), ShouldEqual, 1)

				Convey("The generated JWSKey should be correct", func() {
					ecdsakey := cert1.PublicKey.(*ecdsa.PublicKey)
					So(kk.KID, ShouldEqual, kid1)
					So(kk.Use, ShouldEqual, "sign")
					So(kk.Curve(), ShouldResemble, elliptic.P256())
					So(kk.KTY, ShouldEqual, "EC")
					So(kk.x.String(), ShouldEqual, ecdsakey.X.String())
					So(kk.y.String(), ShouldEqual, ecdsakey.Y.String())
					So(kk.X, ShouldEqual, base64.RawURLEncoding.EncodeToString(ecdsakey.X.Bytes()))
					So(kk.Y, ShouldEqual, base64.RawURLEncoding.EncodeToString(ecdsakey.Y.Bytes()))
				})
			})

			Convey("Getting another kid should fail", func() {
				kk, err := k.Get("kidding")
				So(err, ShouldNotBeNil)
				So(err, ShouldEqual, ErrJWKSNotFound)
				So(kk, ShouldBeNil)
				So(k.keyMap, ShouldContainKey, kid1)
				So(len(k.keyMap), ShouldEqual, 1)
				So(len(k.Keys), ShouldEqual, 1)
			})

			Convey("Deleting the key should work", func() {
				deleted := k.Del(kid1)
				So(deleted, ShouldBeTrue)
				So(k.keyMap, ShouldNotContainKey, kid1)
				So(len(k.keyMap), ShouldEqual, 0)
				So(len(k.Keys), ShouldEqual, 0)
				So(k.GetLast(), ShouldBeNil)

				Convey("Deleting the key again should fail", func() {
					deleted := k.Del(kid1)
					So(deleted, ShouldBeFalse)
					So(k.keyMap, ShouldNotContainKey, kid1)
					So(len(k.keyMap), ShouldEqual, 0)
					So(len(k.Keys), ShouldEqual, 0)
					So(k.GetLast(), ShouldBeNil)
				})
			})

			Convey("Deleting another the key should fail", func() {
				deleted := k.Del("kidding")
				So(deleted, ShouldBeFalse)
				So(k.keyMap, ShouldContainKey, kid1)
				So(len(k.keyMap), ShouldEqual, 1)
				So(len(k.Keys), ShouldEqual, 1)
			})
		})
	})
}

func TestNewRemoteJWKS(t *testing.T) {

	Convey("Given I call the function with a missing context", t, func() {

		jwks, err := NewRemoteJWKS(nil, nil, "toto://not-an-url") // nolint
		So(jwks, ShouldBeNil)
		So(err, ShouldNotBeNil)
		So(errors.As(err, &ErrJWKSRemote{}), ShouldBeTrue)
		So(err.Error(), ShouldEqual, `remote jwks error: unable to build request: net/http: nil Context`)
		So(err.(ErrJWKSRemote).Unwrap().Error(), ShouldEqual, `unable to build request: net/http: nil Context`)
	})

	Convey("Given I call the function pointing to a non existing server", t, func() {

		ctx, cancel := context.WithTimeout(context.Background(), 300*time.Millisecond)
		defer cancel()

		jwks, err := NewRemoteJWKS(ctx, nil, "https://122.33.33.33")
		So(jwks, ShouldBeNil)
		So(err, ShouldNotBeNil)
		So(errors.As(err, &ErrJWKSRemote{}), ShouldBeTrue)
		So(err.Error(), ShouldEqual, `remote jwks error: unable to send request: Get "https://122.33.33.33": context deadline exceeded`)
	})

	Convey("Given a http server that returns invalid body", t, func() {

		ts := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, req *http.Request) {
			w.WriteHeader(http.StatusTeapot)
		}))
		jwks, err := NewRemoteJWKS(context.Background(), nil, ts.URL)
		So(jwks, ShouldBeNil)
		So(err, ShouldNotBeNil)
		So(errors.As(err, &ErrJWKSRemote{}), ShouldBeTrue)
		So(err.Error(), ShouldEqual, `remote jwks error: unable to parse response body: unable to decode application/json: EOF`)
	})

	Convey("Given a http server that returns invalid X", t, func() {

		ts := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, req *http.Request) {
			j := NewJWKS()
			j.Keys = []*JWKSKey{
				{
					X: "oh no..",
					Y: "vs1LjF38O2OOSmy0Nbo45zfroQ1ME7GLuZ69uuiCOkk",
				},
			}

			d, _ := json.Marshal(j)
			w.Write(d) // nolint
		}))

		jwks, err := NewRemoteJWKS(context.Background(), nil, ts.URL)
		So(jwks, ShouldBeNil)
		So(err, ShouldNotBeNil)
		So(errors.As(err, &ErrJWKSRemote{}), ShouldBeTrue)
		So(err.Error(), ShouldEqual, `remote jwks error: unable to decode X: illegal base64 data at input byte 2`)
	})

	Convey("Given a http server that returns invalid X", t, func() {

		ts := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, req *http.Request) {
			j := NewJWKS()
			j.Keys = []*JWKSKey{
				{
					X: "vs1LjF38O2OOSmy0Nbo45zfroQ1ME7GLuZ69uuiCOkk",
					Y: "oh no..",
				},
			}

			d, _ := json.Marshal(j)
			w.Write(d) // nolint
		}))

		jwks, err := NewRemoteJWKS(context.Background(), nil, ts.URL)
		So(jwks, ShouldBeNil)
		So(err, ShouldNotBeNil)
		So(errors.As(err, &ErrJWKSRemote{}), ShouldBeTrue)
		So(err.Error(), ShouldEqual, `remote jwks error: unable to decode Y: illegal base64 data at input byte 2`)
	})

	Convey("Given a http server that returns a valid JWKS", t, func() {

		ts := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, req *http.Request) {
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

		jwks, err := NewRemoteJWKS(context.Background(), nil, ts.URL)
		So(err, ShouldBeNil)
		So(len(jwks.Keys), ShouldEqual, 1)
		So(len(jwks.keyMap), ShouldEqual, 1)
		So(jwks.keyMap, ShouldContainKey, "kid")
		So(jwks.Keys[0].x, ShouldHaveSameTypeAs, big.NewInt(42))
		So(jwks.Keys[0].y, ShouldHaveSameTypeAs, big.NewInt(42))
	})
}

func TestJWKSKeyCurve(t *testing.T) {

	Convey("Calling curve when CRV is P-224 should work", t, func() {
		k := &JWKSKey{CRV: "P-224"}
		So(k.Curve(), ShouldResemble, elliptic.P224())
	})
	Convey("Calling curve when CRV is P-256 should work", t, func() {
		k := &JWKSKey{CRV: "P-256"}
		So(k.Curve(), ShouldResemble, elliptic.P256())
	})
	Convey("Calling curve when CRV is P-384 should work", t, func() {
		k := &JWKSKey{CRV: "P-384"}
		So(k.Curve(), ShouldResemble, elliptic.P384())
	})
	Convey("Calling curve when CRV is P-521 should work", t, func() {
		k := &JWKSKey{CRV: "P-521"}
		So(k.Curve(), ShouldResemble, elliptic.P521())
	})
	Convey("Calling curve when CRV is empty should return nothing", t, func() {
		k := &JWKSKey{CRV: ""}
		So(k.Curve(), ShouldResemble, nil)
	})
}

func TestJWKSKeyPublicKey(t *testing.T) {

	Convey("Calling PublicKey with CRV set to EC should work", t, func() {
		k := &JWKSKey{
			KTY: "EC",
			CRV: "P-224",
			x:   big.NewInt(42),
			y:   big.NewInt(42),
		}
		So(k.PublicKey().(*ecdsa.PublicKey).X.String(), ShouldEqual, k.x.String())
		So(k.PublicKey().(*ecdsa.PublicKey).Y.String(), ShouldResemble, k.y.String())
	})

	Convey("Calling PublicKey with CRV set to something else should fail", t, func() {
		k := &JWKSKey{
			KTY: "not-EC",
		}
		So(k.PublicKey(), ShouldBeNil)
	})
}
