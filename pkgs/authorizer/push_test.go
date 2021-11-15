package authorizer

import (
	"context"
	"fmt"
	"testing"

	. "github.com/smartystreets/goconvey/convey"
	"go.aporeto.io/a3s/pkgs/api"
	"go.aporeto.io/a3s/pkgs/permissions"
	"go.aporeto.io/a3s/pkgs/token"
	"go.aporeto.io/bahamut"
	"go.aporeto.io/elemental"
	"go.aporeto.io/manipulate/maniptest"
)

func TestNewPushDispatcher(t *testing.T) {

	Convey("Given I call NewPushDispatchHandler", t, func() {

		m := maniptest.NewTestManipulator()
		p := bahamut.NewLocalPubSubClient()
		_ = p.Connect(context.Background())
		r := permissions.NewRetriever(m)
		a := New(context.Background(), r, p)
		h := NewPushDispatchHandler(m, a)

		So(h.manipulator, ShouldEqual, m)
		So(func() { h.OnPushSessionStart(bahamut.NewMockSession()) }, ShouldNotPanic)
		So(func() { h.OnPushSessionStop(bahamut.NewMockSession()) }, ShouldNotPanic)
	})
}

func TestOnPushSessionInit(t *testing.T) {

	Convey("Given I have a push handler", t, func() {

		_, key := getECCert()

		m := maniptest.NewTestManipulator()
		p := bahamut.NewLocalPubSubClient()
		_ = p.Connect(context.Background())
		r := permissions.NewMockRetriever()
		a := New(context.Background(), r, p)
		h := NewPushDispatchHandler(m, a)

		Convey("When there are bad restrictions in the token ", func() {

			s := bahamut.NewMockSession()
			s.MockToken = "that's no token"

			ok, err := h.OnPushSessionInit(s)
			So(ok, ShouldBeFalse)
			So(err, ShouldNotBeNil)
			So(err.Error(), ShouldEqual, "unable to compute authz restrictions from token: token contains an invalid number of segments")
		})

		Convey("When authorizer errors", func() {

			r.MockPermissions(t, func(context.Context, []string, string, ...permissions.RetrieverOption) (permissions.PermissionMap, error) {
				return nil, fmt.Errorf("boom")
			})

			s := bahamut.NewMockSession()
			s.MockToken = makeToken(&token.IdentityToken{
				Source: token.Source{
					Type: "test",
				},
			}, key)
			s.MockParameters = map[string]string{"namespace": "/"}

			ok, err := h.OnPushSessionInit(s)
			So(ok, ShouldBeFalse)
			So(err, ShouldNotBeNil)
			So(err.Error(), ShouldEqual, "boom")
		})

		Convey("When authorizer is ok", func() {

			r.MockPermissions(t, func(context.Context, []string, string, ...permissions.RetrieverOption) (permissions.PermissionMap, error) {
				return permissions.PermissionMap{"pushsession": permissions.Permissions{"get": true}}, nil
			})

			s := bahamut.NewMockSession()
			s.MockToken = makeToken(&token.IdentityToken{
				Source: token.Source{
					Type: "test",
				},
			}, key)
			s.MockParameters = map[string]string{"namespace": "/"}

			ok, err := h.OnPushSessionInit(s)
			So(ok, ShouldBeTrue)
			So(err, ShouldBeNil)
		})

		Convey("When authorizer is not ok", func() {

			r.MockPermissions(t, func(context.Context, []string, string, ...permissions.RetrieverOption) (permissions.PermissionMap, error) {
				return permissions.PermissionMap{"not-pushsession": permissions.Permissions{"get": true}}, nil
			})

			s := bahamut.NewMockSession()
			s.MockToken = makeToken(&token.IdentityToken{
				Source: token.Source{
					Type: "test",
				},
			}, key)
			s.MockParameters = map[string]string{"namespace": "/"}

			ok, err := h.OnPushSessionInit(s)
			So(ok, ShouldBeFalse)
			So(err, ShouldBeNil)
		})
	})
}

func TestSummarizeEvent(t *testing.T) {

	Convey("Given I have a push handler", t, func() {

		m := maniptest.NewTestManipulator()
		p := bahamut.NewLocalPubSubClient()
		_ = p.Connect(context.Background())
		r := permissions.NewMockRetriever()
		a := New(context.Background(), r, p)
		h := NewPushDispatchHandler(m, a)

		Convey("Calling SummarizeEvent with an valid event should work", func() {

			evt := elemental.NewEvent(
				elemental.EventCreate,
				&api.Namespace{
					Namespace: "/the/ns",
					Name:      "/the/ns/gros",
					ID:        "x",
				},
			)
			out, err := h.SummarizeEvent(evt)
			So(err, ShouldBeNil)
			So(out, ShouldResemble, pushedEntity{
				Namespace: "/the/ns",
				Name:      "/the/ns/gros",
				ID:        "x",
			})
		})
	})
}

func TestShouldDispatch(t *testing.T) {

	Convey("Given I have a push handler", t, func() {

		_, key := getECCert()

		m := maniptest.NewTestManipulator()
		p := bahamut.NewLocalPubSubClient()
		_ = p.Connect(context.Background())
		r := permissions.NewMockRetriever()
		a := New(context.Background(), r, p)
		h := NewPushDispatchHandler(m, a)

		Convey("When we receive a namespace delete event", func() {

			s := bahamut.NewMockSession()
			s.MockParameters = map[string]string{"namespace": "/test"}

			evt := elemental.NewEvent(elemental.EventDelete, api.NewNamespace())
			ok, err := h.ShouldDispatch(s, evt, pushedEntity{Name: "/test"})
			So(err, ShouldBeNil)
			So(ok, ShouldBeTrue)
		})

		Convey("When we receive a namespace delete event from a ns above", func() {

			s := bahamut.NewMockSession()
			s.MockParameters = map[string]string{"namespace": "/test/yo"}

			evt := elemental.NewEvent(elemental.EventDelete, api.NewNamespace())
			ok, err := h.ShouldDispatch(s, evt, pushedEntity{Name: "/test", Namespace: "/"})
			So(err, ShouldBeNil)
			So(ok, ShouldBeTrue)
		})

		Convey("When we receive an event from a sibling namespace", func() {

			s := bahamut.NewMockSession()
			s.MockParameters = map[string]string{"namespace": "/test"}

			evt := elemental.NewEvent(elemental.EventDelete, api.NewAuthorization())
			ok, err := h.ShouldDispatch(s, evt, pushedEntity{Namespace: "/not-test"})
			So(err, ShouldBeNil)
			So(ok, ShouldBeFalse)
		})

		Convey("When we receive an event from a parent namespace", func() {

			s := bahamut.NewMockSession()
			s.MockParameters = map[string]string{"namespace": "/test"}

			evt := elemental.NewEvent(elemental.EventDelete, api.NewAuthorization())
			ok, err := h.ShouldDispatch(s, evt, pushedEntity{Namespace: "/"})
			So(err, ShouldBeNil)
			So(ok, ShouldBeFalse)
		})

		Convey("When we receive an event from a children namespace without recursive", func() {

			s := bahamut.NewMockSession()
			s.MockParameters = map[string]string{"namespace": "/test"}

			evt := elemental.NewEvent(elemental.EventDelete, api.NewAuthorization())
			ok, err := h.ShouldDispatch(s, evt, pushedEntity{Namespace: "/test/yo"})
			So(err, ShouldBeNil)
			So(ok, ShouldBeFalse)
		})

		Convey("When we receive an event from a parent namespace with permission ok and propagation on", func() {

			r.MockPermissions(t, func(context.Context, []string, string, ...permissions.RetrieverOption) (permissions.PermissionMap, error) {
				return permissions.PermissionMap{"authorization": permissions.Permissions{"get": true}}, nil
			})

			s := bahamut.NewMockSession()
			s.MockToken = makeToken(&token.IdentityToken{
				Source: token.Source{
					Type: "test",
				},
			}, key)
			s.MockParameters = map[string]string{"namespace": "/test"}

			evt := elemental.NewEvent(elemental.EventDelete, api.NewAuthorization())
			ok, err := h.ShouldDispatch(s, evt, pushedEntity{Namespace: "/", Propagate: true})
			So(err, ShouldBeNil)
			So(ok, ShouldBeTrue)
		})

		Convey("When we receive an event from a parent namespace with permission ok and propagation on but hidden true", func() {

			r.MockPermissions(t, func(context.Context, []string, string, ...permissions.RetrieverOption) (permissions.PermissionMap, error) {
				return permissions.PermissionMap{"authorization": permissions.Permissions{"get": true}}, nil
			})

			s := bahamut.NewMockSession()
			s.MockToken = makeToken(&token.IdentityToken{
				Source: token.Source{
					Type: "test",
				},
			}, key)
			s.MockParameters = map[string]string{"namespace": "/test"}

			evt := elemental.NewEvent(elemental.EventDelete, api.NewAuthorization())
			ok, err := h.ShouldDispatch(s, evt, pushedEntity{Namespace: "/", Propagate: true, PropagationHidden: true})
			So(err, ShouldBeNil)
			So(ok, ShouldBeFalse)
		})

		Convey("When we receive an event from the current namespace with permission ko", func() {

			r.MockPermissions(t, func(context.Context, []string, string, ...permissions.RetrieverOption) (permissions.PermissionMap, error) {
				return permissions.PermissionMap{"authorization": permissions.Permissions{"delete": true}}, nil
			})

			s := bahamut.NewMockSession()
			s.MockToken = makeToken(&token.IdentityToken{
				Source: token.Source{
					Type: "test",
				},
			}, key)
			s.MockParameters = map[string]string{"namespace": "/test"}

			evt := elemental.NewEvent(elemental.EventDelete, api.NewAuthorization())
			ok, err := h.ShouldDispatch(s, evt, pushedEntity{Namespace: "/test"})
			So(err, ShouldBeNil)
			So(ok, ShouldBeFalse)
		})

		Convey("When we receive an event from the current namespace with permission ok", func() {

			r.MockPermissions(t, func(context.Context, []string, string, ...permissions.RetrieverOption) (permissions.PermissionMap, error) {
				return permissions.PermissionMap{"authorization": permissions.Permissions{"get": true}}, nil
			})

			s := bahamut.NewMockSession()
			s.MockToken = makeToken(&token.IdentityToken{
				Source: token.Source{
					Type: "test",
				},
			}, key)
			s.MockParameters = map[string]string{"namespace": "/test"}

			evt := elemental.NewEvent(elemental.EventDelete, api.NewAuthorization())
			ok, err := h.ShouldDispatch(s, evt, pushedEntity{Namespace: "/test"})
			So(err, ShouldBeNil)
			So(ok, ShouldBeTrue)
		})

		Convey("When we receive an event from a child namespace with permission ok and recursive", func() {

			r.MockPermissions(t, func(context.Context, []string, string, ...permissions.RetrieverOption) (permissions.PermissionMap, error) {
				return permissions.PermissionMap{"authorization": permissions.Permissions{"get": true}}, nil
			})

			s := bahamut.NewMockSession()
			s.MockToken = makeToken(&token.IdentityToken{
				Source: token.Source{
					Type: "test",
				},
			}, key)
			s.MockParameters = map[string]string{"namespace": "/test", "mode": "all"}

			evt := elemental.NewEvent(elemental.EventDelete, api.NewAuthorization())
			ok, err := h.ShouldDispatch(s, evt, pushedEntity{Namespace: "/test/yo"})
			So(err, ShouldBeNil)
			So(ok, ShouldBeTrue)
		})

		Convey("When the session has an invalid token", func() {

			r.MockPermissions(t, func(context.Context, []string, string, ...permissions.RetrieverOption) (permissions.PermissionMap, error) {
				return permissions.PermissionMap{"authorization": permissions.Permissions{"get": true}}, nil
			})

			s := bahamut.NewMockSession()
			s.MockToken = "not a token"
			s.MockParameters = map[string]string{"namespace": "/test", "mode": "all"}

			evt := elemental.NewEvent(elemental.EventDelete, api.NewAuthorization())
			ok, err := h.ShouldDispatch(s, evt, pushedEntity{Namespace: "/test/yo"})
			So(err, ShouldNotBeNil)
			So(err.Error(), ShouldEqual, "unable to compute authz restrictions from token: token contains an invalid number of segments")
			So(ok, ShouldBeFalse)
		})
	})
}
