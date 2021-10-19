package api

import (
	"fmt"

	"github.com/globalsign/mgo/bson"
	"github.com/mitchellh/copystructure"
	"go.aporeto.io/elemental"
)

// NamespacesIdentity represents the Identity of the object.
var NamespacesIdentity = elemental.Identity{
	Name:     "namespace",
	Category: "namespaces",
	Package:  "policy",
	Private:  false,
}

// NamespacesList represents a list of Namespaces
type NamespacesList []*Namespaces

// Identity returns the identity of the objects in the list.
func (o NamespacesList) Identity() elemental.Identity {

	return NamespacesIdentity
}

// Copy returns a pointer to a copy the NamespacesList.
func (o NamespacesList) Copy() elemental.Identifiables {

	copy := append(NamespacesList{}, o...)
	return &copy
}

// Append appends the objects to the a new copy of the NamespacesList.
func (o NamespacesList) Append(objects ...elemental.Identifiable) elemental.Identifiables {

	out := append(NamespacesList{}, o...)
	for _, obj := range objects {
		out = append(out, obj.(*Namespaces))
	}

	return out
}

// List converts the object to an elemental.IdentifiablesList.
func (o NamespacesList) List() elemental.IdentifiablesList {

	out := make(elemental.IdentifiablesList, len(o))
	for i := 0; i < len(o); i++ {
		out[i] = o[i]
	}

	return out
}

// DefaultOrder returns the default ordering fields of the content.
func (o NamespacesList) DefaultOrder() []string {

	return []string{}
}

// ToSparse returns the NamespacesList converted to SparseNamespacesList.
// Objects in the list will only contain the given fields. No field means entire field set.
func (o NamespacesList) ToSparse(fields ...string) elemental.Identifiables {

	out := make(SparseNamespacesList, len(o))
	for i := 0; i < len(o); i++ {
		out[i] = o[i].ToSparse(fields...).(*SparseNamespaces)
	}

	return out
}

// Version returns the version of the content.
func (o NamespacesList) Version() int {

	return 1
}

// Namespaces represents the model of a namespace
type Namespaces struct {
	// ID is the identifier of the object.
	ID string `json:"ID" msgpack:"ID" bson:"-" mapstructure:"ID,omitempty"`

	// The name of the namespace. When you create a namespace, only put its bare name,
	// not its full path.
	Name string `json:"name" msgpack:"name" bson:"name" mapstructure:"name,omitempty"`

	// The namespace of the object.
	Namespace string `json:"namespace" msgpack:"namespace" bson:"namespace" mapstructure:"namespace,omitempty"`

	ModelVersion int `json:"-" msgpack:"-" bson:"_modelversion"`
}

// NewNamespaces returns a new *Namespaces
func NewNamespaces() *Namespaces {

	return &Namespaces{
		ModelVersion: 1,
	}
}

// Identity returns the Identity of the object.
func (o *Namespaces) Identity() elemental.Identity {

	return NamespacesIdentity
}

// Identifier returns the value of the object's unique identifier.
func (o *Namespaces) Identifier() string {

	return o.ID
}

// SetIdentifier sets the value of the object's unique identifier.
func (o *Namespaces) SetIdentifier(id string) {

	o.ID = id
}

// GetBSON implements the bson marshaling interface.
// This is used to transparently convert ID to MongoDBID as ObectID.
func (o *Namespaces) GetBSON() (interface{}, error) {

	if o == nil {
		return nil, nil
	}

	s := &mongoAttributesNamespaces{}

	if o.ID != "" {
		s.ID = bson.ObjectIdHex(o.ID)
	}
	s.Name = o.Name
	s.Namespace = o.Namespace

	return s, nil
}

// SetBSON implements the bson marshaling interface.
// This is used to transparently convert ID to MongoDBID as ObectID.
func (o *Namespaces) SetBSON(raw bson.Raw) error {

	if o == nil {
		return nil
	}

	s := &mongoAttributesNamespaces{}
	if err := raw.Unmarshal(s); err != nil {
		return err
	}

	o.ID = s.ID.Hex()
	o.Name = s.Name
	o.Namespace = s.Namespace

	return nil
}

// Version returns the hardcoded version of the model.
func (o *Namespaces) Version() int {

	return 1
}

// BleveType implements the bleve.Classifier Interface.
func (o *Namespaces) BleveType() string {

	return "namespace"
}

// DefaultOrder returns the list of default ordering fields.
func (o *Namespaces) DefaultOrder() []string {

	return []string{}
}

// Doc returns the documentation for the object
func (o *Namespaces) Doc() string {

	return `A namespace is grouping object. Every object is part of a namespace, and every
request is made against a namespace. Namespaces form a tree hierarchy.`
}

func (o *Namespaces) String() string {

	return fmt.Sprintf("<%s:%s>", o.Identity().Name, o.Identifier())
}

// GetID returns the ID of the receiver.
func (o *Namespaces) GetID() string {

	return o.ID
}

// SetID sets the property ID of the receiver using the given value.
func (o *Namespaces) SetID(ID string) {

	o.ID = ID
}

// GetName returns the Name of the receiver.
func (o *Namespaces) GetName() string {

	return o.Name
}

// SetName sets the property Name of the receiver using the given value.
func (o *Namespaces) SetName(name string) {

	o.Name = name
}

// GetNamespace returns the Namespace of the receiver.
func (o *Namespaces) GetNamespace() string {

	return o.Namespace
}

// SetNamespace sets the property Namespace of the receiver using the given value.
func (o *Namespaces) SetNamespace(namespace string) {

	o.Namespace = namespace
}

// ToSparse returns the sparse version of the model.
// The returned object will only contain the given fields. No field means entire field set.
func (o *Namespaces) ToSparse(fields ...string) elemental.SparseIdentifiable {

	if len(fields) == 0 {
		// nolint: goimports
		return &SparseNamespaces{
			ID:        &o.ID,
			Name:      &o.Name,
			Namespace: &o.Namespace,
		}
	}

	sp := &SparseNamespaces{}
	for _, f := range fields {
		switch f {
		case "ID":
			sp.ID = &(o.ID)
		case "name":
			sp.Name = &(o.Name)
		case "namespace":
			sp.Namespace = &(o.Namespace)
		}
	}

	return sp
}

// Patch apply the non nil value of a *SparseNamespaces to the object.
func (o *Namespaces) Patch(sparse elemental.SparseIdentifiable) {
	if !sparse.Identity().IsEqual(o.Identity()) {
		panic("cannot patch from a parse with different identity")
	}

	so := sparse.(*SparseNamespaces)
	if so.ID != nil {
		o.ID = *so.ID
	}
	if so.Name != nil {
		o.Name = *so.Name
	}
	if so.Namespace != nil {
		o.Namespace = *so.Namespace
	}
}

// DeepCopy returns a deep copy if the Namespaces.
func (o *Namespaces) DeepCopy() *Namespaces {

	if o == nil {
		return nil
	}

	out := &Namespaces{}
	o.DeepCopyInto(out)

	return out
}

// DeepCopyInto copies the receiver into the given *Namespaces.
func (o *Namespaces) DeepCopyInto(out *Namespaces) {

	target, err := copystructure.Copy(o)
	if err != nil {
		panic(fmt.Sprintf("Unable to deepcopy Namespaces: %s", err))
	}

	*out = *target.(*Namespaces)
}

// Validate valides the current information stored into the structure.
func (o *Namespaces) Validate() error {

	errors := elemental.Errors{}
	requiredErrors := elemental.Errors{}

	if err := elemental.ValidateRequiredString("name", o.Name); err != nil {
		requiredErrors = requiredErrors.Append(err)
	}

	if err := elemental.ValidatePattern("name", o.Name, `^[a-zA-Z0-9_/]+$`, `must only contain alpha numerical characters, '-' or '_'`, true); err != nil {
		errors = errors.Append(err)
	}

	if len(requiredErrors) > 0 {
		return requiredErrors
	}

	if len(errors) > 0 {
		return errors
	}

	return nil
}

// SpecificationForAttribute returns the AttributeSpecification for the given attribute name key.
func (*Namespaces) SpecificationForAttribute(name string) elemental.AttributeSpecification {

	if v, ok := NamespacesAttributesMap[name]; ok {
		return v
	}

	// We could not find it, so let's check on the lower case indexed spec map
	return NamespacesLowerCaseAttributesMap[name]
}

// AttributeSpecifications returns the full attribute specifications map.
func (*Namespaces) AttributeSpecifications() map[string]elemental.AttributeSpecification {

	return NamespacesAttributesMap
}

// ValueForAttribute returns the value for the given attribute.
// This is a very advanced function that you should not need but in some
// very specific use cases.
func (o *Namespaces) ValueForAttribute(name string) interface{} {

	switch name {
	case "ID":
		return o.ID
	case "name":
		return o.Name
	case "namespace":
		return o.Namespace
	}

	return nil
}

// NamespacesAttributesMap represents the map of attribute for Namespaces.
var NamespacesAttributesMap = map[string]elemental.AttributeSpecification{
	"ID": {
		AllowedChoices: []string{},
		Autogenerated:  true,
		BSONFieldName:  "_id",
		ConvertedName:  "ID",
		Description:    `ID is the identifier of the object.`,
		Exposed:        true,
		Getter:         true,
		Identifier:     true,
		Name:           "ID",
		Orderable:      true,
		ReadOnly:       true,
		Setter:         true,
		Stored:         true,
		Type:           "string",
	},
	"Name": {
		AllowedChars:   `^[a-zA-Z0-9_/]+$`,
		AllowedChoices: []string{},
		BSONFieldName:  "name",
		ConvertedName:  "Name",
		CreationOnly:   true,
		Description: `The name of the namespace. When you create a namespace, only put its bare name,
not its full path.`,
		Exposed:  true,
		Getter:   true,
		Name:     "name",
		Required: true,
		Setter:   true,
		Stored:   true,
		Type:     "string",
	},
	"Namespace": {
		AllowedChoices: []string{},
		Autogenerated:  true,
		BSONFieldName:  "namespace",
		ConvertedName:  "Namespace",
		Description:    `The namespace of the object.`,
		Exposed:        true,
		Getter:         true,
		Name:           "namespace",
		Orderable:      true,
		ReadOnly:       true,
		Setter:         true,
		Stored:         true,
		Type:           "string",
	},
}

// NamespacesLowerCaseAttributesMap represents the map of attribute for Namespaces.
var NamespacesLowerCaseAttributesMap = map[string]elemental.AttributeSpecification{
	"id": {
		AllowedChoices: []string{},
		Autogenerated:  true,
		BSONFieldName:  "_id",
		ConvertedName:  "ID",
		Description:    `ID is the identifier of the object.`,
		Exposed:        true,
		Getter:         true,
		Identifier:     true,
		Name:           "ID",
		Orderable:      true,
		ReadOnly:       true,
		Setter:         true,
		Stored:         true,
		Type:           "string",
	},
	"name": {
		AllowedChars:   `^[a-zA-Z0-9_/]+$`,
		AllowedChoices: []string{},
		BSONFieldName:  "name",
		ConvertedName:  "Name",
		CreationOnly:   true,
		Description: `The name of the namespace. When you create a namespace, only put its bare name,
not its full path.`,
		Exposed:  true,
		Getter:   true,
		Name:     "name",
		Required: true,
		Setter:   true,
		Stored:   true,
		Type:     "string",
	},
	"namespace": {
		AllowedChoices: []string{},
		Autogenerated:  true,
		BSONFieldName:  "namespace",
		ConvertedName:  "Namespace",
		Description:    `The namespace of the object.`,
		Exposed:        true,
		Getter:         true,
		Name:           "namespace",
		Orderable:      true,
		ReadOnly:       true,
		Setter:         true,
		Stored:         true,
		Type:           "string",
	},
}

// SparseNamespacesList represents a list of SparseNamespaces
type SparseNamespacesList []*SparseNamespaces

// Identity returns the identity of the objects in the list.
func (o SparseNamespacesList) Identity() elemental.Identity {

	return NamespacesIdentity
}

// Copy returns a pointer to a copy the SparseNamespacesList.
func (o SparseNamespacesList) Copy() elemental.Identifiables {

	copy := append(SparseNamespacesList{}, o...)
	return &copy
}

// Append appends the objects to the a new copy of the SparseNamespacesList.
func (o SparseNamespacesList) Append(objects ...elemental.Identifiable) elemental.Identifiables {

	out := append(SparseNamespacesList{}, o...)
	for _, obj := range objects {
		out = append(out, obj.(*SparseNamespaces))
	}

	return out
}

// List converts the object to an elemental.IdentifiablesList.
func (o SparseNamespacesList) List() elemental.IdentifiablesList {

	out := make(elemental.IdentifiablesList, len(o))
	for i := 0; i < len(o); i++ {
		out[i] = o[i]
	}

	return out
}

// DefaultOrder returns the default ordering fields of the content.
func (o SparseNamespacesList) DefaultOrder() []string {

	return []string{}
}

// ToPlain returns the SparseNamespacesList converted to NamespacesList.
func (o SparseNamespacesList) ToPlain() elemental.IdentifiablesList {

	out := make(elemental.IdentifiablesList, len(o))
	for i := 0; i < len(o); i++ {
		out[i] = o[i].ToPlain()
	}

	return out
}

// Version returns the version of the content.
func (o SparseNamespacesList) Version() int {

	return 1
}

// SparseNamespaces represents the sparse version of a namespace.
type SparseNamespaces struct {
	// ID is the identifier of the object.
	ID *string `json:"ID,omitempty" msgpack:"ID,omitempty" bson:"-" mapstructure:"ID,omitempty"`

	// The name of the namespace. When you create a namespace, only put its bare name,
	// not its full path.
	Name *string `json:"name,omitempty" msgpack:"name,omitempty" bson:"name,omitempty" mapstructure:"name,omitempty"`

	// The namespace of the object.
	Namespace *string `json:"namespace,omitempty" msgpack:"namespace,omitempty" bson:"namespace,omitempty" mapstructure:"namespace,omitempty"`

	ModelVersion int `json:"-" msgpack:"-" bson:"_modelversion"`
}

// NewSparseNamespaces returns a new  SparseNamespaces.
func NewSparseNamespaces() *SparseNamespaces {
	return &SparseNamespaces{}
}

// Identity returns the Identity of the sparse object.
func (o *SparseNamespaces) Identity() elemental.Identity {

	return NamespacesIdentity
}

// Identifier returns the value of the sparse object's unique identifier.
func (o *SparseNamespaces) Identifier() string {

	if o.ID == nil {
		return ""
	}
	return *o.ID
}

// SetIdentifier sets the value of the sparse object's unique identifier.
func (o *SparseNamespaces) SetIdentifier(id string) {

	if id != "" {
		o.ID = &id
	} else {
		o.ID = nil
	}
}

// GetBSON implements the bson marshaling interface.
// This is used to transparently convert ID to MongoDBID as ObectID.
func (o *SparseNamespaces) GetBSON() (interface{}, error) {

	if o == nil {
		return nil, nil
	}

	s := &mongoAttributesSparseNamespaces{}

	if o.ID != nil {
		s.ID = bson.ObjectIdHex(*o.ID)
	}
	if o.Name != nil {
		s.Name = o.Name
	}
	if o.Namespace != nil {
		s.Namespace = o.Namespace
	}

	return s, nil
}

// SetBSON implements the bson marshaling interface.
// This is used to transparently convert ID to MongoDBID as ObectID.
func (o *SparseNamespaces) SetBSON(raw bson.Raw) error {

	if o == nil {
		return nil
	}

	s := &mongoAttributesSparseNamespaces{}
	if err := raw.Unmarshal(s); err != nil {
		return err
	}

	id := s.ID.Hex()
	o.ID = &id
	if s.Name != nil {
		o.Name = s.Name
	}
	if s.Namespace != nil {
		o.Namespace = s.Namespace
	}

	return nil
}

// Version returns the hardcoded version of the model.
func (o *SparseNamespaces) Version() int {

	return 1
}

// ToPlain returns the plain version of the sparse model.
func (o *SparseNamespaces) ToPlain() elemental.PlainIdentifiable {

	out := NewNamespaces()
	if o.ID != nil {
		out.ID = *o.ID
	}
	if o.Name != nil {
		out.Name = *o.Name
	}
	if o.Namespace != nil {
		out.Namespace = *o.Namespace
	}

	return out
}

// GetID returns the ID of the receiver.
func (o *SparseNamespaces) GetID() (out string) {

	if o.ID == nil {
		return
	}

	return *o.ID
}

// SetID sets the property ID of the receiver using the address of the given value.
func (o *SparseNamespaces) SetID(ID string) {

	o.ID = &ID
}

// GetName returns the Name of the receiver.
func (o *SparseNamespaces) GetName() (out string) {

	if o.Name == nil {
		return
	}

	return *o.Name
}

// SetName sets the property Name of the receiver using the address of the given value.
func (o *SparseNamespaces) SetName(name string) {

	o.Name = &name
}

// GetNamespace returns the Namespace of the receiver.
func (o *SparseNamespaces) GetNamespace() (out string) {

	if o.Namespace == nil {
		return
	}

	return *o.Namespace
}

// SetNamespace sets the property Namespace of the receiver using the address of the given value.
func (o *SparseNamespaces) SetNamespace(namespace string) {

	o.Namespace = &namespace
}

// DeepCopy returns a deep copy if the SparseNamespaces.
func (o *SparseNamespaces) DeepCopy() *SparseNamespaces {

	if o == nil {
		return nil
	}

	out := &SparseNamespaces{}
	o.DeepCopyInto(out)

	return out
}

// DeepCopyInto copies the receiver into the given *SparseNamespaces.
func (o *SparseNamespaces) DeepCopyInto(out *SparseNamespaces) {

	target, err := copystructure.Copy(o)
	if err != nil {
		panic(fmt.Sprintf("Unable to deepcopy SparseNamespaces: %s", err))
	}

	*out = *target.(*SparseNamespaces)
}

type mongoAttributesNamespaces struct {
	ID        bson.ObjectId `bson:"_id,omitempty"`
	Name      string        `bson:"name"`
	Namespace string        `bson:"namespace"`
}
type mongoAttributesSparseNamespaces struct {
	ID        bson.ObjectId `bson:"_id,omitempty"`
	Name      *string       `bson:"name,omitempty"`
	Namespace *string       `bson:"namespace,omitempty"`
}
