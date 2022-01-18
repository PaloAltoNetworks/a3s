package processors

import (
	"context"
	"fmt"
	"net/http"

	"go.aporeto.io/a3s/pkgs/api"
	"go.aporeto.io/a3s/pkgs/bearermanip"
	"go.aporeto.io/a3s/pkgs/importing"
	"go.aporeto.io/bahamut"
	"go.aporeto.io/elemental"
	"go.aporeto.io/manipulate"
)

// A HTTPSourcesProcessor is a bahamut processor for HTTPSource.
type ImportProcessor struct {
	bmanipMaker bearermanip.MakerFunc
}

// NewImportProcessor returns a new ImportProcessor .
func NewImportProcessor(bmanipMaker bearermanip.MakerFunc) *ImportProcessor {
	return &ImportProcessor{
		bmanipMaker: bmanipMaker,
	}
}

// ProcessCreate handles the creates requests for HTTPSource.
func (p *ImportProcessor) ProcessCreate(bctx bahamut.Context) error {

	req := bctx.InputData().(*api.Import)
	ns := bctx.Request().Namespace

	for _, lst := range []elemental.Identifiables{
		req.LDAPSources,
		req.OIDCSources,
		req.A3SSources,
		req.MTLSSources,
		req.HTTPSources,
		req.Authorizations,
	} {
		if err := importObjects(
			bctx.Context(),
			api.Manager(),
			p.bmanipMaker(bctx),
			ns,
			req.Label,
			req.Mode,
			lst,
		); err != nil {
			return err
		}
	}

	return nil
}

func importObjects(
	ctx context.Context,
	manager elemental.ModelManager,
	m manipulate.Manipulator,
	namespace string,
	label string,
	mode api.ImportModeValue,
	items elemental.Identifiables,
) error {

	lst := items.List()
	hashed := make(map[string]importing.Importable, len(lst))

	// If the mode is ImportModeRemove, we don't populate
	// the hashed list, which will end up deleting all
	// existing objects.
	if mode == api.ImportModeImport {

		for i, obj := range lst {

			imp, ok := obj.(importing.Importable)
			if !ok {
				return fmt.Errorf("object '%s[%d]' is not Importable", obj.Identity().Name, i)
			}

			h, err := importing.Hash(imp, api.Manager())
			if err != nil {
				return fmt.Errorf("unable to hash '%s[%d]': %w", obj.Identity().Name, i, err)
			}

			imp.SetImportHash(h)
			imp.SetImportLabel(label)

			hashed[h] = imp
		}
	}

	// Now, we retrieve all existing object in the namespace
	// using the same import label.
	currentObjects := manager.Identifiables(items.Identity())
	if err := m.RetrieveMany(
		manipulate.NewContext(
			ctx,
			manipulate.ContextOptionNamespace(namespace),
			manipulate.ContextOptionFilter(
				elemental.NewFilterComposer().
					WithKey("importLabel").Equals(label).
					Done(),
			),
		),
		currentObjects,
	); err != nil {
		return fmt.Errorf("unable to retrieve list of current objects: %w", err)
	}

	// Now we delete all the existing objects that have a hash
	// that is not matching any of the imported objects.
	// We also delete from the list of objects to import all the
	// ones that have a matching hash, since they did not change.
	for _, o := range currentObjects.List() {

		h := o.(importing.Importable).GetImportHash()

		if _, ok := hashed[h]; ok {
			delete(hashed, h)
			continue
		}

		if err := m.Delete(
			manipulate.NewContext(
				ctx,
				manipulate.ContextOptionNamespace(namespace),
				manipulate.ContextOptionOverride(true),
			),
			o,
		); err != nil {
			if elemental.IsErrorWithCode(err, http.StatusNotFound) {
				continue
			}
			return err
		}
	}

	// Finally, we create the remaining objects.
	for _, o := range hashed {
		if err := m.Create(
			manipulate.NewContext(
				ctx,
				manipulate.ContextOptionNamespace(namespace),
			),
			o,
		); err != nil {
			return err
		}
	}

	return nil
}
