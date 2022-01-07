package authorizer

import (
	"context"
	"fmt"
	"testing"
	"time"

	. "github.com/smartystreets/goconvey/convey"
	"go.aporeto.io/a3s/pkgs/api"
	"go.aporeto.io/a3s/pkgs/notification"
	"go.aporeto.io/a3s/pkgs/nscache"
	"go.aporeto.io/bahamut"
	"go.aporeto.io/elemental"
	"go.aporeto.io/manipulate"
	"go.aporeto.io/manipulate/maniptest"
)

func TestNotImplemented(t *testing.T) {
	// free coverage \o/
	Convey("Given an authorizer", t, func() {
		s := webSocketPubSub{}
		So(s.Connect(context.Background()), ShouldBeNil)
		So(s.Disconnect(), ShouldBeNil)
		So(s.Publish(nil), ShouldBeNil)
	})
}

func TestSubscribe(t *testing.T) {

	Convey("I have a subscriber and a ws authorizer", t, func() {

		s := maniptest.NewTestSubscriber()
		a := webSocketPubSub{subscriber: s}
		chInEvents := make(chan *elemental.Event, 2)
		chInStatus := make(chan manipulate.SubscriberStatus, 2)
		chInErrors := make(chan error, 2)

		s.MockEvents(t, func() chan *elemental.Event {
			return chInEvents
		})
		s.MockStatus(t, func() chan manipulate.SubscriberStatus {
			return chInStatus
		})
		s.MockErrors(t, func() chan error {
			return chInErrors
		})

		checkPresent := func(chOutPubs chan *bahamut.Publication) {
			pub := <-chOutPubs
			msg := notification.Message{}
			pub.Decode(&msg) // nolint
			So(msg.Type, ShouldEqual, nscache.NotificationNamespaceChanges)
			So(msg.Data, ShouldEqual, "/the/ns")
		}

		Convey("when I receive a push from a namespace", func() {
			chOutPubs := make(chan *bahamut.Publication, 2)
			chOutErrs := make(chan error, 2)
			a.Subscribe(chOutPubs, chOutErrs, "topic")
			chInEvents <- elemental.NewEvent(elemental.EventUpdate, &api.Namespace{Name: "/the/ns"})
			checkPresent(chOutPubs)
		})

		Convey("when I receive a push from a authorization", func() {
			chOutPubs := make(chan *bahamut.Publication, 2)
			chOutErrs := make(chan error, 2)
			a.Subscribe(chOutPubs, chOutErrs, "topic")
			chInEvents <- elemental.NewEvent(elemental.EventUpdate, &api.Authorization{Namespace: "/the/ns"})
			checkPresent(chOutPubs)
		})

		Convey("when I receive an error", func() {
			chOutPubs := make(chan *bahamut.Publication, 2)
			chOutErrs := make(chan error, 2)
			a.Subscribe(chOutPubs, chOutErrs, "topic")
			chInErrors <- fmt.Errorf("boom")
			e := <-chOutErrs
			So(e.Error(), ShouldEqual, "boom")
		})

		Convey("when I receive a final disconnect", func() {
			chOutPubs := make(chan *bahamut.Publication, 2)
			chOutErrs := make(chan error, 2)
			a.Subscribe(chOutPubs, chOutErrs, "topic")
			chInStatus <- manipulate.SubscriberStatusFinalDisconnection
			time.Sleep(300 * time.Millisecond)
		})

		Convey("when the even it not decodable", func() {
			chOutPubs := make(chan *bahamut.Publication, 2)
			chOutErrs := make(chan error, 2)
			a.Subscribe(chOutPubs, chOutErrs, "topic")
			chInEvents <- &elemental.Event{RawData: []byte("oh no")}
			e := <-chOutErrs
			So(e.Error(), ShouldEqual, "unable to decode application/json: EOF")
		})
	})
}
