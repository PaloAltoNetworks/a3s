package conf

import (
	"crypto/tls"
	"crypto/x509"
	"os"
	"reflect"
	"testing"

	. "github.com/smartystreets/goconvey/convey"
	"go.aporeto.io/a3s/pkgs/api"
	"go.aporeto.io/elemental"
	"go.aporeto.io/tg/tglib"
)

func TestAPIServerConf(t *testing.T) {

	_, clientpool, tlscert := makeFixtures()

	Convey("disabled tls config", t, func() {
		cfg := APIServerConf{
			TLSDisable: true,
		}
		tlscfg, err := cfg.TLSConfig()
		So(err, ShouldBeNil)
		So(tlscfg, ShouldBeNil)
	})

	Convey("tls config with client CA", t, func() {
		cfg := APIServerConf{
			TLSClientCA: "fixtures/ca-cert.pem",
		}
		tlscfg, err := cfg.TLSConfig()
		So(err, ShouldBeNil)
		So(tlscfg, ShouldNotBeNil)
		So(tlscfg.ClientCAs.Subjects(), ShouldResemble, clientpool.Subjects()) // nolint: staticcheck
	})

	Convey("tls config with certificate", t, func() {
		cfg := APIServerConf{
			TLSCertificate: "fixtures/cert-cert.pem",
			TLSKey:         "fixtures/cert-key.pem",
		}
		tlscfg, err := cfg.TLSConfig()
		So(err, ShouldBeNil)
		So(tlscfg, ShouldNotBeNil)
		So(tlscfg.Certificates[0].Certificate, ShouldResemble, tlscert.Certificate)
	})

	Convey("tls config with bad certificate", t, func() {
		cfg := APIServerConf{
			TLSCertificate: "fixtures/cert-key.pem",
			TLSKey:         "fixtures/cert-key.pem",
		}
		tlscfg, err := cfg.TLSConfig()
		So(tlscfg, ShouldBeNil)
		So(err, ShouldNotBeNil)
		So(err.Error(), ShouldStartWith, "unable to load client certificate:")
	})

	Convey("tls config with bad ca", t, func() {
		cfg := APIServerConf{
			TLSClientCA: "nope",
		}
		tlscfg, err := cfg.TLSConfig()
		So(tlscfg, ShouldBeNil)
		So(err, ShouldNotBeNil)
		So(err.Error(), ShouldStartWith, "unable to load ca file:")
	})
}
func makeFixtures() (syspool *x509.CertPool, custompool *x509.CertPool, cert tls.Certificate) {

	var err error

	ccad, err := os.ReadFile("fixtures/ca-cert.pem")
	if err != nil {
		panic(err)
	}
	public, private, err := tglib.ReadCertificatePEM("fixtures/cert-cert.pem", "fixtures/cert-key.pem", "")
	if err != nil {
		panic(err)
	}
	cert, err = tglib.ToTLSCertificate(public, private)
	if err != nil {
		panic(err)
	}

	cuspool := x509.NewCertPool()
	cuspool.AppendCertsFromPEM(ccad)

	syspool, err = x509.SystemCertPool()
	if err != nil {
		panic(err)
	}

	return syspool, cuspool, cert
}

func TestMongoConf(t *testing.T) {

	syspool, cuspool, tlscert := makeFixtures()

	Convey("disabled tls config", t, func() {
		cfg := MongoConf{
			MongoTLSDisable: true,
		}
		tlscfg, err := cfg.TLSConfig()
		So(err, ShouldBeNil)
		So(tlscfg, ShouldBeNil)
	})

	Convey("skipped tls config", t, func() {
		cfg := MongoConf{
			MongoTLSSkip: true,
		}
		tlscfg, err := cfg.TLSConfig()
		So(err, ShouldBeNil)
		So(tlscfg, ShouldNotBeNil)
		So(tlscfg.InsecureSkipVerify, ShouldBeTrue)
	})

	Convey("tls config with system CA", t, func() {
		cfg := MongoConf{}
		tlscfg, err := cfg.TLSConfig()
		So(err, ShouldBeNil)
		So(tlscfg, ShouldNotBeNil)
		So(tlscfg.RootCAs.Subjects(), ShouldResemble, syspool.Subjects()) // nolint: staticcheck
	})

	Convey("tls config with custom CA", t, func() {
		cfg := MongoConf{
			MongoTLSCA: "fixtures/ca-cert.pem",
		}
		tlscfg, err := cfg.TLSConfig()
		So(err, ShouldBeNil)
		So(tlscfg, ShouldNotBeNil)
		So(tlscfg.RootCAs.Subjects(), ShouldResemble, cuspool.Subjects()) // nolint: staticcheck
	})

	Convey("tls config with certificate", t, func() {
		cfg := MongoConf{
			MongoTLSCertificate: "fixtures/cert-cert.pem",
			MongoTLSKey:         "fixtures/cert-key.pem",
		}
		tlscfg, err := cfg.TLSConfig()
		So(err, ShouldBeNil)
		So(tlscfg, ShouldNotBeNil)
		So(tlscfg.Certificates[0].Certificate, ShouldResemble, tlscert.Certificate)
	})

	Convey("tls config with bad certificate", t, func() {
		cfg := MongoConf{
			MongoTLSCertificate: "fixtures/cert-key.pem",
			MongoTLSKey:         "fixtures/cert-key.pem",
		}
		tlscfg, err := cfg.TLSConfig()
		So(tlscfg, ShouldBeNil)
		So(err, ShouldNotBeNil)
		So(err.Error(), ShouldStartWith, "unable to load client certificate:")
	})

	Convey("tls config with bad ca", t, func() {
		cfg := MongoConf{
			MongoTLSCA: "nope",
		}
		tlscfg, err := cfg.TLSConfig()
		So(tlscfg, ShouldBeNil)
		So(err, ShouldNotBeNil)
		So(err.Error(), ShouldStartWith, "unable to load ca file:")
	})
}

func TestNATSConf(t *testing.T) {

	syspool, cuspool, tlscert := makeFixtures()

	Convey("disabled tls config", t, func() {
		cfg := NATSConf{
			NATSTLSDisable: true,
		}
		tlscfg, err := cfg.TLSConfig()
		So(err, ShouldBeNil)
		So(tlscfg, ShouldBeNil)
	})

	Convey("skipped tls config", t, func() {
		cfg := NATSConf{
			NATSTLSSkip: true,
		}
		tlscfg, err := cfg.TLSConfig()
		So(err, ShouldBeNil)
		So(tlscfg, ShouldNotBeNil)
		So(tlscfg.InsecureSkipVerify, ShouldBeTrue)
	})

	Convey("tls config with system CA", t, func() {
		cfg := NATSConf{}
		tlscfg, err := cfg.TLSConfig()
		So(err, ShouldBeNil)
		So(tlscfg, ShouldNotBeNil)
		So(tlscfg.RootCAs.Subjects(), ShouldResemble, syspool.Subjects()) // nolint: staticcheck
	})

	Convey("tls config with custom CA", t, func() {
		cfg := NATSConf{
			NATSTLSCA: "fixtures/ca-cert.pem",
		}
		tlscfg, err := cfg.TLSConfig()
		So(err, ShouldBeNil)
		So(tlscfg, ShouldNotBeNil)
		So(tlscfg.RootCAs.Subjects(), ShouldResemble, cuspool.Subjects()) // nolint: staticcheck
	})

	Convey("tls config with certificate", t, func() {
		cfg := NATSConf{
			NATSTLSCertificate: "fixtures/cert-cert.pem",
			NATSTLSKey:         "fixtures/cert-key.pem",
		}
		tlscfg, err := cfg.TLSConfig()
		So(err, ShouldBeNil)
		So(tlscfg, ShouldNotBeNil)
		So(tlscfg.Certificates[0].Certificate, ShouldResemble, tlscert.Certificate)
	})

	Convey("tls config with bad certificate", t, func() {
		cfg := NATSConf{
			NATSTLSCertificate: "fixtures/cert-key.pem",
			NATSTLSKey:         "fixtures/cert-key.pem",
		}
		tlscfg, err := cfg.TLSConfig()
		So(tlscfg, ShouldBeNil)
		So(err, ShouldNotBeNil)
		So(err.Error(), ShouldStartWith, "unable to load client certificate:")
	})

	Convey("tls config with bad ca", t, func() {
		cfg := NATSConf{
			NATSTLSCA: "nope",
		}
		tlscfg, err := cfg.TLSConfig()
		So(tlscfg, ShouldBeNil)
		So(err, ShouldNotBeNil)
		So(err.Error(), ShouldStartWith, "unable to load ca file:")
	})
}

func TestGatewayConf_GWPrivateOverrides(t *testing.T) {
	tests := []struct {
		name    string
		init    func(t *testing.T) *GatewayConf
		inspect func(r *GatewayConf, t *testing.T) //inspects receiver after test run

		want1 map[elemental.Identity]bool
	}{
		{
			"simple override",
			func(t *testing.T) *GatewayConf {
				return &GatewayConf{
					GWOverridePrivate: []string{"namespace:public", "authorization:private"},
				}
			},
			nil,
			map[elemental.Identity]bool{
				api.NamespaceIdentity:     true,
				api.AuthorizationIdentity: false,
			},
		},
		{
			"* public override",
			func(t *testing.T) *GatewayConf {
				return &GatewayConf{
					GWOverridePrivate: []string{"*:public"},
				}
			},
			nil,
			func() map[elemental.Identity]bool {
				m := map[elemental.Identity]bool{}
				for _, ident := range api.AllIdentities() {
					m[ident] = true
				}
				return m
			}(),
		},
		{
			"* private override",
			func(t *testing.T) *GatewayConf {
				return &GatewayConf{
					GWOverridePrivate: []string{"*:private"},
				}
			},
			nil,
			func() map[elemental.Identity]bool {
				m := map[elemental.Identity]bool{}
				for _, ident := range api.AllIdentities() {
					m[ident] = false
				}
				return m
			}(),
		},
		{
			"mixed private override",
			func(t *testing.T) *GatewayConf {
				return &GatewayConf{
					GWOverridePrivate: []string{"*:private", "namespace:public"},
				}
			},
			nil,
			func() map[elemental.Identity]bool {
				m := map[elemental.Identity]bool{}
				for _, ident := range api.AllIdentities() {
					m[ident] = false
				}
				m[api.NamespaceIdentity] = true
				return m
			}(),
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			receiver := tt.init(t)
			got1 := receiver.GWPrivateOverrides()

			if tt.inspect != nil {
				tt.inspect(receiver, t)
			}

			if !reflect.DeepEqual(got1, tt.want1) {
				t.Errorf("GatewayConf.GWPrivateOverrides got1 = %v, want1: %v", got1, tt.want1)
			}
		})
	}
}
