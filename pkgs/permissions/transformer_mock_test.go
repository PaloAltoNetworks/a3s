package permissions

import (
	"testing"

	. "github.com/smartystreets/goconvey/convey"
)

func TestMockTransformer(t *testing.T) {

	Convey("Given a MockTransformer and permissions map", t, func() {

		mockTransformer := NewMockTransformer()

		permissionMap := PermissionMap{
			"r1": {"get": true, "post": true},
		}

		Convey("Calling Transform without mock should work", func() {
			perms := mockTransformer.Transform(permissionMap)
			So(perms, ShouldNotBeNil)
			So(len(perms), ShouldEqual, 0)
		})

		Convey("Calling Transform with mock should work", func() {
			mockTransformer.MockTransform(t, func(PermissionMap) PermissionMap {
				return PermissionMap{"r1": {"get": true, "post": true}}
			})
			perms := mockTransformer.Transform(permissionMap)
			So(perms, ShouldNotBeNil)
			So(len(perms), ShouldEqual, 1)
		})
	})
}
