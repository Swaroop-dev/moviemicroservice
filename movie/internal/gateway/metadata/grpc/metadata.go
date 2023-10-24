package grpc

import (
	"context"

	"movieapp.com/gen"
	"movieapp.com/internal/grpcutil"
	"movieapp.com/metadata/pkg/model"
	"movieapp.com/pkg/discovery"
)

type Gateway struct {
	registry discovery.Registry
}

func NewGateway(registry discovery.Registry) *Gateway {
	return &Gateway{registry: registry}
}

func (g *Gateway) Get(ctx context.Context, id string) (*model.Metadata, error) {
	conn, err := grpcutil.ServiceConnection(ctx, "metadata", g.registry)
	if err != nil {
		return nil, err
	}

	defer conn.Close()

	client := gen.NewMetadataServiceClient(conn)

	resp, err := client.GetMetadata(ctx, &gen.GetMetadataRequest{Movieid: id})
	if err != nil {
		return nil, err
	}

	return model.MetadataFromProto(resp.Metadata), nil

}
