package nscache

import (
	"testing"

	. "github.com/smartystreets/goconvey/convey"
)

func TestOption(t *testing.T) {

	Convey("Given a new config", t, func() {

		c := newConfig()

		So(c.notificationName, ShouldEqual, NotificationNamespaceChanges)

		Convey("OptionNotificationName should work", func() {
			OptionNotificationName("coucou")(&c)
			So(c.notificationName, ShouldEqual, "coucou")
		})
	})
}
