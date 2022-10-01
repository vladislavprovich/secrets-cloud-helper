package test

import (
	"github.com/vladislavprovich/secrets-cloud-helper/pkg/adapters"
	"testing"
)

func TestAzureKeyVaultSpec(t *testing.T) {
	m := make(map[interface{}]interface{})

	m["url"] = "not-a$url"
	_, err := adapters.NewAzureKeyVaultSpec(m)
	if err == nil {
		t.Error("Expected error for invalid url")
	}

	m["url"] = "http://not.secure.url.com/"
	_, err = adapters.NewAzureKeyVaultSpec(m)
	if err == nil {
		t.Error("Expected error for invalid url")
	}

}
