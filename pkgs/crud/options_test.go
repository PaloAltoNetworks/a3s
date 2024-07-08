package crud

import (
	"fmt"
	"testing"

	. "github.com/smartystreets/goconvey/convey"
)

/*
func TestOptions(t *testing.T) {

	Convey("OptionPreWriteHook should work ", t, func() {
		cfg := cfg{
			preHook: func(elemental.Identifiable, elemental.Identifiable) error { return nil },
		}
		f := func(elemental.Identifiable, elemental.Identifiable) error { return nil }
		OptionPreWriteHook(f)(&cfg)
		So(cfg.preHook, ShouldHaveSameTypeAs, f)
	})

	Convey("OptionPostWriteHook should work ", t, func() {
		cfg := cfg{
			postHook: func(elemental.Identifiable) {},
		}
		f := func(elemental.Identifiable) {}
		OptionPostWriteHook(f)(&cfg)
		So(cfg.postHook, ShouldHaveSameTypeAs, f)
	})
}
*/

func TestErrPreWriteHook(t *testing.T) {

	Convey("Given I have an ErrPreWriteHook", t, func() {
		err := ErrPreWriteHook{Err: fmt.Errorf("boom")}
		So(err.Error(), ShouldEqual, "unable to run pre-write hook: boom")
		So(err.Unwrap().Error(), ShouldEqual, "boom")
	})
}
