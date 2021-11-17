package api

import (
	"fmt"

	"github.com/globalsign/mgo/bson"
	"github.com/mitchellh/copystructure"
	"go.aporeto.io/elemental"
)

// IssueOIDC represents the model of a issueoidc
type IssueOIDC struct {
	// OIDC ceremony code.
	Code string `json:"code" msgpack:"code" bson:"-" mapstructure:"code,omitempty"`

	// OIDC redirect url in case of error.
	RedirectErrorURL string `json:"redirectErrorURL" msgpack:"redirectErrorURL" bson:"-" mapstructure:"redirectErrorURL,omitempty"`

	// OIDC redirect url.
	RedirectURL string `json:"redirectURL" msgpack:"redirectURL" bson:"-" mapstructure:"redirectURL,omitempty"`

	// OIDC ceremony state.
	State string `json:"state" msgpack:"state" bson:"-" mapstructure:"state,omitempty"`

	ModelVersion int `json:"-" msgpack:"-" bson:"_modelversion"`
}

// NewIssueOIDC returns a new *IssueOIDC
func NewIssueOIDC() *IssueOIDC {

	return &IssueOIDC{
		ModelVersion: 1,
	}
}

// GetBSON implements the bson marshaling interface.
// This is used to transparently convert ID to MongoDBID as ObectID.
func (o *IssueOIDC) GetBSON() (interface{}, error) {

	if o == nil {
		return nil, nil
	}

	s := &mongoAttributesIssueOIDC{}

	return s, nil
}

// SetBSON implements the bson marshaling interface.
// This is used to transparently convert ID to MongoDBID as ObectID.
func (o *IssueOIDC) SetBSON(raw bson.Raw) error {

	if o == nil {
		return nil
	}

	s := &mongoAttributesIssueOIDC{}
	if err := raw.Unmarshal(s); err != nil {
		return err
	}

	return nil
}

// BleveType implements the bleve.Classifier Interface.
func (o *IssueOIDC) BleveType() string {

	return "issueoidc"
}

// DeepCopy returns a deep copy if the IssueOIDC.
func (o *IssueOIDC) DeepCopy() *IssueOIDC {

	if o == nil {
		return nil
	}

	out := &IssueOIDC{}
	o.DeepCopyInto(out)

	return out
}

// DeepCopyInto copies the receiver into the given *IssueOIDC.
func (o *IssueOIDC) DeepCopyInto(out *IssueOIDC) {

	target, err := copystructure.Copy(o)
	if err != nil {
		panic(fmt.Sprintf("Unable to deepcopy IssueOIDC: %s", err))
	}

	*out = *target.(*IssueOIDC)
}

// Validate valides the current information stored into the structure.
func (o *IssueOIDC) Validate() error {

	errors := elemental.Errors{}
	requiredErrors := elemental.Errors{}

	if len(requiredErrors) > 0 {
		return requiredErrors
	}

	if len(errors) > 0 {
		return errors
	}

	return nil
}

type mongoAttributesIssueOIDC struct {
}
