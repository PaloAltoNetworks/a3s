package token

import (
	"testing"

	. "github.com/smartystreets/goconvey/convey"
)

func TestOptions(t *testing.T) {

	Convey("PrintOptionRaw should work", t, func() {
		cfg := printCfg{}
		PrintOptionRaw(true)(&cfg)
		So(cfg.raw, ShouldBeTrue)
	})

	Convey("PrintOptionDecoded should work", t, func() {
		cfg := printCfg{}
		PrintOptionDecoded(true)(&cfg)
		So(cfg.decoded, ShouldBeTrue)
	})

	Convey("PrintOptionQRCode should work", t, func() {
		cfg := printCfg{}
		PrintOptionQRCode(true)(&cfg)
		So(cfg.qrcode, ShouldBeTrue)
	})
}
