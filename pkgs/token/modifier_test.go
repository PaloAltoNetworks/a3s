package token

import (
	"context"
	"encoding/json"
	"encoding/pem"
	"io"
	"net/http"
	"net/http/httptest"
	"net/url"
	"testing"

	. "github.com/smartystreets/goconvey/convey"
	"go.aporeto.io/tg/tglib"
)

func TestNew(t *testing.T) {

	Convey("Calling New should work", t, func() {
		xc, xk := getECCert()
		xcb, _ := tglib.CertToPEM(xc)
		xck, _ := tglib.KeyToPEM(xk)
		cert := pem.EncodeToMemory(xcb)
		key := pem.EncodeToMemory(xck)

		Convey("Everything is valid", func() {
			mm, err := NewHTTPIdentityModifier("https://toto.com", "get", nil, cert, key)
			So(err, ShouldBeNil)
			m := mm.(*identityModifier)
			So(m.clientCert, ShouldNotBeNil)
			So(m.caPool, ShouldNotBeNil)
			So(m.url, ShouldEqual, "https://toto.com")
		})

		Convey("Invalid method", func() {
			_, err := NewHTTPIdentityModifier("https://toto.com", "invalid", nil, cert, key)
			So(err, ShouldNotBeNil)
			So(err.Error(), ShouldEqual, "invalid http method: invalid")
		})

		Convey("Invalid cert", func() {
			_, err := NewHTTPIdentityModifier("https://toto.com", "get", nil, []byte("coucou"), key)
			So(err, ShouldNotBeNil)
			So(err.Error(), ShouldEqual, "unable create certificate: tls: failed to find any PEM data in certificate input")
		})
	})
}

func TestModify(t *testing.T) {

	Convey("Given I have a http server", t, func() {

		var incomingBody []byte
		var incomingQuery url.Values
		returnCode := http.StatusOK
		var returnBody []byte
		ts := httptest.NewTLSServer(http.HandlerFunc(func(w http.ResponseWriter, req *http.Request) {
			incomingBody, _ = io.ReadAll(req.Body)
			d, _ := json.Marshal([]string{"aa=aa", "bb=bb"})
			incomingQuery = req.URL.Query()
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
			m, _ := NewHTTPIdentityModifier("nothttp://com", http.MethodPost, ca, cert, key)
			out, err := m.Modify(context.Background(), []string{"a=a", "b=b"})
			So(err, ShouldNotBeNil)
			So(err.Error(), ShouldEqual, `unable to perform request: Post "nothttp://com": unsupported protocol scheme "nothttp"`)
			So(out, ShouldBeNil)
		})

		Convey("Setting a nil context should work", func() {
			m, _ := NewHTTPIdentityModifier(ts.URL, http.MethodGet, ca, cert, key)
			out, err := m.Modify(nil, []string{"a=a", "b=b"})
			So(err, ShouldNotBeNil)
			So(err.Error(), ShouldEqual, `unable to build http request: net/http: nil Context`)
			So(out, ShouldBeNil)
		})

		Convey("Using POST should work", func() {
			m, _ := NewHTTPIdentityModifier(ts.URL, http.MethodPost, ca, cert, key)
			out, err := m.Modify(context.Background(), []string{"a=a", "b=b"})
			So(err, ShouldBeNil)
			So(out, ShouldResemble, []string{"aa=aa", "bb=bb"})
			So(incomingBody, ShouldResemble, []byte(`["a=a","b=b"]`))
			So(len(incomingQuery), ShouldEqual, 0)
		})

		Convey("Using PUT should work", func() {
			m, _ := NewHTTPIdentityModifier(ts.URL, http.MethodPut, ca, cert, key)
			out, err := m.Modify(context.Background(), []string{"a=a", "b=b"})
			So(err, ShouldBeNil)
			So(out, ShouldResemble, []string{"aa=aa", "bb=bb"})
			So(incomingBody, ShouldResemble, []byte(`["a=a","b=b"]`))
			So(len(incomingQuery), ShouldEqual, 0)
		})

		Convey("Using PATCH should work", func() {
			m, _ := NewHTTPIdentityModifier(ts.URL, http.MethodPatch, ca, cert, key)
			out, err := m.Modify(context.Background(), []string{"a=a", "b=b"})
			So(err, ShouldBeNil)
			So(out, ShouldResemble, []string{"aa=aa", "bb=bb"})
			So(incomingBody, ShouldResemble, []byte(`["a=a","b=b"]`))
			So(len(incomingQuery), ShouldEqual, 0)
		})

		Convey("Using GET should work", func() {
			m, _ := NewHTTPIdentityModifier(ts.URL, http.MethodGet, ca, cert, key)
			out, err := m.Modify(context.Background(), []string{"a=a", "b=b"})
			So(err, ShouldBeNil)
			So(out, ShouldResemble, []string{"aa=aa", "bb=bb"})
			So(incomingBody, ShouldBeEmpty)
			So(len(incomingQuery), ShouldEqual, 1)
			So(len(incomingQuery["claim"]), ShouldEqual, 2)
			So(incomingQuery["claim"][0], ShouldEqual, "a=a")
			So(incomingQuery["claim"][1], ShouldEqual, "b=b")
		})

		Convey("When the server returns 204", func() {
			returnCode = http.StatusNoContent
			m, _ := NewHTTPIdentityModifier(ts.URL, http.MethodPost, ca, cert, key)
			out, err := m.Modify(context.Background(), []string{"a=a", "b=b"})
			So(err, ShouldBeNil)
			So(out, ShouldResemble, []string{"a=a", "b=b"})
			So(incomingBody, ShouldResemble, []byte(`["a=a","b=b"]`))
			So(len(incomingQuery), ShouldEqual, 0)
		})

		Convey("When the server returns something else", func() {
			returnCode = http.StatusTeapot
			m, _ := NewHTTPIdentityModifier(ts.URL, http.MethodPost, ca, cert, key)
			out, err := m.Modify(context.Background(), []string{"a=a", "b=b"})
			So(err, ShouldNotBeNil)
			So(err.Error(), ShouldEqual, `service returned an error: 418 I'm a teapot`)
			So(out, ShouldBeNil)
		})

		Convey("When the server returns a forbidden body", func() {
			returnBody = []byte(`["@a=a"]`)
			m, _ := NewHTTPIdentityModifier(ts.URL, http.MethodPost, ca, cert, key)
			out, err := m.Modify(context.Background(), []string{"a=a", "b=b"})
			So(err, ShouldNotBeNil)
			So(err.Error(), ShouldEqual, `invalid returned claim '@a=a': must not be prefixed by @`)
			So(out, ShouldBeNil)
		})
		Convey("When the server returns a bad body", func() {
			returnBody = []byte("oh no")
			m, _ := NewHTTPIdentityModifier(ts.URL, http.MethodPost, ca, cert, key)
			out, err := m.Modify(context.Background(), []string{"a=a", "b=b"})
			So(err, ShouldNotBeNil)
			So(err.Error(), ShouldEqual, `unable to decode response: invalid character 'o' looking for beginning of value`)
			So(out, ShouldBeNil)
		})
	})
}
