package a3sissuer

import (
	"crypto"
	"crypto/x509"
	"crypto/x509/pkix"
	"encoding/pem"
	"reflect"
	"testing"
	"time"

	"github.com/golang-jwt/jwt/v4"
	. "github.com/smartystreets/goconvey/convey"
	"go.aporeto.io/a3s/pkgs/permissions"
	"go.aporeto.io/a3s/pkgs/token"
	"go.aporeto.io/tg/tglib"
)

func getECCert(subject pkix.Name, opts ...tglib.IssueOption) (*x509.Certificate, crypto.PrivateKey) {

	certBlock, keyBlock, err := tglib.Issue(subject, opts...)
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

func TestNewA3SIssuer(t *testing.T) {
	Convey("Given I call newA3SIssuer", t, func() {
		c := newA3SIssuer()
		So(c, ShouldNotBeNil)
	})
}

func TestIssue(t *testing.T) {
	Convey("Given I call ToMidgardClaims", t, func() {
		c := newA3SIssuer()
		ot := token.NewIdentityToken(token.Source{})
		c.token = ot
		So(c.Issue(), ShouldEqual, ot)
	})
}

func TestFromToken(t *testing.T) {

	cert, key := getECCert(pkix.Name{})
	keychain := token.NewJWKS()
	_ = keychain.Append(cert)
	kid := token.Fingerprint(cert)

	Convey("Using a token with an bad restrictions", t, func() {
		token := `eyJhbGciOiJFUzI1NiIsInR5cCI6IkpXVCJ9.eyJyZWFsbSI6IlZpbmNlIiwiZGF0YSI6eyJhY2NvdW50IjoiYXBvbXV4IiwiZW1haWwiOiJhZG1pbkBhcG9tdXguY29tIiwiaWQiOiI1ZTFjZjNlZmEzNzAwMzhmYWY3Zjg3NzciLCJvcmdhbml6YXRpb24iOiJhcG9tdXgiLCJyZWFsbSI6InZpbmNlIiwic3ViamVjdCI6ImFwb211eCJ9LCJyZXN0cmljdGlvbnMiOnsibmV0d29ya3MiOiIxMjcuMC4wLjEvMzIifSwiZXhwIjoxNTkwMDQzMjA1LCJpYXQiOjE1ODk5NTMyMDUsImlzcyI6Imh0dHBzOi8vbG9jYWxob3N0OjQ0NDMiLCJzdWIiOiJhcG9tdXgifQ.dIsnGMSEy961FqXgJH-TBVw8_9VrzH_j4xcQJG4JY0--ekwNuMpLr0CyOJFj_XFuVsY-ZS8Lwj5yJCYHv7TS8Q`
		c := newA3SIssuer()
		err := c.fromToken(token, keychain, "", nil, 0, permissions.Restrictions{})
		So(err, ShouldNotBeNil)
		So(err.Error(), ShouldEqual, `unable to compute restrictions: unable to compute authz restrictions from token: json: cannot unmarshal string into Go struct field Restrictions.restrictions.networks of type []string`)
	})

	Convey("Using a token that is missing kid", t, func() {
		token := `eyJhbGciOiJFUzI1NiIsInR5cCI6IkpXVCJ9.eyJyZWFsbSI6IlZpbmNlIiwiZGF0YSI6eyJhY2NvdW50IjoiYXBvbXV4IiwiZW1haWwiOiJhZG1pbkBhcG9tdXguY29tIiwiaWQiOiI1ZTFjZjNlZmEzNzAwMzhmYWY3Zjg3NzciLCJvcmdhbml6YXRpb24iOiJhcG9tdXgiLCJyZWFsbSI6InZpbmNlIiwic3ViamVjdCI6ImFwb211eCJ9LCJyZXN0cmljdGlvbnMiOnt9LCJleHAiOjE1OTAzMDQzNDgsImlhdCI6MTU5MDIxNDM0OCwiaXNzIjoiaHR0cHM6Ly9sb2NhbGhvc3Q6NDQ0MyIsInN1YiI6ImFwb211eCJ9.7TZEEG-M-Ed-pKTzEGVZnKKZ1fvG0P7kN-VIKnVn_4TkTR2PX0EaToNZViGgcIs6pYXm7SByzjMl63ZiriSYkg`
		c := newA3SIssuer()
		err := c.fromToken(token, keychain, "", nil, 0, permissions.Restrictions{})
		So(err, ShouldNotBeNil)
		So(err.Error(), ShouldEqual, `unable to parse input token: unable to parse jwt: token has no KID in its header`)
	})

	Convey("Using a token that has no restrictions", t, func() {

		mc := token.NewIdentityToken(token.Source{Type: "mtls"})
		mc.ExpiresAt = jwt.NewNumericDate(time.Now().Add(time.Hour))
		mc.Issuer = "iss"

		token, _ := mc.JWT(key, kid, "iss", jwt.ClaimStrings{"aud"}, time.Time{}, nil)
		c := newA3SIssuer()
		err := c.fromToken(token, keychain, "iss", jwt.ClaimStrings{"aud"}, 0, permissions.Restrictions{})

		So(err, ShouldBeNil)
		So(c.token.Restrictions, ShouldBeNil)
	})

	Convey("When token initial audience match", t, func() {

		mc := token.NewIdentityToken(token.Source{Type: "mtls"})
		mc.ExpiresAt = jwt.NewNumericDate(time.Now().Add(time.Hour))
		mc.Issuer = "iss"

		token, _ := mc.JWT(key, kid, "iss", jwt.ClaimStrings{"aud1", "aud2"}, time.Time{}, nil)
		c := newA3SIssuer()
		err := c.fromToken(token, keychain, "iss", jwt.ClaimStrings{"aud1", "aud2"}, 0, permissions.Restrictions{})

		So(err, ShouldBeNil)
		So(c.token.Restrictions, ShouldBeNil)
	})

	Convey("When token initial audience contains the new requested audience", t, func() {

		mc := token.NewIdentityToken(token.Source{Type: "mtls"})
		mc.ExpiresAt = jwt.NewNumericDate(time.Now().Add(time.Hour))
		mc.Issuer = "iss"

		token, _ := mc.JWT(key, kid, "iss", jwt.ClaimStrings{"aud1", "aud2"}, time.Time{}, nil)
		c := newA3SIssuer()
		err := c.fromToken(token, keychain, "iss", jwt.ClaimStrings{"aud2"}, 0, permissions.Restrictions{})

		So(err, ShouldBeNil)
		So(c.token.Restrictions, ShouldBeNil)
	})

	Convey("When token initial audience does not contain the new requested audience", t, func() {

		mc := token.NewIdentityToken(token.Source{Type: "mtls"})
		mc.ExpiresAt = jwt.NewNumericDate(time.Now().Add(time.Hour))
		mc.Issuer = "iss"

		token, _ := mc.JWT(key, kid, "iss", jwt.ClaimStrings{"aud1", "aud2"}, time.Time{}, nil)
		c := newA3SIssuer()
		err := c.fromToken(token, keychain, "iss", jwt.ClaimStrings{"aud2", "aud3"}, 0, permissions.Restrictions{})

		So(err, ShouldNotBeNil)
		So(err.Error(), ShouldEqual, "unable to parse input token: requested audience 'aud3' is not declared in initial token")
	})

	Convey("When requested token audience is empty", t, func() {

		mc := token.NewIdentityToken(token.Source{Type: "mtls"})
		mc.ExpiresAt = jwt.NewNumericDate(time.Now().Add(time.Hour))
		mc.Issuer = "iss"

		token, _ := mc.JWT(key, kid, "iss", jwt.ClaimStrings{"aud1", "aud2"}, time.Time{}, nil)
		c := newA3SIssuer()
		err := c.fromToken(token, keychain, "iss", nil, 0, permissions.Restrictions{})

		So(err, ShouldNotBeNil)
		So(err.Error(), ShouldEqual, "unable to parse input token: you cannot request a token with no audience from a token that has one")
	})
}

func Test_computeNewValidity(t *testing.T) {

	now := jwt.NewNumericDate(time.Now())
	exp := jwt.NewNumericDate(now.Add(time.Hour))

	type args struct {
		originalExpUNIX   *jwt.NumericDate
		requestedValidity time.Duration
		renew             bool
	}
	tests := []struct {
		name    string
		args    args
		want    *jwt.NumericDate
		wantErr bool
	}{
		{
			"no original",
			args{
				nil,
				0,
				false,
			},
			nil,
			true,
		},
		{
			"no requested",
			args{
				exp,
				0,
				false,
			},
			exp,
			false,
		},
		{
			"correct requested",
			args{
				exp,
				30 * time.Minute,
				false,
			},
			jwt.NewNumericDate(now.Add(30 * time.Minute)),
			false,
		},
		{
			"requested too big",
			args{
				exp,
				48 * time.Hour,
				false,
			},
			exp,
			false,
		},
		{
			"requested the same",
			args{
				exp,
				time.Until(exp.Local()),
				false,
			},
			exp,
			false,
		},
		{
			"request bigger but with renew on",
			args{
				exp,
				48 * time.Hour,
				true,
			},
			jwt.NewNumericDate(now.Add(48 * time.Hour)),
			false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := computeNewValidity(tt.args.originalExpUNIX, tt.args.requestedValidity, tt.args.renew)
			if (err != nil) != tt.wantErr {
				t.Errorf("computeNewValidity() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if !reflect.DeepEqual(got, tt.want) {
				t.Errorf("computeNewValidity() = %v, want %v", got, tt.want)
			}
		})
	}
}
