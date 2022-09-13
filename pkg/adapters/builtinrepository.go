package adapters

import (
	"go-secretshelper/pkg/core"
	"sync"
)

// BuiltinRepository stores all item in a map
type BuiltinRepository struct {
	items map[string]interface{}
	m     *sync.Mutex
}

// NewBuiltinRepository creates a new BuiltinRepository
func NewBuiltinRepository() *BuiltinRepository {
	return &BuiltinRepository{
		items: make(map[string]interface{}),
		m:     &sync.Mutex{},
	}
}

// Put places varName with content in repository
func (r *BuiltinRepository) Put(varName string, content interface{}) {
	r.m.Lock()
	defer r.m.Unlock()

	r.items[varName] = content
}

// Get returns varName or an error
func (r *BuiltinRepository) Get(varName string) (interface{}, error) {
	r.m.Lock()
	defer r.m.Unlock()

	res, ex := r.items[varName]
	if !ex {
		return nil, core.RepositoryError{Reason: "No such variable", Info: varName}
	}
	return res, nil
}
