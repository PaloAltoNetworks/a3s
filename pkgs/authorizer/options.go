package authorizer

import "go.aporeto.io/a3s/pkgs/permissions"

type config struct {
	ignoredResources     []string
	operationTransformer OperationTransformer
}

// An Option can be used to configure various options in the Authorizer.
type Option func(*config)

// OptionIgnoredResources sets the list of identities that should skip authorizations.
func OptionIgnoredResources(identities ...string) Option {
	return func(cfg *config) {
		cfg.ignoredResources = identities
	}
}

// OptionOperationTransformer sets operation transformer to apply to each operation.
func OptionOperationTransformer(t OperationTransformer) Option {
	return func(cfg *config) {
		cfg.operationTransformer = t
	}
}

type checkConfig struct {
	sourceIP     string
	id           string
	restrictions permissions.Restrictions
}

// An OptionCheck can be used to configure various options when calling CheckPermissions.
type OptionCheck func(*checkConfig)

// OptionCheckSourceIP sets source IP of the request.
func OptionCheckSourceIP(ip string) OptionCheck {
	return func(cfg *checkConfig) {
		cfg.sourceIP = ip
	}
}

// OptionCheckID sets source IP of the request.
func OptionCheckID(id string) OptionCheck {
	return func(cfg *checkConfig) {
		cfg.id = id
	}
}

// OptionCheckRestrictions sets source restrictions to apply.
func OptionCheckRestrictions(r permissions.Restrictions) OptionCheck {
	return func(cfg *checkConfig) {
		cfg.restrictions = r
	}
}
