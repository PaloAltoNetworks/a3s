package api

import (
	"fmt"

	"github.com/globalsign/mgo/bson"
	"github.com/mitchellh/copystructure"
	"go.aporeto.io/elemental"
)

// AuthorizationIdentity represents the Identity of the object.
var AuthorizationIdentity = elemental.Identity{
	Name:     "authorization",
	Category: "authorizations",
	Package:  "a3s",
	Private:  false,
}

// AuthorizationsList represents a list of Authorizations
type AuthorizationsList []*Authorization

// Identity returns the identity of the objects in the list.
func (o AuthorizationsList) Identity() elemental.Identity {

	return AuthorizationIdentity
}

// Copy returns a pointer to a copy the AuthorizationsList.
func (o AuthorizationsList) Copy() elemental.Identifiables {

	copy := append(AuthorizationsList{}, o...)
	return &copy
}

// Append appends the objects to the a new copy of the AuthorizationsList.
func (o AuthorizationsList) Append(objects ...elemental.Identifiable) elemental.Identifiables {

	out := append(AuthorizationsList{}, o...)
	for _, obj := range objects {
		out = append(out, obj.(*Authorization))
	}

	return out
}

// List converts the object to an elemental.IdentifiablesList.
func (o AuthorizationsList) List() elemental.IdentifiablesList {

	out := make(elemental.IdentifiablesList, len(o))
	for i := 0; i < len(o); i++ {
		out[i] = o[i]
	}

	return out
}

// DefaultOrder returns the default ordering fields of the content.
func (o AuthorizationsList) DefaultOrder() []string {

	return []string{}
}

// ToSparse returns the AuthorizationsList converted to SparseAuthorizationsList.
// Objects in the list will only contain the given fields. No field means entire field set.
func (o AuthorizationsList) ToSparse(fields ...string) elemental.Identifiables {

	out := make(SparseAuthorizationsList, len(o))
	for i := 0; i < len(o); i++ {
		out[i] = o[i].ToSparse(fields...).(*SparseAuthorization)
	}

	return out
}

// Version returns the version of the content.
func (o AuthorizationsList) Version() int {

	return 1
}

// Authorization represents the model of a authorization
type Authorization struct {
	// ID is the identifier of the object.
	ID string `json:"ID" msgpack:"ID" bson:"-" mapstructure:"ID,omitempty"`

	// Description of the Authorization.
	Description string `json:"description" msgpack:"description" bson:"description" mapstructure:"description,omitempty"`

	// Set the authorization to be disabled.
	Disabled bool `json:"disabled" msgpack:"disabled" bson:"disabled" mapstructure:"disabled,omitempty"`

	// This is a set of all subject tags for matching in the DB.
	FlattenedSubject []string `json:"-" msgpack:"-" bson:"flattenedsubject" mapstructure:"-,omitempty"`

	// Hides the policies in children namespaces.
	Hidden bool `json:"hidden" msgpack:"hidden" bson:"hidden" mapstructure:"hidden,omitempty"`

	// The name of the Authorization.
	Name string `json:"name" msgpack:"name" bson:"name" mapstructure:"name,omitempty"`

	// The namespace of the object.
	Namespace string `json:"namespace" msgpack:"namespace" bson:"namespace" mapstructure:"namespace,omitempty"`

	// A list of permissions.
	Permissions []string `json:"permissions" msgpack:"permissions" bson:"permissions" mapstructure:"permissions,omitempty"`

	// Propagates the api authorization to all of its children. This is always true.
	Propagate bool `json:"-" msgpack:"-" bson:"propagate" mapstructure:"-,omitempty"`

	// A tag expression that identifies the authorized user(s).
	Subject [][]string `json:"subject" msgpack:"subject" bson:"subject" mapstructure:"subject,omitempty"`

	// If set, the API authorization will only be valid if the request comes from one
	// the declared subnets.
	Subnets []string `json:"subnets" msgpack:"subnets" bson:"subnets" mapstructure:"subnets,omitempty"`

	// Defines the namespace or namespaces in which the permission for subject should
	// apply. If empty, the object's namespace will be used.
	TargetNamespaces []string `json:"targetNamespaces" msgpack:"targetNamespaces" bson:"targetnamespaces" mapstructure:"targetNamespaces,omitempty"`

	// Hash of the object used to shard the data.
	ZHash int `json:"-" msgpack:"-" bson:"zhash" mapstructure:"-,omitempty"`

	// Sharding zone.
	Zone int `json:"-" msgpack:"-" bson:"zone" mapstructure:"-,omitempty"`

	ModelVersion int `json:"-" msgpack:"-" bson:"_modelversion"`
}

// NewAuthorization returns a new *Authorization
func NewAuthorization() *Authorization {

	return &Authorization{
		ModelVersion:     1,
		FlattenedSubject: []string{},
		Permissions:      []string{},
		Propagate:        true,
		Subject:          [][]string{},
		Subnets:          []string{},
		TargetNamespaces: []string{},
	}
}

// Identity returns the Identity of the object.
func (o *Authorization) Identity() elemental.Identity {

	return AuthorizationIdentity
}

// Identifier returns the value of the object's unique identifier.
func (o *Authorization) Identifier() string {

	return o.ID
}

// SetIdentifier sets the value of the object's unique identifier.
func (o *Authorization) SetIdentifier(id string) {

	o.ID = id
}

// GetBSON implements the bson marshaling interface.
// This is used to transparently convert ID to MongoDBID as ObectID.
func (o *Authorization) GetBSON() (interface{}, error) {

	if o == nil {
		return nil, nil
	}

	s := &mongoAttributesAuthorization{}

	if o.ID != "" {
		s.ID = bson.ObjectIdHex(o.ID)
	}
	s.Description = o.Description
	s.Disabled = o.Disabled
	s.FlattenedSubject = o.FlattenedSubject
	s.Hidden = o.Hidden
	s.Name = o.Name
	s.Namespace = o.Namespace
	s.Permissions = o.Permissions
	s.Propagate = o.Propagate
	s.Subject = o.Subject
	s.Subnets = o.Subnets
	s.TargetNamespaces = o.TargetNamespaces
	s.ZHash = o.ZHash
	s.Zone = o.Zone

	return s, nil
}

// SetBSON implements the bson marshaling interface.
// This is used to transparently convert ID to MongoDBID as ObectID.
func (o *Authorization) SetBSON(raw bson.Raw) error {

	if o == nil {
		return nil
	}

	s := &mongoAttributesAuthorization{}
	if err := raw.Unmarshal(s); err != nil {
		return err
	}

	o.ID = s.ID.Hex()
	o.Description = s.Description
	o.Disabled = s.Disabled
	o.FlattenedSubject = s.FlattenedSubject
	o.Hidden = s.Hidden
	o.Name = s.Name
	o.Namespace = s.Namespace
	o.Permissions = s.Permissions
	o.Propagate = s.Propagate
	o.Subject = s.Subject
	o.Subnets = s.Subnets
	o.TargetNamespaces = s.TargetNamespaces
	o.ZHash = s.ZHash
	o.Zone = s.Zone

	return nil
}

// Version returns the hardcoded version of the model.
func (o *Authorization) Version() int {

	return 1
}

// BleveType implements the bleve.Classifier Interface.
func (o *Authorization) BleveType() string {

	return "authorization"
}

// DefaultOrder returns the list of default ordering fields.
func (o *Authorization) DefaultOrder() []string {

	return []string{}
}

// Doc returns the documentation for the object
func (o *Authorization) Doc() string {

	return `TODO.`
}

func (o *Authorization) String() string {

	return fmt.Sprintf("<%s:%s>", o.Identity().Name, o.Identifier())
}

// GetID returns the ID of the receiver.
func (o *Authorization) GetID() string {

	return o.ID
}

// SetID sets the property ID of the receiver using the given value.
func (o *Authorization) SetID(ID string) {

	o.ID = ID
}

// GetNamespace returns the Namespace of the receiver.
func (o *Authorization) GetNamespace() string {

	return o.Namespace
}

// SetNamespace sets the property Namespace of the receiver using the given value.
func (o *Authorization) SetNamespace(namespace string) {

	o.Namespace = namespace
}

// GetPropagate returns the Propagate of the receiver.
func (o *Authorization) GetPropagate() bool {

	return o.Propagate
}

// SetPropagate sets the property Propagate of the receiver using the given value.
func (o *Authorization) SetPropagate(propagate bool) {

	o.Propagate = propagate
}

// GetZHash returns the ZHash of the receiver.
func (o *Authorization) GetZHash() int {

	return o.ZHash
}

// SetZHash sets the property ZHash of the receiver using the given value.
func (o *Authorization) SetZHash(zHash int) {

	o.ZHash = zHash
}

// GetZone returns the Zone of the receiver.
func (o *Authorization) GetZone() int {

	return o.Zone
}

// SetZone sets the property Zone of the receiver using the given value.
func (o *Authorization) SetZone(zone int) {

	o.Zone = zone
}

// ToSparse returns the sparse version of the model.
// The returned object will only contain the given fields. No field means entire field set.
func (o *Authorization) ToSparse(fields ...string) elemental.SparseIdentifiable {

	if len(fields) == 0 {
		// nolint: goimports
		return &SparseAuthorization{
			ID:               &o.ID,
			Description:      &o.Description,
			Disabled:         &o.Disabled,
			FlattenedSubject: &o.FlattenedSubject,
			Hidden:           &o.Hidden,
			Name:             &o.Name,
			Namespace:        &o.Namespace,
			Permissions:      &o.Permissions,
			Propagate:        &o.Propagate,
			Subject:          &o.Subject,
			Subnets:          &o.Subnets,
			TargetNamespaces: &o.TargetNamespaces,
			ZHash:            &o.ZHash,
			Zone:             &o.Zone,
		}
	}

	sp := &SparseAuthorization{}
	for _, f := range fields {
		switch f {
		case "ID":
			sp.ID = &(o.ID)
		case "description":
			sp.Description = &(o.Description)
		case "disabled":
			sp.Disabled = &(o.Disabled)
		case "flattenedSubject":
			sp.FlattenedSubject = &(o.FlattenedSubject)
		case "hidden":
			sp.Hidden = &(o.Hidden)
		case "name":
			sp.Name = &(o.Name)
		case "namespace":
			sp.Namespace = &(o.Namespace)
		case "permissions":
			sp.Permissions = &(o.Permissions)
		case "propagate":
			sp.Propagate = &(o.Propagate)
		case "subject":
			sp.Subject = &(o.Subject)
		case "subnets":
			sp.Subnets = &(o.Subnets)
		case "targetNamespaces":
			sp.TargetNamespaces = &(o.TargetNamespaces)
		case "zHash":
			sp.ZHash = &(o.ZHash)
		case "zone":
			sp.Zone = &(o.Zone)
		}
	}

	return sp
}

// Patch apply the non nil value of a *SparseAuthorization to the object.
func (o *Authorization) Patch(sparse elemental.SparseIdentifiable) {
	if !sparse.Identity().IsEqual(o.Identity()) {
		panic("cannot patch from a parse with different identity")
	}

	so := sparse.(*SparseAuthorization)
	if so.ID != nil {
		o.ID = *so.ID
	}
	if so.Description != nil {
		o.Description = *so.Description
	}
	if so.Disabled != nil {
		o.Disabled = *so.Disabled
	}
	if so.FlattenedSubject != nil {
		o.FlattenedSubject = *so.FlattenedSubject
	}
	if so.Hidden != nil {
		o.Hidden = *so.Hidden
	}
	if so.Name != nil {
		o.Name = *so.Name
	}
	if so.Namespace != nil {
		o.Namespace = *so.Namespace
	}
	if so.Permissions != nil {
		o.Permissions = *so.Permissions
	}
	if so.Propagate != nil {
		o.Propagate = *so.Propagate
	}
	if so.Subject != nil {
		o.Subject = *so.Subject
	}
	if so.Subnets != nil {
		o.Subnets = *so.Subnets
	}
	if so.TargetNamespaces != nil {
		o.TargetNamespaces = *so.TargetNamespaces
	}
	if so.ZHash != nil {
		o.ZHash = *so.ZHash
	}
	if so.Zone != nil {
		o.Zone = *so.Zone
	}
}

// DeepCopy returns a deep copy if the Authorization.
func (o *Authorization) DeepCopy() *Authorization {

	if o == nil {
		return nil
	}

	out := &Authorization{}
	o.DeepCopyInto(out)

	return out
}

// DeepCopyInto copies the receiver into the given *Authorization.
func (o *Authorization) DeepCopyInto(out *Authorization) {

	target, err := copystructure.Copy(o)
	if err != nil {
		panic(fmt.Sprintf("Unable to deepcopy Authorization: %s", err))
	}

	*out = *target.(*Authorization)
}

// Validate valides the current information stored into the structure.
func (o *Authorization) Validate() error {

	errors := elemental.Errors{}
	requiredErrors := elemental.Errors{}

	if err := elemental.ValidateRequiredString("name", o.Name); err != nil {
		requiredErrors = requiredErrors.Append(err)
	}

	if err := elemental.ValidateRequiredExternal("permissions", o.Permissions); err != nil {
		requiredErrors = requiredErrors.Append(err)
	}

	if err := ValidateAuthorizationSubject("subject", o.Subject); err != nil {
		errors = errors.Append(err)
	}
	if err := ValidateTagsExpression("subject", o.Subject); err != nil {
		errors = errors.Append(err)
	}

	if err := ValidateCIDRListOptional("subnets", o.Subnets); err != nil {
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
func (*Authorization) SpecificationForAttribute(name string) elemental.AttributeSpecification {

	if v, ok := AuthorizationAttributesMap[name]; ok {
		return v
	}

	// We could not find it, so let's check on the lower case indexed spec map
	return AuthorizationLowerCaseAttributesMap[name]
}

// AttributeSpecifications returns the full attribute specifications map.
func (*Authorization) AttributeSpecifications() map[string]elemental.AttributeSpecification {

	return AuthorizationAttributesMap
}

// ValueForAttribute returns the value for the given attribute.
// This is a very advanced function that you should not need but in some
// very specific use cases.
func (o *Authorization) ValueForAttribute(name string) interface{} {

	switch name {
	case "ID":
		return o.ID
	case "description":
		return o.Description
	case "disabled":
		return o.Disabled
	case "flattenedSubject":
		return o.FlattenedSubject
	case "hidden":
		return o.Hidden
	case "name":
		return o.Name
	case "namespace":
		return o.Namespace
	case "permissions":
		return o.Permissions
	case "propagate":
		return o.Propagate
	case "subject":
		return o.Subject
	case "subnets":
		return o.Subnets
	case "targetNamespaces":
		return o.TargetNamespaces
	case "zHash":
		return o.ZHash
	case "zone":
		return o.Zone
	}

	return nil
}

// AuthorizationAttributesMap represents the map of attribute for Authorization.
var AuthorizationAttributesMap = map[string]elemental.AttributeSpecification{
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
	"Description": {
		AllowedChoices: []string{},
		BSONFieldName:  "description",
		ConvertedName:  "Description",
		Description:    `Description of the Authorization.`,
		Exposed:        true,
		Name:           "description",
		Stored:         true,
		Type:           "string",
	},
	"Disabled": {
		AllowedChoices: []string{},
		BSONFieldName:  "disabled",
		ConvertedName:  "Disabled",
		Description:    `Set the authorization to be disabled.`,
		Exposed:        true,
		Name:           "disabled",
		Stored:         true,
		Type:           "boolean",
	},
	"FlattenedSubject": {
		AllowedChoices: []string{},
		BSONFieldName:  "flattenedsubject",
		ConvertedName:  "FlattenedSubject",
		Description:    `This is a set of all subject tags for matching in the DB.`,
		Name:           "flattenedSubject",
		Stored:         true,
		SubType:        "string",
		Type:           "list",
	},
	"Hidden": {
		AllowedChoices: []string{},
		BSONFieldName:  "hidden",
		ConvertedName:  "Hidden",
		Description:    `Hides the policies in children namespaces.`,
		Exposed:        true,
		Name:           "hidden",
		Stored:         true,
		Type:           "boolean",
	},
	"Name": {
		AllowedChoices: []string{},
		BSONFieldName:  "name",
		ConvertedName:  "Name",
		Description:    `The name of the Authorization.`,
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
	"Permissions": {
		AllowedChoices: []string{},
		BSONFieldName:  "permissions",
		ConvertedName:  "Permissions",
		Description:    `A list of permissions.`,
		Exposed:        true,
		Name:           "permissions",
		Required:       true,
		Stored:         true,
		SubType:        "string",
		Type:           "list",
	},
	"Propagate": {
		AllowedChoices: []string{},
		BSONFieldName:  "propagate",
		ConvertedName:  "Propagate",
		DefaultValue:   true,
		Description:    `Propagates the api authorization to all of its children. This is always true.`,
		Getter:         true,
		Name:           "propagate",
		Setter:         true,
		Stored:         true,
		Type:           "boolean",
	},
	"Subject": {
		AllowedChoices: []string{},
		BSONFieldName:  "subject",
		ConvertedName:  "Subject",
		Description:    `A tag expression that identifies the authorized user(s).`,
		Exposed:        true,
		Name:           "subject",
		Orderable:      true,
		Stored:         true,
		SubType:        "[][]string",
		Type:           "external",
	},
	"Subnets": {
		AllowedChoices: []string{},
		BSONFieldName:  "subnets",
		ConvertedName:  "Subnets",
		Description: `If set, the API authorization will only be valid if the request comes from one
the declared subnets.`,
		Exposed: true,
		Name:    "subnets",
		Stored:  true,
		SubType: "string",
		Type:    "list",
	},
	"TargetNamespaces": {
		AllowedChoices: []string{},
		BSONFieldName:  "targetnamespaces",
		ConvertedName:  "TargetNamespaces",
		Description: `Defines the namespace or namespaces in which the permission for subject should
apply. If empty, the object's namespace will be used.`,
		Exposed: true,
		Name:    "targetNamespaces",
		Stored:  true,
		SubType: "string",
		Type:    "list",
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

// AuthorizationLowerCaseAttributesMap represents the map of attribute for Authorization.
var AuthorizationLowerCaseAttributesMap = map[string]elemental.AttributeSpecification{
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
	"description": {
		AllowedChoices: []string{},
		BSONFieldName:  "description",
		ConvertedName:  "Description",
		Description:    `Description of the Authorization.`,
		Exposed:        true,
		Name:           "description",
		Stored:         true,
		Type:           "string",
	},
	"disabled": {
		AllowedChoices: []string{},
		BSONFieldName:  "disabled",
		ConvertedName:  "Disabled",
		Description:    `Set the authorization to be disabled.`,
		Exposed:        true,
		Name:           "disabled",
		Stored:         true,
		Type:           "boolean",
	},
	"flattenedsubject": {
		AllowedChoices: []string{},
		BSONFieldName:  "flattenedsubject",
		ConvertedName:  "FlattenedSubject",
		Description:    `This is a set of all subject tags for matching in the DB.`,
		Name:           "flattenedSubject",
		Stored:         true,
		SubType:        "string",
		Type:           "list",
	},
	"hidden": {
		AllowedChoices: []string{},
		BSONFieldName:  "hidden",
		ConvertedName:  "Hidden",
		Description:    `Hides the policies in children namespaces.`,
		Exposed:        true,
		Name:           "hidden",
		Stored:         true,
		Type:           "boolean",
	},
	"name": {
		AllowedChoices: []string{},
		BSONFieldName:  "name",
		ConvertedName:  "Name",
		Description:    `The name of the Authorization.`,
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
	"permissions": {
		AllowedChoices: []string{},
		BSONFieldName:  "permissions",
		ConvertedName:  "Permissions",
		Description:    `A list of permissions.`,
		Exposed:        true,
		Name:           "permissions",
		Required:       true,
		Stored:         true,
		SubType:        "string",
		Type:           "list",
	},
	"propagate": {
		AllowedChoices: []string{},
		BSONFieldName:  "propagate",
		ConvertedName:  "Propagate",
		DefaultValue:   true,
		Description:    `Propagates the api authorization to all of its children. This is always true.`,
		Getter:         true,
		Name:           "propagate",
		Setter:         true,
		Stored:         true,
		Type:           "boolean",
	},
	"subject": {
		AllowedChoices: []string{},
		BSONFieldName:  "subject",
		ConvertedName:  "Subject",
		Description:    `A tag expression that identifies the authorized user(s).`,
		Exposed:        true,
		Name:           "subject",
		Orderable:      true,
		Stored:         true,
		SubType:        "[][]string",
		Type:           "external",
	},
	"subnets": {
		AllowedChoices: []string{},
		BSONFieldName:  "subnets",
		ConvertedName:  "Subnets",
		Description: `If set, the API authorization will only be valid if the request comes from one
the declared subnets.`,
		Exposed: true,
		Name:    "subnets",
		Stored:  true,
		SubType: "string",
		Type:    "list",
	},
	"targetnamespaces": {
		AllowedChoices: []string{},
		BSONFieldName:  "targetnamespaces",
		ConvertedName:  "TargetNamespaces",
		Description: `Defines the namespace or namespaces in which the permission for subject should
apply. If empty, the object's namespace will be used.`,
		Exposed: true,
		Name:    "targetNamespaces",
		Stored:  true,
		SubType: "string",
		Type:    "list",
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

// SparseAuthorizationsList represents a list of SparseAuthorizations
type SparseAuthorizationsList []*SparseAuthorization

// Identity returns the identity of the objects in the list.
func (o SparseAuthorizationsList) Identity() elemental.Identity {

	return AuthorizationIdentity
}

// Copy returns a pointer to a copy the SparseAuthorizationsList.
func (o SparseAuthorizationsList) Copy() elemental.Identifiables {

	copy := append(SparseAuthorizationsList{}, o...)
	return &copy
}

// Append appends the objects to the a new copy of the SparseAuthorizationsList.
func (o SparseAuthorizationsList) Append(objects ...elemental.Identifiable) elemental.Identifiables {

	out := append(SparseAuthorizationsList{}, o...)
	for _, obj := range objects {
		out = append(out, obj.(*SparseAuthorization))
	}

	return out
}

// List converts the object to an elemental.IdentifiablesList.
func (o SparseAuthorizationsList) List() elemental.IdentifiablesList {

	out := make(elemental.IdentifiablesList, len(o))
	for i := 0; i < len(o); i++ {
		out[i] = o[i]
	}

	return out
}

// DefaultOrder returns the default ordering fields of the content.
func (o SparseAuthorizationsList) DefaultOrder() []string {

	return []string{}
}

// ToPlain returns the SparseAuthorizationsList converted to AuthorizationsList.
func (o SparseAuthorizationsList) ToPlain() elemental.IdentifiablesList {

	out := make(elemental.IdentifiablesList, len(o))
	for i := 0; i < len(o); i++ {
		out[i] = o[i].ToPlain()
	}

	return out
}

// Version returns the version of the content.
func (o SparseAuthorizationsList) Version() int {

	return 1
}

// SparseAuthorization represents the sparse version of a authorization.
type SparseAuthorization struct {
	// ID is the identifier of the object.
	ID *string `json:"ID,omitempty" msgpack:"ID,omitempty" bson:"-" mapstructure:"ID,omitempty"`

	// Description of the Authorization.
	Description *string `json:"description,omitempty" msgpack:"description,omitempty" bson:"description,omitempty" mapstructure:"description,omitempty"`

	// Set the authorization to be disabled.
	Disabled *bool `json:"disabled,omitempty" msgpack:"disabled,omitempty" bson:"disabled,omitempty" mapstructure:"disabled,omitempty"`

	// This is a set of all subject tags for matching in the DB.
	FlattenedSubject *[]string `json:"-" msgpack:"-" bson:"flattenedsubject,omitempty" mapstructure:"-,omitempty"`

	// Hides the policies in children namespaces.
	Hidden *bool `json:"hidden,omitempty" msgpack:"hidden,omitempty" bson:"hidden,omitempty" mapstructure:"hidden,omitempty"`

	// The name of the Authorization.
	Name *string `json:"name,omitempty" msgpack:"name,omitempty" bson:"name,omitempty" mapstructure:"name,omitempty"`

	// The namespace of the object.
	Namespace *string `json:"namespace,omitempty" msgpack:"namespace,omitempty" bson:"namespace,omitempty" mapstructure:"namespace,omitempty"`

	// A list of permissions.
	Permissions *[]string `json:"permissions,omitempty" msgpack:"permissions,omitempty" bson:"permissions,omitempty" mapstructure:"permissions,omitempty"`

	// Propagates the api authorization to all of its children. This is always true.
	Propagate *bool `json:"-" msgpack:"-" bson:"propagate,omitempty" mapstructure:"-,omitempty"`

	// A tag expression that identifies the authorized user(s).
	Subject *[][]string `json:"subject,omitempty" msgpack:"subject,omitempty" bson:"subject,omitempty" mapstructure:"subject,omitempty"`

	// If set, the API authorization will only be valid if the request comes from one
	// the declared subnets.
	Subnets *[]string `json:"subnets,omitempty" msgpack:"subnets,omitempty" bson:"subnets,omitempty" mapstructure:"subnets,omitempty"`

	// Defines the namespace or namespaces in which the permission for subject should
	// apply. If empty, the object's namespace will be used.
	TargetNamespaces *[]string `json:"targetNamespaces,omitempty" msgpack:"targetNamespaces,omitempty" bson:"targetnamespaces,omitempty" mapstructure:"targetNamespaces,omitempty"`

	// Hash of the object used to shard the data.
	ZHash *int `json:"-" msgpack:"-" bson:"zhash,omitempty" mapstructure:"-,omitempty"`

	// Sharding zone.
	Zone *int `json:"-" msgpack:"-" bson:"zone,omitempty" mapstructure:"-,omitempty"`

	ModelVersion int `json:"-" msgpack:"-" bson:"_modelversion"`
}

// NewSparseAuthorization returns a new  SparseAuthorization.
func NewSparseAuthorization() *SparseAuthorization {
	return &SparseAuthorization{}
}

// Identity returns the Identity of the sparse object.
func (o *SparseAuthorization) Identity() elemental.Identity {

	return AuthorizationIdentity
}

// Identifier returns the value of the sparse object's unique identifier.
func (o *SparseAuthorization) Identifier() string {

	if o.ID == nil {
		return ""
	}
	return *o.ID
}

// SetIdentifier sets the value of the sparse object's unique identifier.
func (o *SparseAuthorization) SetIdentifier(id string) {

	if id != "" {
		o.ID = &id
	} else {
		o.ID = nil
	}
}

// GetBSON implements the bson marshaling interface.
// This is used to transparently convert ID to MongoDBID as ObectID.
func (o *SparseAuthorization) GetBSON() (interface{}, error) {

	if o == nil {
		return nil, nil
	}

	s := &mongoAttributesSparseAuthorization{}

	if o.ID != nil {
		s.ID = bson.ObjectIdHex(*o.ID)
	}
	if o.Description != nil {
		s.Description = o.Description
	}
	if o.Disabled != nil {
		s.Disabled = o.Disabled
	}
	if o.FlattenedSubject != nil {
		s.FlattenedSubject = o.FlattenedSubject
	}
	if o.Hidden != nil {
		s.Hidden = o.Hidden
	}
	if o.Name != nil {
		s.Name = o.Name
	}
	if o.Namespace != nil {
		s.Namespace = o.Namespace
	}
	if o.Permissions != nil {
		s.Permissions = o.Permissions
	}
	if o.Propagate != nil {
		s.Propagate = o.Propagate
	}
	if o.Subject != nil {
		s.Subject = o.Subject
	}
	if o.Subnets != nil {
		s.Subnets = o.Subnets
	}
	if o.TargetNamespaces != nil {
		s.TargetNamespaces = o.TargetNamespaces
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
func (o *SparseAuthorization) SetBSON(raw bson.Raw) error {

	if o == nil {
		return nil
	}

	s := &mongoAttributesSparseAuthorization{}
	if err := raw.Unmarshal(s); err != nil {
		return err
	}

	id := s.ID.Hex()
	o.ID = &id
	if s.Description != nil {
		o.Description = s.Description
	}
	if s.Disabled != nil {
		o.Disabled = s.Disabled
	}
	if s.FlattenedSubject != nil {
		o.FlattenedSubject = s.FlattenedSubject
	}
	if s.Hidden != nil {
		o.Hidden = s.Hidden
	}
	if s.Name != nil {
		o.Name = s.Name
	}
	if s.Namespace != nil {
		o.Namespace = s.Namespace
	}
	if s.Permissions != nil {
		o.Permissions = s.Permissions
	}
	if s.Propagate != nil {
		o.Propagate = s.Propagate
	}
	if s.Subject != nil {
		o.Subject = s.Subject
	}
	if s.Subnets != nil {
		o.Subnets = s.Subnets
	}
	if s.TargetNamespaces != nil {
		o.TargetNamespaces = s.TargetNamespaces
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
func (o *SparseAuthorization) Version() int {

	return 1
}

// ToPlain returns the plain version of the sparse model.
func (o *SparseAuthorization) ToPlain() elemental.PlainIdentifiable {

	out := NewAuthorization()
	if o.ID != nil {
		out.ID = *o.ID
	}
	if o.Description != nil {
		out.Description = *o.Description
	}
	if o.Disabled != nil {
		out.Disabled = *o.Disabled
	}
	if o.FlattenedSubject != nil {
		out.FlattenedSubject = *o.FlattenedSubject
	}
	if o.Hidden != nil {
		out.Hidden = *o.Hidden
	}
	if o.Name != nil {
		out.Name = *o.Name
	}
	if o.Namespace != nil {
		out.Namespace = *o.Namespace
	}
	if o.Permissions != nil {
		out.Permissions = *o.Permissions
	}
	if o.Propagate != nil {
		out.Propagate = *o.Propagate
	}
	if o.Subject != nil {
		out.Subject = *o.Subject
	}
	if o.Subnets != nil {
		out.Subnets = *o.Subnets
	}
	if o.TargetNamespaces != nil {
		out.TargetNamespaces = *o.TargetNamespaces
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
func (o *SparseAuthorization) GetID() (out string) {

	if o.ID == nil {
		return
	}

	return *o.ID
}

// SetID sets the property ID of the receiver using the address of the given value.
func (o *SparseAuthorization) SetID(ID string) {

	o.ID = &ID
}

// GetNamespace returns the Namespace of the receiver.
func (o *SparseAuthorization) GetNamespace() (out string) {

	if o.Namespace == nil {
		return
	}

	return *o.Namespace
}

// SetNamespace sets the property Namespace of the receiver using the address of the given value.
func (o *SparseAuthorization) SetNamespace(namespace string) {

	o.Namespace = &namespace
}

// GetPropagate returns the Propagate of the receiver.
func (o *SparseAuthorization) GetPropagate() (out bool) {

	if o.Propagate == nil {
		return
	}

	return *o.Propagate
}

// SetPropagate sets the property Propagate of the receiver using the address of the given value.
func (o *SparseAuthorization) SetPropagate(propagate bool) {

	o.Propagate = &propagate
}

// GetZHash returns the ZHash of the receiver.
func (o *SparseAuthorization) GetZHash() (out int) {

	if o.ZHash == nil {
		return
	}

	return *o.ZHash
}

// SetZHash sets the property ZHash of the receiver using the address of the given value.
func (o *SparseAuthorization) SetZHash(zHash int) {

	o.ZHash = &zHash
}

// GetZone returns the Zone of the receiver.
func (o *SparseAuthorization) GetZone() (out int) {

	if o.Zone == nil {
		return
	}

	return *o.Zone
}

// SetZone sets the property Zone of the receiver using the address of the given value.
func (o *SparseAuthorization) SetZone(zone int) {

	o.Zone = &zone
}

// DeepCopy returns a deep copy if the SparseAuthorization.
func (o *SparseAuthorization) DeepCopy() *SparseAuthorization {

	if o == nil {
		return nil
	}

	out := &SparseAuthorization{}
	o.DeepCopyInto(out)

	return out
}

// DeepCopyInto copies the receiver into the given *SparseAuthorization.
func (o *SparseAuthorization) DeepCopyInto(out *SparseAuthorization) {

	target, err := copystructure.Copy(o)
	if err != nil {
		panic(fmt.Sprintf("Unable to deepcopy SparseAuthorization: %s", err))
	}

	*out = *target.(*SparseAuthorization)
}

type mongoAttributesAuthorization struct {
	ID               bson.ObjectId `bson:"_id,omitempty"`
	Description      string        `bson:"description"`
	Disabled         bool          `bson:"disabled"`
	FlattenedSubject []string      `bson:"flattenedsubject"`
	Hidden           bool          `bson:"hidden"`
	Name             string        `bson:"name"`
	Namespace        string        `bson:"namespace"`
	Permissions      []string      `bson:"permissions"`
	Propagate        bool          `bson:"propagate"`
	Subject          [][]string    `bson:"subject"`
	Subnets          []string      `bson:"subnets"`
	TargetNamespaces []string      `bson:"targetnamespaces"`
	ZHash            int           `bson:"zhash"`
	Zone             int           `bson:"zone"`
}
type mongoAttributesSparseAuthorization struct {
	ID               bson.ObjectId `bson:"_id,omitempty"`
	Description      *string       `bson:"description,omitempty"`
	Disabled         *bool         `bson:"disabled,omitempty"`
	FlattenedSubject *[]string     `bson:"flattenedsubject,omitempty"`
	Hidden           *bool         `bson:"hidden,omitempty"`
	Name             *string       `bson:"name,omitempty"`
	Namespace        *string       `bson:"namespace,omitempty"`
	Permissions      *[]string     `bson:"permissions,omitempty"`
	Propagate        *bool         `bson:"propagate,omitempty"`
	Subject          *[][]string   `bson:"subject,omitempty"`
	Subnets          *[]string     `bson:"subnets,omitempty"`
	TargetNamespaces *[]string     `bson:"targetnamespaces,omitempty"`
	ZHash            *int          `bson:"zhash,omitempty"`
	Zone             *int          `bson:"zone,omitempty"`
}
