package authorizer

import (
	"sync"
	"testing"

	"go.aporeto.io/elemental"
)

type mockedOperationTransformerMethods struct {
	transformMock func(elemental.Operation) string
}

// A MockOperationTransformer allows to mock a transform.OperationTransformer for unit tests.
type MockOperationTransformer interface {
	OperationTransformer
	MockTransform(t *testing.T, impl func(elemental.Operation) string)
}

type mockOperationTransformer struct {
	mocks       map[*testing.T]*mockedOperationTransformerMethods
	currentTest *testing.T

	sync.Mutex
}

// NewMockOperationTransformer returns a MockOperationTransformer.
func NewMockOperationTransformer() MockOperationTransformer {
	return &mockOperationTransformer{
		mocks: map[*testing.T]*mockedOperationTransformerMethods{},
	}
}

// MockTransform replaces the Transform implementation with the given function.
func (r *mockOperationTransformer) MockTransform(t *testing.T, impl func(elemental.Operation) string) {

	r.Lock()
	defer r.Unlock()

	r.currentMocks(t).transformMock = impl
}

func (r *mockOperationTransformer) Transform(operation elemental.Operation) string {

	r.Lock()
	defer r.Unlock()

	if mock := r.currentMocks(r.currentTest); mock != nil && mock.transformMock != nil {
		return mock.transformMock(operation)
	}

	return ""
}

func (r *mockOperationTransformer) currentMocks(t *testing.T) *mockedOperationTransformerMethods {

	mocks := r.mocks[t]

	if mocks == nil {
		mocks = &mockedOperationTransformerMethods{}
		r.mocks[t] = mocks
	}

	r.currentTest = t
	return mocks
}
