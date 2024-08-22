package sharder

import (
	"errors"

	"go.aporeto.io/elemental"
	"go.aporeto.io/manipulate"
	"go.aporeto.io/manipulate/manipmongo"
	"go.mongodb.org/mongo-driver/bson"
)

// ErrNotZonable indicates that the object is Shardable and
// cannot be sharded.
var ErrNotZonable = errors.New("given object is not shardable")

// A Shardable is the interface an object must implement
// in order to be shardable.
type Shardable interface {
	GetZone() int
	SetZone(int)
	GetZHash() int
	SetZHash(int)
	elemental.Identifiable
}

// A Hasher is used to zone/zhash an identifiable.
type Hasher interface {

	// Zone returns the zone for the given identity.
	Zone(elemental.Identity) int

	// Hash performs the zoning/hashing for the given identifiable.
	// This method is responsible to set the Zone and the ZHash of the
	// provided Identifiable.
	Hash(Shardable) error
}

var _ manipmongo.Sharder = &sharder{}

// A Sharder is responsible for computing sharding.
type sharder struct {
	hasher Hasher
}

// New returns a new manipmongo.Sharder.
func New(hasher Hasher) manipmongo.Sharder {
	return &sharder{
		hasher: hasher,
	}
}

// Shard implements the manipmongo.Sharder interface.
func (s *sharder) Shard(m manipulate.TransactionalManipulator, mctx manipulate.Context, o elemental.Identifiable) error {

	z, ok := o.(Shardable)
	if !ok {
		return nil
	}

	return s.hasher.Hash(z)
}

func (s *sharder) OnShardedWrite(m manipulate.TransactionalManipulator, mctx manipulate.Context, op elemental.Operation, o elemental.Identifiable) error {
	return nil
}

// FilterOne implements the manipmongo.Sharder interface.
func (s *sharder) FilterOne(m manipulate.TransactionalManipulator, mctx manipulate.Context, o elemental.Identifiable) (bson.D, error) {

	z, ok := o.(Shardable)
	if !ok || z.GetZHash() == 0 {
		return bson.D{
			{Key: "zone", Value: s.hasher.Zone(o.Identity())},
		}, nil
	}

	return bson.D{
		{Key: "zone", Value: s.hasher.Zone(o.Identity())},
		{Key: "zhash", Value: z.GetZHash()},
	}, nil
}

// FilterMany implements the manipmongo.Sharder interface.
func (s *sharder) FilterMany(m manipulate.TransactionalManipulator, mctx manipulate.Context, identity elemental.Identity) (bson.D, error) {

	return bson.D{
		{Key: "zone", Value: s.hasher.Zone(identity)},
	}, nil
}
