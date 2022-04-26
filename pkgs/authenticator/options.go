package authenticator

import "crypto/x509"

type config struct {
	ignoredResources       []string
	externalTrustedIssuers []RemoteIssuer
}

// An Option can be used to configure various options in the Authenticator.
type Option func(*config)

// OptionIgnoredResources sets the list of identities that should skip authentication.
func OptionIgnoredResources(identities ...string) Option {
	return func(cfg *config) {
		cfg.ignoredResources = identities
	}
}

// A RemoteIssuer holds the URL and the
// CertPool containing a CA to validate the server
type RemoteIssuer struct {
	URL  string
	Pool *x509.CertPool
}

// OptionExternalTrustedIssuers sets the list of additionally trusted issuers.
// This is to trust tokens from other a3s instances as valid and authenticated.
func OptionExternalTrustedIssuers(issuers ...RemoteIssuer) Option {
	return func(cfg *config) {
		cfg.externalTrustedIssuers = issuers
	}
}
