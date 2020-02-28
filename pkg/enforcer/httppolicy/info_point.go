package httppolicy

import (
	"context"
	"encoding/json"
	"fmt"
	"net/http"

	"github.com/tierklinik-dobersberg/identity-server/pkg/enforcer"
)

// InfoPoint implements a HTTP based Policy Information Point (PIP).
type InfoPoint struct {
	url string
	cli *http.Client
}

// InfoPointOption is an option used when creating a HTTP PIP.
type InfoPointOption func(p *InfoPoint)

// NewInfoPoint returns a HTTP based Policy Information Point (PIP).
func NewInfoPoint(url string, opts ...InfoPointOption) *InfoPoint {
	pip := &InfoPoint{
		url: url,
		cli: http.DefaultClient,
	}

	return pip
}

// WithClient configures the HTTP client to use.
func WithClient(cli *http.Client) InfoPointOption {
	return func(p *InfoPoint) {
		p.cli = cli
	}
}

// GetResourceContext implements the InfoPoint interface.
func (pip *InfoPoint) GetResourceContext(ctx context.Context, resource string) (enforcer.Context, error) {
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

	var c enforcer.Context
	if err := json.NewDecoder(res.Body).Decode(&c); err != nil {
		return nil, err
	}

	return c, nil
}
