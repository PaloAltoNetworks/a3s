package api

import (
	"fmt"

	"github.com/globalsign/mgo/bson"
	"github.com/mitchellh/copystructure"
	"go.aporeto.io/elemental"
)

// ImportIdentity represents the Identity of the object.
var ImportIdentity = elemental.Identity{
	Name:     "import",
	Category: "import",
	Package:  "a3s",
	Private:  false,
}

// ImportsList represents a list of Imports
type ImportsList []*Import

// Identity returns the identity of the objects in the list.
func (o ImportsList) Identity() elemental.Identity {

	return ImportIdentity
}

// Copy returns a pointer to a copy the ImportsList.
func (o ImportsList) Copy() elemental.Identifiables {

	copy := append(ImportsList{}, o...)
	return &copy
}

// Append appends the objects to the a new copy of the ImportsList.
func (o ImportsList) Append(objects ...elemental.Identifiable) elemental.Identifiables {

	out := append(ImportsList{}, o...)
	for _, obj := range objects {
		out = append(out, obj.(*Import))
	}

	return out
}

// List converts the object to an elemental.IdentifiablesList.
func (o ImportsList) List() elemental.IdentifiablesList {

	out := make(elemental.IdentifiablesList, len(o))
	for i := 0; i < len(o); i++ {
		out[i] = o[i]
	}

	return out
}

// DefaultOrder returns the default ordering fields of the content.
func (o ImportsList) DefaultOrder() []string {

	return []string{}
}

// ToSparse returns the ImportsList converted to SparseImportsList.
// Objects in the list will only contain the given fields. No field means entire field set.
func (o ImportsList) ToSparse(fields ...string) elemental.Identifiables {

	out := make(SparseImportsList, len(o))
	for i := 0; i < len(o); i++ {
		out[i] = o[i].ToSparse(fields...).(*SparseImport)
	}

	return out
}

// Version returns the version of the content.
func (o ImportsList) Version() int {

	return 1
}

// Import represents the model of a import
type Import struct {
	// A3S sources to import.
	A3SSources A3SSourcesList `json:"A3SSources,omitempty" msgpack:"A3SSources,omitempty" bson:"-" mapstructure:"A3SSources,omitempty"`

	// HTTP sources to import.
	HTTPSources HTTPSourcesList `json:"HTTPSources,omitempty" msgpack:"HTTPSources,omitempty" bson:"-" mapstructure:"HTTPSources,omitempty"`

	// LDAP sources to import.
	LDAPSources LDAPSourcesList `json:"LDAPSources,omitempty" msgpack:"LDAPSources,omitempty" bson:"-" mapstructure:"LDAPSources,omitempty"`

	// MTLS sources to import.
	MTLSSources MTLSSourcesList `json:"MTLSSources,omitempty" msgpack:"MTLSSources,omitempty" bson:"-" mapstructure:"MTLSSources,omitempty"`

	// OIDC sources to import.
	OIDCSources OIDCSourcesList `json:"OIDCSources,omitempty" msgpack:"OIDCSources,omitempty" bson:"-" mapstructure:"OIDCSources,omitempty"`

	// Authorizations to import.
	Authorizations AuthorizationsList `json:"authorizations,omitempty" msgpack:"authorizations,omitempty" bson:"-" mapstructure:"authorizations,omitempty"`

	// Import label that will be used to identify all the resources imported by this
	// resource.
	Label string `json:"label" msgpack:"label" bson:"-" mapstructure:"label,omitempty"`

	ModelVersion int `json:"-" msgpack:"-" bson:"_modelversion"`
}

// NewImport returns a new *Import
func NewImport() *Import {

	return &Import{
		ModelVersion:   1,
		A3SSources:     A3SSourcesList{},
		HTTPSources:    HTTPSourcesList{},
		LDAPSources:    LDAPSourcesList{},
		MTLSSources:    MTLSSourcesList{},
		Authorizations: AuthorizationsList{},
		OIDCSources:    OIDCSourcesList{},
	}
}

// Identity returns the Identity of the object.
func (o *Import) Identity() elemental.Identity {

	return ImportIdentity
}

// Identifier returns the value of the object's unique identifier.
func (o *Import) Identifier() string {

	return ""
}

// SetIdentifier sets the value of the object's unique identifier.
func (o *Import) SetIdentifier(id string) {

}

// GetBSON implements the bson marshaling interface.
// This is used to transparently convert ID to MongoDBID as ObectID.
func (o *Import) GetBSON() (interface{}, error) {

	if o == nil {
		return nil, nil
	}

	s := &mongoAttributesImport{}

	return s, nil
}

// SetBSON implements the bson marshaling interface.
// This is used to transparently convert ID to MongoDBID as ObectID.
func (o *Import) SetBSON(raw bson.Raw) error {

	if o == nil {
		return nil
	}

	s := &mongoAttributesImport{}
	if err := raw.Unmarshal(s); err != nil {
		return err
	}

	return nil
}

// Version returns the hardcoded version of the model.
func (o *Import) Version() int {

	return 1
}

// BleveType implements the bleve.Classifier Interface.
func (o *Import) BleveType() string {

	return "import"
}

// DefaultOrder returns the list of default ordering fields.
func (o *Import) DefaultOrder() []string {

	return []string{}
}

// Doc returns the documentation for the object
func (o *Import) Doc() string {

	return `Import multiple resource at once.`
}

func (o *Import) String() string {

	return fmt.Sprintf("<%s:%s>", o.Identity().Name, o.Identifier())
}

// ToSparse returns the sparse version of the model.
// The returned object will only contain the given fields. No field means entire field set.
func (o *Import) ToSparse(fields ...string) elemental.SparseIdentifiable {

	if len(fields) == 0 {
		// nolint: goimports
		return &SparseImport{
			A3SSources:     &o.A3SSources,
			HTTPSources:    &o.HTTPSources,
			LDAPSources:    &o.LDAPSources,
			MTLSSources:    &o.MTLSSources,
			OIDCSources:    &o.OIDCSources,
			Authorizations: &o.Authorizations,
			Label:          &o.Label,
		}
	}

	sp := &SparseImport{}
	for _, f := range fields {
		switch f {
		case "A3SSources":
			sp.A3SSources = &(o.A3SSources)
		case "HTTPSources":
			sp.HTTPSources = &(o.HTTPSources)
		case "LDAPSources":
			sp.LDAPSources = &(o.LDAPSources)
		case "MTLSSources":
			sp.MTLSSources = &(o.MTLSSources)
		case "OIDCSources":
			sp.OIDCSources = &(o.OIDCSources)
		case "authorizations":
			sp.Authorizations = &(o.Authorizations)
		case "label":
			sp.Label = &(o.Label)
		}
	}

	return sp
}

// Patch apply the non nil value of a *SparseImport to the object.
func (o *Import) Patch(sparse elemental.SparseIdentifiable) {
	if !sparse.Identity().IsEqual(o.Identity()) {
		panic("cannot patch from a parse with different identity")
	}

	so := sparse.(*SparseImport)
	if so.A3SSources != nil {
		o.A3SSources = *so.A3SSources
	}
	if so.HTTPSources != nil {
		o.HTTPSources = *so.HTTPSources
	}
	if so.LDAPSources != nil {
		o.LDAPSources = *so.LDAPSources
	}
	if so.MTLSSources != nil {
		o.MTLSSources = *so.MTLSSources
	}
	if so.OIDCSources != nil {
		o.OIDCSources = *so.OIDCSources
	}
	if so.Authorizations != nil {
		o.Authorizations = *so.Authorizations
	}
	if so.Label != nil {
		o.Label = *so.Label
	}
}

// DeepCopy returns a deep copy if the Import.
func (o *Import) DeepCopy() *Import {

	if o == nil {
		return nil
	}

	out := &Import{}
	o.DeepCopyInto(out)

	return out
}

// DeepCopyInto copies the receiver into the given *Import.
func (o *Import) DeepCopyInto(out *Import) {

	target, err := copystructure.Copy(o)
	if err != nil {
		panic(fmt.Sprintf("Unable to deepcopy Import: %s", err))
	}

	*out = *target.(*Import)
}

// Validate valides the current information stored into the structure.
func (o *Import) Validate() error {

	errors := elemental.Errors{}
	requiredErrors := elemental.Errors{}

	for _, sub := range o.A3SSources {
		if sub == nil {
			continue
		}
		elemental.ResetDefaultForZeroValues(sub)
		if err := sub.Validate(); err != nil {
			errors = errors.Append(err)
		}
	}

	for _, sub := range o.HTTPSources {
		if sub == nil {
			continue
		}
		elemental.ResetDefaultForZeroValues(sub)
		if err := sub.Validate(); err != nil {
			errors = errors.Append(err)
		}
	}

	for _, sub := range o.LDAPSources {
		if sub == nil {
			continue
		}
		elemental.ResetDefaultForZeroValues(sub)
		if err := sub.Validate(); err != nil {
			errors = errors.Append(err)
		}
	}

	for _, sub := range o.MTLSSources {
		if sub == nil {
			continue
		}
		elemental.ResetDefaultForZeroValues(sub)
		if err := sub.Validate(); err != nil {
			errors = errors.Append(err)
		}
	}

	for _, sub := range o.OIDCSources {
		if sub == nil {
			continue
		}
		elemental.ResetDefaultForZeroValues(sub)
		if err := sub.Validate(); err != nil {
			errors = errors.Append(err)
		}
	}

	for _, sub := range o.Authorizations {
		if sub == nil {
			continue
		}
		elemental.ResetDefaultForZeroValues(sub)
		if err := sub.Validate(); err != nil {
			errors = errors.Append(err)
		}
	}

	if err := elemental.ValidateRequiredString("label", o.Label); err != nil {
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
func (*Import) SpecificationForAttribute(name string) elemental.AttributeSpecification {

	if v, ok := ImportAttributesMap[name]; ok {
		return v
	}

	// We could not find it, so let's check on the lower case indexed spec map
	return ImportLowerCaseAttributesMap[name]
}

// AttributeSpecifications returns the full attribute specifications map.
func (*Import) AttributeSpecifications() map[string]elemental.AttributeSpecification {

	return ImportAttributesMap
}

// ValueForAttribute returns the value for the given attribute.
// This is a very advanced function that you should not need but in some
// very specific use cases.
func (o *Import) ValueForAttribute(name string) interface{} {

	switch name {
	case "A3SSources":
		return o.A3SSources
	case "HTTPSources":
		return o.HTTPSources
	case "LDAPSources":
		return o.LDAPSources
	case "MTLSSources":
		return o.MTLSSources
	case "OIDCSources":
		return o.OIDCSources
	case "authorizations":
		return o.Authorizations
	case "label":
		return o.Label
	}

	return nil
}

// ImportAttributesMap represents the map of attribute for Import.
var ImportAttributesMap = map[string]elemental.AttributeSpecification{
	"A3SSources": {
		AllowedChoices: []string{},
		ConvertedName:  "A3SSources",
		Description:    `A3S sources to import.`,
		Exposed:        true,
		Name:           "A3SSources",
		SubType:        "a3ssource",
		Type:           "refList",
	},
	"HTTPSources": {
		AllowedChoices: []string{},
		ConvertedName:  "HTTPSources",
		Description:    `HTTP sources to import.`,
		Exposed:        true,
		Name:           "HTTPSources",
		SubType:        "httpsource",
		Type:           "refList",
	},
	"LDAPSources": {
		AllowedChoices: []string{},
		ConvertedName:  "LDAPSources",
		Description:    `LDAP sources to import.`,
		Exposed:        true,
		Name:           "LDAPSources",
		SubType:        "ldapsource",
		Type:           "refList",
	},
	"MTLSSources": {
		AllowedChoices: []string{},
		ConvertedName:  "MTLSSources",
		Description:    `MTLS sources to import.`,
		Exposed:        true,
		Name:           "MTLSSources",
		SubType:        "mtlssource",
		Type:           "refList",
	},
	"OIDCSources": {
		AllowedChoices: []string{},
		ConvertedName:  "OIDCSources",
		Description:    `OIDC sources to import.`,
		Exposed:        true,
		Name:           "OIDCSources",
		SubType:        "oidcsource",
		Type:           "refList",
	},
	"Authorizations": {
		AllowedChoices: []string{},
		ConvertedName:  "Authorizations",
		Description:    `Authorizations to import.`,
		Exposed:        true,
		Name:           "authorizations",
		SubType:        "authorization",
		Type:           "refList",
	},
	"Label": {
		AllowedChoices: []string{},
		ConvertedName:  "Label",
		Description: `Import label that will be used to identify all the resources imported by this
resource.`,
		Exposed:  true,
		Name:     "label",
		Required: true,
		Type:     "string",
	},
}

// ImportLowerCaseAttributesMap represents the map of attribute for Import.
var ImportLowerCaseAttributesMap = map[string]elemental.AttributeSpecification{
	"a3ssources": {
		AllowedChoices: []string{},
		ConvertedName:  "A3SSources",
		Description:    `A3S sources to import.`,
		Exposed:        true,
		Name:           "A3SSources",
		SubType:        "a3ssource",
		Type:           "refList",
	},
	"httpsources": {
		AllowedChoices: []string{},
		ConvertedName:  "HTTPSources",
		Description:    `HTTP sources to import.`,
		Exposed:        true,
		Name:           "HTTPSources",
		SubType:        "httpsource",
		Type:           "refList",
	},
	"ldapsources": {
		AllowedChoices: []string{},
		ConvertedName:  "LDAPSources",
		Description:    `LDAP sources to import.`,
		Exposed:        true,
		Name:           "LDAPSources",
		SubType:        "ldapsource",
		Type:           "refList",
	},
	"mtlssources": {
		AllowedChoices: []string{},
		ConvertedName:  "MTLSSources",
		Description:    `MTLS sources to import.`,
		Exposed:        true,
		Name:           "MTLSSources",
		SubType:        "mtlssource",
		Type:           "refList",
	},
	"oidcsources": {
		AllowedChoices: []string{},
		ConvertedName:  "OIDCSources",
		Description:    `OIDC sources to import.`,
		Exposed:        true,
		Name:           "OIDCSources",
		SubType:        "oidcsource",
		Type:           "refList",
	},
	"authorizations": {
		AllowedChoices: []string{},
		ConvertedName:  "Authorizations",
		Description:    `Authorizations to import.`,
		Exposed:        true,
		Name:           "authorizations",
		SubType:        "authorization",
		Type:           "refList",
	},
	"label": {
		AllowedChoices: []string{},
		ConvertedName:  "Label",
		Description: `Import label that will be used to identify all the resources imported by this
resource.`,
		Exposed:  true,
		Name:     "label",
		Required: true,
		Type:     "string",
	},
}

// SparseImportsList represents a list of SparseImports
type SparseImportsList []*SparseImport

// Identity returns the identity of the objects in the list.
func (o SparseImportsList) Identity() elemental.Identity {

	return ImportIdentity
}

// Copy returns a pointer to a copy the SparseImportsList.
func (o SparseImportsList) Copy() elemental.Identifiables {

	copy := append(SparseImportsList{}, o...)
	return &copy
}

// Append appends the objects to the a new copy of the SparseImportsList.
func (o SparseImportsList) Append(objects ...elemental.Identifiable) elemental.Identifiables {

	out := append(SparseImportsList{}, o...)
	for _, obj := range objects {
		out = append(out, obj.(*SparseImport))
	}

	return out
}

// List converts the object to an elemental.IdentifiablesList.
func (o SparseImportsList) List() elemental.IdentifiablesList {

	out := make(elemental.IdentifiablesList, len(o))
	for i := 0; i < len(o); i++ {
		out[i] = o[i]
	}

	return out
}

// DefaultOrder returns the default ordering fields of the content.
func (o SparseImportsList) DefaultOrder() []string {

	return []string{}
}

// ToPlain returns the SparseImportsList converted to ImportsList.
func (o SparseImportsList) ToPlain() elemental.IdentifiablesList {

	out := make(elemental.IdentifiablesList, len(o))
	for i := 0; i < len(o); i++ {
		out[i] = o[i].ToPlain()
	}

	return out
}

// Version returns the version of the content.
func (o SparseImportsList) Version() int {

	return 1
}

// SparseImport represents the sparse version of a import.
type SparseImport struct {
	// A3S sources to import.
	A3SSources *A3SSourcesList `json:"A3SSources,omitempty" msgpack:"A3SSources,omitempty" bson:"-" mapstructure:"A3SSources,omitempty"`

	// HTTP sources to import.
	HTTPSources *HTTPSourcesList `json:"HTTPSources,omitempty" msgpack:"HTTPSources,omitempty" bson:"-" mapstructure:"HTTPSources,omitempty"`

	// LDAP sources to import.
	LDAPSources *LDAPSourcesList `json:"LDAPSources,omitempty" msgpack:"LDAPSources,omitempty" bson:"-" mapstructure:"LDAPSources,omitempty"`

	// MTLS sources to import.
	MTLSSources *MTLSSourcesList `json:"MTLSSources,omitempty" msgpack:"MTLSSources,omitempty" bson:"-" mapstructure:"MTLSSources,omitempty"`

	// OIDC sources to import.
	OIDCSources *OIDCSourcesList `json:"OIDCSources,omitempty" msgpack:"OIDCSources,omitempty" bson:"-" mapstructure:"OIDCSources,omitempty"`

	// Authorizations to import.
	Authorizations *AuthorizationsList `json:"authorizations,omitempty" msgpack:"authorizations,omitempty" bson:"-" mapstructure:"authorizations,omitempty"`

	// Import label that will be used to identify all the resources imported by this
	// resource.
	Label *string `json:"label,omitempty" msgpack:"label,omitempty" bson:"-" mapstructure:"label,omitempty"`

	ModelVersion int `json:"-" msgpack:"-" bson:"_modelversion"`
}

// NewSparseImport returns a new  SparseImport.
func NewSparseImport() *SparseImport {
	return &SparseImport{}
}

// Identity returns the Identity of the sparse object.
func (o *SparseImport) Identity() elemental.Identity {

	return ImportIdentity
}

// Identifier returns the value of the sparse object's unique identifier.
func (o *SparseImport) Identifier() string {

	return ""
}

// SetIdentifier sets the value of the sparse object's unique identifier.
func (o *SparseImport) SetIdentifier(id string) {

}

// GetBSON implements the bson marshaling interface.
// This is used to transparently convert ID to MongoDBID as ObectID.
func (o *SparseImport) GetBSON() (interface{}, error) {

	if o == nil {
		return nil, nil
	}

	s := &mongoAttributesSparseImport{}

	return s, nil
}

// SetBSON implements the bson marshaling interface.
// This is used to transparently convert ID to MongoDBID as ObectID.
func (o *SparseImport) SetBSON(raw bson.Raw) error {

	if o == nil {
		return nil
	}

	s := &mongoAttributesSparseImport{}
	if err := raw.Unmarshal(s); err != nil {
		return err
	}

	return nil
}

// Version returns the hardcoded version of the model.
func (o *SparseImport) Version() int {

	return 1
}

// ToPlain returns the plain version of the sparse model.
func (o *SparseImport) ToPlain() elemental.PlainIdentifiable {

	out := NewImport()
	if o.A3SSources != nil {
		out.A3SSources = *o.A3SSources
	}
	if o.HTTPSources != nil {
		out.HTTPSources = *o.HTTPSources
	}
	if o.LDAPSources != nil {
		out.LDAPSources = *o.LDAPSources
	}
	if o.MTLSSources != nil {
		out.MTLSSources = *o.MTLSSources
	}
	if o.OIDCSources != nil {
		out.OIDCSources = *o.OIDCSources
	}
	if o.Authorizations != nil {
		out.Authorizations = *o.Authorizations
	}
	if o.Label != nil {
		out.Label = *o.Label
	}

	return out
}

// DeepCopy returns a deep copy if the SparseImport.
func (o *SparseImport) DeepCopy() *SparseImport {

	if o == nil {
		return nil
	}

	out := &SparseImport{}
	o.DeepCopyInto(out)

	return out
}

// DeepCopyInto copies the receiver into the given *SparseImport.
func (o *SparseImport) DeepCopyInto(out *SparseImport) {

	target, err := copystructure.Copy(o)
	if err != nil {
		panic(fmt.Sprintf("Unable to deepcopy SparseImport: %s", err))
	}

	*out = *target.(*SparseImport)
}

type mongoAttributesImport struct {
}
type mongoAttributesSparseImport struct {
}
