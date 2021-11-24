package a3sissuer

import "fmt"

// An ErrComputeRestrictions represents a generic
// restrictions computation error.
type ErrComputeRestrictions struct {
	Err error
}

func (e ErrComputeRestrictions) Error() string {
	return fmt.Sprintf("unable to compute restrictions: %s", e.Err)
}

// Unwrap returns the wrapped error.
func (e ErrComputeRestrictions) Unwrap() error {
	return e.Err
}

// An ErrInputToken represents a generic input token error.
type ErrInputToken struct {
	Err error
}

func (e ErrInputToken) Error() string {
	return fmt.Sprintf("unable to parse input token: %s", e.Err)
}

// Unwrap returns the wrapped error.
func (e ErrInputToken) Unwrap() error {
	return e.Err
}
