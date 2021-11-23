package api

import (
	"encoding/pem"
	"fmt"
	"net"
	"net/http"
	"regexp"
	"strings"
	"time"

	"go.aporeto.io/elemental"
)

// ValidateDuration valides the given string is a parseable Go duration.
func ValidateDuration(attribute string, duration string) error {

	if duration == "" {
		return nil
	}

	if _, err := time.ParseDuration(duration); err != nil {
		return makeErr("attr", fmt.Sprintf("Attribute '%s' must be a validation duration", attribute))
	}

	return nil
}

// ValidateCIDR validates a CIDR.
func ValidateCIDR(attribute string, network string) error {

	if _, _, err := net.ParseCIDR(network); err == nil {
		return nil
	}

	return makeErr(attribute, fmt.Sprintf("Attribute '%s' must be a CIDR", attribute))
}

// ValidateCIDROptional validates an optional CIDR. It can be empty.
func ValidateCIDROptional(attribute string, network string) error {
	if len(network) == 0 {
		return nil
	}

	return ValidateCIDR(attribute, network)
}

// ValidateCIDRList validates a list of CIDRS.
// The list cannot be empty
func ValidateCIDRList(attribute string, networks []string) error {

	if len(networks) == 0 {
		return makeErr(attribute, fmt.Sprintf("Attribute '%s' must not be empty", attribute))
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

var tagRegex = regexp.MustCompile(`^[^= ]+=.+`)

// ValidateTagsExpression validates an [][]string is a valid tag expression.
func ValidateTagsExpression(attribute string, expression [][]string) error {

	for _, tags := range expression {

		for _, tag := range tags {

			if len([]byte(tag)) >= 1024 {
				return makeErr(attribute, fmt.Sprintf("'%s' must be less than 1024 bytes", tag))
			}
			if !tagRegex.MatchString(tag) {
				return makeErr(attribute, fmt.Sprintf("'%s' must contain at least one '=' symbol separating two valid words", tag))
			}

		}
	}

	return nil
}

// ValidateAuthorizationSubject makes sure api authorization subject is at least secured a bit.
func ValidateAuthorizationSubject(attribute string, subject [][]string) error {

	for i, ands := range subject {

		for _, claim := range ands {

			parts := strings.SplitN(claim, "=", 2)
			if len(parts) != 2 {
				return makeErr(attribute, fmt.Sprintf("Subject claims '%s' on line %d is an invalid tag", claim, i+1))
			}
			if parts[1] == "" {
				return makeErr(attribute, fmt.Sprintf("Subject claims '%s' on line %d has no value", claim, i+1))
			}
		}
	}

	return nil
}

// ValidatePEM validates a string contains a PEM.
func ValidatePEM(attribute string, pemdata string) error {

	if pemdata == "" {
		return nil
	}

	var i int
	var block *pem.Block
	rest := []byte(pemdata)

	for {
		block, rest = pem.Decode(rest)

		if block == nil {
			return makeErr(attribute, fmt.Sprintf("Unable to decode PEM number %d", i))
		}

		if len(rest) == 0 {
			return nil
		}
		i++
	}
}

// ValidateIssue validates a whole issue object.
func ValidateIssue(iss *Issue) error {

	switch iss.SourceType {
	case IssueSourceTypeA3S:
		if iss.InputA3S == nil {
			return makeErr("inputA3S", "You must set inputA3S for the requested sourceType")
		}
	case IssueSourceTypeRemoteA3S:
		if iss.InputRemoteA3S == nil {
			return makeErr("inputRemoteA3S", "You must set inputRemoteA3S for the requested sourceType")
		}
	case IssueSourceTypeAWS:
		if iss.InputAWS == nil {
			return makeErr("inputAWS", "You must set inputAWS for the requested sourceType")
		}
	case IssueSourceTypeLDAP:
		if iss.InputLDAP == nil {
			return makeErr("inputLDAP", "You must set inputLDAP for the requested sourceType")
		}
	case IssueSourceTypeGCP:
		if iss.InputGCP == nil {
			return makeErr("inputGCP", "You must set inputCGP for the requested sourceType")
		}
	case IssueSourceTypeAzure:
		if iss.InputAzure == nil {
			return makeErr("inputAzure", "You must set inputAzure for the requested sourceType")
		}
	}

	return nil
}

func makeErr(attribute string, message string) elemental.Error {

	err := elemental.NewError(
		"Validation Error",
		message,
		"a3s",
		http.StatusUnprocessableEntity,
	)

	if attribute != "" {
		err.Data = map[string]interface{}{"attribute": attribute}
	}

	return err
}
