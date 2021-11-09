package issuer

import (
	"crypto/tls"
	"crypto/x509"
	"fmt"
	"strings"

	"github.com/go-ldap/ldap/v3"
	"go.aporeto.io/a3s/pkgs/api"
	"go.aporeto.io/a3s/pkgs/token"
)

// An ErrLDAP represents an error that can occur
// during interactions with an LDAP server.
type ErrLDAP struct {
	Err error
}

func (e ErrLDAP) Error() string {
	return fmt.Sprintf("ldap error: %s", e.Err)
}

func (e ErrLDAP) Unwrap() error {
	return e.Err
}

// An LDAPIssuer issues identity token from an
// LDAP authentication sources.
type LDAPIssuer struct {
	token  *token.IdentityToken
	source *api.LDAPSource
}

// NewLDAPIssuer returns a new LDAPIssuer.
func NewLDAPIssuer(source *api.LDAPSource) *LDAPIssuer {

	return &LDAPIssuer{
		source: source,
		token: token.NewIdentityToken(token.Source{
			Type:      "ldap",
			Namespace: source.Namespace,
			Name:      source.Name,
		}),
	}
}

// FromCredentials computes the claims information based on the provided username and password.
func (c *LDAPIssuer) FromCredentials(username string, password string) (err error) {

	entry, dn, err := c.retrieveEntry(username, password)
	if err != nil {
		return err
	}

	inc, exc := computeInclusion(c.source)

	c.token.Identity = computeClaims(entry, dn, inc, exc)

	return nil
}

// Issue issues the identity token.
func (c *LDAPIssuer) Issue() *token.IdentityToken {

	return c.token
}

func (c *LDAPIssuer) retrieveEntry(username string, password string) (*ldap.Entry, *ldap.DN, error) {

	var err error

	var caPool *x509.CertPool
	if ca := c.source.CertificateAuthority; ca != "" {
		caPool = x509.NewCertPool()
		caPool.AppendCertsFromPEM([]byte(ca))
	} else {
		caPool, err = x509.SystemCertPool()
		if err != nil {
			return nil, nil, err
		}
	}

	tlsConfig := &tls.Config{
		ServerName: strings.Split(c.source.Address, ":")[0],
		RootCAs:    caPool,
	}

	var conn *ldap.Conn
	if c.source.SecurityProtocol == api.LDAPSourceSecurityProtocolTLS {
		conn, err = ldap.DialTLS("tcp", c.source.Address, tlsConfig)
		if err != nil {
			return nil, nil, ErrLDAP{Err: fmt.Errorf("cannot dial tls: %w", err)}
		}
	} else {
		conn, err = ldap.Dial("tcp", c.source.Address)
		if err != nil {
			return nil, nil, ErrLDAP{Err: fmt.Errorf("cannot dial: %w", err)}
		}
		if err = conn.StartTLS(tlsConfig); err != nil {
			return nil, nil, ErrLDAP{Err: fmt.Errorf("cannot start tls: %w", err)}
		}
	}
	defer conn.Close()

	if err = conn.Bind(c.source.BindDN, c.source.BindPassword); err != nil {
		return nil, nil, ErrLDAP{Err: fmt.Errorf("unable to bind: %w", err)}
	}

	req := ldap.NewSearchRequest(
		c.source.BaseDN,
		ldap.ScopeWholeSubtree,
		ldap.NeverDerefAliases,
		0,
		0,
		false,
		fmt.Sprintf("(&(%s))", strings.Replace(c.source.BindSearchFilter, "{USERNAME}", username, -1)),
		nil,
		nil,
	)

	sr, err := conn.Search(req)
	if err != nil {
		return nil, nil, ErrLDAP{Err: fmt.Errorf("unable to search: %w", err)}
	}

	if len(sr.Entries) != 1 {
		return nil, nil, ErrLDAP{Err: fmt.Errorf("invalid credentials")}
	}

	entry := sr.Entries[0]

	if err = conn.Bind(entry.DN, password); err != nil {
		return nil, nil, ErrLDAP{Err: fmt.Errorf("invalid credentials")}
	}

	dn, err := ldap.ParseDN(entry.DN)
	if err != nil {
		return nil, nil, ErrLDAP{Err: fmt.Errorf("unable to parse entry DN: %w", err)}
	}

	return entry, dn, nil
}

func computeClaims(entry *ldap.Entry, dn *ldap.DN, inc map[string]struct{}, exc map[string]struct{}) (claims []string) {

	claims = append(claims, fmt.Sprintf("dn=%s", entry.DN))

	for _, rdn := range dn.RDNs {

		attr := rdn.Attributes[0]

		if strings.ToLower(attr.Type) == "ou" {
			claims = append(claims, fmt.Sprintf("ou=%s", attr.Value))
		}

		if strings.ToLower(attr.Type) == "dc" {
			claims = append(claims, fmt.Sprintf("dc=%s", attr.Value))
		}
	}

	for _, attr := range entry.Attributes {

		if attr.Name == "userPassword" || attr.Name == "objectClass" || attr.Name == "comment" {
			continue
		}

		if _, ok := exc[strings.ToLower(attr.Name)]; ok {
			continue
		}

		if len(inc) > 0 {
			if _, ok := inc[strings.ToLower(attr.Name)]; !ok {
				continue
			}
		}

		if len(attr.Values) == 0 {
			continue
		}

		for _, v := range attr.Values {
			if v != "" {
				claims = append(claims, fmt.Sprintf("%s=%s", attr.Name, v))
			}
		}
	}

	return claims
}

func computeInclusion(src *api.LDAPSource) (inc map[string]struct{}, exc map[string]struct{}) {

	inc = make(map[string]struct{}, len(src.IncludedKeys))
	for _, key := range src.IncludedKeys {
		inc[strings.ToLower(key)] = struct{}{}
	}

	exc = make(map[string]struct{}, len(src.IgnoredKeys))
	for _, key := range src.IgnoredKeys {
		exc[strings.ToLower(key)] = struct{}{}
	}

	return inc, exc
}
