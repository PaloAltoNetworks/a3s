package httpissuer

import "fmt"

// ErrHTTP represents an error that happened
// during operation related to HTTP source call.
type ErrHTTP struct {
	Err error
}

func (e ErrHTTP) Error() string {
	return fmt.Sprintf("http error: %s", e.Err)
}

// Unwrap returns the warped error.
func (e ErrHTTP) Unwrap() error {
	return e.Err
}

// ErrHTTPResponse represents an error returned by the
// remote server.
type ErrHTTPResponse struct {
	Err error
}

func (e ErrHTTPResponse) Error() string {
	return fmt.Sprintf("http response error: %s", e.Err)
}

// Unwrap returns the warped error.
func (e ErrHTTPResponse) Unwrap() error {
	return e.Err
}
