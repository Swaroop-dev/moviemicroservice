package http

import (
	"context"
	"encoding/json"
	"fmt"
	"net/http"

	"movieapp.com/metadata/pkg/model"
	"movieapp.com/movie/internal/gateway"
)

type Gateway struct {
	addr string
}

func New(addr string) *Gateway {
	return &Gateway{addr: addr}
}

// Get movie metadata
func (g *Gateway) Get(ctx context.Context, id string) (*model.Metadata, error) {
	req, err := http.NewRequest(http.MethodGet, g.addr+"/metadata", nil)

	if err != nil {
		return nil, err
	}

	req = req.WithContext(ctx)
	values := req.URL.Query()

	values.Add("id", id)
	req.URL.RawQuery = values.Encode()

	resp, err := http.DefaultClient.Do(req)
	if err != nil {
		return nil, err
	}

	defer resp.Body.Close()

	if resp.StatusCode == http.StatusFound {
		return nil, gateway.ErrNotFound
	} else if resp.StatusCode/100 != 2 {
		return nil, fmt.Errorf("non 2xx response from the request is %v", resp)
	}

	var metadata *model.Metadata

	if err := json.NewDecoder(resp.Body).Decode(&metadata); err != nil {
		return nil, err
	}

	return metadata, nil
}
