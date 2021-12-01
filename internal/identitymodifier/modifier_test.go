package identitymodifier

import (
	"context"
	"crypto"
	"crypto/x509"
	"crypto/x509/pkix"
	"encoding/json"
	"encoding/pem"
	"io"
	"net/http"
	"net/http/httptest"
	"net/url"
	"testing"

	. "github.com/smartystreets/goconvey/convey"
	"go.aporeto.io/a3s/pkgs/api"
	"go.aporeto.io/a3s/pkgs/token"
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

func TestNew(t *testing.T) {

	Convey("Calling New should work", t, func() {
		xc, xk := getECCert()
		xcb, _ := tglib.CertToPEM(xc)
		xck, _ := tglib.KeyToPEM(xk)
		cert := pem.EncodeToMemory(xcb)
		key := pem.EncodeToMemory(xck)

		Convey("Everything is valid", func() {
			mi := &api.IdentityModifier{
				Certificate: string(cert),
				Key:         string(key),
				URL:         "https://toto.com",
				Method:      api.IdentityModifierMethodGET,
			}
			mm, err := NewRemote(mi, token.Source{})
			So(err, ShouldBeNil)
			m := mm.(*identityModifier)
			So(m.clientCert, ShouldNotBeNil)
			So(m.caPool, ShouldNotBeNil)
			So(m.url, ShouldEqual, "https://toto.com")
		})

		Convey("Invalid method", func() {
			mi := &api.IdentityModifier{
				Certificate: string(cert),
				Key:         string(key),
				URL:         "https://toto.com",
				Method:      api.IdentityModifierMethodValue("invalid"),
			}
			_, err := NewRemote(mi, token.Source{})
			So(err, ShouldNotBeNil)
			So(err.Error(), ShouldEqual, "invalid http method: invalid")
		})

		Convey("Invalid cert", func() {
			mi := &api.IdentityModifier{
				Certificate: "coucou",
				Key:         string(key),
				URL:         "https://toto.com",
				Method:      api.IdentityModifierMethodGET,
			}
			_, err := NewRemote(mi, token.Source{})
			So(err, ShouldNotBeNil)
			So(err.Error(), ShouldEqual, "unable to create certificate: tls: failed to find any PEM data in certificate input")
		})
	})
}

func TestModify(t *testing.T) {

	Convey("Given I have a http server", t, func() {

		var incomingBody []byte
		var incomingQuery url.Values
		var incomingHeader http.Header
		returnCode := http.StatusOK
		var returnBody []byte
		ts := httptest.NewTLSServer(http.HandlerFunc(func(w http.ResponseWriter, req *http.Request) {
			incomingBody, _ = io.ReadAll(req.Body)
			d, _ := json.Marshal([]string{"aa=aa", "bb=bb"})
			incomingQuery = req.URL.Query()
			incomingHeader = req.Header.Clone()
			w.WriteHeader(returnCode)
			if returnBody == nil {
				w.Write(d) // nolint
			} else {
				w.Write(returnBody) // nolint
			}
		}))
		defer ts.Close()

		xc, xk := getECCert()
		xcb, _ := tglib.CertToPEM(xc)
		xck, _ := tglib.KeyToPEM(xk)
		cert := pem.EncodeToMemory(xcb)
		key := pem.EncodeToMemory(xck)
		cab, _ := tglib.CertToPEM(ts.Certificate())
		ca := pem.EncodeToMemory(cab)

		Convey("Setting a bad URL should work", func() {
			mi := &api.IdentityModifier{
				CA:          string(ca),
				Certificate: string(cert),
				Key:         string(key),
				URL:         "nothttp://com",
				Method:      api.IdentityModifierMethodPOST,
			}
			m, _ := NewRemote(mi, token.Source{Type: "type", Namespace: "/ns", Name: "name"})
			out, err := m.Modify(context.Background(), []string{"a=a", "b=b"})
			So(err, ShouldNotBeNil)
			So(err.Error(), ShouldEqual, `unable to perform request: Post "nothttp://com": unsupported protocol scheme "nothttp"`)
			So(out, ShouldBeNil)
		})

		Convey("Setting a nil context should work", func() {
			mi := &api.IdentityModifier{
				CA:          string(ca),
				Certificate: string(cert),
				Key:         string(key),
				URL:         ts.URL,
				Method:      api.IdentityModifierMethodGET,
			}
			m, _ := NewRemote(mi, token.Source{Type: "type", Namespace: "/ns", Name: "name"})
			out, err := m.Modify(nil, []string{"a=a", "b=b"}) // nolint
			So(err, ShouldNotBeNil)
			So(err.Error(), ShouldEqual, `unable to build http request: net/http: nil Context`)
			So(out, ShouldBeNil)
		})

		Convey("Using POST should work", func() {
			mi := &api.IdentityModifier{
				CA:          string(ca),
				Certificate: string(cert),
				Key:         string(key),
				URL:         ts.URL,
				Method:      api.IdentityModifierMethodPOST,
			}
			m, _ := NewRemote(mi, token.Source{Type: "type", Namespace: "/ns", Name: "name"})
			out, err := m.Modify(context.Background(), []string{"a=a", "b=b"})
			So(err, ShouldBeNil)
			So(out, ShouldResemble, []string{"aa=aa", "bb=bb"})
			So(incomingBody, ShouldResemble, []byte(`["a=a","b=b"]`))
			So(len(incomingQuery), ShouldEqual, 0)
			So(incomingHeader.Get("x-a3s-source-type"), ShouldEqual, "type")
			So(incomingHeader.Get("x-a3s-source-namespace"), ShouldEqual, "/ns")
			So(incomingHeader.Get("x-a3s-source-name"), ShouldEqual, "name")
		})

		Convey("Using PUT should work", func() {
			mi := &api.IdentityModifier{
				CA:          string(ca),
				Certificate: string(cert),
				Key:         string(key),
				URL:         ts.URL,
				Method:      api.IdentityModifierMethodPUT,
			}
			m, _ := NewRemote(mi, token.Source{Type: "type", Namespace: "/ns", Name: "name"})
			out, err := m.Modify(context.Background(), []string{"a=a", "b=b"})
			So(err, ShouldBeNil)
			So(out, ShouldResemble, []string{"aa=aa", "bb=bb"})
			So(incomingBody, ShouldResemble, []byte(`["a=a","b=b"]`))
			So(len(incomingQuery), ShouldEqual, 0)
			So(incomingHeader.Get("x-a3s-source-type"), ShouldEqual, "type")
			So(incomingHeader.Get("x-a3s-source-namespace"), ShouldEqual, "/ns")
			So(incomingHeader.Get("x-a3s-source-name"), ShouldEqual, "name")
		})

		Convey("Using PATCH should work", func() {
			mi := &api.IdentityModifier{
				CA:          string(ca),
				Certificate: string(cert),
				Key:         string(key),
				URL:         ts.URL,
				Method:      api.IdentityModifierMethodPATCH,
			}
			m, _ := NewRemote(mi, token.Source{Type: "type", Namespace: "/ns", Name: "name"})
			out, err := m.Modify(context.Background(), []string{"a=a", "b=b"})
			So(err, ShouldBeNil)
			So(out, ShouldResemble, []string{"aa=aa", "bb=bb"})
			So(incomingBody, ShouldResemble, []byte(`["a=a","b=b"]`))
			So(len(incomingQuery), ShouldEqual, 0)
			So(incomingHeader.Get("x-a3s-source-type"), ShouldEqual, "type")
			So(incomingHeader.Get("x-a3s-source-namespace"), ShouldEqual, "/ns")
			So(incomingHeader.Get("x-a3s-source-name"), ShouldEqual, "name")
		})

		Convey("Using GET should work", func() {
			mi := &api.IdentityModifier{
				CA:          string(ca),
				Certificate: string(cert),
				Key:         string(key),
				URL:         ts.URL,
				Method:      api.IdentityModifierMethodGET,
			}
			m, _ := NewRemote(mi, token.Source{Type: "type", Namespace: "/ns", Name: "name"})
			out, err := m.Modify(context.Background(), []string{"a=a", "b=b"})
			So(err, ShouldBeNil)
			So(out, ShouldResemble, []string{"aa=aa", "bb=bb"})
			So(incomingBody, ShouldBeEmpty)
			So(len(incomingQuery), ShouldEqual, 1)
			So(len(incomingQuery["claim"]), ShouldEqual, 2)
			So(incomingQuery["claim"][0], ShouldEqual, "a=a")
			So(incomingQuery["claim"][1], ShouldEqual, "b=b")
			So(incomingHeader.Get("x-a3s-source-type"), ShouldEqual, "type")
			So(incomingHeader.Get("x-a3s-source-namespace"), ShouldEqual, "/ns")
			So(incomingHeader.Get("x-a3s-source-name"), ShouldEqual, "name")
		})

		Convey("When the server returns 204", func() {
			returnCode = http.StatusNoContent
			mi := &api.IdentityModifier{
				CA:          string(ca),
				Certificate: string(cert),
				Key:         string(key),
				URL:         ts.URL,
				Method:      api.IdentityModifierMethodPOST,
			}
			m, _ := NewRemote(mi, token.Source{Type: "type", Namespace: "/ns", Name: "name"})
			out, err := m.Modify(context.Background(), []string{"a=a", "b=b"})
			So(err, ShouldBeNil)
			So(out, ShouldResemble, []string{"a=a", "b=b"})
			So(incomingBody, ShouldResemble, []byte(`["a=a","b=b"]`))
			So(len(incomingQuery), ShouldEqual, 0)
			So(incomingHeader.Get("x-a3s-source-type"), ShouldEqual, "type")
			So(incomingHeader.Get("x-a3s-source-namespace"), ShouldEqual, "/ns")
			So(incomingHeader.Get("x-a3s-source-name"), ShouldEqual, "name")
		})

		Convey("When the server returns something else", func() {
			returnCode = http.StatusTeapot
			mi := &api.IdentityModifier{
				CA:          string(ca),
				Certificate: string(cert),
				Key:         string(key),
				URL:         ts.URL,
				Method:      api.IdentityModifierMethodPOST,
			}
			m, _ := NewRemote(mi, token.Source{Type: "type", Namespace: "/ns", Name: "name"})
			out, err := m.Modify(context.Background(), []string{"a=a", "b=b"})
			So(err, ShouldNotBeNil)
			So(err.Error(), ShouldEqual, `service returned an error: 418 I'm a teapot`)
			So(out, ShouldBeNil)
		})

		Convey("When the server returns a forbidden body", func() {
			returnBody = []byte(`["@a=a"]`)
			mi := &api.IdentityModifier{
				CA:          string(ca),
				Certificate: string(cert),
				Key:         string(key),
				URL:         ts.URL,
				Method:      api.IdentityModifierMethodPOST,
			}
			m, _ := NewRemote(mi, token.Source{Type: "type", Namespace: "/ns", Name: "name"})
			out, err := m.Modify(context.Background(), []string{"a=a", "b=b"})
			So(err, ShouldNotBeNil)
			So(err.Error(), ShouldEqual, `invalid returned claim '@a=a': must not be prefixed by @`)
			So(out, ShouldBeNil)
		})
		Convey("When the server returns a bad body", func() {
			returnBody = []byte("oh no")
			mi := &api.IdentityModifier{
				CA:          string(ca),
				Certificate: string(cert),
				Key:         string(key),
				URL:         ts.URL,
				Method:      api.IdentityModifierMethodPOST,
			}
			m, _ := NewRemote(mi, token.Source{Type: "type", Namespace: "/ns", Name: "name"})
			out, err := m.Modify(context.Background(), []string{"a=a", "b=b"})
			So(err, ShouldNotBeNil)
			So(err.Error(), ShouldEqual, `unable to decode response: invalid character 'o' looking for beginning of value`)
			So(out, ShouldBeNil)
		})
	})
}
