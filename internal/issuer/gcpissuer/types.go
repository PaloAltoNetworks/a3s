package gcpissuer

import (
	"fmt"

	"github.com/golang-jwt/jwt/v4"
)

// ErrGCP represents an error that happened
// during operations related to GCP.
type ErrGCP struct {
	Err error
}

func (e ErrGCP) Error() string {
	return fmt.Sprintf("gcp error: %s", e.Err)
}

// Unwrap returns the warped error.
func (e ErrGCP) Unwrap() error {
	return e.Err
}

type gcpJWT struct {
	Google struct {
		ComputeEngine struct {
			ProjectID     string `json:"project_id"`
			ProjectNumber int    `json:"project_number"`
			Zone          string `json:"zone"`
			InstanceID    string `json:"instance_id"`
			InstanceName  string `json:"instance_name"`
		} `json:"compute_engine"`
	} `json:"google"`
	Email         string `json:"email"`
	EmailVerified bool   `json:"email_verified"`
	jwt.RegisteredClaims
}
