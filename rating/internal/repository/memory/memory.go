package memory

import (
	"context"

	"movieapp.com/rating/internal/repository"
	"movieapp.com/rating/pkg/model"
)

type Repository struct {
	data map[model.RecordType]map[model.RecordId][]model.Rating
}

func New() *Repository {
	return &Repository{map[model.RecordType]map[model.RecordId][]model.Rating{}}
}

func (r *Repository) Get(ctx context.Context, id model.RecordId, typ model.RecordType) ([]model.Rating, error) {
	if _, ok := r.data[typ]; !ok {
		return nil, repository.ErrornotFound
	}

	if ratings, ok := r.data[typ][id]; !ok || len(ratings) == 0 {
		return nil, repository.ErrornotFound
	}

	return r.data[typ][id], nil
}

func (r *Repository) Post(ctx context.Context, id model.RecordId, typ model.RecordType, rating *model.Rating) error {
	if _, ok := r.data[typ]; !ok {
		r.data[typ] = map[model.RecordId][]model.Rating{}
	}
	r.data[typ][id] = append(r.data[typ][id], *rating)
	return nil
}
