package crud

import (
	"go.aporeto.io/a3s/pkgs/api"
	"go.aporeto.io/bahamut"
	"go.aporeto.io/elemental"
	"go.aporeto.io/manipulate"
)

// Create performs the basic create operation.
func Create(bctx bahamut.Context, m manipulate.Manipulator, obj elemental.Identifiable) error {

	if n, ok := obj.(elemental.Namespaceable); ok {
		n.SetNamespace(bctx.Request().Namespace)
	}

	return m.Create(manipulate.NewContext(bctx.Context()), obj)
}

// RetrieveMany performs the basic retrieve many operation.
func RetrieveMany(bctx bahamut.Context, m manipulate.Manipulator, objs elemental.Identifiables) error {

	mctx, err := TranslateContext(bctx)
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

	mctx, err := TranslateContext(bctx)
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

// Update performs the basic update operation.
func Update(bctx bahamut.Context, m manipulate.Manipulator, obj elemental.Identifiable) error {

	obj.SetIdentifier(bctx.Request().ObjectID)

	mctx, err := TranslateContext(bctx)
	if err != nil {
		return err
	}

	eobj := api.Manager().Identifiable(obj.Identity())
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

	if err := m.Update(mctx, obj); err != nil {
		return err
	}

	bctx.SetOutputData(obj)

	return nil
}

// Delete performs the basic delete operation.
func Delete(bctx bahamut.Context, m manipulate.Manipulator, obj elemental.Identifiable) error {

	mctx, err := TranslateContext(bctx)
	if err != nil {
		return err
	}

	obj.SetIdentifier(bctx.Request().ObjectID)
	if err := m.Retrieve(mctx, obj); err != nil {
		return err
	}

	bctx.SetOutputData(obj)

	return m.Delete(mctx, obj)
}

// Info performs the basic info operation.
func Info(bctx bahamut.Context, m manipulate.Manipulator, identity elemental.Identity) error {

	mctx, err := TranslateContext(bctx)
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
