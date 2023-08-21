package bootstrap

import (
	"crypto/x509"

	"go.aporeto.io/a3s/pkgs/token"
	"go.aporeto.io/bahamut"
	"go.aporeto.io/bahamut/authorizer/mtls"
	"go.aporeto.io/bahamut/authorizer/simple"
	"go.aporeto.io/elemental"
)

// MakeCNExcluderFunc returns an mtls.VerifierFunc that rejects mtls verification
// if the certificate has the given Common Name.
func MakeCNExcluderFunc(cn string) mtls.VerifierFunc {
	return func(cert *x509.Certificate) bool { return cert.Subject.CommonName != cn }
}

// MakeMTLSRequestAuthenticator returns a bahamut.RequestAuthenticator that
// will allow requests presenting a client certificate signed by a CA present
// in the given caPool unless the given mtls.VerifierFunc disagrees.
func MakeMTLSRequestAuthenticator(caPool *x509.CertPool, verifier mtls.VerifierFunc) bahamut.RequestAuthenticator {
	return mtls.NewMTLSRequestAuthenticator(defaultOption(caPool, verifier))
}

// MakeMTLSAuthorizer returns a bahamut.Authorizer that will allow requests
// presenting a client certificate signed by a CA present in the given caPool
// unless the given mtls.VerifierFunc disagrees.
func MakeMTLSAuthorizer(caPool *x509.CertPool, verifier mtls.VerifierFunc, ignoredIdentities []elemental.Identity) bahamut.Authorizer {
	v, d, mf, mc := defaultOption(caPool, verifier)
	return mtls.NewMTLSAuthorizer(v, d, ignoredIdentities, mf, mc)
}

// MakeRequestAuthenticatorBypasser returns a bahamut.RequestAuthenticator that returns bahamut.ActionOK
// if the request is for one of the given identities. Otherwise it returns bahamut.Continue.
func MakeRequestAuthenticatorBypasser(identities []elemental.Identity) bahamut.RequestAuthenticator {
	return simple.NewAuthenticator(makeBypasserFunc(identities), nil)
}

// MakeSessionAuthenticatorBypasser returns a bahamut.SessionAuthenticator that returns bahamut.ActionOK
// if the request is for one of the given identities. Otherwise it returns bahamut.Continue.
func MakeSessionAuthenticatorBypasser(identities []elemental.Identity) bahamut.RequestAuthenticator {
	return simple.NewAuthenticator(makeBypasserFunc(identities), nil)
}

// MakeAuthorizerBypasser returns a bahamut.Authorizer that returns bahamut.ActionOK
// if the request is on one of the given identities. Otherwise it returns bahamut.Continue.
func MakeAuthorizerBypasser(identities []elemental.Identity) bahamut.Authorizer {
	return simple.NewAuthorizer(makeBypasserFunc(identities))
}

func makeBypasserFunc(identities []elemental.Identity) func(ctx bahamut.Context) (bahamut.AuthAction, error) {

	return func(ctx bahamut.Context) (bahamut.AuthAction, error) {
		for _, i := range identities {
			if ctx.Request().Identity.IsEqual(i) {
				return bahamut.AuthActionOK, nil
			}
		}
		return bahamut.AuthActionContinue, nil
	}
}

func defaultOption(caPool *x509.CertPool, verifier mtls.VerifierFunc) (x509.VerifyOptions, mtls.DeciderFunc, mtls.VerifierFunc, mtls.CertificateCheckMode) {

	return x509.VerifyOptions{
			Roots:     caPool,
			KeyUsages: []x509.ExtKeyUsage{x509.ExtKeyUsageClientAuth},
		},
		func(a bahamut.AuthAction, c bahamut.Context, s bahamut.Session) bahamut.AuthAction {
			switch a {
			case bahamut.AuthActionKO:
				return bahamut.AuthActionContinue
			case bahamut.AuthActionOK:
				if token.FromRequest(c.Request()) != "" {
					return bahamut.AuthActionContinue
				}
				return bahamut.AuthActionOK
			default:
				panic("should not reach here")
			}
		},
		verifier,
		mtls.CertificateCheckModeHeaderOnly
}
