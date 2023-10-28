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

type ratingIngester interface {
	Ingest(ctx context.Context) (chan model.RatingEvent, error)
}

// Controller defines a rating service controller.
type Controller struct {
	repo     ratingRepository
	ingester ratingIngester
}

// New creates a rating service controller.
func New(repo ratingRepository, ingester ratingIngester) *Controller {
	return &Controller{repo, ingester}
}

// GetAggregatedRating returns the aggregated rating for a
// record or ErrNotFound if there are no ratings for it.

func (c *Controller) GetAgregatedRating(ctx context.Context, id model.RecordId, typ model.RecordType) (float64, error) {
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

func (c *Controller) PutAgregatedRating(ctx context.Context, id model.RecordId, typ model.RecordType, rating *model.Rating) error {
	return c.repo.Post(ctx, id, typ, rating)
}

// StartIngestion starts the ingestion of rating events.
func (s *Controller) StartIngestion(ctx context.Context) error {
	ch, err := s.ingester.Ingest(ctx)
	if err != nil {
		return err
	}
	for e := range ch {
		if err := s.PutAgregatedRating(ctx, e.RecordID, e.RecordType, &model.Rating{UserId: e.UserID, Value: e.Value}); err != nil {
			return err
		}
	}
	return nil
}
