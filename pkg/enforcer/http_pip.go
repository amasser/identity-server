package enforcer

import (
	"context"
	"encoding/json"
	"fmt"
	"net/http"
)

// HTTPInfoPoint implements a HTTP based Policy Information Point (PIP).
type HTTPInfoPoint struct {
	url string
	cli *http.Client
}

// HTTPInfoPointOption is an option used when creating a HTTP PIP.
type HTTPInfoPointOption func(p *HTTPInfoPoint)

// NewHTTPInfoPoint returns a HTTP based Policy Information Point (PIP).
func NewHTTPInfoPoint(url string, opts ...HTTPInfoPointOption) *HTTPInfoPoint {
	pip := &HTTPInfoPoint{
		url: url,
		cli: http.DefaultClient,
	}

	return pip
}

// WithClient configures the HTTP client to use.
func WithClient(cli *http.Client) HTTPInfoPointOption {
	return func(p *HTTPInfoPoint) {
		p.cli = cli
	}
}

// GetResourceContext implements the InfoPoint interface.
func (pip *HTTPInfoPoint) GetResourceContext(ctx context.Context, resource string) (Context, error) {
	req, err := http.NewRequest("GET", pip.url, nil)
	if err != nil {
		return nil, err
	}

	res, err := pip.cli.Do(req)
	if err != nil {
		return nil, err
	}

	if res.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("Unexpected status code %q from %q", res.Status, pip.url)
	}

	var c Context
	if err := json.NewDecoder(res.Body).Decode(&c); err != nil {
		return nil, err
	}

	return c, nil
}
