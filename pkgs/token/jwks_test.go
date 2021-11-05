package token

import (
	"crypto/ecdsa"
	"crypto/elliptic"
	"crypto/sha1"
	"encoding/base64"
	"fmt"
	"math/big"
	"testing"

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
		kid1 := fmt.Sprintf("%02X", sha1.Sum(cert1.Raw))
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
