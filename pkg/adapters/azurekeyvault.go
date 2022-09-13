package adapters

import (
	"context"
	"fmt"
	"github.com/Azure/azure-sdk-for-go/profiles/preview/keyvault/keyvault"
	"github.com/Azure/azure-sdk-for-go/services/keyvault/auth"
	"go-secretshelper/pkg/core"
	"log"
	"net/url"
)

// AzureKeyVaultType is the type name for azure key vaults
const AzureKeyVaultType = "azure-key-vault"

// AzureKeyVault is a core.VaultAccessorPort which pulls secrets from a Key Vault within an Azure subscription
type AzureKeyVault struct {
	log *log.Logger
}

// NewAzureKeyVault creates a new age vault
func NewAzureKeyVault(l *log.Logger) *AzureKeyVault {
	return &AzureKeyVault{
		log: l,
	}
}

// AzureKeyVaultSpec describes access to the vault
type AzureKeyVaultSpec struct {
	URL string `yaml:"url"`
}

// NewAzureKeyVaultSpec creates a new vault spec from the generic interface map
func NewAzureKeyVaultSpec(in map[interface{}]interface{}) (AzureKeyVaultSpec, error) {
	var res AzureKeyVaultSpec

	v, ex := in["url"]
	if ex {
		res.URL = v.(string)

		// check if url is valid
		u, err := url.Parse(res.URL)
		if err != nil {
			return res, fmt.Errorf("invalid url: %s", err)
		}
		if u.Scheme != "https" || u.Host == "" {
			return res, fmt.Errorf("invalid url: %s", err)
		}
	}

	return res, nil
}

// RetrieveSecret decodes both identity and age file according to vault.Spec and
// reads the secret.
func (v *AzureKeyVault) RetrieveSecret(ctx context.Context, defaults *core.Defaults,
	vault *core.Vault, secret *core.Secret) (*core.Secret, error) {

	// parse age vault spec
	spec, err := NewAzureKeyVaultSpec(vault.Spec)
	if err != nil {
		return nil, err
	}

	url := spec.URL
	if len(url) == 0 {
		// compose url from key vault name
		url = fmt.Sprintf("https://%s.vault.azure.net/", vault.Name)
		v.log.Printf("AzureKeyVault: using url: %s", url)
	}

	client := keyvault.New()

	// see https://docs.microsoft.com/de-de/azure/developer/go/azure-sdk-authorization
	// see also https://github.com/Azure-Samples/azure-sdk-for-go-samples/blob/master/internal/iam/authorizers.go
	authorizer, err := auth.NewAuthorizerFromEnvironment()
	if err != nil {
		return nil, err
	}
	client.Authorizer = authorizer

	secretBundle, err := client.GetSecret(ctx,
		url,
		secret.Name,
		"")
	if err != nil {
		return nil, err
	}

	//secretBundeJSON := secretBundle.MarshalJSON()
	//v.log.Printf("%s\n", string(secretBundeJSON))

	if secretBundle.Value == nil || secretBundle.ID == nil {
		return nil, fmt.Errorf("AzureKeyVault[%s]: secret %s was empty or malformed", vault.Name, secret.Name)
	}
	// set content type if transmitted
	var ct = ""
	if secretBundle.ContentType != nil {
		ct = *secretBundle.ContentType
	}

	v.log.Printf("AzureKeyVault[%s]: Retrieved secret name=%s, id=%s, ct=%s", vault.Name, secret.Name, *secretBundle.ID, ct)

	return &core.Secret{
		RawContent:     []byte(*secretBundle.Value),
		RawContentType: ct,
		Name:           secret.Name,
		Type:           secret.Type,
		VaultName:      secret.VaultName,
	}, nil
}
