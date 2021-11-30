package api

import (
	"fmt"

	"github.com/globalsign/mgo/bson"
	"github.com/mitchellh/copystructure"
	"go.aporeto.io/elemental"
)

// IdentityModifierMethodValue represents the possible values for attribute "method".
type IdentityModifierMethodValue string

const (
	// IdentityModifierMethodGET represents the value GET.
	IdentityModifierMethodGET IdentityModifierMethodValue = "GET"

	// IdentityModifierMethodPATCH represents the value PATCH.
	IdentityModifierMethodPATCH IdentityModifierMethodValue = "PATCH"

	// IdentityModifierMethodPOST represents the value POST.
	IdentityModifierMethodPOST IdentityModifierMethodValue = "POST"

	// IdentityModifierMethodPUT represents the value PUT.
	IdentityModifierMethodPUT IdentityModifierMethodValue = "PUT"
)

// IdentityModifier represents the model of a identitymodifier
type IdentityModifier struct {
	// URL of the remote service. This URL will receive a call containing the
	// claims that are about to be delivered. It must reply with 204 if it does not
	// wish to modify the claims, or 200 alongside a body containing the modified
	// claims.
	URL string `json:"URL" msgpack:"URL" bson:"-" mapstructure:"URL,omitempty"`

	// Client certificate required to call URL. A3S will refuse to send data if the
	// endpoint does not support client certificate authentication.
	Certificate string `json:"certificate" msgpack:"certificate" bson:"-" mapstructure:"certificate,omitempty"`

	// CA to use to validate the entity serving the URL.
	CertificateAuthority string `json:"certificateAuthority,omitempty" msgpack:"certificateAuthority,omitempty" bson:"certificateauthority,omitempty" mapstructure:"certificateAuthority,omitempty"`

	// Key associated to the client certificate.
	Key string `json:"key" msgpack:"key" bson:"-" mapstructure:"key,omitempty"`

	// The HTTP method to use to call the endpoint. For POST/PUT/PATCH the remote
	// server will receive the claims as a JSON encoded array in the body. For a GET, the claims will be passed as a query parameter named `claim`.
	Method IdentityModifierMethodValue `json:"method" msgpack:"method" bson:"-" mapstructure:"method,omitempty"`

	ModelVersion int `json:"-" msgpack:"-" bson:"_modelversion"`
}

// NewIdentityModifier returns a new *IdentityModifier
func NewIdentityModifier() *IdentityModifier {

	return &IdentityModifier{
		ModelVersion: 1,
		Method:       IdentityModifierMethodPOST,
	}
}

// GetBSON implements the bson marshaling interface.
// This is used to transparently convert ID to MongoDBID as ObectID.
func (o *IdentityModifier) GetBSON() (interface{}, error) {

	if o == nil {
		return nil, nil
	}

	s := &mongoAttributesIdentityModifier{}

	s.CertificateAuthority = o.CertificateAuthority

	return s, nil
}

// SetBSON implements the bson marshaling interface.
// This is used to transparently convert ID to MongoDBID as ObectID.
func (o *IdentityModifier) SetBSON(raw bson.Raw) error {

	if o == nil {
		return nil
	}

	s := &mongoAttributesIdentityModifier{}
	if err := raw.Unmarshal(s); err != nil {
		return err
	}

	o.CertificateAuthority = s.CertificateAuthority

	return nil
}

// BleveType implements the bleve.Classifier Interface.
func (o *IdentityModifier) BleveType() string {

	return "identitymodifier"
}

// DeepCopy returns a deep copy if the IdentityModifier.
func (o *IdentityModifier) DeepCopy() *IdentityModifier {

	if o == nil {
		return nil
	}

	out := &IdentityModifier{}
	o.DeepCopyInto(out)

	return out
}

// DeepCopyInto copies the receiver into the given *IdentityModifier.
func (o *IdentityModifier) DeepCopyInto(out *IdentityModifier) {

	target, err := copystructure.Copy(o)
	if err != nil {
		panic(fmt.Sprintf("Unable to deepcopy IdentityModifier: %s", err))
	}

	*out = *target.(*IdentityModifier)
}

// Validate valides the current information stored into the structure.
func (o *IdentityModifier) Validate() error {

	errors := elemental.Errors{}
	requiredErrors := elemental.Errors{}

	if err := elemental.ValidateRequiredString("URL", o.URL); err != nil {
		requiredErrors = requiredErrors.Append(err)
	}

	if err := elemental.ValidateRequiredString("certificate", o.Certificate); err != nil {
		requiredErrors = requiredErrors.Append(err)
	}

	if err := elemental.ValidateRequiredString("key", o.Key); err != nil {
		requiredErrors = requiredErrors.Append(err)
	}

	if err := elemental.ValidateRequiredString("method", string(o.Method)); err != nil {
		requiredErrors = requiredErrors.Append(err)
	}

	if err := elemental.ValidateStringInList("method", string(o.Method), []string{"GET", "POST", "PUT", "PATCH"}, false); err != nil {
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
func (*IdentityModifier) SpecificationForAttribute(name string) elemental.AttributeSpecification {

	if v, ok := IdentityModifierAttributesMap[name]; ok {
		return v
	}

	// We could not find it, so let's check on the lower case indexed spec map
	return IdentityModifierLowerCaseAttributesMap[name]
}

// AttributeSpecifications returns the full attribute specifications map.
func (*IdentityModifier) AttributeSpecifications() map[string]elemental.AttributeSpecification {

	return IdentityModifierAttributesMap
}

// ValueForAttribute returns the value for the given attribute.
// This is a very advanced function that you should not need but in some
// very specific use cases.
func (o *IdentityModifier) ValueForAttribute(name string) interface{} {

	switch name {
	case "URL":
		return o.URL
	case "certificate":
		return o.Certificate
	case "certificateAuthority":
		return o.CertificateAuthority
	case "key":
		return o.Key
	case "method":
		return o.Method
	}

	return nil
}

// IdentityModifierAttributesMap represents the map of attribute for IdentityModifier.
var IdentityModifierAttributesMap = map[string]elemental.AttributeSpecification{
	"URL": {
		AllowedChoices: []string{},
		ConvertedName:  "URL",
		Description: `URL of the remote service. This URL will receive a call containing the
claims that are about to be delivered. It must reply with 204 if it does not
wish to modify the claims, or 200 alongside a body containing the modified
claims.`,
		Exposed:  true,
		Name:     "URL",
		Required: true,
		Type:     "string",
	},
	"Certificate": {
		AllowedChoices: []string{},
		ConvertedName:  "Certificate",
		Description: `Client certificate required to call URL. A3S will refuse to send data if the
endpoint does not support client certificate authentication.`,
		Exposed:  true,
		Name:     "certificate",
		Required: true,
		Type:     "string",
	},
	"CertificateAuthority": {
		AllowedChoices: []string{},
		BSONFieldName:  "certificateauthority",
		ConvertedName:  "CertificateAuthority",
		Description:    `CA to use to validate the entity serving the URL.`,
		Exposed:        true,
		Name:           "certificateAuthority",
		Stored:         true,
		Type:           "string",
	},
	"Key": {
		AllowedChoices: []string{},
		ConvertedName:  "Key",
		Description:    `Key associated to the client certificate.`,
		Exposed:        true,
		Name:           "key",
		Required:       true,
		Type:           "string",
	},
	"Method": {
		AllowedChoices: []string{"GET", "POST", "PUT", "PATCH"},
		ConvertedName:  "Method",
		DefaultValue:   IdentityModifierMethodPOST,
		Description: `The HTTP method to use to call the endpoint. For POST/PUT/PATCH the remote
server will receive the claims as a JSON encoded array in the body. For a GET, the claims will be passed as a query parameter named ` + "`" + `claim` + "`" + `.`,
		Exposed:  true,
		Name:     "method",
		Required: true,
		Type:     "enum",
	},
}

// IdentityModifierLowerCaseAttributesMap represents the map of attribute for IdentityModifier.
var IdentityModifierLowerCaseAttributesMap = map[string]elemental.AttributeSpecification{
	"url": {
		AllowedChoices: []string{},
		ConvertedName:  "URL",
		Description: `URL of the remote service. This URL will receive a call containing the
claims that are about to be delivered. It must reply with 204 if it does not
wish to modify the claims, or 200 alongside a body containing the modified
claims.`,
		Exposed:  true,
		Name:     "URL",
		Required: true,
		Type:     "string",
	},
	"certificate": {
		AllowedChoices: []string{},
		ConvertedName:  "Certificate",
		Description: `Client certificate required to call URL. A3S will refuse to send data if the
endpoint does not support client certificate authentication.`,
		Exposed:  true,
		Name:     "certificate",
		Required: true,
		Type:     "string",
	},
	"certificateauthority": {
		AllowedChoices: []string{},
		BSONFieldName:  "certificateauthority",
		ConvertedName:  "CertificateAuthority",
		Description:    `CA to use to validate the entity serving the URL.`,
		Exposed:        true,
		Name:           "certificateAuthority",
		Stored:         true,
		Type:           "string",
	},
	"key": {
		AllowedChoices: []string{},
		ConvertedName:  "Key",
		Description:    `Key associated to the client certificate.`,
		Exposed:        true,
		Name:           "key",
		Required:       true,
		Type:           "string",
	},
	"method": {
		AllowedChoices: []string{"GET", "POST", "PUT", "PATCH"},
		ConvertedName:  "Method",
		DefaultValue:   IdentityModifierMethodPOST,
		Description: `The HTTP method to use to call the endpoint. For POST/PUT/PATCH the remote
server will receive the claims as a JSON encoded array in the body. For a GET, the claims will be passed as a query parameter named ` + "`" + `claim` + "`" + `.`,
		Exposed:  true,
		Name:     "method",
		Required: true,
		Type:     "enum",
	},
}

type mongoAttributesIdentityModifier struct {
	CertificateAuthority string `bson:"certificateauthority,omitempty"`
}
