package jobs

import (
	"context"
	"fmt"
	"testing"
	"time"

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
		makeTimePr := func(t time.Time) *time.Time { return &t }

		now := time.Now()

		Convey("everything works fine", func() {

			var m1ExpectedRecursive bool
			var m1ExpectedFields []string
			var m1ExpectedOrder []string
			m1.MockRetrieveMany(t, func(mctx manipulate.Context, dest elemental.Identifiables) error {

				m1ExpectedRecursive = mctx.Recursive()
				m1ExpectedFields = mctx.Fields()
				m1ExpectedOrder = mctx.Order()

				*dest.(*api.SparseNamespaceDeletionRecordsList) = append(
					*dest.(*api.SparseNamespaceDeletionRecordsList),
					&api.SparseNamespaceDeletionRecord{
						ID:         makeStrPr("1"),
						Namespace:  makeStrPr("/a"),
						DeleteTime: makeTimePr(now),
					},
					&api.SparseNamespaceDeletionRecord{
						ID:         makeStrPr("2"),
						Namespace:  makeStrPr("/a/1"),
						DeleteTime: makeTimePr(now),
					},
					&api.SparseNamespaceDeletionRecord{
						ID:         makeStrPr("3"),
						Namespace:  makeStrPr("/b"),
						DeleteTime: makeTimePr(now),
					},
				)
				return nil
			})

			var m2ExpectedFiler *elemental.Filter
			var m2ExpectedCalls int
			m2.MockDeleteMany(t, func(mctx manipulate.Context, identity elemental.Identity) error {
				m2ExpectedCalls++
				m2ExpectedFiler = mctx.Filter()
				return nil
			})

			err := DeleteOrphanedObjects(
				context.Background(),
				m1,
				m2,
				[]elemental.Identity{
					testmodel.ListIdentity,
					testmodel.TaskIdentity,
					api.NamespaceDeletionRecordIdentity,
				},
			)

			So(err, ShouldBeNil)
			So(m1ExpectedRecursive, ShouldBeTrue)
			So(m1ExpectedFields, ShouldResemble, []string{"namespace", "deleteTime"})
			So(m1ExpectedOrder, ShouldResemble, []string{"ID"})

			So(m2ExpectedCalls, ShouldEqual, 6)
			So(
				m2ExpectedFiler,
				ShouldResemble,
				elemental.NewFilterComposer().And(
					manipulate.NewNamespaceFilter("/b", true),
					elemental.NewFilterComposer().Or(
						elemental.NewFilterComposer().WithKey("createTime").NotExists().Done(),
						elemental.NewFilterComposer().WithKey("createTime").LesserThan(now).Done(),
					).Done(),
				).Done(),
			)
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

		Convey("no deletion records exist", func() {

			var m1ExpectedRecursive bool
			var m1ExpectedFields []string
			var m1ExpectedOrder []string
			m1.MockRetrieveMany(t, func(mctx manipulate.Context, dest elemental.Identifiables) error {

				m1ExpectedRecursive = mctx.Recursive()
				m1ExpectedFields = mctx.Fields()
				m1ExpectedOrder = mctx.Order()

				return nil
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

			So(err, ShouldBeNil)
			So(m1ExpectedRecursive, ShouldBeTrue)
			So(m1ExpectedFields, ShouldResemble, []string{"namespace", "deleteTime"})
			So(m1ExpectedOrder, ShouldResemble, []string{"ID"})
		})

		Convey("When m2 returns an error", func() {

			m1.MockRetrieveMany(t, func(mctx manipulate.Context, dest elemental.Identifiables) error {
				*dest.(*api.SparseNamespaceDeletionRecordsList) = append(
					*dest.(*api.SparseNamespaceDeletionRecordsList),
					&api.SparseNamespaceDeletionRecord{
						ID:         makeStrPr("1"),
						Namespace:  makeStrPr("/a"),
						DeleteTime: makeTimePr(now),
					},
					&api.SparseNamespaceDeletionRecord{
						ID:         makeStrPr("2"),
						Namespace:  makeStrPr("/a/1"),
						DeleteTime: makeTimePr(now),
					},
					&api.SparseNamespaceDeletionRecord{
						ID:         makeStrPr("3"),
						Namespace:  makeStrPr("/b"),
						DeleteTime: makeTimePr(now),
					},
				)
				return nil
			})

			m2.MockDeleteMany(t, func(mctx manipulate.Context, identity elemental.Identity) error {
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
			So(err.Error(), ShouldEqual, "unable to deletemany 'lists' in namespace '/a': bim")
		})
	})
}
