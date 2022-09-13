package test

import (
	"go-secretshelper/pkg/adapters"
	"testing"
)

func TestGCPSecretManagerSpec(t *testing.T) {
	m := make(map[interface{}]interface{})

	m["projectID"] = "someID"
	spec, err := adapters.NewGCPSecretManagerSpec(m)
	if err != nil {
		t.Error("Unexpected error")
	}
	if spec.ProjectID != m["projectID"] {
		t.Error("Expected projectID to be set")
	}

}
