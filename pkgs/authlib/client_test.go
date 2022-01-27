package authlib

import (
	"context"
	"fmt"
	"testing"
	"time"

	. "github.com/smartystreets/goconvey/convey"
	"go.aporeto.io/a3s/pkgs/api"
	"go.aporeto.io/elemental"
	"go.aporeto.io/manipulate"
	"go.aporeto.io/manipulate/maniptest"
)

func TestClient_NewClient(t *testing.T) {

	Convey("Given I create a new Client with a valid URL", t, func() {

		m := maniptest.NewTestManipulator()
		cl := NewClient(m)

		Convey("Then client should be correctly initialized", func() {
			So(cl, ShouldNotBeNil)
		})

		Convey("Then client url should be set", func() {
			So(cl.manipulator, ShouldEqual, m)
		})
	})
}

func TestAuthFromCertificate(t *testing.T) {

	Convey("The function should work", t, func() {

		expectedRequest := api.NewIssue()

		m := maniptest.NewTestManipulator()
		m.MockCreate(t, func(mctx manipulate.Context, object elemental.Identifiable) error {
			expectedRequest = object.(*api.Issue)
			expectedRequest.Token = "yeay!"
			return nil
		})

		cl := NewClient(m)

		token, err := cl.AuthFromCertificate(
			context.Background(),
			"/ns",
			"name",
			OptAudience("aud"),
			OptRefresh(true),
		)

		So(err, ShouldBeNil)
		So(expectedRequest.SourceType, ShouldEqual, api.IssueSourceTypeMTLS)
		So(expectedRequest.SourceNamespace, ShouldEqual, "/ns")
		So(expectedRequest.SourceName, ShouldEqual, "name")
		So(expectedRequest.Audience, ShouldResemble, []string{"aud"})
		So(expectedRequest.Validity, ShouldEqual, time.Hour.String())
		So(expectedRequest.TokenType, ShouldEqual, api.IssueTokenTypeRefresh)
		So(token, ShouldEqual, "yeay!")
	})
}

func TestAuthFromLDAP(t *testing.T) {

	Convey("The function should work", t, func() {

		expectedRequest := api.NewIssue()

		m := maniptest.NewTestManipulator()
		m.MockCreate(t, func(mctx manipulate.Context, object elemental.Identifiable) error {
			expectedRequest = object.(*api.Issue)
			expectedRequest.Token = "yeay!"
			return nil
		})

		cl := NewClient(m)

		token, err := cl.AuthFromLDAP(
			context.Background(),
			"user",
			"pass",
			"/ns",
			"name",
			OptAudience("aud"),
		)

		So(err, ShouldBeNil)
		So(expectedRequest.SourceType, ShouldEqual, api.IssueSourceTypeLDAP)
		So(expectedRequest.SourceNamespace, ShouldEqual, "/ns")
		So(expectedRequest.SourceName, ShouldEqual, "name")
		So(expectedRequest.Audience, ShouldResemble, []string{"aud"})
		So(expectedRequest.Validity, ShouldEqual, time.Hour.String())
		So(expectedRequest.InputLDAP.Username, ShouldEqual, "user")
		So(expectedRequest.InputLDAP.Password, ShouldEqual, "pass")
		So(token, ShouldEqual, "yeay!")
	})
}

func TestAuthFromA3S(t *testing.T) {

	Convey("The function should work", t, func() {

		expectedRequest := api.NewIssue()

		m := maniptest.NewTestManipulator()
		m.MockCreate(t, func(mctx manipulate.Context, object elemental.Identifiable) error {
			expectedRequest = object.(*api.Issue)
			expectedRequest.Token = "yeay!"
			return nil
		})

		cl := NewClient(m)

		token, err := cl.AuthFromA3S(
			context.Background(),
			"token",
			OptAudience("aud"),
		)

		So(err, ShouldBeNil)
		So(expectedRequest.SourceType, ShouldEqual, api.IssueSourceTypeA3S)
		So(expectedRequest.SourceNamespace, ShouldEqual, "")
		So(expectedRequest.SourceName, ShouldEqual, "")
		So(expectedRequest.Audience, ShouldResemble, []string{"aud"})
		So(expectedRequest.Validity, ShouldEqual, time.Hour.String())
		So(expectedRequest.InputA3S.Token, ShouldEqual, "token")
		So(token, ShouldEqual, "yeay!")
	})
}

func TestAuthFromRemoteA3S(t *testing.T) {

	Convey("The function should work", t, func() {

		expectedRequest := api.NewIssue()

		m := maniptest.NewTestManipulator()
		m.MockCreate(t, func(mctx manipulate.Context, object elemental.Identifiable) error {
			expectedRequest = object.(*api.Issue)
			expectedRequest.Token = "yeay!"
			return nil
		})

		cl := NewClient(m)

		token, err := cl.AuthFromRemoteA3S(
			context.Background(),
			"token",
			"/ns",
			"name",
			OptAudience("aud"),
		)

		So(err, ShouldBeNil)
		So(expectedRequest.SourceType, ShouldEqual, api.IssueSourceTypeRemoteA3S)
		So(expectedRequest.SourceNamespace, ShouldEqual, "/ns")
		So(expectedRequest.SourceName, ShouldEqual, "name")
		So(expectedRequest.Audience, ShouldResemble, []string{"aud"})
		So(expectedRequest.Validity, ShouldEqual, time.Hour.String())
		So(expectedRequest.InputRemoteA3S.Token, ShouldEqual, "token")
		So(token, ShouldEqual, "yeay!")
	})
}

func TestAuthFromAWS(t *testing.T) {

	Convey("The function should work", t, func() {

		expectedRequest := api.NewIssue()

		m := maniptest.NewTestManipulator()
		m.MockCreate(t, func(mctx manipulate.Context, object elemental.Identifiable) error {
			expectedRequest = object.(*api.Issue)
			expectedRequest.Token = "yeay!"
			return nil
		})

		cl := NewClient(m)

		token, err := cl.AuthFromAWS(
			context.Background(),
			"aid",
			"sid",
			"token",
			OptAudience("aud"),
		)

		So(err, ShouldBeNil)
		So(expectedRequest.SourceType, ShouldEqual, api.IssueSourceTypeAWS)
		So(expectedRequest.SourceNamespace, ShouldEqual, "")
		So(expectedRequest.SourceName, ShouldEqual, "")
		So(expectedRequest.Audience, ShouldResemble, []string{"aud"})
		So(expectedRequest.Validity, ShouldEqual, time.Hour.String())
		So(expectedRequest.InputAWS.ID, ShouldEqual, "aid")
		So(expectedRequest.InputAWS.Secret, ShouldEqual, "sid")
		So(expectedRequest.InputAWS.Token, ShouldEqual, "token")
		So(token, ShouldEqual, "yeay!")
	})
}

func TestAuthFromGCP(t *testing.T) {

	Convey("The function should work", t, func() {

		expectedRequest := api.NewIssue()

		m := maniptest.NewTestManipulator()
		m.MockCreate(t, func(mctx manipulate.Context, object elemental.Identifiable) error {
			expectedRequest = object.(*api.Issue)
			expectedRequest.Token = "yeay!"
			return nil
		})

		cl := NewClient(m)

		token, err := cl.AuthFromGCP(
			context.Background(),
			"token",
			"audience",
			OptAudience("aud"),
		)

		So(err, ShouldBeNil)
		So(expectedRequest.SourceType, ShouldEqual, api.IssueSourceTypeGCP)
		So(expectedRequest.SourceNamespace, ShouldEqual, "")
		So(expectedRequest.SourceName, ShouldEqual, "")
		So(expectedRequest.Audience, ShouldResemble, []string{"aud"})
		So(expectedRequest.Validity, ShouldEqual, time.Hour.String())
		So(expectedRequest.InputGCP.Token, ShouldEqual, "token")
		So(expectedRequest.InputGCP.Audience, ShouldEqual, "audience")
		So(token, ShouldEqual, "yeay!")
	})
}

func TestAuthFromAzure(t *testing.T) {

	Convey("The function should work", t, func() {

		expectedRequest := api.NewIssue()

		m := maniptest.NewTestManipulator()
		m.MockCreate(t, func(mctx manipulate.Context, object elemental.Identifiable) error {
			expectedRequest = object.(*api.Issue)
			expectedRequest.Token = "yeay!"
			return nil
		})

		cl := NewClient(m)

		token, err := cl.AuthFromAzure(
			context.Background(),
			"token",
			OptAudience("aud"),
		)

		So(err, ShouldBeNil)
		So(expectedRequest.SourceType, ShouldEqual, api.IssueSourceTypeAzure)
		So(expectedRequest.SourceNamespace, ShouldEqual, "")
		So(expectedRequest.SourceName, ShouldEqual, "")
		So(expectedRequest.Audience, ShouldResemble, []string{"aud"})
		So(expectedRequest.Validity, ShouldEqual, time.Hour.String())
		So(expectedRequest.InputAzure.Token, ShouldEqual, "token")
		So(token, ShouldEqual, "yeay!")
	})
}

func TestAuthFromOIDCStep1(t *testing.T) {

	Convey("The function should work", t, func() {

		expectedRequest := api.NewIssue()

		m := maniptest.NewTestManipulator()
		m.MockCreate(t, func(mctx manipulate.Context, object elemental.Identifiable) error {
			expectedRequest = object.(*api.Issue)
			expectedRequest.InputOIDC.AuthURL = "authurl"
			return nil
		})

		cl := NewClient(m)

		url, err := cl.AuthFromOIDCStep1(
			context.Background(),
			"/ns",
			"name",
			"url",
		)

		So(err, ShouldBeNil)
		So(expectedRequest.SourceType, ShouldEqual, api.IssueSourceTypeOIDC)
		So(expectedRequest.SourceNamespace, ShouldEqual, "/ns")
		So(expectedRequest.SourceName, ShouldEqual, "name")
		So(expectedRequest.InputOIDC.RedirectURL, ShouldEqual, "url")
		So(expectedRequest.InputOIDC.NoAuthRedirect, ShouldEqual, true)
		So(url, ShouldEqual, "authurl")
	})
}

func TestAuthFromOIDCStep2(t *testing.T) {

	Convey("The function should work", t, func() {

		expectedRequest := api.NewIssue()

		m := maniptest.NewTestManipulator()
		m.MockCreate(t, func(mctx manipulate.Context, object elemental.Identifiable) error {
			expectedRequest = object.(*api.Issue)
			expectedRequest.Token = "yeay!"
			return nil
		})

		cl := NewClient(m)

		token, err := cl.AuthFromOIDCStep2(
			context.Background(),
			"code",
			"state",
			OptAudience("aud"),
		)

		So(err, ShouldBeNil)
		So(expectedRequest.SourceType, ShouldEqual, api.IssueSourceTypeOIDC)
		So(expectedRequest.SourceNamespace, ShouldEqual, "")
		So(expectedRequest.SourceName, ShouldEqual, "")
		So(expectedRequest.Audience, ShouldResemble, []string{"aud"})
		So(expectedRequest.Validity, ShouldEqual, time.Hour.String())
		So(expectedRequest.InputOIDC.Code, ShouldEqual, "code")
		So(expectedRequest.InputOIDC.State, ShouldEqual, "state")
		So(token, ShouldEqual, "yeay!")
	})
}

func TestAuthFromHTTP(t *testing.T) {

	Convey("The function should work", t, func() {

		expectedRequest := api.NewIssue()

		m := maniptest.NewTestManipulator()
		m.MockCreate(t, func(mctx manipulate.Context, object elemental.Identifiable) error {
			expectedRequest = object.(*api.Issue)
			expectedRequest.Token = "yeay!"
			return nil
		})

		cl := NewClient(m)

		token, err := cl.AuthFromHTTP(
			context.Background(),
			"user",
			"pass",
			"1234",
			"/ns",
			"name",
			OptAudience("aud"),
		)

		So(err, ShouldBeNil)
		So(expectedRequest.SourceType, ShouldEqual, api.IssueSourceTypeHTTP)
		So(expectedRequest.SourceNamespace, ShouldEqual, "/ns")
		So(expectedRequest.SourceName, ShouldEqual, "name")
		So(expectedRequest.Audience, ShouldResemble, []string{"aud"})
		So(expectedRequest.Validity, ShouldEqual, time.Hour.String())
		So(expectedRequest.InputHTTP.Username, ShouldEqual, "user")
		So(expectedRequest.InputHTTP.Password, ShouldEqual, "pass")
		So(expectedRequest.InputHTTP.TOTP, ShouldEqual, "1234")
		So(token, ShouldEqual, "yeay!")
	})
}

func TestSendRequest(t *testing.T) {

	Convey("Calling sendRequest when everything is ok should work. ", t, func() {

		m := maniptest.NewTestManipulator()
		m.MockCreate(t, func(mctx manipulate.Context, object elemental.Identifiable) error {
			object.(*api.Issue).Token = "yeay!"
			return nil
		})

		cl := NewClient(m)

		t, err := cl.sendRequest(context.Background(), &api.Issue{})
		So(err, ShouldBeNil)
		So(t, ShouldEqual, "yeay!")
	})

	Convey("Calling sendRequest when it fail should work. ", t, func() {

		m := maniptest.NewTestManipulator()
		m.MockCreate(t, func(mctx manipulate.Context, object elemental.Identifiable) error {
			return fmt.Errorf("boom")
		})

		cl := NewClient(m)

		t, err := cl.sendRequest(context.Background(), &api.Issue{})
		So(err, ShouldNotBeNil)
		So(err.Error(), ShouldEqual, "boom")
		So(t, ShouldEqual, "")
	})
}
