package oidcissuer

import (
	"context"
	"crypto"
	"crypto/x509"
	"crypto/x509/pkix"
	"encoding/json"
	"encoding/pem"
	"fmt"
	"net/http"
	"net/http/httptest"
	"reflect"
	"testing"

	. "github.com/smartystreets/goconvey/convey"
	"go.aporeto.io/a3s/pkgs/api"
	"go.aporeto.io/tg/tglib"
)

func getECCert(subject pkix.Name, opts ...tglib.IssueOption) (*x509.Certificate, crypto.PrivateKey) {

	certBlock, keyBlock, err := tglib.Issue(subject, opts...)
	if err != nil {
		panic(err)
	}

	cert, err := tglib.ParseCertificate(pem.EncodeToMemory(certBlock))
	if err != nil {
		panic(err)
	}

	key, err := tglib.PEMToKey(keyBlock)
	if err != nil {
		panic(err)
	}

	return cert, key
}

func TestNew(t *testing.T) {

	Convey("Calling New should work ", t, func() {
		src := api.NewOIDCSource()
		src.Name = "name"
		src.Namespace = "/ns"
		iss, _ := New(context.Background(), src, map[string]interface{}{"hello": "world"})
		So(iss.(*oidcIssuer).source, ShouldEqual, src)
		So(iss.Issue().Source.Type, ShouldEqual, "oidc")
		So(iss.Issue().Source.Name, ShouldEqual, "name")
		So(iss.Issue().Source.Namespace, ShouldEqual, "/ns")
	})

	Convey("Calling New with a source and a modifier should work", t, func() {

		src := api.NewOIDCSource()
		ts := httptest.NewTLSServer(http.HandlerFunc(func(w http.ResponseWriter, req *http.Request) {
			d, _ := json.Marshal([]string{"aa=aa", "bb=bb"})
			w.WriteHeader(http.StatusOK)
			w.Write(d) // nolint
		}))
		defer ts.Close()

		usercert1, userkey1 := getECCert(pkix.Name{})
		cab, _ := tglib.CertToPEM(ts.Certificate())
		certb, _ := tglib.CertToPEM(usercert1)
		keyb, _ := tglib.KeyToPEM(userkey1)
		src.Modifier = api.NewIdentityModifier()
		src.Modifier.CA = string(pem.EncodeToMemory(cab))
		src.Modifier.URL = ts.URL
		src.Modifier.Certificate = string(pem.EncodeToMemory(certb))
		src.Modifier.Key = string(pem.EncodeToMemory(keyb))

		iss, _ := New(context.Background(), src, map[string]interface{}{"hello": "world"})
		So(iss.(*oidcIssuer).source, ShouldEqual, src)
		So(iss.Issue().Identity, ShouldResemble, []string{"aa=aa", "bb=bb"})
	})

	Convey("Calling New with a source and a modifier with mussing tls info", t, func() {

		src := api.NewOIDCSource()
		ts := httptest.NewTLSServer(http.HandlerFunc(func(w http.ResponseWriter, req *http.Request) {
			d, _ := json.Marshal([]string{"aa=aa", "bb=bb"})
			w.WriteHeader(http.StatusOK)
			w.Write(d) // nolint
		}))
		defer ts.Close()

		cab, _ := tglib.CertToPEM(ts.Certificate())
		src.Modifier = api.NewIdentityModifier()
		src.Modifier.CA = string(pem.EncodeToMemory(cab))
		src.Modifier.URL = ts.URL

		_, err := New(context.Background(), src, map[string]interface{}{"hello": "world"})
		So(err, ShouldNotBeNil)
		So(err.Error(), ShouldEqual, `unable to prepare source modifier: unable to create certificate: could not read key data from bytes: ''`)
	})

	Convey("Calling New with a source and a modifier that returns an error", t, func() {

		src := api.NewOIDCSource()
		ts := httptest.NewTLSServer(http.HandlerFunc(func(w http.ResponseWriter, req *http.Request) {
			w.WriteHeader(http.StatusForbidden)
		}))
		defer ts.Close()

		usercert1, userkey1 := getECCert(pkix.Name{})
		cab, _ := tglib.CertToPEM(ts.Certificate())
		certb, _ := tglib.CertToPEM(usercert1)
		keyb, _ := tglib.KeyToPEM(userkey1)
		src.Modifier = api.NewIdentityModifier()
		src.Modifier.CA = string(pem.EncodeToMemory(cab))
		src.Modifier.URL = ts.URL
		src.Modifier.Certificate = string(pem.EncodeToMemory(certb))
		src.Modifier.Key = string(pem.EncodeToMemory(keyb))

		_, err := New(context.Background(), src, map[string]interface{}{"hello": "world"})
		So(err, ShouldNotBeNil)
		So(err.Error(), ShouldEqual, `unable to call modifier: service returned an error: 403 Forbidden`)
	})
}

func Test_computeOIDClaims(t *testing.T) {
	type args struct {
		claims map[string]interface{}
	}
	tests := []struct {
		name string
		args func(t *testing.T) args

		want1 []string
	}{
		{
			"standard",
			func(*testing.T) args {
				return args{
					map[string]interface{}{
						"@@string": "value",
						"strings":  []string{"v1", "v2"},
						"int":      42,
						"ints":     []int{1, 2},
						"bool":     true,
						"ifaces":   []interface{}{"a", "b"},
						"map":      map[string]interface{}{},
						"float":    42.42,
						"floats":   []float64{1.2, 3.4},
						"error":    fmt.Errorf("yo"),
					},
				}
			},
			[]string{
				"bool=true",
				"error=yo",
				"float=42.420000",
				"floats=1.200000",
				"floats=3.400000",
				"ifaces=a",
				"ifaces=b",
				"int=42",
				"ints=1",
				"ints=2",
				"map=map[]",
				"string=value",
				"strings=v1",
				"strings=v2",
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			tArgs := tt.args(t)

			got1 := computeOIDClaims(tArgs.claims)

			if !reflect.DeepEqual(got1, tt.want1) {
				t.Errorf("computeOIDClaims got1 = %v, want1: %v", got1, tt.want1)
			}
		})
	}
}
