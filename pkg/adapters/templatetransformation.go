package adapters

import (
	"context"
	"fmt"
	"go-secretshelper/pkg/core"
	"log"
	"strings"
	"text/template"
)

// TemplateTransformationType is the type for a template transformation
const TemplateTransformationType = "template"

// TemplateTransformation is an adapter that transforms a secret using a template
// see TransformationPort interface for more details
type TemplateTransformation struct {
	log *log.Logger
}

// TemplateTransformationSpec contains the specification of a template
type TemplateTransformationSpec struct {
	// Template text to use for transformation
	Template *template.Template

	// Content Type of rendered output (default: text/plain)
	ContentType string
}

// NewTemplateTransformationSpec creates a new TemplateTransformationSpec from a generic map
func NewTemplateTransformationSpec(in map[interface{}]interface{}) (TemplateTransformationSpec, error) {

	templateSource, ok := in["template"]
	if !ok {
		return TemplateTransformationSpec{}, fmt.Errorf("template element is required")
	}
	templateSourceStr, ok := templateSource.(string)
	if !ok {
		return TemplateTransformationSpec{}, fmt.Errorf("template element must be a string")
	}

	tmpl, err := template.New(templateSourceStr).Parse(templateSourceStr)
	if err != nil {
		return TemplateTransformationSpec{}, err
	}

	var contentType = "text/plain"
	cnt, ok := in["contentType"]
	if ok {
		cntStr, ok := cnt.(string)
		if ok {
			contentType = cntStr
		}
	}

	return TemplateTransformationSpec{
		Template:    tmpl,
		ContentType: contentType,
	}, nil
}

// NewTemplateTransformation returns a new instance of TemplateTransformation
func NewTemplateTransformation(log *log.Logger) *TemplateTransformation {
	return &TemplateTransformation{log: log}
}

// ProcessSecret returns a new secret as the result of a template rendering process
func (t *TemplateTransformation) ProcessSecret(ctx context.Context,
	defaults *core.Defaults, in *core.Secrets, transformation *core.Transformation) (*core.Secret, error) {

	spec, err := NewTemplateTransformationSpec(transformation.Spec)
	if err != nil {
		return nil, err
	}

	b := new(strings.Builder)
	data := make(map[string]interface{})
	for _, inVar := range *in {
		data[inVar.Name] = string(inVar.RawContent)
	}

	err = spec.Template.Execute(b, data)
	if err != nil {
		return nil, err
	}

	res := &core.Secret{
		Name:           transformation.Output,
		Type:           "transformed-by:template",
		RawContent:     []byte(b.String()),
		RawContentType: spec.ContentType,
	}

	return res, nil
}
