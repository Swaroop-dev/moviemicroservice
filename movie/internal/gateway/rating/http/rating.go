package http

import (
	"context"
	"encoding/json"
	"fmt"
	"net/http"

	"movieapp.com/movie/internal/gateway"
	"movieapp.com/rating/pkg/model"
)

type Gateway struct {
	addr string
}

func New(addr string) *Gateway {
	return &Gateway{addr: addr}
}

func (g *Gateway) GetAgregatedRating(ctx context.Context, id model.RecordId, typ model.RecordType) (float64, error) {
	req, err := http.NewRequest("GET", g.addr+"/rating", nil)
	if err != nil {
		return 0, err
	}

	req = req.WithContext(ctx)
	values := req.URL.Query()

	values.Add("id", string(id))
	values.Add("typ", string(typ))
	req.URL.RawQuery = values.Encode()

	resp, err := http.DefaultClient.Do(req)
	if err != nil {
		return 0, err
	}
	defer resp.Body.Close()

	if resp.StatusCode == http.StatusFound {
		return 0, gateway.ErrNotFound
	} else if resp.StatusCode/100 != 2 {
		return 0, fmt.Errorf("non 2xx response from the request is %v", resp)
	}

	var rating float64

	if err := json.NewDecoder(resp.Body).Decode(&rating); err != nil {
		return 0, err
	}
	return rating, nil
}

func (g *Gateway) PutRating(ctx context.Context, id model.RecordId, typ model.RecordType, rating *model.Rating) error {
	req, err := http.NewRequest(http.MethodPut, g.addr+"/rating", nil)
	if err != nil {
		return err
	}
	req = req.WithContext(ctx)
	values := req.URL.Query()

	values.Add("id", string(id))
	values.Add("value", fmt.Sprint("%v", rating.Value))
	values.Add("userId", fmt.Sprint("%v", rating.UserId))
	values.Add("type", string(typ))

	req.URL.RawQuery = values.Encode()
	res, err := http.DefaultClient.Do(req)
	if err != nil {
		return err
	}
	defer res.Body.Close()

	if res.StatusCode == http.StatusFound {
		return gateway.ErrNotFound
	} else if res.StatusCode/100 != 2 {
		return fmt.Errorf("some error %v", res)
	}

	return nil
}
