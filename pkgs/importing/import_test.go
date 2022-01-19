package importing

import (
	"context"
	"fmt"
	"net/http"
	"sort"
	"strings"
	"testing"

	. "github.com/smartystreets/goconvey/convey"
	"go.aporeto.io/a3s/pkgs/api"
	"go.aporeto.io/elemental"
	"go.aporeto.io/manipulate"
	"go.aporeto.io/manipulate/maniptest"
)

func TestImport(t *testing.T) {

	Convey("Given a manipulator", t, func() {

		m := maniptest.NewTestManipulator()
		_ = m

		Convey("When I call Import with missing label, it should error", func() {
			err := Import(
				context.Background(),
				api.Manager(),
				m,
				"",
				"",
				nil,
				false,
			)
			So(err, ShouldNotBeNil)
			So(err.Error(), ShouldEqual, "namespace must not be empty")
		})

		Convey("When I call Import with missing namespace, it should error", func() {
			err := Import(
				context.Background(),
				api.Manager(),
				m,
				"ns",
				"",
				nil,
				false,
			)
			So(err, ShouldNotBeNil)
			So(err.Error(), ShouldEqual, "label must not be empty")
		})

		Convey("When I pass a non Importable", func() {
			err := Import(
				context.Background(),
				api.Manager(),
				m,
				"ns",
				"label",
				api.NamespacesList{api.NewNamespace()},
				false,
			)
			So(err, ShouldNotBeNil)
			So(err.Error(), ShouldEqual, "object 'namespace[0]' is not importable")
		})

		Convey("When I pass a a list containing a nil manager", func() {
			err := Import(
				context.Background(),
				nil,
				m,
				"ns",
				"label",
				api.AuthorizationsList{nil},
				false,
			)
			So(err, ShouldNotBeNil)
			So(err.Error(), ShouldEqual, "unable to hash 'authorization[0]': manager must not be nil")
		})

		Convey("When I import, but retrieve many returns an error", func() {
			{

				m.MockRetrieveMany(t, func(mctx manipulate.Context, dest elemental.Identifiables) error {
					return fmt.Errorf("boom")
				})

				objs := api.AuthorizationsList{
					&api.Authorization{
						Name: "1",
					},
					&api.Authorization{
						Name:        "2",
						Description: "new",
					},
					&api.Authorization{
						Name: "4",
					},
				}

				err := Import(context.Background(), api.Manager(), m, "/ns", "label", objs, false)
				So(err, ShouldNotBeNil)
				So(err.Error(), ShouldEqual, "unable to retrieve list of current objects: boom")
			}
		})

		Convey("When I import a list of objects but delete returns an error", func() {

			m.MockRetrieveMany(t, func(mctx manipulate.Context, dest elemental.Identifiables) error {
				*dest.(*api.AuthorizationsList) = append(
					*dest.(*api.AuthorizationsList),
					&api.Authorization{
						ID:          "1",
						Name:        "1",
						ImportHash:  "3132303033343839333331383835343436343834e3b0c44298fc1c149afbf4c8996fb92427ae41e4649b934ca495991b7852b855",
						ImportLabel: "label",
					},
					&api.Authorization{
						ID:          "2",
						Name:        "2",
						ImportHash:  "3132363235373937303539373039393132333639e3b0c44298fc1c149afbf4c8996fb92427ae41e4649b934ca495991b7852b855",
						ImportLabel: "label",
					},
					&api.Authorization{
						ID:          "3",
						Name:        "3",
						ImportHash:  "3",
						ImportLabel: "label",
					},
				)
				return nil
			})

			m.MockDelete(t, func(mctx manipulate.Context, object elemental.Identifiable) error {
				return fmt.Errorf("paf")
			})

			objs := api.AuthorizationsList{
				&api.Authorization{
					Name: "1",
				},
				&api.Authorization{
					Name:        "2",
					Description: "new",
				},
				&api.Authorization{
					Name: "4",
				},
			}

			err := Import(context.Background(), api.Manager(), m, "/ns", "label", objs, false)
			So(err, ShouldNotBeNil)
			So(err.Error(), ShouldEqual, "unable to delete existing object: paf")
		})

		Convey("When I import a list of objects but creates returns an error", func() {

			m.MockRetrieveMany(t, func(mctx manipulate.Context, dest elemental.Identifiables) error {
				*dest.(*api.AuthorizationsList) = append(
					*dest.(*api.AuthorizationsList),
					&api.Authorization{
						ID:          "1",
						Name:        "1",
						ImportHash:  "3132303033343839333331383835343436343834e3b0c44298fc1c149afbf4c8996fb92427ae41e4649b934ca495991b7852b855",
						ImportLabel: "label",
					},
					&api.Authorization{
						ID:          "2",
						Name:        "2",
						ImportHash:  "3132363235373937303539373039393132333639e3b0c44298fc1c149afbf4c8996fb92427ae41e4649b934ca495991b7852b855",
						ImportLabel: "label",
					},
					&api.Authorization{
						ID:          "3",
						Name:        "3",
						ImportHash:  "3",
						ImportLabel: "label",
					},
				)
				return nil
			})

			m.MockDelete(t, func(mctx manipulate.Context, object elemental.Identifiable) error {
				return nil
			})

			m.MockCreate(t, func(mctx manipulate.Context, object elemental.Identifiable) error {
				return fmt.Errorf("bim")
			})

			objs := api.AuthorizationsList{
				&api.Authorization{
					Name: "1",
				},
				&api.Authorization{
					Name:        "2",
					Description: "new",
				},
				&api.Authorization{
					Name: "4",
				},
			}

			err := Import(context.Background(), api.Manager(), m, "/ns", "label", objs, false)
			So(err, ShouldNotBeNil)
			So(err.Error(), ShouldEqual, "unable to create imported object: bim")
		})

		Convey("When I import a list of objects and there are some existing", func() {

			m.MockRetrieveMany(t, func(mctx manipulate.Context, dest elemental.Identifiables) error {
				*dest.(*api.AuthorizationsList) = append(
					*dest.(*api.AuthorizationsList),
					&api.Authorization{
						ID:          "1",
						Name:        "1",
						ImportHash:  "3132303033343839333331383835343436343834e3b0c44298fc1c149afbf4c8996fb92427ae41e4649b934ca495991b7852b855",
						ImportLabel: "label",
					},
					&api.Authorization{
						ID:          "2",
						Name:        "2",
						ImportHash:  "3132363235373937303539373039393132333639e3b0c44298fc1c149afbf4c8996fb92427ae41e4649b934ca495991b7852b855",
						ImportLabel: "label",
					},
					&api.Authorization{
						ID:          "3",
						Name:        "3",
						ImportHash:  "3",
						ImportLabel: "label",
					},
				)
				return nil
			})

			toDelete := elemental.IdentifiablesList{}
			m.MockDelete(t, func(mctx manipulate.Context, object elemental.Identifiable) error {
				// fmt.Println("delete ID:", object.Identifier(), "hash:", object.(Importable).GetImportHash())
				toDelete = append(toDelete, object)
				return nil
			})

			toCreate := elemental.IdentifiablesList{}
			m.MockCreate(t, func(mctx manipulate.Context, object elemental.Identifiable) error {
				// fmt.Println("create ID:", object.Identifier(), "hash:", object.(Importable).GetImportHash())
				toCreate = append(toCreate, object)
				return nil
			})

			objs := api.AuthorizationsList{
				&api.Authorization{
					Name: "1",
				},
				&api.Authorization{
					Name:        "2",
					Description: "new",
				},
				&api.Authorization{
					Name: "4",
				},
			}

			err := Import(context.Background(), api.Manager(), m, "/ns", "label", objs, false)
			So(err, ShouldBeNil)

			sort.Slice(toDelete, func(i, j int) bool {
				return strings.Compare(toDelete[i].(*api.Authorization).Name, toDelete[j].(*api.Authorization).Name) != 1
			})
			sort.Slice(toCreate, func(i, j int) bool {
				return strings.Compare(toCreate[i].(*api.Authorization).Name, toCreate[j].(*api.Authorization).Name) != 1
			})

			So(len(toDelete), ShouldEqual, 2)
			So(toDelete[0].(*api.Authorization).Name, ShouldEqual, "2")
			So(toDelete[1].(*api.Authorization).Name, ShouldEqual, "3")
			So(len(toCreate), ShouldEqual, 2)
			So(toCreate[0].(*api.Authorization).Name, ShouldEqual, "2")
			So(toCreate[1].(*api.Authorization).Name, ShouldEqual, "4")
		})

		Convey("When I import a list of objects and there are some existing but they 404 on delete", func() {

			m.MockRetrieveMany(t, func(mctx manipulate.Context, dest elemental.Identifiables) error {
				*dest.(*api.AuthorizationsList) = append(
					*dest.(*api.AuthorizationsList),
					&api.Authorization{
						ID:          "1",
						Name:        "1",
						ImportHash:  "3132303033343839333331383835343436343834e3b0c44298fc1c149afbf4c8996fb92427ae41e4649b934ca495991b7852b855",
						ImportLabel: "label",
					},
					&api.Authorization{
						ID:          "2",
						Name:        "2",
						ImportHash:  "3132363235373937303539373039393132333639e3b0c44298fc1c149afbf4c8996fb92427ae41e4649b934ca495991b7852b855",
						ImportLabel: "label",
					},
					&api.Authorization{
						ID:          "3",
						Name:        "3",
						ImportHash:  "3",
						ImportLabel: "label",
					},
				)
				return nil
			})

			toDelete := elemental.IdentifiablesList{}
			m.MockDelete(t, func(mctx manipulate.Context, object elemental.Identifiable) error {
				return elemental.NewError("title", "desc", "sub", http.StatusNotFound)
			})

			toCreate := elemental.IdentifiablesList{}
			m.MockCreate(t, func(mctx manipulate.Context, object elemental.Identifiable) error {
				toCreate = append(toCreate, object)
				return nil
			})

			objs := api.AuthorizationsList{
				&api.Authorization{
					Name: "1",
				},
				&api.Authorization{
					Name:        "2",
					Description: "new",
				},
				&api.Authorization{
					Name: "4",
				},
			}

			err := Import(context.Background(), api.Manager(), m, "/ns", "label", objs, false)
			So(err, ShouldBeNil)

			sort.Slice(toDelete, func(i, j int) bool {
				return strings.Compare(toDelete[i].(*api.Authorization).Name, toDelete[j].(*api.Authorization).Name) != 1
			})
			sort.Slice(toCreate, func(i, j int) bool {
				return strings.Compare(toCreate[i].(*api.Authorization).Name, toCreate[j].(*api.Authorization).Name) != 1
			})

			So(len(toDelete), ShouldEqual, 0)
			So(len(toCreate), ShouldEqual, 2)
			So(toCreate[0].(*api.Authorization).Name, ShouldEqual, "2")
			So(toCreate[1].(*api.Authorization).Name, ShouldEqual, "4")
		})
	})
}
