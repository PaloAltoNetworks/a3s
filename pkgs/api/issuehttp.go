package api

import (
	"fmt"

	"github.com/globalsign/mgo/bson"
	"github.com/mitchellh/copystructure"
	"go.aporeto.io/elemental"
)

// IssueHTTP represents the model of a issuehttp
type IssueHTTP struct {
	// The password for the user.
	Password string `json:"password" msgpack:"password" bson:"-" mapstructure:"password,omitempty"`

	// The username.
	Username string `json:"username" msgpack:"username" bson:"-" mapstructure:"username,omitempty"`

	ModelVersion int `json:"-" msgpack:"-" bson:"_modelversion"`
}

// NewIssueHTTP returns a new *IssueHTTP
func NewIssueHTTP() *IssueHTTP {

	return &IssueHTTP{
		ModelVersion: 1,
	}
}

// GetBSON implements the bson marshaling interface.
// This is used to transparently convert ID to MongoDBID as ObectID.
func (o *IssueHTTP) GetBSON() (interface{}, error) {

	if o == nil {
		return nil, nil
	}

	s := &mongoAttributesIssueHTTP{}

	return s, nil
}

// SetBSON implements the bson marshaling interface.
// This is used to transparently convert ID to MongoDBID as ObectID.
func (o *IssueHTTP) SetBSON(raw bson.Raw) error {

	if o == nil {
		return nil
	}

	s := &mongoAttributesIssueHTTP{}
	if err := raw.Unmarshal(s); err != nil {
		return err
	}

	return nil
}

// BleveType implements the bleve.Classifier Interface.
func (o *IssueHTTP) BleveType() string {

	return "issuehttp"
}

// DeepCopy returns a deep copy if the IssueHTTP.
func (o *IssueHTTP) DeepCopy() *IssueHTTP {

	if o == nil {
		return nil
	}

	out := &IssueHTTP{}
	o.DeepCopyInto(out)

	return out
}

// DeepCopyInto copies the receiver into the given *IssueHTTP.
func (o *IssueHTTP) DeepCopyInto(out *IssueHTTP) {

	target, err := copystructure.Copy(o)
	if err != nil {
		panic(fmt.Sprintf("Unable to deepcopy IssueHTTP: %s", err))
	}

	*out = *target.(*IssueHTTP)
}

// Validate valides the current information stored into the structure.
func (o *IssueHTTP) Validate() error {

	errors := elemental.Errors{}
	requiredErrors := elemental.Errors{}

	if err := elemental.ValidateRequiredString("password", o.Password); err != nil {
		requiredErrors = requiredErrors.Append(err)
	}

	if err := elemental.ValidateRequiredString("username", o.Username); err != nil {
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
func (*IssueHTTP) SpecificationForAttribute(name string) elemental.AttributeSpecification {

	if v, ok := IssueHTTPAttributesMap[name]; ok {
		return v
	}

	// We could not find it, so let's check on the lower case indexed spec map
	return IssueHTTPLowerCaseAttributesMap[name]
}

// AttributeSpecifications returns the full attribute specifications map.
func (*IssueHTTP) AttributeSpecifications() map[string]elemental.AttributeSpecification {

	return IssueHTTPAttributesMap
}

// ValueForAttribute returns the value for the given attribute.
// This is a very advanced function that you should not need but in some
// very specific use cases.
func (o *IssueHTTP) ValueForAttribute(name string) interface{} {

	switch name {
	case "password":
		return o.Password
	case "username":
		return o.Username
	}

	return nil
}

// IssueHTTPAttributesMap represents the map of attribute for IssueHTTP.
var IssueHTTPAttributesMap = map[string]elemental.AttributeSpecification{
	"Password": {
		AllowedChoices: []string{},
		ConvertedName:  "Password",
		Description:    `The password for the user.`,
		Exposed:        true,
		Name:           "password",
		Required:       true,
		Type:           "string",
	},
	"Username": {
		AllowedChoices: []string{},
		ConvertedName:  "Username",
		Description:    `The username.`,
		Exposed:        true,
		Name:           "username",
		Required:       true,
		Type:           "string",
	},
}

// IssueHTTPLowerCaseAttributesMap represents the map of attribute for IssueHTTP.
var IssueHTTPLowerCaseAttributesMap = map[string]elemental.AttributeSpecification{
	"password": {
		AllowedChoices: []string{},
		ConvertedName:  "Password",
		Description:    `The password for the user.`,
		Exposed:        true,
		Name:           "password",
		Required:       true,
		Type:           "string",
	},
	"username": {
		AllowedChoices: []string{},
		ConvertedName:  "Username",
		Description:    `The username.`,
		Exposed:        true,
		Name:           "username",
		Required:       true,
		Type:           "string",
	},
}

type mongoAttributesIssueHTTP struct {
}
