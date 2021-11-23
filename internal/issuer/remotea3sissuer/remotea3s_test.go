package remotea3sissuer

import (
	"context"
	"crypto"
	"crypto/sha1"
	"crypto/x509"
	"crypto/x509/pkix"
	"encoding/json"
	"encoding/pem"
	"fmt"
	"net/http"
	"net/http/httptest"
	"testing"
	"time"

	. "github.com/smartystreets/goconvey/convey"
	"go.aporeto.io/a3s/pkgs/api"
	"go.aporeto.io/a3s/pkgs/token"
	"go.aporeto.io/tg/tglib"
)

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

func TestErrRemoteA3S(t *testing.T) {
	Convey("ErrRemoteA3S should work", t, func() {
		e := fmt.Errorf("boom")
		err := ErrRemoteA3S{Err: e}
		So(err.Error(), ShouldEqual, "remote a3s error: boom")
		So(err.Unwrap(), ShouldEqual, e)
	})
}

func TestNewRemoteA3SIssuer(t *testing.T) {
	Convey("Calling newRemoteA3SIssuer should work", t, func() {
		source := &api.A3SSource{
			Name:      "name",
			Namespace: "/ns",
		}
		iss := newRemoteA3SIssuer(source)
		So(iss.source, ShouldEqual, source)
		So(iss.Issue().Source.Type, ShouldEqual, "remotea3s")
		So(iss.Issue().Source.Namespace, ShouldEqual, "/ns")
		So(iss.Issue().Source.Name, ShouldEqual, "name")
	})
}

func TestNew(t *testing.T) {

	Convey("Given an http server and a A3SSource and everything is fine", t, func() {

		remoteCert, remoteKey := getECCert(pkix.Name{CommonName: "local"})
		kid := fmt.Sprintf("%02X", sha1.Sum(remoteCert.Raw))

		ts := httptest.NewServer(
			http.HandlerFunc(func(w http.ResponseWriter, req *http.Request) {
				j := token.NewJWKS()
				_ = j.Append(remoteCert)
				d, _ := json.Marshal(j)
				w.Write(d) // nolint
			}),
		)

		rtok := &token.IdentityToken{
			Source: token.Source{
				Type:      "mtls",
				Name:      "remotename",
				Namespace: "/remotens",
			},
			Identity: []string{"remote=claim"},
		}

		rtokString, _ := rtok.JWT(
			remoteKey,
			kid,
			ts.URL,
			[]string{"local"},
			time.Now().Add(time.Minute),
			nil,
		)

		source := &api.A3SSource{
			Name:      "name",
			Namespace: "/ns",
			Issuer:    ts.URL,
			Audience:  "local",
		}

		iss, err := New(context.Background(), source, rtokString)
		So(err, ShouldBeNil)
		So(iss.Issue().Identity, ShouldResemble, []string{
			"remote=claim",
			"sourcetype=mtls",
			"sourcenamespace=/remotens",
			"sourcename=remotename",
			fmt.Sprintf("iss=%s", ts.URL),
		})
	})

	Convey("Given an http server and a A3SSource but token signature is unknown", t, func() {

		remoteCert, remoteKey := getECCert(pkix.Name{CommonName: "local"})
		kid := fmt.Sprintf("%02X", sha1.Sum(remoteCert.Raw))

		ts := httptest.NewServer(
			http.HandlerFunc(func(w http.ResponseWriter, req *http.Request) {
				j := token.NewJWKS()
				d, _ := json.Marshal(j)
				w.Write(d) // nolint
			}),
		)

		rtok := &token.IdentityToken{
			Source: token.Source{
				Type:      "mtls",
				Name:      "remotename",
				Namespace: "/remotens",
			},
			Identity: []string{"remote=claim"},
		}

		rtokString, _ := rtok.JWT(
			remoteKey,
			kid,
			ts.URL,
			[]string{"local"},
			time.Now().Add(time.Minute),
			nil,
		)

		source := &api.A3SSource{
			Name:      "name",
			Namespace: "/ns",
			Issuer:    ts.URL,
			Audience:  "local",
		}

		iss, err := New(context.Background(), source, rtokString)
		So(iss, ShouldBeNil)
		So(err, ShouldNotBeNil)
		So(err.Error(), ShouldStartWith, `remote a3s error: unable to parse remote a3s token: unable to parse jwt: unable to find kid`)
	})

	Convey("Given an http server and a A3SSource but remote jwks is unavailable", t, func() {

		source := &api.A3SSource{
			Name:      "name",
			Namespace: "/ns",
			Issuer:    "toto:/not-a-url",
			Audience:  "local",
		}

		iss, err := New(context.Background(), source, "rtokString")
		So(iss, ShouldBeNil)
		So(err, ShouldNotBeNil)
		So(err.Error(), ShouldEqual, `remote a3s error: unable to retrieve remote jwks: remote jwks error: unable to send request: Get "toto:///not-a-url/.well-known/jwks.json": unsupported protocol scheme "toto"`)
	})
}
