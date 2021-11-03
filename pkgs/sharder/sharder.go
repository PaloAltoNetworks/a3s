package sharder

import (
	"errors"

	"github.com/globalsign/mgo/bson"
	"go.aporeto.io/elemental"
	"go.aporeto.io/manipulate"
	"go.aporeto.io/manipulate/manipmongo"
)

var _ manipmongo.Sharder = &sharder{}

// A Sharder is responsible for computing sharding.
type sharder struct {
}

// New returns a new manipmongo.Sharder.
func New() manipmongo.Sharder {
	return &sharder{}
}

// Shard implements the manipmongo.Sharder interface.
func (s *sharder) Shard(m manipulate.TransactionalManipulator, mctx manipulate.Context, o elemental.Identifiable) error {

	if err := ApplyZoning(o); !errors.Is(err, ErrNotZonable) {
		return err
	}

	return nil
}

func (s *sharder) OnShardedWrite(m manipulate.TransactionalManipulator, mctx manipulate.Context, op elemental.Operation, o elemental.Identifiable) error {
	return nil
}

// FilterOne implements the manipmongo.Sharder interface.
func (s *sharder) FilterOne(m manipulate.TransactionalManipulator, mctx manipulate.Context, o elemental.Identifiable) (bson.D, error) {

	z, ok := o.(Shardable)
	if !ok || z.GetZHash() == 0 {
		return bson.D{
			{Name: "zone", Value: ZoneForIdentity(o.Identity())},
		}, nil
	}

	return bson.D{
		{Name: "zone", Value: ZoneForIdentity(o.Identity())},
		{Name: "zhash", Value: z.GetZHash()},
	}, nil
}

// FilterMany implements the manipmongo.Sharder interface.
func (s *sharder) FilterMany(m manipulate.TransactionalManipulator, mctx manipulate.Context, identity elemental.Identity) (bson.D, error) {

	return bson.D{
		{Name: "zone", Value: ZoneForIdentity(identity)},
	}, nil
}
