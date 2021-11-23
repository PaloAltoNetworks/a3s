package token

import (
	"crypto"
	"crypto/sha1"
	"crypto/x509"
	"crypto/x509/pkix"
	"encoding/pem"
	"fmt"
	"reflect"
	"testing"
	"time"

	"github.com/golang-jwt/jwt/v4"
	. "github.com/smartystreets/goconvey/convey"
	"go.aporeto.io/a3s/pkgs/permissions"
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

func TestNewIdentityToken(t *testing.T) {

	Convey("Given I create a new Midgard claims", t, func() {

		c := NewIdentityToken(Source{
			Type:      "mtls",
			Namespace: "/my/ns",
			Name:      "mysource",
		})

		So(c.Identity, ShouldBeNil)
		So(c.Source.Type, ShouldEqual, "mtls")
		So(c.Source.Namespace, ShouldEqual, "/my/ns")
		So(c.Source.Name, ShouldEqual, "mysource")
	})
}

func TestParse(t *testing.T) {

	Convey("Given I create an IdentityToken", t, func() {

		cert, key := getECCert()

		token1 := NewIdentityToken(Source{
			Type:      "certificate",
			Namespace: "/my/ns",
			Name:      "mysource",
		})

		token1.Source.Type = "certificate"
		token1.Source.Namespace = "/my/ns"
		token1.Source.Name = "mysource"
		token1.Identity = []string{
			"org=a3s.com",
			"orgunit=admin",
			"commonname=joe",
		}

		keychain := NewJWKS()
		_ = keychain.Append(cert)

		kid := fmt.Sprintf("%02X", sha1.Sum(cert.Raw))

		token, err := token1.JWT(key, kid, "iss", jwt.ClaimStrings{"aud"}, time.Now().Add(10*time.Second), nil)
		So(err, ShouldBeNil)

		Convey("Calling JWT with a missing source type should fail", func() {
			token1.Source.Type = ""
			_, err := token1.JWT(key, "kid", "iss", jwt.ClaimStrings{"aud"}, time.Now().Add(10*time.Second), nil)
			So(err, ShouldNotBeNil)
			So(err.Error(), ShouldEqual, "invalid identity token: missing source type")
		})

		Convey("Calling JWT using cloaking", func() {
			token2, err := token1.JWT(key, kid, "iss", jwt.ClaimStrings{"aud"}, time.Now().Add(10*time.Second), []string{"org="})
			So(err, ShouldBeNil)
			token3, err := Parse(token2, keychain, "iss", "aud")
			So(err, ShouldBeNil)
			So(token3.Identity, ShouldResemble, []string{
				"org=a3s.com",
				"@source:type=certificate",
				"@source:namespace=/my/ns",
				"@source:name=mysource",
			})
		})

		Convey("When I call Parse using the correct signer certificate", func() {

			token2, err := Parse(token, keychain, "iss", "aud")

			So(err, ShouldBeNil)
			So(token2.Source.Type, ShouldEqual, "certificate")
			So(token2.Issuer, ShouldEqual, "iss")
			So(token2.Audience, ShouldResemble, jwt.ClaimStrings{"aud"})
			So(token2.ExpiresAt, ShouldResemble, token1.ExpiresAt)
			So(token2.IssuedAt, ShouldResemble, token1.IssuedAt)
			So(token2.Identity, ShouldResemble, []string{
				"org=a3s.com",
				"orgunit=admin",
				"commonname=joe",
				"@source:type=certificate",
				"@source:namespace=/my/ns",
				"@source:name=mysource",
			})
		})

		Convey("When I call Parse using the wrong issuer", func() {

			token2, err := Parse(token, keychain, "iss2", "aud")

			So(token2, ShouldBeNil)
			So(err, ShouldNotBeNil)
			So(err.Error(), ShouldEqual, "issuer 'iss' is not acceptable. want 'iss2'")
		})

		Convey("When I call Parse on a token missing the @source:type claim", func() {

			// Overwrite the test token
			claims := jwt.NewWithClaims(
				jwt.SigningMethodES256,
				jwt.MapClaims{
					"iss": "iss2",
				},
			)
			claims.Header["kid"] = kid

			token, _ := claims.SignedString(key)

			token2, err := Parse(token, keychain, "iss2", "aud")

			So(token2, ShouldBeNil)
			So(err, ShouldNotBeNil)
			So(err.Error(), ShouldEqual, "invalid token: missing @source:type in identity claims")
		})
		Convey("When I call Parse using the wrong audience", func() {

			token2, err := Parse(token, keychain, "iss", "aud2")

			So(token2, ShouldBeNil)
			So(err, ShouldNotBeNil)
			So(err.Error(), ShouldEqual, "audience '[aud]' is not acceptable. want 'aud2'")
		})

		Convey("When I call Parse using the wrong signer certificate", func() {

			cert2, _ := getECCert()
			keychain2 := NewJWKS()
			_ = keychain2.Append(cert2)
			token2, err := Parse(token, keychain2, "iss", "aud")

			So(token2, ShouldBeNil)
			So(err, ShouldNotBeNil)
			So(err.Error(), ShouldEqual, fmt.Sprintf("unable to parse jwt: unable to find kid '%s': kid not found in JWKS", kid))
		})

		Convey("When I call Parse using a wrong asigning method", func() {

			token = "eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9.eyJyZWFsbSI6IlZpbmNlIiwiaWR0IjpbIkBzb3VyY2U9YS1zb3VyY2UiLCJAbmFtZXNwYWNlPS9hL25hbWVzcGFjZSIsInVzZXJuYW1lPWpvZSIsInRlYW09cmVkIl0sImF1ZCI6Imh0dHBzOi8vYTNzLmNvbSIsImlhdCI6MTU0ODc5MDMwMCwiaXNzIjoiaHR0cHM6Ly9hM3MuY29tIn0.5PYuuULqrMArgdxq5eKSImsNskobw528Gr8Xe7HgPFs"
			token2, err := Parse(token, keychain, "iss", "aud")

			So(token2, ShouldBeNil)
			So(err, ShouldNotBeNil)
			So(err.Error(), ShouldEqual, "unable to parse jwt: unexpected signing method: HS256")
		})
	})
}

func TestIdentityToken_Restrict(t *testing.T) {
	type args struct {
		restrictions permissions.Restrictions
	}
	tests := []struct {
		name    string
		init    func(t *testing.T) *IdentityToken
		inspect func(r *IdentityToken, t *testing.T) //inspects receiver after test run

		args func(t *testing.T) args

		wantErr    bool
		inspectErr func(err error, t *testing.T) //use for more precise error evaluation after test
	}{
		{
			"empty existing restrictions, zero requested",
			func(*testing.T) *IdentityToken {
				return &IdentityToken{
					Restrictions: nil,
				}
			},
			func(r *IdentityToken, t *testing.T) {
				if r.Restrictions != nil {
					t.Fail()
				}
			},
			func(*testing.T) args {
				return args{
					permissions.Restrictions{},
				}
			},
			false,
			nil,
		},
		{
			"zero existing restrictions, zero requested",
			func(*testing.T) *IdentityToken {
				return &IdentityToken{
					Restrictions: &permissions.Restrictions{},
				}
			},
			func(r *IdentityToken, t *testing.T) {
				if r.Restrictions != nil {
					t.Fail()
				}
			},
			func(*testing.T) args {
				return args{
					permissions.Restrictions{},
				}
			},
			false,
			nil,
		},
		{
			"zero existing restrictions, requested",
			func(*testing.T) *IdentityToken {
				return &IdentityToken{
					Restrictions: &permissions.Restrictions{},
				}
			},
			func(r *IdentityToken, t *testing.T) {
				want := &permissions.Restrictions{
					Namespace:   "/the/ns",
					Networks:    []string{"10.0.0.0/24"},
					Permissions: []string{"dog:get,put"},
				}
				if !reflect.DeepEqual(r.Restrictions, want) {
					t.Logf("want %v got %v", want, r)
					t.Fail()
				}
			},
			func(*testing.T) args {
				return args{
					permissions.Restrictions{
						Namespace:   "/the/ns",
						Networks:    []string{"10.0.0.0/24"},
						Permissions: []string{"dog:get,put"},
					},
				}
			},
			false,
			nil,
		},
		{
			"existing restrictions, zero requested",
			func(*testing.T) *IdentityToken {
				return &IdentityToken{
					Restrictions: &permissions.Restrictions{
						Namespace:   "/the/ns",
						Networks:    []string{"10.0.0.0/24"},
						Permissions: []string{"dog:get,put"},
					},
				}
			},
			func(r *IdentityToken, t *testing.T) {
				want := &permissions.Restrictions{
					Namespace:   "/the/ns",
					Networks:    []string{"10.0.0.0/24"},
					Permissions: []string{"dog:get,put"},
				}
				if !reflect.DeepEqual(r.Restrictions, want) {
					t.Logf("want %v got %v", want, r)
					t.Fail()
				}
			},
			func(*testing.T) args {
				return args{
					permissions.Restrictions{},
				}
			},
			false,
			nil,
		},
		{
			"identical",
			func(*testing.T) *IdentityToken {
				return &IdentityToken{
					Restrictions: &permissions.Restrictions{
						Namespace:   "/the/ns",
						Networks:    []string{"10.0.0.0/24"},
						Permissions: []string{"dog:get,put"},
					},
				}
			},
			func(r *IdentityToken, t *testing.T) {
				want := &permissions.Restrictions{
					Namespace:   "/the/ns",
					Networks:    []string{"10.0.0.0/24"},
					Permissions: []string{"dog:get,put"},
				}
				if !reflect.DeepEqual(r.Restrictions, want) {
					t.Logf("want %v got %v", want, r)
					t.Fail()
				}
			},
			func(*testing.T) args {
				return args{
					permissions.Restrictions{
						Namespace:   "/the/ns",
						Networks:    []string{"10.0.0.0/24"},
						Permissions: []string{"dog:get,put"},
					},
				}
			},
			false,
			nil,
		},
		{
			"requested contained in existing",
			func(*testing.T) *IdentityToken {
				return &IdentityToken{
					Restrictions: &permissions.Restrictions{
						Namespace:   "/the/ns",
						Networks:    []string{"10.0.0.0/24"},
						Permissions: []string{"dog:get,put"},
					},
				}
			},
			func(r *IdentityToken, t *testing.T) {
				want := &permissions.Restrictions{
					Namespace:   "/the/ns/2",
					Networks:    []string{"10.0.0.0/32"},
					Permissions: []string{"dog:get"},
				}
				if !reflect.DeepEqual(r.Restrictions, want) {
					t.Logf("want %v got %v", want, r)
					t.Fail()
				}
			},
			func(*testing.T) args {
				return args{
					permissions.Restrictions{
						Namespace:   "/the/ns/2",
						Networks:    []string{"10.0.0.0/32"},
						Permissions: []string{"dog:get"},
					},
				}
			},
			false,
			nil,
		},
		{
			"breaking namespace",
			func(*testing.T) *IdentityToken {
				return &IdentityToken{
					Restrictions: &permissions.Restrictions{
						Namespace:   "/the/ns",
						Networks:    []string{"10.0.0.0/24"},
						Permissions: []string{"dog:get,put"},
					},
				}
			},
			nil,
			func(*testing.T) args {
				return args{
					permissions.Restrictions{
						Namespace:   "/the",
						Networks:    []string{"10.0.0.0/24"},
						Permissions: []string{"dog:get,put"},
					},
				}
			},
			true,
			func(err error, t *testing.T) {
				want := "restriction violation: restricted namespace must be empty, '/the/ns' or one of its children"
				if err.Error() != want {
					t.Logf("want error %s, got %s", want, err)
					t.Fail()
				}
			},
		},
		{
			"breaking networks",
			func(*testing.T) *IdentityToken {
				return &IdentityToken{
					Restrictions: &permissions.Restrictions{
						Namespace:   "/the/ns",
						Networks:    []string{"10.0.0.0/24"},
						Permissions: []string{"dog:get,put"},
					},
				}
			},
			nil,
			func(*testing.T) args {
				return args{
					permissions.Restrictions{
						Namespace:   "/the/ns",
						Networks:    []string{"11.0.0.0/24"},
						Permissions: []string{"dog:get,put"},
					},
				}
			},
			true,
			func(err error, t *testing.T) {
				want := "restriction violation: restricted networks must not overlap the current ones"
				if err.Error() != want {
					t.Logf("want error %s, got %s", want, err)
					t.Fail()
				}
			},
		},
		{
			"breaking permissions",
			func(*testing.T) *IdentityToken {
				return &IdentityToken{
					Restrictions: &permissions.Restrictions{
						Namespace:   "/the/ns",
						Networks:    []string{"10.0.0.0/24"},
						Permissions: []string{"dog:get,put"},
					},
				}
			},
			nil,
			func(*testing.T) args {
				return args{
					permissions.Restrictions{
						Namespace:   "/the/ns",
						Networks:    []string{"10.0.0.0/24"},
						Permissions: []string{"dog:get,put,create"},
					},
				}
			},
			true,
			func(err error, t *testing.T) {
				want := "restriction violation: restricted permissions must not be more permissive than the current ones"
				if err.Error() != want {
					t.Logf("want error %s, got %s", want, err)
					t.Fail()
				}
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			tArgs := tt.args(t)

			receiver := tt.init(t)
			err := receiver.Restrict(tArgs.restrictions)

			if tt.inspect != nil {
				tt.inspect(receiver, t)
			}

			if (err != nil) != tt.wantErr {
				t.Fatalf("IdentityToken.Restrict error = %v, wantErr: %t", err, tt.wantErr)
			}

			if tt.inspectErr != nil {
				tt.inspectErr(err, t)
			}
		})
	}
}
