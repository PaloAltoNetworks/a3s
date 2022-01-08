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
	"fmt"
	"net"
	"strconv"
)

func convertDefaultBool(defaultValue []string) (bools []bool, err error) {

	for _, item := range defaultValue {

		switch item {
		case "true", "True", "TRUE":
			bools = append(bools, true)
		case "false", "False", "FALSE":
			bools = append(bools, false)
		default:
			return nil, fmt.Errorf("default value must a bool got: '%s'", item)
		}
	}

	return
}

func convertDefaultInts(defaultValue []string) (ints []int, err error) {

	for _, item := range defaultValue {

		n, err := strconv.Atoi(item)

		if err != nil {
			return nil, fmt.Errorf("default value must be an int. got '%s'", item)
		}

		ints = append(ints, n)
	}

	return
}

func convertDefaultIPs(defaultValue []string) (ips []net.IP, err error) {

	for _, item := range defaultValue {

		ip := net.ParseIP(item)
		if ip == nil {
			return nil, fmt.Errorf("default value must be an IP. got '%s'", item)
		}

		ips = append(ips, ip)
	}

	return
}
