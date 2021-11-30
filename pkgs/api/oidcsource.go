package api

import (
	"fmt"

	"github.com/globalsign/mgo/bson"
	"github.com/mitchellh/copystructure"
	"go.aporeto.io/elemental"
)

// OIDCSourceIdentity represents the Identity of the object.
var OIDCSourceIdentity = elemental.Identity{
	Name:     "oidcsource",
	Category: "oidcsources",
	Package:  "a3s",
	Private:  false,
}

// OIDCSourcesList represents a list of OIDCSources
type OIDCSourcesList []*OIDCSource

// Identity returns the identity of the objects in the list.
func (o OIDCSourcesList) Identity() elemental.Identity {

	return OIDCSourceIdentity
}

// Copy returns a pointer to a copy the OIDCSourcesList.
func (o OIDCSourcesList) Copy() elemental.Identifiables {

	copy := append(OIDCSourcesList{}, o...)
	return &copy
}

// Append appends the objects to the a new copy of the OIDCSourcesList.
func (o OIDCSourcesList) Append(objects ...elemental.Identifiable) elemental.Identifiables {

	out := append(OIDCSourcesList{}, o...)
	for _, obj := range objects {
		out = append(out, obj.(*OIDCSource))
	}

	return out
}

// List converts the object to an elemental.IdentifiablesList.
func (o OIDCSourcesList) List() elemental.IdentifiablesList {

	out := make(elemental.IdentifiablesList, len(o))
	for i := 0; i < len(o); i++ {
		out[i] = o[i]
	}

	return out
}

// DefaultOrder returns the default ordering fields of the content.
func (o OIDCSourcesList) DefaultOrder() []string {

	return []string{}
}

// ToSparse returns the OIDCSourcesList converted to SparseOIDCSourcesList.
// Objects in the list will only contain the given fields. No field means entire field set.
func (o OIDCSourcesList) ToSparse(fields ...string) elemental.Identifiables {

	out := make(SparseOIDCSourcesList, len(o))
	for i := 0; i < len(o); i++ {
		out[i] = o[i].ToSparse(fields...).(*SparseOIDCSource)
	}

	return out
}

// Version returns the version of the content.
func (o OIDCSourcesList) Version() int {

	return 1
}

// OIDCSource represents the model of a oidcsource
type OIDCSource struct {
	// The Certificate authority to use to validate the authenticity of the OIDC
	// server. If left empty, the system trust stroe will be used. In most of the
	// cases, you don't need to set this.
	CA string `json:"CA" msgpack:"CA" bson:"ca" mapstructure:"CA,omitempty"`

	// ID is the identifier of the object.
	ID string `json:"ID" msgpack:"ID" bson:"-" mapstructure:"ID,omitempty"`

	// Unique client ID.
	ClientID string `json:"clientID" msgpack:"clientID" bson:"clientid" mapstructure:"clientID,omitempty"`

	// Client secret associated with the client ID.
	ClientSecret string `json:"clientSecret" msgpack:"clientSecret" bson:"clientsecret" mapstructure:"clientSecret,omitempty"`

	// The description of the object.
	Description string `json:"description" msgpack:"description" bson:"description" mapstructure:"description,omitempty"`

	// OIDC [discovery
	// endpoint](https://openid.net/specs/openid-connect-discovery-1_0.html#IssuerDiscovery).
	Endpoint string `json:"endpoint" msgpack:"endpoint" bson:"endpoint" mapstructure:"endpoint,omitempty"`

	// Contains optional information about a remote service that can be used to modify
	// the claims that are about to be delivered using this authentication source.
	Modifier *IdentityModifier `json:"modifier,omitempty" msgpack:"modifier,omitempty" bson:"-" mapstructure:"modifier,omitempty"`

	// The name of the source.
	Name string `json:"name" msgpack:"name" bson:"name" mapstructure:"name,omitempty"`

	// The namespace of the object.
	Namespace string `json:"namespace" msgpack:"namespace" bson:"namespace" mapstructure:"namespace,omitempty"`

	// List of scopes to allow.
	Scopes []string `json:"scopes" msgpack:"scopes" bson:"scopes" mapstructure:"scopes,omitempty"`

	// Hash of the object used to shard the data.
	ZHash int `json:"-" msgpack:"-" bson:"zhash" mapstructure:"-,omitempty"`

	// Sharding zone.
	Zone int `json:"-" msgpack:"-" bson:"zone" mapstructure:"-,omitempty"`

	ModelVersion int `json:"-" msgpack:"-" bson:"_modelversion"`
}

// NewOIDCSource returns a new *OIDCSource
func NewOIDCSource() *OIDCSource {

	return &OIDCSource{
		ModelVersion: 1,
		Scopes:       []string{},
	}
}

// Identity returns the Identity of the object.
func (o *OIDCSource) Identity() elemental.Identity {

	return OIDCSourceIdentity
}

// Identifier returns the value of the object's unique identifier.
func (o *OIDCSource) Identifier() string {

	return o.ID
}

// SetIdentifier sets the value of the object's unique identifier.
func (o *OIDCSource) SetIdentifier(id string) {

	o.ID = id
}

// GetBSON implements the bson marshaling interface.
// This is used to transparently convert ID to MongoDBID as ObectID.
func (o *OIDCSource) GetBSON() (interface{}, error) {

	if o == nil {
		return nil, nil
	}

	s := &mongoAttributesOIDCSource{}

	s.CA = o.CA
	if o.ID != "" {
		s.ID = bson.ObjectIdHex(o.ID)
	}
	s.ClientID = o.ClientID
	s.ClientSecret = o.ClientSecret
	s.Description = o.Description
	s.Endpoint = o.Endpoint
	s.Name = o.Name
	s.Namespace = o.Namespace
	s.Scopes = o.Scopes
	s.ZHash = o.ZHash
	s.Zone = o.Zone

	return s, nil
}

// SetBSON implements the bson marshaling interface.
// This is used to transparently convert ID to MongoDBID as ObectID.
func (o *OIDCSource) SetBSON(raw bson.Raw) error {

	if o == nil {
		return nil
	}

	s := &mongoAttributesOIDCSource{}
	if err := raw.Unmarshal(s); err != nil {
		return err
	}

	o.CA = s.CA
	o.ID = s.ID.Hex()
	o.ClientID = s.ClientID
	o.ClientSecret = s.ClientSecret
	o.Description = s.Description
	o.Endpoint = s.Endpoint
	o.Name = s.Name
	o.Namespace = s.Namespace
	o.Scopes = s.Scopes
	o.ZHash = s.ZHash
	o.Zone = s.Zone

	return nil
}

// Version returns the hardcoded version of the model.
func (o *OIDCSource) Version() int {

	return 1
}

// BleveType implements the bleve.Classifier Interface.
func (o *OIDCSource) BleveType() string {

	return "oidcsource"
}

// DefaultOrder returns the list of default ordering fields.
func (o *OIDCSource) DefaultOrder() []string {

	return []string{}
}

// Doc returns the documentation for the object
func (o *OIDCSource) Doc() string {

	return `An OIDC Auth source can be used to issue tokens based on existing OIDC accounts.`
}

func (o *OIDCSource) String() string {

	return fmt.Sprintf("<%s:%s>", o.Identity().Name, o.Identifier())
}

// GetID returns the ID of the receiver.
func (o *OIDCSource) GetID() string {

	return o.ID
}

// SetID sets the property ID of the receiver using the given value.
func (o *OIDCSource) SetID(ID string) {

	o.ID = ID
}

// GetNamespace returns the Namespace of the receiver.
func (o *OIDCSource) GetNamespace() string {

	return o.Namespace
}

// SetNamespace sets the property Namespace of the receiver using the given value.
func (o *OIDCSource) SetNamespace(namespace string) {

	o.Namespace = namespace
}

// GetZHash returns the ZHash of the receiver.
func (o *OIDCSource) GetZHash() int {

	return o.ZHash
}

// SetZHash sets the property ZHash of the receiver using the given value.
func (o *OIDCSource) SetZHash(zHash int) {

	o.ZHash = zHash
}

// GetZone returns the Zone of the receiver.
func (o *OIDCSource) GetZone() int {

	return o.Zone
}

// SetZone sets the property Zone of the receiver using the given value.
func (o *OIDCSource) SetZone(zone int) {

	o.Zone = zone
}

// ToSparse returns the sparse version of the model.
// The returned object will only contain the given fields. No field means entire field set.
func (o *OIDCSource) ToSparse(fields ...string) elemental.SparseIdentifiable {

	if len(fields) == 0 {
		// nolint: goimports
		return &SparseOIDCSource{
			CA:           &o.CA,
			ID:           &o.ID,
			ClientID:     &o.ClientID,
			ClientSecret: &o.ClientSecret,
			Description:  &o.Description,
			Endpoint:     &o.Endpoint,
			Modifier:     o.Modifier,
			Name:         &o.Name,
			Namespace:    &o.Namespace,
			Scopes:       &o.Scopes,
			ZHash:        &o.ZHash,
			Zone:         &o.Zone,
		}
	}

	sp := &SparseOIDCSource{}
	for _, f := range fields {
		switch f {
		case "CA":
			sp.CA = &(o.CA)
		case "ID":
			sp.ID = &(o.ID)
		case "clientID":
			sp.ClientID = &(o.ClientID)
		case "clientSecret":
			sp.ClientSecret = &(o.ClientSecret)
		case "description":
			sp.Description = &(o.Description)
		case "endpoint":
			sp.Endpoint = &(o.Endpoint)
		case "modifier":
			sp.Modifier = o.Modifier
		case "name":
			sp.Name = &(o.Name)
		case "namespace":
			sp.Namespace = &(o.Namespace)
		case "scopes":
			sp.Scopes = &(o.Scopes)
		case "zHash":
			sp.ZHash = &(o.ZHash)
		case "zone":
			sp.Zone = &(o.Zone)
		}
	}

	return sp
}

// EncryptAttributes encrypts the attributes marked as `encrypted` using the given encrypter.
func (o *OIDCSource) EncryptAttributes(encrypter elemental.AttributeEncrypter) (err error) {

	if o.ClientSecret, err = encrypter.EncryptString(o.ClientSecret); err != nil {
		return fmt.Errorf("unable to encrypt attribute 'ClientSecret' for 'OIDCSource' (%s): %s", o.Identifier(), err)
	}

	return nil
}

// DecryptAttributes decrypts the attributes marked as `encrypted` using the given decrypter.
func (o *OIDCSource) DecryptAttributes(encrypter elemental.AttributeEncrypter) (err error) {

	if o.ClientSecret, err = encrypter.DecryptString(o.ClientSecret); err != nil {
		return fmt.Errorf("unable to decrypt attribute 'ClientSecret' for 'OIDCSource' (%s): %s", o.Identifier(), err)
	}

	return nil
}

// Patch apply the non nil value of a *SparseOIDCSource to the object.
func (o *OIDCSource) Patch(sparse elemental.SparseIdentifiable) {
	if !sparse.Identity().IsEqual(o.Identity()) {
		panic("cannot patch from a parse with different identity")
	}

	so := sparse.(*SparseOIDCSource)
	if so.CA != nil {
		o.CA = *so.CA
	}
	if so.ID != nil {
		o.ID = *so.ID
	}
	if so.ClientID != nil {
		o.ClientID = *so.ClientID
	}
	if so.ClientSecret != nil {
		o.ClientSecret = *so.ClientSecret
	}
	if so.Description != nil {
		o.Description = *so.Description
	}
	if so.Endpoint != nil {
		o.Endpoint = *so.Endpoint
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
	if so.Scopes != nil {
		o.Scopes = *so.Scopes
	}
	if so.ZHash != nil {
		o.ZHash = *so.ZHash
	}
	if so.Zone != nil {
		o.Zone = *so.Zone
	}
}

// DeepCopy returns a deep copy if the OIDCSource.
func (o *OIDCSource) DeepCopy() *OIDCSource {

	if o == nil {
		return nil
	}

	out := &OIDCSource{}
	o.DeepCopyInto(out)

	return out
}

// DeepCopyInto copies the receiver into the given *OIDCSource.
func (o *OIDCSource) DeepCopyInto(out *OIDCSource) {

	target, err := copystructure.Copy(o)
	if err != nil {
		panic(fmt.Sprintf("Unable to deepcopy OIDCSource: %s", err))
	}

	*out = *target.(*OIDCSource)
}

// Validate valides the current information stored into the structure.
func (o *OIDCSource) Validate() error {

	errors := elemental.Errors{}
	requiredErrors := elemental.Errors{}

	if err := ValidatePEM("CA", o.CA); err != nil {
		errors = errors.Append(err)
	}

	if err := elemental.ValidateRequiredString("clientID", o.ClientID); err != nil {
		requiredErrors = requiredErrors.Append(err)
	}

	if err := elemental.ValidateRequiredString("clientSecret", o.ClientSecret); err != nil {
		requiredErrors = requiredErrors.Append(err)
	}

	if err := elemental.ValidateRequiredString("endpoint", o.Endpoint); err != nil {
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
func (*OIDCSource) SpecificationForAttribute(name string) elemental.AttributeSpecification {

	if v, ok := OIDCSourceAttributesMap[name]; ok {
		return v
	}

	// We could not find it, so let's check on the lower case indexed spec map
	return OIDCSourceLowerCaseAttributesMap[name]
}

// AttributeSpecifications returns the full attribute specifications map.
func (*OIDCSource) AttributeSpecifications() map[string]elemental.AttributeSpecification {

	return OIDCSourceAttributesMap
}

// ValueForAttribute returns the value for the given attribute.
// This is a very advanced function that you should not need but in some
// very specific use cases.
func (o *OIDCSource) ValueForAttribute(name string) interface{} {

	switch name {
	case "CA":
		return o.CA
	case "ID":
		return o.ID
	case "clientID":
		return o.ClientID
	case "clientSecret":
		return o.ClientSecret
	case "description":
		return o.Description
	case "endpoint":
		return o.Endpoint
	case "modifier":
		return o.Modifier
	case "name":
		return o.Name
	case "namespace":
		return o.Namespace
	case "scopes":
		return o.Scopes
	case "zHash":
		return o.ZHash
	case "zone":
		return o.Zone
	}

	return nil
}

// OIDCSourceAttributesMap represents the map of attribute for OIDCSource.
var OIDCSourceAttributesMap = map[string]elemental.AttributeSpecification{
	"CA": {
		AllowedChoices: []string{},
		BSONFieldName:  "ca",
		ConvertedName:  "CA",
		Description: `The Certificate authority to use to validate the authenticity of the OIDC
server. If left empty, the system trust stroe will be used. In most of the
cases, you don't need to set this.`,
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
	"ClientID": {
		AllowedChoices: []string{},
		BSONFieldName:  "clientid",
		ConvertedName:  "ClientID",
		Description:    `Unique client ID.`,
		Exposed:        true,
		Name:           "clientID",
		Required:       true,
		Stored:         true,
		Type:           "string",
	},
	"ClientSecret": {
		AllowedChoices: []string{},
		BSONFieldName:  "clientsecret",
		ConvertedName:  "ClientSecret",
		Description:    `Client secret associated with the client ID.`,
		Encrypted:      true,
		Exposed:        true,
		Name:           "clientSecret",
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
	"Endpoint": {
		AllowedChoices: []string{},
		BSONFieldName:  "endpoint",
		ConvertedName:  "Endpoint",
		Description: `OIDC [discovery
endpoint](https://openid.net/specs/openid-connect-discovery-1_0.html#IssuerDiscovery).`,
		Exposed:  true,
		Name:     "endpoint",
		Required: true,
		Stored:   true,
		Type:     "string",
	},
	"Modifier": {
		AllowedChoices: []string{},
		ConvertedName:  "Modifier",
		Description: `Contains optional information about a remote service that can be used to modify
the claims that are about to be delivered using this authentication source.`,
		Exposed: true,
		Name:    "modifier",
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
	"Scopes": {
		AllowedChoices: []string{},
		BSONFieldName:  "scopes",
		ConvertedName:  "Scopes",
		Description:    `List of scopes to allow.`,
		Exposed:        true,
		Name:           "scopes",
		Stored:         true,
		SubType:        "string",
		Type:           "list",
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

// OIDCSourceLowerCaseAttributesMap represents the map of attribute for OIDCSource.
var OIDCSourceLowerCaseAttributesMap = map[string]elemental.AttributeSpecification{
	"ca": {
		AllowedChoices: []string{},
		BSONFieldName:  "ca",
		ConvertedName:  "CA",
		Description: `The Certificate authority to use to validate the authenticity of the OIDC
server. If left empty, the system trust stroe will be used. In most of the
cases, you don't need to set this.`,
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
	"clientid": {
		AllowedChoices: []string{},
		BSONFieldName:  "clientid",
		ConvertedName:  "ClientID",
		Description:    `Unique client ID.`,
		Exposed:        true,
		Name:           "clientID",
		Required:       true,
		Stored:         true,
		Type:           "string",
	},
	"clientsecret": {
		AllowedChoices: []string{},
		BSONFieldName:  "clientsecret",
		ConvertedName:  "ClientSecret",
		Description:    `Client secret associated with the client ID.`,
		Encrypted:      true,
		Exposed:        true,
		Name:           "clientSecret",
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
	"endpoint": {
		AllowedChoices: []string{},
		BSONFieldName:  "endpoint",
		ConvertedName:  "Endpoint",
		Description: `OIDC [discovery
endpoint](https://openid.net/specs/openid-connect-discovery-1_0.html#IssuerDiscovery).`,
		Exposed:  true,
		Name:     "endpoint",
		Required: true,
		Stored:   true,
		Type:     "string",
	},
	"modifier": {
		AllowedChoices: []string{},
		ConvertedName:  "Modifier",
		Description: `Contains optional information about a remote service that can be used to modify
the claims that are about to be delivered using this authentication source.`,
		Exposed: true,
		Name:    "modifier",
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
	"scopes": {
		AllowedChoices: []string{},
		BSONFieldName:  "scopes",
		ConvertedName:  "Scopes",
		Description:    `List of scopes to allow.`,
		Exposed:        true,
		Name:           "scopes",
		Stored:         true,
		SubType:        "string",
		Type:           "list",
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

// SparseOIDCSourcesList represents a list of SparseOIDCSources
type SparseOIDCSourcesList []*SparseOIDCSource

// Identity returns the identity of the objects in the list.
func (o SparseOIDCSourcesList) Identity() elemental.Identity {

	return OIDCSourceIdentity
}

// Copy returns a pointer to a copy the SparseOIDCSourcesList.
func (o SparseOIDCSourcesList) Copy() elemental.Identifiables {

	copy := append(SparseOIDCSourcesList{}, o...)
	return &copy
}

// Append appends the objects to the a new copy of the SparseOIDCSourcesList.
func (o SparseOIDCSourcesList) Append(objects ...elemental.Identifiable) elemental.Identifiables {

	out := append(SparseOIDCSourcesList{}, o...)
	for _, obj := range objects {
		out = append(out, obj.(*SparseOIDCSource))
	}

	return out
}

// List converts the object to an elemental.IdentifiablesList.
func (o SparseOIDCSourcesList) List() elemental.IdentifiablesList {

	out := make(elemental.IdentifiablesList, len(o))
	for i := 0; i < len(o); i++ {
		out[i] = o[i]
	}

	return out
}

// DefaultOrder returns the default ordering fields of the content.
func (o SparseOIDCSourcesList) DefaultOrder() []string {

	return []string{}
}

// ToPlain returns the SparseOIDCSourcesList converted to OIDCSourcesList.
func (o SparseOIDCSourcesList) ToPlain() elemental.IdentifiablesList {

	out := make(elemental.IdentifiablesList, len(o))
	for i := 0; i < len(o); i++ {
		out[i] = o[i].ToPlain()
	}

	return out
}

// Version returns the version of the content.
func (o SparseOIDCSourcesList) Version() int {

	return 1
}

// SparseOIDCSource represents the sparse version of a oidcsource.
type SparseOIDCSource struct {
	// The Certificate authority to use to validate the authenticity of the OIDC
	// server. If left empty, the system trust stroe will be used. In most of the
	// cases, you don't need to set this.
	CA *string `json:"CA,omitempty" msgpack:"CA,omitempty" bson:"ca,omitempty" mapstructure:"CA,omitempty"`

	// ID is the identifier of the object.
	ID *string `json:"ID,omitempty" msgpack:"ID,omitempty" bson:"-" mapstructure:"ID,omitempty"`

	// Unique client ID.
	ClientID *string `json:"clientID,omitempty" msgpack:"clientID,omitempty" bson:"clientid,omitempty" mapstructure:"clientID,omitempty"`

	// Client secret associated with the client ID.
	ClientSecret *string `json:"clientSecret,omitempty" msgpack:"clientSecret,omitempty" bson:"clientsecret,omitempty" mapstructure:"clientSecret,omitempty"`

	// The description of the object.
	Description *string `json:"description,omitempty" msgpack:"description,omitempty" bson:"description,omitempty" mapstructure:"description,omitempty"`

	// OIDC [discovery
	// endpoint](https://openid.net/specs/openid-connect-discovery-1_0.html#IssuerDiscovery).
	Endpoint *string `json:"endpoint,omitempty" msgpack:"endpoint,omitempty" bson:"endpoint,omitempty" mapstructure:"endpoint,omitempty"`

	// Contains optional information about a remote service that can be used to modify
	// the claims that are about to be delivered using this authentication source.
	Modifier *IdentityModifier `json:"modifier,omitempty" msgpack:"modifier,omitempty" bson:"-" mapstructure:"modifier,omitempty"`

	// The name of the source.
	Name *string `json:"name,omitempty" msgpack:"name,omitempty" bson:"name,omitempty" mapstructure:"name,omitempty"`

	// The namespace of the object.
	Namespace *string `json:"namespace,omitempty" msgpack:"namespace,omitempty" bson:"namespace,omitempty" mapstructure:"namespace,omitempty"`

	// List of scopes to allow.
	Scopes *[]string `json:"scopes,omitempty" msgpack:"scopes,omitempty" bson:"scopes,omitempty" mapstructure:"scopes,omitempty"`

	// Hash of the object used to shard the data.
	ZHash *int `json:"-" msgpack:"-" bson:"zhash,omitempty" mapstructure:"-,omitempty"`

	// Sharding zone.
	Zone *int `json:"-" msgpack:"-" bson:"zone,omitempty" mapstructure:"-,omitempty"`

	ModelVersion int `json:"-" msgpack:"-" bson:"_modelversion"`
}

// NewSparseOIDCSource returns a new  SparseOIDCSource.
func NewSparseOIDCSource() *SparseOIDCSource {
	return &SparseOIDCSource{}
}

// Identity returns the Identity of the sparse object.
func (o *SparseOIDCSource) Identity() elemental.Identity {

	return OIDCSourceIdentity
}

// Identifier returns the value of the sparse object's unique identifier.
func (o *SparseOIDCSource) Identifier() string {

	if o.ID == nil {
		return ""
	}
	return *o.ID
}

// SetIdentifier sets the value of the sparse object's unique identifier.
func (o *SparseOIDCSource) SetIdentifier(id string) {

	if id != "" {
		o.ID = &id
	} else {
		o.ID = nil
	}
}

// GetBSON implements the bson marshaling interface.
// This is used to transparently convert ID to MongoDBID as ObectID.
func (o *SparseOIDCSource) GetBSON() (interface{}, error) {

	if o == nil {
		return nil, nil
	}

	s := &mongoAttributesSparseOIDCSource{}

	if o.CA != nil {
		s.CA = o.CA
	}
	if o.ID != nil {
		s.ID = bson.ObjectIdHex(*o.ID)
	}
	if o.ClientID != nil {
		s.ClientID = o.ClientID
	}
	if o.ClientSecret != nil {
		s.ClientSecret = o.ClientSecret
	}
	if o.Description != nil {
		s.Description = o.Description
	}
	if o.Endpoint != nil {
		s.Endpoint = o.Endpoint
	}
	if o.Name != nil {
		s.Name = o.Name
	}
	if o.Namespace != nil {
		s.Namespace = o.Namespace
	}
	if o.Scopes != nil {
		s.Scopes = o.Scopes
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
func (o *SparseOIDCSource) SetBSON(raw bson.Raw) error {

	if o == nil {
		return nil
	}

	s := &mongoAttributesSparseOIDCSource{}
	if err := raw.Unmarshal(s); err != nil {
		return err
	}

	if s.CA != nil {
		o.CA = s.CA
	}
	id := s.ID.Hex()
	o.ID = &id
	if s.ClientID != nil {
		o.ClientID = s.ClientID
	}
	if s.ClientSecret != nil {
		o.ClientSecret = s.ClientSecret
	}
	if s.Description != nil {
		o.Description = s.Description
	}
	if s.Endpoint != nil {
		o.Endpoint = s.Endpoint
	}
	if s.Name != nil {
		o.Name = s.Name
	}
	if s.Namespace != nil {
		o.Namespace = s.Namespace
	}
	if s.Scopes != nil {
		o.Scopes = s.Scopes
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
func (o *SparseOIDCSource) Version() int {

	return 1
}

// ToPlain returns the plain version of the sparse model.
func (o *SparseOIDCSource) ToPlain() elemental.PlainIdentifiable {

	out := NewOIDCSource()
	if o.CA != nil {
		out.CA = *o.CA
	}
	if o.ID != nil {
		out.ID = *o.ID
	}
	if o.ClientID != nil {
		out.ClientID = *o.ClientID
	}
	if o.ClientSecret != nil {
		out.ClientSecret = *o.ClientSecret
	}
	if o.Description != nil {
		out.Description = *o.Description
	}
	if o.Endpoint != nil {
		out.Endpoint = *o.Endpoint
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
	if o.Scopes != nil {
		out.Scopes = *o.Scopes
	}
	if o.ZHash != nil {
		out.ZHash = *o.ZHash
	}
	if o.Zone != nil {
		out.Zone = *o.Zone
	}

	return out
}

// EncryptAttributes encrypts the attributes marked as `encrypted` using the given encrypter.
func (o *SparseOIDCSource) EncryptAttributes(encrypter elemental.AttributeEncrypter) (err error) {

	if *o.ClientSecret, err = encrypter.EncryptString(*o.ClientSecret); err != nil {
		return fmt.Errorf("unable to encrypt attribute 'ClientSecret' for 'SparseOIDCSource' (%s): %s", o.Identifier(), err)
	}

	return nil
}

// DecryptAttributes decrypts the attributes marked as `encrypted` using the given decrypter.
func (o *SparseOIDCSource) DecryptAttributes(encrypter elemental.AttributeEncrypter) (err error) {

	if *o.ClientSecret, err = encrypter.DecryptString(*o.ClientSecret); err != nil {
		return fmt.Errorf("unable to decrypt attribute 'ClientSecret' for 'SparseOIDCSource' (%s): %s", o.Identifier(), err)
	}

	return nil
}

// GetID returns the ID of the receiver.
func (o *SparseOIDCSource) GetID() (out string) {

	if o.ID == nil {
		return
	}

	return *o.ID
}

// SetID sets the property ID of the receiver using the address of the given value.
func (o *SparseOIDCSource) SetID(ID string) {

	o.ID = &ID
}

// GetNamespace returns the Namespace of the receiver.
func (o *SparseOIDCSource) GetNamespace() (out string) {

	if o.Namespace == nil {
		return
	}

	return *o.Namespace
}

// SetNamespace sets the property Namespace of the receiver using the address of the given value.
func (o *SparseOIDCSource) SetNamespace(namespace string) {

	o.Namespace = &namespace
}

// GetZHash returns the ZHash of the receiver.
func (o *SparseOIDCSource) GetZHash() (out int) {

	if o.ZHash == nil {
		return
	}

	return *o.ZHash
}

// SetZHash sets the property ZHash of the receiver using the address of the given value.
func (o *SparseOIDCSource) SetZHash(zHash int) {

	o.ZHash = &zHash
}

// GetZone returns the Zone of the receiver.
func (o *SparseOIDCSource) GetZone() (out int) {

	if o.Zone == nil {
		return
	}

	return *o.Zone
}

// SetZone sets the property Zone of the receiver using the address of the given value.
func (o *SparseOIDCSource) SetZone(zone int) {

	o.Zone = &zone
}

// DeepCopy returns a deep copy if the SparseOIDCSource.
func (o *SparseOIDCSource) DeepCopy() *SparseOIDCSource {

	if o == nil {
		return nil
	}

	out := &SparseOIDCSource{}
	o.DeepCopyInto(out)

	return out
}

// DeepCopyInto copies the receiver into the given *SparseOIDCSource.
func (o *SparseOIDCSource) DeepCopyInto(out *SparseOIDCSource) {

	target, err := copystructure.Copy(o)
	if err != nil {
		panic(fmt.Sprintf("Unable to deepcopy SparseOIDCSource: %s", err))
	}

	*out = *target.(*SparseOIDCSource)
}

type mongoAttributesOIDCSource struct {
	CA           string        `bson:"ca"`
	ID           bson.ObjectId `bson:"_id,omitempty"`
	ClientID     string        `bson:"clientid"`
	ClientSecret string        `bson:"clientsecret"`
	Description  string        `bson:"description"`
	Endpoint     string        `bson:"endpoint"`
	Name         string        `bson:"name"`
	Namespace    string        `bson:"namespace"`
	Scopes       []string      `bson:"scopes"`
	ZHash        int           `bson:"zhash"`
	Zone         int           `bson:"zone"`
}
type mongoAttributesSparseOIDCSource struct {
	CA           *string       `bson:"ca,omitempty"`
	ID           bson.ObjectId `bson:"_id,omitempty"`
	ClientID     *string       `bson:"clientid,omitempty"`
	ClientSecret *string       `bson:"clientsecret,omitempty"`
	Description  *string       `bson:"description,omitempty"`
	Endpoint     *string       `bson:"endpoint,omitempty"`
	Name         *string       `bson:"name,omitempty"`
	Namespace    *string       `bson:"namespace,omitempty"`
	Scopes       *[]string     `bson:"scopes,omitempty"`
	ZHash        *int          `bson:"zhash,omitempty"`
	Zone         *int          `bson:"zone,omitempty"`
}
