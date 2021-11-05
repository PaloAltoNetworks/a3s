package processors

import (
	"net/http"

	"go.aporeto.io/a3s/pkgs/api"
	"go.aporeto.io/a3s/pkgs/authorizer"
	"go.aporeto.io/a3s/pkgs/permissions"
	"go.aporeto.io/a3s/pkgs/token"
	"go.aporeto.io/bahamut"
	"go.aporeto.io/elemental"
)

// A AuthzProcessor is a bahamut processor for Authzs.
type AuthzProcessor struct {
	authorizer authorizer.Authorizer
	jwks       *token.JWKS
}

// NewAuthzProcessor returns a new AuthzProcessor.
func NewAuthzProcessor(authorizer authorizer.Authorizer, jwks *token.JWKS) *AuthzProcessor {
	return &AuthzProcessor{
		authorizer: authorizer,
		jwks:       jwks,
	}
}

// ProcessCreate handles the creates requests for Authzs.
func (p *AuthzProcessor) ProcessCreate(bctx bahamut.Context) error {

	req := bctx.InputData().(*api.Authz)

	idt := &token.IdentityToken{}
	if err := idt.Parse(req.Token, p.jwks, "", ""); err != nil {
		return elemental.NewError(
			"Bad Request",
			err.Error(),
			"a3s:authz",
			http.StatusBadRequest,
		)
	}

	var r permissions.Restrictions
	if idt.Restrictions != nil {
		r = *idt.Restrictions
	}

	ok, err := p.authorizer.CheckAuthorization(
		bctx.Context(),
		idt.Identity,
		req.Action,
		req.Namespace,
		req.Resource,
		authorizer.OptionCheckID(req.ID),
		authorizer.OptionCheckSourceIP(req.IP),
		authorizer.OptionCheckRestrictions(r),
	)
	if err != nil {
		return err
	}

	if ok {
		bctx.SetStatusCode(http.StatusOK)
	} else {
		bctx.SetStatusCode(http.StatusForbidden)
	}

	return nil
}
