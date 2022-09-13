package core

import (
	"fmt"
	"github.com/drone/envsubst"
	"github.com/go-playground/validator/v10"
	"gopkg.in/yaml.v2"
	"io"
	"os"
	"strings"
)

// Config is the main configuration struct.
type Config struct {
	// Defaults contains settings valid for some or all other sections
	Defaults Defaults `yaml:"defaults" validate:""`

	// Vaults define all the vaults where secrets are pulled from
	Vaults Vaults `yaml:"vaults" validate:"required,dive"`

	// Secrets define the name and location of the secrets
	Secrets Secrets `yaml:"secrets" validate:"required,dive"`

	// Transformations define optional transformation steps for secrets
	Transformations Transformations `yaml:"transformations" validate:""`

	// Sinks define the output sinks for the (transformed) secrets
	Sinks Sinks `yaml:"sinks" validate:"required,dive"`
}

// NewDefaultConfig returns a configuration struct with valid default settings
func NewDefaultConfig() *Config {
	return &Config{}
}

// NewConfig is the default way of reading configuration from yaml stream
func NewConfig(in io.Reader) (*Config, error) {
	yamlDec := yaml.NewDecoder(in)

	res := NewDefaultConfig()
	if err := yamlDec.Decode(res); err != nil {
		return res, err
	}

	return res, nil
}

// NewConfigWithEnvSubst works like NewConfig with environment variable substitution
func NewConfigWithEnvSubst(in io.Reader) (*Config, error) {
	// read all of in and run substitution on it
	buf := new(strings.Builder)
	_, err := io.Copy(buf, in)
	if err != nil {
		return nil, err
	}

	inSubst, err := envsubst.EvalEnv(buf.String())
	if err != nil {
		return nil, err
	}

	res := NewDefaultConfig()
	if err := yaml.Unmarshal([]byte(inSubst), res); err != nil {
		return res, err
	}

	return res, nil
}

// NewConfigFromFile creates a configuration from yaml file
func NewConfigFromFile(fileName string, withEnvSubst bool) (*Config, error) {
	f, err := os.Open(fileName)
	if err != nil {
		return nil, err
	}

	if withEnvSubst {
		return NewConfigWithEnvSubst(f)
	}
	return NewConfig(f)
}

// IsVarDefined checks if given variable name is defined, either in
// secrets or as the result of a transformation step
func (c *Config) IsVarDefined(varName string) bool {
	for _, secret := range c.Secrets {
		if secret.Name == varName {
			return true
		}
	}

	for _, transformation := range c.Transformations {
		if transformation.Output == varName {
			return true
		}
	}

	return false
}

func validateSecretType(fl validator.FieldLevel) bool {
	for _, secretType := range ValidSecretTypes() {
		if secretType == fl.Field().String() {
			return true
		}
	}
	return false
}

// Validate validates a configuration using the validator and
// additional cross checks.
func (c *Config) Validate(f Factory) error {
	v := validator.New()
	v.RegisterValidation("valid-secret-type", validateSecretType)

	st := make(map[string]struct{})
	for _, e := range f.SinkTypes() {
		st[e] = struct{}{}
	}
	vat := make(map[string]struct{})
	for _, e := range f.VaultAccessorTypes() {
		vat[e] = struct{}{}
	}
	tt := make(map[string]struct{})
	for _, e := range f.TransformationTypes() {
		tt[e] = struct{}{}
	}

	for _, vault := range c.Vaults {
		if err := v.Struct(vault); err != nil {
			return err
		}

		if _, ex := vat[vault.Type]; !ex {
			return fmt.Errorf("unknown vault type: %s in vault: %s", vault.Type, vault.Name)
		}
	}

	if err := c.validateSecrets(f); err != nil {
		return err
	}

	if err := c.validateTransformations(f); err != nil {
		return err
	}

	if err := c.validateSinks(f); err != nil {
		return err
	}

	return nil
}

func (c *Config) validateSecrets(f Factory) error {
	v := validator.New()
	v.RegisterValidation("valid-secret-type", validateSecretType)

	for _, secret := range c.Secrets {
		if err := v.Struct(secret); err != nil {
			return err
		}

		if v := c.Vaults.GetVaultByName(secret.VaultName); v == nil {
			return fmt.Errorf("invalid vault %s referenced in secret %s", secret.VaultName, secret.Name)
		}
	}

	return nil
}

func (c *Config) validateTransformations(f Factory) error {
	v := validator.New()

	tt := make(map[string]struct{})
	for _, e := range f.TransformationTypes() {
		tt[e] = struct{}{}
	}

	for _, transformation := range c.Transformations {
		if err := v.Struct(transformation); err != nil {
			return err
		}

		if _, ex := tt[transformation.Type]; !ex {
			return fmt.Errorf("unknown transformation type: %s", transformation.Type)
		}

		// all input variables have to be defined
		for _, inputVar := range transformation.Input {
			if !c.IsVarDefined(inputVar) {
				return fmt.Errorf("unknown input variable: %s", inputVar)
			}
		}
	}

	return nil
}

func (c *Config) validateSinks(f Factory) error {
	v := validator.New()

	st := make(map[string]struct{})
	for _, e := range f.SinkTypes() {
		st[e] = struct{}{}
	}

	for _, sink := range c.Sinks {
		if err := v.Struct(sink); err != nil {
			return err
		}

		if _, ex := st[sink.Type]; !ex {
			return fmt.Errorf("unknown sink type: %s", sink.Type)
		}

		if !c.IsVarDefined(sink.Var) {
			return fmt.Errorf("invalid variable %s referenced in a sink", sink.Var)
		}
	}

	return nil
}
