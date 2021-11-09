package issuer

import (
	"fmt"
	"reflect"
	"testing"

	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/aws/arn"
	"github.com/aws/aws-sdk-go/service/sts"
	. "github.com/smartystreets/goconvey/convey"
)

func TestErrAWSSTS(t *testing.T) {
	Convey("ErrAWSSTS should work", t, func() {
		e := fmt.Errorf("boom")
		err := ErrAWSSTS{Err: e}
		So(err.Error(), ShouldEqual, "aws error: boom")
		So(err.Unwrap(), ShouldEqual, e)
	})
}

func TestNewAWSSTSIssuer(t *testing.T) {
	Convey("Calling NewAWSSTSIssuer should work", t, func() {
		iss := NewAWSSTSIssuer()
		So(iss.token, ShouldNotBeNil)
		So(iss.token.Source.Type, ShouldEqual, "awssts")
		So(iss.Issue(), ShouldEqual, iss.token)
	})
}

func TestFromAWSSTS(t *testing.T) {
	Convey("Given an AWSSTSIssuer", t, func() {
		iss := NewAWSSTSIssuer()
		err := iss.FromSTS("id", "secret", "token")
		So(err, ShouldNotBeNil)
		So(err.Error(), ShouldStartWith, "aws error: unable to retrieve aws identity: InvalidClientTokenId: The security token included in the request is invalid.")
	})
}

func Test_computeSTSClaims(t *testing.T) {
	type args struct {
		out  *sts.GetCallerIdentityOutput
		parn arn.ARN
	}
	tests := []struct {
		name string
		args func(t *testing.T) args

		want1 []string
	}{
		{
			"standard test",
			func(*testing.T) args {
				return args{
					&sts.GetCallerIdentityOutput{
						Account: aws.String("account"),
						Arn:     aws.String("arn"),
						UserId:  aws.String("userid"),
					},
					arn.ARN{
						Partition: "partition",
						Service:   "service",
						Resource:  "a/b/c",
					},
				}
			},
			[]string{
				"account=account",
				"arn=arn",
				"userid=userid",
				"partition=partition",
				"service=service",
				"resource=a/b/c",
				"resourcetype=a",
				"rolename=b",
				"rolesessionname=c",
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			tArgs := tt.args(t)

			got1 := computeSTSClaims(tArgs.out, tArgs.parn)

			if !reflect.DeepEqual(got1, tt.want1) {
				t.Errorf("computeSTSClaims got1 = %v, want1: %v", got1, tt.want1)
			}
		})
	}
}
