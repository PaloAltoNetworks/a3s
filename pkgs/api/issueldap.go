package api

import (
	"fmt"

	"github.com/globalsign/mgo/bson"
	"github.com/mitchellh/copystructure"
	"go.aporeto.io/elemental"
)

// IssueLDAP represents the model of a issueldap
type IssueLDAP struct {
	// The password for the user.
	Password string `json:"password" msgpack:"password" bson:"-" mapstructure:"password,omitempty"`

	// The LDAP username.
	Username string `json:"username" msgpack:"username" bson:"-" mapstructure:"username,omitempty"`

	ModelVersion int `json:"-" msgpack:"-" bson:"_modelversion"`
}

// NewIssueLDAP returns a new *IssueLDAP
func NewIssueLDAP() *IssueLDAP {

	return &IssueLDAP{
		ModelVersion: 1,
	}
}

// GetBSON implements the bson marshaling interface.
// This is used to transparently convert ID to MongoDBID as ObectID.
func (o *IssueLDAP) GetBSON() (interface{}, error) {

	if o == nil {
		return nil, nil
	}

	s := &mongoAttributesIssueLDAP{}

	return s, nil
}

// SetBSON implements the bson marshaling interface.
// This is used to transparently convert ID to MongoDBID as ObectID.
func (o *IssueLDAP) SetBSON(raw bson.Raw) error {

	if o == nil {
		return nil
	}

	s := &mongoAttributesIssueLDAP{}
	if err := raw.Unmarshal(s); err != nil {
		return err
	}

	return nil
}

// BleveType implements the bleve.Classifier Interface.
func (o *IssueLDAP) BleveType() string {

	return "issueldap"
}

// DeepCopy returns a deep copy if the IssueLDAP.
func (o *IssueLDAP) DeepCopy() *IssueLDAP {

	if o == nil {
		return nil
	}

	out := &IssueLDAP{}
	o.DeepCopyInto(out)

	return out
}

// DeepCopyInto copies the receiver into the given *IssueLDAP.
func (o *IssueLDAP) DeepCopyInto(out *IssueLDAP) {

	target, err := copystructure.Copy(o)
	if err != nil {
		panic(fmt.Sprintf("Unable to deepcopy IssueLDAP: %s", err))
	}

	*out = *target.(*IssueLDAP)
}

// Validate valides the current information stored into the structure.
func (o *IssueLDAP) Validate() error {

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

type mongoAttributesIssueLDAP struct {
}
