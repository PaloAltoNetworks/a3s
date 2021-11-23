package api

import (
	"fmt"

	"github.com/globalsign/mgo/bson"
	"github.com/mitchellh/copystructure"
	"go.aporeto.io/elemental"
)

// IssueA3S represents the model of a issuea3s
type IssueA3S struct {
	// The original token.
	Token string `json:"token" msgpack:"token" bson:"-" mapstructure:"token,omitempty"`

	ModelVersion int `json:"-" msgpack:"-" bson:"_modelversion"`
}

// NewIssueA3S returns a new *IssueA3S
func NewIssueA3S() *IssueA3S {

	return &IssueA3S{
		ModelVersion: 1,
	}
}

// GetBSON implements the bson marshaling interface.
// This is used to transparently convert ID to MongoDBID as ObectID.
func (o *IssueA3S) GetBSON() (interface{}, error) {

	if o == nil {
		return nil, nil
	}

	s := &mongoAttributesIssueA3S{}

	return s, nil
}

// SetBSON implements the bson marshaling interface.
// This is used to transparently convert ID to MongoDBID as ObectID.
func (o *IssueA3S) SetBSON(raw bson.Raw) error {

	if o == nil {
		return nil
	}

	s := &mongoAttributesIssueA3S{}
	if err := raw.Unmarshal(s); err != nil {
		return err
	}

	return nil
}

// BleveType implements the bleve.Classifier Interface.
func (o *IssueA3S) BleveType() string {

	return "issuea3s"
}

// DeepCopy returns a deep copy if the IssueA3S.
func (o *IssueA3S) DeepCopy() *IssueA3S {

	if o == nil {
		return nil
	}

	out := &IssueA3S{}
	o.DeepCopyInto(out)

	return out
}

// DeepCopyInto copies the receiver into the given *IssueA3S.
func (o *IssueA3S) DeepCopyInto(out *IssueA3S) {

	target, err := copystructure.Copy(o)
	if err != nil {
		panic(fmt.Sprintf("Unable to deepcopy IssueA3S: %s", err))
	}

	*out = *target.(*IssueA3S)
}

// Validate valides the current information stored into the structure.
func (o *IssueA3S) Validate() error {

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

type mongoAttributesIssueA3S struct {
}
