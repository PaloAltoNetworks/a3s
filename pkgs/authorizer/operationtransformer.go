package authorizer

import "go.aporeto.io/elemental"

// A OperationTransformer is an interface that can transform the operation being evaluated.
type OperationTransformer interface {
	Transform(operation elemental.Operation) string
}
