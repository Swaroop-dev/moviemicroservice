package rating

import (
	"context"
	"errors"

	"movieapp.com/rating/internal/repository"
	"movieapp.com/rating/pkg/model"
)

var ErrNotFound = errors.New("ratings not found for a record")

type ratingRepository interface {
	Get(ctx context.Context, id model.RecordId, typ model.RecordType) ([]model.Rating, error)
	Post(ctx context.Context, id model.RecordId, typ model.RecordType, rating *model.Rating) error
}

type Controller struct {
	repo ratingRepository
}

func NewController(repo ratingRepository) *Controller {
	return &Controller{repo}
}

// GetAggregatedRating returns the aggregated rating for a
// record or ErrNotFound if there are no ratings for it.

func (c *Controller) Get(ctx context.Context, id model.RecordId, typ model.RecordType) (float64, error) {
	ratings, err := c.repo.Get(ctx, id, typ)

	if err != nil && err != repository.ErrornotFound {
		return 0, ErrNotFound
	}

	average := float64(0)
	cnt := 0
	for _, r := range ratings {
		average = average + float64(r.Value)
		cnt = cnt + 1
	}

	return average / float64(cnt), err
}
