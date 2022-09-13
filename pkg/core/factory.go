package core

// Factory is able to create interfaces of type Repository, SinkWriterPort, TransformationPort
// and VaultAccessorPort, depending on the given type (e.g. a template-based sink as a SinkWriterPort).
// It also returns the types it is able to create.
type Factory interface {
	// NewRepository creates a new internal Repository for variables
	NewRepository() Repository

	// SinkTypes returns the list of sinks that this Factory produces
	SinkTypes() []string

	// NewSinkWriter creates a sink writer of given type
	NewSinkWriter(sinkType string) SinkWriterPort

	// TransformationTypes returns the list of transformation that this Factory produces
	TransformationTypes() []string

	// NewTransformation creates a Transformation of given type
	NewTransformation(transformationType string) TransformationPort

	// VaultAccessorTypes returns the list of vault accessors that this Factory produces
	VaultAccessorTypes() []string

	// NewVaultAccessor creates a vault accessor for a given
	NewVaultAccessor(vaultType string) VaultAccessorPort
}
