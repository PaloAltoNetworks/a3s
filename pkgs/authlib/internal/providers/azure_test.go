package providers

import (
	"encoding/json"
	"fmt"
	"net/http"
	"net/http/httptest"
	"testing"

	. "github.com/smartystreets/goconvey/convey" //revive:disable-line:dot-imports
)

func newValidAzureToken() string {
	token := struct {
		AccessToken string `json:"access_token"`
	}{
		AccessToken: "the role",
	}

	data, _ := json.Marshal(token) // nolint errcheck

	return string(data)
}

func Test_AzureServiceIdentityToken(t *testing.T) {

	Convey("When I call AzureServiceIdentityToken with no errors", t, func() {
		ts := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			if r.URL.Path == "/" {
				fmt.Fprintln(w, newValidAzureToken())
			}
		}))
		defer ts.Close()

		azureServiceTokenURL = ts.URL
		token, err := AzureServiceIdentityToken()

		Convey("Then err should be nil", func() {
			So(err, ShouldBeNil)
		})

		Convey("Then the token should be correct", func() {
			So(token, ShouldResemble, "the role")
		})

	})

	Convey("When I call AzureServiceIdentityToken and the token cannot be decoded", t, func() {
		ts := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			if r.URL.Path == "/" {
				fmt.Fprintln(w, `bad data`)
			}
		}))
		defer ts.Close()

		azureServiceTokenURL = ts.URL
		_, err := AzureServiceIdentityToken()

		Convey("Then err should  not be nil", func() {
			So(err, ShouldNotBeNil)
		})

	})

	Convey("When I call AzureServiceIdentityToken without info (calling Azure) but can't retrieve token (comm error)", t, func() {

		ts2 := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			if r.URL.Path == "/latest/meta-data/iam/security-credentials/" {
				fmt.Fprintln(w, `the-role`)
			}
		}))
		defer ts2.Close()

		azureServiceTokenURL = "nope"
		_, err := AzureServiceIdentityToken()

		Convey("Then err should not be nil", func() {
			So(err, ShouldNotBeNil)
		})
	})

}
