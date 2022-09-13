package adapters

import (
	secretmanager "cloud.google.com/go/secretmanager/apiv1"
	"context"
	"fmt"
	"go-secretshelper/pkg/core"
	secretmanagerpb "google.golang.org/genproto/googleapis/cloud/secretmanager/v1"
	"log"
)

// GCPSecretManagerType is the type name as it appears in the configuration
const GCPSecretManagerType = "gcp-secretmanager"

// GCPSecretManagerSpeccProjectID specifies the gcp project id
const GCPSecretManagerSpeccProjectID = "projectID"

// GCPSecretManager is an adapter for GCP Secret Manager
type GCPSecretManager struct {
	log *log.Logger
}

// NewGCPSecretManager returns a new instance of GCPSecretManager.
func NewGCPSecretManager(log *log.Logger) *GCPSecretManager {
	return &GCPSecretManager{log: log}
}

// GCPSecretManagerSpec is the configuration for the GCP Secret Manager adapter
type GCPSecretManagerSpec struct {
	ProjectID string `json:"project_id" validate:"required"`
}

// NewGCPSecretManagerSpec creates a new instance of GCPSecretManagerSpec from a generic map
func NewGCPSecretManagerSpec(in map[interface{}]interface{}) (*GCPSecretManagerSpec, error) {
	var res GCPSecretManagerSpec

	v, ex := in[GCPSecretManagerSpeccProjectID]
	if !ex {
		return nil, fmt.Errorf("%s is required", GCPSecretManagerSpeccProjectID)
	}
	res.ProjectID = v.(string)

	return &res, nil
}

// RetrieveSecret retrieves a secret from GCP Secret Manager.
func (v *GCPSecretManager) RetrieveSecret(ctx context.Context, defaults *core.Defaults,
	vault *core.Vault, secret *core.Secret) (*core.Secret, error) {

	spec, err := NewGCPSecretManagerSpec(vault.Spec)
	if err != nil {
		return nil, err
	}

	client, err := secretmanager.NewClient(ctx)
	if err != nil {
		return nil, fmt.Errorf("GCPSecretManager: failed to create secretmanager client: %v", err)
	}
	defer client.Close()

	req := &secretmanagerpb.AccessSecretVersionRequest{
		Name: fmt.Sprintf("projects/%s/secrets/%s/versions/latest", spec.ProjectID, secret.Name),
	}

	result, err := client.AccessSecretVersion(ctx, req)
	if err != nil {
		return nil, fmt.Errorf("GCPSecretManager: failed to access secret version: %v", err)
	}

	v.log.Printf("GCPSecretManager[%s]: Retrieved secret name=%s\n", vault.Name, secret.Name)

	return &core.Secret{
		RawContent:     []byte(result.Payload.Data),
		RawContentType: "",
		Name:           secret.Name,
		Type:           secret.Type,
		VaultName:      secret.VaultName,
	}, nil
}
