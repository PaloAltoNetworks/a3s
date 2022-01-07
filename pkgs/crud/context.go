package crud

import (
	"net/http"

	"go.aporeto.io/bahamut"
	"go.aporeto.io/elemental"
	"go.aporeto.io/manipulate"
)

// translateContext translates the given bahamut.Context to a manipulate.Context
// It handles, namespace, recursive, propagate and the `q` parameter.
// If your code needs to apply another filter, it will override the filter
// created from the query parameter.
func translateContext(bctx bahamut.Context) (manipulate.Context, error) {

	opts := []manipulate.ContextOption{
		manipulate.ContextOptionNamespace(bctx.Request().Namespace),
		manipulate.ContextOptionRecursive(bctx.Request().Recursive),
		manipulate.ContextOptionPropagated(bctx.Request().Propagated),
	}

	qfilter, err := manipulate.NewFiltersFromQueryParameters(bctx.Request().Parameters)
	if err != nil {
		return nil, elemental.NewError(
			"Bad Request",
			err.Error(),
			"a3s:policy",
			http.StatusBadRequest,
		)
	}
	if qfilter != nil {
		opts = append(opts, manipulate.ContextOptionFilter(qfilter))
	}

	return manipulate.NewContext(bctx.Context(), opts...), err
}
