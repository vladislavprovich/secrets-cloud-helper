package adapters

import (
	"github.com/spf13/afero"
	"go-secretshelper/pkg/core"
	"log"
)

// BuiltinFactory is able to create all builtin components of
// the adapters package
type BuiltinFactory struct {
	log *log.Logger
	fs  afero.Fs
}

// NewBuiltinFactory creates the Builtin Factory
func NewBuiltinFactory(log *log.Logger, fs afero.Fs) *BuiltinFactory {
	return &BuiltinFactory{
		log: log,
		fs:  fs,
	}
}

// SinkTypes returns valid sink types
func (f *BuiltinFactory) SinkTypes() []string {
	return []string{
		FileSinkType,
	}
}

// TransformationTypes returns valid transformation types
func (f *BuiltinFactory) TransformationTypes() []string {
	return []string{
		TemplateTransformationType,
		AgeEncryptTransformationType,
		JQTransformationType,
	}
}

// VaultAccessorTypes returns valid vault types
func (f *BuiltinFactory) VaultAccessorTypes() []string {
	return []string{
		AgeVaultType,
		AzureKeyVaultType,
		AWSSecretsManagerType,
		GCPSecretManagerType,
	}
}

// NewRepository creates a new repository
func (f *BuiltinFactory) NewRepository() core.Repository {
	return NewBuiltinRepository()
}

// NewSinkWriter creates a new sink writer for a supported type
func (f *BuiltinFactory) NewSinkWriter(sinkType string) core.SinkWriterPort {
	switch sinkType {
	case FileSinkType:
		return NewFileSink(f.log, f.fs)
	}
	return nil
}

// NewTransformation creates a new transformation for a supported type
func (f *BuiltinFactory) NewTransformation(transformationType string) core.TransformationPort {
	switch transformationType {
	case TemplateTransformationType:
		return NewTemplateTransformation(f.log)
	case AgeEncryptTransformationType:
		return NewAgeEncryptTransformation(f.log)
	case JQTransformationType:
		return NewJQTransformation(f.log)
	}
	return nil
}

// NewVaultAccessor creates a new vault accessor for a supported type
func (f *BuiltinFactory) NewVaultAccessor(vaultType string) core.VaultAccessorPort {
	switch vaultType {
	case AgeVaultType:
		return NewAgeVault(f.log, f.fs)
	case AzureKeyVaultType:
		return NewAzureKeyVault(f.log)
	case AWSSecretsManagerType:
		return NewAWSSecretsManager(f.log)
	case GCPSecretManagerType:
		return NewGCPSecretManager(f.log)
	}
	return nil
}
