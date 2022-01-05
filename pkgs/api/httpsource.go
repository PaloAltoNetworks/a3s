package api

import (
	"fmt"

	"github.com/globalsign/mgo/bson"
	"github.com/mitchellh/copystructure"
	"go.aporeto.io/elemental"
)

// HTTPSourceIdentity represents the Identity of the object.
var HTTPSourceIdentity = elemental.Identity{
	Name:     "httpsource",
	Category: "httpsources",
	Package:  "a3s",
	Private:  false,
}

// HTTPSourcesList represents a list of HTTPSources
type HTTPSourcesList []*HTTPSource

// Identity returns the identity of the objects in the list.
func (o HTTPSourcesList) Identity() elemental.Identity {

	return HTTPSourceIdentity
}

// Copy returns a pointer to a copy the HTTPSourcesList.
func (o HTTPSourcesList) Copy() elemental.Identifiables {

	copy := append(HTTPSourcesList{}, o...)
	return &copy
}

// Append appends the objects to the a new copy of the HTTPSourcesList.
func (o HTTPSourcesList) Append(objects ...elemental.Identifiable) elemental.Identifiables {

	out := append(HTTPSourcesList{}, o...)
	for _, obj := range objects {
		out = append(out, obj.(*HTTPSource))
	}

	return out
}

// List converts the object to an elemental.IdentifiablesList.
func (o HTTPSourcesList) List() elemental.IdentifiablesList {

	out := make(elemental.IdentifiablesList, len(o))
	for i := 0; i < len(o); i++ {
		out[i] = o[i]
	}

	return out
}

// DefaultOrder returns the default ordering fields of the content.
func (o HTTPSourcesList) DefaultOrder() []string {

	return []string{}
}

// ToSparse returns the HTTPSourcesList converted to SparseHTTPSourcesList.
// Objects in the list will only contain the given fields. No field means entire field set.
func (o HTTPSourcesList) ToSparse(fields ...string) elemental.Identifiables {

	out := make(SparseHTTPSourcesList, len(o))
	for i := 0; i < len(o); i++ {
		out[i] = o[i].ToSparse(fields...).(*SparseHTTPSource)
	}

	return out
}

// Version returns the version of the content.
func (o HTTPSourcesList) Version() int {

	return 1
}

// HTTPSource represents the model of a httpsource
type HTTPSource struct {
	// The certificate authority to use to validate the remote http server.
	CA string `json:"CA" msgpack:"CA" bson:"ca" mapstructure:"CA,omitempty"`

	// ID is the identifier of the object.
	ID string `json:"ID" msgpack:"ID" bson:"-" mapstructure:"ID,omitempty"`

	// Client certificate required to call URL. A3S will refuse to send data if the
	// endpoint does not support client certificate authentication.
	Certificate string `json:"certificate" msgpack:"certificate" bson:"certificate" mapstructure:"certificate,omitempty"`

	// The description of the object.
	Description string `json:"description" msgpack:"description" bson:"description" mapstructure:"description,omitempty"`

	// URL of the remote service. This URL will receive a POST containing the
	// credentials information that must be validated. It must reply with 200 with a
	// body containing a json array that will be used as claims for the token. Any
	// other error code will be returned as a 401 error.
	Endpoint string `json:"endpoint" msgpack:"endpoint" bson:"endpoint" mapstructure:"endpoint,omitempty"`

	// Key associated to the client certificate.
	Key string `json:"key" msgpack:"key" bson:"key" mapstructure:"key,omitempty"`

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

// NewHTTPSource returns a new *HTTPSource
func NewHTTPSource() *HTTPSource {

	return &HTTPSource{
		ModelVersion: 1,
	}
}

// Identity returns the Identity of the object.
func (o *HTTPSource) Identity() elemental.Identity {

	return HTTPSourceIdentity
}

// Identifier returns the value of the object's unique identifier.
func (o *HTTPSource) Identifier() string {

	return o.ID
}

// SetIdentifier sets the value of the object's unique identifier.
func (o *HTTPSource) SetIdentifier(id string) {

	o.ID = id
}

// GetBSON implements the bson marshaling interface.
// This is used to transparently convert ID to MongoDBID as ObectID.
func (o *HTTPSource) GetBSON() (interface{}, error) {

	if o == nil {
		return nil, nil
	}

	s := &mongoAttributesHTTPSource{}

	s.CA = o.CA
	if o.ID != "" {
		s.ID = bson.ObjectIdHex(o.ID)
	}
	s.Certificate = o.Certificate
	s.Description = o.Description
	s.Endpoint = o.Endpoint
	s.Key = o.Key
	s.Modifier = o.Modifier
	s.Name = o.Name
	s.Namespace = o.Namespace
	s.ZHash = o.ZHash
	s.Zone = o.Zone

	return s, nil
}

// SetBSON implements the bson marshaling interface.
// This is used to transparently convert ID to MongoDBID as ObectID.
func (o *HTTPSource) SetBSON(raw bson.Raw) error {

	if o == nil {
		return nil
	}

	s := &mongoAttributesHTTPSource{}
	if err := raw.Unmarshal(s); err != nil {
		return err
	}

	o.CA = s.CA
	o.ID = s.ID.Hex()
	o.Certificate = s.Certificate
	o.Description = s.Description
	o.Endpoint = s.Endpoint
	o.Key = s.Key
	o.Modifier = s.Modifier
	o.Name = s.Name
	o.Namespace = s.Namespace
	o.ZHash = s.ZHash
	o.Zone = s.Zone

	return nil
}

// Version returns the hardcoded version of the model.
func (o *HTTPSource) Version() int {

	return 1
}

// BleveType implements the bleve.Classifier Interface.
func (o *HTTPSource) BleveType() string {

	return "httpsource"
}

// DefaultOrder returns the list of default ordering fields.
func (o *HTTPSource) DefaultOrder() []string {

	return []string{}
}

// Doc returns the documentation for the object
func (o *HTTPSource) Doc() string {

	return `A source that can call a remote service to validate generic credentials.`
}

func (o *HTTPSource) String() string {

	return fmt.Sprintf("<%s:%s>", o.Identity().Name, o.Identifier())
}

// GetID returns the ID of the receiver.
func (o *HTTPSource) GetID() string {

	return o.ID
}

// SetID sets the property ID of the receiver using the given value.
func (o *HTTPSource) SetID(ID string) {

	o.ID = ID
}

// GetNamespace returns the Namespace of the receiver.
func (o *HTTPSource) GetNamespace() string {

	return o.Namespace
}

// SetNamespace sets the property Namespace of the receiver using the given value.
func (o *HTTPSource) SetNamespace(namespace string) {

	o.Namespace = namespace
}

// GetZHash returns the ZHash of the receiver.
func (o *HTTPSource) GetZHash() int {

	return o.ZHash
}

// SetZHash sets the property ZHash of the receiver using the given value.
func (o *HTTPSource) SetZHash(zHash int) {

	o.ZHash = zHash
}

// GetZone returns the Zone of the receiver.
func (o *HTTPSource) GetZone() int {

	return o.Zone
}

// SetZone sets the property Zone of the receiver using the given value.
func (o *HTTPSource) SetZone(zone int) {

	o.Zone = zone
}

// ToSparse returns the sparse version of the model.
// The returned object will only contain the given fields. No field means entire field set.
func (o *HTTPSource) ToSparse(fields ...string) elemental.SparseIdentifiable {

	if len(fields) == 0 {
		// nolint: goimports
		return &SparseHTTPSource{
			CA:          &o.CA,
			ID:          &o.ID,
			Certificate: &o.Certificate,
			Description: &o.Description,
			Endpoint:    &o.Endpoint,
			Key:         &o.Key,
			Modifier:    o.Modifier,
			Name:        &o.Name,
			Namespace:   &o.Namespace,
			ZHash:       &o.ZHash,
			Zone:        &o.Zone,
		}
	}

	sp := &SparseHTTPSource{}
	for _, f := range fields {
		switch f {
		case "CA":
			sp.CA = &(o.CA)
		case "ID":
			sp.ID = &(o.ID)
		case "certificate":
			sp.Certificate = &(o.Certificate)
		case "description":
			sp.Description = &(o.Description)
		case "endpoint":
			sp.Endpoint = &(o.Endpoint)
		case "key":
			sp.Key = &(o.Key)
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

// Patch apply the non nil value of a *SparseHTTPSource to the object.
func (o *HTTPSource) Patch(sparse elemental.SparseIdentifiable) {
	if !sparse.Identity().IsEqual(o.Identity()) {
		panic("cannot patch from a parse with different identity")
	}

	so := sparse.(*SparseHTTPSource)
	if so.CA != nil {
		o.CA = *so.CA
	}
	if so.ID != nil {
		o.ID = *so.ID
	}
	if so.Certificate != nil {
		o.Certificate = *so.Certificate
	}
	if so.Description != nil {
		o.Description = *so.Description
	}
	if so.Endpoint != nil {
		o.Endpoint = *so.Endpoint
	}
	if so.Key != nil {
		o.Key = *so.Key
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

// DeepCopy returns a deep copy if the HTTPSource.
func (o *HTTPSource) DeepCopy() *HTTPSource {

	if o == nil {
		return nil
	}

	out := &HTTPSource{}
	o.DeepCopyInto(out)

	return out
}

// DeepCopyInto copies the receiver into the given *HTTPSource.
func (o *HTTPSource) DeepCopyInto(out *HTTPSource) {

	target, err := copystructure.Copy(o)
	if err != nil {
		panic(fmt.Sprintf("Unable to deepcopy HTTPSource: %s", err))
	}

	*out = *target.(*HTTPSource)
}

// Validate valides the current information stored into the structure.
func (o *HTTPSource) Validate() error {

	errors := elemental.Errors{}
	requiredErrors := elemental.Errors{}

	if err := elemental.ValidateRequiredString("CA", o.CA); err != nil {
		requiredErrors = requiredErrors.Append(err)
	}

	if err := ValidatePEM("CA", o.CA); err != nil {
		errors = errors.Append(err)
	}

	if err := elemental.ValidateRequiredString("certificate", o.Certificate); err != nil {
		requiredErrors = requiredErrors.Append(err)
	}

	if err := ValidatePEM("certificate", o.Certificate); err != nil {
		errors = errors.Append(err)
	}

	if err := elemental.ValidateRequiredString("endpoint", o.Endpoint); err != nil {
		requiredErrors = requiredErrors.Append(err)
	}

	if err := ValidateURL("endpoint", o.Endpoint); err != nil {
		errors = errors.Append(err)
	}

	if err := elemental.ValidateRequiredString("key", o.Key); err != nil {
		requiredErrors = requiredErrors.Append(err)
	}

	if err := ValidatePEM("key", o.Key); err != nil {
		errors = errors.Append(err)
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
func (*HTTPSource) SpecificationForAttribute(name string) elemental.AttributeSpecification {

	if v, ok := HTTPSourceAttributesMap[name]; ok {
		return v
	}

	// We could not find it, so let's check on the lower case indexed spec map
	return HTTPSourceLowerCaseAttributesMap[name]
}

// AttributeSpecifications returns the full attribute specifications map.
func (*HTTPSource) AttributeSpecifications() map[string]elemental.AttributeSpecification {

	return HTTPSourceAttributesMap
}

// ValueForAttribute returns the value for the given attribute.
// This is a very advanced function that you should not need but in some
// very specific use cases.
func (o *HTTPSource) ValueForAttribute(name string) interface{} {

	switch name {
	case "CA":
		return o.CA
	case "ID":
		return o.ID
	case "certificate":
		return o.Certificate
	case "description":
		return o.Description
	case "endpoint":
		return o.Endpoint
	case "key":
		return o.Key
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

// HTTPSourceAttributesMap represents the map of attribute for HTTPSource.
var HTTPSourceAttributesMap = map[string]elemental.AttributeSpecification{
	"CA": {
		AllowedChoices: []string{},
		BSONFieldName:  "ca",
		ConvertedName:  "CA",
		Description:    `The certificate authority to use to validate the remote http server.`,
		Exposed:        true,
		Name:           "CA",
		Required:       true,
		Stored:         true,
		Type:           "string",
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
	"Certificate": {
		AllowedChoices: []string{},
		BSONFieldName:  "certificate",
		ConvertedName:  "Certificate",
		Description: `Client certificate required to call URL. A3S will refuse to send data if the
endpoint does not support client certificate authentication.`,
		Exposed:  true,
		Name:     "certificate",
		Required: true,
		Stored:   true,
		Type:     "string",
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
		Description: `URL of the remote service. This URL will receive a POST containing the
credentials information that must be validated. It must reply with 200 with a
body containing a json array that will be used as claims for the token. Any
other error code will be returned as a 401 error.`,
		Exposed:  true,
		Name:     "endpoint",
		Required: true,
		Stored:   true,
		Type:     "string",
	},
	"Key": {
		AllowedChoices: []string{},
		BSONFieldName:  "key",
		ConvertedName:  "Key",
		Description:    `Key associated to the client certificate.`,
		Exposed:        true,
		Name:           "key",
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

// HTTPSourceLowerCaseAttributesMap represents the map of attribute for HTTPSource.
var HTTPSourceLowerCaseAttributesMap = map[string]elemental.AttributeSpecification{
	"ca": {
		AllowedChoices: []string{},
		BSONFieldName:  "ca",
		ConvertedName:  "CA",
		Description:    `The certificate authority to use to validate the remote http server.`,
		Exposed:        true,
		Name:           "CA",
		Required:       true,
		Stored:         true,
		Type:           "string",
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
	"certificate": {
		AllowedChoices: []string{},
		BSONFieldName:  "certificate",
		ConvertedName:  "Certificate",
		Description: `Client certificate required to call URL. A3S will refuse to send data if the
endpoint does not support client certificate authentication.`,
		Exposed:  true,
		Name:     "certificate",
		Required: true,
		Stored:   true,
		Type:     "string",
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
		Description: `URL of the remote service. This URL will receive a POST containing the
credentials information that must be validated. It must reply with 200 with a
body containing a json array that will be used as claims for the token. Any
other error code will be returned as a 401 error.`,
		Exposed:  true,
		Name:     "endpoint",
		Required: true,
		Stored:   true,
		Type:     "string",
	},
	"key": {
		AllowedChoices: []string{},
		BSONFieldName:  "key",
		ConvertedName:  "Key",
		Description:    `Key associated to the client certificate.`,
		Exposed:        true,
		Name:           "key",
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

// SparseHTTPSourcesList represents a list of SparseHTTPSources
type SparseHTTPSourcesList []*SparseHTTPSource

// Identity returns the identity of the objects in the list.
func (o SparseHTTPSourcesList) Identity() elemental.Identity {

	return HTTPSourceIdentity
}

// Copy returns a pointer to a copy the SparseHTTPSourcesList.
func (o SparseHTTPSourcesList) Copy() elemental.Identifiables {

	copy := append(SparseHTTPSourcesList{}, o...)
	return &copy
}

// Append appends the objects to the a new copy of the SparseHTTPSourcesList.
func (o SparseHTTPSourcesList) Append(objects ...elemental.Identifiable) elemental.Identifiables {

	out := append(SparseHTTPSourcesList{}, o...)
	for _, obj := range objects {
		out = append(out, obj.(*SparseHTTPSource))
	}

	return out
}

// List converts the object to an elemental.IdentifiablesList.
func (o SparseHTTPSourcesList) List() elemental.IdentifiablesList {

	out := make(elemental.IdentifiablesList, len(o))
	for i := 0; i < len(o); i++ {
		out[i] = o[i]
	}

	return out
}

// DefaultOrder returns the default ordering fields of the content.
func (o SparseHTTPSourcesList) DefaultOrder() []string {

	return []string{}
}

// ToPlain returns the SparseHTTPSourcesList converted to HTTPSourcesList.
func (o SparseHTTPSourcesList) ToPlain() elemental.IdentifiablesList {

	out := make(elemental.IdentifiablesList, len(o))
	for i := 0; i < len(o); i++ {
		out[i] = o[i].ToPlain()
	}

	return out
}

// Version returns the version of the content.
func (o SparseHTTPSourcesList) Version() int {

	return 1
}

// SparseHTTPSource represents the sparse version of a httpsource.
type SparseHTTPSource struct {
	// The certificate authority to use to validate the remote http server.
	CA *string `json:"CA,omitempty" msgpack:"CA,omitempty" bson:"ca,omitempty" mapstructure:"CA,omitempty"`

	// ID is the identifier of the object.
	ID *string `json:"ID,omitempty" msgpack:"ID,omitempty" bson:"-" mapstructure:"ID,omitempty"`

	// Client certificate required to call URL. A3S will refuse to send data if the
	// endpoint does not support client certificate authentication.
	Certificate *string `json:"certificate,omitempty" msgpack:"certificate,omitempty" bson:"certificate,omitempty" mapstructure:"certificate,omitempty"`

	// The description of the object.
	Description *string `json:"description,omitempty" msgpack:"description,omitempty" bson:"description,omitempty" mapstructure:"description,omitempty"`

	// URL of the remote service. This URL will receive a POST containing the
	// credentials information that must be validated. It must reply with 200 with a
	// body containing a json array that will be used as claims for the token. Any
	// other error code will be returned as a 401 error.
	Endpoint *string `json:"endpoint,omitempty" msgpack:"endpoint,omitempty" bson:"endpoint,omitempty" mapstructure:"endpoint,omitempty"`

	// Key associated to the client certificate.
	Key *string `json:"key,omitempty" msgpack:"key,omitempty" bson:"key,omitempty" mapstructure:"key,omitempty"`

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

// NewSparseHTTPSource returns a new  SparseHTTPSource.
func NewSparseHTTPSource() *SparseHTTPSource {
	return &SparseHTTPSource{}
}

// Identity returns the Identity of the sparse object.
func (o *SparseHTTPSource) Identity() elemental.Identity {

	return HTTPSourceIdentity
}

// Identifier returns the value of the sparse object's unique identifier.
func (o *SparseHTTPSource) Identifier() string {

	if o.ID == nil {
		return ""
	}
	return *o.ID
}

// SetIdentifier sets the value of the sparse object's unique identifier.
func (o *SparseHTTPSource) SetIdentifier(id string) {

	if id != "" {
		o.ID = &id
	} else {
		o.ID = nil
	}
}

// GetBSON implements the bson marshaling interface.
// This is used to transparently convert ID to MongoDBID as ObectID.
func (o *SparseHTTPSource) GetBSON() (interface{}, error) {

	if o == nil {
		return nil, nil
	}

	s := &mongoAttributesSparseHTTPSource{}

	if o.CA != nil {
		s.CA = o.CA
	}
	if o.ID != nil {
		s.ID = bson.ObjectIdHex(*o.ID)
	}
	if o.Certificate != nil {
		s.Certificate = o.Certificate
	}
	if o.Description != nil {
		s.Description = o.Description
	}
	if o.Endpoint != nil {
		s.Endpoint = o.Endpoint
	}
	if o.Key != nil {
		s.Key = o.Key
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
func (o *SparseHTTPSource) SetBSON(raw bson.Raw) error {

	if o == nil {
		return nil
	}

	s := &mongoAttributesSparseHTTPSource{}
	if err := raw.Unmarshal(s); err != nil {
		return err
	}

	if s.CA != nil {
		o.CA = s.CA
	}
	id := s.ID.Hex()
	o.ID = &id
	if s.Certificate != nil {
		o.Certificate = s.Certificate
	}
	if s.Description != nil {
		o.Description = s.Description
	}
	if s.Endpoint != nil {
		o.Endpoint = s.Endpoint
	}
	if s.Key != nil {
		o.Key = s.Key
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
func (o *SparseHTTPSource) Version() int {

	return 1
}

// ToPlain returns the plain version of the sparse model.
func (o *SparseHTTPSource) ToPlain() elemental.PlainIdentifiable {

	out := NewHTTPSource()
	if o.CA != nil {
		out.CA = *o.CA
	}
	if o.ID != nil {
		out.ID = *o.ID
	}
	if o.Certificate != nil {
		out.Certificate = *o.Certificate
	}
	if o.Description != nil {
		out.Description = *o.Description
	}
	if o.Endpoint != nil {
		out.Endpoint = *o.Endpoint
	}
	if o.Key != nil {
		out.Key = *o.Key
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
func (o *SparseHTTPSource) GetID() (out string) {

	if o.ID == nil {
		return
	}

	return *o.ID
}

// SetID sets the property ID of the receiver using the address of the given value.
func (o *SparseHTTPSource) SetID(ID string) {

	o.ID = &ID
}

// GetNamespace returns the Namespace of the receiver.
func (o *SparseHTTPSource) GetNamespace() (out string) {

	if o.Namespace == nil {
		return
	}

	return *o.Namespace
}

// SetNamespace sets the property Namespace of the receiver using the address of the given value.
func (o *SparseHTTPSource) SetNamespace(namespace string) {

	o.Namespace = &namespace
}

// GetZHash returns the ZHash of the receiver.
func (o *SparseHTTPSource) GetZHash() (out int) {

	if o.ZHash == nil {
		return
	}

	return *o.ZHash
}

// SetZHash sets the property ZHash of the receiver using the address of the given value.
func (o *SparseHTTPSource) SetZHash(zHash int) {

	o.ZHash = &zHash
}

// GetZone returns the Zone of the receiver.
func (o *SparseHTTPSource) GetZone() (out int) {

	if o.Zone == nil {
		return
	}

	return *o.Zone
}

// SetZone sets the property Zone of the receiver using the address of the given value.
func (o *SparseHTTPSource) SetZone(zone int) {

	o.Zone = &zone
}

// DeepCopy returns a deep copy if the SparseHTTPSource.
func (o *SparseHTTPSource) DeepCopy() *SparseHTTPSource {

	if o == nil {
		return nil
	}

	out := &SparseHTTPSource{}
	o.DeepCopyInto(out)

	return out
}

// DeepCopyInto copies the receiver into the given *SparseHTTPSource.
func (o *SparseHTTPSource) DeepCopyInto(out *SparseHTTPSource) {

	target, err := copystructure.Copy(o)
	if err != nil {
		panic(fmt.Sprintf("Unable to deepcopy SparseHTTPSource: %s", err))
	}

	*out = *target.(*SparseHTTPSource)
}

type mongoAttributesHTTPSource struct {
	CA          string            `bson:"ca"`
	ID          bson.ObjectId     `bson:"_id,omitempty"`
	Certificate string            `bson:"certificate"`
	Description string            `bson:"description"`
	Endpoint    string            `bson:"endpoint"`
	Key         string            `bson:"key"`
	Modifier    *IdentityModifier `bson:"modifier,omitempty"`
	Name        string            `bson:"name"`
	Namespace   string            `bson:"namespace"`
	ZHash       int               `bson:"zhash"`
	Zone        int               `bson:"zone"`
}
type mongoAttributesSparseHTTPSource struct {
	CA          *string           `bson:"ca,omitempty"`
	ID          bson.ObjectId     `bson:"_id,omitempty"`
	Certificate *string           `bson:"certificate,omitempty"`
	Description *string           `bson:"description,omitempty"`
	Endpoint    *string           `bson:"endpoint,omitempty"`
	Key         *string           `bson:"key,omitempty"`
	Modifier    *IdentityModifier `bson:"modifier,omitempty"`
	Name        *string           `bson:"name,omitempty"`
	Namespace   *string           `bson:"namespace,omitempty"`
	ZHash       *int              `bson:"zhash,omitempty"`
	Zone        *int              `bson:"zone,omitempty"`
}
