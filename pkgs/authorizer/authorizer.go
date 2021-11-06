package authorizer

import (
	"context"
	"fmt"
	"math/rand"
	"net/http"
	"time"

	"github.com/spaolacci/murmur3"
	"go.aporeto.io/a3s/pkgs/nscache"
	"go.aporeto.io/a3s/pkgs/permissions"
	"go.aporeto.io/a3s/pkgs/token"
	"go.aporeto.io/bahamut"
	"go.aporeto.io/elemental"
)

// Various Authorizer errors.
var (
	ErrMissingNamespace = elemental.NewError(
		"Forbidden",
		"Missing X-Namespace header",
		"a3s:authorizer",
		http.StatusForbidden,
	)

	ErrInvalidNamespace = elemental.NewError(
		"Forbidden",
		"Invalid X-Namespace header. A namespace must start with /",
		"a3s:authorizer",
		http.StatusForbidden,
	)

	ErrMissingToken = elemental.NewError(
		"Forbidden",
		"Missing token in either Authorization header or X-A3S-Token in cookies",
		"a3s:authorizer",
		http.StatusForbidden,
	)
)

// An Authorizer is a bahamut.Authorizer that provides
// additional methods.
type Authorizer interface {
	bahamut.Authorizer

	CheckAuthorization(
		ctx context.Context,
		claims []string,
		op string,
		ns string,
		resource string,
		opts ...OptionCheck,
	) (bool, error)
}

// An Authorizer is the enforcer of the authorizations of all API calls.
//
// It implements the bahamut.Authorizer interface.
type authorizer struct {
	retriever        permissions.Retriever
	ignoredResources map[string]struct{}
	cache            *nscache.NamespacedCache
}

// New creates a new Authorizer using cid.
func New(ctx context.Context, retriever permissions.Retriever, pubsub bahamut.PubSubClient, options ...Option) Authorizer {

	cfg := config{}
	for _, opt := range options {
		opt(&cfg)
	}

	ignored := map[string]struct{}{}
	for _, i := range cfg.ignoredResources {
		ignored[i] = struct{}{}
	}

	authCache := nscache.New(pubsub, 24000)
	if pubsub != nil {
		authCache.Start(ctx)
	}

	return &authorizer{
		retriever:        retriever,
		ignoredResources: ignored,
		cache:            authCache,
	}
}

// IsAuthorized is the main method that returns whether the API call is authorized or not.
func (a *authorizer) IsAuthorized(ctx bahamut.Context) (bahamut.AuthAction, error) {

	req := ctx.Request()

	if _, ok := a.ignoredResources[req.Identity.Category]; ok {
		return bahamut.AuthActionOK, nil
	}

	token := token.FromRequest(req)
	if token == "" {
		return bahamut.AuthActionKO, ErrMissingToken
	}

	restrictions, err := permissions.GetRestrictions(token)
	if err != nil {
		return bahamut.AuthActionKO, elemental.NewError(
			"Forbidden",
			err.Error(),
			"a3s:authorizer",
			http.StatusForbidden,
		)
	}

	ok, err := a.CheckAuthorization(
		ctx.Context(),
		ctx.Claims(),
		string(req.Operation),
		req.Namespace,
		req.Identity.Category,
		OptionCheckRestrictions(restrictions),
		OptionCheckID(req.ObjectID),
		OptionCheckSourceIP(req.ClientIP),
	)
	if err != nil {
		return bahamut.AuthActionKO, err
	}

	if ok {
		return bahamut.AuthActionOK, nil
	}

	return bahamut.AuthActionKO, nil
}

func (a *authorizer) CheckAuthorization(ctx context.Context, claims []string, operation string, ns string, resource string, opts ...OptionCheck) (bool, error) {

	cfg := checkConfig{}
	for _, o := range opts {
		o(&cfg)
	}

	if _, ok := a.ignoredResources[resource]; ok {
		return true, nil
	}

	if ns == "" {
		return false, ErrMissingNamespace
	}

	if ns[0] != '/' {
		return false, ErrInvalidNamespace
	}

	key := hash(claims, cfg.sourceIP, cfg.id, cfg.restrictions)

	if r := a.cache.Get(ns, key); r != nil && !r.Expired() {
		perms := r.Value().(permissions.PermissionMap)
		return perms.Allows(operation, resource), nil
	}

	ropts := []permissions.RetrieverOption{
		permissions.OptionRetrieverSourceIP(cfg.sourceIP),
		permissions.OptionRetrieverID(cfg.id),
		permissions.OptionRetrieverRestrictions(cfg.restrictions),
	}

	perms, err := a.retriever.Permissions(ctx, claims, ns, ropts...)
	if err != nil {
		return false, err
	}

	a.cache.Set(
		ns,
		key,
		perms,
		time.Hour+time.Duration(rand.Int63n(60*30))*time.Second,
	)

	return perms.Allows(operation, resource), nil
}

func hash(claims []string, remoteaddr string, id string, restrictions permissions.Restrictions) string {
	return fmt.Sprintf("%d",
		murmur3.Sum64(
			[]byte(
				fmt.Sprintf("%s:%s:%s:%s:%s:%s",
					claims,
					remoteaddr,
					id,
					restrictions.Namespace,
					restrictions.Networks,
					restrictions.Permissions,
				),
			),
		),
	)
}
