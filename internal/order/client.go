package order

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"net/http"
	"time"
)

// ErrProductNotFound is returned when product-service has no such product.
var ErrProductNotFound = errors.New("product not found")

// Product mirrors the fields order-service needs from product-service.
type Product struct {
	ID    string `json:"id"`
	Name  string `json:"name"`
	Price int64  `json:"price"`
	Stock int    `json:"stock"`
}

// ProductClient is the service-to-service boundary (easy to fake in tests).
type ProductClient interface {
	GetProduct(ctx context.Context, id string) (Product, error)
}

// HTTPProductClient calls product-service over HTTP.
type HTTPProductClient struct {
	BaseURL string
	HTTP    *http.Client
}

func NewHTTPProductClient(baseURL string) *HTTPProductClient {
	return &HTTPProductClient{BaseURL: baseURL, HTTP: &http.Client{Timeout: 3 * time.Second}}
}

func (c *HTTPProductClient) GetProduct(ctx context.Context, id string) (Product, error) {
	req, err := http.NewRequestWithContext(ctx, http.MethodGet, c.BaseURL+"/products/"+id, nil)
	if err != nil {
		return Product{}, err
	}
	resp, err := c.HTTP.Do(req)
	if err != nil {
		return Product{}, err
	}
	defer resp.Body.Close()

	switch resp.StatusCode {
	case http.StatusOK:
		var p Product
		if err := json.NewDecoder(resp.Body).Decode(&p); err != nil {
			return Product{}, err
		}
		return p, nil
	case http.StatusNotFound:
		return Product{}, ErrProductNotFound
	default:
		return Product{}, fmt.Errorf("product-service returned %d", resp.StatusCode)
	}
}
