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
}
