package oidcissuer

import (
	"fmt"
	"reflect"
	"testing"

	. "github.com/smartystreets/goconvey/convey"
)

func TestNew(t *testing.T) {

	Convey("Calling New should work ", t, func() {
		iss := New(map[string]interface{}{"hello": "world"})
		So(iss.Issue().Source.Type, ShouldEqual, "oidc")
	})
}

func Test_computeOIDClaims(t *testing.T) {
	type args struct {
		claims map[string]interface{}
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
					map[string]interface{}{
						"@@string": "value",
						"strings":  []string{"v1", "v2"},
						"int":      42,
						"ints":     []int{1, 2},
						"bool":     true,
						"ifaces":   []interface{}{"a", "b"},
						"map":      map[string]interface{}{},
						"float":    42.42,
						"floats":   []float64{1.2, 3.4},
						"error":    fmt.Errorf("yo"),
					},
				}
			},
			[]string{
				"bool=true",
				"error=yo",
				"float=42.420000",
				"floats=1.200000",
				"floats=3.400000",
				"ifaces=a",
				"ifaces=b",
				"int=42",
				"ints=1",
				"ints=2",
				"map=map[]",
				"string=value",
				"strings=v1",
				"strings=v2",
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			tArgs := tt.args(t)

			got1 := computeOIDClaims(tArgs.claims)

			if !reflect.DeepEqual(got1, tt.want1) {
				t.Errorf("computeOIDClaims got1 = %v, want1: %v", got1, tt.want1)
			}
		})
	}
}
