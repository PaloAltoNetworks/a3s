package api

import (
	"fmt"

	"github.com/globalsign/mgo/bson"
	"github.com/mitchellh/copystructure"
	"go.aporeto.io/elemental"
)

// IssueSourceTypeValue represents the possible values for attribute "sourceType".
type IssueSourceTypeValue string

const (
	// IssueSourceTypeA3S represents the value A3S.
	IssueSourceTypeA3S IssueSourceTypeValue = "A3S"

	// IssueSourceTypeAWS represents the value AWS.
	IssueSourceTypeAWS IssueSourceTypeValue = "AWS"

	// IssueSourceTypeAzure represents the value Azure.
	IssueSourceTypeAzure IssueSourceTypeValue = "Azure"

	// IssueSourceTypeGCP represents the value GCP.
	IssueSourceTypeGCP IssueSourceTypeValue = "GCP"

	// IssueSourceTypeHTTP represents the value HTTP.
	IssueSourceTypeHTTP IssueSourceTypeValue = "HTTP"

	// IssueSourceTypeLDAP represents the value LDAP.
	IssueSourceTypeLDAP IssueSourceTypeValue = "LDAP"

	// IssueSourceTypeMTLS represents the value MTLS.
	IssueSourceTypeMTLS IssueSourceTypeValue = "MTLS"

	// IssueSourceTypeOIDC represents the value OIDC.
	IssueSourceTypeOIDC IssueSourceTypeValue = "OIDC"

	// IssueSourceTypeRemoteA3S represents the value RemoteA3S.
	IssueSourceTypeRemoteA3S IssueSourceTypeValue = "RemoteA3S"

	// IssueSourceTypeSAML represents the value SAML.
	IssueSourceTypeSAML IssueSourceTypeValue = "SAML"
)

// IssueTokenTypeValue represents the possible values for attribute "tokenType".
type IssueTokenTypeValue string

const (
	// IssueTokenTypeIdentity represents the value Identity.
	IssueTokenTypeIdentity IssueTokenTypeValue = "Identity"

	// IssueTokenTypeRefresh represents the value Refresh.
	IssueTokenTypeRefresh IssueTokenTypeValue = "Refresh"
)

// IssueIdentity represents the Identity of the object.
var IssueIdentity = elemental.Identity{
	Name:     "issue",
	Category: "issue",
	Package:  "authn",
	Private:  false,
}

// IssuesList represents a list of Issues
type IssuesList []*Issue

// Identity returns the identity of the objects in the list.
func (o IssuesList) Identity() elemental.Identity {

	return IssueIdentity
}

// Copy returns a pointer to a copy the IssuesList.
func (o IssuesList) Copy() elemental.Identifiables {

	copy := append(IssuesList{}, o...)
	return &copy
}

// Append appends the objects to the a new copy of the IssuesList.
func (o IssuesList) Append(objects ...elemental.Identifiable) elemental.Identifiables {

	out := append(IssuesList{}, o...)
	for _, obj := range objects {
		out = append(out, obj.(*Issue))
	}

	return out
}

// List converts the object to an elemental.IdentifiablesList.
func (o IssuesList) List() elemental.IdentifiablesList {

	out := make(elemental.IdentifiablesList, len(o))
	for i := 0; i < len(o); i++ {
		out[i] = o[i]
	}

	return out
}

// DefaultOrder returns the default ordering fields of the content.
func (o IssuesList) DefaultOrder() []string {

	return []string{}
}

// ToSparse returns the IssuesList converted to SparseIssuesList.
// Objects in the list will only contain the given fields. No field means entire field set.
func (o IssuesList) ToSparse(fields ...string) elemental.Identifiables {

	out := make(SparseIssuesList, len(o))
	for i := 0; i < len(o); i++ {
		out[i] = o[i].ToSparse(fields...).(*SparseIssue)
	}

	return out
}

// Version returns the version of the content.
func (o IssuesList) Version() int {

	return 1
}

// Issue represents the model of a issue
type Issue struct {
	// Requested audience for the delivered token.
	Audience []string `json:"audience,omitempty" msgpack:"audience,omitempty" bson:"-" mapstructure:"audience,omitempty"`

	// Sets a list of identity claim prefix to allow in the final token. This can be
	// used to hide some information when asking for a token as not all systems need to
	// know all of the claims.
	Cloak []string `json:"cloak,omitempty" msgpack:"cloak,omitempty" bson:"-" mapstructure:"cloak,omitempty"`

	// If set, return the token as a secure cookie.
	Cookie bool `json:"cookie,omitempty" msgpack:"cookie,omitempty" bson:"-" mapstructure:"cookie,omitempty"`

	// If set, use the provided domain for the delivered cookie.
	CookieDomain string `json:"cookieDomain,omitempty" msgpack:"cookieDomain,omitempty" bson:"-" mapstructure:"cookieDomain,omitempty"`

	// Contains additional information for an A3S token source.
	InputA3S *IssueA3S `json:"inputA3S,omitempty" msgpack:"inputA3S,omitempty" bson:"-" mapstructure:"inputA3S,omitempty"`

	// Contains additional information for an AWS STS token source.
	InputAWS *IssueAWS `json:"inputAWS,omitempty" msgpack:"inputAWS,omitempty" bson:"-" mapstructure:"inputAWS,omitempty"`

	// Contains additional information for an Azure token source.
	InputAzure *IssueAzure `json:"inputAzure,omitempty" msgpack:"inputAzure,omitempty" bson:"-" mapstructure:"inputAzure,omitempty"`

	// Contains additional information for an GCP token source.
	InputGCP *IssueGCP `json:"inputGCP,omitempty" msgpack:"inputGCP,omitempty" bson:"-" mapstructure:"inputGCP,omitempty"`

	// Contains additional information for an HTTP source.
	InputHTTP *IssueHTTP `json:"inputHTTP,omitempty" msgpack:"inputHTTP,omitempty" bson:"-" mapstructure:"inputHTTP,omitempty"`

	// Contains additional information for an LDAP source.
	InputLDAP *IssueLDAP `json:"inputLDAP,omitempty" msgpack:"inputLDAP,omitempty" bson:"-" mapstructure:"inputLDAP,omitempty"`

	// Contains additional information for an OIDC source.
	InputOIDC *IssueOIDC `json:"inputOIDC,omitempty" msgpack:"inputOIDC,omitempty" bson:"-" mapstructure:"inputOIDC,omitempty"`

	// Contains additional information for a remote A3S token source.
	InputRemoteA3S *IssueRemoteA3S `json:"inputRemoteA3S,omitempty" msgpack:"inputRemoteA3S,omitempty" bson:"-" mapstructure:"inputRemoteA3S,omitempty"`

	// Opaque data that will be included in the issued token.
	Opaque map[string]string `json:"opaque,omitempty" msgpack:"opaque,omitempty" bson:"-" mapstructure:"opaque,omitempty"`

	// Restricts the namespace where the token can be used.
	//
	// For instance, if you have have access to `/namespace` and below, you can
	// tell the policy engine that it should restrict further more to
	// `/namespace/child`.
	//
	// Restricting to a namespace you don't have initially access according to the
	// policy engine has no effect and may end up making the token unusable.
	RestrictedNamespace string `json:"restrictedNamespace,omitempty" msgpack:"restrictedNamespace,omitempty" bson:"-" mapstructure:"restrictedNamespace,omitempty"`

	// Restricts the networks from where the token can be used. This will reduce the
	// existing set of authorized networks that normally apply to the token according
	// to the policy engine.
	//
	// For instance, If you have authorized access from `0.0.0.0/0` (by default) or
	// from
	// `10.0.0.0/8`, you can ask for a token that will only be valid if used from
	// `10.1.0.0/16`.
	//
	// Restricting to a network that is not initially authorized by the policy
	// engine has no effect and may end up making the token unusable.
	RestrictedNetworks []string `json:"restrictedNetworks,omitempty" msgpack:"restrictedNetworks,omitempty" bson:"-" mapstructure:"restrictedNetworks,omitempty"`

	// Restricts the permissions of token. This will reduce the existing permissions
	// that normally apply to the token according to the policy engine.
	//
	// For instance, if you have administrative role, you can ask for a token that will
	// tell the policy engine to reduce the permission it would have granted to what is
	// given defined in the token.
	//
	// Restricting to some permissions you don't initially have according to the policy
	// engine has no effect and may end up making the token unusable.
	RestrictedPermissions []string `json:"restrictedPermissions,omitempty" msgpack:"restrictedPermissions,omitempty" bson:"-" mapstructure:"restrictedPermissions,omitempty"`

	// The name of the source to use.
	SourceName string `json:"sourceName,omitempty" msgpack:"sourceName,omitempty" bson:"-" mapstructure:"sourceName,omitempty"`

	// The namespace of the source to use.
	SourceNamespace string `json:"sourceNamespace,omitempty" msgpack:"sourceNamespace,omitempty" bson:"-" mapstructure:"sourceNamespace,omitempty"`

	// The authentication source. This will define how to verify
	// credentials from internal or external source of authentication.
	SourceType IssueSourceTypeValue `json:"sourceType" msgpack:"sourceType" bson:"-" mapstructure:"sourceType,omitempty"`

	// Issued token.
	Token string `json:"token,omitempty" msgpack:"token,omitempty" bson:"-" mapstructure:"token,omitempty"`

	// The type of token to issue.
	TokenType IssueTokenTypeValue `json:"tokenType,omitempty" msgpack:"tokenType,omitempty" bson:"-" mapstructure:"tokenType,omitempty"`

	// Configures the maximum length of validity for a token, using
	// [Golang duration syntax](https://golang.org/pkg/time/#example_Duration). If it
	// is bigger than the configured max validity, it will be capped. Default: `24h`.
	Validity string `json:"validity,omitempty" msgpack:"validity,omitempty" bson:"-" mapstructure:"validity,omitempty"`

	ModelVersion int `json:"-" msgpack:"-" bson:"_modelversion"`
}

// NewIssue returns a new *Issue
func NewIssue() *Issue {

	return &Issue{
		ModelVersion:          1,
		Audience:              []string{},
		Opaque:                map[string]string{},
		RestrictedNetworks:    []string{},
		RestrictedPermissions: []string{},
		Cloak:                 []string{},
		TokenType:             IssueTokenTypeIdentity,
		Validity:              "24h",
	}
}

// Identity returns the Identity of the object.
func (o *Issue) Identity() elemental.Identity {

	return IssueIdentity
}

// Identifier returns the value of the object's unique identifier.
func (o *Issue) Identifier() string {

	return ""
}

// SetIdentifier sets the value of the object's unique identifier.
func (o *Issue) SetIdentifier(id string) {

}

// GetBSON implements the bson marshaling interface.
// This is used to transparently convert ID to MongoDBID as ObectID.
func (o *Issue) GetBSON() (interface{}, error) {

	if o == nil {
		return nil, nil
	}

	s := &mongoAttributesIssue{}

	return s, nil
}

// SetBSON implements the bson marshaling interface.
// This is used to transparently convert ID to MongoDBID as ObectID.
func (o *Issue) SetBSON(raw bson.Raw) error {

	if o == nil {
		return nil
	}

	s := &mongoAttributesIssue{}
	if err := raw.Unmarshal(s); err != nil {
		return err
	}

	return nil
}

// Version returns the hardcoded version of the model.
func (o *Issue) Version() int {

	return 1
}

// BleveType implements the bleve.Classifier Interface.
func (o *Issue) BleveType() string {

	return "issue"
}

// DefaultOrder returns the list of default ordering fields.
func (o *Issue) DefaultOrder() []string {

	return []string{}
}

// Doc returns the documentation for the object
func (o *Issue) Doc() string {

	return `Issues a new a normalized token using various authentication sources.`
}

func (o *Issue) String() string {

	return fmt.Sprintf("<%s:%s>", o.Identity().Name, o.Identifier())
}

// ToSparse returns the sparse version of the model.
// The returned object will only contain the given fields. No field means entire field set.
func (o *Issue) ToSparse(fields ...string) elemental.SparseIdentifiable {

	if len(fields) == 0 {
		// nolint: goimports
		return &SparseIssue{
			Audience:              &o.Audience,
			Cloak:                 &o.Cloak,
			Cookie:                &o.Cookie,
			CookieDomain:          &o.CookieDomain,
			InputA3S:              o.InputA3S,
			InputAWS:              o.InputAWS,
			InputAzure:            o.InputAzure,
			InputGCP:              o.InputGCP,
			InputHTTP:             o.InputHTTP,
			InputLDAP:             o.InputLDAP,
			InputOIDC:             o.InputOIDC,
			InputRemoteA3S:        o.InputRemoteA3S,
			Opaque:                &o.Opaque,
			RestrictedNamespace:   &o.RestrictedNamespace,
			RestrictedNetworks:    &o.RestrictedNetworks,
			RestrictedPermissions: &o.RestrictedPermissions,
			SourceName:            &o.SourceName,
			SourceNamespace:       &o.SourceNamespace,
			SourceType:            &o.SourceType,
			Token:                 &o.Token,
			TokenType:             &o.TokenType,
			Validity:              &o.Validity,
		}
	}

	sp := &SparseIssue{}
	for _, f := range fields {
		switch f {
		case "audience":
			sp.Audience = &(o.Audience)
		case "cloak":
			sp.Cloak = &(o.Cloak)
		case "cookie":
			sp.Cookie = &(o.Cookie)
		case "cookieDomain":
			sp.CookieDomain = &(o.CookieDomain)
		case "inputA3S":
			sp.InputA3S = o.InputA3S
		case "inputAWS":
			sp.InputAWS = o.InputAWS
		case "inputAzure":
			sp.InputAzure = o.InputAzure
		case "inputGCP":
			sp.InputGCP = o.InputGCP
		case "inputHTTP":
			sp.InputHTTP = o.InputHTTP
		case "inputLDAP":
			sp.InputLDAP = o.InputLDAP
		case "inputOIDC":
			sp.InputOIDC = o.InputOIDC
		case "inputRemoteA3S":
			sp.InputRemoteA3S = o.InputRemoteA3S
		case "opaque":
			sp.Opaque = &(o.Opaque)
		case "restrictedNamespace":
			sp.RestrictedNamespace = &(o.RestrictedNamespace)
		case "restrictedNetworks":
			sp.RestrictedNetworks = &(o.RestrictedNetworks)
		case "restrictedPermissions":
			sp.RestrictedPermissions = &(o.RestrictedPermissions)
		case "sourceName":
			sp.SourceName = &(o.SourceName)
		case "sourceNamespace":
			sp.SourceNamespace = &(o.SourceNamespace)
		case "sourceType":
			sp.SourceType = &(o.SourceType)
		case "token":
			sp.Token = &(o.Token)
		case "tokenType":
			sp.TokenType = &(o.TokenType)
		case "validity":
			sp.Validity = &(o.Validity)
		}
	}

	return sp
}

// Patch apply the non nil value of a *SparseIssue to the object.
func (o *Issue) Patch(sparse elemental.SparseIdentifiable) {
	if !sparse.Identity().IsEqual(o.Identity()) {
		panic("cannot patch from a parse with different identity")
	}

	so := sparse.(*SparseIssue)
	if so.Audience != nil {
		o.Audience = *so.Audience
	}
	if so.Cloak != nil {
		o.Cloak = *so.Cloak
	}
	if so.Cookie != nil {
		o.Cookie = *so.Cookie
	}
	if so.CookieDomain != nil {
		o.CookieDomain = *so.CookieDomain
	}
	if so.InputA3S != nil {
		o.InputA3S = so.InputA3S
	}
	if so.InputAWS != nil {
		o.InputAWS = so.InputAWS
	}
	if so.InputAzure != nil {
		o.InputAzure = so.InputAzure
	}
	if so.InputGCP != nil {
		o.InputGCP = so.InputGCP
	}
	if so.InputHTTP != nil {
		o.InputHTTP = so.InputHTTP
	}
	if so.InputLDAP != nil {
		o.InputLDAP = so.InputLDAP
	}
	if so.InputOIDC != nil {
		o.InputOIDC = so.InputOIDC
	}
	if so.InputRemoteA3S != nil {
		o.InputRemoteA3S = so.InputRemoteA3S
	}
	if so.Opaque != nil {
		o.Opaque = *so.Opaque
	}
	if so.RestrictedNamespace != nil {
		o.RestrictedNamespace = *so.RestrictedNamespace
	}
	if so.RestrictedNetworks != nil {
		o.RestrictedNetworks = *so.RestrictedNetworks
	}
	if so.RestrictedPermissions != nil {
		o.RestrictedPermissions = *so.RestrictedPermissions
	}
	if so.SourceName != nil {
		o.SourceName = *so.SourceName
	}
	if so.SourceNamespace != nil {
		o.SourceNamespace = *so.SourceNamespace
	}
	if so.SourceType != nil {
		o.SourceType = *so.SourceType
	}
	if so.Token != nil {
		o.Token = *so.Token
	}
	if so.TokenType != nil {
		o.TokenType = *so.TokenType
	}
	if so.Validity != nil {
		o.Validity = *so.Validity
	}
}

// DeepCopy returns a deep copy if the Issue.
func (o *Issue) DeepCopy() *Issue {

	if o == nil {
		return nil
	}

	out := &Issue{}
	o.DeepCopyInto(out)

	return out
}

// DeepCopyInto copies the receiver into the given *Issue.
func (o *Issue) DeepCopyInto(out *Issue) {

	target, err := copystructure.Copy(o)
	if err != nil {
		panic(fmt.Sprintf("Unable to deepcopy Issue: %s", err))
	}

	*out = *target.(*Issue)
}

// Validate valides the current information stored into the structure.
func (o *Issue) Validate() error {

	errors := elemental.Errors{}
	requiredErrors := elemental.Errors{}

	if o.InputA3S != nil {
		elemental.ResetDefaultForZeroValues(o.InputA3S)
		if err := o.InputA3S.Validate(); err != nil {
			errors = errors.Append(err)
		}
	}

	if o.InputAWS != nil {
		elemental.ResetDefaultForZeroValues(o.InputAWS)
		if err := o.InputAWS.Validate(); err != nil {
			errors = errors.Append(err)
		}
	}

	if o.InputAzure != nil {
		elemental.ResetDefaultForZeroValues(o.InputAzure)
		if err := o.InputAzure.Validate(); err != nil {
			errors = errors.Append(err)
		}
	}

	if o.InputGCP != nil {
		elemental.ResetDefaultForZeroValues(o.InputGCP)
		if err := o.InputGCP.Validate(); err != nil {
			errors = errors.Append(err)
		}
	}

	if o.InputHTTP != nil {
		elemental.ResetDefaultForZeroValues(o.InputHTTP)
		if err := o.InputHTTP.Validate(); err != nil {
			errors = errors.Append(err)
		}
	}

	if o.InputLDAP != nil {
		elemental.ResetDefaultForZeroValues(o.InputLDAP)
		if err := o.InputLDAP.Validate(); err != nil {
			errors = errors.Append(err)
		}
	}

	if o.InputOIDC != nil {
		elemental.ResetDefaultForZeroValues(o.InputOIDC)
		if err := o.InputOIDC.Validate(); err != nil {
			errors = errors.Append(err)
		}
	}

	if o.InputRemoteA3S != nil {
		elemental.ResetDefaultForZeroValues(o.InputRemoteA3S)
		if err := o.InputRemoteA3S.Validate(); err != nil {
			errors = errors.Append(err)
		}
	}

	if err := ValidateCIDRListOptional("restrictedNetworks", o.RestrictedNetworks); err != nil {
		errors = errors.Append(err)
	}

	if err := elemental.ValidateRequiredString("sourceType", string(o.SourceType)); err != nil {
		requiredErrors = requiredErrors.Append(err)
	}

	if err := elemental.ValidateStringInList("sourceType", string(o.SourceType), []string{"A3S", "AWS", "Azure", "GCP", "HTTP", "LDAP", "MTLS", "OIDC", "RemoteA3S", "SAML"}, false); err != nil {
		errors = errors.Append(err)
	}

	if err := elemental.ValidateStringInList("tokenType", string(o.TokenType), []string{"Identity", "Refresh"}, false); err != nil {
		errors = errors.Append(err)
	}

	if err := ValidateDuration("validity", o.Validity); err != nil {
		errors = errors.Append(err)
	}

	// Custom object validation.
	if err := ValidateIssue(o); err != nil {
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
func (*Issue) SpecificationForAttribute(name string) elemental.AttributeSpecification {

	if v, ok := IssueAttributesMap[name]; ok {
		return v
	}

	// We could not find it, so let's check on the lower case indexed spec map
	return IssueLowerCaseAttributesMap[name]
}

// AttributeSpecifications returns the full attribute specifications map.
func (*Issue) AttributeSpecifications() map[string]elemental.AttributeSpecification {

	return IssueAttributesMap
}

// ValueForAttribute returns the value for the given attribute.
// This is a very advanced function that you should not need but in some
// very specific use cases.
func (o *Issue) ValueForAttribute(name string) interface{} {

	switch name {
	case "audience":
		return o.Audience
	case "cloak":
		return o.Cloak
	case "cookie":
		return o.Cookie
	case "cookieDomain":
		return o.CookieDomain
	case "inputA3S":
		return o.InputA3S
	case "inputAWS":
		return o.InputAWS
	case "inputAzure":
		return o.InputAzure
	case "inputGCP":
		return o.InputGCP
	case "inputHTTP":
		return o.InputHTTP
	case "inputLDAP":
		return o.InputLDAP
	case "inputOIDC":
		return o.InputOIDC
	case "inputRemoteA3S":
		return o.InputRemoteA3S
	case "opaque":
		return o.Opaque
	case "restrictedNamespace":
		return o.RestrictedNamespace
	case "restrictedNetworks":
		return o.RestrictedNetworks
	case "restrictedPermissions":
		return o.RestrictedPermissions
	case "sourceName":
		return o.SourceName
	case "sourceNamespace":
		return o.SourceNamespace
	case "sourceType":
		return o.SourceType
	case "token":
		return o.Token
	case "tokenType":
		return o.TokenType
	case "validity":
		return o.Validity
	}

	return nil
}

// IssueAttributesMap represents the map of attribute for Issue.
var IssueAttributesMap = map[string]elemental.AttributeSpecification{
	"Audience": {
		AllowedChoices: []string{},
		ConvertedName:  "Audience",
		Description:    `Requested audience for the delivered token.`,
		Exposed:        true,
		Name:           "audience",
		SubType:        "string",
		Type:           "list",
	},
	"Cloak": {
		AllowedChoices: []string{},
		ConvertedName:  "Cloak",
		Description: `Sets a list of identity claim prefix to allow in the final token. This can be
used to hide some information when asking for a token as not all systems need to
know all of the claims.`,
		Exposed: true,
		Name:    "cloak",
		SubType: "string",
		Type:    "list",
	},
	"Cookie": {
		AllowedChoices: []string{},
		ConvertedName:  "Cookie",
		Description:    `If set, return the token as a secure cookie.`,
		Exposed:        true,
		Name:           "cookie",
		Type:           "boolean",
	},
	"CookieDomain": {
		AllowedChoices: []string{},
		ConvertedName:  "CookieDomain",
		Description:    `If set, use the provided domain for the delivered cookie.`,
		Exposed:        true,
		Name:           "cookieDomain",
		Type:           "string",
	},
	"InputA3S": {
		AllowedChoices: []string{},
		ConvertedName:  "InputA3S",
		Description:    `Contains additional information for an A3S token source.`,
		Exposed:        true,
		Name:           "inputA3S",
		SubType:        "issuea3s",
		Type:           "ref",
	},
	"InputAWS": {
		AllowedChoices: []string{},
		ConvertedName:  "InputAWS",
		Description:    `Contains additional information for an AWS STS token source.`,
		Exposed:        true,
		Name:           "inputAWS",
		SubType:        "issueaws",
		Type:           "ref",
	},
	"InputAzure": {
		AllowedChoices: []string{},
		ConvertedName:  "InputAzure",
		Description:    `Contains additional information for an Azure token source.`,
		Exposed:        true,
		Name:           "inputAzure",
		SubType:        "issueazure",
		Type:           "ref",
	},
	"InputGCP": {
		AllowedChoices: []string{},
		ConvertedName:  "InputGCP",
		Description:    `Contains additional information for an GCP token source.`,
		Exposed:        true,
		Name:           "inputGCP",
		SubType:        "issuegcp",
		Type:           "ref",
	},
	"InputHTTP": {
		AllowedChoices: []string{},
		ConvertedName:  "InputHTTP",
		Description:    `Contains additional information for an HTTP source.`,
		Exposed:        true,
		Name:           "inputHTTP",
		SubType:        "issuehttp",
		Type:           "ref",
	},
	"InputLDAP": {
		AllowedChoices: []string{},
		ConvertedName:  "InputLDAP",
		Description:    `Contains additional information for an LDAP source.`,
		Exposed:        true,
		Name:           "inputLDAP",
		SubType:        "issueldap",
		Type:           "ref",
	},
	"InputOIDC": {
		AllowedChoices: []string{},
		ConvertedName:  "InputOIDC",
		Description:    `Contains additional information for an OIDC source.`,
		Exposed:        true,
		Name:           "inputOIDC",
		SubType:        "issueoidc",
		Type:           "ref",
	},
	"InputRemoteA3S": {
		AllowedChoices: []string{},
		ConvertedName:  "InputRemoteA3S",
		Description:    `Contains additional information for a remote A3S token source.`,
		Exposed:        true,
		Name:           "inputRemoteA3S",
		SubType:        "issueremotea3s",
		Type:           "ref",
	},
	"Opaque": {
		AllowedChoices: []string{},
		ConvertedName:  "Opaque",
		Description:    `Opaque data that will be included in the issued token.`,
		Exposed:        true,
		Name:           "opaque",
		SubType:        "map[string]string",
		Type:           "external",
	},
	"RestrictedNamespace": {
		AllowedChoices: []string{},
		ConvertedName:  "RestrictedNamespace",
		Description: `Restricts the namespace where the token can be used.

For instance, if you have have access to ` + "`" + `/namespace` + "`" + ` and below, you can
tell the policy engine that it should restrict further more to
` + "`" + `/namespace/child` + "`" + `.

Restricting to a namespace you don't have initially access according to the
policy engine has no effect and may end up making the token unusable.`,
		Exposed: true,
		Name:    "restrictedNamespace",
		Type:    "string",
	},
	"RestrictedNetworks": {
		AllowedChoices: []string{},
		ConvertedName:  "RestrictedNetworks",
		Description: `Restricts the networks from where the token can be used. This will reduce the
existing set of authorized networks that normally apply to the token according
to the policy engine.

For instance, If you have authorized access from ` + "`" + `0.0.0.0/0` + "`" + ` (by default) or
from
` + "`" + `10.0.0.0/8` + "`" + `, you can ask for a token that will only be valid if used from
` + "`" + `10.1.0.0/16` + "`" + `.

Restricting to a network that is not initially authorized by the policy
engine has no effect and may end up making the token unusable.`,
		Exposed: true,
		Name:    "restrictedNetworks",
		SubType: "string",
		Type:    "list",
	},
	"RestrictedPermissions": {
		AllowedChoices: []string{},
		ConvertedName:  "RestrictedPermissions",
		Description: `Restricts the permissions of token. This will reduce the existing permissions
that normally apply to the token according to the policy engine.

For instance, if you have administrative role, you can ask for a token that will
tell the policy engine to reduce the permission it would have granted to what is
given defined in the token.

Restricting to some permissions you don't initially have according to the policy
engine has no effect and may end up making the token unusable.`,
		Exposed: true,
		Name:    "restrictedPermissions",
		SubType: "string",
		Type:    "list",
	},
	"SourceName": {
		AllowedChoices: []string{},
		ConvertedName:  "SourceName",
		Description:    `The name of the source to use.`,
		Exposed:        true,
		Name:           "sourceName",
		Type:           "string",
	},
	"SourceNamespace": {
		AllowedChoices: []string{},
		ConvertedName:  "SourceNamespace",
		Description:    `The namespace of the source to use.`,
		Exposed:        true,
		Name:           "sourceNamespace",
		Type:           "string",
	},
	"SourceType": {
		AllowedChoices: []string{"A3S", "AWS", "Azure", "GCP", "HTTP", "LDAP", "MTLS", "OIDC", "RemoteA3S", "SAML"},
		ConvertedName:  "SourceType",
		Description: `The authentication source. This will define how to verify
credentials from internal or external source of authentication.`,
		Exposed:  true,
		Name:     "sourceType",
		Required: true,
		Type:     "enum",
	},
	"Token": {
		AllowedChoices: []string{},
		Autogenerated:  true,
		ConvertedName:  "Token",
		Description:    `Issued token.`,
		Exposed:        true,
		Name:           "token",
		ReadOnly:       true,
		Type:           "string",
	},
	"TokenType": {
		AllowedChoices: []string{"Identity", "Refresh"},
		ConvertedName:  "TokenType",
		DefaultValue:   IssueTokenTypeIdentity,
		Description:    `The type of token to issue.`,
		Exposed:        true,
		Name:           "tokenType",
		Type:           "enum",
	},
	"Validity": {
		AllowedChoices: []string{},
		ConvertedName:  "Validity",
		DefaultValue:   "24h",
		Description: `Configures the maximum length of validity for a token, using
[Golang duration syntax](https://golang.org/pkg/time/#example_Duration). If it
is bigger than the configured max validity, it will be capped. Default: ` + "`" + `24h` + "`" + `.`,
		Exposed: true,
		Name:    "validity",
		Type:    "string",
	},
}

// IssueLowerCaseAttributesMap represents the map of attribute for Issue.
var IssueLowerCaseAttributesMap = map[string]elemental.AttributeSpecification{
	"audience": {
		AllowedChoices: []string{},
		ConvertedName:  "Audience",
		Description:    `Requested audience for the delivered token.`,
		Exposed:        true,
		Name:           "audience",
		SubType:        "string",
		Type:           "list",
	},
	"cloak": {
		AllowedChoices: []string{},
		ConvertedName:  "Cloak",
		Description: `Sets a list of identity claim prefix to allow in the final token. This can be
used to hide some information when asking for a token as not all systems need to
know all of the claims.`,
		Exposed: true,
		Name:    "cloak",
		SubType: "string",
		Type:    "list",
	},
	"cookie": {
		AllowedChoices: []string{},
		ConvertedName:  "Cookie",
		Description:    `If set, return the token as a secure cookie.`,
		Exposed:        true,
		Name:           "cookie",
		Type:           "boolean",
	},
	"cookiedomain": {
		AllowedChoices: []string{},
		ConvertedName:  "CookieDomain",
		Description:    `If set, use the provided domain for the delivered cookie.`,
		Exposed:        true,
		Name:           "cookieDomain",
		Type:           "string",
	},
	"inputa3s": {
		AllowedChoices: []string{},
		ConvertedName:  "InputA3S",
		Description:    `Contains additional information for an A3S token source.`,
		Exposed:        true,
		Name:           "inputA3S",
		SubType:        "issuea3s",
		Type:           "ref",
	},
	"inputaws": {
		AllowedChoices: []string{},
		ConvertedName:  "InputAWS",
		Description:    `Contains additional information for an AWS STS token source.`,
		Exposed:        true,
		Name:           "inputAWS",
		SubType:        "issueaws",
		Type:           "ref",
	},
	"inputazure": {
		AllowedChoices: []string{},
		ConvertedName:  "InputAzure",
		Description:    `Contains additional information for an Azure token source.`,
		Exposed:        true,
		Name:           "inputAzure",
		SubType:        "issueazure",
		Type:           "ref",
	},
	"inputgcp": {
		AllowedChoices: []string{},
		ConvertedName:  "InputGCP",
		Description:    `Contains additional information for an GCP token source.`,
		Exposed:        true,
		Name:           "inputGCP",
		SubType:        "issuegcp",
		Type:           "ref",
	},
	"inputhttp": {
		AllowedChoices: []string{},
		ConvertedName:  "InputHTTP",
		Description:    `Contains additional information for an HTTP source.`,
		Exposed:        true,
		Name:           "inputHTTP",
		SubType:        "issuehttp",
		Type:           "ref",
	},
	"inputldap": {
		AllowedChoices: []string{},
		ConvertedName:  "InputLDAP",
		Description:    `Contains additional information for an LDAP source.`,
		Exposed:        true,
		Name:           "inputLDAP",
		SubType:        "issueldap",
		Type:           "ref",
	},
	"inputoidc": {
		AllowedChoices: []string{},
		ConvertedName:  "InputOIDC",
		Description:    `Contains additional information for an OIDC source.`,
		Exposed:        true,
		Name:           "inputOIDC",
		SubType:        "issueoidc",
		Type:           "ref",
	},
	"inputremotea3s": {
		AllowedChoices: []string{},
		ConvertedName:  "InputRemoteA3S",
		Description:    `Contains additional information for a remote A3S token source.`,
		Exposed:        true,
		Name:           "inputRemoteA3S",
		SubType:        "issueremotea3s",
		Type:           "ref",
	},
	"opaque": {
		AllowedChoices: []string{},
		ConvertedName:  "Opaque",
		Description:    `Opaque data that will be included in the issued token.`,
		Exposed:        true,
		Name:           "opaque",
		SubType:        "map[string]string",
		Type:           "external",
	},
	"restrictednamespace": {
		AllowedChoices: []string{},
		ConvertedName:  "RestrictedNamespace",
		Description: `Restricts the namespace where the token can be used.

For instance, if you have have access to ` + "`" + `/namespace` + "`" + ` and below, you can
tell the policy engine that it should restrict further more to
` + "`" + `/namespace/child` + "`" + `.

Restricting to a namespace you don't have initially access according to the
policy engine has no effect and may end up making the token unusable.`,
		Exposed: true,
		Name:    "restrictedNamespace",
		Type:    "string",
	},
	"restrictednetworks": {
		AllowedChoices: []string{},
		ConvertedName:  "RestrictedNetworks",
		Description: `Restricts the networks from where the token can be used. This will reduce the
existing set of authorized networks that normally apply to the token according
to the policy engine.

For instance, If you have authorized access from ` + "`" + `0.0.0.0/0` + "`" + ` (by default) or
from
` + "`" + `10.0.0.0/8` + "`" + `, you can ask for a token that will only be valid if used from
` + "`" + `10.1.0.0/16` + "`" + `.

Restricting to a network that is not initially authorized by the policy
engine has no effect and may end up making the token unusable.`,
		Exposed: true,
		Name:    "restrictedNetworks",
		SubType: "string",
		Type:    "list",
	},
	"restrictedpermissions": {
		AllowedChoices: []string{},
		ConvertedName:  "RestrictedPermissions",
		Description: `Restricts the permissions of token. This will reduce the existing permissions
that normally apply to the token according to the policy engine.

For instance, if you have administrative role, you can ask for a token that will
tell the policy engine to reduce the permission it would have granted to what is
given defined in the token.

Restricting to some permissions you don't initially have according to the policy
engine has no effect and may end up making the token unusable.`,
		Exposed: true,
		Name:    "restrictedPermissions",
		SubType: "string",
		Type:    "list",
	},
	"sourcename": {
		AllowedChoices: []string{},
		ConvertedName:  "SourceName",
		Description:    `The name of the source to use.`,
		Exposed:        true,
		Name:           "sourceName",
		Type:           "string",
	},
	"sourcenamespace": {
		AllowedChoices: []string{},
		ConvertedName:  "SourceNamespace",
		Description:    `The namespace of the source to use.`,
		Exposed:        true,
		Name:           "sourceNamespace",
		Type:           "string",
	},
	"sourcetype": {
		AllowedChoices: []string{"A3S", "AWS", "Azure", "GCP", "HTTP", "LDAP", "MTLS", "OIDC", "RemoteA3S", "SAML"},
		ConvertedName:  "SourceType",
		Description: `The authentication source. This will define how to verify
credentials from internal or external source of authentication.`,
		Exposed:  true,
		Name:     "sourceType",
		Required: true,
		Type:     "enum",
	},
	"token": {
		AllowedChoices: []string{},
		Autogenerated:  true,
		ConvertedName:  "Token",
		Description:    `Issued token.`,
		Exposed:        true,
		Name:           "token",
		ReadOnly:       true,
		Type:           "string",
	},
	"tokentype": {
		AllowedChoices: []string{"Identity", "Refresh"},
		ConvertedName:  "TokenType",
		DefaultValue:   IssueTokenTypeIdentity,
		Description:    `The type of token to issue.`,
		Exposed:        true,
		Name:           "tokenType",
		Type:           "enum",
	},
	"validity": {
		AllowedChoices: []string{},
		ConvertedName:  "Validity",
		DefaultValue:   "24h",
		Description: `Configures the maximum length of validity for a token, using
[Golang duration syntax](https://golang.org/pkg/time/#example_Duration). If it
is bigger than the configured max validity, it will be capped. Default: ` + "`" + `24h` + "`" + `.`,
		Exposed: true,
		Name:    "validity",
		Type:    "string",
	},
}

// SparseIssuesList represents a list of SparseIssues
type SparseIssuesList []*SparseIssue

// Identity returns the identity of the objects in the list.
func (o SparseIssuesList) Identity() elemental.Identity {

	return IssueIdentity
}

// Copy returns a pointer to a copy the SparseIssuesList.
func (o SparseIssuesList) Copy() elemental.Identifiables {

	copy := append(SparseIssuesList{}, o...)
	return &copy
}

// Append appends the objects to the a new copy of the SparseIssuesList.
func (o SparseIssuesList) Append(objects ...elemental.Identifiable) elemental.Identifiables {

	out := append(SparseIssuesList{}, o...)
	for _, obj := range objects {
		out = append(out, obj.(*SparseIssue))
	}

	return out
}

// List converts the object to an elemental.IdentifiablesList.
func (o SparseIssuesList) List() elemental.IdentifiablesList {

	out := make(elemental.IdentifiablesList, len(o))
	for i := 0; i < len(o); i++ {
		out[i] = o[i]
	}

	return out
}

// DefaultOrder returns the default ordering fields of the content.
func (o SparseIssuesList) DefaultOrder() []string {

	return []string{}
}

// ToPlain returns the SparseIssuesList converted to IssuesList.
func (o SparseIssuesList) ToPlain() elemental.IdentifiablesList {

	out := make(elemental.IdentifiablesList, len(o))
	for i := 0; i < len(o); i++ {
		out[i] = o[i].ToPlain()
	}

	return out
}

// Version returns the version of the content.
func (o SparseIssuesList) Version() int {

	return 1
}

// SparseIssue represents the sparse version of a issue.
type SparseIssue struct {
	// Requested audience for the delivered token.
	Audience *[]string `json:"audience,omitempty" msgpack:"audience,omitempty" bson:"-" mapstructure:"audience,omitempty"`

	// Sets a list of identity claim prefix to allow in the final token. This can be
	// used to hide some information when asking for a token as not all systems need to
	// know all of the claims.
	Cloak *[]string `json:"cloak,omitempty" msgpack:"cloak,omitempty" bson:"-" mapstructure:"cloak,omitempty"`

	// If set, return the token as a secure cookie.
	Cookie *bool `json:"cookie,omitempty" msgpack:"cookie,omitempty" bson:"-" mapstructure:"cookie,omitempty"`

	// If set, use the provided domain for the delivered cookie.
	CookieDomain *string `json:"cookieDomain,omitempty" msgpack:"cookieDomain,omitempty" bson:"-" mapstructure:"cookieDomain,omitempty"`

	// Contains additional information for an A3S token source.
	InputA3S *IssueA3S `json:"inputA3S,omitempty" msgpack:"inputA3S,omitempty" bson:"-" mapstructure:"inputA3S,omitempty"`

	// Contains additional information for an AWS STS token source.
	InputAWS *IssueAWS `json:"inputAWS,omitempty" msgpack:"inputAWS,omitempty" bson:"-" mapstructure:"inputAWS,omitempty"`

	// Contains additional information for an Azure token source.
	InputAzure *IssueAzure `json:"inputAzure,omitempty" msgpack:"inputAzure,omitempty" bson:"-" mapstructure:"inputAzure,omitempty"`

	// Contains additional information for an GCP token source.
	InputGCP *IssueGCP `json:"inputGCP,omitempty" msgpack:"inputGCP,omitempty" bson:"-" mapstructure:"inputGCP,omitempty"`

	// Contains additional information for an HTTP source.
	InputHTTP *IssueHTTP `json:"inputHTTP,omitempty" msgpack:"inputHTTP,omitempty" bson:"-" mapstructure:"inputHTTP,omitempty"`

	// Contains additional information for an LDAP source.
	InputLDAP *IssueLDAP `json:"inputLDAP,omitempty" msgpack:"inputLDAP,omitempty" bson:"-" mapstructure:"inputLDAP,omitempty"`

	// Contains additional information for an OIDC source.
	InputOIDC *IssueOIDC `json:"inputOIDC,omitempty" msgpack:"inputOIDC,omitempty" bson:"-" mapstructure:"inputOIDC,omitempty"`

	// Contains additional information for a remote A3S token source.
	InputRemoteA3S *IssueRemoteA3S `json:"inputRemoteA3S,omitempty" msgpack:"inputRemoteA3S,omitempty" bson:"-" mapstructure:"inputRemoteA3S,omitempty"`

	// Opaque data that will be included in the issued token.
	Opaque *map[string]string `json:"opaque,omitempty" msgpack:"opaque,omitempty" bson:"-" mapstructure:"opaque,omitempty"`

	// Restricts the namespace where the token can be used.
	//
	// For instance, if you have have access to `/namespace` and below, you can
	// tell the policy engine that it should restrict further more to
	// `/namespace/child`.
	//
	// Restricting to a namespace you don't have initially access according to the
	// policy engine has no effect and may end up making the token unusable.
	RestrictedNamespace *string `json:"restrictedNamespace,omitempty" msgpack:"restrictedNamespace,omitempty" bson:"-" mapstructure:"restrictedNamespace,omitempty"`

	// Restricts the networks from where the token can be used. This will reduce the
	// existing set of authorized networks that normally apply to the token according
	// to the policy engine.
	//
	// For instance, If you have authorized access from `0.0.0.0/0` (by default) or
	// from
	// `10.0.0.0/8`, you can ask for a token that will only be valid if used from
	// `10.1.0.0/16`.
	//
	// Restricting to a network that is not initially authorized by the policy
	// engine has no effect and may end up making the token unusable.
	RestrictedNetworks *[]string `json:"restrictedNetworks,omitempty" msgpack:"restrictedNetworks,omitempty" bson:"-" mapstructure:"restrictedNetworks,omitempty"`

	// Restricts the permissions of token. This will reduce the existing permissions
	// that normally apply to the token according to the policy engine.
	//
	// For instance, if you have administrative role, you can ask for a token that will
	// tell the policy engine to reduce the permission it would have granted to what is
	// given defined in the token.
	//
	// Restricting to some permissions you don't initially have according to the policy
	// engine has no effect and may end up making the token unusable.
	RestrictedPermissions *[]string `json:"restrictedPermissions,omitempty" msgpack:"restrictedPermissions,omitempty" bson:"-" mapstructure:"restrictedPermissions,omitempty"`

	// The name of the source to use.
	SourceName *string `json:"sourceName,omitempty" msgpack:"sourceName,omitempty" bson:"-" mapstructure:"sourceName,omitempty"`

	// The namespace of the source to use.
	SourceNamespace *string `json:"sourceNamespace,omitempty" msgpack:"sourceNamespace,omitempty" bson:"-" mapstructure:"sourceNamespace,omitempty"`

	// The authentication source. This will define how to verify
	// credentials from internal or external source of authentication.
	SourceType *IssueSourceTypeValue `json:"sourceType,omitempty" msgpack:"sourceType,omitempty" bson:"-" mapstructure:"sourceType,omitempty"`

	// Issued token.
	Token *string `json:"token,omitempty" msgpack:"token,omitempty" bson:"-" mapstructure:"token,omitempty"`

	// The type of token to issue.
	TokenType *IssueTokenTypeValue `json:"tokenType,omitempty" msgpack:"tokenType,omitempty" bson:"-" mapstructure:"tokenType,omitempty"`

	// Configures the maximum length of validity for a token, using
	// [Golang duration syntax](https://golang.org/pkg/time/#example_Duration). If it
	// is bigger than the configured max validity, it will be capped. Default: `24h`.
	Validity *string `json:"validity,omitempty" msgpack:"validity,omitempty" bson:"-" mapstructure:"validity,omitempty"`

	ModelVersion int `json:"-" msgpack:"-" bson:"_modelversion"`
}

// NewSparseIssue returns a new  SparseIssue.
func NewSparseIssue() *SparseIssue {
	return &SparseIssue{}
}

// Identity returns the Identity of the sparse object.
func (o *SparseIssue) Identity() elemental.Identity {

	return IssueIdentity
}

// Identifier returns the value of the sparse object's unique identifier.
func (o *SparseIssue) Identifier() string {

	return ""
}

// SetIdentifier sets the value of the sparse object's unique identifier.
func (o *SparseIssue) SetIdentifier(id string) {

}

// GetBSON implements the bson marshaling interface.
// This is used to transparently convert ID to MongoDBID as ObectID.
func (o *SparseIssue) GetBSON() (interface{}, error) {

	if o == nil {
		return nil, nil
	}

	s := &mongoAttributesSparseIssue{}

	return s, nil
}

// SetBSON implements the bson marshaling interface.
// This is used to transparently convert ID to MongoDBID as ObectID.
func (o *SparseIssue) SetBSON(raw bson.Raw) error {

	if o == nil {
		return nil
	}

	s := &mongoAttributesSparseIssue{}
	if err := raw.Unmarshal(s); err != nil {
		return err
	}

	return nil
}

// Version returns the hardcoded version of the model.
func (o *SparseIssue) Version() int {

	return 1
}

// ToPlain returns the plain version of the sparse model.
func (o *SparseIssue) ToPlain() elemental.PlainIdentifiable {

	out := NewIssue()
	if o.Audience != nil {
		out.Audience = *o.Audience
	}
	if o.Cloak != nil {
		out.Cloak = *o.Cloak
	}
	if o.Cookie != nil {
		out.Cookie = *o.Cookie
	}
	if o.CookieDomain != nil {
		out.CookieDomain = *o.CookieDomain
	}
	if o.InputA3S != nil {
		out.InputA3S = o.InputA3S
	}
	if o.InputAWS != nil {
		out.InputAWS = o.InputAWS
	}
	if o.InputAzure != nil {
		out.InputAzure = o.InputAzure
	}
	if o.InputGCP != nil {
		out.InputGCP = o.InputGCP
	}
	if o.InputHTTP != nil {
		out.InputHTTP = o.InputHTTP
	}
	if o.InputLDAP != nil {
		out.InputLDAP = o.InputLDAP
	}
	if o.InputOIDC != nil {
		out.InputOIDC = o.InputOIDC
	}
	if o.InputRemoteA3S != nil {
		out.InputRemoteA3S = o.InputRemoteA3S
	}
	if o.Opaque != nil {
		out.Opaque = *o.Opaque
	}
	if o.RestrictedNamespace != nil {
		out.RestrictedNamespace = *o.RestrictedNamespace
	}
	if o.RestrictedNetworks != nil {
		out.RestrictedNetworks = *o.RestrictedNetworks
	}
	if o.RestrictedPermissions != nil {
		out.RestrictedPermissions = *o.RestrictedPermissions
	}
	if o.SourceName != nil {
		out.SourceName = *o.SourceName
	}
	if o.SourceNamespace != nil {
		out.SourceNamespace = *o.SourceNamespace
	}
	if o.SourceType != nil {
		out.SourceType = *o.SourceType
	}
	if o.Token != nil {
		out.Token = *o.Token
	}
	if o.TokenType != nil {
		out.TokenType = *o.TokenType
	}
	if o.Validity != nil {
		out.Validity = *o.Validity
	}

	return out
}

// DeepCopy returns a deep copy if the SparseIssue.
func (o *SparseIssue) DeepCopy() *SparseIssue {

	if o == nil {
		return nil
	}

	out := &SparseIssue{}
	o.DeepCopyInto(out)

	return out
}

// DeepCopyInto copies the receiver into the given *SparseIssue.
func (o *SparseIssue) DeepCopyInto(out *SparseIssue) {

	target, err := copystructure.Copy(o)
	if err != nil {
		panic(fmt.Sprintf("Unable to deepcopy SparseIssue: %s", err))
	}

	*out = *target.(*SparseIssue)
}

type mongoAttributesIssue struct {
}
type mongoAttributesSparseIssue struct {
}
