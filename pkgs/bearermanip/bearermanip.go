package bearermanip

import (
	"context"
	"crypto/tls"

	"go.aporeto.io/a3s/pkgs/token"
	"go.aporeto.io/bahamut"
	"go.aporeto.io/elemental"
	"go.aporeto.io/manipulate"
	"go.aporeto.io/manipulate/maniphttp"
)

// MakerFunc is a functon you can use to create a bearer manipulator.
type MakerFunc func(bahamut.Context) manipulate.Manipulator

// Configure returns a function that can be used to create bearer manipulators that will act on
// behalf of the calling user.
func Configure(ctx context.Context, api string, tlsConfig *tls.Config, options ...maniphttp.Option) (MakerFunc, error) {

	opts := append(
		[]maniphttp.Option{
			maniphttp.OptionDisableCompression(),
		},
		options...,
	)

	// Always add back given tls config after everything else.
	opts = append(opts, maniphttp.OptionTLSConfig(tlsConfig))

	m, err := maniphttp.New(ctx, api, opts...)
	if err != nil {
		return nil, err
	}

	return func(ctx bahamut.Context) manipulate.Manipulator {
		return &bearerManipulator{
			apiManipulator: m,
			token:          token.FromRequest(ctx.Request()),
			namespace:      ctx.Request().Namespace,
			clientIP:       ctx.Request().ClientIP,
			requestContext: ctx.Context(),
		}
	}, nil
}

// bearerManipulator extends the basic API manipuator by allowing custom credentials
// for the API calls.
type bearerManipulator struct {
	apiManipulator manipulate.Manipulator
	token          string
	clientIP       string
	namespace      string
	requestContext context.Context
}

// RetrieveMany retrieves the a list of objects with the given elemental.Identity and put them in the given dest.
func (b *bearerManipulator) RetrieveMany(mctx manipulate.Context, dest elemental.Identifiables) error {
	return b.apiManipulator.RetrieveMany(b.extendContext(mctx), dest)
}

// Retrieve retrieves one or multiple elemental.Identifiables.
func (b *bearerManipulator) Retrieve(mctx manipulate.Context, object elemental.Identifiable) error {
	return b.apiManipulator.Retrieve(b.extendContext(mctx), object)
}

// Create creates a the given elemental.Identifiables.
func (b *bearerManipulator) Create(mctx manipulate.Context, object elemental.Identifiable) error {
	return b.apiManipulator.Create(b.extendContext(mctx), object)
}

// Update updates one or multiple elemental.Identifiables.
func (b *bearerManipulator) Update(mctx manipulate.Context, object elemental.Identifiable) error {
	return b.apiManipulator.Update(b.extendContext(mctx), object)
}

// Delete deletes one or multiple elemental.Identifiables.
func (b *bearerManipulator) Delete(mctx manipulate.Context, object elemental.Identifiable) error {
	return b.apiManipulator.Delete(b.extendContext(mctx), object)
}

// DeleteMany deletes all objects of with the given identity.
func (b *bearerManipulator) DeleteMany(mctx manipulate.Context, identity elemental.Identity) error {
	return b.apiManipulator.DeleteMany(b.extendContext(mctx), identity)
}

// Count returns the number of objects with the given identity.
func (b *bearerManipulator) Count(mctx manipulate.Context, identity elemental.Identity) (int, error) {
	return b.apiManipulator.Count(b.extendContext(mctx), identity)
}

// extendContext extends the manipulate context based on the current request processed.
func (b *bearerManipulator) extendContext(mctx manipulate.Context) manipulate.Context {

	if mctx == nil {
		mctx = manipulate.NewContext(b.requestContext)
	}

	namespace := mctx.Namespace()
	if namespace == "" {
		namespace = b.namespace
	}

	return mctx.Derive(
		manipulate.ContextOptionNamespace(namespace),
		manipulate.ContextOptionToken(b.token),
		manipulate.ContextOptionClientIP(b.clientIP),
	)
}
