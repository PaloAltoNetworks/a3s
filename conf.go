package main

import (
	"fmt"

	"go.aporeto.io/a3s/internal/conf"
	"go.aporeto.io/a3s/srv/authn"
	"go.aporeto.io/addedeffect/lombric"
)

// Conf holds the main configuration flags.
type Conf struct {
	AuthNConf authn.Conf `mapstructure:",squash"`

	conf.APIServerConf       `mapstructure:",squash"`
	conf.HealthConfiguration `mapstructure:",squash"`
	conf.HTTPTimeoutsConf    `mapstructure:",squash"`
	conf.LoggingConf         `mapstructure:",squash"`
	conf.NATSConf            `mapstructure:",squash"`
	conf.ProfilingConf       `mapstructure:",squash"`
	conf.RateLimitingConf    `mapstructure:",squash"`
	conf.MongoConf           `mapstructure:",squash" override:"mongo-db=authn"`
}

// Prefix returns the configuration prefix.
func (c *Conf) Prefix() string { return "authn" }

// PrintVersion prints the current version.
func (c *Conf) PrintVersion() {
	fmt.Printf("authn 0.0.1")
}

func newConf() Conf {
	c := Conf{}
	lombric.Initialize(&c)
	return c
}
