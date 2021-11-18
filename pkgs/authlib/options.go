package authlib

import "time"

type config struct {
	validity              time.Duration
	opaque                map[string]string
	audience              []string
	restrictedNamespace   string
	restrictedPermissions []string
	restrictedNetworks    []string
	cloak                 []string
}

func newConfig() config {
	return config{
		validity: 1 * time.Hour,
	}
}

// An Option is the type of various options
// You can add the issue requests.
type Option func(*config)

// OptValidity sets the validity to request for the token.
func OptValidity(validity time.Duration) Option {

	return func(opts *config) {
		opts.validity = validity
	}
}

// OptCloak sets the claims cloaking option for the token.
func OptCloak(cloaking ...string) Option {

	return func(opts *config) {
		opts.cloak = cloaking
	}
}

// OptOpaque passes opaque data that will be
// included in the JWT.
func OptOpaque(opaque map[string]string) Option {

	return func(opts *config) {
		opts.opaque = opaque
	}
}

// OptAudience passes the requested audience for the token.
func OptAudience(audience ...string) Option {

	return func(opts *config) {
		opts.audience = audience
	}
}

// OptRestrictNamespace asks for a restricted token on the given namespace.
func OptRestrictNamespace(namespace string) Option {

	return func(opts *config) {
		opts.restrictedNamespace = namespace
	}
}

// OptRestrictPermissions asks for a restricted token on the given permissions.
func OptRestrictPermissions(permissions []string) Option {

	return func(opts *config) {
		opts.restrictedPermissions = permissions
	}
}

// OptRestrictNetworks asks for a restricted token on the given networks.
func OptRestrictNetworks(networks []string) Option {

	return func(opts *config) {
		opts.restrictedNetworks = networks
	}
}
