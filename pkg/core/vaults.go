package core

import (
	"context"
	"fmt"
)

// Vaults is an array of Vault structs
type Vaults []*Vault

// Vault defines one source for secrets
type Vault struct {
	// Name of the vault, referenced by secrets
	Name string `yaml:"name" validate:"required"`

	// Type of vault
	Type string `yaml:"type" validate:"required"`

	// Detailed specification
	Spec VaultSpec `yaml:"spec" validate:""`
}

// VaultSpec declares details of how to connect to the vault
type VaultSpec map[interface{}]interface{}

// String creates a string representation for a vault
func (v Vault) String() string {
	return fmt.Sprintf("Vault:[Name=%s, Type=%s]", v.Name, v.Type)
}

// GetVaultByName searches for a single named vault from an array
func (vaults *Vaults) GetVaultByName(name string) *Vault {
	for _, vault := range *vaults {
		if vault.Name == name {
			return vault
		}
	}
	return nil
}

// VaultAccessorPort is able to pull secrets from a Vault
type VaultAccessorPort interface {

	// RetrieveSecret retrieves a secret from given vault
	RetrieveSecret(context.Context, *Defaults, *Vault, *Secret) (*Secret, error)
}
