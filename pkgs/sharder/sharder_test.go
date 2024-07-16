package sharder

import (
	"context"
	"fmt"
	"reflect"
	"testing"

	. "github.com/smartystreets/goconvey/convey"
	"github.com/spaolacci/murmur3"
	"go.aporeto.io/a3s/pkgs/api"
	"go.aporeto.io/elemental"
	testmodel "go.aporeto.io/elemental/test/model"
	"go.aporeto.io/manipulate"
	"go.aporeto.io/manipulate/manipmongo"
	"go.aporeto.io/manipulate/maniptest"
	"go.mongodb.org/mongo-driver/bson"
)

type testHasher struct{}

func (t *testHasher) Zone(identity elemental.Identity) int {
	return 0
}

func (t *testHasher) Hash(z Shardable) error {

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

func TestShard(t *testing.T) {

	Convey("Given I have a sharder", t, func() {

		s := New(&testHasher{})

		So(s.OnShardedWrite(nil, nil, elemental.OperationCreate, nil), ShouldBeNil)

		aString := "a-name"

		Convey("Then sharding an a non zonable object", func() {
			o := testmodel.NewList()
			o.Name = aString
			So(s.Shard(nil, nil, o), ShouldBeNil)
		})

		Convey("Then sharding a Namespace should work", func() {
			o := api.NewNamespace()
			o.Name = aString
			So(s.Shard(nil, nil, o), ShouldBeNil)
			So(o.Zone, ShouldEqual, 0)
			So(o.ZHash, ShouldEqual, hash(aString))

			so := api.NewSparseNamespace()
			so.Name = &aString
			So(s.Shard(nil, nil, so), ShouldBeNil)
			So(*so.Zone, ShouldEqual, 0)
			So(*so.ZHash, ShouldEqual, hash(aString))
		})

		Convey("Then sharding an Authorization should work", func() {
			o := api.NewAuthorization()
			o.ID = aString
			So(s.Shard(nil, nil, o), ShouldBeNil)
			So(o.Zone, ShouldEqual, 0)
			So(o.ZHash, ShouldEqual, hash(aString))

			so := api.NewSparseAuthorization()
			so.ID = &aString
			So(s.Shard(nil, nil, so), ShouldBeNil)
			So(*so.Zone, ShouldEqual, 0)
			So(*so.ZHash, ShouldEqual, hash(aString))
		})

		Convey("Then sharding an MTLSSource should work", func() {
			o := api.NewMTLSSource()
			o.Namespace = aString
			o.Name = aString
			So(s.Shard(nil, nil, o), ShouldBeNil)
			So(o.Zone, ShouldEqual, 0)
			So(o.ZHash, ShouldEqual, hash(fmt.Sprintf("%s:%s", aString, aString)))

			so := api.NewSparseMTLSSource()
			so.Namespace = &aString
			so.Name = &aString
			So(s.Shard(nil, nil, so), ShouldBeNil)
			So(*so.Zone, ShouldEqual, 0)
			So(o.ZHash, ShouldEqual, hash(fmt.Sprintf("%s:%s", aString, aString)))
		})

		Convey("Then sharding an LDAPSource should work", func() {
			o := api.NewLDAPSource()
			o.Namespace = aString
			o.Name = aString
			So(s.Shard(nil, nil, o), ShouldBeNil)
			So(o.Zone, ShouldEqual, 0)
			So(o.ZHash, ShouldEqual, hash(fmt.Sprintf("%s:%s", aString, aString)))

			so := api.NewSparseLDAPSource()
			so.Namespace = &aString
			so.Name = &aString
			So(s.Shard(nil, nil, so), ShouldBeNil)
			So(*so.Zone, ShouldEqual, 0)
			So(o.ZHash, ShouldEqual, hash(fmt.Sprintf("%s:%s", aString, aString)))
		})
	})
}

func TestFilterOne(t *testing.T) {
	type args struct {
		m    manipulate.TransactionalManipulator
		mctx manipulate.Context
		o    elemental.Identifiable
	}
	tests := []struct {
		name    string
		s       *sharder
		args    args
		want    bson.D
		wantErr bool
	}{
		{
			"zonable with zhash",
			&sharder{hasher: &testHasher{}},
			args{
				maniptest.NewTestManipulator(),
				manipulate.NewContext(context.Background()),
				&api.Namespace{
					Zone:  2, // should be reset to 0 no matter what
					ZHash: 43,
				},
			},
			bson.D{{Key: "zone", Value: 0}, {Key: "zhash", Value: 43}},
			false,
		},
		{
			"zonable with no zhash",
			&sharder{hasher: &testHasher{}},
			args{
				maniptest.NewTestManipulator(),
				manipulate.NewContext(context.Background()),
				&api.Namespace{
					Zone:  2, // should be reset to 0 no matter what
					ZHash: 0,
				},
			},
			bson.D{{Key: "zone", Value: 0}},
			false,
		},
		{
			"zonable with zhash and mongo upsert for an Identifiable with custom sharding zhash",
			&sharder{hasher: &testHasher{}},
			args{
				maniptest.NewTestManipulator(),
				manipulate.NewContext(context.Background(), manipmongo.ContextOptionUpsert(nil)),
				&api.Namespace{
					Zone:  0, // should be reset to 0 no matter what
					ZHash: 43,
					ID:    "abcd",
					Name:  "abcd",
				},
			},
			bson.D{{Key: "zone", Value: 0}, {Key: "zhash", Value: 43}},
			false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			s := &sharder{hasher: &testHasher{}}
			got, err := s.FilterOne(tt.args.m, tt.args.mctx, tt.args.o)
			if (err != nil) != tt.wantErr {
				t.Errorf("sharder.FilterOne() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if !reflect.DeepEqual(got, tt.want) {
				t.Errorf("sharder.FilterOne() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestFilterMany(t *testing.T) {
	type args struct {
		m        manipulate.TransactionalManipulator
		mctx     manipulate.Context
		identity elemental.Identity
	}
	tests := []struct {
		name    string
		s       *sharder
		args    args
		want    bson.D
		wantErr bool
	}{
		{
			"z0",
			&sharder{hasher: &testHasher{}},
			args{
				nil,
				nil,
				elemental.Identity{},
			},
			bson.D{{Key: "zone", Value: 0}},
			false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			s := &sharder{hasher: &testHasher{}}
			got, err := s.FilterMany(tt.args.m, tt.args.mctx, tt.args.identity)
			if (err != nil) != tt.wantErr {
				t.Errorf("sharder.FilterMany() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if !reflect.DeepEqual(got, tt.want) {
				t.Errorf("sharder.FilterMany() = %v, want %v", got, tt.want)
			}
		})
	}
}
