package token

import (
	"net/http"
	"testing"

	. "github.com/smartystreets/goconvey/convey"
	"go.aporeto.io/a3s/pkgs/api"
	"go.aporeto.io/bahamut"
	"go.aporeto.io/elemental"
)

func TestFromRequest(t *testing.T) {

	Convey("Calling FromRequest with a bearer token", t, func() {
		httpreq, _ := http.NewRequest(http.MethodGet, "https://toto.com/namespaces", nil)
		httpreq.Header.Set("Authorization", "Bearer token")
		req, _ := elemental.NewRequestFromHTTPRequest(httpreq, api.Manager())

		t := FromRequest(req)
		So(t, ShouldEqual, "token")
	})

	Convey("Calling FromRequest with a cookie token", t, func() {
		httpreq, _ := http.NewRequest(http.MethodGet, "https://toto.com/namespaces", nil)
		httpreq.AddCookie(&http.Cookie{Name: "x-a3s-token", Value: "token"})
		req, _ := elemental.NewRequestFromHTTPRequest(httpreq, api.Manager())
		t := FromRequest(req)
		So(t, ShouldEqual, "token")
	})

	Convey("Calling FromRequest with both bearer and  cookie token", t, func() {
		httpreq, _ := http.NewRequest(http.MethodGet, "https://toto.com/namespaces", nil)
		httpreq.Header.Set("Authorization", "Bearer token1")
		httpreq.AddCookie(&http.Cookie{Name: "x-a3s-token", Value: "token2"})
		req, _ := elemental.NewRequestFromHTTPRequest(httpreq, api.Manager())
		t := FromRequest(req)
		So(t, ShouldEqual, "token2")
	})
}

func TestFromSession(t *testing.T) {

	Convey("Calling FromRequest with a bearer token", t, func() {
		s := bahamut.NewMockSession()
		s.MockToken = "token"
		t := FromSession(s)
		So(t, ShouldEqual, "token")
	})

	Convey("Calling FromRequest with a cookie token", t, func() {
		s := bahamut.NewMockSession()
		s.MockCookies = map[string]*http.Cookie{"x-a3s-token": {Name: "x-a3s-token", Value: "token"}}
		t := FromSession(s)
		So(t, ShouldEqual, "token")
	})

	Convey("Calling FromRequest with both bearer and  cookie token", t, func() {
		s := bahamut.NewMockSession()
		s.MockCookies = map[string]*http.Cookie{"x-a3s-token": {Name: "x-a3s-token", Value: "token2"}}
		s.MockToken = "token"
		t := FromSession(s)
		So(t, ShouldEqual, "token2")
	})
}
