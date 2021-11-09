package api

import (
	"fmt"

	"github.com/globalsign/mgo/bson"
	"github.com/mitchellh/copystructure"
	"go.aporeto.io/elemental"
)

// MTLSSourceIdentity represents the Identity of the object.
var MTLSSourceIdentity = elemental.Identity{
	Name:     "mtlssource",
	Category: "mtlssources",
	Package:  "a3s",
	Private:  false,
}

// MTLSSourcesList represents a list of MTLSSources
type MTLSSourcesList []*MTLSSource

// Identity returns the identity of the objects in the list.
func (o MTLSSourcesList) Identity() elemental.Identity {

	return MTLSSourceIdentity
}

// Copy returns a pointer to a copy the MTLSSourcesList.
func (o MTLSSourcesList) Copy() elemental.Identifiables {

	copy := append(MTLSSourcesList{}, o...)
	return &copy
}

// Append appends the objects to the a new copy of the MTLSSourcesList.
func (o MTLSSourcesList) Append(objects ...elemental.Identifiable) elemental.Identifiables {

	out := append(MTLSSourcesList{}, o...)
	for _, obj := range objects {
		out = append(out, obj.(*MTLSSource))
	}

	return out
}

// List converts the object to an elemental.IdentifiablesList.
func (o MTLSSourcesList) List() elemental.IdentifiablesList {

	out := make(elemental.IdentifiablesList, len(o))
	for i := 0; i < len(o); i++ {
		out[i] = o[i]
	}

	return out
}

// DefaultOrder returns the default ordering fields of the content.
func (o MTLSSourcesList) DefaultOrder() []string {

	return []string{}
}

// ToSparse returns the MTLSSourcesList converted to SparseMTLSSourcesList.
// Objects in the list will only contain the given fields. No field means entire field set.
func (o MTLSSourcesList) ToSparse(fields ...string) elemental.Identifiables {

	out := make(SparseMTLSSourcesList, len(o))
	for i := 0; i < len(o); i++ {
		out[i] = o[i].ToSparse(fields...).(*SparseMTLSSource)
	}

	return out
}

// Version returns the version of the content.
func (o MTLSSourcesList) Version() int {

	return 1
}

// MTLSSource represents the model of a mtlssource
type MTLSSource struct {
	// ID is the identifier of the object.
	ID string `json:"ID" msgpack:"ID" bson:"-" mapstructure:"ID,omitempty"`

	// The Certificate authority to use to validate user certificates in PEM format.
	CertificateAuthority string `json:"certificateAuthority" msgpack:"certificateAuthority" bson:"certificateauthority" mapstructure:"certificateAuthority,omitempty"`

	// The description of the object.
	Description string `json:"description" msgpack:"description" bson:"description" mapstructure:"description,omitempty"`

	// The name of the source.
	Name string `json:"name" msgpack:"name" bson:"name" mapstructure:"name,omitempty"`

	// The namespace of the object.
	Namespace string `json:"namespace" msgpack:"namespace" bson:"namespace" mapstructure:"namespace,omitempty"`

	// Hash of the object used to shard the data.
	ZHash int `json:"-" msgpack:"-" bson:"zhash" mapstructure:"-,omitempty"`

	// Sharding zone.
	Zone int `json:"-" msgpack:"-" bson:"zone" mapstructure:"-,omitempty"`

	ModelVersion int `json:"-" msgpack:"-" bson:"_modelversion"`
}

// NewMTLSSource returns a new *MTLSSource
func NewMTLSSource() *MTLSSource {

	return &MTLSSource{
		ModelVersion: 1,
	}
}

// Identity returns the Identity of the object.
func (o *MTLSSource) Identity() elemental.Identity {

	return MTLSSourceIdentity
}

// Identifier returns the value of the object's unique identifier.
func (o *MTLSSource) Identifier() string {

	return o.ID
}

// SetIdentifier sets the value of the object's unique identifier.
func (o *MTLSSource) SetIdentifier(id string) {

	o.ID = id
}

// GetBSON implements the bson marshaling interface.
// This is used to transparently convert ID to MongoDBID as ObectID.
func (o *MTLSSource) GetBSON() (interface{}, error) {

	if o == nil {
		return nil, nil
	}

	s := &mongoAttributesMTLSSource{}

	if o.ID != "" {
		s.ID = bson.ObjectIdHex(o.ID)
	}
	s.CertificateAuthority = o.CertificateAuthority
	s.Description = o.Description
	s.Name = o.Name
	s.Namespace = o.Namespace
	s.ZHash = o.ZHash
	s.Zone = o.Zone

	return s, nil
}

// SetBSON implements the bson marshaling interface.
// This is used to transparently convert ID to MongoDBID as ObectID.
func (o *MTLSSource) SetBSON(raw bson.Raw) error {

	if o == nil {
		return nil
	}

	s := &mongoAttributesMTLSSource{}
	if err := raw.Unmarshal(s); err != nil {
		return err
	}

	o.ID = s.ID.Hex()
	o.CertificateAuthority = s.CertificateAuthority
	o.Description = s.Description
	o.Name = s.Name
	o.Namespace = s.Namespace
	o.ZHash = s.ZHash
	o.Zone = s.Zone

	return nil
}

// Version returns the hardcoded version of the model.
func (o *MTLSSource) Version() int {

	return 1
}

// BleveType implements the bleve.Classifier Interface.
func (o *MTLSSource) BleveType() string {

	return "mtlssource"
}

// DefaultOrder returns the list of default ordering fields.
func (o *MTLSSource) DefaultOrder() []string {

	return []string{}
}

// Doc returns the documentation for the object
func (o *MTLSSource) Doc() string {

	return `An MTLS Auth source can be used to issue tokens based on user certificates.`
}

func (o *MTLSSource) String() string {

	return fmt.Sprintf("<%s:%s>", o.Identity().Name, o.Identifier())
}

// GetID returns the ID of the receiver.
func (o *MTLSSource) GetID() string {

	return o.ID
}

// SetID sets the property ID of the receiver using the given value.
func (o *MTLSSource) SetID(ID string) {

	o.ID = ID
}

// GetNamespace returns the Namespace of the receiver.
func (o *MTLSSource) GetNamespace() string {

	return o.Namespace
}

// SetNamespace sets the property Namespace of the receiver using the given value.
func (o *MTLSSource) SetNamespace(namespace string) {

	o.Namespace = namespace
}

// GetZHash returns the ZHash of the receiver.
func (o *MTLSSource) GetZHash() int {

	return o.ZHash
}

// SetZHash sets the property ZHash of the receiver using the given value.
func (o *MTLSSource) SetZHash(zHash int) {

	o.ZHash = zHash
}

// GetZone returns the Zone of the receiver.
func (o *MTLSSource) GetZone() int {

	return o.Zone
}

// SetZone sets the property Zone of the receiver using the given value.
func (o *MTLSSource) SetZone(zone int) {

	o.Zone = zone
}

// ToSparse returns the sparse version of the model.
// The returned object will only contain the given fields. No field means entire field set.
func (o *MTLSSource) ToSparse(fields ...string) elemental.SparseIdentifiable {

	if len(fields) == 0 {
		// nolint: goimports
		return &SparseMTLSSource{
			ID:                   &o.ID,
			CertificateAuthority: &o.CertificateAuthority,
			Description:          &o.Description,
			Name:                 &o.Name,
			Namespace:            &o.Namespace,
			ZHash:                &o.ZHash,
			Zone:                 &o.Zone,
		}
	}

	sp := &SparseMTLSSource{}
	for _, f := range fields {
		switch f {
		case "ID":
			sp.ID = &(o.ID)
		case "certificateAuthority":
			sp.CertificateAuthority = &(o.CertificateAuthority)
		case "description":
			sp.Description = &(o.Description)
		case "name":
			sp.Name = &(o.Name)
		case "namespace":
			sp.Namespace = &(o.Namespace)
		case "zHash":
			sp.ZHash = &(o.ZHash)
		case "zone":
			sp.Zone = &(o.Zone)
		}
	}

	return sp
}

// Patch apply the non nil value of a *SparseMTLSSource to the object.
func (o *MTLSSource) Patch(sparse elemental.SparseIdentifiable) {
	if !sparse.Identity().IsEqual(o.Identity()) {
		panic("cannot patch from a parse with different identity")
	}

	so := sparse.(*SparseMTLSSource)
	if so.ID != nil {
		o.ID = *so.ID
	}
	if so.CertificateAuthority != nil {
		o.CertificateAuthority = *so.CertificateAuthority
	}
	if so.Description != nil {
		o.Description = *so.Description
	}
	if so.Name != nil {
		o.Name = *so.Name
	}
	if so.Namespace != nil {
		o.Namespace = *so.Namespace
	}
	if so.ZHash != nil {
		o.ZHash = *so.ZHash
	}
	if so.Zone != nil {
		o.Zone = *so.Zone
	}
}

// DeepCopy returns a deep copy if the MTLSSource.
func (o *MTLSSource) DeepCopy() *MTLSSource {

	if o == nil {
		return nil
	}

	out := &MTLSSource{}
	o.DeepCopyInto(out)

	return out
}

// DeepCopyInto copies the receiver into the given *MTLSSource.
func (o *MTLSSource) DeepCopyInto(out *MTLSSource) {

	target, err := copystructure.Copy(o)
	if err != nil {
		panic(fmt.Sprintf("Unable to deepcopy MTLSSource: %s", err))
	}

	*out = *target.(*MTLSSource)
}

// Validate valides the current information stored into the structure.
func (o *MTLSSource) Validate() error {

	errors := elemental.Errors{}
	requiredErrors := elemental.Errors{}

	if err := elemental.ValidateRequiredString("certificateAuthority", o.CertificateAuthority); err != nil {
		requiredErrors = requiredErrors.Append(err)
	}

	if err := ValidatePEM("certificateAuthority", o.CertificateAuthority); err != nil {
		errors = errors.Append(err)
	}

	if err := elemental.ValidateRequiredString("name", o.Name); err != nil {
		requiredErrors = requiredErrors.Append(err)
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
func (*MTLSSource) SpecificationForAttribute(name string) elemental.AttributeSpecification {

	if v, ok := MTLSSourceAttributesMap[name]; ok {
		return v
	}

	// We could not find it, so let's check on the lower case indexed spec map
	return MTLSSourceLowerCaseAttributesMap[name]
}

// AttributeSpecifications returns the full attribute specifications map.
func (*MTLSSource) AttributeSpecifications() map[string]elemental.AttributeSpecification {

	return MTLSSourceAttributesMap
}

// ValueForAttribute returns the value for the given attribute.
// This is a very advanced function that you should not need but in some
// very specific use cases.
func (o *MTLSSource) ValueForAttribute(name string) interface{} {

	switch name {
	case "ID":
		return o.ID
	case "certificateAuthority":
		return o.CertificateAuthority
	case "description":
		return o.Description
	case "name":
		return o.Name
	case "namespace":
		return o.Namespace
	case "zHash":
		return o.ZHash
	case "zone":
		return o.Zone
	}

	return nil
}

// MTLSSourceAttributesMap represents the map of attribute for MTLSSource.
var MTLSSourceAttributesMap = map[string]elemental.AttributeSpecification{
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
	"CertificateAuthority": {
		AllowedChoices: []string{},
		BSONFieldName:  "certificateauthority",
		ConvertedName:  "CertificateAuthority",
		Description:    `The Certificate authority to use to validate user certificates in PEM format.`,
		Exposed:        true,
		Name:           "certificateAuthority",
		Required:       true,
		Stored:         true,
		Type:           "string",
	},
	"Description": {
		AllowedChoices: []string{},
		BSONFieldName:  "description",
		ConvertedName:  "Description",
		Description:    `The description of the object.`,
		Exposed:        true,
		Name:           "description",
		Stored:         true,
		Type:           "string",
	},
	"Name": {
		AllowedChoices: []string{},
		BSONFieldName:  "name",
		ConvertedName:  "Name",
		Description:    `The name of the source.`,
		Exposed:        true,
		Name:           "name",
		Required:       true,
		Stored:         true,
		Type:           "string",
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
	"ZHash": {
		AllowedChoices: []string{},
		Autogenerated:  true,
		BSONFieldName:  "zhash",
		ConvertedName:  "ZHash",
		Description:    `Hash of the object used to shard the data.`,
		Getter:         true,
		Name:           "zHash",
		ReadOnly:       true,
		Setter:         true,
		Stored:         true,
		Type:           "integer",
	},
	"Zone": {
		AllowedChoices: []string{},
		Autogenerated:  true,
		BSONFieldName:  "zone",
		ConvertedName:  "Zone",
		Description:    `Sharding zone.`,
		Getter:         true,
		Name:           "zone",
		ReadOnly:       true,
		Setter:         true,
		Stored:         true,
		Transient:      true,
		Type:           "integer",
	},
}

// MTLSSourceLowerCaseAttributesMap represents the map of attribute for MTLSSource.
var MTLSSourceLowerCaseAttributesMap = map[string]elemental.AttributeSpecification{
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
	"certificateauthority": {
		AllowedChoices: []string{},
		BSONFieldName:  "certificateauthority",
		ConvertedName:  "CertificateAuthority",
		Description:    `The Certificate authority to use to validate user certificates in PEM format.`,
		Exposed:        true,
		Name:           "certificateAuthority",
		Required:       true,
		Stored:         true,
		Type:           "string",
	},
	"description": {
		AllowedChoices: []string{},
		BSONFieldName:  "description",
		ConvertedName:  "Description",
		Description:    `The description of the object.`,
		Exposed:        true,
		Name:           "description",
		Stored:         true,
		Type:           "string",
	},
	"name": {
		AllowedChoices: []string{},
		BSONFieldName:  "name",
		ConvertedName:  "Name",
		Description:    `The name of the source.`,
		Exposed:        true,
		Name:           "name",
		Required:       true,
		Stored:         true,
		Type:           "string",
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
	"zhash": {
		AllowedChoices: []string{},
		Autogenerated:  true,
		BSONFieldName:  "zhash",
		ConvertedName:  "ZHash",
		Description:    `Hash of the object used to shard the data.`,
		Getter:         true,
		Name:           "zHash",
		ReadOnly:       true,
		Setter:         true,
		Stored:         true,
		Type:           "integer",
	},
	"zone": {
		AllowedChoices: []string{},
		Autogenerated:  true,
		BSONFieldName:  "zone",
		ConvertedName:  "Zone",
		Description:    `Sharding zone.`,
		Getter:         true,
		Name:           "zone",
		ReadOnly:       true,
		Setter:         true,
		Stored:         true,
		Transient:      true,
		Type:           "integer",
	},
}

// SparseMTLSSourcesList represents a list of SparseMTLSSources
type SparseMTLSSourcesList []*SparseMTLSSource

// Identity returns the identity of the objects in the list.
func (o SparseMTLSSourcesList) Identity() elemental.Identity {

	return MTLSSourceIdentity
}

// Copy returns a pointer to a copy the SparseMTLSSourcesList.
func (o SparseMTLSSourcesList) Copy() elemental.Identifiables {

	copy := append(SparseMTLSSourcesList{}, o...)
	return &copy
}

// Append appends the objects to the a new copy of the SparseMTLSSourcesList.
func (o SparseMTLSSourcesList) Append(objects ...elemental.Identifiable) elemental.Identifiables {

	out := append(SparseMTLSSourcesList{}, o...)
	for _, obj := range objects {
		out = append(out, obj.(*SparseMTLSSource))
	}

	return out
}

// List converts the object to an elemental.IdentifiablesList.
func (o SparseMTLSSourcesList) List() elemental.IdentifiablesList {

	out := make(elemental.IdentifiablesList, len(o))
	for i := 0; i < len(o); i++ {
		out[i] = o[i]
	}

	return out
}

// DefaultOrder returns the default ordering fields of the content.
func (o SparseMTLSSourcesList) DefaultOrder() []string {

	return []string{}
}

// ToPlain returns the SparseMTLSSourcesList converted to MTLSSourcesList.
func (o SparseMTLSSourcesList) ToPlain() elemental.IdentifiablesList {

	out := make(elemental.IdentifiablesList, len(o))
	for i := 0; i < len(o); i++ {
		out[i] = o[i].ToPlain()
	}

	return out
}

// Version returns the version of the content.
func (o SparseMTLSSourcesList) Version() int {

	return 1
}

// SparseMTLSSource represents the sparse version of a mtlssource.
type SparseMTLSSource struct {
	// ID is the identifier of the object.
	ID *string `json:"ID,omitempty" msgpack:"ID,omitempty" bson:"-" mapstructure:"ID,omitempty"`

	// The Certificate authority to use to validate user certificates in PEM format.
	CertificateAuthority *string `json:"certificateAuthority,omitempty" msgpack:"certificateAuthority,omitempty" bson:"certificateauthority,omitempty" mapstructure:"certificateAuthority,omitempty"`

	// The description of the object.
	Description *string `json:"description,omitempty" msgpack:"description,omitempty" bson:"description,omitempty" mapstructure:"description,omitempty"`

	// The name of the source.
	Name *string `json:"name,omitempty" msgpack:"name,omitempty" bson:"name,omitempty" mapstructure:"name,omitempty"`

	// The namespace of the object.
	Namespace *string `json:"namespace,omitempty" msgpack:"namespace,omitempty" bson:"namespace,omitempty" mapstructure:"namespace,omitempty"`

	// Hash of the object used to shard the data.
	ZHash *int `json:"-" msgpack:"-" bson:"zhash,omitempty" mapstructure:"-,omitempty"`

	// Sharding zone.
	Zone *int `json:"-" msgpack:"-" bson:"zone,omitempty" mapstructure:"-,omitempty"`

	ModelVersion int `json:"-" msgpack:"-" bson:"_modelversion"`
}

// NewSparseMTLSSource returns a new  SparseMTLSSource.
func NewSparseMTLSSource() *SparseMTLSSource {
	return &SparseMTLSSource{}
}

// Identity returns the Identity of the sparse object.
func (o *SparseMTLSSource) Identity() elemental.Identity {

	return MTLSSourceIdentity
}

// Identifier returns the value of the sparse object's unique identifier.
func (o *SparseMTLSSource) Identifier() string {

	if o.ID == nil {
		return ""
	}
	return *o.ID
}

// SetIdentifier sets the value of the sparse object's unique identifier.
func (o *SparseMTLSSource) SetIdentifier(id string) {

	if id != "" {
		o.ID = &id
	} else {
		o.ID = nil
	}
}

// GetBSON implements the bson marshaling interface.
// This is used to transparently convert ID to MongoDBID as ObectID.
func (o *SparseMTLSSource) GetBSON() (interface{}, error) {

	if o == nil {
		return nil, nil
	}

	s := &mongoAttributesSparseMTLSSource{}

	if o.ID != nil {
		s.ID = bson.ObjectIdHex(*o.ID)
	}
	if o.CertificateAuthority != nil {
		s.CertificateAuthority = o.CertificateAuthority
	}
	if o.Description != nil {
		s.Description = o.Description
	}
	if o.Name != nil {
		s.Name = o.Name
	}
	if o.Namespace != nil {
		s.Namespace = o.Namespace
	}
	if o.ZHash != nil {
		s.ZHash = o.ZHash
	}
	if o.Zone != nil {
		s.Zone = o.Zone
	}

	return s, nil
}

// SetBSON implements the bson marshaling interface.
// This is used to transparently convert ID to MongoDBID as ObectID.
func (o *SparseMTLSSource) SetBSON(raw bson.Raw) error {

	if o == nil {
		return nil
	}

	s := &mongoAttributesSparseMTLSSource{}
	if err := raw.Unmarshal(s); err != nil {
		return err
	}

	id := s.ID.Hex()
	o.ID = &id
	if s.CertificateAuthority != nil {
		o.CertificateAuthority = s.CertificateAuthority
	}
	if s.Description != nil {
		o.Description = s.Description
	}
	if s.Name != nil {
		o.Name = s.Name
	}
	if s.Namespace != nil {
		o.Namespace = s.Namespace
	}
	if s.ZHash != nil {
		o.ZHash = s.ZHash
	}
	if s.Zone != nil {
		o.Zone = s.Zone
	}

	return nil
}

// Version returns the hardcoded version of the model.
func (o *SparseMTLSSource) Version() int {

	return 1
}

// ToPlain returns the plain version of the sparse model.
func (o *SparseMTLSSource) ToPlain() elemental.PlainIdentifiable {

	out := NewMTLSSource()
	if o.ID != nil {
		out.ID = *o.ID
	}
	if o.CertificateAuthority != nil {
		out.CertificateAuthority = *o.CertificateAuthority
	}
	if o.Description != nil {
		out.Description = *o.Description
	}
	if o.Name != nil {
		out.Name = *o.Name
	}
	if o.Namespace != nil {
		out.Namespace = *o.Namespace
	}
	if o.ZHash != nil {
		out.ZHash = *o.ZHash
	}
	if o.Zone != nil {
		out.Zone = *o.Zone
	}

	return out
}

// GetID returns the ID of the receiver.
func (o *SparseMTLSSource) GetID() (out string) {

	if o.ID == nil {
		return
	}

	return *o.ID
}

// SetID sets the property ID of the receiver using the address of the given value.
func (o *SparseMTLSSource) SetID(ID string) {

	o.ID = &ID
}

// GetNamespace returns the Namespace of the receiver.
func (o *SparseMTLSSource) GetNamespace() (out string) {

	if o.Namespace == nil {
		return
	}

	return *o.Namespace
}

// SetNamespace sets the property Namespace of the receiver using the address of the given value.
func (o *SparseMTLSSource) SetNamespace(namespace string) {

	o.Namespace = &namespace
}

// GetZHash returns the ZHash of the receiver.
func (o *SparseMTLSSource) GetZHash() (out int) {

	if o.ZHash == nil {
		return
	}

	return *o.ZHash
}

// SetZHash sets the property ZHash of the receiver using the address of the given value.
func (o *SparseMTLSSource) SetZHash(zHash int) {

	o.ZHash = &zHash
}

// GetZone returns the Zone of the receiver.
func (o *SparseMTLSSource) GetZone() (out int) {

	if o.Zone == nil {
		return
	}

	return *o.Zone
}

// SetZone sets the property Zone of the receiver using the address of the given value.
func (o *SparseMTLSSource) SetZone(zone int) {

	o.Zone = &zone
}

// DeepCopy returns a deep copy if the SparseMTLSSource.
func (o *SparseMTLSSource) DeepCopy() *SparseMTLSSource {

	if o == nil {
		return nil
	}

	out := &SparseMTLSSource{}
	o.DeepCopyInto(out)

	return out
}

// DeepCopyInto copies the receiver into the given *SparseMTLSSource.
func (o *SparseMTLSSource) DeepCopyInto(out *SparseMTLSSource) {

	target, err := copystructure.Copy(o)
	if err != nil {
		panic(fmt.Sprintf("Unable to deepcopy SparseMTLSSource: %s", err))
	}

	*out = *target.(*SparseMTLSSource)
}

type mongoAttributesMTLSSource struct {
	ID                   bson.ObjectId `bson:"_id,omitempty"`
	CertificateAuthority string        `bson:"certificateauthority"`
	Description          string        `bson:"description"`
	Name                 string        `bson:"name"`
	Namespace            string        `bson:"namespace"`
	ZHash                int           `bson:"zhash"`
	Zone                 int           `bson:"zone"`
}
type mongoAttributesSparseMTLSSource struct {
	ID                   bson.ObjectId `bson:"_id,omitempty"`
	CertificateAuthority *string       `bson:"certificateauthority,omitempty"`
	Description          *string       `bson:"description,omitempty"`
	Name                 *string       `bson:"name,omitempty"`
	Namespace            *string       `bson:"namespace,omitempty"`
	ZHash                *int          `bson:"zhash,omitempty"`
	Zone                 *int          `bson:"zone,omitempty"`
}
