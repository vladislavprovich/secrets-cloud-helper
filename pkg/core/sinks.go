package core

import (
	"context"
	"fmt"
)

// Sinks is an array of Sink structs
type Sinks []*Sink

// Sink defines where a secret, indicated by Var, should be written to.
type Sink struct {
	// Type of sink. determines structure of spec.
	Type string `yaml:"type" validate:"required"`

	// Var defines which variable is written into the sink
	Var string `yaml:"var" validate:"required"`

	// Spec optionally defines properties of the sink
	Spec SinkSpec `yaml:"spec" validate:""`
}

// SinkSpec contains details about where and how it should be written
type SinkSpec map[interface{}]interface{}

// String creates a string representation of a sink
func (s Sink) String() string {
	return fmt.Sprintf("Sink:[Var=%s, Type=%s]", s.Var, s.Type)
}

// SinkWriterPort is able to write a secret into a defined sink
type SinkWriterPort interface {

	// Write takes the raw content of given secret and writes it to the sink using the defaults
	Write(context.Context, *Defaults, *Secret, *Sink) error
}
