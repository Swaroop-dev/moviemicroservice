package movie

import (
	"context"
	"errors"

	metadata "movieapp.com/metadata/pkg/model"
	"movieapp.com/movie/internal/gateway"
	"movieapp.com/movie/pkg/model"
	rating "movieapp.com/rating/pkg/model"
)

var ErrNotFound = errors.New("movie metadata not found")

type ratingGateway interface {
	GetAgregatedRating(ctx context.Context, id rating.RecordId, typ rating.RecordType) (float64, error)
	PutRating(ctx context.Context, id rating.RecordId, typ rating.RecordType, rating *rating.Rating) error
}

type metadataGateway interface {
	Get(ctx context.Context, id string) (*metadata.Metadata, error)
}

type Controller struct {
	ratingGateway   ratingGateway
	metadataGateway metadataGateway
}

// New creates a new movie service controller
func New(r ratingGateway, m metadataGateway) *Controller {
	return &Controller{ratingGateway: r, metadataGateway: m}
}

// Get returns the movie details including the aggregated
// rating and movie metadata.

func (c *Controller) Get(ctx context.Context, id string) (*model.MovieDetails, error) {
	metadata, err := c.metadataGateway.Get(ctx, id)
	if err != nil && errors.Is(err, gateway.ErrNotFound) {
		return nil, ErrNotFound
	} else if err != nil {
		return nil, err
	}

	details := &model.MovieDetails{MetaData: *metadata}
	rating, err := c.ratingGateway.GetAgregatedRating(ctx, rating.RecordId(id), rating.RecordTypeMovie)
	if err != nil && errors.Is(err, gateway.ErrNotFound) {
		//dont do anything since there might be  a case where movie is not rated yet
	} else if err != nil {
		return nil, err
	}
	details.Rating = rating

	return details, nil
}
