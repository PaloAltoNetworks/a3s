package permissions

import (
	"sync"
	"testing"
)

type mockedTransformerMethods struct {
	transformMock func(PermissionMap, Restrictions) PermissionMap
}

// A MockTransformer allows to mock a transform.Transformer for unit tests.
type MockTransformer interface {
	Transformer
	MockTransform(t *testing.T, impl func(PermissionMap, Restrictions) PermissionMap)
}

type mockTransformer struct {
	mocks       map[*testing.T]*mockedTransformerMethods
	currentTest *testing.T

	sync.Mutex
}

// NewMockTransformer returns a MockTransformer.
func NewMockTransformer() MockTransformer {
	return &mockTransformer{
		mocks: map[*testing.T]*mockedTransformerMethods{},
	}
}

// MockTransform replaces the Transform implementation with the given function.
func (r *mockTransformer) MockTransform(t *testing.T, impl func(PermissionMap, Restrictions) PermissionMap) {

	r.Lock()
	defer r.Unlock()

	r.currentMocks(t).transformMock = impl
}

func (r *mockTransformer) Transform(permissions PermissionMap, restrictions Restrictions) PermissionMap {

	r.Lock()
	defer r.Unlock()

	if mock := r.currentMocks(r.currentTest); mock != nil && mock.transformMock != nil {
		return mock.transformMock(permissions, restrictions)
	}

	return PermissionMap{}
}

func (r *mockTransformer) currentMocks(t *testing.T) *mockedTransformerMethods {

	mocks := r.mocks[t]

	if mocks == nil {
		mocks = &mockedTransformerMethods{}
		r.mocks[t] = mocks
	}

	r.currentTest = t
	return mocks
}
