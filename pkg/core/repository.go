// Package core contains the components for a repository
//go:generate mockgen -package mocks -destination=mocks/mock_repository.go go-secretshelper/pkg/core Repository
package core

import "fmt"

// RepositoryError ...
type RepositoryError struct {
	Reason string
	Info   string
}

func (e RepositoryError) Error() string { return fmt.Sprintf("%s: %s", e.Reason, e.Info) }

var (
	// RepositoryErrorNoSuchVariable ...
	RepositoryErrorNoSuchVariable = RepositoryError{Reason: "No such variable"}
)

// Repository stores secrets under variable names
type Repository interface {

	// Put (re)places a secret given by the variable name and its content
	Put(varName string, content interface{})

	// Get retrieves a secret by its variable name, or an error
	Get(varName string) (interface{}, error)
}
