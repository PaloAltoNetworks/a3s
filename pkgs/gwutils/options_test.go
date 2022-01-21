package gwutils

import (
	"testing"
	"time"

	. "github.com/smartystreets/goconvey/convey"
)

func TestOptions(t *testing.T) {

	Convey("OptionCacheDuration should work", t, func() {
		cfg := newVerifierConf()
		OptionCacheDuration(time.Second)(&cfg)
		So(cfg.cacheDuration, ShouldEqual, time.Second)
	})

	Convey("OptionCacheSize should work", t, func() {
		cfg := newVerifierConf()
		OptionCacheSize(42)(&cfg)
		So(cfg.cacheMaxSize, ShouldEqual, 42)
	})

	Convey("OptionTimeout should work", t, func() {
		cfg := newVerifierConf()
		OptionTimeout(time.Second)(&cfg)
		So(cfg.timeout, ShouldEqual, time.Second)
	})
}
