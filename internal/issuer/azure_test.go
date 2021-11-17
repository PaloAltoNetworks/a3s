package issuer

import (
	"context"
	"fmt"
	"reflect"
	"testing"

	. "github.com/smartystreets/goconvey/convey"
)

func TestErrAzure(t *testing.T) {
	Convey("ErrAzure should behave correctly ", t, func() {
		e := fmt.Errorf("boom")
		err := ErrAzure{Err: e}
		So(err.Error(), ShouldEqual, "azure error: boom")
		So(err.Unwrap(), ShouldEqual, e)
	})
}

func TestNewAzureIssuer(t *testing.T) {
	Convey("NewAzureIssuer should work", t, func() {
		iss := NewAzureIssuer()
		So(iss.Issue().Source.Type, ShouldEqual, "azure")
	})
}

func TestAzureFromToken(t *testing.T) {
	Convey("Call FromAzureToken should work", t, func() {
		iss := NewAzureIssuer()
		err := iss.FromAzureToken(context.Background(), "not a token")
		So(err, ShouldNotBeNil)
	})
}

func Test_computeAzureClaims(t *testing.T) {
	type args struct {
		token azureJWT
	}
	tests := []struct {
		name string
		args func(t *testing.T) args

		want1 []string
	}{
		{
			"standard",
			func(*testing.T) args {
				return args{
					azureJWT{
						AIO:      "aio",
						AppID:    "appid",
						AppIDAcr: "appidacr",
						IDP:      "idp",
						OID:      "oid",
						RH:       "rh",
						TID:      "tid",
						UTI:      "uti",
						XmsMIRID: "/subscriptions/sub/resourcegroups/grp/providers/prov/7/8",
					},
				}
			},
			[]string{
				"aio=aio",
				"appid=appid",
				"appidacr=appidacr",
				"idp=idp",
				"oid=oid",
				"rh=rh",
				"tid=tid",
				"uti=uti",
				"subscriptions=sub",
				"resourcegroups=grp",
				"providers=prov",
				"providertype=7",
				"identity=8",
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			tArgs := tt.args(t)

			got1 := computeAzureClaims(tArgs.token)

			if !reflect.DeepEqual(got1, tt.want1) {
				t.Errorf("computeAzureClaims got1 = %v, want1: %v", got1, tt.want1)
			}
		})
	}
}
