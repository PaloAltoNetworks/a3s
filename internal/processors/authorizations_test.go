package processors

import (
	"reflect"
	"testing"
)

func Test_flattenTags(t *testing.T) {
	type args struct {
		term [][]string
	}
	tests := []struct {
		name string
		args func(t *testing.T) args

		want1 []string
	}{
		{
			"standard test",
			func(*testing.T) args {
				return args{[][]string{{"a=a", "b=b"}, {"c=c", "a=a"}, {"a=a"}}}
			},
			[]string{"a=a", "b=b", "c=c"},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			tArgs := tt.args(t)

			got1 := flattenTags(tArgs.term)

			if !reflect.DeepEqual(got1, tt.want1) {
				t.Errorf("flattenTags got1 = %v, want1: %v", got1, tt.want1)
			}
		})
	}
}
