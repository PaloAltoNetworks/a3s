package authlib

import (
	"context"
	"fmt"
	"net/http"
	"net/http/httptest"
	"sync/atomic"
	"testing"
	"time"

	. "github.com/smartystreets/goconvey/convey"
	"go.aporeto.io/a3s/pkgs/api"
	"go.aporeto.io/elemental"
	"go.aporeto.io/manipulate"
	"go.aporeto.io/manipulate/maniptest"
)

func TestTokenManager_Issue(t *testing.T) {

	Convey("Calling NewPeriodicTokenManager with a nil function should panic", t, func() {
		So(func() { NewPeriodicTokenManager(10*time.Second, nil) }, ShouldPanicWith, "issuerFunc cannot be nil")
	})

	Convey("Given I have TokenIssuerFunc that works and a token manager", t, func() {

		tf := func(ctx context.Context, v time.Duration) (string, error) {
			return "token!", nil
		}

		tm := NewPeriodicTokenManager(10*time.Second, tf)

		t, err := tm.Issue(context.Background())

		So(err, ShouldBeNil)
		So(t, ShouldEqual, "token!")
	})
}

func TestTokenManager_Run(t *testing.T) {

	tickDuration = 1 * time.Millisecond

	Convey("Given I have TokenIssuerFunc that works and a token manager", t, func() {

		var called int32
		tf := func(ctx context.Context, v time.Duration) (string, error) {
			atomic.AddInt32(&called, 1)
			return "token!", nil
		}

		tm := NewPeriodicTokenManager(2*time.Millisecond, tf)

		Convey("When I call Run and wait for a few", func() {

			ctx, cancel := context.WithTimeout(context.Background(), time.Second)
			defer cancel()

			tokenCh := make(chan string)
			go tm.Run(ctx, tokenCh)

			var c int
			var lastToken string
		L:
			for {
				select {
				case lastToken = <-tokenCh:
					c++
					if c == 4 {
						break L
					}
				case <-ctx.Done():
					panic("timeout exceeded")
				}
			}

			So(c, ShouldEqual, 4)
			So(atomic.LoadInt32(&called), ShouldEqual, 4)
			So(lastToken, ShouldEqual, "token!")
		})
	})

	Convey("Given I have TokenIssuerFunc that fails and a token manager", t, func() {

		var called int32
		tf := func(ctx context.Context, v time.Duration) (string, error) {
			atomic.AddInt32(&called, 1)
			return "", fmt.Errorf("bim")
		}

		tm := NewPeriodicTokenManager(2*time.Millisecond, tf)

		Convey("When I call Run and wait for a few", func() {

			ctx, cancel := context.WithTimeout(context.Background(), 1*time.Second)

			tokenCh := make(chan string)
			go tm.Run(ctx, tokenCh)

		L:
			for {
				select {
				case <-tokenCh:
					panic("received a token")
				case <-ctx.Done():
					break L
				}
			}

			cancel()

			So(atomic.LoadInt32(&called), ShouldBeGreaterThan, 0)
		})
	})
}

func TestNewX509TokenManager(t *testing.T) {

	Convey("Given I can NewX509TokenManager ", t, func() {

		ts := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			http.Error(w, "{}", http.StatusForbidden)
		}))
		defer ts.Close()

		tm := NewX509TokenManager("sourceName", "/ns", OptValidity(10*time.Second))
		m := maniptest.NewTestManipulator()
		tm.SetManipulator(m)

		So(tm.(*x509TokenManager).validity, ShouldEqual, 10*time.Second)
		So(tm.(*x509TokenManager).issuerFunc, ShouldNotBeNil)

		Convey("When I call the the issue func", func() {

			m.MockCreate(t, func(mctx manipulate.Context, object elemental.Identifiable) error {
				*object.(*api.Issue) = api.Issue{
					Token: "hello",
				}
				return nil
			})

			token, err := tm.(*x509TokenManager).issuerFunc(context.Background(), 10*time.Second)

			So(err, ShouldBeNil)
			So(token, ShouldEqual, "hello")
		})
	})
}
