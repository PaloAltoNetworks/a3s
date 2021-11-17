package issuer

import (
	"fmt"
	"reflect"
	"testing"

	. "github.com/smartystreets/goconvey/convey"
)

func TestErrGCP(t *testing.T) {
	Convey("ErrGCP should behave correctly ", t, func() {
		e := fmt.Errorf("boom")
		err := ErrGCP{Err: e}
		So(err.Error(), ShouldEqual, "gcp error: boom")
		So(err.Unwrap(), ShouldEqual, e)
	})
}

func TestNewGCPIssuer(t *testing.T) {
	Convey("NewGCPIssuer should work", t, func() {
		iss := NewGCPIssuer()
		So(iss.Issue().Source.Type, ShouldEqual, "gcp")
	})
}

func TestGCPFromToken(t *testing.T) {
	Convey("Call FromToken should work", t, func() {
		iss := NewGCPIssuer()
		err := iss.FromToken("not a token", "aud")
		So(err, ShouldNotBeNil)
	})
}
func Test_computeGCPClaims(t *testing.T) {
	type args struct {
		token gcpJWT
	}
	tests := []struct {
		name string
		args func(t *testing.T) args

		want1 []string
	}{
		{
			"standard",
			func(*testing.T) args {
				token := gcpJWT{}
				token.Subject = "sub"
				token.Google.ComputeEngine.ProjectID = "projectid"
				token.Google.ComputeEngine.ProjectNumber = 42
				token.Google.ComputeEngine.Zone = "zone"
				token.Google.ComputeEngine.InstanceID = "instanceid"
				token.Google.ComputeEngine.InstanceName = "name"
				token.Email = "email@email.com"
				token.EmailVerified = true
				return args{
					token: token,
				}
			},
			[]string{
				"subject=sub",
				"projectid=projectid",
				"projectnumber=42",
				"zone=zone",
				"instanceid=instanceid",
				"instancename=name",
				"email=email@email.com",
				"emailverified=true",
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			tArgs := tt.args(t)

			got1 := computeGCPClaims(tArgs.token)

			if !reflect.DeepEqual(got1, tt.want1) {
				t.Errorf("computeGCPClaims got1 = %v, want1: %v", got1, tt.want1)
			}
		})
	}
}
