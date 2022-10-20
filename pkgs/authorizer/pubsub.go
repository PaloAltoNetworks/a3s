package authorizer

import (
	"context"

	"go.aporeto.io/a3s/pkgs/api"
	"go.aporeto.io/a3s/pkgs/notification"
	"go.aporeto.io/a3s/pkgs/nscache"
	"go.aporeto.io/bahamut"
	"go.aporeto.io/manipulate"
)

type eventData struct {
	Namespace string `json:"namespace" msgpack:"namespace"`
	Name      string `json:"name" msgpack:"name"`
}

// webSocketPubSub is a naive bahamut.PubSubClient internal implementation
// that is backed by a manipulate.Subscriber. This is used to
// make the Authorizer working when used by third party clients that
// won't have access to the internal NATS notification topic.
// It basically acts a shim layer that translates classic elemental.Events
// into the relevant notification.Message used by the authorizer internal
// namespace cache.
type webSocketPubSub struct {
	subscriber manipulate.Subscriber
}

// not implemented. These are just here to satisfy the bahamut.PubSubClient interface.
func (w *webSocketPubSub) Connect(context.Context) error { return nil }
func (w *webSocketPubSub) Disconnect() error             { return nil }
func (w *webSocketPubSub) Publish(*bahamut.Publication, ...bahamut.PubSubOptPublish) error {
	return nil
}

func (w *webSocketPubSub) Subscribe(pubs chan *bahamut.Publication, errors chan error, topic string, opts ...bahamut.PubSubOptSubscribe) func() {

	sendErr := func(err error) {
		select {
		case errors <- err:
		default:
		}
	}

	sendPub := func(pub *bahamut.Publication) {
		select {
		case pubs <- pub:
		default:
		}
	}

	go func() {

		for {
			select {

			case evt := <-w.subscriber.Events():

				// We decode the vent in a generic container structure.
				d := &eventData{}
				if err := evt.Decode(d); err != nil {
					sendErr(err)
					break
				}

				// We prepare a notification Message that the authorizer
				// nscache will understand.
				msg := notification.Message{
					Type: nscache.NotificationNamespaceChanges,
				}

				// We populate the namespace name based on the
				// event identity.
				switch evt.Identity {
				case api.NamespaceIdentity.Name:
					msg.Data = d.Name
				case api.AuthorizationIdentity.Name:
					msg.Data = d.Namespace
				}

				// Then we create a publication and wrap the msg inside.
				p := bahamut.NewPublication(topic)
				if err := p.Encode(msg); err != nil {
					sendErr(err)
					break
				}

				sendPub(p)

			case st := <-w.subscriber.Status():
				if st == manipulate.SubscriberStatusFinalDisconnection {
					return
				}

			case err := <-w.subscriber.Errors():
				sendErr(err)
			}
		}

	}()

	return func() { w.subscriber.Status() <- manipulate.SubscriberStatusFinalDisconnection }
}
