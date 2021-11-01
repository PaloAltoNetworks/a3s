package processors

import (
	"crypto"
	"crypto/x509"

	"go.aporeto.io/bahamut"
	"go.aporeto.io/manipulate"
)

// A Issues is a bahamut processor for Issuess.
type IssueProcessor struct {
	manipulator    manipulate.Manipulator
	jwtSigningKey  crypto.PrivateKey
	jwtSigningCert *x509.Certificate
}

// NewIssueProcessor returns a new IssuessProcessor.
func NewIssueProcessor(manipulator manipulate.Manipulator, cert *x509.Certificate, key crypto.PrivateKey) *IssueProcessor {

	return &IssueProcessor{
		manipulator:    manipulator,
		jwtSigningCert: cert,
		jwtSigningKey:  key,
	}
}

// ProcessCreate handles the creates requests for Issuess.
func (p *IssueProcessor) ProcessCreate(ctx bahamut.Context) error {

	return nil
}
