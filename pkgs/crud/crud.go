package crud

import (
	"errors"
	"fmt"
	"reflect"

	"go.aporeto.io/a3s/pkgs/importing"
	"go.aporeto.io/bahamut"
	"go.aporeto.io/elemental"
	"go.aporeto.io/manipulate"
)

// Create performs the basic create operation.
func Create(bctx bahamut.Context, m manipulate.Manipulator, obj elemental.Identifiable, opts ...Option) error {

	cfg := cfg{}
	for _, o := range opts {
		o(&cfg)
	}

	if as, ok := obj.(elemental.AttributeSpecifiable); ok {
		if err := elemental.ValidateAdvancedSpecification(as, nil, elemental.OperationCreate); err != nil {
			return err
		}
	}

	if n, ok := obj.(elemental.Namespaceable); ok {
		n.SetNamespace(bctx.Request().Namespace)
	}

	if cfg.preHook != nil {
		if err := cfg.preHook(obj, nil); err != nil {
			return ErrPreWriteHook{Err: err}
		}
	}

	if err := m.Create(manipulate.NewContext(bctx.Context()), obj); err != nil {
		return err
	}

	if cfg.postHook != nil {
		cfg.postHook(obj)
	}

	bctx.SetOutputData(obj)

	return nil
}

// RetrieveMany performs the basic retrieve many operation.
func RetrieveMany(bctx bahamut.Context, m manipulate.Manipulator, objs elemental.Identifiables, opts ...Option) error {

	mctx, err := translateContext(bctx)
	if err != nil {
		return err
	}

	if err := m.RetrieveMany(mctx, objs); err != nil {
		return err
	}

	bctx.SetOutputData(objs)

	return nil
}

// Retrieve performs the basic retrieve operation.
func Retrieve(bctx bahamut.Context, m manipulate.Manipulator, obj elemental.Identifiable) error {

	mctx, err := translateContext(bctx)
	if err != nil {
		return err
	}

	obj.SetIdentifier(bctx.Request().ObjectID)
	if err := m.Retrieve(mctx, obj); err != nil {
		return err
	}

	bctx.SetOutputData(obj)

	return nil
}

func newIdentifiable(obj elemental.Identifiable) (elemental.Identifiable, error) {

	objType := reflect.TypeOf(obj).Elem()
	newObjType := reflect.New(objType).Interface()
	identifiable, ok := newObjType.(elemental.Identifiable)
	if !ok {
		return nil, errors.New("object is not elemental.Identifiable")
	}

	return identifiable, nil
}

// Update performs the basic update operation.
func Update(bctx bahamut.Context, m manipulate.Manipulator, obj elemental.Identifiable, opts ...Option) error {

	cfg := cfg{}
	for _, o := range opts {
		o(&cfg)
	}

	obj.SetIdentifier(bctx.Request().ObjectID)

	mctx, err := translateContext(bctx)
	if err != nil {
		return fmt.Errorf("unable to create new identifiable to retrieve state: %w", err)
	}

	eobj, err := newIdentifiable(obj)
	if err != nil {
		return err
	}

	eobj.SetIdentifier(obj.Identifier())

	if err := m.Retrieve(mctx, eobj); err != nil {
		return err
	}

	if a, ok := obj.(elemental.AttributeSpecifiable); ok {
		elemental.BackportUnexposedFields(
			eobj.(elemental.AttributeSpecifiable),
			a,
		)
	}

	if _, ok := obj.(elemental.AttributeSpecifiable); ok {
		if err = elemental.ValidateAdvancedSpecification(
			obj.(elemental.AttributeSpecifiable),
			eobj.(elemental.AttributeSpecifiable),
			elemental.OperationUpdate,
		); err != nil {
			return err
		}
	}

	if cfg.preHook != nil {
		if err := cfg.preHook(obj, eobj); err != nil {
			return ErrPreWriteHook{Err: err}
		}
	}

	if err := m.Update(mctx, obj); err != nil {
		return err
	}

	if cfg.postHook != nil {
		cfg.postHook(obj)
	}

	// We now reset the import hash, if any
	if imp, ok := obj.(importing.Importable); ok {
		imp.SetImportHash("")
	}

	bctx.SetOutputData(obj)

	return nil
}

// Delete performs the basic delete operation.
func Delete(bctx bahamut.Context, m manipulate.Manipulator, obj elemental.Identifiable, opts ...Option) error {

	cfg := cfg{}
	for _, o := range opts {
		o(&cfg)
	}

	mctx, err := translateContext(bctx)
	if err != nil {
		return err
	}

	obj.SetIdentifier(bctx.Request().ObjectID)
	if err := m.Retrieve(mctx, obj); err != nil {
		return err
	}

	bctx.SetOutputData(obj)

	if cfg.preHook != nil {
		if err := cfg.preHook(obj, nil); err != nil {
			return ErrPreWriteHook{Err: err}
		}
	}

	if err := m.Delete(mctx, obj); err != nil {
		return err
	}

	if cfg.postHook != nil {
		cfg.postHook(obj)
	}

	return nil
}

// Info performs the basic info operation.
func Info(bctx bahamut.Context, m manipulate.Manipulator, identity elemental.Identity) error {

	mctx, err := translateContext(bctx)
	if err != nil {
		return err
	}

	c, err := m.Count(mctx, identity)
	if err != nil {
		return err
	}

	bctx.SetCount(c)

	return nil
}
