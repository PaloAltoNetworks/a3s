package mtlsissuer

import (
	"context"
	"crypto"
	"crypto/x509"
	"crypto/x509/pkix"
	"encoding/json"
	"encoding/pem"
	"fmt"
	"math/big"
	"net/http"
	"net/http/httptest"
	"testing"

	. "github.com/smartystreets/goconvey/convey"
	"go.aporeto.io/a3s/pkgs/api"
	"go.aporeto.io/a3s/pkgs/token"
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
		usercert1, userkey1 := getECCert(
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

			block, _ := tglib.CertToPEM(cacert1)

			src := api.NewMTLSSource()
			src.Name = "mysource"
			src.Namespace = "/my/ns"
			src.CA = string(pem.EncodeToMemory(block))
			iss := newMTLSIssuer(src)
			So(iss.token, ShouldNotBeNil)
			So(iss.source, ShouldEqual, src)

			Convey("Calling FromCertificate with a source missing a CA", func() {
				src.CA = ""
				_, err := New(context.Background(), src, usercert2)
				So(err, ShouldNotBeNil)
				So(err.Error(), ShouldStartWith, "unable to prepare x509 verifier: could not append cert from source.CA")
			})

			Convey("Calling FromCertificate with a valid user cert should work", func() {

				iss, err := New(context.Background(), src, usercert1)
				So(err, ShouldBeNil)

				idt := iss.Issue()
				So(len(idt.Identity), ShouldEqual, 22)
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
				So(idt.Identity, ShouldContain, fmt.Sprintf("akid=%02X", cacert1.SubjectKeyId))
				So(idt.Identity, ShouldContain, fmt.Sprintf("fingerprint=%s", token.Fingerprint(usercert1)))
				So(idt.Identity, ShouldContain, fmt.Sprintf("issuerchain=%s", token.Fingerprint(cacert1)))
			})

			Convey("Calling FromCertificate with a valid modifier should work", func() {

				ts := httptest.NewTLSServer(http.HandlerFunc(func(w http.ResponseWriter, req *http.Request) {
					d, _ := json.Marshal([]string{"aa=aa", "bb=bb"})
					w.WriteHeader(http.StatusOK)
					w.Write(d) // nolint
				}))
				defer ts.Close()

				cab, _ := tglib.CertToPEM(ts.Certificate())
				certb, _ := tglib.CertToPEM(usercert1)
				keyb, _ := tglib.KeyToPEM(userkey1)
				src.Modifier = api.NewIdentityModifier()
				src.Modifier.CA = string(pem.EncodeToMemory(cab))
				src.Modifier.URL = ts.URL
				src.Modifier.Certificate = string(pem.EncodeToMemory(certb))
				src.Modifier.Key = string(pem.EncodeToMemory(keyb))

				iss, err := New(context.Background(), src, usercert1)
				So(err, ShouldBeNil)

				idt := iss.Issue()
				So(len(idt.Identity), ShouldEqual, 2)
				So(idt.Identity, ShouldContain, "aa=aa")
				So(idt.Identity, ShouldContain, "bb=bb")
			})

			Convey("Calling FromCertificate with a modifier with missing tls info", func() {

				ts := httptest.NewTLSServer(http.HandlerFunc(func(w http.ResponseWriter, req *http.Request) {
					d, _ := json.Marshal([]string{"aa=aa", "bb=bb"})
					w.WriteHeader(http.StatusOK)
					w.Write(d) // nolint
				}))
				defer ts.Close()

				cab, _ := tglib.CertToPEM(ts.Certificate())
				src.Modifier = api.NewIdentityModifier()
				src.Modifier.CA = string(pem.EncodeToMemory(cab))
				src.Modifier.URL = ts.URL

				err := iss.fromCertificate(context.Background(), usercert1)
				So(err, ShouldNotBeNil)
				So(err.Error(), ShouldEqual, `unable to prepare source modifier: unable to create certificate: could not read key data from bytes: ''`)
			})

			Convey("Calling FromCertificate with a valid modifier returning an error", func() {

				ts := httptest.NewTLSServer(http.HandlerFunc(func(w http.ResponseWriter, req *http.Request) {
					w.WriteHeader(http.StatusForbidden)
				}))
				defer ts.Close()

				cab, _ := tglib.CertToPEM(ts.Certificate())
				certb, _ := tglib.CertToPEM(usercert1)
				keyb, _ := tglib.KeyToPEM(userkey1)
				src.Modifier = api.NewIdentityModifier()
				src.Modifier.CA = string(pem.EncodeToMemory(cab))
				src.Modifier.URL = ts.URL
				src.Modifier.Certificate = string(pem.EncodeToMemory(certb))
				src.Modifier.Key = string(pem.EncodeToMemory(keyb))

				err := iss.fromCertificate(context.Background(), usercert1)
				So(err, ShouldNotBeNil)
				So(err.Error(), ShouldEqual, `unable to call modifier: service returned an error: 403 Forbidden`)
			})

			Convey("Calling FromCertificate with a invalid user cert should work", func() {

				err := iss.fromCertificate(context.Background(), usercert2)
				So(err, ShouldNotBeNil)
				So(err.Error(), ShouldStartWith, "unable to verify certificate: x509: certificate signed by unknown authority")
			})
		})
	})
}
