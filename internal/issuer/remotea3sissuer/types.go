package remotea3sissuer

import "fmt"

// ErrRemoteA3S represents an error that happened
// during operation related to Azure.
type ErrRemoteA3S struct {
	Err error
}

func (e ErrRemoteA3S) Error() string {
	return fmt.Sprintf("remote a3s error: %s", e.Err)
}

// Unwrap returns the warped error.
func (e ErrRemoteA3S) Unwrap() error {
	return e.Err
}
