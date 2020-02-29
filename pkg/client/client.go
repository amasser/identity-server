package client

import (
	"bytes"
	"context"
	"encoding/json"
	"errors"
	"io/ioutil"
	"net/http"
)

// Option is used to configure the IdentityClient
type Option func(cli *IdentityClient)

// IdentityClient talks to the identity-server using it's
// HTTP API.
type IdentityClient struct {
	cli   *http.Client
	url   string
	token TokenLoader
}

// NewIdentityClient returns a new IdentityClient that talks to the
// identity-server running at url. By default, IdentityClient uses
// a &http.Client{} that can be modified by various Options. The
// access token is read from the IAM_ACCESS_TOKEN environment variable
// by default. See WithTokenLoader() for more information.
func NewIdentityClient(url string, opts ...Option) *IdentityClient {
	c := &IdentityClient{
		cli:   &http.Client{},
		url:   url,
		token: NewEnvLoader("IAM_ACCESS_TOKEN"),
	}

	for _, opt := range opts {
		opt(c)
	}

	return c
}

func (cli *IdentityClient) newRequest(ctx context.Context, method string, endpoint string, body interface{}) (*http.Request, error) {
	req, err := http.NewRequest(method, cli.url+method, nil)
	if err != nil {
		return nil, err
	}

	req = req.Clone(ctx)

	token, err := cli.token.Load()
	if err != nil {
		return nil, err
	}

	req.Header.Add("Authorization", "Bearer "+token)

	if body != nil {
		blob, err := json.Marshal(body)
		if err != nil {
			return nil, err
		}

		req.Body = ioutil.NopCloser(bytes.NewBuffer(blob))
	}

	return req, nil
}

func (cli *IdentityClient) parseResponse(res *http.Response, target interface{}) error {
	if res.StatusCode >= 300 {
		return errors.New(res.Status)
	}

	if target != nil {
		if err := json.NewDecoder(res.Body).Decode(target); err != nil {
			return err
		}
		defer res.Body.Close()
	}

	return nil
}

// Users returns a UsersClient using this IdentityClient.
func (cli *IdentityClient) Users() *UserClient {
	return &UserClient{cli}
}

// WithClient sets the http.Client that should be used.
func WithClient(cli *http.Client) Option {
	return func(c *IdentityClient) {
		c.cli = cli
	}
}

// WithTokenLoader configures the access token loader to
// use.
func WithTokenLoader(loader TokenLoader) Option {
	return func(c *IdentityClient) {
		c.token = loader
	}
}
