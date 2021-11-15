package authorizer

import (
	"fmt"

	"go.aporeto.io/a3s/pkgs/api"
	"go.aporeto.io/a3s/pkgs/permissions"
	"go.aporeto.io/a3s/pkgs/token"
	"go.aporeto.io/bahamut"
	"go.aporeto.io/elemental"
	"go.aporeto.io/manipulate"
	"go.uber.org/zap"
)

var pushSessionIdentity = "pushsession"

type pushedEntity struct {
	ID                string `msgpack:"ID" json:"ID"`
	Namespace         string `msgpack:"namespace" json:"namespace"`
	Name              string `msgpack:"name" json:"name"`
	Propagate         bool   `msgpack:"propagate" json:"propagate"`
	PropagationHidden bool   `msgpack:"propagationHidden" json:"propagationHidden"`
}

// A PushDispatchHandler handles dispatching events to push sessions.
type PushDispatchHandler struct {
	authorizer  Authorizer
	manipulator manipulate.Manipulator
}

// NewPushDispatchHandler returns a new PushDispatchHandler.
func NewPushDispatchHandler(manipulator manipulate.Manipulator, authorizer Authorizer) *PushDispatchHandler {

	return &PushDispatchHandler{
		manipulator: manipulator,
		authorizer:  authorizer,
	}
}

// OnPushSessionInit is part of the bahamut.PushDispatchHandler interface
func (g *PushDispatchHandler) OnPushSessionInit(session bahamut.PushSession) (bool, error) {

	restrictions, err := permissions.GetRestrictions(token.FromSession(session))
	if err != nil {
		return false, err
	}

	ok, err := g.authorizer.CheckAuthorization(
		session.Context(),
		session.Claims(),
		"get",
		session.Parameter("namespace"),
		pushSessionIdentity,
		OptionCheckSourceIP(session.ClientIP()),
		OptionCheckRestrictions(restrictions),
	)

	if err != nil {
		zap.L().Error("Unable to authorize session", zap.Error(err))
		return false, err
	}

	return ok, nil
}

// OnPushSessionStart is part of the bahamut.PushDispatchHandler interface
func (g *PushDispatchHandler) OnPushSessionStart(session bahamut.PushSession) {
	zap.L().Debug("Push session started", zap.Strings("claims", session.Claims()))
}

// OnPushSessionStop is part of the bahamut.PushDispatchHandler interface
func (g *PushDispatchHandler) OnPushSessionStop(session bahamut.PushSession) {
	zap.L().Debug("Push session stopped", zap.Strings("claims", session.Claims()))
}

// SummarizeEvent is part of the bahamut.PushDispatchHandler interface
func (g *PushDispatchHandler) SummarizeEvent(event *elemental.Event) (interface{}, error) {

	entity := pushedEntity{}
	if err := event.Decode(&entity); err != nil {
		return nil, fmt.Errorf("unable to summarize event entity: %s", err)
	}

	return entity, nil
}

// RelatedEventIdentities is part of the bahamut.PushDispatchHandler interface
func (g *PushDispatchHandler) RelatedEventIdentities(identity string) []string {
	return nil
}

// ShouldDispatch is part of the bahamut.PushDispatchHandler interface
func (g *PushDispatchHandler) ShouldDispatch(session bahamut.PushSession, event *elemental.Event, summary interface{}) (bool, error) {

	entity := summary.(pushedEntity)
	sessionNS := session.Parameter("namespace")

	isFromCurrentNS := entity.Namespace == sessionNS
	isFromParentNS := elemental.IsNamespaceParentOfNamespace(entity.Namespace, sessionNS)
	isFromChildNS := elemental.IsNamespaceChildrenOfNamespace(entity.Namespace, sessionNS)
	isRecursive := session.Parameter("mode") == "all"

	// If it's a ns delete, then all authorization are already voided,
	// we just push the event to the clients so they know it's gone.
	if event.Identity == api.NamespaceIdentity.Name &&
		event.Type == elemental.EventDelete &&
		(entity.Name == sessionNS || isFromParentNS) {
		return true, nil
	}

	// If the object is in a parent namespace or in a child namespace
	// and it's not in recursive mode, we don't push unless it is propagating.
	if !(isFromCurrentNS || (isFromChildNS && isRecursive)) {

		// If the object is not from a parent NS, we don't push.
		if !isFromParentNS {
			return false, nil
		}

		// If the object does not propagate or propagation is hidden, we don't push.
		if !entity.Propagate || entity.PropagationHidden {
			return false, nil
		}
	}

	// Finally we check if the session has the right to read
	// the object that is about to be pushed.
	restrictions, err := permissions.GetRestrictions(token.FromSession(session))
	if err != nil {
		return false, err
	}

	return g.authorizer.CheckAuthorization(
		session.Context(),
		session.Claims(),
		"get",
		sessionNS,
		string(event.Identity),
		OptionCheckRestrictions(restrictions),
		OptionCheckSourceIP(session.ClientIP()),
	)
}
