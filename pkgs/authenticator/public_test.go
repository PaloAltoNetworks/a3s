package authenticator

import (
	"context"
	"testing"

	. "github.com/smartystreets/goconvey/convey"
	"go.aporeto.io/bahamut"
	"go.aporeto.io/elemental"
)

func TestNewPublicAuthenticator(t *testing.T) {

	Convey("Calling NewPublicAuthenticator with public resources should work", t, func() {

		a := NewPublic("r1", "r2")
		So(a, ShouldNotBeNil)
		So(len(a.publicResources), ShouldEqual, 2)
		So(a.publicResources, ShouldContainKey, "r1")
		So(a.publicResources, ShouldContainKey, "r2")
	})

	Convey("Calling NewPublicAuthenticator without public resources should work", t, func() {

		a := NewPublic()
		So(a, ShouldNotBeNil)
		So(len(a.publicResources), ShouldEqual, 0)
	})
}

func TestPublicAuthenticateSession(t *testing.T) {

	Convey("Given I have a Public Authenticator ", t, func() {

		a := NewPublic("r1", "r2")

		Convey("Calling AuthenticateSession should always work", func() {

			session := bahamut.NewMockSession()
			action, err := a.AuthenticateSession(session)
			So(err, ShouldBeNil)
			So(action, ShouldEqual, bahamut.AuthActionContinue)
		})
	})
}

func TestPublicAuthenticateRequest(t *testing.T) {

	Convey("Given I have a Public Authenticator ", t, func() {

		a := NewPublic("r1", "r2")

		Convey("Calling AuthenticateSession on public resource should work", func() {

			bctx := bahamut.NewMockContext(context.Background())
			bctx.MockRequest = &elemental.Request{
				Identity: elemental.MakeIdentity("r1", "r1"),
			}
			action, err := a.AuthenticateRequest(bctx)
			So(err, ShouldBeNil)
			So(action, ShouldEqual, bahamut.AuthActionOK)
		})

		Convey("Calling AuthenticateSession on private resource should work", func() {

			bctx := bahamut.NewMockContext(context.Background())
			bctx.MockRequest = &elemental.Request{
				Identity: elemental.MakeIdentity("r11", "r11"),
			}
			action, err := a.AuthenticateRequest(bctx)
			So(err, ShouldBeNil)
			So(action, ShouldEqual, bahamut.AuthActionContinue)
		})
	})
}
