package authlib

import (
	"time"

	"go.aporeto.io/a3s/pkgs/permissions"
)

type config struct {
	validity     time.Duration
	opaque       map[string]string
	audience     []string
	restrictions permissions.Restrictions
	cloak        []string
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

// OptRestrictions sets the request restrictions for the token.
func OptRestrictions(restrictions permissions.Restrictions) Option {
	return func(opts *config) {
		opts.restrictions = restrictions
	}
}
