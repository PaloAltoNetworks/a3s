package processors

import (
	"context"
	"crypto/tls"
	"crypto/x509"
	"fmt"
	"net/http"
	"time"

	"go.aporeto.io/a3s/internal/issuer"
	"go.aporeto.io/a3s/pkgs/api"
	"go.aporeto.io/a3s/pkgs/token"
	"go.aporeto.io/bahamut"
	"go.aporeto.io/elemental"
	"go.aporeto.io/manipulate"
)

// A IssueProcessor is a bahamut processor for Issue.
type IssueProcessor struct {
	manipulator manipulate.Manipulator
	jwks        *token.JWKS
	maxValidity time.Duration
}

// NewIssueProcessor returns a new IssueProcessor.
func NewIssueProcessor(manipulator manipulate.Manipulator, jwks *token.JWKS, maxValidity time.Duration) *IssueProcessor {

	return &IssueProcessor{
		manipulator: manipulator,
		jwks:        jwks,
		maxValidity: maxValidity,
	}
}

// ProcessCreate handles the creates requests for Issue.
func (p *IssueProcessor) ProcessCreate(bctx bahamut.Context) (err error) {

	req := bctx.InputData().(*api.Issue)
	validity, _ := time.ParseDuration(req.Validity) // elemental already validated this

	switch req.SourceType {

	case api.IssueSourceTypeMTLS:
		tlsState := bctx.Request().TLSConnectionState
		if err := p.handleCertificateIssue(bctx.Context(), req, tlsState, validity); err != nil {
			return err
		}
	}

	bctx.SetOutputData(req)

	return nil
}

func (p *IssueProcessor) handleCertificateIssue(ctx context.Context, req *api.Issue, tlsState *tls.ConnectionState, validity time.Duration) (err error) {

	if tlsState == nil || len(tlsState.PeerCertificates) == 0 {
		return elemental.NewError("Bad Request", "No client certificates", "a3s", http.StatusBadRequest)
	}

	out, err := retrieveSource(ctx, p.manipulator, req.SourceNamespace, req.SourceName, api.MTLSSourceIdentity)
	if err != nil {
		return err
	}
	src := out.(*api.MTLSSource)

	pool := x509.NewCertPool()
	pool.AppendCertsFromPEM([]byte(src.CertificateAuthority))

	userCert := tlsState.PeerCertificates[0]
	iss := issuer.NewMTLSIssuer(pool, req.SourceNamespace, req.SourceName)
	if err := iss.FromCertificate(userCert); err != nil {
		return err
	}

	idt := iss.Issue()

	req.Token, err = idt.JWT(p.jwks.GetLast().PrivateKey(), "kid-placeholder", time.Now().Add(validity))
	if err != nil {
		return err
	}

	return nil
}

func retrieveSource(ctx context.Context, m manipulate.Manipulator, namespace string, name string, identity elemental.Identity) (elemental.Identifiable, error) {

	if namespace == "" {
		return nil, elemental.NewError(
			"Bad Request",
			"You must set sourceNamespace and sourceName",
			"a3s:auth",
			http.StatusBadRequest,
		)
	}

	if name == "" {
		return nil, elemental.NewError(
			"Bad Request",
			"You must set sourceNamespace and sourceName",
			"a3s:auth",
			http.StatusBadRequest,
		)
	}

	mctx := manipulate.NewContext(ctx,
		manipulate.ContextOptionNamespace(namespace),
		manipulate.ContextOptionFilter(
			elemental.NewFilterComposer().WithKey("name").Equals(name).
				Done(),
		),
	)

	identifiables := api.Manager().IdentifiablesFromString(identity.Name)
	if err := m.RetrieveMany(mctx, identifiables); err != nil {
		return nil, err
	}

	lst := identifiables.List()
	switch len(lst) {
	case 0:
		return nil, elemental.NewError(
			"Not Found",
			"Unable to find the request auth source",
			"a3s:authn",
			http.StatusNotFound,
		)
	case 1:
	default:
		return nil, fmt.Errorf("more than one auth source found")
	}

	return lst[0], nil
}
