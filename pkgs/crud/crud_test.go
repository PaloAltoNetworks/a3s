package crud

import (
	"context"
	"errors"
	"fmt"
	"testing"

	. "github.com/smartystreets/goconvey/convey"
	"go.aporeto.io/a3s/pkgs/api"
	"go.aporeto.io/bahamut"
	"go.aporeto.io/elemental"
	testmodel "go.aporeto.io/elemental/test/model"
	"go.aporeto.io/manipulate"
	"go.aporeto.io/manipulate/maniptest"
)

func TestCreate(t *testing.T) {

	Convey("Given I have a bahamut context and a manipulator", t, func() {

		bctx := bahamut.NewMockContext(context.Background())
		m := maniptest.NewTestManipulator()

		Convey("When everything is ok with a non namespaceable", func() {
			obj := testmodel.NewList()
			err := Create(bctx, m, obj)
			So(err, ShouldBeNil)
		})

		Convey("When everything is ok with a namespaceable", func() {

			obj := api.NewNamespace()
			bctx.MockRequest = &elemental.Request{Namespace: "/hello"}

			err := Create(bctx, m, obj)

			So(err, ShouldBeNil)
			So(obj.Namespace, ShouldEqual, "/hello")
		})

		Convey("When read only object is set", func() {

			obj := api.NewNamespace()
			obj.Namespace = "/hello"

			err := Create(bctx, m, obj)

			So(err, ShouldNotBeNil)
			So(errors.As(err, &elemental.Errors{}), ShouldBeTrue)
			So(err.Error(), ShouldResemble, "error 422 (elemental): Read Only Error: Field namespace is read only. You cannot set its value.")
		})

		Convey("When manipulate fails", func() {

			obj := api.NewNamespace()
			bctx.MockRequest = &elemental.Request{Namespace: "/hello"}

			m.MockCreate(t, func(mctx manipulate.Context, object elemental.Identifiable) error {
				return fmt.Errorf("boom")
			})

			err := Create(bctx, m, obj)

			So(err, ShouldNotBeNil)
			So(err.Error(), ShouldEqual, "boom")
		})

		Convey("When I use hooks that work", func() {

			obj := api.NewNamespace()
			bctx.MockRequest = &elemental.Request{Namespace: "/hello"}

			var preWriteHookCalled bool
			var preWriteObj elemental.Identifiable
			var preWriteOrig elemental.Identifiable
			preWrite := func(obj elemental.Identifiable, eobj elemental.Identifiable) error {
				preWriteHookCalled = true
				preWriteObj = obj
				preWriteOrig = eobj
				return nil
			}

			var postWriteHookCalled bool
			var postWriteObj elemental.Identifiable
			postWrite := func(obj elemental.Identifiable) {
				postWriteHookCalled = true
				postWriteObj = obj
			}
			err := Create(bctx, m, obj,
				OptionPreWriteHook(preWrite),
				OptionPostWriteHook(postWrite),
			)

			So(err, ShouldBeNil)
			So(preWriteHookCalled, ShouldBeTrue)
			So(preWriteObj, ShouldNotBeNil)
			So(preWriteOrig, ShouldBeNil)
			So(postWriteHookCalled, ShouldBeTrue)
			So(postWriteObj, ShouldNotBeNil)
		})

		Convey("When I use hooks that fails", func() {

			obj := api.NewNamespace()
			bctx.MockRequest = &elemental.Request{Namespace: "/hello"}

			preWrite := func(obj elemental.Identifiable, eobj elemental.Identifiable) error {
				return fmt.Errorf("boom")
			}

			err := Create(bctx, m, obj,
				OptionPreWriteHook(preWrite),
			)

			So(err, ShouldNotBeNil)
			So(err.Error(), ShouldEqual, "unable to run pre-write hook: boom")
			So(errors.As(err, &ErrPreWriteHook{}), ShouldBeTrue)
		})
	})
}

func TestRetrieveMany(t *testing.T) {

	Convey("Given I have a bahamut context and a manipulator", t, func() {

		bctx := bahamut.NewMockContext(context.Background())
		m := maniptest.NewTestManipulator()

		Convey("When everything is ok", func() {

			objs := api.NamespacesList{}
			bctx.MockRequest = &elemental.Request{Namespace: "/hello"}

			var expectedNamespace string
			m.MockRetrieveMany(t, func(mctx manipulate.Context, dest elemental.Identifiables) error {
				expectedNamespace = mctx.Namespace()
				*dest.(*api.NamespacesList) = append(
					*dest.(*api.NamespacesList),
					&api.Namespace{Name: "/hello/a"},
					&api.Namespace{Name: "/hello/b"},
				)
				return nil
			})

			err := RetrieveMany(bctx, m, &objs)

			So(err, ShouldBeNil)
			So(expectedNamespace, ShouldEqual, "/hello")
			So(len(objs), ShouldEqual, 2)
		})

		Convey("When translating context fails", func() {

			objs := api.NamespacesList{}
			bctx.MockRequest = &elemental.Request{
				Namespace: "/hello",
				Parameters: elemental.Parameters{
					"q": elemental.NewParameter(elemental.ParameterTypeString, "oops"),
				},
			}

			err := RetrieveMany(bctx, m, &objs)

			So(err, ShouldNotBeNil)
			So(err.Error(), ShouldStartWith, "error 400")
		})

		Convey("When manipulate fails", func() {

			objs := api.NamespacesList{}
			bctx.MockRequest = &elemental.Request{Namespace: "/hello"}

			m.MockRetrieveMany(t, func(mctx manipulate.Context, dest elemental.Identifiables) error {
				return fmt.Errorf("boom")
			})

			err := RetrieveMany(bctx, m, &objs)

			So(err, ShouldNotBeNil)
			So(err.Error(), ShouldEqual, "boom")
		})
	})
}

func TestUpdate(t *testing.T) {

	Convey("Given I have a bahamut context and a manipulator", t, func() {

		bctx := bahamut.NewMockContext(context.Background())
		m := maniptest.NewTestManipulator()

		Convey("When everything is ok", func() {

			obj := api.NewNamespace()
			bctx.MockRequest = &elemental.Request{
				Namespace: "/hello",
				ObjectID:  "xyz",
			}

			var expectedNamespace string
			var expectedID string
			m.MockUpdate(t, func(mctx manipulate.Context, object elemental.Identifiable) error {
				expectedNamespace = mctx.Namespace()
				expectedID = object.Identifier()
				return nil
			})

			err := Update(bctx, m, obj)

			So(err, ShouldBeNil)
			So(expectedNamespace, ShouldEqual, "/hello")
			So(expectedID, ShouldEqual, "xyz")
		})

		Convey("When read only object is set", func() {

			obj := api.NewNamespace()
			bctx.MockRequest = &elemental.Request{
				Namespace: "/hello",
				ObjectID:  "xyz",
			}

			obj.Namespace = "/new"
			err := Update(bctx, m, obj)

			So(err, ShouldNotBeNil)
			So(errors.As(err, &elemental.Errors{}), ShouldBeTrue)
			So(err.Error(), ShouldResemble, "error 422 (elemental): Read Only Error: Field namespace is read only. You cannot modify its value.")
		})

		Convey("When translating context fails", func() {

			obj := api.NewNamespace()
			bctx.MockRequest = &elemental.Request{
				Namespace: "/hello",
				ObjectID:  "xyz",
				Parameters: elemental.Parameters{
					"q": elemental.NewParameter(elemental.ParameterTypeString, "oops"),
				},
			}

			err := Update(bctx, m, obj)

			So(err, ShouldNotBeNil)
			So(err.Error(), ShouldStartWith, "error 400")
		})

		Convey("When manipulate fails to retrieve", func() {

			obj := api.NewNamespace()
			bctx.MockRequest = &elemental.Request{Namespace: "/hello"}

			m.MockRetrieve(t, func(mctx manipulate.Context, object elemental.Identifiable) error {
				return fmt.Errorf("boom")
			})

			err := Update(bctx, m, obj)

			So(err, ShouldNotBeNil)
			So(err.Error(), ShouldEqual, "boom")
		})

		Convey("When manipulate fails to delete", func() {

			obj := api.NewNamespace()
			bctx.MockRequest = &elemental.Request{Namespace: "/hello"}

			m.MockUpdate(t, func(mctx manipulate.Context, object elemental.Identifiable) error {
				return fmt.Errorf("boom")
			})

			err := Update(bctx, m, obj)

			So(err, ShouldNotBeNil)
			So(err.Error(), ShouldEqual, "boom")
		})

		Convey("When I use hooks that work", func() {

			obj := api.NewNamespace()
			bctx.MockRequest = &elemental.Request{Namespace: "/hello"}

			var preWriteHookCalled bool
			var preWriteObj elemental.Identifiable
			var preWriteOrig elemental.Identifiable
			preWrite := func(obj elemental.Identifiable, eobj elemental.Identifiable) error {
				preWriteHookCalled = true
				preWriteObj = obj
				preWriteOrig = eobj
				return nil
			}

			var postWriteHookCalled bool
			var postWriteObj elemental.Identifiable
			postWrite := func(obj elemental.Identifiable) {
				postWriteHookCalled = true
				postWriteObj = obj
			}
			err := Update(bctx, m, obj,
				OptionPreWriteHook(preWrite),
				OptionPostWriteHook(postWrite),
			)

			So(err, ShouldBeNil)
			So(preWriteHookCalled, ShouldBeTrue)
			So(preWriteObj, ShouldNotBeNil)
			So(preWriteOrig, ShouldNotBeNil)
			So(postWriteHookCalled, ShouldBeTrue)
			So(postWriteObj, ShouldNotBeNil)
		})

		Convey("When I use hooks that fails", func() {

			obj := api.NewNamespace()
			bctx.MockRequest = &elemental.Request{Namespace: "/hello"}

			preWrite := func(obj elemental.Identifiable, eobj elemental.Identifiable) error {
				return fmt.Errorf("boom")
			}

			err := Update(bctx, m, obj,
				OptionPreWriteHook(preWrite),
			)

			So(err, ShouldNotBeNil)
			So(err.Error(), ShouldEqual, "unable to run pre-write hook: boom")
			So(errors.As(err, &ErrPreWriteHook{}), ShouldBeTrue)
		})
	})
}

func TestRetrieve(t *testing.T) {

	Convey("Given I have a bahamut context and a manipulator", t, func() {

		bctx := bahamut.NewMockContext(context.Background())
		m := maniptest.NewTestManipulator()

		Convey("When everything is ok", func() {

			obj := api.NewNamespace()
			bctx.MockRequest = &elemental.Request{
				Namespace: "/hello",
				ObjectID:  "xyz",
			}

			var expectedNamespace string
			var expectedID string
			m.MockRetrieve(t, func(mctx manipulate.Context, object elemental.Identifiable) error {
				expectedNamespace = mctx.Namespace()
				expectedID = object.Identifier()
				return nil
			})

			err := Retrieve(bctx, m, obj)

			So(err, ShouldBeNil)
			So(expectedNamespace, ShouldEqual, "/hello")
			So(expectedID, ShouldEqual, "xyz")
		})

		Convey("When translating context fails", func() {

			obj := api.NewNamespace()
			bctx.MockRequest = &elemental.Request{
				Namespace: "/hello",
				ObjectID:  "xyz",
				Parameters: elemental.Parameters{
					"q": elemental.NewParameter(elemental.ParameterTypeString, "oops"),
				},
			}

			err := Retrieve(bctx, m, obj)

			So(err, ShouldNotBeNil)
			So(err.Error(), ShouldStartWith, "error 400")
		})

		Convey("When manipulate fails", func() {

			obj := api.NewNamespace()
			bctx.MockRequest = &elemental.Request{Namespace: "/hello"}

			m.MockRetrieve(t, func(mctx manipulate.Context, object elemental.Identifiable) error {
				return fmt.Errorf("boom")
			})

			err := Retrieve(bctx, m, obj)

			So(err, ShouldNotBeNil)
			So(err.Error(), ShouldEqual, "boom")
		})
	})
}

func TestDelete(t *testing.T) {

	Convey("Given I have a bahamut context and a manipulator", t, func() {

		bctx := bahamut.NewMockContext(context.Background())
		m := maniptest.NewTestManipulator()

		Convey("When everything is ok", func() {

			obj := api.NewNamespace()
			bctx.MockRequest = &elemental.Request{
				Namespace: "/hello",
				ObjectID:  "xyz",
			}

			var expectedNamespace string
			var expectedID string
			m.MockDelete(t, func(mctx manipulate.Context, object elemental.Identifiable) error {
				expectedNamespace = mctx.Namespace()
				expectedID = object.Identifier()
				return nil
			})

			err := Delete(bctx, m, obj)

			So(err, ShouldBeNil)
			So(expectedNamespace, ShouldEqual, "/hello")
			So(expectedID, ShouldEqual, "xyz")
		})

		Convey("When translating context fails", func() {

			obj := api.NewNamespace()
			bctx.MockRequest = &elemental.Request{
				Namespace: "/hello",
				ObjectID:  "xyz",
				Parameters: elemental.Parameters{
					"q": elemental.NewParameter(elemental.ParameterTypeString, "oops"),
				},
			}

			err := Delete(bctx, m, obj)

			So(err, ShouldNotBeNil)
			So(err.Error(), ShouldStartWith, "error 400")
		})

		Convey("When manipulate fails to retrieve", func() {

			obj := api.NewNamespace()
			bctx.MockRequest = &elemental.Request{Namespace: "/hello"}

			m.MockRetrieve(t, func(mctx manipulate.Context, object elemental.Identifiable) error {
				return fmt.Errorf("boom")
			})

			err := Delete(bctx, m, obj)

			So(err, ShouldNotBeNil)
			So(err.Error(), ShouldEqual, "boom")
		})

		Convey("When manipulate fails to delete", func() {

			obj := api.NewNamespace()
			bctx.MockRequest = &elemental.Request{Namespace: "/hello"}

			m.MockDelete(t, func(mctx manipulate.Context, object elemental.Identifiable) error {
				return fmt.Errorf("boom")
			})

			err := Delete(bctx, m, obj)

			So(err, ShouldNotBeNil)
			So(err.Error(), ShouldEqual, "boom")
		})

		Convey("When I use hooks that work", func() {

			obj := api.NewNamespace()
			bctx.MockRequest = &elemental.Request{Namespace: "/hello"}

			var preWriteHookCalled bool
			var preWriteObj elemental.Identifiable
			var preWriteOrig elemental.Identifiable
			preWrite := func(obj elemental.Identifiable, eobj elemental.Identifiable) error {
				preWriteHookCalled = true
				preWriteObj = obj
				preWriteOrig = eobj
				return nil
			}

			var postWriteHookCalled bool
			var postWriteObj elemental.Identifiable
			postWrite := func(obj elemental.Identifiable) {
				postWriteHookCalled = true
				postWriteObj = obj
			}
			err := Delete(bctx, m, obj,
				OptionPreWriteHook(preWrite),
				OptionPostWriteHook(postWrite),
			)

			So(err, ShouldBeNil)
			So(preWriteHookCalled, ShouldBeTrue)
			So(preWriteObj, ShouldNotBeNil)
			So(preWriteOrig, ShouldBeNil)
			So(postWriteHookCalled, ShouldBeTrue)
			So(postWriteObj, ShouldNotBeNil)
		})

		Convey("When I use hooks that fails", func() {

			obj := api.NewNamespace()
			bctx.MockRequest = &elemental.Request{Namespace: "/hello"}

			preWrite := func(obj elemental.Identifiable, eobj elemental.Identifiable) error {
				return fmt.Errorf("boom")
			}

			err := Delete(bctx, m, obj,
				OptionPreWriteHook(preWrite),
			)

			So(err, ShouldNotBeNil)
			So(err.Error(), ShouldEqual, "unable to run pre-write hook: boom")
			So(errors.As(err, &ErrPreWriteHook{}), ShouldBeTrue)
		})
	})
}

func TestInfo(t *testing.T) {

	Convey("Given I have a bahamut context and a manipulator", t, func() {

		bctx := bahamut.NewMockContext(context.Background())
		m := maniptest.NewTestManipulator()

		Convey("When everything is ok", func() {

			bctx.MockRequest = &elemental.Request{
				Namespace: "/hello",
			}

			var expectedNamespace string
			var expectedIdentity elemental.Identity
			m.MockCount(t, func(mctx manipulate.Context, identity elemental.Identity) (int, error) {
				expectedNamespace = mctx.Namespace()
				expectedIdentity = identity
				return 42, nil
			})

			err := Info(bctx, m, api.NamespaceIdentity)

			So(err, ShouldBeNil)
			So(expectedNamespace, ShouldEqual, "/hello")
			So(expectedIdentity.Name, ShouldEqual, "namespace")
		})

		Convey("When translating context fails", func() {

			bctx.MockRequest = &elemental.Request{
				Namespace: "/hello",
				Parameters: elemental.Parameters{
					"q": elemental.NewParameter(elemental.ParameterTypeString, "oops"),
				},
			}

			err := Info(bctx, m, api.NamespaceIdentity)

			So(err, ShouldNotBeNil)
			So(err.Error(), ShouldStartWith, "error 400")
		})

		Convey("When manipulate fails", func() {

			bctx.MockRequest = &elemental.Request{Namespace: "/hello"}

			m.MockCount(t, func(mctx manipulate.Context, identity elemental.Identity) (int, error) {
				return 0, fmt.Errorf("boom")
			})

			err := Info(bctx, m, api.NamespaceIdentity)

			So(err, ShouldNotBeNil)
			So(err.Error(), ShouldEqual, "boom")
		})
	})
}
