package ldapissuer

import (
	"fmt"
	"reflect"
	"testing"

	"github.com/go-ldap/ldap/v3"
	. "github.com/smartystreets/goconvey/convey"
	"go.aporeto.io/a3s/pkgs/api"
)

func TestErrLDAP(t *testing.T) {
	Convey("ErrLDAP should work", t, func() {
		e := fmt.Errorf("boom")
		err := ErrLDAP{Err: e}
		So(err.Error(), ShouldEqual, "ldap error: boom")
		So(err.Unwrap(), ShouldEqual, e)
	})
}

func TestNewLDAPIssuer(t *testing.T) {
	Convey("Calling NewLDAPIssuer should work", t, func() {
		src := api.NewLDAPSource()
		src.Namespace = "/my/ns"
		src.Name = "my-src"
		iss := newLDAPIssuer(src)
		So(iss.source, ShouldEqual, src)
		So(iss.token.Source.Type, ShouldEqual, "ldap")
		So(iss.token.Source.Namespace, ShouldEqual, "/my/ns")
		So(iss.token.Source.Name, ShouldEqual, "my-src")
		So(iss.Issue(), ShouldEqual, iss.token)
	})
}

func TestFromCredential(t *testing.T) {

	Convey("Given a LDAP Issuer and a source with no address", t, func() {
		src := api.NewLDAPSource()
		src.Namespace = "/my/ns"
		src.Name = "my-src"
		iss := newLDAPIssuer(src)
		err := iss.fromCredentials("", "")
		So(err, ShouldNotBeNil)
		So(err.Error(), ShouldEqual, `ldap error: cannot dial: LDAP Result Code 200 "Network Error": dial tcp: missing address`)
	})

	Convey("Given a LDAP Issuer and a source with a ca", t, func() {
		src := api.NewLDAPSource()
		src.Namespace = "/my/ns"
		src.Name = "my-src"
		src.CertificateAuthority = "a-ca"
		iss := newLDAPIssuer(src)
		err := iss.fromCredentials("", "")
		So(err, ShouldNotBeNil)
		So(err.Error(), ShouldEqual, `ldap error: cannot dial: LDAP Result Code 200 "Network Error": dial tcp: missing address`)
	})

	Convey("Given a LDAP Issuer and a source using TLS", t, func() {
		src := api.NewLDAPSource()
		src.Namespace = "/my/ns"
		src.Name = "my-src"
		src.SecurityProtocol = api.LDAPSourceSecurityProtocolTLS
		iss := newLDAPIssuer(src)
		err := iss.fromCredentials("", "")
		So(err, ShouldNotBeNil)
		So(err.Error(), ShouldEqual, `ldap error: cannot dial tls: LDAP Result Code 200 "Network Error": dial tcp: missing address`)
	})
}

func Test_computeLDAPClaims(t *testing.T) {
	type args struct {
		entry *ldap.Entry
		dn    *ldap.DN
		inc   map[string]struct{}
		exc   map[string]struct{}
	}
	tests := []struct {
		name string
		args func(t *testing.T) args

		want1 []string
	}{
		{
			"standard test",
			func(*testing.T) args {
				return args{
					&ldap.Entry{
						DN: "hello",
						Attributes: []*ldap.EntryAttribute{
							{Name: "userPassword", Values: []string{"skipped"}},
							{Name: "objectClass", Values: []string{"skipped"}},
							{Name: "comment", Values: []string{"skipped"}},
							{Name: "key1", Values: []string{"value1-1", "value1-2", ""}},
							{Name: "@@key2", Values: []string{"value1-1", ""}},
							{Name: "novalues", Values: nil},
						},
					},
					&ldap.DN{
						RDNs: []*ldap.RelativeDN{
							{
								Attributes: []*ldap.AttributeTypeAndValue{
									{Type: "ou", Value: "the-ou"},
								},
							},
							{
								Attributes: []*ldap.AttributeTypeAndValue{
									{Type: "dc", Value: "the-dc"},
								},
							},
						},
					},
					nil,
					nil,
				}
			},
			[]string{"dn=hello", "ou=the-ou", "dc=the-dc", "key1=value1-1", "key1=value1-2", "key2=value1-1"},
		},
		{
			"with exclude keys",
			func(*testing.T) args {
				return args{
					&ldap.Entry{
						DN: "hello",
						Attributes: []*ldap.EntryAttribute{
							{Name: "userPassword", Values: []string{"skipped"}},
							{Name: "objectClass", Values: []string{"skipped"}},
							{Name: "comment", Values: []string{"skipped"}},
							{Name: "key1", Values: []string{"value1-1", "value1-2", ""}},
							{Name: "key2", Values: []string{"value2-1", "value2-2", ""}},
							{Name: "novalues", Values: nil},
						},
					},
					&ldap.DN{
						RDNs: []*ldap.RelativeDN{
							{
								Attributes: []*ldap.AttributeTypeAndValue{
									{Type: "ou", Value: "the-ou"},
								},
							},
							{
								Attributes: []*ldap.AttributeTypeAndValue{
									{Type: "dc", Value: "the-dc"},
								},
							},
						},
					},
					nil,
					map[string]struct{}{"key2": {}},
				}
			},
			[]string{"dn=hello", "ou=the-ou", "dc=the-dc", "key1=value1-1", "key1=value1-2"},
		},
		{
			"with include keys",
			func(*testing.T) args {
				return args{
					&ldap.Entry{
						DN: "hello",
						Attributes: []*ldap.EntryAttribute{
							{Name: "key1", Values: []string{"value1-1", "value1-2", ""}},
							{Name: "key2", Values: []string{"value2-1", "value2-2", ""}},
						},
					},
					&ldap.DN{
						RDNs: []*ldap.RelativeDN{
							{
								Attributes: []*ldap.AttributeTypeAndValue{
									{Type: "ou", Value: "the-ou"},
								},
							},
							{
								Attributes: []*ldap.AttributeTypeAndValue{
									{Type: "dc", Value: "the-dc"},
								},
							},
						},
					},
					map[string]struct{}{"key2": {}},
					nil,
				}
			},
			[]string{"dn=hello", "ou=the-ou", "dc=the-dc", "key2=value2-1", "key2=value2-2"},
		},
		{
			"with both included and excludec keys",
			func(*testing.T) args {
				return args{
					&ldap.Entry{
						DN: "hello",
						Attributes: []*ldap.EntryAttribute{
							{Name: "key1", Values: []string{"value1-1", "value1-2", ""}},
							{Name: "key2", Values: []string{"value2-1", "value2-2", ""}},
						},
					},
					&ldap.DN{
						RDNs: []*ldap.RelativeDN{
							{
								Attributes: []*ldap.AttributeTypeAndValue{
									{Type: "ou", Value: "the-ou"},
								},
							},
							{
								Attributes: []*ldap.AttributeTypeAndValue{
									{Type: "dc", Value: "the-dc"},
								},
							},
						},
					},
					map[string]struct{}{"key2": {}},
					map[string]struct{}{"key2": {}},
				}
			},
			[]string{"dn=hello", "ou=the-ou", "dc=the-dc"},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			tArgs := tt.args(t)

			got1 := computeLDAPClaims(tArgs.entry, tArgs.dn, tArgs.inc, tArgs.exc)

			if !reflect.DeepEqual(got1, tt.want1) {
				t.Errorf("computeClaims got1 = %v, want1: %v", got1, tt.want1)
			}
		})
	}
}

func Test_computeLDAPInclusion(t *testing.T) {
	type args struct {
		src *api.LDAPSource
	}
	tests := []struct {
		name string
		args func(t *testing.T) args

		want1 map[string]struct{}
		want2 map[string]struct{}
	}{
		{
			"simple test",
			func(*testing.T) args {
				return args{
					&api.LDAPSource{
						IncludedKeys: []string{"a", "b"},
						IgnoredKeys:  []string{"b", "c"},
					},
				}
			},
			map[string]struct{}{"a": {}, "b": {}},
			map[string]struct{}{"b": {}, "c": {}},
		},
		{
			"nil keys",
			func(*testing.T) args {
				return args{
					&api.LDAPSource{},
				}
			},
			map[string]struct{}{},
			map[string]struct{}{},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			tArgs := tt.args(t)

			got1, got2 := computeLDPInclusion(tArgs.src)

			if !reflect.DeepEqual(got1, tt.want1) {
				t.Errorf("computeInclusion got1 = %v, want1: %v", got1, tt.want1)
			}

			if !reflect.DeepEqual(got2, tt.want2) {
				t.Errorf("computeInclusion got2 = %v, want2: %v", got2, tt.want2)
			}
		})
	}
}
