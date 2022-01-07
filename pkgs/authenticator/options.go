package authenticator

type config struct {
	ignoredResources []string
}

// An Option can be used to configure various options in the Authenticator.
type Option func(*config)

// OptionIgnoredResources sets the list of identities that should skip authentication.
func OptionIgnoredResources(identities ...string) Option {
	return func(cfg *config) {
		cfg.ignoredResources = identities
	}
}
