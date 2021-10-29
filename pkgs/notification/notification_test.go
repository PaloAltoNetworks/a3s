package notification

import (
	"context"
	"testing"
	"time"

	. "github.com/smartystreets/goconvey/convey"
	"github.com/ugorji/go/codec"
	"go.aporeto.io/bahamut"
	"go.aporeto.io/elemental"
)

type unencodable struct{}

func (u *unencodable) CodecEncodeSelf(e *codec.Encoder) {
	panic("bim")
}
func (u *unencodable) CodecDecodeSelf(e *codec.Decoder) {}

type undecodable struct{}

func (u *undecodable) CodecEncodeSelf(e *codec.Encoder) {}
func (u *undecodable) CodecDecodeSelf(e *codec.Decoder) {
	panic("bam")
}

func TestPublish(t *testing.T) {

	Convey("Given I have a pubsub client", t, func() {

		pubsub := bahamut.NewLocalPubSubClient()
		_ = pubsub.Connect(context.Background())

		pubs := make(chan *bahamut.Publication)
		pubsub.Subscribe(pubs, nil, "test")

		Convey("When everything works", func() {

			msg := Message{
				Type: "type",
				Data: "hello world",
			}

			err := Publish(pubsub, "test", &msg)
			So(err, ShouldBeNil)

			var p *bahamut.Publication
			select {
			case p = <-pubs:
			case <-time.After(300 * time.Millisecond):
				panic("test did not get response in time")
			}

			d, _ := elemental.Encode(elemental.EncodingTypeMSGPACK, msg)
			So(string(p.Data), ShouldResemble, string(d))
		})

		Convey("When encoding fails", func() {

			msg := Message{
				Type: "type",
				Data: &unencodable{},
			}

			err := Publish(pubsub, "test", &msg)
			So(err, ShouldNotBeNil)
			So(err.Error(), ShouldEqual, "unable to encode notification publication: unable to encode application/msgpack: msgpack encode error: bim")
		})
	})
}

func TestSubscribe(t *testing.T) {

	Convey("Given I have a pubsub client", t, func() {

		ctx, cancel := context.WithCancel(context.Background())
		defer cancel()

		pubsub := bahamut.NewLocalPubSubClient()
		_ = pubsub.Connect(ctx)

		recvmsg := make(chan *Message, 2)
		h := func(msg *Message) {
			recvmsg <- msg
		}

		Convey("When everything works", func() {

			Subscribe(ctx, pubsub, "test", h)

			pub := bahamut.NewPublication("test")
			_ = pub.Encode(&Message{Type: "type"})
			_ = pubsub.Publish(pub)

			var msg *Message
			select {
			case msg = <-recvmsg:
			case <-time.After(300 * time.Millisecond):
				panic("test did not get response in time")
			}

			So(msg, ShouldNotBeNil)
		})

		Convey("When message is not decodable", func() {

			Subscribe(ctx, pubsub, "test", h)

			pub := bahamut.NewPublication("test")
			_ = pub.Encode(&Message{Data: &undecodable{}})
			_ = pubsub.Publish(pub)

			select {
			case <-recvmsg:
				t.Log("received unwanted message")
				t.Fail()
			case <-time.After(300 * time.Millisecond):
			}
		})
	})
}
