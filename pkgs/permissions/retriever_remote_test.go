package permissions

import (
	"context"
	"fmt"
	"testing"

	. "github.com/smartystreets/goconvey/convey"
	"go.aporeto.io/a3s/pkgs/api"
	"go.aporeto.io/elemental"
	"go.aporeto.io/manipulate"
	"go.aporeto.io/manipulate/maniptest"
)

func TestNewRemoteRetriever(t *testing.T) {
	Convey("Calling NewRemoteRetriever should work", t, func() {
		m := maniptest.NewTestManipulator()
		r := NewRemoteRetriever(m)
		So(r.(*remoteRetriever).manipulator, ShouldEqual, m)
		So(r.(*remoteRetriever).transformer, ShouldEqual, nil)
	})
}

func TestNewRemoteRetrieverWithTransformer(t *testing.T) {
	Convey("Calling NewRemoteRetrieverWithTransformer should work", t, func() {
		m := maniptest.NewTestManipulator()
		mockTransformer := NewMockTransformer()
		r := NewRemoteRetrieverWithTransformer(m, mockTransformer)
		So(r.(*remoteRetriever).manipulator, ShouldEqual, m)
		So(r.(*remoteRetriever).transformer, ShouldEqual, mockTransformer)
	})
}

func TestPermissions(t *testing.T) {

	Convey("Given a remote permissions retriever", t, func() {

		m := maniptest.NewTestManipulator()
		r := NewRemoteRetriever(m)

		Convey("When retrieving subscriptions is OK", func() {

			var expectedClaims []string
			var expectedNamespace string
			var expectedRestrictions Restrictions
			var expectedOffloadPermissionsRestrictions bool
			var expectedID string
			var expectedIP string
			m.MockCreate(t, func(mctx manipulate.Context, object elemental.Identifiable) error {
				o := object.(*api.Permissions)
				o.Permissions = map[string]map[string]bool{
					"cat": {"pet": false},
					"dog": {"pet": true},
				}
				expectedClaims = o.Claims
				expectedNamespace = o.Namespace
				expectedID = o.ID
				expectedIP = o.IP
				expectedOffloadPermissionsRestrictions = o.OffloadPermissionsRestrictions
				expectedRestrictions = Restrictions{
					Namespace:   o.RestrictedNamespace,
					Permissions: o.RestrictedPermissions,
					Networks:    o.RestrictedNetworks,
				}

				return nil
			})

			perms, err := r.Permissions(
				context.Background(),
				[]string{"a=a"},
				"/the/ns",
				OptionRetrieverID("id"),
				OptionRetrieverSourceIP("1.1.1.1"),
				OptionRetrieverRestrictions(Restrictions{
					Namespace:   "/the/ns/sub",
					Networks:    []string{"1.1.1.1/32", "2.2.2.2/32"},
					Permissions: []string{"cat:pet"},
				}),
			)

			So(err, ShouldBeNil)
			So(perms, ShouldResemble, PermissionMap{
				"cat": Permissions{"pet": false},
				"dog": Permissions{"pet": true},
			})
			So(expectedClaims, ShouldResemble, []string{"a=a"})
			So(expectedNamespace, ShouldResemble, "/the/ns")
			So(expectedID, ShouldEqual, "id")
			So(expectedIP, ShouldEqual, "1.1.1.1")
			So(expectedOffloadPermissionsRestrictions, ShouldBeFalse)
			So(expectedRestrictions, ShouldResemble, Restrictions{
				Namespace:   "/the/ns/sub",
				Networks:    []string{"1.1.1.1/32", "2.2.2.2/32"},
				Permissions: []string{"cat:pet"},
			})
		})

		Convey("When retrieving permissions fails", func() {

			m.MockCreate(t, func(mctx manipulate.Context, object elemental.Identifiable) error {
				return fmt.Errorf("boom")
			})

			_, err := r.Permissions(
				context.Background(),
				[]string{"a=a"},
				"/the/ns",
				OptionRetrieverID("id"),
			)

			So(err, ShouldNotBeNil)
			So(err.Error(), ShouldEqual, "boom")
		})

		Convey("When retrieving subscriptions with a defined transformer", func() {

			mockTransformer := NewMockTransformer()
			mockTransformer.MockTransform(t, func(permissions PermissionMap) PermissionMap {
				return PermissionMap{
					"cat": Permissions{
						"pet":  false,
						"feed": true,
					},
					"dog": Permissions{
						"pet":  true,
						"feed": true,
					},
				}
			})

			r = NewRemoteRetrieverWithTransformer(m, mockTransformer)

			var expectedClaims []string
			var expectedNamespace string
			var expectedRestrictions Restrictions
			var expectedOffloadPermissionsRestrictions bool
			var expectedID string
			var expectedIP string
			m.MockCreate(t, func(mctx manipulate.Context, object elemental.Identifiable) error {
				o := object.(*api.Permissions)
				o.Permissions = map[string]map[string]bool{
					"cat": {"pet": false},
					"dog": {"pet": true},
				}
				expectedClaims = o.Claims
				expectedNamespace = o.Namespace
				expectedID = o.ID
				expectedIP = o.IP
				expectedOffloadPermissionsRestrictions = o.OffloadPermissionsRestrictions
				expectedRestrictions = Restrictions{
					Namespace:   o.RestrictedNamespace,
					Permissions: o.RestrictedPermissions,
					Networks:    o.RestrictedNetworks,
				}

				return nil
			})

			perms, err := r.Permissions(
				context.Background(),
				[]string{"a=a"},
				"/the/ns",
				OptionRetrieverID("id"),
				OptionRetrieverSourceIP("1.1.1.1"),
				OptionRetrieverRestrictions(Restrictions{
					Namespace:   "/the/ns/sub",
					Networks:    []string{"1.1.1.1/32", "2.2.2.2/32"},
					Permissions: []string{"cat:pet"},
				}),
			)

			So(err, ShouldBeNil)
			So(perms, ShouldResemble, PermissionMap{
				"cat": Permissions{
					"feed": true,
				},
				"dog": Permissions{
					"pet":  true,
					"feed": true,
				},
			})
			So(expectedClaims, ShouldResemble, []string{"a=a"})
			So(expectedNamespace, ShouldResemble, "/the/ns")
			So(expectedID, ShouldEqual, "id")
			So(expectedIP, ShouldEqual, "1.1.1.1")
			So(expectedOffloadPermissionsRestrictions, ShouldBeTrue)
			So(expectedRestrictions, ShouldResemble, Restrictions{
				Namespace:   "/the/ns/sub",
				Networks:    []string{"1.1.1.1/32", "2.2.2.2/32"},
				Permissions: []string{"cat:pet"},
			})
		})
	})
}
