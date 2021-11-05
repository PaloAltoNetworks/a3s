package token

import (
	"crypto"
	"crypto/sha1"
	"crypto/x509"
	"crypto/x509/pkix"
	"encoding/pem"
	"fmt"
	"testing"
	"time"

	"github.com/golang-jwt/jwt/v4"
	. "github.com/smartystreets/goconvey/convey"
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

func TestNewIdentityToken(t *testing.T) {

	Convey("Given I create a new Midgard claims", t, func() {

		c := NewIdentityToken(Source{
			Type:      "mtls",
			Namespace: "/my/ns",
			Name:      "mysource",
		})

		So(c.Identity, ShouldBeNil)
		So(c.Source.Type, ShouldEqual, "mtls")
		So(c.Source.Namespace, ShouldEqual, "/my/ns")
		So(c.Source.Name, ShouldEqual, "mysource")
	})
}

func TestParse(t *testing.T) {

	Convey("Given I create an IdentityToken", t, func() {

		cert, key := getECCert()

		token1 := NewIdentityToken(Source{
			Type:      "certificate",
			Namespace: "/my/ns",
			Name:      "mysource",
		})

		token1.Source.Type = "certificate"
		token1.Source.Namespace = "/my/ns"
		token1.Source.Name = "mysource"
		token1.Identity = []string{
			"org=a3s.com",
			"orgunit=admin",
			"commonname=joe",
		}

		keychain := NewJWKS()
		_ = keychain.Append(cert)

		kid := fmt.Sprintf("%02X", sha1.Sum(cert.Raw))

		token, err := token1.JWT(key, kid, "iss", "aud", time.Now().Add(10*time.Second))
		So(err, ShouldBeNil)

		Convey("Calling JWT with a missing source type should fail", func() {
			token1.Source.Type = ""
			_, err := token1.JWT(key, "kid", "iss", "aud", time.Now().Add(10*time.Second))
			So(err, ShouldNotBeNil)
			So(err.Error(), ShouldEqual, "invalid identity token: missing source type")
		})

		Convey("When I call Parse using the correct signer certificate", func() {

			token2 := NewIdentityToken(Source{
				Type:      "certificate",
				Namespace: "/my/ns",
				Name:      "mysource",
			})
			token2.Issuer = "iss"
			token2.Audience = jwt.ClaimStrings{"aud"}

			err = token2.Parse(token, keychain, "iss", "aud")

			So(err, ShouldBeNil)
			So(token2.Source.Type, ShouldEqual, "certificate")
			So(token2.Issuer, ShouldEqual, "iss")
			So(token2.Audience, ShouldResemble, jwt.ClaimStrings{"aud"})
			So(token2.ExpiresAt, ShouldResemble, token1.ExpiresAt)
			So(token2.IssuedAt, ShouldResemble, token1.IssuedAt)
			So(token2.Identity, ShouldResemble, []string{
				"org=a3s.com",
				"orgunit=admin",
				"commonname=joe",
				"@sourcetype=certificate",
				"@sourcenamespace=/my/ns",
				"@sourcename=mysource",
			})
		})

		Convey("When I call Parse using the wrong issuer", func() {

			token2 := NewIdentityToken(Source{
				Type:      "certificate",
				Namespace: "/my/ns",
				Name:      "mysource",
			})
			token2.Issuer = "not-iss"
			token2.Audience = jwt.ClaimStrings{"aud"}
			err = token2.Parse(token, keychain, "iss", "aud")

			So(err, ShouldNotBeNil)
			So(err.Error(), ShouldEqual, "issuer 'not-iss' is not acceptable. want 'iss'")
		})

		Convey("When I call Parse using the wrong audience", func() {

			token2 := NewIdentityToken(Source{
				Type:      "certificate",
				Namespace: "/my/ns",
				Name:      "mysource",
			})
			token2.Issuer = "iss"
			token2.Audience = jwt.ClaimStrings{"not-aud"}
			err = token2.Parse(token, keychain, "iss", "aud")

			So(err, ShouldNotBeNil)
			So(err.Error(), ShouldEqual, "audience '[not-aud]' is not acceptable. want 'aud'")
		})

		Convey("When I call Parse using the wrong signer certificate", func() {

			cert2, _ := getECCert()
			keychain2 := NewJWKS()
			_ = keychain2.Append(cert2)
			token2 := NewIdentityToken(Source{
				Type:      "certificate",
				Namespace: "/my/ns",
				Name:      "mysource",
			})
			err = token2.Parse(token, keychain2, "iss", "aud")

			So(err, ShouldNotBeNil)
			So(err.Error(), ShouldEqual, fmt.Sprintf("unable to parse jwt: unable to find kid '%s': kid not found in JWKS", kid))
		})

		Convey("When I call Parse using a wrong asigning method", func() {

			token = "eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9.eyJyZWFsbSI6IlZpbmNlIiwiaWR0IjpbIkBzb3VyY2U9YS1zb3VyY2UiLCJAbmFtZXNwYWNlPS9hL25hbWVzcGFjZSIsInVzZXJuYW1lPWpvZSIsInRlYW09cmVkIl0sImF1ZCI6Imh0dHBzOi8vYTNzLmNvbSIsImlhdCI6MTU0ODc5MDMwMCwiaXNzIjoiaHR0cHM6Ly9hM3MuY29tIn0.5PYuuULqrMArgdxq5eKSImsNskobw528Gr8Xe7HgPFs"
			token2 := NewIdentityToken(Source{
				Type:      "certificate",
				Namespace: "/my/ns",
				Name:      "mysource",
			})
			err = token2.Parse(token, keychain, "iss", "aud")

			So(err, ShouldNotBeNil)
			So(err.Error(), ShouldEqual, "unable to parse jwt: unexpected signing method: HS256")
		})
	})
}
