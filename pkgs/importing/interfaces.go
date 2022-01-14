package importing

import "go.aporeto.io/elemental"

// An Importable is the interface an object
// must satisfy in order to be importable.
type Importable interface {
	GetImportHash() string
	SetImportHash(string)
	GetImportLabel() string
	SetImportLabel(string)

	elemental.Namespaceable
	elemental.Identifiable
	elemental.AttributeSpecifiable
}
