package conf

import (
	"crypto/tls"
	"crypto/x509"
	"os"
	"testing"

	. "github.com/smartystreets/goconvey/convey"
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
		So(tlscfg.ClientCAs.Subjects(), ShouldResemble, clientpool.Subjects())
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
		So(tlscfg.RootCAs.Subjects(), ShouldResemble, syspool.Subjects())
	})

	Convey("tls config with custom CA", t, func() {
		cfg := MongoConf{
			MongoTLSCA: "fixtures/ca-cert.pem",
		}
		tlscfg, err := cfg.TLSConfig()
		So(err, ShouldBeNil)
		So(tlscfg, ShouldNotBeNil)
		So(tlscfg.RootCAs.Subjects(), ShouldResemble, cuspool.Subjects())
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
		So(tlscfg.RootCAs.Subjects(), ShouldResemble, syspool.Subjects())
	})

	Convey("tls config with custom CA", t, func() {
		cfg := NATSConf{
			NATSTLSCA: "fixtures/ca-cert.pem",
		}
		tlscfg, err := cfg.TLSConfig()
		So(err, ShouldBeNil)
		So(tlscfg, ShouldNotBeNil)
		So(tlscfg.RootCAs.Subjects(), ShouldResemble, cuspool.Subjects())
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
