package test

import (
	"github.com/golang/mock/gomock"
	"go-secretshelper/pkg/core"
	"go-secretshelper/pkg/core/mocks"
	"testing"
)

// MockFactory produces mocks only
type MockFactory struct {
	mockCtrl *gomock.Controller
	t        *testing.T

	repo            *mocks.MockRepository
	vaults          map[string]*mocks.MockVaultAccessorPort
	sinks           map[string]*mocks.MockSinkWriterPort
	transformations map[string]*mocks.MockTransformationPort
}

// NewMockFactory creates a new mock factors
func NewMockFactory(mockCtrl *gomock.Controller, t *testing.T) *MockFactory {
	mf := &MockFactory{
		mockCtrl:        mockCtrl,
		t:               t,
		vaults:          make(map[string]*mocks.MockVaultAccessorPort),
		sinks:           make(map[string]*mocks.MockSinkWriterPort),
		transformations: make(map[string]*mocks.MockTransformationPort),
		repo:            mocks.NewMockRepository(mockCtrl),
	}

	// auto set up mock port
	mf.newVaultAccessorInternal("mock")
	mf.newSinkWriterInternal("mock")
	mf.newTransformationInternal("mock")

	return mf
}

// SinkTypes returns valid sink types
func (df *MockFactory) SinkTypes() []string {
	return []string{
		"mock",
	}
}

// TransformationTypes returns valid transformation types
func (df *MockFactory) TransformationTypes() []string {
	return []string{
		"mock",
	}
}

// VaultAccessorTypes returns valid vault types
func (df *MockFactory) VaultAccessorTypes() []string {
	return []string{
		"mock",
	}
}

// NewRepository creates a new repository
func (df *MockFactory) NewRepository() core.Repository {
	return df.repo
}

// GetMockRepository returns created repo
func (df *MockFactory) GetMockRepository() *mocks.MockRepository {
	return df.repo
}

func (df *MockFactory) newSinkWriterInternal(sinkType string) core.SinkWriterPort {
	s := mocks.NewMockSinkWriterPort(df.mockCtrl)
	df.sinks[sinkType] = s
	return s
}

// GetMockSinkWriter returns the mock sink writer for a given type
func (df *MockFactory) GetMockSinkWriter(sinkType string) *mocks.MockSinkWriterPort {
	return df.sinks[sinkType]
}

// NewSinkWriter creates a new sink writer for a supported type
func (df *MockFactory) NewSinkWriter(sinkType string) core.SinkWriterPort {
	return df.sinks[sinkType]
}

// NewTransformation creates a new transformation for a supported type
func (df *MockFactory) NewTransformation(transformationType string) core.TransformationPort {
	return df.transformations[transformationType]
}

func (df *MockFactory) newTransformationInternal(transformationType string) core.TransformationPort {
	s := mocks.NewMockTransformationPort(df.mockCtrl)
	df.transformations[transformationType] = s
	return s
}

// GetMockTransformation returns the mock transformation for a given type
func (df *MockFactory) GetMockTransformation(transformationType string) *mocks.MockTransformationPort {
	return df.transformations[transformationType]
}

// NewVaultAccessor creates a new vault accessor for a supported type
func (df *MockFactory) NewVaultAccessor(vaultType string) core.VaultAccessorPort {
	return df.vaults[vaultType]
}

func (df *MockFactory) newVaultAccessorInternal(vaultType string) core.VaultAccessorPort {
	va := mocks.NewMockVaultAccessorPort(df.mockCtrl)
	df.vaults[vaultType] = va
	return va
}

// GetMockVaultAccessor returns a vault accessor for a type
func (df *MockFactory) GetMockVaultAccessor(t string) *mocks.MockVaultAccessorPort {
	return df.vaults[t]
}
