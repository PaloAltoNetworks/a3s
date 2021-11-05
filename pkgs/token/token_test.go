package token

import (
	"crypto"
	"crypto/x509"
	"crypto/x509/pkix"
	"encoding/pem"
	"testing"
	"time"

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
		token1.Issuer = "https://a3s.com"
		token1.Audience = "https://a3s.com"
		token1.ExpiresAt = time.Now().Add(10 * time.Second).Unix()
		token1.IssuedAt = time.Now().Unix()
		token1.Identity = []string{
			"org=a3s.com",
			"orgunit=admin",
			"commonname=joe",
		}

		token, err := token1.JWT(key, "kid", time.Now().Add(10*time.Second))
		So(err, ShouldBeNil)

		Convey("Calling JWT with a missing source type should fail", func() {
			token1.Source.Type = ""
			_, err := token1.JWT(key, "kid", time.Now().Add(10*time.Second))
			So(err, ShouldNotBeNil)
			So(err.Error(), ShouldEqual, "invalid identity token: missing source type")
		})

		Convey("When I call Parse using the correct signer certificate", func() {

			token2 := NewIdentityToken(Source{
				Type:      "certificate",
				Namespace: "/my/ns",
				Name:      "mysource",
			})
			err = token2.Parse(token, cert, "https://a3s.com", "https://a3s.com")

			So(err, ShouldBeNil)
			So(token2.Source.Type, ShouldEqual, "certificate")
			So(token2.Issuer, ShouldEqual, "https://a3s.com")
			So(token2.Audience, ShouldEqual, "https://a3s.com")
			So(token2.ExpiresAt, ShouldEqual, token1.ExpiresAt)
			So(token2.IssuedAt, ShouldEqual, token1.IssuedAt)
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
			err = token2.Parse(token, cert, "https://not-a3s.com", "https://a3s.com")

			So(err, ShouldNotBeNil)
			So(err.Error(), ShouldEqual, "issuer 'https://a3s.com' is not acceptable. want 'https://not-a3s.com'")
		})

		Convey("When I call Parse using the wrong audience", func() {

			token2 := NewIdentityToken(Source{
				Type:      "certificate",
				Namespace: "/my/ns",
				Name:      "mysource",
			})
			err = token2.Parse(token, cert, "https://a3s.com", "https://not-a3s.com")

			So(err, ShouldNotBeNil)
			So(err.Error(), ShouldEqual, "audience 'https://a3s.com' is not acceptable. want 'https://not-a3s.com'")
		})

		Convey("When I call Parse using the wrong signer certificate", func() {

			cert2, _ := getECCert()
			token2 := NewIdentityToken(Source{
				Type:      "certificate",
				Namespace: "/my/ns",
				Name:      "mysource",
			})
			err = token2.Parse(token, cert2, "https://a3s.com", "https://a3s.com")

			So(err, ShouldNotBeNil)
			So(err.Error(), ShouldEqual, "unable to parse jwt: crypto/ecdsa: verification error")
		})

		Convey("When I call Parse using a wrong asigning method", func() {

			token = "eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9.eyJyZWFsbSI6IlZpbmNlIiwiaWR0IjpbIkBzb3VyY2U9YS1zb3VyY2UiLCJAbmFtZXNwYWNlPS9hL25hbWVzcGFjZSIsInVzZXJuYW1lPWpvZSIsInRlYW09cmVkIl0sImF1ZCI6Imh0dHBzOi8vYTNzLmNvbSIsImlhdCI6MTU0ODc5MDMwMCwiaXNzIjoiaHR0cHM6Ly9hM3MuY29tIn0.5PYuuULqrMArgdxq5eKSImsNskobw528Gr8Xe7HgPFs"
			token2 := NewIdentityToken(Source{
				Type:      "certificate",
				Namespace: "/my/ns",
				Name:      "mysource",
			})
			err = token2.Parse(token, cert, "https://a3s.com", "https://a3s.com")

			So(err, ShouldNotBeNil)
			So(err.Error(), ShouldEqual, "unable to parse jwt: unexpected signing method: HS256")
		})
	})
}
