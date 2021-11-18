package providers

import (
	"context"
	"fmt"

	"cloud.google.com/go/compute/metadata"
)

// GCPServiceAccountToken will retrieve the service account
// token and call the midgard library.
func GCPServiceAccountToken(ctx context.Context, audience string) (string, error) {

	return metadata.Get(fmt.Sprintf("instance/service-accounts/default/identity?audience=%s&format=full", audience))
}
