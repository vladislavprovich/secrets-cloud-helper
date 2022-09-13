package test

import (
	"context"
	"go-secretshelper/pkg/adapters"
	"go-secretshelper/pkg/core"
	"io/ioutil"
	"log"
	"reflect"
	"testing"
)

func TestTemplateTransformationSpec(t *testing.T) {
	ts := "Name: {{ .s1 }}\nValue: {{ .s2 }}\n"
	in := core.TransformationSpec{
		"template": ts,
	}

	spec, err := adapters.NewTemplateTransformationSpec(in)
	if err != nil {
		t.Errorf("Error creating spec: %s", err)
	}
	if spec.Template == nil {
		t.Errorf("Expected template to be present, got nil")
	}

	spec, err = adapters.NewTemplateTransformationSpec(core.TransformationSpec{})
	if err == nil {
		t.Errorf("Expected error creating spec")
	}

}

func TestTemplateTransformation(t *testing.T) {
	secrets := &core.Secrets{
		{
			Name:       "s1",
			RawContent: []byte("123"),
		},
		{
			Name:       "s2",
			RawContent: []byte("456"),
		},
	}
	ts := "Name: {{ .s1 }}\nValue: {{ .s2 }}\n"
	tr := "Name: 123\nValue: 456\n"

	transformation := &core.Transformation{
		Input:  []string{"s1", "s2"},
		Output: "result",
		Type:   "template",
		Spec: core.TransformationSpec{
			"template": ts,
		},
	}

	tt := adapters.NewTemplateTransformation(log.New(ioutil.Discard, "", 0))
	s, err := tt.ProcessSecret(context.TODO(), &core.Defaults{}, secrets, transformation)
	if err != nil {
		t.Errorf("Error processing template transformation: %s", err)
	}

	if s == nil {
		t.Errorf("Expected secret to be present, got nil")
	}

	if !reflect.DeepEqual(s.RawContent, []byte(tr)) {
		t.Errorf("Expected secret to be %s, got %s", tr, string(s.RawContent))
	}
}
