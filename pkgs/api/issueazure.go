package api

import (
	"fmt"

	"github.com/globalsign/mgo/bson"
	"github.com/mitchellh/copystructure"
	"go.aporeto.io/elemental"
)

// IssueAzure represents the model of a issueazure
type IssueAzure struct {
	// The original token.
	Token string `json:"token" msgpack:"token" bson:"-" mapstructure:"token,omitempty"`

	ModelVersion int `json:"-" msgpack:"-" bson:"_modelversion"`
}

// NewIssueAzure returns a new *IssueAzure
func NewIssueAzure() *IssueAzure {

	return &IssueAzure{
		ModelVersion: 1,
	}
}

// GetBSON implements the bson marshaling interface.
// This is used to transparently convert ID to MongoDBID as ObectID.
func (o *IssueAzure) GetBSON() (interface{}, error) {

	if o == nil {
		return nil, nil
	}

	s := &mongoAttributesIssueAzure{}

	return s, nil
}

// SetBSON implements the bson marshaling interface.
// This is used to transparently convert ID to MongoDBID as ObectID.
func (o *IssueAzure) SetBSON(raw bson.Raw) error {

	if o == nil {
		return nil
	}

	s := &mongoAttributesIssueAzure{}
	if err := raw.Unmarshal(s); err != nil {
		return err
	}

	return nil
}

// BleveType implements the bleve.Classifier Interface.
func (o *IssueAzure) BleveType() string {

	return "issueazure"
}

// DeepCopy returns a deep copy if the IssueAzure.
func (o *IssueAzure) DeepCopy() *IssueAzure {

	if o == nil {
		return nil
	}

	out := &IssueAzure{}
	o.DeepCopyInto(out)

	return out
}

// DeepCopyInto copies the receiver into the given *IssueAzure.
func (o *IssueAzure) DeepCopyInto(out *IssueAzure) {

	target, err := copystructure.Copy(o)
	if err != nil {
		panic(fmt.Sprintf("Unable to deepcopy IssueAzure: %s", err))
	}

	*out = *target.(*IssueAzure)
}

// Validate valides the current information stored into the structure.
func (o *IssueAzure) Validate() error {

	errors := elemental.Errors{}
	requiredErrors := elemental.Errors{}

	if err := elemental.ValidateRequiredString("token", o.Token); err != nil {
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
func (*IssueAzure) SpecificationForAttribute(name string) elemental.AttributeSpecification {

	if v, ok := IssueAzureAttributesMap[name]; ok {
		return v
	}

	// We could not find it, so let's check on the lower case indexed spec map
	return IssueAzureLowerCaseAttributesMap[name]
}

// AttributeSpecifications returns the full attribute specifications map.
func (*IssueAzure) AttributeSpecifications() map[string]elemental.AttributeSpecification {

	return IssueAzureAttributesMap
}

// ValueForAttribute returns the value for the given attribute.
// This is a very advanced function that you should not need but in some
// very specific use cases.
func (o *IssueAzure) ValueForAttribute(name string) interface{} {

	switch name {
	case "token":
		return o.Token
	}

	return nil
}

// IssueAzureAttributesMap represents the map of attribute for IssueAzure.
var IssueAzureAttributesMap = map[string]elemental.AttributeSpecification{
	"Token": {
		AllowedChoices: []string{},
		ConvertedName:  "Token",
		Description:    `The original token.`,
		Exposed:        true,
		Name:           "token",
		Required:       true,
		Type:           "string",
	},
}

// IssueAzureLowerCaseAttributesMap represents the map of attribute for IssueAzure.
var IssueAzureLowerCaseAttributesMap = map[string]elemental.AttributeSpecification{
	"token": {
		AllowedChoices: []string{},
		ConvertedName:  "Token",
		Description:    `The original token.`,
		Exposed:        true,
		Name:           "token",
		Required:       true,
		Type:           "string",
	},
}

type mongoAttributesIssueAzure struct {
}
