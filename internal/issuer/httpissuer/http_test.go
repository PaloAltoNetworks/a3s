package httpissuer

import (
	"context"
	"crypto"
	"crypto/x509"
	"crypto/x509/pkix"
	"encoding/json"
	"encoding/pem"
	"fmt"
	"io"
	"net/http"
	"net/http/httptest"
	"testing"

	. "github.com/smartystreets/goconvey/convey"
	"go.aporeto.io/a3s/pkgs/api"
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

func TestErrs(t *testing.T) {

	Convey("ErrHTTP should work", t, func() {
		e := fmt.Errorf("boom")
		err := ErrHTTP{Err: e}
		So(err.Error(), ShouldEqual, "http error: boom")
	})

	Convey("ErrHTTPResponse should work", t, func() {
		e := fmt.Errorf("boom")
		err := ErrHTTPResponse{Err: e}
		So(err.Error(), ShouldEqual, "http response error: boom")
	})
}

func TestNewHTTPIssuer(t *testing.T) {

	Convey("Calling newHTTPIssue should work", t, func() {
		source := &api.HTTPSource{
			Name:      "name",
			Namespace: "/ns",
		}
		iss := newHTTPIssuer(source)
		So(iss.source, ShouldEqual, source)
		So(iss.Issue().Source.Type, ShouldEqual, "http")
		So(iss.Issue().Source.Name, ShouldEqual, "name")
		So(iss.Issue().Source.Namespace, ShouldEqual, "/ns")
	})
}

func TestNew(t *testing.T) {

	Convey("Given an http server and an HTTPSource", t, func() {

		remoteCert, remoteKey := getECCert(pkix.Name{CommonName: "local"})

		Convey("When url is not correctly formed", func() {

			source := &api.HTTPSource{
				Name:        "name",
				Namespace:   "/ns",
				Endpoint:    "nothttp://coucou.com",
				Certificate: func() string { c, _ := tglib.CertToPEM(remoteCert); return string(pem.EncodeToMemory(c)) }(),
				Key:         func() string { c, _ := tglib.KeyToPEM(remoteKey); return string(pem.EncodeToMemory(c)) }(),
			}

			_, err := New(context.Background(), source, "user", "pass")
			So(err, ShouldNotBeNil)
			So(err.Error(), ShouldEqual, `http error: unable to send request: Post "nothttp://coucou.com": unsupported protocol scheme "nothttp"`)
		})

		Convey("When certificate is not correctly formed", func() {

			source := &api.HTTPSource{
				Name:        "name",
				Namespace:   "/ns",
				Endpoint:    "nothttp://coucou.com",
				Certificate: "not-pem",
				Key:         func() string { c, _ := tglib.KeyToPEM(remoteKey); return string(pem.EncodeToMemory(c)) }(),
			}

			_, err := New(context.Background(), source, "user", "pass")
			So(err, ShouldNotBeNil)
			So(err.Error(), ShouldEqual, `http error: unable to read certificate: tls: failed to find any PEM data in certificate input`)
		})

		Convey("When key is not correctly formed", func() {

			source := &api.HTTPSource{
				Name:        "name",
				Namespace:   "/ns",
				Endpoint:    "nothttp://coucou.com",
				Certificate: func() string { c, _ := tglib.CertToPEM(remoteCert); return string(pem.EncodeToMemory(c)) }(),
				Key:         "not-pem",
			}

			_, err := New(context.Background(), source, "user", "pass")
			So(err, ShouldNotBeNil)
			So(err.Error(), ShouldEqual, `http error: unable to read certificate: could not read key data from bytes: 'not-pem'`)
		})

		Convey("When everything is fine and there is no modifier", func() {

			var expectedUser string
			var expectedPass string

			ts := httptest.NewServer(
				http.HandlerFunc(func(w http.ResponseWriter, req *http.Request) {
					data, _ := io.ReadAll(req.Body)
					d := map[string]string{}
					_ = json.Unmarshal(data, &d)
					expectedUser = d["username"]
					expectedPass = d["password"]
					w.Write([]byte(`["k=v"]`)) // nolint
				}),
			)

			source := &api.HTTPSource{
				Name:        "name",
				Namespace:   "/ns",
				Endpoint:    ts.URL,
				Certificate: func() string { c, _ := tglib.CertToPEM(remoteCert); return string(pem.EncodeToMemory(c)) }(),
				Key:         func() string { c, _ := tglib.KeyToPEM(remoteKey); return string(pem.EncodeToMemory(c)) }(),
			}

			iss, err := New(context.Background(), source, "user", "pass")
			So(err, ShouldBeNil)
			So(iss.Issue().Identity, ShouldResemble, []string{"k=v"})
			So(expectedUser, ShouldEqual, "user")
			So(expectedPass, ShouldEqual, "pass")
		})

		Convey("When server returns an error", func() {

			ts := httptest.NewServer(
				http.HandlerFunc(func(w http.ResponseWriter, req *http.Request) {
					w.WriteHeader(http.StatusForbidden)
				}),
			)

			source := &api.HTTPSource{
				Name:        "name",
				Namespace:   "/ns",
				Endpoint:    ts.URL,
				Certificate: func() string { c, _ := tglib.CertToPEM(remoteCert); return string(pem.EncodeToMemory(c)) }(),
				Key:         func() string { c, _ := tglib.KeyToPEM(remoteKey); return string(pem.EncodeToMemory(c)) }(),
			}

			_, err := New(context.Background(), source, "user", "pass")
			So(err, ShouldNotBeNil)
			So(err.Error(), ShouldEqual, `http response error: server responded with '403 Forbidden'`)
		})

		Convey("When server returns a malformed body", func() {

			ts := httptest.NewServer(
				http.HandlerFunc(func(w http.ResponseWriter, req *http.Request) {
					w.Write([]byte("looool")) // nolint
				}),
			)

			source := &api.HTTPSource{
				Name:        "name",
				Namespace:   "/ns",
				Endpoint:    ts.URL,
				Certificate: func() string { c, _ := tglib.CertToPEM(remoteCert); return string(pem.EncodeToMemory(c)) }(),
				Key:         func() string { c, _ := tglib.KeyToPEM(remoteKey); return string(pem.EncodeToMemory(c)) }(),
			}

			_, err := New(context.Background(), source, "user", "pass")
			So(err, ShouldNotBeNil)
			So(err.Error(), ShouldEqual, `http response error: unable to decode response body: invalid character 'l' looking for beginning of value`)
		})

		Convey("When there is a modifier and everything is fine ", func() {

			ts := httptest.NewServer(
				http.HandlerFunc(func(w http.ResponseWriter, req *http.Request) {
					w.Write([]byte(`["k=v"]`)) // nolint
				}),
			)

			ts2 := httptest.NewTLSServer(http.HandlerFunc(func(w http.ResponseWriter, req *http.Request) {
				d, _ := json.Marshal([]string{"aa=aa", "bb=bb"})
				w.WriteHeader(http.StatusOK)
				w.Write(d) // nolint
			}))
			defer ts2.Close()

			usercert1, userkey1 := getECCert(pkix.Name{})

			source := &api.HTTPSource{
				Name:        "name",
				Namespace:   "/ns",
				Endpoint:    ts.URL,
				Certificate: func() string { c, _ := tglib.CertToPEM(remoteCert); return string(pem.EncodeToMemory(c)) }(),
				Key:         func() string { c, _ := tglib.KeyToPEM(remoteKey); return string(pem.EncodeToMemory(c)) }(),
				Modifier: &api.IdentityModifier{
					CA:          func() string { c, _ := tglib.CertToPEM(ts2.Certificate()); return string(pem.EncodeToMemory(c)) }(),
					URL:         ts2.URL,
					Certificate: func() string { c, _ := tglib.CertToPEM(usercert1); return string(pem.EncodeToMemory(c)) }(),
					Key:         func() string { c, _ := tglib.KeyToPEM(userkey1); return string(pem.EncodeToMemory(c)) }(),
					Method:      api.IdentityModifierMethodPOST,
				},
			}

			iss, err := New(context.Background(), source, "user", "pass")
			So(err, ShouldBeNil)
			So(iss.Issue().Identity, ShouldResemble, []string{"aa=aa", "bb=bb"})
		})

		Convey("When there is a modifier with missing TLS ", func() {

			ts := httptest.NewServer(
				http.HandlerFunc(func(w http.ResponseWriter, req *http.Request) {
					w.Write([]byte(`["k=v"]`)) // nolint
				}),
			)

			ts2 := httptest.NewTLSServer(http.HandlerFunc(func(w http.ResponseWriter, req *http.Request) {
				d, _ := json.Marshal([]string{"aa=aa", "bb=bb"})
				w.WriteHeader(http.StatusOK)
				w.Write(d) // nolint
			}))
			defer ts2.Close()

			source := &api.HTTPSource{
				Name:        "name",
				Namespace:   "/ns",
				Endpoint:    ts.URL,
				Certificate: func() string { c, _ := tglib.CertToPEM(remoteCert); return string(pem.EncodeToMemory(c)) }(),
				Key:         func() string { c, _ := tglib.KeyToPEM(remoteKey); return string(pem.EncodeToMemory(c)) }(),
				Modifier: &api.IdentityModifier{
					CA:     func() string { c, _ := tglib.CertToPEM(ts2.Certificate()); return string(pem.EncodeToMemory(c)) }(),
					URL:    ts2.URL,
					Method: api.IdentityModifierMethodPOST,
				},
			}

			_, err := New(context.Background(), source, "user", "pass")
			So(err, ShouldNotBeNil)
			So(err.Error(), ShouldEqual, `unable to prepare source modifier: unable to create certificate: could not read key data from bytes: ''`)
		})

		Convey("When there is a modifier but it returns an error", func() {

			ts := httptest.NewServer(
				http.HandlerFunc(func(w http.ResponseWriter, req *http.Request) {
					w.Write([]byte(`["k=v"]`)) // nolint
				}),
			)

			ts2 := httptest.NewTLSServer(http.HandlerFunc(func(w http.ResponseWriter, req *http.Request) {
				w.WriteHeader(http.StatusForbidden)
			}))
			defer ts2.Close()

			usercert1, userkey1 := getECCert(pkix.Name{})

			source := &api.HTTPSource{
				Name:        "name",
				Namespace:   "/ns",
				Endpoint:    ts.URL,
				Certificate: func() string { c, _ := tglib.CertToPEM(remoteCert); return string(pem.EncodeToMemory(c)) }(),
				Key:         func() string { c, _ := tglib.KeyToPEM(remoteKey); return string(pem.EncodeToMemory(c)) }(),
				Modifier: &api.IdentityModifier{
					CA:          func() string { c, _ := tglib.CertToPEM(ts2.Certificate()); return string(pem.EncodeToMemory(c)) }(),
					URL:         ts2.URL,
					Certificate: func() string { c, _ := tglib.CertToPEM(usercert1); return string(pem.EncodeToMemory(c)) }(),
					Key:         func() string { c, _ := tglib.KeyToPEM(userkey1); return string(pem.EncodeToMemory(c)) }(),
					Method:      api.IdentityModifierMethodPOST,
				},
			}

			_, err := New(context.Background(), source, "user", "pass")
			So(err, ShouldNotBeNil)
			So(err.Error(), ShouldEqual, `unable to call modifier: service returned an error: 403 Forbidden`)
		})

	})
}
