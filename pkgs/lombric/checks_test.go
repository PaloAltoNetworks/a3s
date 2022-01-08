// Copyright 2019 Aporeto Inc.
// Licensed under the Apache License, Version 2.0 (the "License");
// you may not use this file except in compliance with the License.
// You may obtain a copy of the License at
//     http://www.apache.org/licenses/LICENSE-2.0
// Unless required by applicable law or agreed to in writing, software
// distributed under the License is distributed on an "AS IS" BASIS,
// WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
// See the License for the specific language governing permissions and
// limitations under the License.

package lombric

import (
	"testing"
)

func Test_stringInSlice(t *testing.T) {
	type args struct {
		str  string
		list []string
	}
	tests := []struct {
		name string
		args args
		want bool
	}{
		{
			"test string in slice",
			args{
				"a",
				[]string{"a", "b", "c"},
			},
			true,
		},
		{
			"test string not in slice",
			args{
				"z",
				[]string{"a", "b", "c"},
			},
			false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := stringInSlice(tt.args.str, tt.args.list); got != tt.want {
				t.Errorf("stringInSlice() = %v, want %v", got, tt.want)
			}
		})
	}
}

func Test_checkRequired(t *testing.T) {
	type args struct {
		keys []string
	}
	tests := []struct {
		name    string
		args    args
		wantErr bool
	}{
		{
			"test failure",
			args{
				[]string{"a", "b"},
			},
			true,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if err := checkRequired(tt.args.keys...); (err != nil) != tt.wantErr {
				t.Errorf("checkRequired() error = %v, wantErr %v", err, tt.wantErr)
			}
		})
	}
}

func Test_checkAllowedValues(t *testing.T) {
	type args struct {
		allowedValues map[string][]string
	}
	tests := []struct {
		name    string
		args    args
		wantErr bool
	}{
		{
			"test failure",
			args{
				map[string][]string{"a": {"1", "2"}},
			},
			true,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if err := checkAllowedValues(tt.args.allowedValues); (err != nil) != tt.wantErr {
				t.Errorf("checkAllowedValues() error = %v, wantErr %v", err, tt.wantErr)
			}
		})
	}
}
