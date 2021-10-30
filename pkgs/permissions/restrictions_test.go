package permissions

import (
	"reflect"
	"testing"

	. "github.com/smartystreets/goconvey/convey"
)

func TestGetRestrictions(t *testing.T) {
	type args struct {
		token string
	}
	tests := []struct {
		name    string
		args    args
		want    Restrictions
		wantErr bool
	}{
		{
			"token with restrictions",
			args{
				`eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9.eyJyZWFsbSI6IlZpbmNlIiwicmVzdHJpY3Rpb25zIjp7InBlcm1zIjpbInIxLGdldCxwb3N0Il0sIm5hbWVzcGFjZSI6Ii9hcG9tdXgvY2hpbGQiLCJuZXR3b3JrcyI6WyIxMjcuMC4wLjEvMzIiXX0sImV4cCI6MTU5MDA0Mjk5OCwiaWF0IjoxNTg5OTUyOTk4LCJpc3MiOiJodHRwczovL2xvY2FsaG9zdDo0NDQzIiwic3ViIjoiYXBvbXV4In0.8q9wEwRAj2JHqGUhrlKrkymf_xF6rIQkvKXu4YcyI-Q`,
			},
			Restrictions{Namespace: "/apomux/child", Permissions: []string{"r1,get,post"}, Networks: []string{"127.0.0.1/32"}},
			false,
		},
		{
			"token with no restriction",
			args{
				`eyJhbGciOiJFUzI1NiIsInR5cCI6IkpXVCJ9.eyJyZWFsbSI6IlZpbmNlIiwiZGF0YSI6eyJhY2NvdW50IjoiYXBvbXV4IiwiZW1haWwiOiJhZG1pbkBhcG9tdXguY29tIiwiaWQiOiI1ZTFjZjNlZmEzNzAwMzhmYWY3Zjg3NzciLCJvcmdhbml6YXRpb24iOiJhcG9tdXgiLCJyZWFsbSI6InZpbmNlIiwic3ViamVjdCI6ImFwb211eCJ9LCJleHAiOjE1OTAwMTUzNTIsImlhdCI6MTU4OTkyNTM1MiwiaXNzIjoiaHR0cHM6Ly9sb2NhbGhvc3Q6NDQ0MyIsInN1YiI6ImFwb211eCJ9.agqImtfkfjJugJH59XfQwkasIayYtvG6tz3p84jMulfbgwZzTLzgfRDLNIcfnfqfUix_702BUJxvdlsaSsgeUg`,
			},
			Restrictions{Namespace: "", Permissions: nil},
			false,
		},
		{
			"invalid token",
			args{
				`eyJhbGciOiJFUzI1NiIsInR5cCI6IkpXVCJ9.eyJyZWFsbSI6IlZpbmNlIiwiZGF0YSI6eyJ1pbkBhcG9tdXguY29tIiwiaWQiOiI1ZTFjZjNlZmEzNzAwMzhmYWY3Zjg3NzciLCJvcmdhbml6YXRpb24iOiJhcG9tdXgiLCJyZWFsbSI6InZpbmNlIiwic3ViamVjdCI6ImFwb211eCJ9LCJleHAiOjE1OTAwMTUzNTIsImlhdCI6MTU4OTkyNTM1MiwiaXNzIjoiaHR0cHM6Ly9sb2NhbGhvc3Q6NDQ0MyIsInN1YiI6ImFwb211eCJ9.agqImtfkfjJugJH59XfQwkasIayYtvG6tz3p84jMulfbgwZzTLzgfRDLNIcfnfqfUix_702BUJxvdlsaSsgeUg`,
			},
			Restrictions{Namespace: "", Permissions: nil},
			true,
		},
		{
			"token with invalid namespace type",
			args{
				`eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9.eyJyZWFsbSI6IlZpbmNlIiwiZGF0YSI6eyJhY2NvdW50IjoiYXBvbXV4IiwiZW1haWwiOiJhZG1pbkBhcG9tdXguY29tIiwiaWQiOiI1ZTFjZjNlZmEzNzAwMzhmYWY3Zjg3NzciLCJvcmdhbml6YXRpb24iOiJhcG9tdXgiLCJyZWFsbSI6InZpbmNlIiwic3ViamVjdCI6ImFwb211eCJ9LCJyZXN0cmljdGlvbnMiOnsicGVybXMiOlsiQGF1dGg6cm9sZT10ZXN0Il0sIm5hbWVzcGFjZSI6NDIsIm5ldHdvcmtzIjpbIjEyNy4wLjAuMS8zMiJdfSwiZXhwIjoxNTkwMDQyOTk4LCJpYXQiOjE1ODk5NTI5OTgsImlzcyI6Imh0dHBzOi8vbG9jYWxob3N0OjQ0NDMiLCJzdWIiOiJhcG9tdXgifQ.FsYFkIzR5XXoiujjaAiYLyhIW1j0bQHuEhX8eEgIb-M`,
			},
			Restrictions{Namespace: "", Permissions: nil},
			true,
		},
		{
			"token with invalid perms type",
			args{
				`eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9.eyJyZWFsbSI6IlZpbmNlIiwiZGF0YSI6eyJhY2NvdW50IjoiYXBvbXV4IiwiZW1haWwiOiJhZG1pbkBhcG9tdXguY29tIiwiaWQiOiI1ZTFjZjNlZmEzNzAwMzhmYWY3Zjg3NzciLCJvcmdhbml6YXRpb24iOiJhcG9tdXgiLCJyZWFsbSI6InZpbmNlIiwic3ViamVjdCI6ImFwb211eCJ9LCJyZXN0cmljdGlvbnMiOnsicGVybXMiOjQyLCJuYW1lc3BhY2UiOiIvYXBvbXV4L2NoaWxkIiwibmV0d29ya3MiOlsiMTI3LjAuMC4xLzMyIl19LCJleHAiOjE1OTAwNDI5OTgsImlhdCI6MTU4OTk1Mjk5OCwiaXNzIjoiaHR0cHM6Ly9sb2NhbGhvc3Q6NDQ0MyIsInN1YiI6ImFwb211eCJ9.CXF5OH5nsutTDKceZELYxjTJi1MxRrBSatt2WdqUia4`,
			},
			Restrictions{Namespace: "", Permissions: nil},
			true,
		},
		{
			"token with invalid perms content type",
			args{
				`eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9.eyJyZWFsbSI6IlZpbmNlIiwiZGF0YSI6eyJhY2NvdW50IjoiYXBvbXV4IiwiZW1haWwiOiJhZG1pbkBhcG9tdXguY29tIiwiaWQiOiI1ZTFjZjNlZmEzNzAwMzhmYWY3Zjg3NzciLCJvcmdhbml6YXRpb24iOiJhcG9tdXgiLCJyZWFsbSI6InZpbmNlIiwic3ViamVjdCI6ImFwb211eCJ9LCJyZXN0cmljdGlvbnMiOnsicGVybXMiOls0Ml0sIm5hbWVzcGFjZSI6Ii9hcG9tdXgvY2hpbGQiLCJuZXR3b3JrcyI6WyIxMjcuMC4wLjEvMzIiXX0sImV4cCI6MTU5MDA0Mjk5OCwiaWF0IjoxNTg5OTUyOTk4LCJpc3MiOiJodHRwczovL2xvY2FsaG9zdDo0NDQzIiwic3ViIjoiYXBvbXV4In0.JIg_iFiiWnpqkvWejomrofR3R_YY5h3r3SQFmmriR7g`,
			},
			Restrictions{Namespace: "", Permissions: nil},
			true,
		},
		{
			"token with invalid networks type",
			args{
				`eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9.eyJyZWFsbSI6IlZpbmNlIiwiZGF0YSI6eyJhY2NvdW50IjoiYXBvbXV4IiwiZW1haWwiOiJhZG1pbkBhcG9tdXguY29tIiwiaWQiOiI1ZTFjZjNlZmEzNzAwMzhmYWY3Zjg3NzciLCJvcmdhbml6YXRpb24iOiJhcG9tdXgiLCJyZWFsbSI6InZpbmNlIiwic3ViamVjdCI6ImFwb211eCJ9LCJyZXN0cmljdGlvbnMiOnsicGVybXMiOlsiQGF1dGg6cm9sZT10ZXN0Il0sIm5hbWVzcGFjZSI6Ii9hcG9tdXgvY2hpbGQiLCJuZXR3b3JrcyI6NDJ9LCJleHAiOjE1OTAwNDI5OTgsImlhdCI6MTU4OTk1Mjk5OCwiaXNzIjoiaHR0cHM6Ly9sb2NhbGhvc3Q6NDQ0MyIsInN1YiI6ImFwb211eCJ9.RffIbk1AJOxOPz_Gr1SAnqdanNDnOnNtGuEAIPU5Hk4`,
			},
			Restrictions{Namespace: "", Permissions: nil},
			true,
		},
		{
			"token with invalid networks content type",
			args{
				`eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9.eyJyZWFsbSI6IlZpbmNlIiwiZGF0YSI6eyJhY2NvdW50IjoiYXBvbXV4IiwiZW1haWwiOiJhZG1pbkBhcG9tdXguY29tIiwiaWQiOiI1ZTFjZjNlZmEzNzAwMzhmYWY3Zjg3NzciLCJvcmdhbml6YXRpb24iOiJhcG9tdXgiLCJyZWFsbSI6InZpbmNlIiwic3ViamVjdCI6ImFwb211eCJ9LCJyZXN0cmljdGlvbnMiOnsicGVybXMiOlsiQGF1dGg6cm9sZT10ZXN0Il0sIm5hbWVzcGFjZSI6Ii9hcG9tdXgvY2hpbGQiLCJuZXR3b3JrcyI6WzQyXX0sImV4cCI6MTU5MDA0Mjk5OCwiaWF0IjoxNTg5OTUyOTk4LCJpc3MiOiJodHRwczovL2xvY2FsaG9zdDo0NDQzIiwic3ViIjoiYXBvbXV4In0.zJJzHJsQu6dsIDhvtp3O-zDb6W1LeLgA76_1BBX8enE`,
			},
			Restrictions{Namespace: "", Permissions: nil},
			true,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := GetRestrictions(tt.args.token)
			if (err != nil) != tt.wantErr {
				t.Errorf("GetRestrictions() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if !reflect.DeepEqual(got, tt.want) {
				t.Errorf("GetRestrictions() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestResolveRestrictions(t *testing.T) {
	type args struct {
		restrictions Restrictions
	}
	tests := []struct {
		name string
		args args
		want PermissionMap
	}{

		{
			"simple unique permission",
			args{
				Restrictions{
					Permissions: []string{
						"api1,get,post",
					},
				},
			},
			PermissionMap{
				"api1": {"get": true, "post": true},
			},
		},

		{
			"simple double permissions",
			args{
				Restrictions{
					Permissions: []string{
						"api1,get,post",
						"api2,get",
					},
				},
			},
			PermissionMap{
				"api1": {"get": true, "post": true},
				"api2": {"get": true},
			},
		},

		{
			"simple overlapping permissions",
			args{
				Restrictions{
					Permissions: []string{
						"api1,get,post",
						"api1,delete",
					},
				},
			},
			PermissionMap{
				"api1": {"get": true, "post": true, "delete": true},
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := ResolveRestrictions(tt.args.restrictions); !reflect.DeepEqual(got, tt.want) {
				t.Errorf("ResolveRestrictions() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestRestrictions_ComputeNamespaceRestriction(t *testing.T) {
	type fields struct {
		Namespace   string
		Permissions []string
		Networks    []string
	}
	type args struct {
		requested string
	}
	tests := []struct {
		name    string
		fields  fields
		args    args
		want    string
		wantErr bool
	}{
		{
			"no original, no requested",
			fields{
				"",
				nil,
				nil,
			},
			args{
				"",
			},
			"",
			false,
		},
		{
			"original, no requested",
			fields{
				"/ns",
				nil,
				nil,
			},
			args{
				"",
			},
			"/ns",
			false,
		},
		{
			"original, identical requested",
			fields{
				"/ns",
				nil,
				nil,
			},
			args{
				"/ns",
			},
			"/ns",
			false,
		},
		{
			"original, child requested",
			fields{
				"/ns",
				nil,
				nil,
			},
			args{
				"/ns/child",
			},
			"/ns/child",
			false,
		},
		{
			"original, root requested",
			fields{
				"/ns",
				nil,
				nil,
			},
			args{
				"/",
			},
			"",
			true,
		},
		{
			"original, / requested",
			fields{
				"/parent/ns",
				nil,
				nil,
			},
			args{
				"/parent",
			},
			"",
			true,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			r := Restrictions{
				Namespace:   tt.fields.Namespace,
				Permissions: tt.fields.Permissions,
				Networks:    tt.fields.Networks,
			}
			got, err := r.ComputeNamespaceRestriction(tt.args.requested)
			if (err != nil) != tt.wantErr {
				t.Errorf("Restrictions.ComputeNamespaceRestriction() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if got != tt.want {
				t.Errorf("Restrictions.ComputeNamespaceRestriction() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestRestrictions_ComputeNetworkRestrictions(t *testing.T) {
	type fields struct {
		Namespace   string
		Permissions []string
		Networks    []string
	}
	type args struct {
		requested []string
	}
	tests := []struct {
		name    string
		fields  fields
		args    args
		want    []string
		wantErr bool
	}{
		{
			"no original, no requested",
			fields{
				"",
				nil,
				nil,
			},
			args{
				nil,
			},
			nil,
			false,
		},
		{
			"no original, requested",
			fields{
				"",
				nil,
				nil,
			},
			args{
				[]string{"1.0.0.0/8"},
			},
			[]string{"1.0.0.0/8"},
			false,
		},

		{
			"single original, single valid requested",
			fields{
				"",
				nil,
				[]string{"1.0.0.0/8"},
			},
			args{
				[]string{"1.1.0.0/16"},
			},
			[]string{"1.1.0.0/16"},
			false,
		},
		{
			"single original, dual valid requested",
			fields{
				"",
				nil,
				[]string{"1.0.0.0/8"},
			},
			args{
				[]string{"1.1.0.0/16", "1.2.0.0/16"},
			},
			[]string{"1.1.0.0/16", "1.2.0.0/16"},
			false,
		},
		{
			"single original, dual invalid requested",
			fields{
				"",
				nil,
				[]string{"1.0.0.0/8"},
			},
			args{
				[]string{"10.1.0.0/16", "10.2.0.0/16"},
			},
			nil,
			true,
		},
		{
			"single original, one valid and one invalid requested",
			fields{
				"",
				nil,
				[]string{"1.0.0.0/8"},
			},
			args{
				[]string{"1.1.0.0/16", "10.2.0.0/16"},
			},
			nil,
			true,
		},
		{
			"single original, identical requested",
			fields{
				"",
				nil,
				[]string{"1.0.0.0/8"},
			},
			args{
				[]string{"1.0.0.0/8"},
			},
			[]string{"1.0.0.0/8"},
			false,
		},

		{
			"dual original, single valid requested",
			fields{
				"",
				nil,
				[]string{"1.0.0.0/8", "2.0.0.0/8"},
			},
			args{
				[]string{"1.1.0.0/16"},
			},
			[]string{"1.1.0.0/16"},
			false,
		},
		{
			"dual original, single invalid requested",
			fields{
				"",
				nil,
				[]string{"1.0.0.0/8", "2.0.0.0/8"},
			},
			args{
				[]string{"3.1.0.0/16"},
			},
			nil,
			true,
		},
		{
			"dual original, dual valid requested",
			fields{
				"",
				nil,
				[]string{"1.0.0.0/8", "2.0.0.0/8"},
			},
			args{
				[]string{"1.1.0.0/16", "2.1.0.0/16"},
			},
			[]string{"1.1.0.0/16", "2.1.0.0/16"},
			false,
		},
		{
			"dual original, dual one valid and on invalid requested",
			fields{
				"",
				nil,
				[]string{"1.0.0.0/8", "2.0.0.0/8"},
			},
			args{
				[]string{"1.1.0.0/16", "10.1.0.0/16"},
			},
			nil,
			true,
		},
		{
			"dual original, identical dual requested",
			fields{
				"",
				nil,
				[]string{"1.0.0.0/8", "2.0.0.0/8"},
			},
			args{
				[]string{"1.0.0.0/8", "2.0.0.0/8"},
			},
			[]string{"1.0.0.0/8", "2.0.0.0/8"},
			false,
		},

		{
			"invalid original",
			fields{
				"",
				nil,
				[]string{"chien"},
			},
			args{
				[]string{"1.1.0.0/16", "10.1.0.0/16"},
			},
			nil,
			true,
		},
		{
			"invalid requested",
			fields{
				"",
				nil,
				[]string{"1.1.0.0/16", "10.1.0.0/16"},
			},
			args{
				[]string{"chien"},
			},
			nil,
			true,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			r := Restrictions{
				Namespace:   tt.fields.Namespace,
				Permissions: tt.fields.Permissions,
				Networks:    tt.fields.Networks,
			}
			got, err := r.ComputeNetworkRestrictions(tt.args.requested)
			if (err != nil) != tt.wantErr {
				t.Errorf("Restrictions.ComputeNetworkRestrictions() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if !reflect.DeepEqual(got, tt.want) {
				t.Errorf("Restrictions.ComputeNetworkRestrictions() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestRestrictions_ComputePermissionsRestrictions(t *testing.T) {
	type fields struct {
		Namespace   string
		Permissions []string
		Networks    []string
	}
	type args struct {
		requested []string
	}
	tests := []struct {
		name    string
		fields  fields
		args    args
		want    []string
		wantErr bool
	}{
		{
			"no original, no requested",
			fields{
				"",
				nil,
				nil,
			},
			args{
				nil,
			},
			nil,
			false,
		},
		{
			"original, no requested",
			fields{
				"",
				[]string{"r,get"},
				nil,
			},
			args{
				nil,
			},
			[]string{"r,get"},
			false,
		},
		{
			"no original, requested",
			fields{
				"",
				nil,
				nil,
			},
			args{
				[]string{"r,get"},
			},
			[]string{"r,get"},
			false,
		},

		{
			"single original, single valid requested",
			fields{
				"",
				[]string{"r,get,post"},
				nil,
			},
			args{
				[]string{"r,get"},
			},
			[]string{"r,get"},
			false,
		},
		{
			"single original, single invalid requested",
			fields{
				"",
				[]string{"r,get"},
				nil,
			},
			args{
				[]string{"r,post"},
			},
			nil,
			true,
		},
		{
			"single original, identical requested",
			fields{
				"",
				[]string{"r,get"},
				nil,
			},
			args{
				[]string{"r,get"},
			},
			[]string{"r,get"},
			false,
		},

		{
			"single original, dual valid requested",
			fields{
				"",
				[]string{"r,get,post"},
				nil,
			},
			args{
				[]string{"r,get", "r,post"},
			},
			[]string{"r,get", "r,post"},
			false,
		},
		{
			"single original, dual invalid requested",
			fields{
				"",
				[]string{"r,get"},
				nil,
			},
			args{
				[]string{"r,post", "r,put"},
			},
			nil,
			true,
		},
		{
			"single original, one valid and one invalid requested",
			fields{
				"",
				[]string{"r,get"},
				nil,
			},
			args{
				[]string{"r,get", "r,delete"},
			},
			nil,
			true,
		},

		{
			"dual original, dual valid requested",
			fields{
				"",
				[]string{"r1,get,post", "r2,get,post"},
				nil,
			},
			args{
				[]string{"r1,get", "r2,post"},
			},
			[]string{"r1,get", "r2,post"},
			false,
		},
		{
			"dual original, dual invalid requested",
			fields{
				"",
				[]string{"r1,get", "r2,get,post"},
				nil,
			},
			args{
				[]string{"r1,delete", "r2,delete"},
			},
			nil,
			true,
		},
		{
			"dual original, one valid and one invalid requested",
			fields{
				"",
				[]string{"r1,get,post", "r2,get,post"},
				nil,
			},
			args{
				[]string{"r1,get", "r2,delete"},
			},
			nil,
			true,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			r := Restrictions{
				Namespace:   tt.fields.Namespace,
				Permissions: tt.fields.Permissions,
				Networks:    tt.fields.Networks,
			}
			got, err := r.ComputePermissionsRestrictions(tt.args.requested)
			if (err != nil) != tt.wantErr {
				t.Errorf("Restrictions.ComputePermissionsRestrictions() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if !reflect.DeepEqual(got, tt.want) {
				t.Errorf("Restrictions.ComputePermissionsRestrictions() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestTokenUtils_UnsecureClaimsMap(t *testing.T) {

	Convey("Given I have a valid token", t, func() {

		token := "eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9.eyJyZWFsbSI6IlZpbmNlIiwiZGF0YSI6eyJhY2NvdW50IjoiYXBvbXV4IiwiZW1haWwiOiJhZG1pbkBhcG9tdXguY29tIiwiaWQiOiI1YTZhNTUxMTdkZGYxZjIxMmY4ZWIwY2UiLCJvcmdhbml6YXRpb24iOiJhcG9tdXgiLCJyZWFsbSI6InZpbmNlIn0sImF1ZCI6ImFwb3JldG8uY29tIiwiZXhwIjoxNTIwNjQ5MTAyLCJpYXQiOjE1MTgwNTcxMDIsImlzcyI6Im1pZGdhcmQuYXBvbXV4LmNvbSIsInN1YiI6ImFwb211eCJ9.jvh034mNSV-Fy--GIGnnYeWouluV6CexC9_8IHJ-IR4"

		Convey("When I UnsecureClaimsMap", func() {

			claims, err := UnsecureClaimsMap(token)

			Convey("Then err should be nil", func() {
				So(err, ShouldBeNil)
			})

			Convey("Then claims should be correct", func() {
				So(claims, ShouldNotBeNil)
				So(claims["data"].(map[string]interface{})["realm"].(string), ShouldEqual, "vince")
				So(claims["sub"].(string), ShouldEqual, "apomux")
			})
		})
	})

	Convey("Given I have a token an invalid token", t, func() {

		token := "not good"

		Convey("When I UnsecureClaimsMap", func() {

			claims, err := UnsecureClaimsMap(token)

			Convey("Then err should be nil", func() {
				So(err, ShouldNotBeNil)
				So(err.Error(), ShouldEqual, "invalid jwt: not enough segments")
			})

			Convey("Then claims should be nil", func() {
				So(claims, ShouldBeNil)
			})
		})
	})

	Convey("Given I have a token an empty token", t, func() {

		token := ""

		Convey("When I UnsecureClaimsMap", func() {

			claims, err := UnsecureClaimsMap(token)

			Convey("Then err should be nil", func() {
				So(err, ShouldNotBeNil)
				So(err.Error(), ShouldEqual, "invalid jwt: empty")
			})

			Convey("Then claims should be nil", func() {
				So(claims, ShouldBeNil)
			})
		})
	})

	Convey("Given I have a token a token with invalid base64", t, func() {

		token := "eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9.XXX.jvh034mNSV-Fy--GIGnnYeWouluV6CexC9_8IHJ-IR4"

		Convey("When I UnsecureClaimsMap", func() {

			claims, err := UnsecureClaimsMap(token)

			Convey("Then err should not be nil", func() {
				So(err, ShouldNotBeNil)
				So(err.Error(), ShouldEqual, "invalid jwt: invalid character ']' looking for beginning of value")
			})

			Convey("Then claims should be nil", func() {
				So(claims, ShouldBeNil)
			})
		})
	})

	Convey("Given I have a token a token with valid base64 URL", t, func() {

		token := "eyJhbGciOiJFUzI1NiIsInR5cCI6IkpXVCJ9.ew0KICAicmVhbG0iOiAiPz8_Pz8_Ig0KfQ0KDQo.hJwIQMyvQ1-VqfImMnaNYDPq6uxy8fscLLVji9loOW9OrPLZfqqfMxrg0tko9CPLxoB4wG3_KWshPPeyUBUspQ"

		Convey("When I UnsecureClaimsMap", func() {

			claims, err := UnsecureClaimsMap(token)

			Convey("Then err should be nil", func() {
				So(err, ShouldBeNil)
			})

			Convey("Then claims should not be nil", func() {
				So(claims, ShouldNotBeNil)
			})
		})
	})

	Convey("Given I have a token a token with invalid json data", t, func() {

		token := "eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9.eyJicm9rZW46ICJqc29u.jvh034mNSV-Fy--GIGnnYeWouluV6CexC9_8IHJ-IR4"

		Convey("When I UnsecureClaimsMap", func() {

			alg, err := UnsecureClaimsMap(token)

			Convey("Then err should be nil", func() {
				So(err, ShouldNotBeNil)
				So(err.Error(), ShouldEqual, "invalid jwt: invalid character 'j' after object key")
			})

			Convey("Then alg should be empty", func() {
				So(alg, ShouldBeEmpty)
			})
		})
	})
}
