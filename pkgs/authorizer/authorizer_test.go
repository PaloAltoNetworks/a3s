package authorizer

import (
	"context"
	"crypto"
	"crypto/x509"
	"crypto/x509/pkix"
	"encoding/pem"
	"fmt"
	"testing"
	"time"

	. "github.com/smartystreets/goconvey/convey"
	"go.aporeto.io/a3s/pkgs/notification"
	"go.aporeto.io/a3s/pkgs/nscache"
	"go.aporeto.io/a3s/pkgs/permissions"
	"go.aporeto.io/a3s/pkgs/token"
	"go.aporeto.io/bahamut"
	"go.aporeto.io/elemental"
	"go.aporeto.io/tg/tglib"
)

func getECCert() (*x509.Certificate, crypto.PrivateKey) {

	certBlock, keyBlock, err := tglib.Issue(pkix.Name{})
	if err != nil {
		panic(err)
	}

	cert, err := tglib.ParseCertificate(pem.EncodeToMemory(certBlock))
	if err != nil {
		panic(err)
	}

	key, err := tglib.PEMToKey(keyBlock)
	if err != nil {
		panic(err)
	}

	return cert, key
}

func makeToken(claims *token.IdentityToken, key crypto.PrivateKey) string {
	token, err := claims.JWT(key, "kid", "iss", "aud", time.Now().Add(time.Minute))
	if err != nil {
		panic(err)
	}
	return token
}

func TestNewAuthorizer(t *testing.T) {

	Convey("New should work", t, func() {

		p := bahamut.NewLocalPubSubClient()
		_ = p.Connect(context.Background())

		r := permissions.NewMockRetriever()
		a := New(context.Background(), r, p, OptionIgnoredResources("r1", "r2")).(*authorizer)

		Convey("Then a should be correct", func() {
			So(a.retriever, ShouldEqual, r)
			So(len(a.ignoredResources), ShouldEqual, 2)
			So(a.cache, ShouldNotBeNil)
		})
	})
}

func TestIsAuthorized(t *testing.T) {

	Convey("Given an Authorizer", t, func() {

		p := bahamut.NewLocalPubSubClient()
		_ = p.Connect(context.Background())

		r := permissions.NewMockRetriever()
		a := New(context.Background(), r, p, OptionIgnoredResources("r1", "r2")).(*authorizer)

		Convey("Callin on an ignored resource should worlk", func() {

			bctx := bahamut.NewMockContext(context.Background())
			bctx.MockRequest = &elemental.Request{Identity: elemental.MakeIdentity("r1", "r1")}

			action, err := a.IsAuthorized(bctx)
			So(err, ShouldBeNil)
			So(action, ShouldEqual, bahamut.AuthActionOK)
		})

		Convey("Calling with a missing token should fail", func() {

			bctx := bahamut.NewMockContext(context.Background())
			bctx.MockRequest = &elemental.Request{
				Identity:  elemental.MakeIdentity("r0", "r0"),
				Namespace: "/",
			}

			action, err := a.IsAuthorized(bctx)
			So(err, ShouldNotBeNil)
			So(err.Error(), ShouldEqual, "error 403 (a3s:authorizer): Forbidden: Missing token in either Authorization header or X-A3S-Token in cookies")
			So(action, ShouldEqual, bahamut.AuthActionKO)
		})

		Convey("Calling with a token with invalid token should fail", func() {

			bctx := bahamut.NewMockContext(context.Background())
			bctx.MockRequest = &elemental.Request{
				Identity:  elemental.MakeIdentity("r0", "r0"),
				Namespace: "/",
				Password:  "not.a.token",
			}

			action, err := a.IsAuthorized(bctx)
			So(err, ShouldNotBeNil)
			So(err.Error(), ShouldEqual, "error 403 (a3s:authorizer): Forbidden: unable to compute authz restrictions from token: invalid character '\\u009e' looking for beginning of value")
			So(action, ShouldEqual, bahamut.AuthActionKO)
		})

		Convey("Calling with a token with valid token and permissions are granted should work", func() {

			_, key := getECCert()
			token := makeToken(&token.IdentityToken{
				Source: token.Source{Type: "mtls"},
			}, key)

			bctx := bahamut.NewMockContext(context.Background())
			bctx.MockRequest = &elemental.Request{
				Identity:  elemental.MakeIdentity("r0", "r0"),
				Namespace: "/",
				Password:  token,
				Operation: elemental.OperationRetrieveMany,
			}

			r.MockPermissions(t, func(context.Context, []string, string, ...permissions.RetrieverOption) (permissions.PermissionMap, error) {

				return permissions.PermissionMap{"r0": permissions.Permissions{"retrieve-many": true}}, nil
			})
			action, err := a.IsAuthorized(bctx)
			So(err, ShouldBeNil)
			So(action, ShouldEqual, bahamut.AuthActionOK)
		})

		Convey("Calling with a token with valid token and permissions are denied should fail", func() {

			_, key := getECCert()
			token := makeToken(&token.IdentityToken{
				Source: token.Source{Type: "mtls"},
			}, key)

			bctx := bahamut.NewMockContext(context.Background())
			bctx.MockRequest = &elemental.Request{
				Identity:  elemental.MakeIdentity("r0", "r0"),
				Namespace: "/",
				Password:  token,
				Operation: elemental.OperationRetrieveMany,
			}

			r.MockPermissions(t, func(context.Context, []string, string, ...permissions.RetrieverOption) (permissions.PermissionMap, error) {

				return permissions.PermissionMap{"r0": permissions.Permissions{}}, nil
			})

			action, err := a.IsAuthorized(bctx)
			So(err, ShouldBeNil)
			So(action, ShouldEqual, bahamut.AuthActionKO)
		})

		Convey("Calling with a token with valid token and retriever errors should fail", func() {

			_, key := getECCert()
			token := makeToken(&token.IdentityToken{
				Source: token.Source{Type: "mtls"},
			}, key)

			bctx := bahamut.NewMockContext(context.Background())
			bctx.MockRequest = &elemental.Request{
				Identity:  elemental.MakeIdentity("r0", "r0"),
				Namespace: "/",
				Password:  token,
				Operation: elemental.OperationRetrieveMany,
			}

			r.MockPermissions(t, func(context.Context, []string, string, ...permissions.RetrieverOption) (permissions.PermissionMap, error) {

				return nil, fmt.Errorf("boom")
			})
			action, err := a.IsAuthorized(bctx)
			So(err, ShouldNotBeNil)
			So(err.Error(), ShouldEqual, "boom")
			So(action, ShouldEqual, bahamut.AuthActionKO)
		})
	})
}

func TestCheckPermissions(t *testing.T) {

	Convey("Given an Authorizer", t, func() {

		p := bahamut.NewLocalPubSubClient()
		_ = p.Connect(context.Background())

		r := permissions.NewMockRetriever()
		a := New(context.Background(), r, p, OptionIgnoredResources("r1", "r2")).(*authorizer)

		Convey("Calling on an ignored resource should work", func() {
			ok, err := a.CheckAuthorization(context.Background(), []string{}, "retrieve-many", "/", "r1")
			So(err, ShouldBeNil)
			So(ok, ShouldBeTrue)
		})

		Convey("Calling with a missing namespace should fail", func() {

			ok, err := a.CheckAuthorization(context.Background(), []string{}, "retrieve-many", "", "r0")

			So(err, ShouldNotBeNil)
			So(err.Error(), ShouldEqual, "error 403 (a3s:authorizer): Forbidden: Missing X-Namespace header")
			So(ok, ShouldBeFalse)
		})

		Convey("Calling with an incorrect namespace should fail", func() {

			ok, err := a.CheckAuthorization(context.Background(), []string{}, "retrieve-many", "no", "r0")

			So(err, ShouldNotBeNil)
			So(err.Error(), ShouldEqual, "error 403 (a3s:authorizer): Forbidden: Invalid X-Namespace header. A namespace must start with /")
			So(ok, ShouldBeFalse)
		})

		Convey("Calling when retriever grants permission should work", func() {

			r.MockPermissions(t, func(context.Context, []string, string, ...permissions.RetrieverOption) (permissions.PermissionMap, error) {
				return permissions.PermissionMap{"r0": permissions.Permissions{"retrieve-many": true}}, nil
			})

			ok, err := a.CheckAuthorization(context.Background(), []string{}, "retrieve-many", "/", "r0")

			So(err, ShouldBeNil)
			So(ok, ShouldBeTrue)

			Convey("Calling one more time should use the cache and work", func() {

				r.MockPermissions(t, func(context.Context, []string, string, ...permissions.RetrieverOption) (permissions.PermissionMap, error) {
					return permissions.PermissionMap{}, nil // this simulates a changes that was not pushed, so cache will be used.
				})

				ok, err := a.CheckAuthorization(context.Background(), []string{}, "retrieve-many", "/", "r0")

				So(err, ShouldBeNil)
				So(ok, ShouldBeTrue)

				Convey("Calling once again after a push should fail", func() {

					pub := bahamut.NewPublication(nscache.NotificationNamespaceChanges)
					_ = pub.Encode(notification.Message{Data: "/"})
					_ = p.Publish(pub)

					time.Sleep(300 * time.Millisecond) // give a bit of time

					ok, err := a.CheckAuthorization(context.Background(), []string{}, "retrieve-many", "/", "r0")

					So(err, ShouldBeNil)
					So(ok, ShouldBeFalse)
				})
			})
		})

		Convey("Calling when retriever denies permission should fail", func() {

			r.MockPermissions(t, func(context.Context, []string, string, ...permissions.RetrieverOption) (permissions.PermissionMap, error) {
				return permissions.PermissionMap{}, nil
			})

			ok, err := a.CheckAuthorization(context.Background(), []string{}, "retrieve-many", "/", "r0")

			So(err, ShouldBeNil)
			So(ok, ShouldBeFalse)
		})

		Convey("Calling when retriever errors should fail", func() {

			r.MockPermissions(t, func(context.Context, []string, string, ...permissions.RetrieverOption) (permissions.PermissionMap, error) {
				return nil, fmt.Errorf("bim")
			})

			ok, err := a.CheckAuthorization(context.Background(), []string{}, "retrieve-many", "/", "r0")

			So(err, ShouldNotBeNil)
			So(err.Error(), ShouldEqual, "bim")
			So(ok, ShouldBeFalse)
		})

		Convey("Calling with options should forward them to the retriever", func() {

			var expectedOptions []permissions.RetrieverOption
			var expectedNamespace string
			var expectedClaims []string
			r.MockPermissions(t, func(ctx context.Context, claims []string, ns string, opts ...permissions.RetrieverOption) (permissions.PermissionMap, error) {
				expectedClaims = claims
				expectedNamespace = ns
				expectedOptions = opts
				return permissions.PermissionMap{"r0": permissions.Permissions{"retrieve-many": true}}, nil
			})

			_, _ = a.CheckAuthorization(context.Background(), []string{"a=a"}, "retrieve-many", "/", "r0",
				OptionCheckSourceIP("1.1.1.1"),
				OptionCheckID("id"),
				OptionCheckRestrictions(permissions.Restrictions{}),
			)

			So(expectedNamespace, ShouldEqual, "/")
			So(expectedClaims, ShouldResemble, []string{"a=a"})
			So(len(expectedOptions), ShouldEqual, 3)
		})
	})
}
