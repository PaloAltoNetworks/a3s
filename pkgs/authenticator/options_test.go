package authenticator

import (
	"testing"

	. "github.com/smartystreets/goconvey/convey"
)

func TestOption(t *testing.T) {

	Convey("OptionIgnoredResources should work", t, func() {
		cfg := &config{}
		OptionIgnoredResources("r1", "r2")(cfg)
		So(cfg.ignoredResources, ShouldResemble, []string{"r1", "r2"})
	})

	Convey("OptionExternalTrustedIssuers should work", t, func() {
		cfg := &config{}
		i1 := RemoteIssuer{URL: "a"}
		i2 := RemoteIssuer{URL: "b"}
		OptionExternalTrustedIssuers(i1, i2)(cfg)
		So(cfg.externalTrustedIssuers, ShouldResemble, []RemoteIssuer{i1, i2})
	})
}
