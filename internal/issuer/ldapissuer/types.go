package ldapissuer

import "fmt"

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
