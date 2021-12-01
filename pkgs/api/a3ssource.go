package api

import (
	"fmt"

	"github.com/globalsign/mgo/bson"
	"github.com/mitchellh/copystructure"
	"go.aporeto.io/elemental"
)

// A3SSourceIdentity represents the Identity of the object.
var A3SSourceIdentity = elemental.Identity{
	Name:     "a3ssource",
	Category: "a3ssources",
	Package:  "a3s",
	Private:  false,
}

// A3SSourcesList represents a list of A3SSources
type A3SSourcesList []*A3SSource

// Identity returns the identity of the objects in the list.
func (o A3SSourcesList) Identity() elemental.Identity {

	return A3SSourceIdentity
}

// Copy returns a pointer to a copy the A3SSourcesList.
func (o A3SSourcesList) Copy() elemental.Identifiables {

	copy := append(A3SSourcesList{}, o...)
	return &copy
}

// Append appends the objects to the a new copy of the A3SSourcesList.
func (o A3SSourcesList) Append(objects ...elemental.Identifiable) elemental.Identifiables {

	out := append(A3SSourcesList{}, o...)
	for _, obj := range objects {
		out = append(out, obj.(*A3SSource))
	}

	return out
}

// List converts the object to an elemental.IdentifiablesList.
func (o A3SSourcesList) List() elemental.IdentifiablesList {

	out := make(elemental.IdentifiablesList, len(o))
	for i := 0; i < len(o); i++ {
		out[i] = o[i]
	}

	return out
}

// DefaultOrder returns the default ordering fields of the content.
func (o A3SSourcesList) DefaultOrder() []string {

	return []string{}
}

// ToSparse returns the A3SSourcesList converted to SparseA3SSourcesList.
// Objects in the list will only contain the given fields. No field means entire field set.
func (o A3SSourcesList) ToSparse(fields ...string) elemental.Identifiables {

	out := make(SparseA3SSourcesList, len(o))
	for i := 0; i < len(o); i++ {
		out[i] = o[i].ToSparse(fields...).(*SparseA3SSource)
	}

	return out
}

// Version returns the version of the content.
func (o A3SSourcesList) Version() int {

	return 1
}

// A3SSource represents the model of a a3ssource
type A3SSource struct {
	// The Certificate authority to use to validate the authenticity of the A3S
	// server. If left empty, the system trust stroe will be used.
	CA string `json:"CA" msgpack:"CA" bson:"ca" mapstructure:"CA,omitempty"`

	// ID is the identifier of the object.
	ID string `json:"ID" msgpack:"ID" bson:"-" mapstructure:"ID,omitempty"`

	// The audience that must be present in the remote a3s token.
	Audience string `json:"audience" msgpack:"audience" bson:"audience" mapstructure:"audience,omitempty"`

	// The description of the object.
	Description string `json:"description" msgpack:"description" bson:"description" mapstructure:"description,omitempty"`

	// Endpoint of the remote a3s server, in case it is different from the issuer. If
	// left empty, the issuer value will be used.
	Endpoint string `json:"endpoint" msgpack:"endpoint" bson:"endpoint" mapstructure:"endpoint,omitempty"`

	// The issuer that represents the remote a3s server.
	Issuer string `json:"issuer" msgpack:"issuer" bson:"issuer" mapstructure:"issuer,omitempty"`

	// Contains optional information about a remote service that can be used to modify
	// the claims that are about to be delivered using this authentication source.
	Modifier *IdentityModifier `json:"modifier,omitempty" msgpack:"modifier,omitempty" bson:"modifier,omitempty" mapstructure:"modifier,omitempty"`

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

// NewA3SSource returns a new *A3SSource
func NewA3SSource() *A3SSource {

	return &A3SSource{
		ModelVersion: 1,
	}
}

// Identity returns the Identity of the object.
func (o *A3SSource) Identity() elemental.Identity {

	return A3SSourceIdentity
}

// Identifier returns the value of the object's unique identifier.
func (o *A3SSource) Identifier() string {

	return o.ID
}

// SetIdentifier sets the value of the object's unique identifier.
func (o *A3SSource) SetIdentifier(id string) {

	o.ID = id
}

// GetBSON implements the bson marshaling interface.
// This is used to transparently convert ID to MongoDBID as ObectID.
func (o *A3SSource) GetBSON() (interface{}, error) {

	if o == nil {
		return nil, nil
	}

	s := &mongoAttributesA3SSource{}

	s.CA = o.CA
	if o.ID != "" {
		s.ID = bson.ObjectIdHex(o.ID)
	}
	s.Audience = o.Audience
	s.Description = o.Description
	s.Endpoint = o.Endpoint
	s.Issuer = o.Issuer
	s.Modifier = o.Modifier
	s.Name = o.Name
	s.Namespace = o.Namespace
	s.ZHash = o.ZHash
	s.Zone = o.Zone

	return s, nil
}

// SetBSON implements the bson marshaling interface.
// This is used to transparently convert ID to MongoDBID as ObectID.
func (o *A3SSource) SetBSON(raw bson.Raw) error {

	if o == nil {
		return nil
	}

	s := &mongoAttributesA3SSource{}
	if err := raw.Unmarshal(s); err != nil {
		return err
	}

	o.CA = s.CA
	o.ID = s.ID.Hex()
	o.Audience = s.Audience
	o.Description = s.Description
	o.Endpoint = s.Endpoint
	o.Issuer = s.Issuer
	o.Modifier = s.Modifier
	o.Name = s.Name
	o.Namespace = s.Namespace
	o.ZHash = s.ZHash
	o.Zone = s.Zone

	return nil
}

// Version returns the hardcoded version of the model.
func (o *A3SSource) Version() int {

	return 1
}

// BleveType implements the bleve.Classifier Interface.
func (o *A3SSource) BleveType() string {

	return "a3ssource"
}

// DefaultOrder returns the list of default ordering fields.
func (o *A3SSource) DefaultOrder() []string {

	return []string{}
}

// Doc returns the documentation for the object
func (o *A3SSource) Doc() string {

	return `A source allowing to trust a remote instance of A3S.`
}

func (o *A3SSource) String() string {

	return fmt.Sprintf("<%s:%s>", o.Identity().Name, o.Identifier())
}

// GetID returns the ID of the receiver.
func (o *A3SSource) GetID() string {

	return o.ID
}

// SetID sets the property ID of the receiver using the given value.
func (o *A3SSource) SetID(ID string) {

	o.ID = ID
}

// GetNamespace returns the Namespace of the receiver.
func (o *A3SSource) GetNamespace() string {

	return o.Namespace
}

// SetNamespace sets the property Namespace of the receiver using the given value.
func (o *A3SSource) SetNamespace(namespace string) {

	o.Namespace = namespace
}

// GetZHash returns the ZHash of the receiver.
func (o *A3SSource) GetZHash() int {

	return o.ZHash
}

// SetZHash sets the property ZHash of the receiver using the given value.
func (o *A3SSource) SetZHash(zHash int) {

	o.ZHash = zHash
}

// GetZone returns the Zone of the receiver.
func (o *A3SSource) GetZone() int {

	return o.Zone
}

// SetZone sets the property Zone of the receiver using the given value.
func (o *A3SSource) SetZone(zone int) {

	o.Zone = zone
}

// ToSparse returns the sparse version of the model.
// The returned object will only contain the given fields. No field means entire field set.
func (o *A3SSource) ToSparse(fields ...string) elemental.SparseIdentifiable {

	if len(fields) == 0 {
		// nolint: goimports
		return &SparseA3SSource{
			CA:          &o.CA,
			ID:          &o.ID,
			Audience:    &o.Audience,
			Description: &o.Description,
			Endpoint:    &o.Endpoint,
			Issuer:      &o.Issuer,
			Modifier:    o.Modifier,
			Name:        &o.Name,
			Namespace:   &o.Namespace,
			ZHash:       &o.ZHash,
			Zone:        &o.Zone,
		}
	}

	sp := &SparseA3SSource{}
	for _, f := range fields {
		switch f {
		case "CA":
			sp.CA = &(o.CA)
		case "ID":
			sp.ID = &(o.ID)
		case "audience":
			sp.Audience = &(o.Audience)
		case "description":
			sp.Description = &(o.Description)
		case "endpoint":
			sp.Endpoint = &(o.Endpoint)
		case "issuer":
			sp.Issuer = &(o.Issuer)
		case "modifier":
			sp.Modifier = o.Modifier
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

// Patch apply the non nil value of a *SparseA3SSource to the object.
func (o *A3SSource) Patch(sparse elemental.SparseIdentifiable) {
	if !sparse.Identity().IsEqual(o.Identity()) {
		panic("cannot patch from a parse with different identity")
	}

	so := sparse.(*SparseA3SSource)
	if so.CA != nil {
		o.CA = *so.CA
	}
	if so.ID != nil {
		o.ID = *so.ID
	}
	if so.Audience != nil {
		o.Audience = *so.Audience
	}
	if so.Description != nil {
		o.Description = *so.Description
	}
	if so.Endpoint != nil {
		o.Endpoint = *so.Endpoint
	}
	if so.Issuer != nil {
		o.Issuer = *so.Issuer
	}
	if so.Modifier != nil {
		o.Modifier = so.Modifier
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

// DeepCopy returns a deep copy if the A3SSource.
func (o *A3SSource) DeepCopy() *A3SSource {

	if o == nil {
		return nil
	}

	out := &A3SSource{}
	o.DeepCopyInto(out)

	return out
}

// DeepCopyInto copies the receiver into the given *A3SSource.
func (o *A3SSource) DeepCopyInto(out *A3SSource) {

	target, err := copystructure.Copy(o)
	if err != nil {
		panic(fmt.Sprintf("Unable to deepcopy A3SSource: %s", err))
	}

	*out = *target.(*A3SSource)
}

// Validate valides the current information stored into the structure.
func (o *A3SSource) Validate() error {

	errors := elemental.Errors{}
	requiredErrors := elemental.Errors{}

	if err := ValidatePEM("CA", o.CA); err != nil {
		errors = errors.Append(err)
	}

	if err := elemental.ValidateRequiredString("issuer", o.Issuer); err != nil {
		requiredErrors = requiredErrors.Append(err)
	}

	if o.Modifier != nil {
		elemental.ResetDefaultForZeroValues(o.Modifier)
		if err := o.Modifier.Validate(); err != nil {
			errors = errors.Append(err)
		}
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
func (*A3SSource) SpecificationForAttribute(name string) elemental.AttributeSpecification {

	if v, ok := A3SSourceAttributesMap[name]; ok {
		return v
	}

	// We could not find it, so let's check on the lower case indexed spec map
	return A3SSourceLowerCaseAttributesMap[name]
}

// AttributeSpecifications returns the full attribute specifications map.
func (*A3SSource) AttributeSpecifications() map[string]elemental.AttributeSpecification {

	return A3SSourceAttributesMap
}

// ValueForAttribute returns the value for the given attribute.
// This is a very advanced function that you should not need but in some
// very specific use cases.
func (o *A3SSource) ValueForAttribute(name string) interface{} {

	switch name {
	case "CA":
		return o.CA
	case "ID":
		return o.ID
	case "audience":
		return o.Audience
	case "description":
		return o.Description
	case "endpoint":
		return o.Endpoint
	case "issuer":
		return o.Issuer
	case "modifier":
		return o.Modifier
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

// A3SSourceAttributesMap represents the map of attribute for A3SSource.
var A3SSourceAttributesMap = map[string]elemental.AttributeSpecification{
	"CA": {
		AllowedChoices: []string{},
		BSONFieldName:  "ca",
		ConvertedName:  "CA",
		Description: `The Certificate authority to use to validate the authenticity of the A3S
server. If left empty, the system trust stroe will be used.`,
		Exposed: true,
		Name:    "CA",
		Stored:  true,
		Type:    "string",
	},
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
	"Audience": {
		AllowedChoices: []string{},
		BSONFieldName:  "audience",
		ConvertedName:  "Audience",
		Description:    `The audience that must be present in the remote a3s token.`,
		Exposed:        true,
		Name:           "audience",
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
	"Endpoint": {
		AllowedChoices: []string{},
		BSONFieldName:  "endpoint",
		ConvertedName:  "Endpoint",
		Description: `Endpoint of the remote a3s server, in case it is different from the issuer. If
left empty, the issuer value will be used.`,
		Exposed: true,
		Name:    "endpoint",
		Stored:  true,
		Type:    "string",
	},
	"Issuer": {
		AllowedChoices: []string{},
		BSONFieldName:  "issuer",
		ConvertedName:  "Issuer",
		Description:    `The issuer that represents the remote a3s server.`,
		Exposed:        true,
		Name:           "issuer",
		Required:       true,
		Stored:         true,
		Type:           "string",
	},
	"Modifier": {
		AllowedChoices: []string{},
		BSONFieldName:  "modifier",
		ConvertedName:  "Modifier",
		Description: `Contains optional information about a remote service that can be used to modify
the claims that are about to be delivered using this authentication source.`,
		Exposed: true,
		Name:    "modifier",
		Stored:  true,
		SubType: "identitymodifier",
		Type:    "ref",
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

// A3SSourceLowerCaseAttributesMap represents the map of attribute for A3SSource.
var A3SSourceLowerCaseAttributesMap = map[string]elemental.AttributeSpecification{
	"ca": {
		AllowedChoices: []string{},
		BSONFieldName:  "ca",
		ConvertedName:  "CA",
		Description: `The Certificate authority to use to validate the authenticity of the A3S
server. If left empty, the system trust stroe will be used.`,
		Exposed: true,
		Name:    "CA",
		Stored:  true,
		Type:    "string",
	},
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
	"audience": {
		AllowedChoices: []string{},
		BSONFieldName:  "audience",
		ConvertedName:  "Audience",
		Description:    `The audience that must be present in the remote a3s token.`,
		Exposed:        true,
		Name:           "audience",
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
	"endpoint": {
		AllowedChoices: []string{},
		BSONFieldName:  "endpoint",
		ConvertedName:  "Endpoint",
		Description: `Endpoint of the remote a3s server, in case it is different from the issuer. If
left empty, the issuer value will be used.`,
		Exposed: true,
		Name:    "endpoint",
		Stored:  true,
		Type:    "string",
	},
	"issuer": {
		AllowedChoices: []string{},
		BSONFieldName:  "issuer",
		ConvertedName:  "Issuer",
		Description:    `The issuer that represents the remote a3s server.`,
		Exposed:        true,
		Name:           "issuer",
		Required:       true,
		Stored:         true,
		Type:           "string",
	},
	"modifier": {
		AllowedChoices: []string{},
		BSONFieldName:  "modifier",
		ConvertedName:  "Modifier",
		Description: `Contains optional information about a remote service that can be used to modify
the claims that are about to be delivered using this authentication source.`,
		Exposed: true,
		Name:    "modifier",
		Stored:  true,
		SubType: "identitymodifier",
		Type:    "ref",
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

// SparseA3SSourcesList represents a list of SparseA3SSources
type SparseA3SSourcesList []*SparseA3SSource

// Identity returns the identity of the objects in the list.
func (o SparseA3SSourcesList) Identity() elemental.Identity {

	return A3SSourceIdentity
}

// Copy returns a pointer to a copy the SparseA3SSourcesList.
func (o SparseA3SSourcesList) Copy() elemental.Identifiables {

	copy := append(SparseA3SSourcesList{}, o...)
	return &copy
}

// Append appends the objects to the a new copy of the SparseA3SSourcesList.
func (o SparseA3SSourcesList) Append(objects ...elemental.Identifiable) elemental.Identifiables {

	out := append(SparseA3SSourcesList{}, o...)
	for _, obj := range objects {
		out = append(out, obj.(*SparseA3SSource))
	}

	return out
}

// List converts the object to an elemental.IdentifiablesList.
func (o SparseA3SSourcesList) List() elemental.IdentifiablesList {

	out := make(elemental.IdentifiablesList, len(o))
	for i := 0; i < len(o); i++ {
		out[i] = o[i]
	}

	return out
}

// DefaultOrder returns the default ordering fields of the content.
func (o SparseA3SSourcesList) DefaultOrder() []string {

	return []string{}
}

// ToPlain returns the SparseA3SSourcesList converted to A3SSourcesList.
func (o SparseA3SSourcesList) ToPlain() elemental.IdentifiablesList {

	out := make(elemental.IdentifiablesList, len(o))
	for i := 0; i < len(o); i++ {
		out[i] = o[i].ToPlain()
	}

	return out
}

// Version returns the version of the content.
func (o SparseA3SSourcesList) Version() int {

	return 1
}

// SparseA3SSource represents the sparse version of a a3ssource.
type SparseA3SSource struct {
	// The Certificate authority to use to validate the authenticity of the A3S
	// server. If left empty, the system trust stroe will be used.
	CA *string `json:"CA,omitempty" msgpack:"CA,omitempty" bson:"ca,omitempty" mapstructure:"CA,omitempty"`

	// ID is the identifier of the object.
	ID *string `json:"ID,omitempty" msgpack:"ID,omitempty" bson:"-" mapstructure:"ID,omitempty"`

	// The audience that must be present in the remote a3s token.
	Audience *string `json:"audience,omitempty" msgpack:"audience,omitempty" bson:"audience,omitempty" mapstructure:"audience,omitempty"`

	// The description of the object.
	Description *string `json:"description,omitempty" msgpack:"description,omitempty" bson:"description,omitempty" mapstructure:"description,omitempty"`

	// Endpoint of the remote a3s server, in case it is different from the issuer. If
	// left empty, the issuer value will be used.
	Endpoint *string `json:"endpoint,omitempty" msgpack:"endpoint,omitempty" bson:"endpoint,omitempty" mapstructure:"endpoint,omitempty"`

	// The issuer that represents the remote a3s server.
	Issuer *string `json:"issuer,omitempty" msgpack:"issuer,omitempty" bson:"issuer,omitempty" mapstructure:"issuer,omitempty"`

	// Contains optional information about a remote service that can be used to modify
	// the claims that are about to be delivered using this authentication source.
	Modifier *IdentityModifier `json:"modifier,omitempty" msgpack:"modifier,omitempty" bson:"modifier,omitempty" mapstructure:"modifier,omitempty"`

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

// NewSparseA3SSource returns a new  SparseA3SSource.
func NewSparseA3SSource() *SparseA3SSource {
	return &SparseA3SSource{}
}

// Identity returns the Identity of the sparse object.
func (o *SparseA3SSource) Identity() elemental.Identity {

	return A3SSourceIdentity
}

// Identifier returns the value of the sparse object's unique identifier.
func (o *SparseA3SSource) Identifier() string {

	if o.ID == nil {
		return ""
	}
	return *o.ID
}

// SetIdentifier sets the value of the sparse object's unique identifier.
func (o *SparseA3SSource) SetIdentifier(id string) {

	if id != "" {
		o.ID = &id
	} else {
		o.ID = nil
	}
}

// GetBSON implements the bson marshaling interface.
// This is used to transparently convert ID to MongoDBID as ObectID.
func (o *SparseA3SSource) GetBSON() (interface{}, error) {

	if o == nil {
		return nil, nil
	}

	s := &mongoAttributesSparseA3SSource{}

	if o.CA != nil {
		s.CA = o.CA
	}
	if o.ID != nil {
		s.ID = bson.ObjectIdHex(*o.ID)
	}
	if o.Audience != nil {
		s.Audience = o.Audience
	}
	if o.Description != nil {
		s.Description = o.Description
	}
	if o.Endpoint != nil {
		s.Endpoint = o.Endpoint
	}
	if o.Issuer != nil {
		s.Issuer = o.Issuer
	}
	if o.Modifier != nil {
		s.Modifier = o.Modifier
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
func (o *SparseA3SSource) SetBSON(raw bson.Raw) error {

	if o == nil {
		return nil
	}

	s := &mongoAttributesSparseA3SSource{}
	if err := raw.Unmarshal(s); err != nil {
		return err
	}

	if s.CA != nil {
		o.CA = s.CA
	}
	id := s.ID.Hex()
	o.ID = &id
	if s.Audience != nil {
		o.Audience = s.Audience
	}
	if s.Description != nil {
		o.Description = s.Description
	}
	if s.Endpoint != nil {
		o.Endpoint = s.Endpoint
	}
	if s.Issuer != nil {
		o.Issuer = s.Issuer
	}
	if s.Modifier != nil {
		o.Modifier = s.Modifier
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
func (o *SparseA3SSource) Version() int {

	return 1
}

// ToPlain returns the plain version of the sparse model.
func (o *SparseA3SSource) ToPlain() elemental.PlainIdentifiable {

	out := NewA3SSource()
	if o.CA != nil {
		out.CA = *o.CA
	}
	if o.ID != nil {
		out.ID = *o.ID
	}
	if o.Audience != nil {
		out.Audience = *o.Audience
	}
	if o.Description != nil {
		out.Description = *o.Description
	}
	if o.Endpoint != nil {
		out.Endpoint = *o.Endpoint
	}
	if o.Issuer != nil {
		out.Issuer = *o.Issuer
	}
	if o.Modifier != nil {
		out.Modifier = o.Modifier
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
func (o *SparseA3SSource) GetID() (out string) {

	if o.ID == nil {
		return
	}

	return *o.ID
}

// SetID sets the property ID of the receiver using the address of the given value.
func (o *SparseA3SSource) SetID(ID string) {

	o.ID = &ID
}

// GetNamespace returns the Namespace of the receiver.
func (o *SparseA3SSource) GetNamespace() (out string) {

	if o.Namespace == nil {
		return
	}

	return *o.Namespace
}

// SetNamespace sets the property Namespace of the receiver using the address of the given value.
func (o *SparseA3SSource) SetNamespace(namespace string) {

	o.Namespace = &namespace
}

// GetZHash returns the ZHash of the receiver.
func (o *SparseA3SSource) GetZHash() (out int) {

	if o.ZHash == nil {
		return
	}

	return *o.ZHash
}

// SetZHash sets the property ZHash of the receiver using the address of the given value.
func (o *SparseA3SSource) SetZHash(zHash int) {

	o.ZHash = &zHash
}

// GetZone returns the Zone of the receiver.
func (o *SparseA3SSource) GetZone() (out int) {

	if o.Zone == nil {
		return
	}

	return *o.Zone
}

// SetZone sets the property Zone of the receiver using the address of the given value.
func (o *SparseA3SSource) SetZone(zone int) {

	o.Zone = &zone
}

// DeepCopy returns a deep copy if the SparseA3SSource.
func (o *SparseA3SSource) DeepCopy() *SparseA3SSource {

	if o == nil {
		return nil
	}

	out := &SparseA3SSource{}
	o.DeepCopyInto(out)

	return out
}

// DeepCopyInto copies the receiver into the given *SparseA3SSource.
func (o *SparseA3SSource) DeepCopyInto(out *SparseA3SSource) {

	target, err := copystructure.Copy(o)
	if err != nil {
		panic(fmt.Sprintf("Unable to deepcopy SparseA3SSource: %s", err))
	}

	*out = *target.(*SparseA3SSource)
}

type mongoAttributesA3SSource struct {
	CA          string            `bson:"ca"`
	ID          bson.ObjectId     `bson:"_id,omitempty"`
	Audience    string            `bson:"audience"`
	Description string            `bson:"description"`
	Endpoint    string            `bson:"endpoint"`
	Issuer      string            `bson:"issuer"`
	Modifier    *IdentityModifier `bson:"modifier,omitempty"`
	Name        string            `bson:"name"`
	Namespace   string            `bson:"namespace"`
	ZHash       int               `bson:"zhash"`
	Zone        int               `bson:"zone"`
}
type mongoAttributesSparseA3SSource struct {
	CA          *string           `bson:"ca,omitempty"`
	ID          bson.ObjectId     `bson:"_id,omitempty"`
	Audience    *string           `bson:"audience,omitempty"`
	Description *string           `bson:"description,omitempty"`
	Endpoint    *string           `bson:"endpoint,omitempty"`
	Issuer      *string           `bson:"issuer,omitempty"`
	Modifier    *IdentityModifier `bson:"modifier,omitempty"`
	Name        *string           `bson:"name,omitempty"`
	Namespace   *string           `bson:"namespace,omitempty"`
	ZHash       *int              `bson:"zhash,omitempty"`
	Zone        *int              `bson:"zone,omitempty"`
}
