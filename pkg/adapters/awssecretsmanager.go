package adapters

import (
	"context"
	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/aws/session"
	"github.com/aws/aws-sdk-go/service/secretsmanager"
	"go-secretshelper/pkg/core"
	"log"
)

// AWSSecretsManagerType is the type of this adapter, to be used in configuration files
const AWSSecretsManagerType = "aws-secretsmanager"

// AWSSecretsManager is a VaultAccessPort for the AWS Secrets Manager service.
type AWSSecretsManager struct {
	log *log.Logger
}

// NewAWSSecretsManager returns a new AWSSecretsManager.
func NewAWSSecretsManager(log *log.Logger) *AWSSecretsManager {
	return &AWSSecretsManager{log: log}
}

// AWSSecretsManagerSpec specifies the configuration for an AWSSecretsManager.
type AWSSecretsManagerSpec struct {
	Region string `yaml:"region"`
}

// NewAWSSecretsManagerSpec returns a new AWSSecretsManagerSpec.
func NewAWSSecretsManagerSpec(in map[interface{}]interface{}) (*AWSSecretsManagerSpec, error) {
	var res AWSSecretsManagerSpec

	v, ex := in["region"]
	if ex {
		res.Region = v.(string)
	}

	return &res, nil
}

// RetrieveSecret retrieves a secret from the aws' secrets manager
func (v *AWSSecretsManager) RetrieveSecret(ctx context.Context, defaults *core.Defaults,
	vault *core.Vault, secret *core.Secret) (*core.Secret, error) {

	spec, err := NewAWSSecretsManagerSpec(vault.Spec)
	if err != nil {
		return nil, err
	}

	config := aws.NewConfig()
	if spec.Region != "" {
		config.Region = &spec.Region
	}

	session, err := session.NewSession(config)
	if err != nil {
		return nil, err
	}

	svc := secretsmanager.New(session)

	/*secretDescription, err := svc.DescribeSecret(&secretsmanager.DescribeSecretInput{
		SecretId:     aws.String(secret.Name),
	})
	if err != nil {
		return nil, err
	}
	v.log.Printf("AWSSecretsManager[%s]: %#v\n", vault.Name, *secretDescription)*/

	result, err := svc.GetSecretValue(&secretsmanager.GetSecretValueInput{
		SecretId:     aws.String(secret.Name),
		VersionStage: aws.String("AWSCURRENT"),
	})
	if err != nil {
		return nil, err
	}

	v.log.Printf("AWSSecretsManager[%s]: Retrieved secret name=%s, arn=%s, v=%s\n", vault.Name, secret.Name, *result.ARN, *result.VersionId)

	return &core.Secret{
		RawContent:     []byte(*result.SecretString),
		RawContentType: "",
		Name:           secret.Name,
		Type:           secret.Type,
		VaultName:      secret.VaultName,
	}, nil
}
