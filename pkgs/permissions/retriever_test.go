package permissions

import (
	"context"
	"fmt"
	"testing"

	. "github.com/smartystreets/goconvey/convey"
	"go.aporeto.io/a3s/pkgs/api"
	"go.aporeto.io/elemental"
	"go.aporeto.io/manipulate"
	"go.aporeto.io/manipulate/maniptest"
)

func flattenTags(term [][]string) (set []string) {
	for _, rows := range term {
		set = append(set, rows...)
	}
	return set
}

func TestNewRetriever(t *testing.T) {
	Convey("Given have a subscriber and a manipulator", t, func() {
		m := maniptest.NewTestManipulator()
		a := NewRetriever(m).(*retriever)
		So(a.manipulator, ShouldNotBeNil)
	})
}

func TestIsAuthorizedWithToken(t *testing.T) {

	var (
		permSetAllowAll = "*,*"
		permSetOnBla    = "bla"
		ctx             = context.Background()
	)

	makeAPIPol := func(perms []string, subnets []string) *api.Authorization {
		apiauth := api.NewAuthorization()
		apiauth.ID = "1"
		apiauth.Namespace = "/a"
		apiauth.Subject = [][]string{{"color=blue"}}
		apiauth.TargetNamespace = "/a"
		apiauth.Permissions = perms
		apiauth.Subnets = subnets
		apiauth.FlattenedSubject = flattenTags(apiauth.Subject)

		return apiauth
	}

	Convey("Given I have an authorizer and a token", t, func() {

		m := maniptest.NewTestManipulator()

		r := NewRetriever(m).(*retriever)

		m.MockCount(t, func(mctx manipulate.Context, identity elemental.Identity) (int, error) {
			return 1, nil
		})

		Convey("When there is no policy matching", func() {

			perms, err := r.Permissions(ctx, []string{"color=blue"}, "/a")

			So(err, ShouldBeNil)
			So(IsAllowed(perms, "get", "things"), ShouldEqual, false)
		})

		Convey("When I retrieving the ns fails", func() {

			m.MockCount(t, func(mctx manipulate.Context, identity elemental.Identity) (int, error) {
				return 0, fmt.Errorf("noooo")
			})

			perms, err := r.Permissions(ctx, []string{"color=blue"}, "/a")

			So(err, ShouldNotBeNil)
			So(err.Error(), ShouldEqual, "noooo")
			So(IsAllowed(perms, "get", "things"), ShouldEqual, false)
		})

		Convey("When I retrieving the ns is not found", func() {

			m.MockCount(t, func(mctx manipulate.Context, identity elemental.Identity) (int, error) {
				return 0, nil
			})

			perms, err := r.Permissions(ctx, []string{"color=blue"}, "/a")

			So(err, ShouldBeNil)
			So(IsAllowed(perms, "get", "things"), ShouldEqual, false)
		})

		Convey("When there is a policy matching *,*", func() {

			m.MockRetrieveMany(t, func(mctx manipulate.Context, dest elemental.Identifiables) error {
				*dest.(*api.AuthorizationsList) = append(
					*dest.(*api.AuthorizationsList),
					makeAPIPol([]string{permSetAllowAll}, nil),
				)
				return nil
			})

			perms, err := r.Permissions(ctx, []string{"color=blue"}, "/a")

			So(err, ShouldBeNil)
			So(IsAllowed(perms, "get", "things"), ShouldEqual, true)
		})

		Convey("When there is a policy matching twice using twice the same set", func() {

			m.MockRetrieveMany(t, func(mctx manipulate.Context, dest elemental.Identifiables) error {
				*dest.(*api.AuthorizationsList) = append(
					*dest.(*api.AuthorizationsList),
					makeAPIPol([]string{"things,delete"}, nil),
					makeAPIPol([]string{"things,get"}, nil),
				)
				return nil
			})

			perms, err := r.Permissions(ctx, []string{"color=blue"}, "/a")

			So(err, ShouldBeNil)
			So(IsAllowed(perms, "get", "things"), ShouldEqual, true)
		})

		Convey("When there is a policy matching with target namespace outside of restricted ns", func() {

			m.MockRetrieveMany(t, func(mctx manipulate.Context, dest elemental.Identifiables) error {
				*dest.(*api.AuthorizationsList) = append(
					*dest.(*api.AuthorizationsList),
					makeAPIPol([]string{permSetAllowAll}, nil),
				)
				return nil
			})

			perms, err := r.Permissions(ctx, []string{"color=blue"}, "/a",
				OptionRetrieverRestrictions(Restrictions{Namespace: "/b"}),
			)

			So(err, ShouldBeNil)
			So(IsAllowed(perms, "get", "things"), ShouldEqual, false)
		})

		Convey("When there is a policy matching with target namespace equals to restricted ns", func() {

			m.MockRetrieveMany(t, func(mctx manipulate.Context, dest elemental.Identifiables) error {
				*dest.(*api.AuthorizationsList) = append(
					*dest.(*api.AuthorizationsList),
					makeAPIPol([]string{permSetAllowAll}, nil),
				)
				return nil
			})

			perms, err := r.Permissions(ctx, []string{"color=blue"}, "/a",
				OptionRetrieverRestrictions(Restrictions{Namespace: "/a"}),
			)

			So(err, ShouldBeNil)
			So(IsAllowed(perms, "get", "things"), ShouldEqual, true)
		})

		Convey("When there is a policy matching with target namespace a child of restricted ns", func() {

			m.MockRetrieveMany(t, func(mctx manipulate.Context, dest elemental.Identifiables) error {
				*dest.(*api.AuthorizationsList) = append(
					*dest.(*api.AuthorizationsList),
					makeAPIPol([]string{permSetAllowAll}, nil),
				)
				return nil
			})

			perms, err := r.Permissions(ctx, []string{"color=blue"}, "/a/b",
				OptionRetrieverRestrictions(Restrictions{Namespace: "/a"}),
			)

			So(err, ShouldBeNil)
			So(IsAllowed(perms, "get", "things"), ShouldEqual, true)
		})

		Convey("When there is a policy matching with target namespace a parent of restricted ns", func() {

			m.MockRetrieveMany(t, func(mctx manipulate.Context, dest elemental.Identifiables) error {
				*dest.(*api.AuthorizationsList) = append(
					*dest.(*api.AuthorizationsList),
					makeAPIPol([]string{permSetAllowAll}, nil),
				)
				return nil
			})

			perms, err := r.Permissions(ctx, []string{"color=blue"}, "/a",
				OptionRetrieverRestrictions(Restrictions{Namespace: "/a/b"}),
			)

			So(err, ShouldBeNil)
			So(IsAllowed(perms, "get", "things"), ShouldEqual, false)
		})

		Convey("When there is a policy matching with a bad targetNs", func() {

			pol := makeAPIPol([]string{permSetAllowAll}, nil)
			pol.TargetNamespace = "/az/b/c"

			m.MockRetrieveMany(t, func(mctx manipulate.Context, dest elemental.Identifiables) error {
				*dest.(*api.AuthorizationsList) = append(
					*dest.(*api.AuthorizationsList),
					pol,
				)
				return nil
			})

			perms, err := r.Permissions(ctx, []string{"color=blue"}, "/a")

			So(err, ShouldBeNil)
			So(IsAllowed(perms, "get", "things"), ShouldEqual, false)
		})

		Convey("When there is a policy that is not matching", func() {

			m.MockRetrieveMany(t, func(mctx manipulate.Context, dest elemental.Identifiables) error {
				*dest.(*api.AuthorizationsList) = append(
					*dest.(*api.AuthorizationsList),
					makeAPIPol([]string{"nope,*"}, nil),
				)
				return nil
			})

			perms, err := r.Permissions(ctx, []string{"color=blue"}, "/a")

			So(err, ShouldBeNil)
			So(IsAllowed(perms, "get", "things"), ShouldEqual, false)
		})

		Convey("When there is a policy matching but not on the correct permission set", func() {

			m.MockRetrieveMany(t, func(mctx manipulate.Context, dest elemental.Identifiables) error {
				*dest.(*api.AuthorizationsList) = append(
					*dest.(*api.AuthorizationsList),
					makeAPIPol([]string{permSetOnBla}, nil),
				)
				return nil
			})

			perms, err := r.Permissions(ctx, []string{"color=blue"}, "/a")

			So(err, ShouldBeNil)
			So(IsAllowed(perms, "get", "things"), ShouldEqual, false)
		})

		Convey("When there is a policy matching", func() {

			m.MockRetrieveMany(t, func(mctx manipulate.Context, dest elemental.Identifiables) error {
				*dest.(*api.AuthorizationsList) = append(
					*dest.(*api.AuthorizationsList),
					makeAPIPol([]string{"things,get"}, nil),
				)
				return nil
			})

			perms, err := r.Permissions(ctx, []string{"color=blue"}, "/a")

			So(err, ShouldBeNil)
			So(IsAllowed(perms, "get", "things"), ShouldEqual, true)
		})

		Convey("When there is a policy with matching restriction", func() {

			m.MockRetrieveMany(t, func(mctx manipulate.Context, dest elemental.Identifiables) error {
				*dest.(*api.AuthorizationsList) = append(
					*dest.(*api.AuthorizationsList),
					makeAPIPol([]string{"things,get"}, []string{"10.0.0.0/8", "11.0.0.0/8"}),
				)
				return nil
			})

			perms, err := r.Permissions(ctx, []string{"color=blue"}, "/a",
				OptionRetrieverIPAddr("11.2.2.2"),
			)

			So(err, ShouldBeNil)
			So(IsAllowed(perms, "get", "things"), ShouldEqual, true)
		})

		Convey("When there is a policy with not matching restriction", func() {

			m.MockRetrieveMany(t, func(mctx manipulate.Context, dest elemental.Identifiables) error {
				*dest.(*api.AuthorizationsList) = append(
					*dest.(*api.AuthorizationsList),
					makeAPIPol([]string{"things,get"}, []string{"10.0.0.0/8", "11.0.0.0/8"}),
				)
				return nil
			})

			perms, err := r.Permissions(ctx, []string{"color=blue"}, "/a",
				OptionRetrieverIPAddr("13.2.2.2"),
			)

			So(err, ShouldBeNil)
			So(IsAllowed(perms, "get", "things"), ShouldEqual, false)
		})

		Convey("When there is a policy invalid IP", func() {

			m.MockRetrieveMany(t, func(mctx manipulate.Context, dest elemental.Identifiables) error {
				*dest.(*api.AuthorizationsList) = append(
					*dest.(*api.AuthorizationsList),
					makeAPIPol([]string{"things,get"}, []string{"10.0.0.0/8", "11.0.0.0/8"}),
				)
				return nil
			})

			perms, err := r.Permissions(ctx, []string{"color=blue"}, "/a",
				OptionRetrieverIPAddr(".2.2.2"),
			)

			So(err, ShouldNotBeNil)
			So(err.Error(), ShouldEqual, "missing or invalid origin IP '.2.2.2'")
			So(IsAllowed(perms, "get", "things"), ShouldEqual, false)
		})

		Convey("When there is a policy invalid declared CIDR", func() {

			m.MockRetrieveMany(t, func(mctx manipulate.Context, dest elemental.Identifiables) error {
				*dest.(*api.AuthorizationsList) = append(
					*dest.(*api.AuthorizationsList),
					makeAPIPol([]string{"things,get"}, []string{"dawf"}),
				)
				return nil
			})

			perms, err := r.Permissions(ctx, []string{"color=blue"}, "/a",
				OptionRetrieverIPAddr("2.2.2.2"),
			)

			So(err, ShouldNotBeNil)
			So(err.Error(), ShouldEqual, "invalid CIDR address: dawf")
			So(IsAllowed(perms, "get", "things"), ShouldEqual, false)
		})

		Convey("When there is a policy with empty subject", func() {

			pol := makeAPIPol([]string{}, nil)
			pol.Subject = [][]string{{}}

			m.MockRetrieveMany(t, func(mctx manipulate.Context, dest elemental.Identifiables) error {
				*dest.(*api.AuthorizationsList) = append(
					*dest.(*api.AuthorizationsList),
					pol,
				)
				return nil
			})

			perms, err := r.Permissions(ctx, []string{"color=blue"}, "/a")

			So(err, ShouldBeNil)
			So(IsAllowed(perms, "get", "things"), ShouldEqual, false)
		})

		Convey("When there is a policy with empty string subject", func() {

			pol := makeAPIPol([]string{}, nil)
			pol.Subject = [][]string{{""}}

			m.MockRetrieveMany(t, func(mctx manipulate.Context, dest elemental.Identifiables) error {
				*dest.(*api.AuthorizationsList) = append(
					*dest.(*api.AuthorizationsList),
					pol,
				)
				return nil
			})

			perms, err := r.Permissions(ctx, []string{"color=blue"}, "/a")

			So(err, ShouldBeNil)
			So(IsAllowed(perms, "get", "things"), ShouldEqual, false)
		})

		Convey("When there is a policy with nil subject", func() {

			pol := makeAPIPol([]string{}, nil)
			pol.Subject = [][]string{nil}

			m.MockRetrieveMany(t, func(mctx manipulate.Context, dest elemental.Identifiables) error {
				*dest.(*api.AuthorizationsList) = append(
					*dest.(*api.AuthorizationsList),
					pol,
				)
				return nil
			})

			perms, err := r.Permissions(ctx, []string{"color=blue"}, "/a")

			So(err, ShouldBeNil)
			So(IsAllowed(perms, "get", "things"), ShouldEqual, false)
		})

		Convey("When retrieving the policy fails", func() {

			m.MockRetrieveMany(t, func(mctx manipulate.Context, dest elemental.Identifiables) error {
				return fmt.Errorf("boom")
			})

			perms, err := r.Permissions(ctx, []string{"color=blue"}, "/a")

			So(err, ShouldNotBeNil)
			So(err.Error(), ShouldEqual, "unable to retrieve api authorizations: boom")
			So(IsAllowed(perms, "get", "things"), ShouldEqual, false)
		})

		Convey("When retrieving the policy with an invalid allowedSubnet", func() {

			m.MockRetrieveMany(t, func(mctx manipulate.Context, dest elemental.Identifiables) error {
				*dest.(*api.AuthorizationsList) = append(
					*dest.(*api.AuthorizationsList),
					makeAPIPol([]string{}, []string{".2.2.2."}),
				)
				return nil
			})

			perms, err := r.Permissions(ctx, []string{"color=blue"}, "/a")

			So(err, ShouldNotBeNil)
			So(err.Error(), ShouldEqual, "missing or invalid origin IP ''")
			So(IsAllowed(perms, "get", "things"), ShouldEqual, false)
		})

		// Restrictions

		Convey("When there is a policy matching but the namespace is restricted", func() {

			m.MockRetrieveMany(t, func(mctx manipulate.Context, dest elemental.Identifiables) error {
				*dest.(*api.AuthorizationsList) = append(
					*dest.(*api.AuthorizationsList),
					makeAPIPol([]string{"@auth:role=testrole2"}, nil),
				)
				return nil
			})

			perms, err := r.Permissions(ctx, []string{"color=blue"}, "/a",
				OptionRetrieverRestrictions(Restrictions{Namespace: "/b"}),
			)

			So(err, ShouldBeNil)
			So(IsAllowed(perms, "get", "things"), ShouldEqual, false)
		})

		Convey("When there is a policy matching but the permissions are restricted", func() {

			m.MockRetrieveMany(t, func(mctx manipulate.Context, dest elemental.Identifiables) error {
				*dest.(*api.AuthorizationsList) = append(
					*dest.(*api.AuthorizationsList),
					makeAPIPol([]string{permSetAllowAll}, nil),
				)
				return nil
			})

			perms, err := r.Permissions(ctx, []string{"color=blue"}, "/a",
				OptionRetrieverRestrictions(Restrictions{Permissions: []string{"dog,get"}}),
			)

			So(err, ShouldBeNil)
			So(IsAllowed(perms, "get", "things"), ShouldEqual, false)
		})

		Convey("When there is a policy matching but the networks are restricted", func() {

			m.MockRetrieveMany(t, func(mctx manipulate.Context, dest elemental.Identifiables) error {
				*dest.(*api.AuthorizationsList) = append(
					*dest.(*api.AuthorizationsList),
					makeAPIPol([]string{permSetAllowAll}, nil),
				)
				return nil
			})

			Convey("When I the networks are correct", func() {

				perms, err := r.Permissions(ctx, []string{"color=blue"}, "/a",
					OptionRetrieverIPAddr("127.0.0.1"),
					OptionRetrieverRestrictions(Restrictions{Networks: []string{"10.0.0.0/8"}}),
				)

				So(err, ShouldBeNil)
				So(IsAllowed(perms, "get", "things"), ShouldEqual, false)
			})

			Convey("When I the networks are incorrect", func() {

				perms, err := r.Permissions(ctx, []string{"color=blue"}, "/a",
					OptionRetrieverIPAddr("1.1.1.1"),
					OptionRetrieverRestrictions(Restrictions{Networks: []string{"how-come?"}}),
				)

				So(err, ShouldNotBeNil)
				So(err.Error(), ShouldEqual, `invalid CIDR address: how-come?`)
				So(IsAllowed(perms, "get", "things"), ShouldEqual, false)
			})
		})

		// Single ID target

		Convey("When there is a policy with id restriction and not id provided", func() {

			m.MockRetrieveMany(t, func(mctx manipulate.Context, dest elemental.Identifiables) error {
				*dest.(*api.AuthorizationsList) = append(
					*dest.(*api.AuthorizationsList),
					makeAPIPol([]string{"things,get:xyz,abc"}, nil),
				)
				return nil
			})

			perms, err := r.Permissions(ctx, []string{"color=blue"}, "/a")

			So(err, ShouldBeNil)
			So(IsAllowed(perms, "get", "things"), ShouldEqual, false)
		})

		Convey("When there is a policy with id restriction and not matching id provided", func() {

			m.MockRetrieveMany(t, func(mctx manipulate.Context, dest elemental.Identifiables) error {
				*dest.(*api.AuthorizationsList) = append(
					*dest.(*api.AuthorizationsList),
					makeAPIPol([]string{"things,get:xyz,abc"}, nil),
				)
				return nil
			})

			perms, err := r.Permissions(ctx, []string{"color=blue"}, "/a",
				OptionRetrieverID("nope-id"),
			)

			So(err, ShouldBeNil)
			So(IsAllowed(perms, "get", "things"), ShouldEqual, false)
		})

		Convey("When there is a policy with id restriction and matching id provided", func() {

			m.MockRetrieveMany(t, func(mctx manipulate.Context, dest elemental.Identifiables) error {
				*dest.(*api.AuthorizationsList) = append(
					*dest.(*api.AuthorizationsList),
					makeAPIPol([]string{"things,get:xyz,abc"}, nil),
				)
				return nil
			})

			perms, err := r.Permissions(ctx, []string{"color=blue"}, "/a",
				OptionRetrieverID("xyz"),
			)

			So(err, ShouldBeNil)
			So(IsAllowed(perms, "get", "things"), ShouldEqual, true)
		})
	})
}

func TestPermissionsWithToken(t *testing.T) {

	var (
		testrole1 = "stuff,*"
		testrole2 = "*,*"
		testrole3 = "bla,get,post,put"
		ctx       = context.Background()
	)

	makeAPIPol := func(perms []string, subnets []string) *api.Authorization {
		apiauth := api.NewAuthorization()
		apiauth.ID = "1"
		apiauth.Namespace = "/a"
		apiauth.Subject = [][]string{{"color=blue"}}
		apiauth.TargetNamespace = "/"
		apiauth.Permissions = perms
		apiauth.Subnets = subnets
		apiauth.FlattenedSubject = flattenTags(apiauth.Subject)

		return apiauth
	}

	Convey("Given I have an authorizer and a token", t, func() {

		m := maniptest.NewTestManipulator()

		r := NewRetriever(m).(*retriever)

		m.MockCount(t, func(mctx manipulate.Context, identity elemental.Identity) (int, error) {
			return 1, nil
		})

		Convey("When I call Authorizations when I have no policy", func() {

			perms, err := r.Permissions(ctx, []string{"color=blue"}, "/a")

			So(err, ShouldBeNil)
			So(perms, ShouldResemble, PermissionMap{})
		})

		Convey("When I call Authorizations when I have the role testroles2", func() {

			m.MockRetrieveMany(t, func(mctx manipulate.Context, dest elemental.Identifiables) error {
				*dest.(*api.AuthorizationsList) = append(
					*dest.(*api.AuthorizationsList),
					makeAPIPol([]string{testrole2}, nil),
				)
				return nil
			})

			perms, err := r.Permissions(ctx, []string{"color=blue"}, "/a")

			So(err, ShouldBeNil)
			So(perms, ShouldResemble, PermissionMap{
				"*": {"*": true},
			})
		})

		Convey("When I call Authorizations when I have the role testroles1 and testrole3", func() {

			m.MockRetrieveMany(t, func(mctx manipulate.Context, dest elemental.Identifiables) error {
				*dest.(*api.AuthorizationsList) = append(
					*dest.(*api.AuthorizationsList),
					makeAPIPol([]string{testrole1, testrole3}, nil),
				)
				return nil
			})

			perms, err := r.Permissions(ctx, []string{"color=blue"}, "/a")

			So(err, ShouldBeNil)
			So(perms, ShouldResemble, PermissionMap{
				"bla":   {"get": true, "post": true, "put": true},
				"stuff": {"*": true},
			})
		})

		Convey("When I call Authorizations when I have with individual authotization", func() {

			m.MockRetrieveMany(t, func(mctx manipulate.Context, dest elemental.Identifiables) error {
				*dest.(*api.AuthorizationsList) = append(
					*dest.(*api.AuthorizationsList),
					makeAPIPol([]string{"r1,get,post", "r2,put"}, nil),
				)
				return nil
			})

			perms, err := r.Permissions(ctx, []string{"color=blue"}, "/a")

			So(err, ShouldBeNil)
			So(perms, ShouldResemble, PermissionMap{
				"r1": {"get": true, "post": true},
				"r2": {"put": true},
			})
		})
	})
}

func TestCountNamespace(t *testing.T) {

	Convey("Given I have a http manipulator and an authorizer", t, func() {

		m := maniptest.NewTestManipulator()

		r := NewRetriever(m).(*retriever)

		Convey("When I call countNamespace", func() {

			attempt := -1
			consistency := manipulate.ReadConsistencyDefault
			m.MockCount(t, func(mctx manipulate.Context, identity elemental.Identity) (int, error) {
				consistency = mctx.ReadConsistency()
				attempt++
				return attempt, nil
			})

			count, err := r.countNamespace(context.Background(), "ns")

			So(err, ShouldBeNil)
			So(count, ShouldEqual, 1)
			So(consistency, ShouldEqual, manipulate.ReadConsistencyStrong)
		})
	})
}
