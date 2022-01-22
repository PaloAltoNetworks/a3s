package processors

import (
	"fmt"

	"go.aporeto.io/a3s/pkgs/api"
	"go.aporeto.io/a3s/pkgs/crud"
	"go.aporeto.io/a3s/pkgs/token"
	"go.aporeto.io/bahamut"
	"go.aporeto.io/elemental"
	"go.aporeto.io/manipulate"
	"go.aporeto.io/tg/tglib"
)

// A MTLSSourcesProcessor is a bahamut processor for MTLSSource.
type MTLSSourcesProcessor struct {
	manipulator manipulate.Manipulator
}

// NewMTLSSourcesProcessor returns a new MTLSSourcesProcessor.
func NewMTLSSourcesProcessor(manipulator manipulate.Manipulator) *MTLSSourcesProcessor {
	return &MTLSSourcesProcessor{
		manipulator: manipulator,
	}
}

// ProcessCreate handles the creates requests for MTLSSource.
func (p *MTLSSourcesProcessor) ProcessCreate(bctx bahamut.Context) error {
	return crud.Create(bctx, p.manipulator, bctx.InputData().(*api.MTLSSource),
		crud.OptionPreWriteHook(func(obj elemental.Identifiable, orig elemental.Identifiable) error {
			return insertReferences(obj.(*api.MTLSSource))
		}),
	)
}

// ProcessRetrieveMany handles the retrieve many requests for MTLSSource.
func (p *MTLSSourcesProcessor) ProcessRetrieveMany(bctx bahamut.Context) error {
	return crud.RetrieveMany(bctx, p.manipulator, &api.MTLSSourcesList{})
}

// ProcessRetrieve handles the retrieve requests for MTLSSource.
func (p *MTLSSourcesProcessor) ProcessRetrieve(bctx bahamut.Context) error {
	return crud.Retrieve(bctx, p.manipulator, api.NewMTLSSource())
}

// ProcessUpdate handles the update requests for MTLSSource.
func (p *MTLSSourcesProcessor) ProcessUpdate(bctx bahamut.Context) error {
	return crud.Update(bctx, p.manipulator, bctx.InputData().(*api.MTLSSource),
		crud.OptionPreWriteHook(func(obj elemental.Identifiable, orig elemental.Identifiable) error {
			return insertReferences(obj.(*api.MTLSSource))
		}),
	)
}

// ProcessDelete handles the delete requests for MTLSSource.
func (p *MTLSSourcesProcessor) ProcessDelete(bctx bahamut.Context) error {
	return crud.Delete(bctx, p.manipulator, api.NewMTLSSource())
}

// ProcessInfo handles the info request for MTLSSource.
func (p *MTLSSourcesProcessor) ProcessInfo(bctx bahamut.Context) error {
	return crud.Info(bctx, p.manipulator, api.MTLSSourceIdentity)
}

func insertReferences(src *api.MTLSSource) error {

	certs, err := tglib.ParseCertificates([]byte(src.CA))
	if err != nil {
		return err
	}

	src.Fingerprints = make([]string, len(certs))
	src.SubjectKeyIDs = make([]string, len(certs))
	for i, cert := range certs {
		src.Fingerprints[i] = token.Fingerprint(cert)
		src.SubjectKeyIDs[i] = fmt.Sprintf("%02X", cert.SubjectKeyId)
	}

	return nil
}
