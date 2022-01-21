package gwutils

import (
	"context"
	"crypto/x509"
	"encoding/pem"
	"fmt"
	"net/http"
	"strings"

	"github.com/karlseguin/ccache/v2"
	"go.aporeto.io/a3s/pkgs/api"
	"go.aporeto.io/bahamut/gateway"
	"go.aporeto.io/elemental"
	"go.aporeto.io/manipulate"
	"go.aporeto.io/tg/tglib"
)

// MakeTLSPeerCertificateVerifier returns a function you can use as
// tls.Config.VerifyPeerCertificate. You will need to do this if you want to
// support user authentication through MTLS while you are behind a
// bahamut.Gateway.
//
// This is the first step of the necessary dance to securely forward the client
// certificate in a trusted header. You will then need to add an interceptor
// using MakeTLSPeerCertificateForwarder
//
// The returned function will use the provided manipulator to search A3S for an
// mtls source that holds the CA that has issued the presented client
// certificates by matching the certificate AuthorityKeyID. If it can find one,
// the certificate signature will be checked using the matching CA.
//
// The results are cached for the provided cacheDuration and a maximum of
// cacheSize items will be kept.
func MakeTLSPeerCertificateVerifier(
	ctx context.Context,
	m manipulate.Manipulator,
	opts ...VerifierOption,
) func([][]byte, [][]*x509.Certificate) error {

	cfg := newVerifierConf()
	for _, o := range opts {
		o(&cfg)
	}

	cache := ccache.New(ccache.Configure().MaxSize(cfg.cacheMaxSize))

	return func(
		rawCerts [][]byte,
		verifiedChains [][]*x509.Certificate,
	) error {

		if len(rawCerts) == 0 {
			return nil
		}

		cert, err := x509.ParseCertificate(rawCerts[0])
		if err != nil {
			return fmt.Errorf("tls: failed to parse certificate from server: %w", err)
		}

		authorityKeyId := fmt.Sprintf("%02X", cert.AuthorityKeyId)
		item := cache.Get(authorityKeyId)
		var pool *x509.CertPool

		if item == nil || item.Expired() {

			tctx, cancel := context.WithTimeout(ctx, cfg.timeout)
			defer cancel()

			mctx := manipulate.NewContext(
				tctx,
				manipulate.ContextOptionFilter(
					elemental.NewFilterComposer().
						WithKey("subjectKeyIDs").Equals(authorityKeyId).
						Done(),
				),
			)

			sources := api.MTLSSourcesList{}
			if err := m.RetrieveMany(mctx, &sources); err != nil {
				return fmt.Errorf("unable to retrieve mtlssources: %w", err)
			}

			switch len(sources) {
			case 1:
			case 0:
				return fmt.Errorf("unable to retrieve any matching mtlssource")
			default:
				return fmt.Errorf("more than one mtls sources hold the signing CA. this is not supported")
			}

			pool = x509.NewCertPool()
			pool.AppendCertsFromPEM([]byte(sources[0].CA))
			cache.Set(authorityKeyId, pool, cfg.cacheDuration)
		} else {
			pool = item.Value().(*x509.CertPool)
		}

		if _, err := cert.Verify(
			x509.VerifyOptions{
				Roots: pool,
				KeyUsages: []x509.ExtKeyUsage{
					x509.ExtKeyUsageClientAuth,
				},
			},
		); err != nil {
			return fmt.Errorf("unable to validate client certificate: %w", err)
		}

		return nil
	}
}

// MakeTLSPeerCertificateForwarder returns a bahamut gateway.InterceptorFunc
// that you will need to add to the bahamut.Gateway in order to intercept any
// calls going to the A3S /issue endpoint (or any other one you would have as a
// proxy) in order to pass the user certificates as a secure header.
//
// The encryptionPassphrase is necessary as A3S will refuse to trust a header
// containing a user certificate if it is not encrypted with that key. The key
// must be exactly 16, 24 or 32 bytes long to encrypt respectively to AES-128,
// A3S-192 or AES-256.
//
// WARNING: You MUST NOT use this function without installing a custom peer
// certificate verifier with MakeTLSPeerCertificateVerifier in the gateway's
// server TLS config. A3S will blindly trust the certificate in the header,
// which is why you MUST verify it before.
func MakeTLSPeerCertificateForwarder(encryptionPassphrase string) gateway.InterceptorFunc {

	return func(
		w http.ResponseWriter,
		req *http.Request,
		writeError gateway.ErrorWriter,
		corsInjector func(),
	) (gateway.InterceptorAction, string, error) {

		if len(req.TLS.PeerCertificates) == 0 {
			return gateway.InterceptorActionForward, "", nil
		}

		cert := req.TLS.PeerCertificates[0]

		block, err := tglib.CertToPEM(cert)
		if err != nil {
			return gateway.InterceptorActionStop, "", err
		}

		enc, err := elemental.NewAESAttributeEncrypter(encryptionPassphrase)
		if err != nil {
			return gateway.InterceptorActionStop, "", err
		}

		h, err := enc.EncryptString(
			strings.ReplaceAll(
				string(pem.EncodeToMemory(block)),
				"\n",
				" ",
			),
		)
		if err != nil {
			return gateway.InterceptorActionStop, "", err
		}

		req.Header.Set("X-TLS-Certificate", h)

		return gateway.InterceptorActionForward, "", nil
	}
}
