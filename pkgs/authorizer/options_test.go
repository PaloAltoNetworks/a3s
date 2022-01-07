package authorizer

import (
	"testing"

	. "github.com/smartystreets/goconvey/convey"
	"go.aporeto.io/a3s/pkgs/permissions"
)

func TestOption(t *testing.T) {

	Convey("OptionIgnoredResources should work", t, func() {
		cfg := &config{}
		OptionIgnoredResources("r1", "r2")(cfg)
		So(cfg.ignoredResources, ShouldResemble, []string{"r1", "r2"})
	})
}

func TestOptionCheck(t *testing.T) {

	Convey("OptionCheckSourceIP should work", t, func() {
		cfg := &checkConfig{}
		OptionCheckSourceIP("1.1.1.1")(cfg)
		So(cfg.sourceIP, ShouldEqual, "1.1.1.1")
	})

	Convey("OptionCheckID should work", t, func() {
		cfg := &checkConfig{}
		OptionCheckID("id")(cfg)
		So(cfg.id, ShouldEqual, "id")
	})

	Convey("OptionCheckRestrictions should work", t, func() {
		cfg := &checkConfig{}
		r := permissions.Restrictions{Namespace: "/a"}
		OptionCheckRestrictions(r)(cfg)
		So(cfg.restrictions, ShouldResemble, r)
	})
}
