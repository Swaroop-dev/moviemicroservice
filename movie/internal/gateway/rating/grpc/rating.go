package grpc

import (
	"context"

	"movieapp.com/gen"
	"movieapp.com/internal/grpcutil"
	"movieapp.com/pkg/discovery"
	"movieapp.com/rating/pkg/model"
)

type Gateway struct {
	registry discovery.Registry
}

func NewGateway(registry discovery.Registry) *Gateway {
	return &Gateway{registry: registry}
}

func (g *Gateway) GetAgregatedRating(ctx context.Context, id model.RecordId, typ model.RecordType) (float64, error) {
	conn, err := grpcutil.ServiceConnection(ctx, "rating", g.registry)
	if err != nil {
		return 0, err
	}

	defer conn.Close()

	client := gen.NewRatingServiceClient(conn)

	resp, err := client.GetAggregatedRating(ctx, &gen.GetAggregatedRatingRequest{RecordId: string(id), RecordType: string(typ)})

	if err != nil {
		return 0, err
	}

	return resp.RatingValue, nil

}
