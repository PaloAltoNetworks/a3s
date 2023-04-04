package jobs

import (
	"context"
	"fmt"
	"testing"

	. "github.com/smartystreets/goconvey/convey"
	"go.aporeto.io/a3s/pkgs/api"
	"go.aporeto.io/elemental"
	testmodel "go.aporeto.io/elemental/test/model"
	"go.aporeto.io/manipulate"
	"go.aporeto.io/manipulate/maniptest"
)

func TestDeleteOrphanedJobs(t *testing.T) {

	Convey("I have some manipulators and a model", t, func() {

		m1 := maniptest.NewTestManipulator()
		m2 := maniptest.NewTestManipulator()
		makeStrPr := func(str string) *string { return &str }

		Convey("everything works fine", func() {

			var m1ExpectedRecursive bool
			var m1ExpectedFields []string
			var m1ExpectedOrder []string
			m1.MockRetrieveMany(t, func(mctx manipulate.Context, dest elemental.Identifiables) error {

				m1ExpectedRecursive = mctx.Recursive()
				m1ExpectedFields = mctx.Fields()
				m1ExpectedOrder = mctx.Order()

				*dest.(*api.SparseNamespacesList) = append(
					*dest.(*api.SparseNamespacesList),
					&api.SparseNamespace{
						ID:   makeStrPr("1"),
						Name: makeStrPr("/a"),
					},
					&api.SparseNamespace{
						ID:   makeStrPr("2"),
						Name: makeStrPr("/a/1"),
					},
					&api.SparseNamespace{
						ID:   makeStrPr("3"),
						Name: makeStrPr("/b"),
					},
				)
				return nil
			})

			So(func() {
				_ = DeleteOrphanedObjects(
					context.Background(),
					m1,
					m2,
					[]elemental.Identity{
						testmodel.ListIdentity,
						testmodel.TaskIdentity,
					},
				)
			}, ShouldPanicWith, "you can only pass a mongo manipulator to GetDatabase")

			So(m1ExpectedRecursive, ShouldBeTrue)
			So(m1ExpectedFields, ShouldResemble, []string{"name"})
			So(m1ExpectedOrder, ShouldResemble, []string{"ID"})
		})

		Convey("When m1 returns an error", func() {

			m1.MockRetrieveMany(t, func(mctx manipulate.Context, dest elemental.Identifiables) error {
				return fmt.Errorf("bim")
			})

			err := DeleteOrphanedObjects(
				context.Background(),
				m1,
				m2,
				[]elemental.Identity{
					testmodel.ListIdentity,
					testmodel.TaskIdentity,
				},
			)

			So(err, ShouldNotBeNil)
			So(err.Error(), ShouldEqual, "unable to retrieve list of namespaces: unable to retrieve objects for iteration 1: bim")
		})
	})
}
