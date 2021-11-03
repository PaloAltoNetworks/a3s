package sharder

import "testing"

func Test_hash(t *testing.T) {
	type args struct {
		v string
	}
	tests := []struct {
		name string
		args args
		want int
	}{
		{
			"test",
			args{
				"hello",
			},
			5465302536158026498,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := hash(tt.args.v); got != tt.want {
				t.Errorf("hash() = %v, want %v", got, tt.want)
			}
		})
	}
}
