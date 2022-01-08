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
	"bytes"
	"fmt"
	"io/ioutil"
	"net/url"
	"os"
	"reflect"
	"strconv"
	"strings"
	"time"

	"github.com/spf13/pflag"
	"github.com/spf13/viper"
)

const enabledKey = "true"

// Configurable is the interface of a configuration.
type Configurable interface {
}

// EnvPrexixer is the interface to implement in order to
// support arguments from env.
type EnvPrexixer interface {
	Prefix() string
}

// VersionPrinter is an extension to Configurable that can print its version.
type VersionPrinter interface {
	PrintVersion()
}

// Initialize does all the basic job of bindings
func Initialize(conf Configurable) {

	requiredFlags, secretFlags, allowedValues, filesValues := installFlags(conf)

	pflag.VisitAll(func(f *pflag.Flag) {

		var v interface{}
		var err error

		switch f.Value.Type() {

		case "stringSlice":
			v, err = pflag.CommandLine.GetStringSlice(f.Name)

		case "boolSlice":
			v, err = pflag.CommandLine.GetBoolSlice(f.Name)

		case "intSlice":
			v, err = pflag.CommandLine.GetIntSlice(f.Name)

		case "ipSlice":
			v, err = pflag.CommandLine.GetIPSlice(f.Name)
		}

		if err != nil {
			panic("Unable to trick viper with the defaults: %s" + err.Error())
		}

		viper.SetDefault(f.Name, v)
	})

	if err := viper.BindPFlags(pflag.CommandLine); err != nil {
		panic("Unable to bind flags: " + err.Error())
	}

	if p, ok := conf.(EnvPrexixer); ok {
		viper.SetEnvPrefix(p.Prefix())
		viper.SetEnvKeyReplacer(strings.NewReplacer("-", "_"))
		viper.AutomaticEnv()
		viper.SetTypeByDefaultValue(true)
	}

	if vp, ok := conf.(VersionPrinter); ok && viper.GetBool("version") {
		vp.PrintVersion()
		os.Exit(0)
	}

	if err := checkRequired(requiredFlags...); err != nil {
		fail()
	}

	if err := checkAllowedValues(allowedValues); err != nil {
		fail()
	}

	// Replace secret from content of files if needed.
	if _, ok := conf.(EnvPrexixer); ok {

		for _, key := range filesValues {

			value := viper.GetString(key)

			if strings.HasPrefix(value, "file://") {

				u, err := url.Parse(value)
				if err != nil {
					panic(fmt.Sprintf("invalid url for secret: %s", err))
				}

				data, err := ioutil.ReadFile(u.Path)
				if err != nil {
					panic(fmt.Sprintf("unable to read secret file for key '%s': %s", key, err))
				}
				viper.Set(key, string(bytes.TrimSpace(data)))

				if u.Query().Get("delete") != "" {
					if err := os.Remove(u.Path); err != nil {
						panic(fmt.Sprintf("unable to delete secret file: %s", err))
					}
				}
			}
		}
	}

	if err := viper.Unmarshal(conf); err != nil {
		panic("Unable to unmarshal configuration: " + err.Error())
	}

	// Clean up all secrets
	if p, ok := conf.(EnvPrexixer); ok {
		for _, key := range secretFlags {
			env := strings.Replace(strings.ToUpper(p.Prefix()+"_"+key), "-", "_", -1)
			if err := os.Unsetenv(env); err != nil {
				panic("Unable to unset secret env variable " + env)
			}
		}
	}
}

func deepFields(ift reflect.Type) ([]reflect.StructField, []string) {

	fields := make([]reflect.StructField, 0)
	overrides := []string{}

	for i := 0; i < ift.NumField(); i++ {

		field := ift.Field(i)

		switch field.Type.Kind() {

		case reflect.Struct:

			if overrideString := field.Tag.Get("override"); overrideString != "" {
				overrides = append(overrides, overrideString)
			}

			f, o := deepFields(field.Type)
			overrides = append(overrides, o...)
			fields = append(fields, f...)

		default:
			fields = append(fields, field)
		}
	}

	return fields, overrides
}

func installFlags(conf Configurable) (requiredFlags []string, secretFlags []string, allowedValues map[string][]string, fromFilesFlags []string) {

	t := reflect.ValueOf(conf).Elem().Type()

	fields, overrides := deepFields(t)
	defaultOverrides := map[string]string{}
	allowedValues = map[string][]string{}
	hiddenFlags := []string{}

	for _, raw := range overrides {

		for _, innerOverride := range strings.Split(raw, ",") {

			parts := strings.SplitN(innerOverride, "=", 2)
			defaultOverrides[parts[0]] = parts[1]
		}
	}

	for _, field := range fields {

		key := field.Tag.Get("mapstructure")
		if key == "" || key == "-" {
			continue
		}

		description := field.Tag.Get("desc")

		def := field.Tag.Get("default")
		if o, ok := defaultOverrides[key]; ok {
			if o == "-" {
				continue
			}
			def = o
		}

		if field.Tag.Get("secret") == enabledKey {
			secretFlags = append(secretFlags, key)
		}

		if field.Tag.Get("file") == enabledKey {
			fromFilesFlags = append(fromFilesFlags, key)
		}

		if field.Tag.Get("required") == enabledKey {
			requiredFlags = append(requiredFlags, key)
			description += " [required]"
		}

		if field.Tag.Get("hidden") == enabledKey {
			hiddenFlags = append(hiddenFlags, key)
		}

		if field.Type.Kind() != reflect.Slice {

			switch field.Type.Name() {

			case "bool":
				pflag.Bool(key, def == enabledKey, description)

			case "string":
				if allowed := field.Tag.Get("allowed"); allowed != "" {
					allowedValues[key] = strings.Split(allowed, ",")
					description += fmt.Sprintf(" [allowed: %s]", allowed)
				}
				pflag.String(key, def, description)

			case "Duration":
				if def == "" {
					pflag.Duration(key, 0, description)
					break
				}
				d, err := time.ParseDuration(def)
				if err != nil {
					panic("Unable to parse duration from: " + def)
				}
				pflag.Duration(key, d, description)

			case "int":
				if def == "" {
					pflag.Int(key, 0, description)
					break
				}
				d, err := strconv.Atoi(def)
				if err != nil {
					panic("Unable to parse int from: " + def)
				}
				pflag.Int(key, d, description)

			default:
				panic("Unsupported type: " + field.Type.Name())
			}

		} else {

			defaultValues := strings.Split(def, ",")

			switch field.Type.Elem().Name() {

			case "string":

				pflag.StringSlice(key, defaultValues, description)

			case "bool":
				def, err := convertDefaultBool(defaultValues)
				if err != nil {
					panic(err)
				}
				pflag.BoolSlice(key, def, description)

			case "int":

				def, err := convertDefaultInts(defaultValues)
				if err != nil {
					panic(err)
				}
				pflag.IntSlice(key, def, description)

			case "IP":

				def, err := convertDefaultIPs(defaultValues)
				if err != nil {
					panic(err)
				}
				pflag.IPSlice(key, def, description)

			default:
				panic("Unsupported type: " + field.Type.Elem().Name())
			}
		}
	}

	if _, ok := conf.(VersionPrinter); ok {
		pflag.BoolP("version", "v", false, "Display the version")
	}

	for _, hiddenFlag := range hiddenFlags {
		if err := pflag.CommandLine.MarkHidden(hiddenFlag); err != nil {
			panic("Unable to hide flags: " + err.Error())
		}
	}

	pflag.Parse()

	return requiredFlags, secretFlags, allowedValues, fromFilesFlags
}
