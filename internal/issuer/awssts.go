package issuer

import (
	"fmt"
	"strings"

	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/aws/arn"
	"github.com/aws/aws-sdk-go/aws/credentials"
	"github.com/aws/aws-sdk-go/aws/session"
	"github.com/aws/aws-sdk-go/service/sts"
	"go.aporeto.io/a3s/pkgs/token"
)

// A ErrAWSSTS represents an error during interactions
// with AWS.
type ErrAWSSTS struct {
	Err error
}

func (e ErrAWSSTS) Error() string {
	return fmt.Sprintf("aws error: %s", e.Err)
}

func (e ErrAWSSTS) Unwrap() error {
	return e.Err
}

// An AWSSTSIssuer issues identity token from an AWS sts token.
type AWSSTSIssuer struct {
	token *token.IdentityToken
}

// NewAWSSTSIssuer returns a new AWSSTSIssuer.
func NewAWSSTSIssuer() *AWSSTSIssuer {

	return &AWSSTSIssuer{
		token: token.NewIdentityToken(token.Source{
			Type: "awssts",
		}),
	}
}

// FromSTS populates the identity token using the given aws STS token information.
func (c *AWSSTSIssuer) FromSTS(ID string, secret string, token string) error {

	config := &aws.Config{
		Credentials:                   credentials.NewStaticCredentials(ID, secret, token),
		CredentialsChainVerboseErrors: aws.Bool(true),
	}

	session, err := session.NewSession(config)
	if err != nil {
		return ErrAWSSTS{Err: fmt.Errorf("unable to start aws session: %w", err)}
	}

	out, err := sts.New(session).GetCallerIdentity(&sts.GetCallerIdentityInput{})
	if err != nil {
		return ErrAWSSTS{Err: fmt.Errorf("unable to retrieve aws identity: %w", err)}
	}

	if out.Account == nil {
		return ErrAWSSTS{Err: fmt.Errorf("invalid aws response: missing account")}
	}

	if out.Arn == nil {
		return ErrAWSSTS{Err: fmt.Errorf("invalid aws response: missing arn")}
	}

	if out.UserId == nil {
		return ErrAWSSTS{Err: fmt.Errorf("invalid aws response: missing user id")}
	}

	parn, err := arn.Parse(*out.Arn)
	if err != nil {
		return ErrAWSSTS{Err: fmt.Errorf("unable to parse arn '%s': %w", *out.Arn, err)}
	}

	c.token.Identity = computeSTSClaims(out, parn)

	return nil
}

// Issue returns the token.IdentityToken.
func (c *AWSSTSIssuer) Issue() *token.IdentityToken {
	return c.token
}

func computeSTSClaims(out *sts.GetCallerIdentityOutput, parn arn.ARN) (claims []string) {

	claims = []string{
		fmt.Sprintf("account=%s", *out.Account),
		fmt.Sprintf("arn=%s", *out.Arn),
		fmt.Sprintf("userid=%s", *out.UserId),
	}

	if v := parn.Partition; v != "" {
		claims = append(claims, fmt.Sprintf("partition=%s", v))
	}

	if v := parn.Service; v != "" {
		claims = append(claims, fmt.Sprintf("service=%s", v))
	}

	if v := parn.Resource; v != "" {
		claims = append(claims, fmt.Sprintf("resource=%s", v))
		if strings.Count(v, "/") == 2 {
			parts := strings.SplitN(v, "/", 3)
			if len(parts) == 3 {
				claims = append(claims, fmt.Sprintf("resourcetype=%s", parts[0]))
				claims = append(claims, fmt.Sprintf("rolename=%s", parts[1]))
				claims = append(claims, fmt.Sprintf("rolesessionname=%s", parts[2]))
			}
		}
	}

	return claims
}
