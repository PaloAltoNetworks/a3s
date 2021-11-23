package api

import (
	"fmt"

	"github.com/globalsign/mgo/bson"
	"github.com/mitchellh/copystructure"
	"go.aporeto.io/elemental"
)

// IssueRemoteA3S represents the model of a issueremotea3s
type IssueRemoteA3S struct {
	// The remote a3s token.
	Token string `json:"token" msgpack:"token" bson:"-" mapstructure:"token,omitempty"`

	ModelVersion int `json:"-" msgpack:"-" bson:"_modelversion"`
}

// NewIssueRemoteA3S returns a new *IssueRemoteA3S
func NewIssueRemoteA3S() *IssueRemoteA3S {

	return &IssueRemoteA3S{
		ModelVersion: 1,
	}
}

// GetBSON implements the bson marshaling interface.
// This is used to transparently convert ID to MongoDBID as ObectID.
func (o *IssueRemoteA3S) GetBSON() (interface{}, error) {

	if o == nil {
		return nil, nil
	}

	s := &mongoAttributesIssueRemoteA3S{}

	return s, nil
}

// SetBSON implements the bson marshaling interface.
// This is used to transparently convert ID to MongoDBID as ObectID.
func (o *IssueRemoteA3S) SetBSON(raw bson.Raw) error {

	if o == nil {
		return nil
	}

	s := &mongoAttributesIssueRemoteA3S{}
	if err := raw.Unmarshal(s); err != nil {
		return err
	}

	return nil
}

// BleveType implements the bleve.Classifier Interface.
func (o *IssueRemoteA3S) BleveType() string {

	return "issueremotea3s"
}

// DeepCopy returns a deep copy if the IssueRemoteA3S.
func (o *IssueRemoteA3S) DeepCopy() *IssueRemoteA3S {

	if o == nil {
		return nil
	}

	out := &IssueRemoteA3S{}
	o.DeepCopyInto(out)

	return out
}

// DeepCopyInto copies the receiver into the given *IssueRemoteA3S.
func (o *IssueRemoteA3S) DeepCopyInto(out *IssueRemoteA3S) {

	target, err := copystructure.Copy(o)
	if err != nil {
		panic(fmt.Sprintf("Unable to deepcopy IssueRemoteA3S: %s", err))
	}

	*out = *target.(*IssueRemoteA3S)
}

// Validate valides the current information stored into the structure.
func (o *IssueRemoteA3S) Validate() error {

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

type mongoAttributesIssueRemoteA3S struct {
}
