package hasher

import (
	"fmt"

	"github.com/spaolacci/murmur3"
	"go.aporeto.io/a3s/pkgs/api"
	"go.aporeto.io/a3s/pkgs/sharder"
	"go.aporeto.io/elemental"
)

// A Hasher computes the zone and zhash of
// identifiables and is used by a sharder.Sharder.
type Hasher struct{}

// Zone returns the zone for the given identity.
func (t *Hasher) Zone(identity elemental.Identity) int {
	return 0
}

// Hash computes and sets the zone and zhash for the given
// sharder.Shardable.
func (t *Hasher) Hash(z sharder.Shardable) error {

	z.SetZone(t.Zone(z.Identity()))

	switch oo := z.(type) {

	case *api.Namespace:
		z.SetZHash(hash(oo.Name))
	case *api.SparseNamespace:
		z.SetZHash(hash(*oo.Name))

	case *api.MTLSSource:
		z.SetZHash(hash(fmt.Sprintf("%s:%s", oo.Namespace, oo.Name)))
	case *api.SparseMTLSSource:
		z.SetZHash(hash(fmt.Sprintf("%s:%s", *oo.Namespace, *oo.Name)))

	case *api.LDAPSource:
		z.SetZHash(hash(fmt.Sprintf("%s:%s", oo.Namespace, oo.Name)))
	case *api.SparseLDAPSource:
		z.SetZHash(hash(fmt.Sprintf("%s:%s", *oo.Namespace, *oo.Name)))

	default:
		z.SetZHash(hash(oo.Identifier()))
	}

	return nil
}

func hash(v string) int {
	return int(murmur3.Sum64([]byte(v)) & 0x7FFFFFFFFFFFFFFF)
}
