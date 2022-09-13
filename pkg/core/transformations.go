// Package core contains the core components for a transformation
//go:generate mockgen -package mocks -destination=mocks/mock_transformationport.go go-secretshelper/pkg/core TransformationPort
package core

import "context"

// Transformations is an array of Transformation structs
type Transformations []*Transformation

// Transformation describe a single transformation
type Transformation struct {
	// Input is the list of input variables for this transformation. These must have
	// been defined as secrets or must have been processed before by other transformations
	Input []string `yaml:"in" validate:"required,dive,required"`

	// Output is the name of output variable. The result of the transformation will go here.
	Output string `yaml:"out" validate:"required"`

	// Type is the type of transformation
	Type string `yaml:"type" validate:"required"`

	// Spec is the generic specification for a transformation of a given type
	Spec TransformationSpec `yaml:"spec" validate:""`
}

// TransformationConfigOpts is the fluent-style configuration option func
type TransformationConfigOpts func(*Transformation)

// TransformationConfig creates a Transformation struct from the given options
func TransformationConfig(opts ...TransformationConfigOpts) *Transformation {
	res := &Transformation{
		Input: make([]string, 0),
		Spec:  TransformationSpec{},
	}
	for _, opt := range opts {
		opt(res)
	}
	return res
}

// WithTransformationType sets the type of the transformation
func WithTransformationType(tp string) TransformationConfigOpts {
	return func(t *Transformation) {
		t.Type = tp
	}
}

// WithOutput sets the output variable name
func WithOutput(outp string) TransformationConfigOpts {
	return func(t *Transformation) {
		t.Output = outp
	}
}

// WithSpec sets a key/value pair of the spec of the transformation
func WithSpec(key string, value interface{}) TransformationConfigOpts {
	return func(t *Transformation) {
		t.Spec[key] = value
	}
}

// WithInput sets the input variables by adding the given one.
func WithInput(inp string) TransformationConfigOpts {
	return func(t *Transformation) {
		t.Input = append(t.Input, inp)
	}
}

// TransformationSpec is the generic specification for a transformation of a given type (simple k/v pairs)
type TransformationSpec map[interface{}]interface{}

// TransformationPort is the interface for a single transformation
type TransformationPort interface {
	// ProcessSecret applies the Transformation, using given Secrets and returns an updated Secret
	ProcessSecret(context.Context, *Defaults, *Secrets, *Transformation) (*Secret, error)
}
