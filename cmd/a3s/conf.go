package main

import (
	"fmt"

	"go.aporeto.io/a3s/internal/conf"
	"go.aporeto.io/a3s/internal/srv/authn"
	"go.aporeto.io/a3s/internal/srv/policy"
	"go.aporeto.io/addedeffect/lombric"
)

// Conf holds the main configuration flags.
type Conf struct {
	AuthNConf  authn.Conf  `mapstructure:",squash"`
	PolicyConf policy.Conf `mapstructure:",squash"`

	Init               bool   `mapstructure:"init-root-permissions"  desc:"If set, initialize the root permissions usingf the CA passed in --init-root-ca and exit"`
	InitRootUserCAPath string `mapstructure:"init-root-ca"           desc:"Path to the root CA to use to initialize root permissions"`

	conf.APIServerConf       `mapstructure:",squash"`
	conf.HealthConfiguration `mapstructure:",squash"`
	conf.HTTPTimeoutsConf    `mapstructure:",squash"`
	conf.LoggingConf         `mapstructure:",squash"`
	conf.NATSConf            `mapstructure:",squash"`
	conf.ProfilingConf       `mapstructure:",squash"`
	conf.RateLimitingConf    `mapstructure:",squash"`
	conf.MongoConf           `mapstructure:",squash" override:"mongo-db=a3s"`
}

// Prefix returns the configuration prefix.
func (c *Conf) Prefix() string { return "a3s" }

// PrintVersion prints the current version.
func (c *Conf) PrintVersion() {
	fmt.Printf("authn 0.0.1")
}

func newConf() Conf {
	c := Conf{}
	lombric.Initialize(&c)
	return c
}
