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
	"io/ioutil"
	"net"
	"os"
	"testing"
	"time"

	. "github.com/smartystreets/goconvey/convey"
	"github.com/spf13/viper"
)

func init() {
	testMode = true
}

type testConf struct {
	ABool                   bool          `mapstructure:"a-bool"                    desc:"This is a boolean"            default:"true"`
	ARequiredBool           bool          `mapstructure:"a-required-bool"           desc:"This is a boolean"            required:"true"`
	ABoolNoDef              bool          `mapstructure:"a-bool-nodef"              desc:"This is a no def boolean"     `
	ABoolSlice              []bool        `mapstructure:"a-bool-slice"              desc:"This is a bool slice"         default:"true,false,true"`
	ADuration               time.Duration `mapstructure:"a-duration"                desc:"This is a duration"           default:"10s"`
	ADurationNoDef          time.Duration `mapstructure:"a-duration-nodef"          desc:"This is a no def duration"    `
	AInteger                int           `mapstructure:"a-integer"                 desc:"This is a number"             default:"42"`
	AIntegerNoDef           int           `mapstructure:"a-integer-nodef"           desc:"This is a no def number"      `
	AIntSlice               []int         `mapstructure:"a-int-slice"               desc:"This is a int slice"          default:"1,2,3"`
	AnEnum                  string        `mapstructure:"a-enum"                    desc:"This is an enum"              allowed:"a,b,c" default:"a"`
	AnIPSlice               []net.IP      `mapstructure:"a-ip-slice"                desc:"This is an ip slice"          default:"127.0.0.1,192.168.100.1"`
	AnotherStringSliceNoDef []string      `mapstructure:"a-string-slice-from-var"   desc:"This is a no def string"      `
	ASecret                 string        `mapstructure:"a-secret-from-var"         desc:"This is a secret"             secret:"true"`
	AString                 string        `mapstructure:"a-string"                  desc:"This is a string"             default:"hello"`
	AStringNoDef            string        `mapstructure:"a-string-nodef"            desc:"This is a no def string"      `
	AStringSlice            []string      `mapstructure:"a-string-slice"            desc:"This is a string slice"       default:"a,b,c"`
	AStringSliceNoDef       []string      `mapstructure:"a-string-slice-nodef"      desc:"This is a no def string slice"`
	ASimpleFromFile         string        `mapstructure:"a-simple-from-file"        desc:"This is a simple from file"   file:"true"`
	ASimpleFromFileDelete   string        `mapstructure:"a-simple-from-file-del"    desc:"This is a simple from file"   file:"true"`

	embedTestConf `mapstructure:",squash" override:"embedded-string-a=outter1,embedded-ignored-string=-"`
}

type embedTestConf struct {
	EmbeddedStringA        string `mapstructure:"embedded-string-a"        desc:"This is a string"       default:"inner1"`
	EmbeddedStringB        string `mapstructure:"embedded-string-b"        desc:"This is a string"       default:"inner2"`
	EmbeddedIgnoredStringB string `mapstructure:"embedded-ignored-string"  desc:"This is a string"       default:"inner3"`
}

// Prefix return the configuration prefix.
func (c *testConf) Prefix() string { return "lombric" }
func (c *testConf) PrintVersion()  {}

func TestLombric_Initialize(t *testing.T) {

	Convey("Given have a conf", t, func() {

		sfile1, err := ioutil.TempFile(os.TempDir(), "secret")
		if err != nil {
			panic(err)
		}
		defer sfile1.Close() // nolint
		if _, err := sfile1.WriteString("this-is-super=s3cr3t\n\n"); err != nil {
			panic(err)
		}
		spath1 := fmt.Sprintf("file://%s", sfile1.Name())

		sfile2, err := ioutil.TempFile(os.TempDir(), "secret2")
		if err != nil {
			panic(err)
		}
		defer sfile2.Close() // nolint
		if _, err := sfile2.WriteString("42\n\n"); err != nil {
			panic(err)
		}
		spath2 := fmt.Sprintf("file://%s?delete=true", sfile2.Name())

		conf := &testConf{}
		os.Setenv("LOMBRIC_A_STRING_SLICE_FROM_VAR", "x y z") // nolint: errcheck
		os.Setenv("LOMBRIC_A_REQUIRED_BOOL", "true")          // nolint: errcheck
		os.Setenv("LOMBRIC_A_SECRET_FROM_VAR", "secret")      // nolint: errcheck
		os.Setenv("LOMBRIC_A_SIMPLE_FROM_FILE", spath1)       // nolint: errcheck
		os.Setenv("LOMBRIC_A_SIMPLE_FROM_FILE_DEL", spath2)   // nolint: errcheck

		Initialize(conf)

		Convey("Then the flags should be correctly set", func() {

			So(conf.ABool, ShouldEqual, true)
			So(conf.ARequiredBool, ShouldEqual, true)
			So(conf.ABoolNoDef, ShouldEqual, false)
			So(conf.ABoolSlice, ShouldResemble, []bool{true, false, true})
			So(conf.ADuration, ShouldEqual, 10*time.Second)
			So(conf.ADurationNoDef, ShouldEqual, 0)
			So(conf.AInteger, ShouldEqual, 42)
			So(conf.AIntegerNoDef, ShouldEqual, 0)
			So(conf.AIntSlice, ShouldResemble, []int{1, 2, 3})
			So(conf.AnIPSlice, ShouldResemble, []net.IP{net.IPv4(127, 0, 0, 1), net.IPv4(192, 168, 100, 1)})
			So(conf.AnotherStringSliceNoDef, ShouldResemble, []string{"x", "y", "z"})
			So(conf.ASecret, ShouldEqual, "secret")
			So(conf.ASimpleFromFile, ShouldEqual, "this-is-super=s3cr3t")
			So(conf.ASimpleFromFileDelete, ShouldEqual, "42")
			So(viper.GetString("a-simple-from-file"), ShouldEqual, "this-is-super=s3cr3t")
			So(viper.GetInt("a-simple-from-file-del"), ShouldEqual, 42)
			So(conf.AString, ShouldEqual, "hello")
			So(conf.AStringNoDef, ShouldEqual, "")
			So(conf.AStringSlice, ShouldResemble, []string{"a", "b", "c"})
			So(conf.AStringSliceNoDef, ShouldResemble, []string(nil))
			So(conf.EmbeddedIgnoredStringB, ShouldEqual, "")
			So(conf.EmbeddedStringA, ShouldEqual, "outter1")
			So(conf.EmbeddedStringB, ShouldEqual, "inner2")
			So(os.Getenv("LOMBRIC_A_SECRET_FROM_VAR"), ShouldEqual, "")
			So(viper.AllKeys(), ShouldNotContain, "embedded-ignored-string")

			_, err := os.Stat(sfile1.Name())
			So(os.IsNotExist(err), ShouldBeFalse)

			_, err = os.Stat(sfile2.Name())
			So(os.IsNotExist(err), ShouldBeTrue)
		})
	})
}

func TestBadDefaults(t *testing.T) {

	Convey("Given I have struct with bad default duration", t, func() {

		c := &struct {
			A time.Duration `mapstructure:"BadDefaultDuration" desc:"" default:"toto"`
		}{}

		Convey("Then calling Initialize should panic", func() {
			So(func() { Initialize(c) }, ShouldPanicWith, "Unable to parse duration from: toto")
		})
	})

	Convey("Given I have struct with bad default int", t, func() {

		c := &struct {
			A int `mapstructure:"badDefaultInt" desc:"" default:"toto"`
		}{}

		Convey("Then calling Initialize should panic", func() {
			So(func() { Initialize(c) }, ShouldPanicWith, "Unable to parse int from: toto")
		})
	})

	Convey("Given I have struct with unsuported type", t, func() {

		c := &struct {
			A float64 `mapstructure:"badFloat" desc:"" default:"toto"`
		}{}

		Convey("Then calling Initialize should panic", func() {
			So(func() { Initialize(c) }, ShouldPanicWith, "Unsupported type: float64")
		})
	})

	Convey("Given I have struct with bad default bool slice", t, func() {

		c := &struct {
			A []bool `mapstructure:"badBools" desc:"" default:"a,b,c"`
		}{}

		Convey("Then calling Initialize should panic", func() {
			So(func() { Initialize(c) }, ShouldPanic)
			// So(func() { Initialize(c) }, ShouldPanicWith, `default value must a bool got: 'a'`)
		})
	})

	Convey("Given I have struct with bad default int slice", t, func() {

		c := &struct {
			A []int `mapstructure:"badInts" desc:"" default:"a,b,c"`
		}{}

		Convey("Then calling Initialize should panic", func() {
			So(func() { Initialize(c) }, ShouldPanic)
			// So(func() { Initialize(c) }, ShouldPanicWith, "default value must be an int. got 'a'")
		})
	})

	Convey("Given I have struct with bad default int slice", t, func() {

		c := &struct {
			A []net.IP `mapstructure:"badIPs" desc:"" default:"a,b,c"`
		}{}

		Convey("Then calling Initialize should panic", func() {
			So(func() { Initialize(c) }, ShouldPanic)
			// So(func() { Initialize(c) }, ShouldPanicWith, "default value must be an int. got 'a'")
		})
	})

	Convey("Given I have struct with unsuported slice", t, func() {

		c := &struct {
			A []float64 `mapstructure:"badFloats" desc:"" default:"a,b,c"`
		}{}

		Convey("Then calling Initialize should panic", func() {
			So(func() { Initialize(c) }, ShouldPanicWith, "Unsupported type: float64")
		})
	})
}
