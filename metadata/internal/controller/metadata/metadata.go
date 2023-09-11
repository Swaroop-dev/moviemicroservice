package metadata

import (
	"context"
	"errors"

	"movieapp.com/metadata/internal/repository"
	"movieapp.com/metadata/pkg/model"
)

var ErrornotFound = errors.New("not found")

type metadataRepository interface {
	Get(ctx context.Context, id string) (*model.Metadata, error)
}

type Controller struct {
	repo metadataRepository
}

func New(repo metadataRepository) *Controller {
	return &Controller{repo}
}

func (c *Controller) Get(ctx context.Context, id string) (*model.Metadata, error) {
	res, err := c.repo.Get(ctx, id)

	if err != nil && errors.Is(err, repository.ErrornotFound) {
		return nil, ErrornotFound
	}

	return res, nil

}
