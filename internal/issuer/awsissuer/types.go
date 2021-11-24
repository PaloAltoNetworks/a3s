package awsissuer

import "fmt"

// A ErrAWSSTS represents an error during interactions with AWS.
type ErrAWSSTS struct {
	Err error
}

func (e ErrAWSSTS) Error() string {
	return fmt.Sprintf("aws error: %s", e.Err)
}

// Unwrap returns the wrapped error.
func (e ErrAWSSTS) Unwrap() error {
	return e.Err
}
