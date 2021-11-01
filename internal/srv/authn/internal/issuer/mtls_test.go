package claims

import (
	"crypto"
	"crypto/x509"
	"crypto/x509/pkix"
	"encoding/pem"
	"math/big"
	"testing"

	. "github.com/smartystreets/goconvey/convey"
	"go.aporeto.io/tg/tglib"
)

func getCA() (*x509.Certificate, crypto.PrivateKey) {

	certBlock, keyBlock, err := tglib.Issue(pkix.Name{}, tglib.OptIssueTypeCA())
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

func getECCert(subject pkix.Name, opts ...tglib.IssueOption) (*x509.Certificate, crypto.PrivateKey) {

	certBlock, keyBlock, err := tglib.Issue(subject, opts...)
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

func TestMTLSIssuer(t *testing.T) {

	Convey("Given I have some certificates", t, func() {

		cacert1, cakey1 := getCA()
		usercert1, _ := getECCert(
			pkix.Name{
				CommonName:         "jean-mich",
				Country:            []string{"US", "France"},
				Province:           []string{"CA", "IDF"},
				Locality:           []string{"Sunnyvale", "Paris"},
				PostalCode:         []string{"94085", "75009"},
				Organization:       []string{"o1", "o2"},
				OrganizationalUnit: []string{"ou1", "ou2"},
				StreetAddress:      []string{"3000 Tanery way", "13 Rue de Mogador"},
				SerialNumber:       "1234567890",
			},
			tglib.OptIssueSigner(cacert1, cakey1),
			tglib.OptIssueDNSSANs("toto.com", "titi.com"),
			tglib.OptIssueSerialNumber(big.NewInt(42)),
			tglib.OptIssueEmailAddresses([]string{"me@me.com"}),
		)
		usercert2, _ := getECCert(pkix.Name{})

		Convey("Calling NewMTLSIssuer should work", func() {

			pool := x509.NewCertPool()
			pool.AddCert(cacert1)

			iss := NewMTLSIssuer(pool, "/my/ns", "mysource")
			So(iss.caPool, ShouldEqual, pool)
			So(iss.token, ShouldNotBeNil)

			Convey("Calling FromCertificate with a valid user cert should work", func() {

				err := iss.FromCertificate(usercert1)
				So(err, ShouldBeNil)

				Convey("Calling Issue should work", func() {

					idt := iss.Issue()
					So(len(idt.Identity), ShouldEqual, 19)
					So(idt.Identity, ShouldContain, "commonname=jean-mich")
					So(idt.Identity, ShouldContain, "country=US")
					So(idt.Identity, ShouldContain, "country=France")
					So(idt.Identity, ShouldContain, "province=CA")
					So(idt.Identity, ShouldContain, "province=IDF")
					So(idt.Identity, ShouldContain, "locality=Sunnyvale")
					So(idt.Identity, ShouldContain, "locality=Paris")
					So(idt.Identity, ShouldContain, "postalcode=94085")
					So(idt.Identity, ShouldContain, "postalcode=75009")
					So(idt.Identity, ShouldContain, "organization=o1")
					So(idt.Identity, ShouldContain, "organization=o2")
					So(idt.Identity, ShouldContain, "organizationalunit=ou1")
					So(idt.Identity, ShouldContain, "organizationalunit=ou2")
					So(idt.Identity, ShouldContain, "streetaddress=3000 Tanery way")
					So(idt.Identity, ShouldContain, "streetaddress=13 Rue de Mogador")
					So(idt.Identity, ShouldContain, "email=me@me.com")
					So(idt.Identity, ShouldContain, "serialnumber=42")
				})
			})

			Convey("Calling FromCertificate with a invalid user cert should work", func() {

				err := iss.FromCertificate(usercert2)
				So(err, ShouldNotBeNil)
				So(err.Error(), ShouldStartWith, "unable Verify certificate: x509: certificate signed by unknown authority")
			})
		})
	})
}
