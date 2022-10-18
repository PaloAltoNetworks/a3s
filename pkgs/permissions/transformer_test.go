package permissions

import (
	"testing"

	. "github.com/smartystreets/goconvey/convey"
)

func TestNewTransformer(t *testing.T) {
	Convey("Given have a subscriber and a manipulator", t, func() {
		roles := map[string][]string{
			"r1": {"something:get,post"},
		}
		t := NewTransformer(roles).(*transformer)
		So(t.roleExpander, ShouldNotBeNil)
	})
}

func TestTransform(t *testing.T) {

	Convey("Given I have a transformer and permission map", t, func() {

		roles := map[string][]string{
			"r1": {"something:get,post"},
		}

		t := NewTransformer(roles).(*transformer)

		Convey("When I call Transform with permissions that need expanded", func() {

			perms := []string{
				"dog:pet,walk",
				"r1",
			}

			newPermissions := t.Transform(perms)

			So(len(newPermissions), ShouldEqual, 2)
			So(newPermissions, ShouldResemble, []string{
				"dog:pet,walk",
				"something:get,post",
			})
		})

		Convey("When I call Transform with permissions that don't match expander", func() {

			perms := []string{
				"dog:pet,walk",
				"r2",
			}

			newPermissions := t.Transform(perms)

			So(len(newPermissions), ShouldEqual, 2)
			So(newPermissions, ShouldResemble, []string{
				"dog:pet,walk",
				"r2",
			})
		})

		Convey("When I call Transform with no permissions", func() {

			newPermissions := t.Transform(nil)

			So(len(newPermissions), ShouldEqual, 0)
		})

		Convey("When I call Transform with no role expander defined", func() {

			t := NewTransformer(nil).(*transformer)

			perms := []string{
				"dog:pet,walk",
				"r1",
			}

			newPermissions := t.Transform(perms)

			So(len(newPermissions), ShouldEqual, 2)
			So(newPermissions, ShouldResemble, []string{
				"dog:pet,walk",
				"r1",
			})
		})
	})
}
