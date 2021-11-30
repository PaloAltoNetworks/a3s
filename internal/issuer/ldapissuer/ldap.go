package ldapissuer

import (
	"context"
	"crypto/tls"
	"crypto/x509"
	"fmt"
	"strings"

	"github.com/go-ldap/ldap/v3"
	"go.aporeto.io/a3s/pkgs/api"
	"go.aporeto.io/a3s/pkgs/token"
)

// New returns a new LDAP issuer.
func New(ctx context.Context, source *api.LDAPSource, username string, password string) (token.Issuer, error) {

	c := newLDAPIssuer(source)
	if err := c.fromCredentials(ctx, username, password); err != nil {
		return nil, err
	}

	return c, nil
}

type ldapIssuer struct {
	token  *token.IdentityToken
	source *api.LDAPSource
}

func newLDAPIssuer(source *api.LDAPSource) *ldapIssuer {

	return &ldapIssuer{
		source: source,
		token: token.NewIdentityToken(token.Source{
			Type:      "ldap",
			Namespace: source.Namespace,
			Name:      source.Name,
		}),
	}
}

func (c *ldapIssuer) Issue() *token.IdentityToken {

	return c.token
}

func (c *ldapIssuer) fromCredentials(ctx context.Context, username string, password string) (err error) {

	entry, dn, err := c.retrieveEntry(username, password)
	if err != nil {
		return err
	}

	inc, exc := computeLDPInclusion(c.source)

	c.token.Identity = computeLDAPClaims(entry, dn, inc, exc)

	if srcmod := c.source.Modifier; srcmod != nil {

		m, err := token.NewHTTPIdentityModifier(
			srcmod.URL,
			string(srcmod.Method),
			[]byte(srcmod.CA),
			[]byte(srcmod.Certificate),
			[]byte(srcmod.Key),
		)
		if err != nil {
			return fmt.Errorf("unable to prepare source modifier: %w", err)
		}

		if c.token.Identity, err = m.Modify(ctx, c.token.Identity); err != nil {
			return fmt.Errorf("unable to call modifier: %w", err)
		}
	}

	return nil
}

func (c *ldapIssuer) retrieveEntry(username string, password string) (*ldap.Entry, *ldap.DN, error) {

	var err error

	var caPool *x509.CertPool
	if ca := c.source.CA; ca != "" {
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

func computeLDAPClaims(entry *ldap.Entry, dn *ldap.DN, inc map[string]struct{}, exc map[string]struct{}) (claims []string) {

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
				claims = append(claims, fmt.Sprintf("%s=%s", strings.TrimLeft(attr.Name, "@"), v))
			}
		}
	}

	return claims
}

func computeLDPInclusion(src *api.LDAPSource) (inc map[string]struct{}, exc map[string]struct{}) {

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
