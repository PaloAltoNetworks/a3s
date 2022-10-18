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

	Convey("OptionRetrieverRestrictions should work", t, func() {
		cfg := &config{}
		r := Restrictions{Namespace: "/a"}
		OptionRetrieverRestrictions(r)(cfg)
		So(cfg.restrictions, ShouldResemble, r)
	})

	Convey("OptionRetrieverTransformer should work", t, func() {
		cfg := &config{}
		t := NewTransformer(
			map[string][]string{
				"r1": {"something:get,post"},
				"r2": {"else:put"},
			},
		)
		OptionRetrieverTransformer(t)(cfg)
		So(cfg.transformer, ShouldResemble, t)
	})
}
