package authlib

import (
	"testing"
	"time"

	. "github.com/smartystreets/goconvey/convey"
)

func TestBahamut_Options(t *testing.T) {

	c := config{}

	Convey("Calling OptValidity should work", t, func() {
		OptValidity(2 * time.Hour)(&c)
		So(c.validity, ShouldResemble, 2*time.Hour)
	})

	Convey("Calling OptCloak should work", t, func() {
		OptCloak("c1", "c2")(&c)
		So(c.cloak, ShouldResemble, []string{"c1", "c2"})
	})
	Convey("Calling OptOpaque should work", t, func() {
		OptOpaque(map[string]string{"a": "b"})(&c)
		So(c.opaque, ShouldResemble, map[string]string{"a": "b"})
	})

	Convey("Calling OptAudience should work", t, func() {
		OptAudience("audience")(&c)
		So(c.audience, ShouldResemble, []string{"audience"})
	})

	Convey("Calling OptRestrictNamespace should work", t, func() {
		OptRestrictNamespace("/ns")(&c)
		So(c.restrictedNamespace, ShouldEqual, "/ns")
	})

	Convey("Calling OptRestrictPermissions should work", t, func() {
		OptRestrictPermissions([]string{"@auth:role=toto", "test,get,post,put"})(&c)
		So(c.restrictedPermissions, ShouldResemble, []string{"@auth:role=toto", "test,get,post,put"})
	})

	Convey("Calling OptRestrictNetworks should work", t, func() {
		OptRestrictNetworks([]string{"1.0.0.0/8", "2.0.0.0/8"})(&c)
		So(c.restrictedNetworks, ShouldResemble, []string{"1.0.0.0/8", "2.0.0.0/8"})
	})
}
