package enforcer

import (
	"bytes"
	"context"
	"encoding/json"
	"errors"
	"net/http"

	"github.com/ory/ladon"
)

// RequestEncoder encodes the given permission request so it can be transmitted to a remote
// policy enforcement point.
type RequestEncoder func(ctx context.Context, subject, action, resource string, context Context) ([]byte, error)

// HTTPEnforcerOption applies custom options to the HTTP enforcer
type HTTPEnforcerOption func(h *HTTPEnforcer)

// HTTPEnforcer implements the Enforcer interface fowarding permission requests
// to a remote HTTP endpoint.
type HTTPEnforcer struct {
	url         string
	cli         *http.Client
	encoder     RequestEncoder
	contentType string
}

// NewHTTPEnforcer returns a new policy enforcer that queries a remote endpoint.
func NewHTTPEnforcer(url string, opts ...HTTPEnforcerOption) *HTTPEnforcer {
	e := &HTTPEnforcer{
		url:         url,
		cli:         http.DefaultClient,
		encoder:     DefaultRequestEncoder,
		contentType: "application/json",
	}

	for _, opt := range opts {
		opt(e)
	}

	return e
}

// Enforce implements the Enforcer interface and sends a HTTP POST request to the URL configured in
// NewHTTPEnforcer.
func (e *HTTPEnforcer) Enforce(ctx context.Context, subject, action, resource string, context Context) error {
	payload, err := e.encoder(ctx, subject, action, resource, context)
	if err != nil {
		return err
	}

	req, err := http.NewRequest("POST", e.url, bytes.NewBuffer(payload))
	if err != nil {
		return err
	}

	res, err := e.cli.Do(req)
	if err != nil {
		return err
	}

	if res.StatusCode >= http.StatusOK && res.StatusCode < http.StatusMultipleChoices {
		return nil
	}

	return errors.New(res.Status)
}

// DefaultRequestEncoder is used as the default RequestEncoder in NewHTTPEnforcer and directly encodes
// a ladon.Request.
func DefaultRequestEncoder(_ context.Context, subject, action, resource string, context Context) ([]byte, error) {
	return json.Marshal(ladon.Request{
		Action:   action,
		Subject:  subject,
		Resource: resource,
		Context:  ladon.Context(context),
	})
}

// WithHTTPClient configures the http.Client to use
// for the HTTPEnforcer.
func WithHTTPClient(cli *http.Client) HTTPEnforcerOption {
	return func(e *HTTPEnforcer) {
		e.cli = cli
	}
}

// WithRequestEncoder configures the RequestEncoder a HTTPEnforcer should use.
func WithRequestEncoder(encoder RequestEncoder) HTTPEnforcerOption {
	return func(e *HTTPEnforcer) {
		e.encoder = encoder
	}
}

// WithContentType configures the content type to use for enforcement requests.
func WithContentType(contentType string) HTTPEnforcerOption {
	return func(e *HTTPEnforcer) {
		e.contentType = contentType
	}
}
