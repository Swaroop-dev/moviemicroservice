package memory

import (
	"context"
	"sync"

	"movieapp.com/metadata/internal/repository"
	"movieapp.com/metadata/pkg/model"
)

type Repository struct {
	sync.RWMutex
	data map[string]*model.Metadata
}

// New func creates new memory repository
func New() *Repository {
	return &Repository{data: map[string]*model.Metadata{}}
}

// Get returns movie metadata from memory
func (r *Repository) Get(_ context.Context, ID string) (*model.Metadata, error) {
	r.Lock()
	defer r.Unlock()

	m, ok := r.data[ID]

	if !ok {
		return nil, repository.ErrornotFound
	}

	return m, nil
}

func (r *Repository) Post(_ context.Context, ID string, metadata *model.Metadata) error {
	r.Lock()
	defer r.Unlock()

	r.data[ID] = metadata
	return nil
}
