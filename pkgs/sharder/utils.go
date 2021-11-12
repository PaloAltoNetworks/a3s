package sharder

import (
	"errors"
	"fmt"

	"github.com/spaolacci/murmur3"
	"go.aporeto.io/a3s/pkgs/api"
	"go.aporeto.io/elemental"
)

var (
	// ErrNotZonable indicates that the object is Shardable and
	// cannot be sharded.
	ErrNotZonable = errors.New("given object is not shardable")
)

// A Shardable is the interface an object must implement
// in order to be shardable.
type Shardable interface {
	GetZone() int
	SetZone(int)
	GetZHash() int
	SetZHash(int)
}

// ZoneForIdentity returns the sharding zone for the given identity.
func ZoneForIdentity(identity elemental.Identity) int {
	return 0
}

// ApplyZoning applies sharding on a elemental object.
// If the object is not a Shardabled, ErrNotZonable will be returned.
func ApplyZoning(o elemental.Identifiable) error {

	z, ok := o.(Shardable)
	if !ok {
		return ErrNotZonable
	}

	z.SetZone(ZoneForIdentity(o.Identity()))

	switch oo := o.(type) {

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
