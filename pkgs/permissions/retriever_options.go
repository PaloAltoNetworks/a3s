package permissions

type config struct {
	id           string
	addr         string
	restrictions Restrictions
}

type RetrieverOption func(*config)

// OptionRetrieverID sets the ID to use to compute permissions.
func OptionRetrieverID(id string) RetrieverOption {
	return func(c *config) {
		c.id = id
	}
}

// OptionRetrieverIPAddr sets the source IP to use to compute permissions.
func OptionRetrieverIPAddr(ip string) RetrieverOption {
	return func(c *config) {
		c.addr = ip
	}
}

// OptionRetrieverRestrictions sets the restrictions to apply on the retrieved permissions.
func OptionRetrieverRestrictions(r Restrictions) RetrieverOption {
	return func(c *config) {
		c.restrictions = r
	}
}
