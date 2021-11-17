package azureissuer

import (
	"context"
	"fmt"
	"strings"

	"github.com/coreos/go-oidc/v3/oidc"
	"go.aporeto.io/a3s/pkgs/token"
)

const (
	azureJWTCertURL  = "https://login.microsoftonline.com/common/discovery/keys"
	azureJWTAudience = "https://management.azure.com/"
	azureJWTIssuer   = "https://sts.windows.net/65888785-a93c-4c8f-89eb-d42bf7d03244/"
)

func New(ctx context.Context, tokenString string) (token.Issuer, error) {

	c := newAzureIssuer()
	if err := c.fromToken(ctx, tokenString); err != nil {
		return nil, err
	}

	return c, nil
}

type azureIssuer struct {
	token *token.IdentityToken
}

func newAzureIssuer() *azureIssuer {
	return &azureIssuer{
		token: token.NewIdentityToken(token.Source{
			Type: "azure",
		}),
	}
}

// Issue returns the IdentityToken.
func (c *azureIssuer) Issue() *token.IdentityToken {

	return c.token
}

func (c *azureIssuer) fromToken(ctx context.Context, tokenString string) (err error) {

	ks := oidc.NewRemoteKeySet(ctx, azureJWTCertURL)
	verifier := oidc.NewVerifier(azureJWTIssuer, ks, &oidc.Config{ClientID: azureJWTAudience})
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
