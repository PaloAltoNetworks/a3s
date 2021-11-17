package issuer

import (
	"context"
	"fmt"
	"strings"

	"github.com/coreos/go-oidc/v3/oidc"
	"go.aporeto.io/a3s/pkgs/token"
)

const (
	azureCertURL  = "https://login.microsoftonline.com/common/discovery/keys"
	azureAudience = "https://management.azure.com/"
	azureIssuer   = "https://sts.windows.net/65888785-a93c-4c8f-89eb-d42bf7d03244/"
)

// ErrAzure represents an error that happened
// during operation related to Azure.
type ErrAzure struct {
	Err error
}

func (e ErrAzure) Error() string {
	return fmt.Sprintf("azure error: %s", e.Err)
}

func (e ErrAzure) Unwrap() error {
	return e.Err
}

type azureJWT struct {
	AIO      string `json:"aio"`
	AppID    string `json:"appid"`
	AppIDAcr string `json:"appidacr"`
	IDP      string `json:"idp"`
	OID      string `json:"oid"`
	RH       string `json:"rh"`
	TID      string `json:"tid"`
	UTI      string `json:"uti"`
	XmsMIRID string `json:"xms_mirid"`
}

// An AzureIssuer issues an IdentityToken from
// an existing valid Azure token.
type AzureIssuer struct {
	token *token.IdentityToken
}

// NewAzureIssuer returns a new AzureIssuer.
func NewAzureIssuer() *AzureIssuer {
	return &AzureIssuer{
		token: token.NewIdentityToken(token.Source{
			Type: "azure",
		}),
	}
}

// FromAzureToken computes and verifies the given azure token.
func (c *AzureIssuer) FromAzureToken(ctx context.Context, tokenString string) (err error) {

	ks := oidc.NewRemoteKeySet(ctx, azureCertURL)
	verifier := oidc.NewVerifier(azureIssuer, ks, &oidc.Config{ClientID: azureAudience})
	idt, err := verifier.Verify(ctx, tokenString)
	if err != nil {
		return ErrAzure{Err: err}
	}

	atoken := azureJWT{}
	if err := idt.Claims(&atoken); err != nil {
		return ErrAzure{Err: err}
	}

	c.token.Identity = computeAzureClaims(atoken)

	return nil
}

// Issue returns the IdentityToken.
func (c *AzureIssuer) Issue() *token.IdentityToken {

	return c.token
}

func computeAzureClaims(token azureJWT) []string {

	var out []string

	if token.AIO != "" {
		out = append(out, fmt.Sprintf("aio=%s", token.AIO))
	}

	if token.AppID != "" {
		out = append(out, fmt.Sprintf("appid=%s", token.AppID))
	}

	if token.AppIDAcr != "" {
		out = append(out, fmt.Sprintf("appidacr=%s", token.AppIDAcr))
	}

	if token.IDP != "" {
		out = append(out, fmt.Sprintf("idp=%s", token.IDP))
	}

	if token.OID != "" {
		out = append(out, fmt.Sprintf("oid=%s", token.OID))
	}

	if token.RH != "" {
		out = append(out, fmt.Sprintf("rh=%s", token.RH))
	}

	if token.TID != "" {
		out = append(out, fmt.Sprintf("tid=%s", token.TID))
	}

	if token.UTI != "" {
		out = append(out, fmt.Sprintf("uti=%s", token.UTI))
	}

	if parts := strings.Split(token.XmsMIRID, "/"); len(parts) == 9 {

		if parts[1] == "subscriptions" {
			out = append(out, fmt.Sprintf("subscriptions=%s", parts[2]))
		}

		if parts[3] == "resourcegroups" {
			out = append(out, fmt.Sprintf("resourcegroups=%s", parts[4]))
		}

		if parts[5] == "providers" {
			out = append(out, fmt.Sprintf("providers=%s", parts[6]))
		}

		out = append(out, fmt.Sprintf("providertype=%s", parts[7]))
		out = append(out, fmt.Sprintf("identity=%s", parts[8]))
	}

	return out
}
