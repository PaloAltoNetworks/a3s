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
	"errors"
	"fmt"
	"os"
	"reflect"

	"github.com/spf13/pflag"
	"github.com/spf13/viper"
)

func checkRequired(keys ...string) error {

	var failed bool
	for _, key := range keys {
		if !viper.IsSet(key) || reflect.DeepEqual(viper.Get(key), reflect.Zero(reflect.TypeOf(viper.Get(key))).Interface()) {
			fmt.Printf("Error: Parameter '--%s' is required.\n", key)
			failed = true
		}
	}

	if failed {
		return errors.New("missing required parameter")
	}

	return nil
}

func checkAllowedValues(allowedValues map[string][]string) error {

	var failed bool
	for key, values := range allowedValues {

		if !stringInSlice(viper.GetString(key), values) {
			fmt.Printf("Error: Parameter '--%s' must be one of %s. '%s' is invalid.\n", key, values, viper.GetString(key))
			failed = true
		}
	}

	if failed {
		return errors.New("wrong allowed values")
	}

	return nil
}

var testMode bool

func fail() {
	fmt.Println()
	pflag.Usage()

	if testMode {
		os.Exit(0)
	}

	os.Exit(1)
}

func stringInSlice(str string, list []string) bool {

	for _, s := range list {
		if s == str {
			return true
		}
	}

	return false
}
