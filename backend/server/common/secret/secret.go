package secret

import (
	"context"

	secretmanager "cloud.google.com/go/secretmanager/apiv1"
	"cloud.google.com/go/secretmanager/apiv1/secretmanagerpb"
	"go.fabra.io/server/common/errors"
)

func FetchSecret(ctx context.Context, name string) (*string, error) {
	client, err := secretmanager.NewClient(ctx)
	if err != nil {
		return nil, errors.Wrap(err, "failed to create secretmanager client")
	}

	req := &secretmanagerpb.AccessSecretVersionRequest{
		Name: name,
	}

	result, err := client.AccessSecretVersion(ctx, req)
	if err != nil {
		return nil, errors.Wrap(err, "failed to access secret version")
	}

	secret := string(result.Payload.Data)
	return &secret, nil
}
