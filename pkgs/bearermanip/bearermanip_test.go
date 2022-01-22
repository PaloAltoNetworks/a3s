package bearermanip

import (
	"context"
	"crypto/tls"
	"testing"

	. "github.com/smartystreets/goconvey/convey"
	"go.aporeto.io/bahamut"
	"go.aporeto.io/elemental"
	"go.aporeto.io/manipulate/maniphttp"
)

func TestConfigure(t *testing.T) {

	Convey("Calling Configure should work", t, func() {
		api := "https://toto.com"
		tlsConfig := &tls.Config{}

		f := Configure(context.Background(), api, tlsConfig)
		So(f, ShouldHaveSameTypeAs, (MakerFunc)(nil))

		bctx := bahamut.NewMockContext(context.Background())
		bctx.MockRequest = elemental.NewRequest()
		bctx.MockRequest.Password = "the-token"
		bctx.MockRequest.Namespace = "/ns"

		m := f(bctx)
		user, pass := maniphttp.ExtractCredentials(m)
		So(user, ShouldEqual, "Bearer")
		So(pass, ShouldEqual, "the-token")
		So(maniphttp.ExtractTLSConfig(m), ShouldResemble, tlsConfig)
		So(maniphttp.ExtractNamespace(m), ShouldEqual, "/ns")
		So(maniphttp.ExtractEndpoint(m), ShouldEqual, api)
	})
}
