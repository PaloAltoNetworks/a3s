package importing

import (
	"encoding/json"
	"reflect"
	"testing"

	"go.aporeto.io/a3s/pkgs/api"
	"go.aporeto.io/elemental"
)

func TestHash(t *testing.T) {
	type args struct {
		obj     Importable
		manager elemental.ModelManager
	}
	tests := []struct {
		name string
		args func(t *testing.T) args

		want1      string
		wantErr    bool
		inspectErr func(err error, t *testing.T) //use for more precise error evaluation after test
	}{
		{
			"nil obj",
			func(*testing.T) args {
				return args{
					nil,
					nil,
				}
			},
			"",
			true,
			nil,
		},
		{
			"nil manager",
			func(*testing.T) args {
				return args{
					api.NewHTTPSource(),
					nil,
				}
			},
			"",
			true,
			nil,
		},
		{
			"basic",
			func(*testing.T) args {
				return args{
					api.NewHTTPSource(),
					api.Manager(),
				}
			},
			"2d30e3b0c44298fc1c149afbf4c8996fb92427ae41e4649b934ca495991b7852b855",
			false,
			nil,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			tArgs := tt.args(t)

			got1, err := Hash(tArgs.obj, tArgs.manager)

			if !reflect.DeepEqual(got1, tt.want1) {
				t.Errorf("Hash got1 = %v, want1: %v", got1, tt.want1)
			}

			if (err != nil) != tt.wantErr {
				t.Fatalf("Hash error = %v, wantErr: %t", err, tt.wantErr)
			}

			if tt.inspectErr != nil {
				tt.inspectErr(err, t)
			}
		})
	}
}

func Test_sanitize(t *testing.T) {
	type args struct {
		obj     Importable
		manager elemental.ModelManager
	}
	tests := []struct {
		name string
		args func(t *testing.T) args

		want1      map[string]interface{}
		wantErr    bool
		inspectErr func(err error, t *testing.T) // use for more precise error evaluation after test
	}{
		{
			"basic without nested object",
			func(*testing.T) args {
				obj := api.NewHTTPSource()
				obj.Name = "name"
				obj.CA = "ca"
				obj.ImportHash = "should be removed"
				obj.ImportLabel = "should be removed"
				obj.Namespace = "should be removed because it's autogen"
				return args{
					obj,
					api.Manager(),
				}
			},
			map[string]interface{}{
				"name": "name",
				"CA":   "ca",
			},
			false,
			nil,
		},
		{
			"with zero value nested object",
			func(*testing.T) args {
				obj := api.NewHTTPSource()
				obj.Name = "name"
				obj.CA = "ca"
				obj.ImportHash = "should be removed"
				obj.ImportLabel = "should be removed"
				obj.Namespace = "should be removed because it's autogen"
				obj.Modifier = api.NewIdentityModifier()
				return args{
					obj,
					api.Manager(),
				}
			},
			map[string]interface{}{
				"name": "name",
				"CA":   "ca",
			},
			false,
			nil,
		},
		{
			"with non zero nested object",
			func(*testing.T) args {
				obj := api.NewHTTPSource()
				obj.Name = "name"
				obj.CA = "ca"
				obj.ImportHash = "should be removed"
				obj.ImportLabel = "should be removed"
				obj.Namespace = "should be removed because it's autogen"
				obj.Modifier = api.NewIdentityModifier()
				obj.Modifier.Certificate = "cert"
				return args{
					obj,
					api.Manager(),
				}
			},
			map[string]interface{}{
				"name": "name",
				"CA":   "ca",
				"modifier": map[string]interface{}{
					"certificate": "cert",
				},
			},
			false,
			nil,
		},
		{
			"with default enum",
			func(*testing.T) args {
				obj := api.NewLDAPSource()
				obj.Name = "name"
				obj.CA = "ca"
				obj.ImportHash = "should be removed"
				obj.ImportLabel = "should be removed"
				obj.Namespace = "should be removed because it's autogen"
				obj.SecurityProtocol = api.LDAPSourceSecurityProtocolInbandTLS
				return args{
					obj,
					api.Manager(),
				}
			},
			map[string]interface{}{
				"name": "name",
				"CA":   "ca",
			},
			false,
			nil,
		},
		{
			"with non default enum",
			func(*testing.T) args {
				obj := api.NewLDAPSource()
				obj.Name = "name"
				obj.CA = "ca"
				obj.ImportHash = "should be removed"
				obj.ImportLabel = "should be removed"
				obj.Namespace = "should be removed because it's autogen"
				obj.SecurityProtocol = api.LDAPSourceSecurityProtocolTLS
				return args{
					obj,
					api.Manager(),
				}
			},
			map[string]interface{}{
				"name":             "name",
				"CA":               "ca",
				"securityProtocol": api.LDAPSourceSecurityProtocolTLS,
			},
			false,
			nil,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			tArgs := tt.args(t)

			got1, err := sanitize(tArgs.obj, tArgs.manager)

			if !reflect.DeepEqual(got1, tt.want1) {
				t.Errorf("sanitize got1 = %v, want1: %v", got1, tt.want1)
			}

			if (err != nil) != tt.wantErr {
				t.Fatalf("sanitize error = %v, wantErr: %t", err, tt.wantErr)
			}

			if tt.inspectErr != nil {
				tt.inspectErr(err, t)
			}
		})
	}
}

func Test_cleanIrrelevantValues(t *testing.T) {
	type args struct {
		data     map[string]interface{}
		template map[string]interface{}
	}
	tests := []struct {
		name string
		args func(t *testing.T) args

		want1 map[string]interface{}
	}{
		{
			"empty",
			func(*testing.T) args {
				return args{
					map[string]interface{}{},
					map[string]interface{}{},
				}
			},
			map[string]interface{}{},
		},
		{
			"basic",
			func(*testing.T) args {
				return args{
					map[string]interface{}{
						"zero-string":          "",
						"default-string":       "default",
						"string":               "string",
						"other-string":         "other-string",
						"zero-string-array":    nil,
						"default-string-array": []string{"default"},
						"string-array":         []string{"string"},
						"other-string-array":   []string{"other-string"},
						"not-matching-type":    "a",
						"sub": map[string]interface{}{
							"zero-string":          "",
							"default-string":       "default",
							"string":               "string",
							"other-string":         "other-string",
							"zero-string-array":    nil,
							"default-string-array": []string{"default"},
							"string-array":         []string{"string"},
							"other-string-array":   []string{"other-string"},
							"not-matching-type":    "a",
						},
						"not-matching-sub": map[string]interface{}{"a": "a"},
						"equal-sub":        map[string]interface{}{"a": "a"},
					},
					map[string]interface{}{
						"default-string":       "default",
						"string":               "not-string",
						"default-string-array": []string{"default"},
						"string-array":         []string{"not-string"},
						"not-matching-type":    1,
						"sub": map[string]interface{}{
							"default-string":       "default",
							"string":               "not-string",
							"default-string-array": []string{"default"},
							"string-array":         []string{"not-string"},
							"not-matching-type":    1,
						},
						"not-matching-sub": "a",
						"equal-sub":        map[string]interface{}{"a": "a"},
					},
				}
			},
			map[string]interface{}{
				"string":             "string",
				"other-string":       "other-string",
				"string-array":       []string{"string"},
				"other-string-array": []string{"other-string"},
				"not-matching-type":  "a",
				"sub": map[string]interface{}{
					"not-matching-type":  "a",
					"string":             "string",
					"other-string":       "other-string",
					"string-array":       []string{"string"},
					"other-string-array": []string{"other-string"},
				},
				"not-matching-sub": map[string]interface{}{"a": "a"},
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			tArgs := tt.args(t)

			got1 := cleanIrrelevantValues(tArgs.data, tArgs.template)

			if !reflect.DeepEqual(got1, tt.want1) {
				a, _ := json.MarshalIndent(got1, "", "  ")
				b, _ := json.MarshalIndent(tt.want1, "", "  ")
				t.Errorf("clean got1 = %s\nwant1: %s", string(a), string(b))
			}
		})
	}
}

func Test_hash(t *testing.T) {
	type args struct {
		data map[string]interface{}
		ns   string
	}
	tests := []struct {
		name string
		args func(t *testing.T) args

		want1      string
		wantErr    bool
		inspectErr func(err error, t *testing.T) //use for more precise error evaluation after test
	}{
		{
			"nil",
			func(*testing.T) args {
				return args{
					nil,
					"",
				}
			},
			"2d30e3b0c44298fc1c149afbf4c8996fb92427ae41e4649b934ca495991b7852b855",
			false,
			nil,
		},
		{
			"basic map",
			func(*testing.T) args {
				return args{
					map[string]interface{}{
						"a": true,
						"b": 1,
						"c": "c",
						"d": []string{"a", "b"},
						"e": []interface{}{"a", "b"},
						"f": map[string]interface{}{"a": "b"},
					},
					"",
				}
			},
			"2d36353939373434343439343034313732353632e3b0c44298fc1c149afbf4c8996fb92427ae41e4649b934ca495991b7852b855",
			false,
			nil,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			tArgs := tt.args(t)

			got1, err := hash(tArgs.data, tArgs.ns)

			if !reflect.DeepEqual(got1, tt.want1) {
				t.Errorf("hash got1 = %v, want1: %v", got1, tt.want1)
			}

			if (err != nil) != tt.wantErr {
				t.Fatalf("hash error = %v, wantErr: %t", err, tt.wantErr)
			}

			if tt.inspectErr != nil {
				tt.inspectErr(err, t)
			}
		})
	}
}
