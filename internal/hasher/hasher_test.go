package hasher

import (
	"fmt"
	"testing"

	. "github.com/smartystreets/goconvey/convey"
	"go.aporeto.io/a3s/pkgs/api"
)

func TestShard(t *testing.T) {

	Convey("Given I have a hasher", t, func() {

		s := &Hasher{}

		aString := "a-name"

		Convey("Then sharding a Namespace should work", func() {
			o := api.NewNamespace()
			o.Name = aString
			So(s.Hash(o), ShouldBeNil)
			So(o.Zone, ShouldEqual, 0)
			So(o.ZHash, ShouldEqual, hash(aString))

			so := api.NewSparseNamespace()
			so.Name = &aString
			So(s.Hash(so), ShouldBeNil)
			So(*so.Zone, ShouldEqual, 0)
			So(*so.ZHash, ShouldEqual, hash(aString))
		})

		Convey("Then sharding an Authorization should work", func() {
			o := api.NewAuthorization()
			o.ID = aString
			So(s.Hash(o), ShouldBeNil)
			So(o.Zone, ShouldEqual, 0)
			So(o.ZHash, ShouldEqual, hash(aString))

			so := api.NewSparseAuthorization()
			so.ID = &aString
			So(s.Hash(so), ShouldBeNil)
			So(*so.Zone, ShouldEqual, 0)
			So(*so.ZHash, ShouldEqual, hash(aString))
		})

		Convey("Then sharding an MTLSSource should work", func() {
			o := api.NewMTLSSource()
			o.Namespace = aString
			o.Name = aString
			So(s.Hash(o), ShouldBeNil)
			So(o.Zone, ShouldEqual, 0)
			So(o.ZHash, ShouldEqual, hash(fmt.Sprintf("%s:%s", aString, aString)))

			so := api.NewSparseMTLSSource()
			so.Namespace = &aString
			so.Name = &aString
			So(s.Hash(so), ShouldBeNil)
			So(*so.Zone, ShouldEqual, 0)
			So(o.ZHash, ShouldEqual, hash(fmt.Sprintf("%s:%s", aString, aString)))
		})

		Convey("Then sharding an LDAPSource should work", func() {
			o := api.NewLDAPSource()
			o.Namespace = aString
			o.Name = aString
			So(s.Hash(o), ShouldBeNil)
			So(o.Zone, ShouldEqual, 0)
			So(o.ZHash, ShouldEqual, hash(fmt.Sprintf("%s:%s", aString, aString)))

			so := api.NewSparseLDAPSource()
			so.Namespace = &aString
			so.Name = &aString
			So(s.Hash(so), ShouldBeNil)
			So(*so.Zone, ShouldEqual, 0)
			So(o.ZHash, ShouldEqual, hash(fmt.Sprintf("%s:%s", aString, aString)))
		})

		Convey("Then sharding an A3SSource should work", func() {
			o := api.NewA3SSource()
			o.Namespace = aString
			o.Name = aString
			So(s.Hash(o), ShouldBeNil)
			So(o.Zone, ShouldEqual, 0)
			So(o.ZHash, ShouldEqual, hash(fmt.Sprintf("%s:%s", aString, aString)))

			so := api.NewSparseA3SSource()
			so.Namespace = &aString
			so.Name = &aString
			So(s.Hash(so), ShouldBeNil)
			So(*so.Zone, ShouldEqual, 0)
			So(o.ZHash, ShouldEqual, hash(fmt.Sprintf("%s:%s", aString, aString)))
		})
	})
}
