package importing

import (
	"context"
	"fmt"
	"net/http"

	"go.aporeto.io/elemental"
	"go.aporeto.io/manipulate"
)

// Import preforms the importing of the given
// objects, in the given namespace, with the given label
// using the given manipulator.
// If removeMode is true, all the objects with the given
// label will be deleted.
// This function does not make any permission check, and will
// fail if the given manipulator does not bear sufficient permissions.
func Import(
	ctx context.Context,
	manager elemental.ModelManager,
	m manipulate.Manipulator,
	namespace string,
	label string,
	objects elemental.Identifiables,
	removeMode bool,
) error {

	if namespace == "" {
		return fmt.Errorf("namespace must not be empty")
	}

	if label == "" {
		return fmt.Errorf("label must not be empty")
	}

	lst := objects.List()
	hashed := make(map[string]Importable, len(lst))

	// If the mode is ImportModeRemove, we don't populate
	// the hashed list, which will end up deleting all
	// existing objects.
	if !removeMode {

		for i, obj := range lst {

			imp, ok := obj.(Importable)
			if !ok {
				return fmt.Errorf("object '%s[%d]' is not importable", obj.Identity().Name, i)
			}

			h, err := Hash(imp, manager)
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
	currentObjects := manager.Identifiables(objects.Identity())
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

	// Then, we delete all the existing objects that have a hash
	// that is not matching any of the imported objects.
	// We also delete from the list of objects to import all the
	// ones that have a matching hash, since they did not change.
	for _, o := range currentObjects.List() {

		h := o.(Importable).GetImportHash()

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
			return fmt.Errorf("unable to delete existing object: %w", err)
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
			return fmt.Errorf("unable to create imported object: %w", err)
		}
	}

	return nil
}
