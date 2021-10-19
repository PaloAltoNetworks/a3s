package api

import (
	"fmt"
	"net"
	"net/http"

	"go.aporeto.io/elemental"
)

// ValidateCIDROptional validates an optional CIDR. It can be empty.
func ValidateCIDROptional(attribute string, network string) error {
	if len(network) == 0 {
		return nil
	}

	return ValidateCIDR(attribute, network)
}

// ValidateCIDR validates a CIDR.
func ValidateCIDR(attribute string, network string) error {

	if _, _, err := net.ParseCIDR(network); err == nil {
		return nil
	}

	return makeValidationError(attribute, fmt.Sprintf("Attribute '%s' must be a CIDR", attribute))
}

// ValidateCIDRList validates a list of CIDRS.
// The list cannot be empty
func ValidateCIDRList(attribute string, networks []string) error {

	if len(networks) == 0 {
		return makeValidationError(attribute, fmt.Sprintf("Attribute '%s' must not be empty", attribute))
	}

	return ValidateCIDRListOptional(attribute, networks)
}

// ValidateCIDRListOptional validates a list of CIDRs.
// It can be empty.
func ValidateCIDRListOptional(attribute string, networks []string) error {

	for _, network := range networks {
		if err := ValidateCIDR(attribute, network); err != nil {
			return err
		}
	}

	return nil
}

func makeValidationError(attribute string, message string) elemental.Error {

	err := elemental.NewError(
		"Validation Error",
		message,
		"gaia",
		http.StatusUnprocessableEntity,
	)

	if attribute != "" {
		err.Data = map[string]interface{}{"attribute": attribute}
	}

	return err
}
