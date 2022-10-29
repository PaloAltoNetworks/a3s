package permissions

import (
	"testing"

	. "github.com/smartystreets/goconvey/convey"
)

func TestRetrieverOptions(t *testing.T) {

	Convey("OptionRetrievedID should work", t, func() {
		cfg := &config{}
		OptionRetrieverID("xxx")(cfg)
		So(cfg.id, ShouldEqual, "xxx")
	})

	Convey("OptionRetrievedSourceIP should work", t, func() {
		cfg := &config{}
		OptionRetrieverSourceIP("1.2.3.4")(cfg)
		So(cfg.addr, ShouldEqual, "1.2.3.4")
	})

	Convey("Option should work", t, func() {
		cfg := &config{}
		r := Restrictions{Namespace: "/a"}
		OptionRetrieverRestrictions(r)(cfg)
		So(cfg.restrictions, ShouldResemble, r)
	})

	Convey("OptionOffloadPermissionsRestrictions should work", t, func() {
		cfg := &config{}
		OptionOffloadPermissionsRestrictions(true)(cfg)
		So(cfg.offloadPermissionsRestrictions, ShouldBeTrue)
	})

}
