package processors

import (
	"crypto"
	"crypto/tls"
	"crypto/x509"
	"net/http"
	"os"
	"time"

	"go.aporeto.io/a3s/internal/srv/authn/internal/issuer"
	"go.aporeto.io/a3s/pkgs/api"
	"go.aporeto.io/bahamut"
	"go.aporeto.io/elemental"
	"go.aporeto.io/manipulate"
)

// A IssueProcessor is a bahamut processor for Issue.
type IssueProcessor struct {
	manipulator manipulate.Manipulator
	jwtKey      crypto.PrivateKey
	jwtCert     *x509.Certificate
	maxValidity time.Duration
}

// NewIssueProcessor returns a new IssueProcessor.
func NewIssueProcessor(manipulator manipulate.Manipulator, cert *x509.Certificate, key crypto.PrivateKey, maxValidity time.Duration) *IssueProcessor {

	return &IssueProcessor{
		manipulator: manipulator,
		jwtCert:     cert,
		jwtKey:      key,
		maxValidity: maxValidity,
	}
}

// ProcessCreate handles the creates requests for Issue.
func (p *IssueProcessor) ProcessCreate(bctx bahamut.Context) (err error) {

	req := bctx.InputData().(*api.Issue)
	validity, _ := time.ParseDuration(req.Validity) // elemental already validated this

	switch req.SourceType {

	case api.IssueSourceTypeCertificate:
		tlsState := bctx.Request().TLSConnectionState
		if err := p.handleCertificateIssue(req, tlsState, validity); err != nil {
			return err
		}
	}

	bctx.SetOutputData(req)

	return nil
}

func (p *IssueProcessor) handleCertificateIssue(req *api.Issue, tlsState *tls.ConnectionState, validity time.Duration) (err error) {

	if tlsState == nil || len(tlsState.PeerCertificates) == 0 {
		return elemental.NewError("Bad Request", "No client certificates", "a3s", http.StatusBadRequest)
	}

	userCert := tlsState.PeerCertificates[0]

	pool := getDevPool()

	iss := issuer.NewMTLSIssuer(pool, req.SourceNamespace, req.SourceName)
	if err := iss.FromCertificate(userCert); err != nil {
		return err
	}

	idt := iss.Issue()

	req.Token, err = idt.JWT(p.jwtKey, time.Now().Add(validity))
	if err != nil {
		return err
	}
	return nil
}

func getDevPool() *x509.CertPool {
	pool := x509.NewCertPool()
	cadata, err := os.ReadFile("dev/.data/certificates/ca-acme-cert.pem")
	if err != nil {
		panic(err)
	}

	pool.AppendCertsFromPEM(cadata)

	return pool
}
