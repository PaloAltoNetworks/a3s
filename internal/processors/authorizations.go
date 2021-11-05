package processors

import (
	"fmt"
	"net/http"
	"sort"

	"go.aporeto.io/a3s/pkgs/api"
	"go.aporeto.io/a3s/pkgs/crud"
	"go.aporeto.io/a3s/pkgs/notification"
	"go.aporeto.io/a3s/pkgs/nscache"
	"go.aporeto.io/a3s/pkgs/permissions"
	"go.aporeto.io/a3s/pkgs/token"
	"go.aporeto.io/bahamut"
	"go.aporeto.io/elemental"
	"go.aporeto.io/manipulate"
)

// A AuthorizationsProcessor is a bahamut processor for Authorizations.
type AuthorizationsProcessor struct {
	manipulator manipulate.Manipulator
	retriever   permissions.Retriever
	pubsub      bahamut.PubSubClient
}

// NewAuthorizationProcessor returns a new AuthorizationsProcessor.
func NewAuthorizationProcessor(manipulator manipulate.Manipulator, pubsub bahamut.PubSubClient, retriever permissions.Retriever) *AuthorizationsProcessor {
	return &AuthorizationsProcessor{
		manipulator: manipulator,
		pubsub:      pubsub,
		retriever:   retriever,
	}
}

// ProcessCreate handles the creates requests for Authorizations.
func (p *AuthorizationsProcessor) ProcessCreate(bctx bahamut.Context) error {
	return crud.Create(bctx, p.manipulator, bctx.InputData().(*api.Authorization),
		crud.OptionPreWriteHook(p.makePreHook(bctx)),
		crud.OptionPostWriteHook(p.makeNotify()),
	)
}

// ProcessRetrieveMany handles the retrieve many requests for Authorizations.
func (p *AuthorizationsProcessor) ProcessRetrieveMany(bctx bahamut.Context) error {
	return crud.RetrieveMany(bctx, p.manipulator, &api.AuthorizationsList{})
}

// ProcessRetrieve handles the retrieve requests for Authorizations.
func (p *AuthorizationsProcessor) ProcessRetrieve(bctx bahamut.Context) error {
	return crud.Retrieve(bctx, p.manipulator, api.NewAuthorization())
}

// ProcessUpdate handles the update requests for Authorizations.
func (p *AuthorizationsProcessor) ProcessUpdate(bctx bahamut.Context) error {
	return crud.Update(bctx, p.manipulator, bctx.InputData().(*api.Authorization),
		crud.OptionPostWriteHook(p.makeNotify()),
	)
}

// ProcessDelete handles the delete requests for Authorizations.
func (p *AuthorizationsProcessor) ProcessDelete(bctx bahamut.Context) error {
	return crud.Delete(bctx, p.manipulator, api.NewAuthorization(),
		crud.OptionPostWriteHook(p.makeNotify()),
	)
}

// ProcessInfo handles the info request for Authorizations.
func (p *AuthorizationsProcessor) ProcessInfo(bctx bahamut.Context) error {
	return crud.Info(bctx, p.manipulator, api.AuthorizationIdentity)
}

func (p *AuthorizationsProcessor) makeNotify() crud.PostWriteHook {
	return func(obj elemental.Identifiable) {
		_ = notification.Publish(
			p.pubsub,
			nscache.NotificationNamespaceChanges,
			&notification.Message{
				Data: obj.(*api.Authorization).Namespace,
			},
		)
	}
}

func (p *AuthorizationsProcessor) makePreHook(ctx bahamut.Context) crud.PreWriteHook {

	return func(obj elemental.Identifiable, original elemental.Identifiable) error {

		auth := obj.(*api.Authorization)
		auth.FlattenedSubject = flattenTags(auth.Subject)

		req := ctx.Request()
		token := token.FromRequest(req)

		restrictions, err := permissions.GetRestrictions(token)
		if err != nil {
			return fmt.Errorf("unable to retrieve restrictions: %s", err)
		}

		perms, err := p.retriever.Permissions(
			ctx.Context(),
			ctx.Claims(),
			req.Namespace,
			permissions.OptionRetrieverSourceIP(req.ClientIP),
			permissions.OptionRetrieverRestrictions(restrictions),
		)
		if err != nil {
			return err
		}

		if !permissions.Contains(perms, permissions.Parse(auth.Permissions, "")) {
			return elemental.NewErrorWithData(
				"Validation Error",
				"You cannot create an APIAuthorization with more privileges than your current ones.",
				"a3s:policy",
				http.StatusUnprocessableEntity,
				map[string]interface{}{"attribute": "authorizedIdentities"},
			)
		}

		return validatePolicyTargetNamespace(auth.TargetNamespace, ctx.Request().Namespace)
	}
}

func validatePolicyTargetNamespace(targetNamespace string, requestNamespace string) error {

	if targetNamespace == requestNamespace {
		return nil
	}

	if elemental.IsNamespaceParentOfNamespace(targetNamespace, requestNamespace) {
		return elemental.NewErrorWithData(
			"Invalid Target Namespace",
			"You cannot set TargetNamespace to a parent namespace",
			"gaia",
			http.StatusUnprocessableEntity,
			map[string]interface{}{"attribute": "targetNamespace"},
		)
	}

	if !elemental.IsNamespaceChildrenOfNamespace(targetNamespace, requestNamespace) {
		return elemental.NewErrorWithData(
			"Invalid Target Namespace",
			"You cannot set TargetNamespace to a sibling namespace",
			"gaia",
			http.StatusUnprocessableEntity,
			map[string]interface{}{"attribute": "targetNamespace"},
		)
	}

	return nil
}

func flattenTags(term [][]string) (out []string) {

	set := map[string]struct{}{}

	for _, rows := range term {
		for _, r := range rows {
			set[r] = struct{}{}
		}
	}

	out = make([]string, len(set))
	var i int
	for k := range set {
		out[i] = k
		i++
	}

	sort.Strings(out)

	return out
}
