package permissions

import (
	"context"
	"fmt"
	"net"
	"strings"

	mapset "github.com/deckarep/golang-set"
	"go.aporeto.io/a3s/pkgs/api"
	"go.aporeto.io/elemental"
	"go.aporeto.io/manipulate"
)

// A Retriever is an object that can retrieve permissions
// for the given informations.
type Retriever interface {

	// Permissions returns the PermissionMap for the given
	// clams on the given namespace for the given id (optional)
	// from the given address with the given restrictions.
	Permissions(ctx context.Context, claims []string, ns string, opts ...RetrieverOption) (PermissionMap, error)
}

type retriever struct {
	manipulator manipulate.Manipulator
}

// NewRetriever returns a new Retriever.
func NewRetriever(manipulator manipulate.Manipulator) Retriever {
	return &retriever{
		manipulator: manipulator,
	}
}

func (a *retriever) Permissions(ctx context.Context, claims []string, ns string, opts ...RetrieverOption) (PermissionMap, error) {

	cfg := &config{}
	for _, o := range opts {
		o(cfg)
	}

	// Handle token's authorizedNamespace.
	if cfg.restrictions.Namespace != "" {
		if cfg.restrictions.Namespace != ns && !elemental.IsNamespaceParentOfNamespace(cfg.restrictions.Namespace, ns) {
			return nil, nil
		}
	}

	if ns != "/" {
		count, err := a.countNamespace(ctx, ns)

		if err != nil {
			return nil, err
		}

		if count != 1 {
			return nil, nil // we don't return the error to the client or some namespace names may leak.
		}
	}

	policies, err := a.resolvePoliciesMatchingClaims(ctx, claims, ns)
	if err != nil {
		return nil, fmt.Errorf("unable to retrieve api authorizations: %s", err)
	}

	out := PermissionMap{}
	for _, p := range policies {

		if len(p.Subject) == 0 || len(p.Subject[0]) == 0 {
			continue
		}

		var nsMatch bool
		for _, targetNS := range p.TargetNamespaces {
			if ns == targetNS || elemental.IsNamespaceChildrenOfNamespace(ns, targetNS) {
				nsMatch = true
				break
			}
		}

		if !nsMatch {
			continue
		}

		if l := len(p.Subnets); l > 0 {

			allowedSubnets := map[string]interface{}{}
			for _, sub := range p.Subnets {
				allowedSubnets[sub] = struct{}{}
			}

			valid, err := validateClientIP(cfg.addr, allowedSubnets)
			if err != nil {
				return nil, err
			}
			if !valid {
				continue
			}
		}

		for identity, perms := range Parse(p.Permissions, cfg.id) {
			if _, ok := out[identity]; !ok {
				out[identity] = perms
			} else {
				for verb := range perms {
					out[identity][verb] = true
				}
			}
		}
	}

	// If we have restrictions on permission from the token,
	// we reduce the
	if len(cfg.restrictions.Permissions) > 0 {
		out = out.Intersect(Parse(cfg.restrictions.Permissions, cfg.id))
	}

	// If we have restrictions on the origin networks from the token
	// we verify here.
	if len(cfg.restrictions.Networks) > 0 {
		allowedSubnets := map[string]interface{}{}
		for _, net := range cfg.restrictions.Networks {
			allowedSubnets[net] = struct{}{}
		}
		valid, err := validateClientIP(cfg.addr, allowedSubnets)
		if err != nil {
			return nil, err
		}
		if !valid {
			return nil, nil
		}
	}

	return out, nil
}

func (a *retriever) resolvePoliciesMatchingClaims(ctx context.Context, claims []string, ns string) (api.AuthorizationsList, error) {

	mctx := manipulate.NewContext(
		ctx,
		manipulate.ContextOptionNamespace(ns),
		manipulate.ContextOptionPropagated(true),
		manipulate.ContextOptionFilter(
			makeAPIAuthorizationPolicyRetrieveFilter(claims),
		),
	)

	// Find all policies that are matching at least one claim
	policies := api.AuthorizationsList{}
	if err := a.manipulator.RetrieveMany(mctx, &policies); err != nil {
		return nil, err
	}

	// Ignore policies that are not matching all claims
	matchingPolicies := []*api.Authorization{}
	for _, p := range policies {
		if match(p.Subject, claims) {
			matchingPolicies = append(matchingPolicies, p)
		}
	}

	return matchingPolicies, nil
}

func validateClientIP(remoteAddr string, allowedSubnets map[string]interface{}) (bool, error) {

	ipStr, _, err := net.SplitHostPort(remoteAddr)
	if err != nil {
		ipStr = remoteAddr
	}

	ip := net.ParseIP(ipStr)
	if ip == nil {
		return false, fmt.Errorf("missing or invalid origin IP '%s'", ipStr)
	}

	if ip.IsLoopback() {
		ip = net.IPv4(127, 0, 0, 1)
	} else {
		ip = ip.To4()
	}

	for sub := range allowedSubnets {

		_, subnet, err := net.ParseCIDR(sub)
		if err != nil {
			return false, err
		}

		if subnet.Contains(ip) {
			return true, nil
		}
	}

	return false, nil
}

// makeAPIAuthorizationPolicyRetrieveFilter creates a manipulate filter to retrieve the api authorization policies matching the claims.
func makeAPIAuthorizationPolicyRetrieveFilter(claims []string) *elemental.Filter {

	itags := []interface{}{}
	set := mapset.NewSet()
	var issuer string
	for _, tag := range claims {
		if set.Add(tag) {
			itags = append(itags, tag)
			if strings.HasPrefix(tag, "@issuer=") {
				issuer = strings.TrimPrefix(tag, "@issuer=")
			}
		}
	}

	filter := elemental.NewFilterComposer().
		WithKey("flattenedsubject").In(itags...).
		WithKey("disabled").Equals(false)

	if issuer != "" {
		filter.WithKey("trustedissuers").Contains(issuer)
	}

	return filter.Done()
}

// countNamespace tries to find the namespace in a two step process.
func (a *retriever) countNamespace(ctx context.Context, ns string) (int, error) {

	var count int
	var err error

	filter := elemental.NewFilterComposer().WithKey("name").Equals(ns).Done()

	if count, err = a.manipulator.Count(
		manipulate.NewContext(
			ctx,
			manipulate.ContextOptionFilter(filter),
			manipulate.ContextOptionRecursive(true),
		),
		api.NamespaceIdentity,
	); err != nil {
		return 0, err
	}

	if count == 0 {
		// If we could not find a namespace on the first attempt
		// try it a second time with strong read consistency,
		count, err = a.manipulator.Count(
			manipulate.NewContext(
				ctx,
				manipulate.ContextOptionReadConsistency(manipulate.ReadConsistencyStrong),
				manipulate.ContextOptionFilter(filter),
				manipulate.ContextOptionRecursive(true),
			),
			api.NamespaceIdentity,
		)
	}

	return count, err
}

func match(expression [][]string, tags []string) bool {

	tm := mapset.NewSetFromSlice(stringListToInterfaceList(tags))

	for _, ands := range expression {
		if mapset.NewSetFromSlice(stringListToInterfaceList(ands)).IsSubset(tm) {
			return true
		}
	}

	return false
}

func stringListToInterfaceList(in []string) (out []interface{}) {
	out = make([]interface{}, len(in))
	for i, s := range in {
		out[i] = s
	}

	return out
}
