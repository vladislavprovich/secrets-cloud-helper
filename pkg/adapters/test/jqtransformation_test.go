package test

import (
	"context"
	"encoding/json"
	"go-secretshelper/pkg/adapters"
	"go-secretshelper/pkg/core"
	"io/ioutil"
	"log"
	"reflect"
	"testing"
)

func TestJQTransformationSpec(t *testing.T) {
	q := ".foo | .."
	in := core.TransformationSpec{
		"q":           q,
		"contentType": "text/plain",
		"raw":         true,
	}

	spec, err := adapters.NewJQTransformationSpec(in)
	if err != nil {
		t.Errorf("Error creating spec: %s", err)
	}
	if spec.Query == "" {
		t.Errorf("Expected query to be present, got nil")
	}
	if spec.Raw != true {
		t.Errorf("Expected raw flag to be present, got nil")
	}

	spec, err = adapters.NewJQTransformationSpec(core.TransformationSpec{})
	if err == nil {
		t.Errorf("Expected error creating spec")
	}

}

func TestJQTransformation(t *testing.T) {
	secrets := &core.Secrets{
		{
			Name: "s1",
			RawContent: []byte(`{
	"first": "s3cr3t",
	"then": [ "the", "rest" ]
}`),
		},
	}

	transformation := &core.Transformation{
		Input:  []string{"s1"},
		Output: "result",
		Type:   "jq",
		Spec: core.TransformationSpec{
			"q":   ".first",
			"raw": true,
		},
	}

	tt := adapters.NewJQTransformation(log.New(ioutil.Discard, "***", 0))
	s, err := tt.ProcessSecret(context.TODO(), &core.Defaults{}, secrets, transformation)
	if err != nil {
		t.Errorf("Error processing jq transformation: %s", err)
	}

	if s == nil {
		t.Errorf("Expected secret to be present, got nil")
	}

	if !reflect.DeepEqual(s.RawContent, []byte("s3cr3t")) {
		t.Errorf("Unexpected secret, got %s", string(s.RawContent))
	}

	transformation = core.TransformationConfig(core.WithTransformationType(adapters.JQTransformationType),
		core.WithInput("s1"),
		core.WithOutput("result"),
		core.WithSpec("q", ".then"),
		core.WithSpec("raw", false),
	)

	tt = adapters.NewJQTransformation(log.New(ioutil.Discard, "***", 0))
	s, err = tt.ProcessSecret(context.TODO(), &core.Defaults{}, secrets, transformation)
	if err != nil {
		t.Errorf("Error processing jq transformation: %s", err)
	}

	if s == nil {
		t.Errorf("Expected secret to be present, got nil")
	}

	if !reflect.DeepEqual(s.RawContent, []byte(`["the","rest"]`)) {
		t.Errorf("Unexpected secret, got %s", string(s.RawContent))
	}

	secrets = &core.Secrets{
		{
			Name: "s2",
			RawContent: []byte(`{
  "db1": {
	"endpoint": "foo",
    "apikey": "bar"
  },
  "db2": {
    "endpoint": "baz",
    "apikey": "qux"
  }
}`),
		},
	}

	transformation = &core.Transformation{
		Input:  []string{"s2"},
		Output: "result",
		Type:   "jq",
		Spec: core.TransformationSpec{
			"q":   ".db2",
			"raw": false,
		},
	}

	tt = adapters.NewJQTransformation(log.New(ioutil.Discard, "***", 0))
	s, err = tt.ProcessSecret(context.TODO(), &core.Defaults{}, secrets, transformation)
	if err != nil {
		t.Errorf("Error processing jq transformation: %s", err)
	}

	if s == nil {
		t.Errorf("Expected secret to be present, got nil")
	}

	// don't compare deepequal of byte sequences because of ordering.
	// parse and check values
	res := make(map[string]string)
	err = json.Unmarshal(s.RawContent, &res)
	if err != nil {
		t.Errorf("Unexpected err: %s", err)
	}

	if !(res["endpoint"] == "baz" && res["apikey"] == "qux") {
		t.Errorf("Unexpected secret, got %s", string(s.RawContent))
	}

}
