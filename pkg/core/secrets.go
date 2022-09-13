package core

import "fmt"

// Secrets is an array of Secret structs
type Secrets []*Secret

// Secret defines a named secrets, referenced in a named Vault
type Secret struct {
	// Name of the secret within the vault
	Name string `yaml:"name" validate:"required"`

	// VaultName specifies in which vault the secret is stored
	VaultName string `yaml:"vault" validate:"required"`

	// Type of secret
	Type string `yaml:"type" validate:"required,valid-secret-type"`

	// RawContent contains the secret
	RawContent []byte

	// RawContentType is the content-type of RawContent
	RawContentType string
}

// ValidSecretTypes is a list of valid types of secrets that can be queried from vaults
func ValidSecretTypes() []string {
	return []string{
		"secret", // secret is the default, text/plain retrievable secret, e.g. a password
	}
}

// String returns a string representation of a secret
func (s Secret) String() string {
	set := false
	if len(s.RawContent) > 0 {
		set = true
	}
	return fmt.Sprintf("Secret:[name=%s, Type=%s, set=%t, content-type=%s]", s.Name, s.Type, set, s.RawContentType)
}
