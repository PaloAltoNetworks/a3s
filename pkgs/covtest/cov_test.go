package cov

import (
	"testing"
)

func Test_f(t *testing.T) {
	tests := []struct {
		name string
	}{
		{
			"cool",
		},
		//TODO: Add test cases
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			f()

		})
	}
}
