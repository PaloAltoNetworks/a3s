package claims

import (
	"crypto/x509/pkix"
	"testing"
	"time"

	. "github.com/smartystreets/goconvey/convey"
	"go.aporeto.io/a3s/pkgs/permissions"
	"go.aporeto.io/a3s/pkgs/token"
)

func TestNewTokenIssuer(t *testing.T) {
	Convey("Given I call NewTokenIssuer", t, func() {
		c := NewTokenIssuer()
		So(c, ShouldNotBeNil)
	})
}

func TestIssue(t *testing.T) {
	Convey("Given I call ToMidgardClaims", t, func() {
		c := NewTokenIssuer()
		ot := token.NewIdentityToken(token.Source{})
		c.token = ot
		So(c.Issue(), ShouldEqual, ot)
	})
}

func TestFromToken(t *testing.T) {

	cert, key := getECCert(pkix.Name{})

	Convey("Using a token with an bad restrictions", t, func() {
		token := `eyJhbGciOiJFUzI1NiIsInR5cCI6IkpXVCJ9.eyJyZWFsbSI6IlZpbmNlIiwiZGF0YSI6eyJhY2NvdW50IjoiYXBvbXV4IiwiZW1haWwiOiJhZG1pbkBhcG9tdXguY29tIiwiaWQiOiI1ZTFjZjNlZmEzNzAwMzhmYWY3Zjg3NzciLCJvcmdhbml6YXRpb24iOiJhcG9tdXgiLCJyZWFsbSI6InZpbmNlIiwic3ViamVjdCI6ImFwb211eCJ9LCJyZXN0cmljdGlvbnMiOnsibmV0d29ya3MiOiIxMjcuMC4wLjEvMzIifSwiZXhwIjoxNTkwMDQzMjA1LCJpYXQiOjE1ODk5NTMyMDUsImlzcyI6Imh0dHBzOi8vbG9jYWxob3N0OjQ0NDMiLCJzdWIiOiJhcG9tdXgifQ.dIsnGMSEy961FqXgJH-TBVw8_9VrzH_j4xcQJG4JY0--ekwNuMpLr0CyOJFj_XFuVsY-ZS8Lwj5yJCYHv7TS8Q`
		c := NewTokenIssuer()
		err := c.FromToken(token, cert, "", "", "", permissions.Restrictions{})
		So(err, ShouldNotBeNil)
		So(err.Error(), ShouldEqual, `unable to compute restrictions: unable to compute authz restrictions from token: invalid restrictions.networks claim type`)
	})

	Convey("Using a token that is not correctly signed", t, func() {
		token := `eyJhbGciOiJFUzI1NiIsInR5cCI6IkpXVCJ9.eyJyZWFsbSI6IlZpbmNlIiwiZGF0YSI6eyJhY2NvdW50IjoiYXBvbXV4IiwiZW1haWwiOiJhZG1pbkBhcG9tdXguY29tIiwiaWQiOiI1ZTFjZjNlZmEzNzAwMzhmYWY3Zjg3NzciLCJvcmdhbml6YXRpb24iOiJhcG9tdXgiLCJyZWFsbSI6InZpbmNlIiwic3ViamVjdCI6ImFwb211eCJ9LCJyZXN0cmljdGlvbnMiOnt9LCJleHAiOjE1OTAzMDQzNDgsImlhdCI6MTU5MDIxNDM0OCwiaXNzIjoiaHR0cHM6Ly9sb2NhbGhvc3Q6NDQ0MyIsInN1YiI6ImFwb211eCJ9.7TZEEG-M-Ed-pKTzEGVZnKKZ1fvG0P7kN-VIKnVn_4TkTR2PX0EaToNZViGgcIs6pYXm7SByzjMl63ZiriSYkg`
		c := NewTokenIssuer()
		err := c.FromToken(token, cert, "", "", "", permissions.Restrictions{})
		So(err, ShouldNotBeNil)
		So(err.Error(), ShouldEqual, `unable to parse input token: unable to parse jwt: crypto/ecdsa: verification error`)
	})

	Convey("Using a token that has no restrictions", t, func() {

		mc := token.NewIdentityToken(token.Source{Type: "mtls"})
		mc.ExpiresAt = time.Now().Add(time.Hour).Unix()
		mc.Issuer = "iss"

		token, _ := mc.JWT(key, time.Time{})
		c := NewTokenIssuer()
		err := c.FromToken(token, cert, "iss", "", "", permissions.Restrictions{})

		So(err, ShouldBeNil)
		So(c.token.Restrictions, ShouldBeNil)
	})

	Convey("Using a token that has bad expiration", t, func() {

		mc := token.NewIdentityToken(token.Source{Type: "mtls"})
		mc.ExpiresAt = time.Now().Add(time.Hour).Unix()
		mc.Issuer = "iss"

		token, _ := mc.JWT(key, time.Time{})
		c := NewTokenIssuer()
		err := c.FromToken(token, cert, "iss", "", "chien", permissions.Restrictions{})

		So(err, ShouldNotBeNil)
		So(err.Error(), ShouldEqual, `unable to compute restrictions: time: invalid duration "chien"`)
	})

	Convey("Using a token that has all valid restrictions", t, func() {

		mc := token.NewIdentityToken(token.Source{Type: "mtls"})
		mc.ExpiresAt = time.Now().Add(time.Hour).Unix()
		mc.Issuer = "iss"
		mc.Restrictions = &permissions.Restrictions{
			Namespace:   "/a",
			Networks:    []string{"1.0.0.0/8"},
			Permissions: []string{"res,get,post"},
		}

		token, _ := mc.JWT(key, time.Time{})
		c := NewTokenIssuer()
		err := c.FromToken(token, cert, "iss", "", "", permissions.Restrictions{
			Namespace:   "/a/b",
			Networks:    []string{"1.1.0.0/16"},
			Permissions: []string{"res,get"},
		})

		So(err, ShouldBeNil)
		So(c.token.Restrictions, ShouldNotBeNil)
		So(c.token.Restrictions.Namespace, ShouldEqual, "/a/b")
		So(c.token.Restrictions.Networks, ShouldResemble, []string{"1.1.0.0/16"})
		So(c.token.Restrictions.Permissions, ShouldResemble, []string{"res,get"})
	})

	Convey("Using a token that has bad ns restrictions", t, func() {

		mc := token.NewIdentityToken(token.Source{Type: "mtls"})
		mc.ExpiresAt = time.Now().Add(time.Hour).Unix()
		mc.Issuer = "iss"
		mc.Restrictions = &permissions.Restrictions{
			Namespace:   "/a",
			Networks:    []string{"1.0.0.0/8"},
			Permissions: []string{"res,get,post"},
		}

		token, _ := mc.JWT(key, time.Time{})
		c := NewTokenIssuer()
		err := c.FromToken(token, cert, "iss", "", "", permissions.Restrictions{
			Namespace:   "/",
			Networks:    []string{"1.1.0.0/16"},
			Permissions: []string{"res,post"},
		})

		So(err, ShouldNotBeNil)
		So(err.Error(), ShouldEqual, `unable to compute restrictions: the new namespace restriction must be empty, '/a' or one of its children`)
	})

	Convey("Using a token that has bad net restrictions", t, func() {

		mc := token.NewIdentityToken(token.Source{Type: "mtls"})
		mc.ExpiresAt = time.Now().Add(time.Hour).Unix()
		mc.Issuer = "iss"
		mc.Restrictions = &permissions.Restrictions{
			Namespace:   "/a",
			Networks:    []string{"1.0.0.0/8"},
			Permissions: []string{"res,get,post"},
		}

		token, _ := mc.JWT(key, time.Time{})
		c := NewTokenIssuer()
		err := c.FromToken(token, cert, "iss", "", "", permissions.Restrictions{
			Namespace:   "/a",
			Networks:    []string{"10.1.0.0/16"},
			Permissions: []string{"res,get"},
		})

		So(err, ShouldNotBeNil)
		So(err.Error(), ShouldEqual, `unable to compute restrictions: the new network restrictions must not overlap any of the original ones`)
	})

	Convey("Using a token that has bad perms restrictions", t, func() {

		mc := token.NewIdentityToken(token.Source{Type: "mtls"})
		mc.ExpiresAt = time.Now().Add(time.Hour).Unix()
		mc.Issuer = "issuer"
		mc.Restrictions = &permissions.Restrictions{
			Namespace:   "/a",
			Networks:    []string{"1.0.0.0/8"},
			Permissions: []string{"@auth:role=enforcer"},
		}

		token, _ := mc.JWT(key, time.Time{})
		c := NewTokenIssuer()
		err := c.FromToken(token, cert, "issuer", "", "", permissions.Restrictions{
			Namespace:   "/a",
			Networks:    []string{"1.1.0.0/16"},
			Permissions: []string{"@auth:role=namespace.administrator"},
		})

		So(err, ShouldNotBeNil)
		So(err.Error(), ShouldEqual, `unable to compute restrictions: the new permissions restrictions must not be broader than the existing ones`)
	})

}

func Test_computeNewValidity(t *testing.T) {

	now := time.Now()
	exp := now.Add(time.Hour)

	type args struct {
		originalExpUNIX      int64
		requestedValidityStr string
	}
	tests := []struct {
		name    string
		args    args
		want    int64
		wantErr bool
	}{
		{
			"no original",
			args{
				0,
				"",
			},
			0,
			true,
		},
		{
			"no requested",
			args{
				exp.Unix(),
				"",
			},
			exp.Unix(),
			false,
		},
		{
			"bad requested",
			args{
				exp.Unix(),
				"chien",
			},
			0,
			true,
		},
		{
			"correct requested",
			args{
				exp.Unix(),
				"30m",
			},
			now.Add(30 * time.Minute).Unix(),
			false,
		},
		{
			"requested too big",
			args{
				exp.Unix(),
				"48h",
			},
			exp.Unix(),
			false,
		},
		{
			"requested the same",
			args{
				exp.Unix(),
				time.Until(exp).String(),
			},
			exp.Unix(),
			false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := computeNewValidity(tt.args.originalExpUNIX, tt.args.requestedValidityStr)
			if (err != nil) != tt.wantErr {
				t.Errorf("computeNewValidity() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if got != tt.want {
				t.Errorf("computeNewValidity() = %v, want %v", got, tt.want)
			}
		})
	}
}
