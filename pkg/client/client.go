package client

import "net/http"

// Option is used to configure the IdentityClient
type Option func(cli *IdentityClient) 

// IdentityClient talks to the identity-server using it's
// HTTP API.
type IdentityClient struct {
	cli *http.Client
	token TokenLoader
}

// NewIdentityClient returns a new IdentityClient that talks to the 
// identity-server running at url. By default, IdentityClient uses
// a &http.Client{} that can be modified by various Options. The 
// access token is read form the IAM_ACCESS_TOKEN environment variable
// by default. See WithTokenLoader() for more information.
func NewIdentityClient(url string, opts ...Option) *IdentityClient {
	c := &IdentityClient{
		cli: &http.Client{},
		token: NewEnvLoader("IAM_ACCESS_TOKEN"),
	}
	
	for _, opt := range opts {
		opt(c)
	}
	
	return c
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