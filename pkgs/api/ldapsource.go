package api

import (
	"fmt"

	"github.com/globalsign/mgo/bson"
	"github.com/mitchellh/copystructure"
	"go.aporeto.io/elemental"
)

// LDAPSourceSecurityProtocolValue represents the possible values for attribute "securityProtocol".
type LDAPSourceSecurityProtocolValue string

const (
	// LDAPSourceSecurityProtocolInbandTLS represents the value InbandTLS.
	LDAPSourceSecurityProtocolInbandTLS LDAPSourceSecurityProtocolValue = "InbandTLS"

	// LDAPSourceSecurityProtocolTLS represents the value TLS.
	LDAPSourceSecurityProtocolTLS LDAPSourceSecurityProtocolValue = "TLS"
)

// LDAPSourceIdentity represents the Identity of the object.
var LDAPSourceIdentity = elemental.Identity{
	Name:     "ldapsource",
	Category: "ldapsources",
	Package:  "a3s",
	Private:  false,
}

// LDAPSourcesList represents a list of LDAPSources
type LDAPSourcesList []*LDAPSource

// Identity returns the identity of the objects in the list.
func (o LDAPSourcesList) Identity() elemental.Identity {

	return LDAPSourceIdentity
}

// Copy returns a pointer to a copy the LDAPSourcesList.
func (o LDAPSourcesList) Copy() elemental.Identifiables {

	copy := append(LDAPSourcesList{}, o...)
	return &copy
}

// Append appends the objects to the a new copy of the LDAPSourcesList.
func (o LDAPSourcesList) Append(objects ...elemental.Identifiable) elemental.Identifiables {

	out := append(LDAPSourcesList{}, o...)
	for _, obj := range objects {
		out = append(out, obj.(*LDAPSource))
	}

	return out
}

// List converts the object to an elemental.IdentifiablesList.
func (o LDAPSourcesList) List() elemental.IdentifiablesList {

	out := make(elemental.IdentifiablesList, len(o))
	for i := 0; i < len(o); i++ {
		out[i] = o[i]
	}

	return out
}

// DefaultOrder returns the default ordering fields of the content.
func (o LDAPSourcesList) DefaultOrder() []string {

	return []string{}
}

// ToSparse returns the LDAPSourcesList converted to SparseLDAPSourcesList.
// Objects in the list will only contain the given fields. No field means entire field set.
func (o LDAPSourcesList) ToSparse(fields ...string) elemental.Identifiables {

	out := make(SparseLDAPSourcesList, len(o))
	for i := 0; i < len(o); i++ {
		out[i] = o[i].ToSparse(fields...).(*SparseLDAPSource)
	}

	return out
}

// Version returns the version of the content.
func (o LDAPSourcesList) Version() int {

	return 1
}

// LDAPSource represents the model of a ldapsource
type LDAPSource struct {
	// ID is the identifier of the object.
	ID string `json:"ID" msgpack:"ID" bson:"-" mapstructure:"ID,omitempty"`

	// IP address or FQDN of the LDAP server.
	Address string `json:"address" msgpack:"address" bson:"address" mapstructure:"address,omitempty"`

	// The base distinguished name (DN) to use for LDAP queries.
	BaseDN string `json:"baseDN" msgpack:"baseDN" bson:"basedn" mapstructure:"baseDN,omitempty"`

	// The DN to use to bind to the LDAP server.
	BindDN string `json:"bindDN" msgpack:"bindDN" bson:"binddn" mapstructure:"bindDN,omitempty"`

	// Password to be used with the `bindDN` to authenticate to the LDAP server.
	BindPassword string `json:"bindPassword" msgpack:"bindPassword" bson:"bindpassword" mapstructure:"bindPassword,omitempty"`

	// The filter to use to locate the relevant user accounts. For Windows-based
	// systems, the value may be `sAMAccountName={USERNAME}`. For Linux and other
	// systems, the value may be `uid={USERNAME}`.
	BindSearchFilter string `json:"bindSearchFilter" msgpack:"bindSearchFilter" bson:"bindsearchfilter" mapstructure:"bindSearchFilter,omitempty"`

	// Can be left empty if the LDAP server's certificate is signed by a public,
	// trusted certificate authority. Otherwise, include the public key of the
	// certificate authority that signed the LDAP server's certificate.
	CertificateAuthority string `json:"certificateAuthority,omitempty" msgpack:"certificateAuthority,omitempty" bson:"certificateauthority,omitempty" mapstructure:"certificateAuthority,omitempty"`

	// The description of the object.
	Description string `json:"description" msgpack:"description" bson:"description" mapstructure:"description,omitempty"`

	// Contains optional information about a remote service that can be used to modify
	// the claims that are about to be delivered using this authentication source.
	IdentityModifier *IdentityModifier `json:"identityModifier,omitempty" msgpack:"identityModifier,omitempty" bson:"-" mapstructure:"identityModifier,omitempty"`

	// A list of keys that must not be imported into the identity token. If
	// `includedKeys` is also set, and a key is in both lists, the key will be ignored.
	IgnoredKeys []string `json:"ignoredKeys,omitempty" msgpack:"ignoredKeys,omitempty" bson:"ignoredkeys,omitempty" mapstructure:"ignoredKeys,omitempty"`

	// A list of keys that must be imported into the identity token. If `ignoredKeys`
	// is also set, and a key is in both lists, the key will be ignored.
	IncludedKeys []string `json:"includedKeys,omitempty" msgpack:"includedKeys,omitempty" bson:"includedkeys,omitempty" mapstructure:"includedKeys,omitempty"`

	// The name of the source.
	Name string `json:"name" msgpack:"name" bson:"name" mapstructure:"name,omitempty"`

	// The namespace of the object.
	Namespace string `json:"namespace" msgpack:"namespace" bson:"namespace" mapstructure:"namespace,omitempty"`

	// Specifies the connection type for the LDAP provider.
	SecurityProtocol LDAPSourceSecurityProtocolValue `json:"securityProtocol" msgpack:"securityProtocol" bson:"securityprotocol" mapstructure:"securityProtocol,omitempty"`

	// Hash of the object used to shard the data.
	ZHash int `json:"-" msgpack:"-" bson:"zhash" mapstructure:"-,omitempty"`

	// Sharding zone.
	Zone int `json:"-" msgpack:"-" bson:"zone" mapstructure:"-,omitempty"`

	ModelVersion int `json:"-" msgpack:"-" bson:"_modelversion"`
}

// NewLDAPSource returns a new *LDAPSource
func NewLDAPSource() *LDAPSource {

	return &LDAPSource{
		ModelVersion:     1,
		BindSearchFilter: "uid={USERNAME}",
		IgnoredKeys:      []string{},
		IncludedKeys:     []string{},
		SecurityProtocol: LDAPSourceSecurityProtocolInbandTLS,
	}
}

// Identity returns the Identity of the object.
func (o *LDAPSource) Identity() elemental.Identity {

	return LDAPSourceIdentity
}

// Identifier returns the value of the object's unique identifier.
func (o *LDAPSource) Identifier() string {

	return o.ID
}

// SetIdentifier sets the value of the object's unique identifier.
func (o *LDAPSource) SetIdentifier(id string) {

	o.ID = id
}

// GetBSON implements the bson marshaling interface.
// This is used to transparently convert ID to MongoDBID as ObectID.
func (o *LDAPSource) GetBSON() (interface{}, error) {

	if o == nil {
		return nil, nil
	}

	s := &mongoAttributesLDAPSource{}

	if o.ID != "" {
		s.ID = bson.ObjectIdHex(o.ID)
	}
	s.Address = o.Address
	s.BaseDN = o.BaseDN
	s.BindDN = o.BindDN
	s.BindPassword = o.BindPassword
	s.BindSearchFilter = o.BindSearchFilter
	s.CertificateAuthority = o.CertificateAuthority
	s.Description = o.Description
	s.IgnoredKeys = o.IgnoredKeys
	s.IncludedKeys = o.IncludedKeys
	s.Name = o.Name
	s.Namespace = o.Namespace
	s.SecurityProtocol = o.SecurityProtocol
	s.ZHash = o.ZHash
	s.Zone = o.Zone

	return s, nil
}

// SetBSON implements the bson marshaling interface.
// This is used to transparently convert ID to MongoDBID as ObectID.
func (o *LDAPSource) SetBSON(raw bson.Raw) error {

	if o == nil {
		return nil
	}

	s := &mongoAttributesLDAPSource{}
	if err := raw.Unmarshal(s); err != nil {
		return err
	}

	o.ID = s.ID.Hex()
	o.Address = s.Address
	o.BaseDN = s.BaseDN
	o.BindDN = s.BindDN
	o.BindPassword = s.BindPassword
	o.BindSearchFilter = s.BindSearchFilter
	o.CertificateAuthority = s.CertificateAuthority
	o.Description = s.Description
	o.IgnoredKeys = s.IgnoredKeys
	o.IncludedKeys = s.IncludedKeys
	o.Name = s.Name
	o.Namespace = s.Namespace
	o.SecurityProtocol = s.SecurityProtocol
	o.ZHash = s.ZHash
	o.Zone = s.Zone

	return nil
}

// Version returns the hardcoded version of the model.
func (o *LDAPSource) Version() int {

	return 1
}

// BleveType implements the bleve.Classifier Interface.
func (o *LDAPSource) BleveType() string {

	return "ldapsource"
}

// DefaultOrder returns the list of default ordering fields.
func (o *LDAPSource) DefaultOrder() []string {

	return []string{}
}

// Doc returns the documentation for the object
func (o *LDAPSource) Doc() string {

	return `Defines a remote LDAP to use as an authentication source.`
}

func (o *LDAPSource) String() string {

	return fmt.Sprintf("<%s:%s>", o.Identity().Name, o.Identifier())
}

// GetID returns the ID of the receiver.
func (o *LDAPSource) GetID() string {

	return o.ID
}

// SetID sets the property ID of the receiver using the given value.
func (o *LDAPSource) SetID(ID string) {

	o.ID = ID
}

// GetNamespace returns the Namespace of the receiver.
func (o *LDAPSource) GetNamespace() string {

	return o.Namespace
}

// SetNamespace sets the property Namespace of the receiver using the given value.
func (o *LDAPSource) SetNamespace(namespace string) {

	o.Namespace = namespace
}

// GetZHash returns the ZHash of the receiver.
func (o *LDAPSource) GetZHash() int {

	return o.ZHash
}

// SetZHash sets the property ZHash of the receiver using the given value.
func (o *LDAPSource) SetZHash(zHash int) {

	o.ZHash = zHash
}

// GetZone returns the Zone of the receiver.
func (o *LDAPSource) GetZone() int {

	return o.Zone
}

// SetZone sets the property Zone of the receiver using the given value.
func (o *LDAPSource) SetZone(zone int) {

	o.Zone = zone
}

// ToSparse returns the sparse version of the model.
// The returned object will only contain the given fields. No field means entire field set.
func (o *LDAPSource) ToSparse(fields ...string) elemental.SparseIdentifiable {

	if len(fields) == 0 {
		// nolint: goimports
		return &SparseLDAPSource{
			ID:                   &o.ID,
			Address:              &o.Address,
			BaseDN:               &o.BaseDN,
			BindDN:               &o.BindDN,
			BindPassword:         &o.BindPassword,
			BindSearchFilter:     &o.BindSearchFilter,
			CertificateAuthority: &o.CertificateAuthority,
			Description:          &o.Description,
			IdentityModifier:     o.IdentityModifier,
			IgnoredKeys:          &o.IgnoredKeys,
			IncludedKeys:         &o.IncludedKeys,
			Name:                 &o.Name,
			Namespace:            &o.Namespace,
			SecurityProtocol:     &o.SecurityProtocol,
			ZHash:                &o.ZHash,
			Zone:                 &o.Zone,
		}
	}

	sp := &SparseLDAPSource{}
	for _, f := range fields {
		switch f {
		case "ID":
			sp.ID = &(o.ID)
		case "address":
			sp.Address = &(o.Address)
		case "baseDN":
			sp.BaseDN = &(o.BaseDN)
		case "bindDN":
			sp.BindDN = &(o.BindDN)
		case "bindPassword":
			sp.BindPassword = &(o.BindPassword)
		case "bindSearchFilter":
			sp.BindSearchFilter = &(o.BindSearchFilter)
		case "certificateAuthority":
			sp.CertificateAuthority = &(o.CertificateAuthority)
		case "description":
			sp.Description = &(o.Description)
		case "identityModifier":
			sp.IdentityModifier = o.IdentityModifier
		case "ignoredKeys":
			sp.IgnoredKeys = &(o.IgnoredKeys)
		case "includedKeys":
			sp.IncludedKeys = &(o.IncludedKeys)
		case "name":
			sp.Name = &(o.Name)
		case "namespace":
			sp.Namespace = &(o.Namespace)
		case "securityProtocol":
			sp.SecurityProtocol = &(o.SecurityProtocol)
		case "zHash":
			sp.ZHash = &(o.ZHash)
		case "zone":
			sp.Zone = &(o.Zone)
		}
	}

	return sp
}

// EncryptAttributes encrypts the attributes marked as `encrypted` using the given encrypter.
func (o *LDAPSource) EncryptAttributes(encrypter elemental.AttributeEncrypter) (err error) {

	if o.BindPassword, err = encrypter.EncryptString(o.BindPassword); err != nil {
		return fmt.Errorf("unable to encrypt attribute 'BindPassword' for 'LDAPSource' (%s): %s", o.Identifier(), err)
	}

	return nil
}

// DecryptAttributes decrypts the attributes marked as `encrypted` using the given decrypter.
func (o *LDAPSource) DecryptAttributes(encrypter elemental.AttributeEncrypter) (err error) {

	if o.BindPassword, err = encrypter.DecryptString(o.BindPassword); err != nil {
		return fmt.Errorf("unable to decrypt attribute 'BindPassword' for 'LDAPSource' (%s): %s", o.Identifier(), err)
	}

	return nil
}

// Patch apply the non nil value of a *SparseLDAPSource to the object.
func (o *LDAPSource) Patch(sparse elemental.SparseIdentifiable) {
	if !sparse.Identity().IsEqual(o.Identity()) {
		panic("cannot patch from a parse with different identity")
	}

	so := sparse.(*SparseLDAPSource)
	if so.ID != nil {
		o.ID = *so.ID
	}
	if so.Address != nil {
		o.Address = *so.Address
	}
	if so.BaseDN != nil {
		o.BaseDN = *so.BaseDN
	}
	if so.BindDN != nil {
		o.BindDN = *so.BindDN
	}
	if so.BindPassword != nil {
		o.BindPassword = *so.BindPassword
	}
	if so.BindSearchFilter != nil {
		o.BindSearchFilter = *so.BindSearchFilter
	}
	if so.CertificateAuthority != nil {
		o.CertificateAuthority = *so.CertificateAuthority
	}
	if so.Description != nil {
		o.Description = *so.Description
	}
	if so.IdentityModifier != nil {
		o.IdentityModifier = so.IdentityModifier
	}
	if so.IgnoredKeys != nil {
		o.IgnoredKeys = *so.IgnoredKeys
	}
	if so.IncludedKeys != nil {
		o.IncludedKeys = *so.IncludedKeys
	}
	if so.Name != nil {
		o.Name = *so.Name
	}
	if so.Namespace != nil {
		o.Namespace = *so.Namespace
	}
	if so.SecurityProtocol != nil {
		o.SecurityProtocol = *so.SecurityProtocol
	}
	if so.ZHash != nil {
		o.ZHash = *so.ZHash
	}
	if so.Zone != nil {
		o.Zone = *so.Zone
	}
}

// DeepCopy returns a deep copy if the LDAPSource.
func (o *LDAPSource) DeepCopy() *LDAPSource {

	if o == nil {
		return nil
	}

	out := &LDAPSource{}
	o.DeepCopyInto(out)

	return out
}

// DeepCopyInto copies the receiver into the given *LDAPSource.
func (o *LDAPSource) DeepCopyInto(out *LDAPSource) {

	target, err := copystructure.Copy(o)
	if err != nil {
		panic(fmt.Sprintf("Unable to deepcopy LDAPSource: %s", err))
	}

	*out = *target.(*LDAPSource)
}

// Validate valides the current information stored into the structure.
func (o *LDAPSource) Validate() error {

	errors := elemental.Errors{}
	requiredErrors := elemental.Errors{}

	if err := elemental.ValidateRequiredString("address", o.Address); err != nil {
		requiredErrors = requiredErrors.Append(err)
	}

	if err := elemental.ValidateRequiredString("baseDN", o.BaseDN); err != nil {
		requiredErrors = requiredErrors.Append(err)
	}

	if err := elemental.ValidateRequiredString("bindDN", o.BindDN); err != nil {
		requiredErrors = requiredErrors.Append(err)
	}

	if o.IdentityModifier != nil {
		elemental.ResetDefaultForZeroValues(o.IdentityModifier)
		if err := o.IdentityModifier.Validate(); err != nil {
			errors = errors.Append(err)
		}
	}

	if err := elemental.ValidateRequiredString("name", o.Name); err != nil {
		requiredErrors = requiredErrors.Append(err)
	}

	if err := elemental.ValidateStringInList("securityProtocol", string(o.SecurityProtocol), []string{"TLS", "InbandTLS"}, false); err != nil {
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
func (*LDAPSource) SpecificationForAttribute(name string) elemental.AttributeSpecification {

	if v, ok := LDAPSourceAttributesMap[name]; ok {
		return v
	}

	// We could not find it, so let's check on the lower case indexed spec map
	return LDAPSourceLowerCaseAttributesMap[name]
}

// AttributeSpecifications returns the full attribute specifications map.
func (*LDAPSource) AttributeSpecifications() map[string]elemental.AttributeSpecification {

	return LDAPSourceAttributesMap
}

// ValueForAttribute returns the value for the given attribute.
// This is a very advanced function that you should not need but in some
// very specific use cases.
func (o *LDAPSource) ValueForAttribute(name string) interface{} {

	switch name {
	case "ID":
		return o.ID
	case "address":
		return o.Address
	case "baseDN":
		return o.BaseDN
	case "bindDN":
		return o.BindDN
	case "bindPassword":
		return o.BindPassword
	case "bindSearchFilter":
		return o.BindSearchFilter
	case "certificateAuthority":
		return o.CertificateAuthority
	case "description":
		return o.Description
	case "identityModifier":
		return o.IdentityModifier
	case "ignoredKeys":
		return o.IgnoredKeys
	case "includedKeys":
		return o.IncludedKeys
	case "name":
		return o.Name
	case "namespace":
		return o.Namespace
	case "securityProtocol":
		return o.SecurityProtocol
	case "zHash":
		return o.ZHash
	case "zone":
		return o.Zone
	}

	return nil
}

// LDAPSourceAttributesMap represents the map of attribute for LDAPSource.
var LDAPSourceAttributesMap = map[string]elemental.AttributeSpecification{
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
	"Address": {
		AllowedChoices: []string{},
		BSONFieldName:  "address",
		ConvertedName:  "Address",
		Description:    `IP address or FQDN of the LDAP server.`,
		Exposed:        true,
		Name:           "address",
		Required:       true,
		Stored:         true,
		Type:           "string",
	},
	"BaseDN": {
		AllowedChoices: []string{},
		BSONFieldName:  "basedn",
		ConvertedName:  "BaseDN",
		Description:    `The base distinguished name (DN) to use for LDAP queries.`,
		Exposed:        true,
		Name:           "baseDN",
		Required:       true,
		Stored:         true,
		Type:           "string",
	},
	"BindDN": {
		AllowedChoices: []string{},
		BSONFieldName:  "binddn",
		ConvertedName:  "BindDN",
		Description:    `The DN to use to bind to the LDAP server.`,
		Exposed:        true,
		Name:           "bindDN",
		Required:       true,
		Stored:         true,
		Type:           "string",
	},
	"BindPassword": {
		AllowedChoices: []string{},
		BSONFieldName:  "bindpassword",
		ConvertedName:  "BindPassword",
		Description:    `Password to be used with the ` + "`" + `bindDN` + "`" + ` to authenticate to the LDAP server.`,
		Encrypted:      true,
		Exposed:        true,
		Name:           "bindPassword",
		Required:       true,
		Secret:         true,
		Stored:         true,
		Transient:      true,
		Type:           "string",
	},
	"BindSearchFilter": {
		AllowedChoices: []string{},
		BSONFieldName:  "bindsearchfilter",
		ConvertedName:  "BindSearchFilter",
		DefaultValue:   "uid={USERNAME}",
		Description: `The filter to use to locate the relevant user accounts. For Windows-based
systems, the value may be ` + "`" + `sAMAccountName={USERNAME}` + "`" + `. For Linux and other
systems, the value may be ` + "`" + `uid={USERNAME}` + "`" + `.`,
		Exposed:   true,
		Name:      "bindSearchFilter",
		Orderable: true,
		Stored:    true,
		Type:      "string",
	},
	"CertificateAuthority": {
		AllowedChoices: []string{},
		BSONFieldName:  "certificateauthority",
		ConvertedName:  "CertificateAuthority",
		Description: `Can be left empty if the LDAP server's certificate is signed by a public,
trusted certificate authority. Otherwise, include the public key of the
certificate authority that signed the LDAP server's certificate.`,
		Exposed: true,
		Name:    "certificateAuthority",
		Stored:  true,
		Type:    "string",
	},
	"Description": {
		AllowedChoices: []string{},
		BSONFieldName:  "description",
		ConvertedName:  "Description",
		Description:    `The description of the object.`,
		Exposed:        true,
		Name:           "description",
		Stored:         true,
		Type:           "string",
	},
	"IdentityModifier": {
		AllowedChoices: []string{},
		ConvertedName:  "IdentityModifier",
		Description: `Contains optional information about a remote service that can be used to modify
the claims that are about to be delivered using this authentication source.`,
		Exposed: true,
		Name:    "identityModifier",
		SubType: "identitymodifier",
		Type:    "ref",
	},
	"IgnoredKeys": {
		AllowedChoices: []string{},
		BSONFieldName:  "ignoredkeys",
		ConvertedName:  "IgnoredKeys",
		Description: `A list of keys that must not be imported into the identity token. If
` + "`" + `includedKeys` + "`" + ` is also set, and a key is in both lists, the key will be ignored.`,
		Exposed: true,
		Name:    "ignoredKeys",
		Stored:  true,
		SubType: "string",
		Type:    "list",
	},
	"IncludedKeys": {
		AllowedChoices: []string{},
		BSONFieldName:  "includedkeys",
		ConvertedName:  "IncludedKeys",
		Description: `A list of keys that must be imported into the identity token. If ` + "`" + `ignoredKeys` + "`" + `
is also set, and a key is in both lists, the key will be ignored.`,
		Exposed: true,
		Name:    "includedKeys",
		Stored:  true,
		SubType: "string",
		Type:    "list",
	},
	"Name": {
		AllowedChoices: []string{},
		BSONFieldName:  "name",
		ConvertedName:  "Name",
		Description:    `The name of the source.`,
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
	"SecurityProtocol": {
		AllowedChoices: []string{"TLS", "InbandTLS"},
		BSONFieldName:  "securityprotocol",
		ConvertedName:  "SecurityProtocol",
		DefaultValue:   LDAPSourceSecurityProtocolInbandTLS,
		Description:    `Specifies the connection type for the LDAP provider.`,
		Exposed:        true,
		Name:           "securityProtocol",
		Stored:         true,
		Type:           "enum",
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

// LDAPSourceLowerCaseAttributesMap represents the map of attribute for LDAPSource.
var LDAPSourceLowerCaseAttributesMap = map[string]elemental.AttributeSpecification{
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
	"address": {
		AllowedChoices: []string{},
		BSONFieldName:  "address",
		ConvertedName:  "Address",
		Description:    `IP address or FQDN of the LDAP server.`,
		Exposed:        true,
		Name:           "address",
		Required:       true,
		Stored:         true,
		Type:           "string",
	},
	"basedn": {
		AllowedChoices: []string{},
		BSONFieldName:  "basedn",
		ConvertedName:  "BaseDN",
		Description:    `The base distinguished name (DN) to use for LDAP queries.`,
		Exposed:        true,
		Name:           "baseDN",
		Required:       true,
		Stored:         true,
		Type:           "string",
	},
	"binddn": {
		AllowedChoices: []string{},
		BSONFieldName:  "binddn",
		ConvertedName:  "BindDN",
		Description:    `The DN to use to bind to the LDAP server.`,
		Exposed:        true,
		Name:           "bindDN",
		Required:       true,
		Stored:         true,
		Type:           "string",
	},
	"bindpassword": {
		AllowedChoices: []string{},
		BSONFieldName:  "bindpassword",
		ConvertedName:  "BindPassword",
		Description:    `Password to be used with the ` + "`" + `bindDN` + "`" + ` to authenticate to the LDAP server.`,
		Encrypted:      true,
		Exposed:        true,
		Name:           "bindPassword",
		Required:       true,
		Secret:         true,
		Stored:         true,
		Transient:      true,
		Type:           "string",
	},
	"bindsearchfilter": {
		AllowedChoices: []string{},
		BSONFieldName:  "bindsearchfilter",
		ConvertedName:  "BindSearchFilter",
		DefaultValue:   "uid={USERNAME}",
		Description: `The filter to use to locate the relevant user accounts. For Windows-based
systems, the value may be ` + "`" + `sAMAccountName={USERNAME}` + "`" + `. For Linux and other
systems, the value may be ` + "`" + `uid={USERNAME}` + "`" + `.`,
		Exposed:   true,
		Name:      "bindSearchFilter",
		Orderable: true,
		Stored:    true,
		Type:      "string",
	},
	"certificateauthority": {
		AllowedChoices: []string{},
		BSONFieldName:  "certificateauthority",
		ConvertedName:  "CertificateAuthority",
		Description: `Can be left empty if the LDAP server's certificate is signed by a public,
trusted certificate authority. Otherwise, include the public key of the
certificate authority that signed the LDAP server's certificate.`,
		Exposed: true,
		Name:    "certificateAuthority",
		Stored:  true,
		Type:    "string",
	},
	"description": {
		AllowedChoices: []string{},
		BSONFieldName:  "description",
		ConvertedName:  "Description",
		Description:    `The description of the object.`,
		Exposed:        true,
		Name:           "description",
		Stored:         true,
		Type:           "string",
	},
	"identitymodifier": {
		AllowedChoices: []string{},
		ConvertedName:  "IdentityModifier",
		Description: `Contains optional information about a remote service that can be used to modify
the claims that are about to be delivered using this authentication source.`,
		Exposed: true,
		Name:    "identityModifier",
		SubType: "identitymodifier",
		Type:    "ref",
	},
	"ignoredkeys": {
		AllowedChoices: []string{},
		BSONFieldName:  "ignoredkeys",
		ConvertedName:  "IgnoredKeys",
		Description: `A list of keys that must not be imported into the identity token. If
` + "`" + `includedKeys` + "`" + ` is also set, and a key is in both lists, the key will be ignored.`,
		Exposed: true,
		Name:    "ignoredKeys",
		Stored:  true,
		SubType: "string",
		Type:    "list",
	},
	"includedkeys": {
		AllowedChoices: []string{},
		BSONFieldName:  "includedkeys",
		ConvertedName:  "IncludedKeys",
		Description: `A list of keys that must be imported into the identity token. If ` + "`" + `ignoredKeys` + "`" + `
is also set, and a key is in both lists, the key will be ignored.`,
		Exposed: true,
		Name:    "includedKeys",
		Stored:  true,
		SubType: "string",
		Type:    "list",
	},
	"name": {
		AllowedChoices: []string{},
		BSONFieldName:  "name",
		ConvertedName:  "Name",
		Description:    `The name of the source.`,
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
	"securityprotocol": {
		AllowedChoices: []string{"TLS", "InbandTLS"},
		BSONFieldName:  "securityprotocol",
		ConvertedName:  "SecurityProtocol",
		DefaultValue:   LDAPSourceSecurityProtocolInbandTLS,
		Description:    `Specifies the connection type for the LDAP provider.`,
		Exposed:        true,
		Name:           "securityProtocol",
		Stored:         true,
		Type:           "enum",
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

// SparseLDAPSourcesList represents a list of SparseLDAPSources
type SparseLDAPSourcesList []*SparseLDAPSource

// Identity returns the identity of the objects in the list.
func (o SparseLDAPSourcesList) Identity() elemental.Identity {

	return LDAPSourceIdentity
}

// Copy returns a pointer to a copy the SparseLDAPSourcesList.
func (o SparseLDAPSourcesList) Copy() elemental.Identifiables {

	copy := append(SparseLDAPSourcesList{}, o...)
	return &copy
}

// Append appends the objects to the a new copy of the SparseLDAPSourcesList.
func (o SparseLDAPSourcesList) Append(objects ...elemental.Identifiable) elemental.Identifiables {

	out := append(SparseLDAPSourcesList{}, o...)
	for _, obj := range objects {
		out = append(out, obj.(*SparseLDAPSource))
	}

	return out
}

// List converts the object to an elemental.IdentifiablesList.
func (o SparseLDAPSourcesList) List() elemental.IdentifiablesList {

	out := make(elemental.IdentifiablesList, len(o))
	for i := 0; i < len(o); i++ {
		out[i] = o[i]
	}

	return out
}

// DefaultOrder returns the default ordering fields of the content.
func (o SparseLDAPSourcesList) DefaultOrder() []string {

	return []string{}
}

// ToPlain returns the SparseLDAPSourcesList converted to LDAPSourcesList.
func (o SparseLDAPSourcesList) ToPlain() elemental.IdentifiablesList {

	out := make(elemental.IdentifiablesList, len(o))
	for i := 0; i < len(o); i++ {
		out[i] = o[i].ToPlain()
	}

	return out
}

// Version returns the version of the content.
func (o SparseLDAPSourcesList) Version() int {

	return 1
}

// SparseLDAPSource represents the sparse version of a ldapsource.
type SparseLDAPSource struct {
	// ID is the identifier of the object.
	ID *string `json:"ID,omitempty" msgpack:"ID,omitempty" bson:"-" mapstructure:"ID,omitempty"`

	// IP address or FQDN of the LDAP server.
	Address *string `json:"address,omitempty" msgpack:"address,omitempty" bson:"address,omitempty" mapstructure:"address,omitempty"`

	// The base distinguished name (DN) to use for LDAP queries.
	BaseDN *string `json:"baseDN,omitempty" msgpack:"baseDN,omitempty" bson:"basedn,omitempty" mapstructure:"baseDN,omitempty"`

	// The DN to use to bind to the LDAP server.
	BindDN *string `json:"bindDN,omitempty" msgpack:"bindDN,omitempty" bson:"binddn,omitempty" mapstructure:"bindDN,omitempty"`

	// Password to be used with the `bindDN` to authenticate to the LDAP server.
	BindPassword *string `json:"bindPassword,omitempty" msgpack:"bindPassword,omitempty" bson:"bindpassword,omitempty" mapstructure:"bindPassword,omitempty"`

	// The filter to use to locate the relevant user accounts. For Windows-based
	// systems, the value may be `sAMAccountName={USERNAME}`. For Linux and other
	// systems, the value may be `uid={USERNAME}`.
	BindSearchFilter *string `json:"bindSearchFilter,omitempty" msgpack:"bindSearchFilter,omitempty" bson:"bindsearchfilter,omitempty" mapstructure:"bindSearchFilter,omitempty"`

	// Can be left empty if the LDAP server's certificate is signed by a public,
	// trusted certificate authority. Otherwise, include the public key of the
	// certificate authority that signed the LDAP server's certificate.
	CertificateAuthority *string `json:"certificateAuthority,omitempty" msgpack:"certificateAuthority,omitempty" bson:"certificateauthority,omitempty" mapstructure:"certificateAuthority,omitempty"`

	// The description of the object.
	Description *string `json:"description,omitempty" msgpack:"description,omitempty" bson:"description,omitempty" mapstructure:"description,omitempty"`

	// Contains optional information about a remote service that can be used to modify
	// the claims that are about to be delivered using this authentication source.
	IdentityModifier *IdentityModifier `json:"identityModifier,omitempty" msgpack:"identityModifier,omitempty" bson:"-" mapstructure:"identityModifier,omitempty"`

	// A list of keys that must not be imported into the identity token. If
	// `includedKeys` is also set, and a key is in both lists, the key will be ignored.
	IgnoredKeys *[]string `json:"ignoredKeys,omitempty" msgpack:"ignoredKeys,omitempty" bson:"ignoredkeys,omitempty" mapstructure:"ignoredKeys,omitempty"`

	// A list of keys that must be imported into the identity token. If `ignoredKeys`
	// is also set, and a key is in both lists, the key will be ignored.
	IncludedKeys *[]string `json:"includedKeys,omitempty" msgpack:"includedKeys,omitempty" bson:"includedkeys,omitempty" mapstructure:"includedKeys,omitempty"`

	// The name of the source.
	Name *string `json:"name,omitempty" msgpack:"name,omitempty" bson:"name,omitempty" mapstructure:"name,omitempty"`

	// The namespace of the object.
	Namespace *string `json:"namespace,omitempty" msgpack:"namespace,omitempty" bson:"namespace,omitempty" mapstructure:"namespace,omitempty"`

	// Specifies the connection type for the LDAP provider.
	SecurityProtocol *LDAPSourceSecurityProtocolValue `json:"securityProtocol,omitempty" msgpack:"securityProtocol,omitempty" bson:"securityprotocol,omitempty" mapstructure:"securityProtocol,omitempty"`

	// Hash of the object used to shard the data.
	ZHash *int `json:"-" msgpack:"-" bson:"zhash,omitempty" mapstructure:"-,omitempty"`

	// Sharding zone.
	Zone *int `json:"-" msgpack:"-" bson:"zone,omitempty" mapstructure:"-,omitempty"`

	ModelVersion int `json:"-" msgpack:"-" bson:"_modelversion"`
}

// NewSparseLDAPSource returns a new  SparseLDAPSource.
func NewSparseLDAPSource() *SparseLDAPSource {
	return &SparseLDAPSource{}
}

// Identity returns the Identity of the sparse object.
func (o *SparseLDAPSource) Identity() elemental.Identity {

	return LDAPSourceIdentity
}

// Identifier returns the value of the sparse object's unique identifier.
func (o *SparseLDAPSource) Identifier() string {

	if o.ID == nil {
		return ""
	}
	return *o.ID
}

// SetIdentifier sets the value of the sparse object's unique identifier.
func (o *SparseLDAPSource) SetIdentifier(id string) {

	if id != "" {
		o.ID = &id
	} else {
		o.ID = nil
	}
}

// GetBSON implements the bson marshaling interface.
// This is used to transparently convert ID to MongoDBID as ObectID.
func (o *SparseLDAPSource) GetBSON() (interface{}, error) {

	if o == nil {
		return nil, nil
	}

	s := &mongoAttributesSparseLDAPSource{}

	if o.ID != nil {
		s.ID = bson.ObjectIdHex(*o.ID)
	}
	if o.Address != nil {
		s.Address = o.Address
	}
	if o.BaseDN != nil {
		s.BaseDN = o.BaseDN
	}
	if o.BindDN != nil {
		s.BindDN = o.BindDN
	}
	if o.BindPassword != nil {
		s.BindPassword = o.BindPassword
	}
	if o.BindSearchFilter != nil {
		s.BindSearchFilter = o.BindSearchFilter
	}
	if o.CertificateAuthority != nil {
		s.CertificateAuthority = o.CertificateAuthority
	}
	if o.Description != nil {
		s.Description = o.Description
	}
	if o.IgnoredKeys != nil {
		s.IgnoredKeys = o.IgnoredKeys
	}
	if o.IncludedKeys != nil {
		s.IncludedKeys = o.IncludedKeys
	}
	if o.Name != nil {
		s.Name = o.Name
	}
	if o.Namespace != nil {
		s.Namespace = o.Namespace
	}
	if o.SecurityProtocol != nil {
		s.SecurityProtocol = o.SecurityProtocol
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
func (o *SparseLDAPSource) SetBSON(raw bson.Raw) error {

	if o == nil {
		return nil
	}

	s := &mongoAttributesSparseLDAPSource{}
	if err := raw.Unmarshal(s); err != nil {
		return err
	}

	id := s.ID.Hex()
	o.ID = &id
	if s.Address != nil {
		o.Address = s.Address
	}
	if s.BaseDN != nil {
		o.BaseDN = s.BaseDN
	}
	if s.BindDN != nil {
		o.BindDN = s.BindDN
	}
	if s.BindPassword != nil {
		o.BindPassword = s.BindPassword
	}
	if s.BindSearchFilter != nil {
		o.BindSearchFilter = s.BindSearchFilter
	}
	if s.CertificateAuthority != nil {
		o.CertificateAuthority = s.CertificateAuthority
	}
	if s.Description != nil {
		o.Description = s.Description
	}
	if s.IgnoredKeys != nil {
		o.IgnoredKeys = s.IgnoredKeys
	}
	if s.IncludedKeys != nil {
		o.IncludedKeys = s.IncludedKeys
	}
	if s.Name != nil {
		o.Name = s.Name
	}
	if s.Namespace != nil {
		o.Namespace = s.Namespace
	}
	if s.SecurityProtocol != nil {
		o.SecurityProtocol = s.SecurityProtocol
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
func (o *SparseLDAPSource) Version() int {

	return 1
}

// ToPlain returns the plain version of the sparse model.
func (o *SparseLDAPSource) ToPlain() elemental.PlainIdentifiable {

	out := NewLDAPSource()
	if o.ID != nil {
		out.ID = *o.ID
	}
	if o.Address != nil {
		out.Address = *o.Address
	}
	if o.BaseDN != nil {
		out.BaseDN = *o.BaseDN
	}
	if o.BindDN != nil {
		out.BindDN = *o.BindDN
	}
	if o.BindPassword != nil {
		out.BindPassword = *o.BindPassword
	}
	if o.BindSearchFilter != nil {
		out.BindSearchFilter = *o.BindSearchFilter
	}
	if o.CertificateAuthority != nil {
		out.CertificateAuthority = *o.CertificateAuthority
	}
	if o.Description != nil {
		out.Description = *o.Description
	}
	if o.IdentityModifier != nil {
		out.IdentityModifier = o.IdentityModifier
	}
	if o.IgnoredKeys != nil {
		out.IgnoredKeys = *o.IgnoredKeys
	}
	if o.IncludedKeys != nil {
		out.IncludedKeys = *o.IncludedKeys
	}
	if o.Name != nil {
		out.Name = *o.Name
	}
	if o.Namespace != nil {
		out.Namespace = *o.Namespace
	}
	if o.SecurityProtocol != nil {
		out.SecurityProtocol = *o.SecurityProtocol
	}
	if o.ZHash != nil {
		out.ZHash = *o.ZHash
	}
	if o.Zone != nil {
		out.Zone = *o.Zone
	}

	return out
}

// EncryptAttributes encrypts the attributes marked as `encrypted` using the given encrypter.
func (o *SparseLDAPSource) EncryptAttributes(encrypter elemental.AttributeEncrypter) (err error) {

	if *o.BindPassword, err = encrypter.EncryptString(*o.BindPassword); err != nil {
		return fmt.Errorf("unable to encrypt attribute 'BindPassword' for 'SparseLDAPSource' (%s): %s", o.Identifier(), err)
	}

	return nil
}

// DecryptAttributes decrypts the attributes marked as `encrypted` using the given decrypter.
func (o *SparseLDAPSource) DecryptAttributes(encrypter elemental.AttributeEncrypter) (err error) {

	if *o.BindPassword, err = encrypter.DecryptString(*o.BindPassword); err != nil {
		return fmt.Errorf("unable to decrypt attribute 'BindPassword' for 'SparseLDAPSource' (%s): %s", o.Identifier(), err)
	}

	return nil
}

// GetID returns the ID of the receiver.
func (o *SparseLDAPSource) GetID() (out string) {

	if o.ID == nil {
		return
	}

	return *o.ID
}

// SetID sets the property ID of the receiver using the address of the given value.
func (o *SparseLDAPSource) SetID(ID string) {

	o.ID = &ID
}

// GetNamespace returns the Namespace of the receiver.
func (o *SparseLDAPSource) GetNamespace() (out string) {

	if o.Namespace == nil {
		return
	}

	return *o.Namespace
}

// SetNamespace sets the property Namespace of the receiver using the address of the given value.
func (o *SparseLDAPSource) SetNamespace(namespace string) {

	o.Namespace = &namespace
}

// GetZHash returns the ZHash of the receiver.
func (o *SparseLDAPSource) GetZHash() (out int) {

	if o.ZHash == nil {
		return
	}

	return *o.ZHash
}

// SetZHash sets the property ZHash of the receiver using the address of the given value.
func (o *SparseLDAPSource) SetZHash(zHash int) {

	o.ZHash = &zHash
}

// GetZone returns the Zone of the receiver.
func (o *SparseLDAPSource) GetZone() (out int) {

	if o.Zone == nil {
		return
	}

	return *o.Zone
}

// SetZone sets the property Zone of the receiver using the address of the given value.
func (o *SparseLDAPSource) SetZone(zone int) {

	o.Zone = &zone
}

// DeepCopy returns a deep copy if the SparseLDAPSource.
func (o *SparseLDAPSource) DeepCopy() *SparseLDAPSource {

	if o == nil {
		return nil
	}

	out := &SparseLDAPSource{}
	o.DeepCopyInto(out)

	return out
}

// DeepCopyInto copies the receiver into the given *SparseLDAPSource.
func (o *SparseLDAPSource) DeepCopyInto(out *SparseLDAPSource) {

	target, err := copystructure.Copy(o)
	if err != nil {
		panic(fmt.Sprintf("Unable to deepcopy SparseLDAPSource: %s", err))
	}

	*out = *target.(*SparseLDAPSource)
}

type mongoAttributesLDAPSource struct {
	ID                   bson.ObjectId                   `bson:"_id,omitempty"`
	Address              string                          `bson:"address"`
	BaseDN               string                          `bson:"basedn"`
	BindDN               string                          `bson:"binddn"`
	BindPassword         string                          `bson:"bindpassword"`
	BindSearchFilter     string                          `bson:"bindsearchfilter"`
	CertificateAuthority string                          `bson:"certificateauthority,omitempty"`
	Description          string                          `bson:"description"`
	IgnoredKeys          []string                        `bson:"ignoredkeys,omitempty"`
	IncludedKeys         []string                        `bson:"includedkeys,omitempty"`
	Name                 string                          `bson:"name"`
	Namespace            string                          `bson:"namespace"`
	SecurityProtocol     LDAPSourceSecurityProtocolValue `bson:"securityprotocol"`
	ZHash                int                             `bson:"zhash"`
	Zone                 int                             `bson:"zone"`
}
type mongoAttributesSparseLDAPSource struct {
	ID                   bson.ObjectId                    `bson:"_id,omitempty"`
	Address              *string                          `bson:"address,omitempty"`
	BaseDN               *string                          `bson:"basedn,omitempty"`
	BindDN               *string                          `bson:"binddn,omitempty"`
	BindPassword         *string                          `bson:"bindpassword,omitempty"`
	BindSearchFilter     *string                          `bson:"bindsearchfilter,omitempty"`
	CertificateAuthority *string                          `bson:"certificateauthority,omitempty"`
	Description          *string                          `bson:"description,omitempty"`
	IgnoredKeys          *[]string                        `bson:"ignoredkeys,omitempty"`
	IncludedKeys         *[]string                        `bson:"includedkeys,omitempty"`
	Name                 *string                          `bson:"name,omitempty"`
	Namespace            *string                          `bson:"namespace,omitempty"`
	SecurityProtocol     *LDAPSourceSecurityProtocolValue `bson:"securityprotocol,omitempty"`
	ZHash                *int                             `bson:"zhash,omitempty"`
	Zone                 *int                             `bson:"zone,omitempty"`
}
