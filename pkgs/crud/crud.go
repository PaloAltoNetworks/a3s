package crud

import (
	"net/http"

	"go.aporeto.io/a3s/pkgs/api"
	"go.aporeto.io/bahamut"
	"go.aporeto.io/elemental"
	"go.aporeto.io/manipulate"
)

func Create(bctx bahamut.Context, m manipulate.Manipulator, obj elemental.Identifiable) error {

	if n, ok := obj.(elemental.Namespaceable); ok {
		n.SetNamespace(bctx.Request().Namespace)
	}

	return m.Create(manipulate.NewContext(bctx.Context()), obj)
}

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

func Update(bctx bahamut.Context, m manipulate.Manipulator, obj elemental.Identifiable) error {

	obj.SetIdentifier(bctx.Request().ObjectID)

	mctx, err := TranslateContext(bctx)
	if err != nil {
		return err
	}

	eobj := api.Manager().Identifiable(obj.Identity())
	eobj.SetIdentifier(obj.Identifier())

	if err := m.Retrieve(mctx, eobj); err != nil {
		return elemental.NewError(
			"Not Found",
			"Object not found",
			"a3s:policy",
			http.StatusNotFound,
		)
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
