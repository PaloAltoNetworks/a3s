package gwutils

import (
	"crypto"
	"crypto/tls"
	"crypto/x509"
	"crypto/x509/pkix"
	"encoding/pem"
	"fmt"
	"net/http"
	"testing"
	"time"

	. "github.com/smartystreets/goconvey/convey"
	"go.aporeto.io/a3s/pkgs/api"
	"go.aporeto.io/bahamut/gateway"
	"go.aporeto.io/elemental"
	"go.aporeto.io/manipulate"
	"go.aporeto.io/manipulate/maniptest"
	"go.aporeto.io/tg/tglib"
)

func getECCert() (*x509.Certificate, crypto.PrivateKey) {

	certBlock, keyBlock, err := tglib.Issue(pkix.Name{})
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

func TestMakeTLSVerifyPeerCertificate(t *testing.T) {

	Convey("Given a cert and a manipulator", t, func() {

		cert, _ := getECCert()
		block, _ := tglib.CertToPEM(cert)
		wrongCert, _ := getECCert()
		wrongBlock, _ := tglib.CertToPEM(wrongCert)

		m := maniptest.NewTestManipulator()

		verifier := MakeTLSVerifyPeerCertificate(m, 10, time.Second)

		Convey("When no cert are passed, nothing should happen", func() {
			err := verifier(nil, nil)
			So(err, ShouldBeNil)
		})

		Convey("When given certs are garbage, it should fail", func() {
			err := verifier([][]byte{[]byte("garbage")}, nil)
			So(err, ShouldNotBeNil)
			So(err.Error(), ShouldEqual, "tls: failed to parse certificate from server: x509: malformed certificate")
		})

		Convey("When the manipulator returns an error", func() {

			m.MockRetrieveMany(t, func(mctx manipulate.Context, dest elemental.Identifiables) error {
				return fmt.Errorf("boom")
			})

			err := verifier([][]byte{block.Bytes}, nil)
			So(err, ShouldNotBeNil)
			So(err.Error(), ShouldEqual, "unable to retrieve mtlssources: boom")
		})

		Convey("Calling with a certificate that matches an exising mtls source should work", func() {

			called := 0
			m.MockRetrieveMany(t, func(mctx manipulate.Context, dest elemental.Identifiables) error {
				called++
				*dest.(*api.MTLSSourcesList) = append(
					*dest.(*api.MTLSSourcesList),
					&api.MTLSSource{
						CA: string(pem.EncodeToMemory(block)),
					},
				)
				return nil
			})

			err := verifier([][]byte{block.Bytes}, nil)
			So(err, ShouldBeNil)
			So(called, ShouldEqual, 1)

			Convey("When I call it again, it should work from the cache", func() {
				err := verifier([][]byte{block.Bytes}, nil)
				So(err, ShouldBeNil)
				So(called, ShouldEqual, 1)
			})
		})

		Convey("Calling with a certificate that matches an exising mtls with a wrong cert should fail", func() {

			m.MockRetrieveMany(t, func(mctx manipulate.Context, dest elemental.Identifiables) error {
				*dest.(*api.MTLSSourcesList) = append(
					*dest.(*api.MTLSSourcesList),
					&api.MTLSSource{
						CA: string(pem.EncodeToMemory(wrongBlock)),
					},
				)
				return nil
			})

			err := verifier([][]byte{block.Bytes}, nil)
			So(err, ShouldNotBeNil)
			So(err.Error(), ShouldStartWith, "unable to validate client certificate: x509: certificate signed by unknown authority")
		})

		Convey("Calling with a certificate that does not matche any exising mtls source should fail", func() {

			m.MockRetrieveMany(t, func(mctx manipulate.Context, dest elemental.Identifiables) error {
				return nil
			})

			err := verifier([][]byte{block.Bytes}, nil)
			So(err, ShouldNotBeNil)
			So(err.Error(), ShouldEqual, "unable to retrieve any matching mtlssource")
		})

		Convey("When more than one source match, it should fail", func() {

			m.MockRetrieveMany(t, func(mctx manipulate.Context, dest elemental.Identifiables) error {

				*dest.(*api.MTLSSourcesList) = append(
					*dest.(*api.MTLSSourcesList),
					&api.MTLSSource{
						CA: string(pem.EncodeToMemory(block)),
					},
					&api.MTLSSource{
						CA: string(pem.EncodeToMemory(block)),
					},
				)
				return nil
			})

			err := verifier([][]byte{block.Bytes}, nil)
			So(err, ShouldNotBeNil)
			So(err.Error(), ShouldEqual, "more than one mtls sources hold the signing CA. this is not supported")
		})
	})
}

func TestMakeTLSPeerCertificateForwarder(t *testing.T) {

	cert, _ := getECCert()

	Convey("Given I have a passphrase and an interceptor and no peer cert", t, func() {
		f := MakeTLSPeerCertificateForwarder("1234567890abcdef")
		action, upstream, err := f(nil, &http.Request{TLS: &tls.ConnectionState{}}, nil, nil)
		So(action, ShouldEqual, gateway.InterceptorActionForward)
		So(upstream, ShouldBeEmpty)
		So(err, ShouldBeNil)
	})

	Convey("A bad certificate should be handled", t, func() {
		f := MakeTLSPeerCertificateForwarder("1234567890abcdef")
		action, upstream, err := f(
			nil,
			&http.Request{
				TLS: &tls.ConnectionState{
					PeerCertificates: []*x509.Certificate{nil},
				},
			},
			nil,
			nil,
		)
		So(action, ShouldEqual, gateway.InterceptorActionStop)
		So(upstream, ShouldBeEmpty)
		So(err, ShouldNotBeNil)
		So(err.Error(), ShouldEqual, "nil certificate provided")
	})

	Convey("A bad encryption pass should be handled", t, func() {
		f := MakeTLSPeerCertificateForwarder("oh no")
		action, upstream, err := f(
			nil,
			&http.Request{
				TLS: &tls.ConnectionState{
					PeerCertificates: []*x509.Certificate{cert},
				},
			},
			nil,
			nil,
		)
		So(action, ShouldEqual, gateway.InterceptorActionStop)
		So(upstream, ShouldBeEmpty)
		So(err, ShouldNotBeNil)
		So(err.Error(), ShouldEqual, "invalid passphrase: size must be exactly 16 bytes")
	})

	Convey("When everything is fine, the header should be provided", t, func() {
		f := MakeTLSPeerCertificateForwarder("1234567890abcdef")
		req := &http.Request{
			Header: http.Header{},
			TLS: &tls.ConnectionState{
				PeerCertificates: []*x509.Certificate{cert},
			},
		}
		action, upstream, err := f(
			nil,
			req,
			nil,
			nil,
		)
		So(action, ShouldEqual, gateway.InterceptorActionForward)
		So(upstream, ShouldBeEmpty)
		So(err, ShouldBeEmpty)
		So(req.Header.Get("X-TLS-Certificate"), ShouldNotBeEmpty)
	})
}
