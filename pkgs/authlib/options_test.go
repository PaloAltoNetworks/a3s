package authlib

import (
	"testing"
	"time"

	. "github.com/smartystreets/goconvey/convey"
	"go.aporeto.io/a3s/pkgs/permissions"
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

	Convey("Calling OptRestrictions should work", t, func() {
		r := permissions.Restrictions{
			Namespace:   "/ns",
			Networks:    []string{"10.0.0.1/32"},
			Permissions: []string{"r:a1,a2"},
		}
		OptRestrictions(r)(&c)
		So(c.restrictions, ShouldResemble, r)
	})
}
