package test

import (
	"github.com/vladislavprovich/secrets-cloud-helper/pkg/adapters"
	"testing"
)

func TestAWSSecretsManagerSpec(t *testing.T) {
	m := make(map[interface{}]interface{})

	m["region"] = "us-east-1"
	spec, err := adapters.NewAWSSecretsManagerSpec(m)
	if err != nil {
		t.Error("Unexpected error")
	}
	if spec.Region != m["region"] {
		t.Error("Expected region to be set")
	}

}
