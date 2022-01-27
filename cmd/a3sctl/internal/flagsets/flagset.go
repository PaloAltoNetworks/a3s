package flagsets

import (
	"github.com/spf13/pflag"
)

// MakeAutoAuthFlags returns the flag set to handle auto auth
func MakeAutoAuthFlags() *pflag.FlagSet {
	fs := pflag.NewFlagSet("", pflag.ExitOnError)
	fs.Bool("renew-cached-token", false, "If set, the cached token will be refreshed")
	fs.String("auto-auth-method", "", "If set, override config's file autoauth.enable")
	return fs
}
