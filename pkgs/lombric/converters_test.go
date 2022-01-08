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
	"net"
	"reflect"
	"testing"
)

func Test_convertDefaultBool(t *testing.T) {
	type args struct {
		defaultValue []string
	}
	tests := []struct {
		name      string
		args      args
		wantBools []bool
		wantErr   bool
	}{
		{
			"test with valid list of bools",
			args{defaultValue: []string{"true", "TRUE", "True", "false", "FALSE", "False"}},
			[]bool{true, true, true, false, false, false},
			false,
		},
		{
			"test invalid list of bools",
			args{defaultValue: []string{"no", "TRUE", "True", "false", "FALSE", "False"}},
			nil,
			true,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			gotBools, err := convertDefaultBool(tt.args.defaultValue)
			if (err != nil) != tt.wantErr {
				t.Errorf("convertDefaultBool() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if !reflect.DeepEqual(gotBools, tt.wantBools) {
				t.Errorf("convertDefaultBool() = %v, want %v", gotBools, tt.wantBools)
			}
		})
	}
}

func Test_convertDefaultInts(t *testing.T) {
	type args struct {
		defaultValue []string
	}
	tests := []struct {
		name     string
		args     args
		wantInts []int
		wantErr  bool
	}{
		{
			"test with valid ints",
			args{defaultValue: []string{"1", "2", "3"}},
			[]int{1, 2, 3},
			false,
		},
		{
			"test with invalid ints",
			args{defaultValue: []string{"1", "2", "not3"}},
			nil,
			true,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			gotInts, err := convertDefaultInts(tt.args.defaultValue)
			if (err != nil) != tt.wantErr {
				t.Errorf("convertDefaultInts() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if !reflect.DeepEqual(gotInts, tt.wantInts) {
				t.Errorf("convertDefaultInts() = %v, want %v", gotInts, tt.wantInts)
			}
		})
	}
}

func Test_convertDefaultIPs(t *testing.T) {
	type args struct {
		defaultValue []string
	}
	tests := []struct {
		name    string
		args    args
		wantIps []net.IP
		wantErr bool
	}{
		{
			"test with valid ips",
			args{defaultValue: []string{"1.2.3.4", "2.2.2.2"}},
			[]net.IP{net.IPv4(1, 2, 3, 4), net.IPv4(2, 2, 2, 2)},
			false,
		},
		{
			"test with invalid ips",
			args{defaultValue: []string{"not1.2.3.4", "2.2.2.2"}},
			nil,
			true,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			gotIps, err := convertDefaultIPs(tt.args.defaultValue)
			if (err != nil) != tt.wantErr {
				t.Errorf("convertDefaultIPs() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if !reflect.DeepEqual(gotIps, tt.wantIps) {
				t.Errorf("convertDefaultIPs() = %v, want %v", gotIps, tt.wantIps)
			}
		})
	}
}
