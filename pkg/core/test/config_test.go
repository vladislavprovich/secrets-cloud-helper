package test

import (
	"fmt"
	"github.com/go-playground/validator/v10"
	"github.com/golang/mock/gomock"
	"go-secretshelper/pkg/core"
	"os"
	"strings"
	"testing"
)

func TestConfig(t *testing.T) {
	if core.NewDefaultConfig() == nil {
		t.Error("Must provide a default config")
	}

	_, err := core.NewConfigFromFile("no.such.config.json", false)
	if err == nil {
		t.Errorf("Expected error, got nil")
	}

	cfg, err := core.NewConfigFromFile("../../../tests/fixtures/fixture-1.yaml", false)
	if err != nil {
		t.Errorf("Expected err=nil, got err=%s", err)
	}
	if cfg == nil {
		t.Error("Expected config result, got nil")
	}

	inp := `
vaults:
  - name: kv1
    type: mock

secrets:
  - type: secret
    vault: kv1
    name: test

sinks:
  - type: mock
    var: test
    spec:
      path: ./test.txt
      mode: 400
      user: 1000
`

	cfg, err = core.NewConfig(strings.NewReader(inp))
	if err != nil {
		t.Errorf("Expected err=nil, got err=%s", err)
	}
	if cfg == nil {
		t.Error("Expected config result, got nil")
	}

	inpWithEnv := `
vaults:
  - name: kv1
    type: mock

secrets:
  - type: secret
    vault: ${vaultname}
    name: test

sinks:
  - type: ${nonex}
    var: ${varname:=test}
    spec:
      path: ./test.txt
      mode: 400
      user: 1000
`

	os.Setenv("vaultname", "kv1")

	cfg, err = core.NewConfigWithEnvSubst(strings.NewReader(inpWithEnv))
	if err != nil {
		t.Errorf("Expected err=nil, got err=%s", err)
	}
	if cfg == nil {
		t.Error("Expected config result, got nil")
	}

	if cfg.Secrets[0].VaultName != "kv1" {
		t.Errorf("Expected vault=kv1, got vault=%s", cfg.Secrets[0].VaultName)
	}
	if cfg.Sinks[0].Var != "test" {
		t.Errorf("Expected var=test, got var=%s", cfg.Sinks[0].Var)
	}
	if cfg.Sinks[0].Type != "" {
		t.Errorf("Expected empty type, got type=%s", cfg.Sinks[0].Type)
	}
}

func DumpValidationErrors(err error) {
	if err != nil {

		if _, ok := err.(*validator.InvalidValidationError); ok {
			fmt.Println(err)
			return
		}

		for _, err := range err.(validator.ValidationErrors) {

			fmt.Println(err.Namespace())
			fmt.Println(err.Field())
			fmt.Println(err.StructNamespace())
			fmt.Println(err.StructField())
			fmt.Println(err.Tag())
			fmt.Println(err.ActualTag())
			fmt.Println(err.Kind())
			fmt.Println(err.Type())
			fmt.Println(err.Value())
			fmt.Println(err.Param())
			fmt.Println()
		}

		return
	}
}

func TestValidation(t *testing.T) {
	cfg, err := core.NewConfigFromFile("../../../tests/fixtures/fixture-1.yaml", false)
	if err != nil {
		t.Errorf("Expected err=null, got err=%s", err)
	}
	if cfg == nil {
		t.Error("Expected config result, got nil")
	}

	mockCtrl := gomock.NewController(t)
	defer mockCtrl.Finish()

	mf := NewMockFactory(mockCtrl, t)

	err = cfg.Validate(mf)
	if err != nil {
		t.Errorf("Expected nil got err=%#v", err)
		DumpValidationErrors(err)
	}

	// simple validation (missing elements)

	cfg = &core.Config{
		Vaults: []*core.Vault{
			{},
		},
	}
	err = cfg.Validate(mf)
	if err == nil {
		t.Errorf("Expected validation error, got nil")
	}

	// referential validation
	cfg = &core.Config{
		Vaults: []*core.Vault{
			{
				Name: "a",
				Type: "nonex",
				Spec: core.VaultSpec{},
			},
		},
		Secrets: []*core.Secret{
			{
				Name:      "b",
				VaultName: "nonex", // this vault is not defined above
				Type:      "secret",
			},
		},
	}
	err = cfg.Validate(mf)
	if err == nil {
		t.Errorf("Expected validation error, got nil")
	}

}

func TestValidationForTransformation(t *testing.T) {
	mockCtrl := gomock.NewController(t)
	defer mockCtrl.Finish()

	mf := NewMockFactory(mockCtrl, t)

	// check if variables are defined
	cfg, err := core.NewConfig(strings.NewReader(`
vaults:
  - name: kv1
    type: mock

secrets:
  - type: secret
    vault: kv1
    name: test

transformations:
  - in:
    - test
    out: testout
    type: mock
  - in:
    - testout
    out: devnull
    type: mock

sinks:
  - type: mock
    var: test
    spec:
      path: ./test.txt
      mode: 400
      user: 1000
`))
	if err != nil {
		t.Errorf("Expected nil got err=%#v", err)
	}

	err = cfg.Validate(mf)
	if err != nil {
		t.Errorf("Expected nil got err=%s", err)
	}

	// test fur unknown input var
	cfg, err = core.NewConfig(strings.NewReader(`
vaults:
  - name: kv1
    type: mock

secrets:
  - type: secret
    vault: kv1
    name: test

transformations:
  - in:
    - nonex
    out: testout
    type: mock

sinks:
  - type: mock
    var: test
    spec:
      path: ./test.txt
      mode: 400
      user: 1000
`))
	if err != nil {
		t.Errorf("Expected nil got err=%#v", err)
	}

	err = cfg.Validate(mf)
	if err == nil {
		t.Errorf("Expected validation error, got nil")
	}

	// test fur empty input var
	cfg, err = core.NewConfig(strings.NewReader(`
vaults:
  - name: kv1
    type: mock

secrets:
  - type: secret
    vault: kv1
    name: test

transformations:
  - in:
    out: testout
    type: mock

sinks:
  - type: mock
    var: test
    spec:
      path: ./test.txt
      mode: 400
      user: 1000
`))
	if err != nil {
		t.Errorf("Expected nil got err=%#v", err)
	}

	err = cfg.Validate(mf)
	if err == nil {
		t.Errorf("Expected validation error, got nil")
	}
}

func TestValidationForTypes(t *testing.T) {
	mockCtrl := gomock.NewController(t)
	defer mockCtrl.Finish()

	mf := NewMockFactory(mockCtrl, t)

	// check if variables are defined
	cfg, err := core.NewConfig(strings.NewReader(`
vaults:
  - name: kv1
    type: mock

secrets:
  - type: noSuchSecretType
    vault: kv1
    name: test

sinks:
  - type: mock
    var: test
    spec:
      path: ./test.txt
      mode: 400
      user: 1000
`))
	if err != nil {
		t.Errorf("Expected nil got err=%#v", err)
	}

	err = cfg.Validate(mf)
	if err == nil {
		t.Errorf("Expected validation error, got nil")
	}

}
