package client

import (
	"context"
	"errors"

	"github.com/tierklinik-dobersberg/identity-server/pkg/iam"
)

// PolicyClient implements a HTTP client for the policy management
// endpoints.
type PolicyClient struct {
	*IdentityClient
}

// Create creates a new access policy under name.
func (pc *PolicyClient) Create(ctx context.Context, name string, policy iam.Policy) (iam.PolicyURN, error) {
	body := struct {
		Name   string     `json:"name"`
		Policy iam.Policy `json:"policy"`
	}{
		Name:   name,
		Policy: policy,
	}

	req, err := pc.newRequest(ctx, "POST", "/v1/policies/", body)
	if err != nil {
		return "", err
	}

	res, err := pc.cli.Do(req)
	if err != nil {
		return "", err
	}

	var response struct {
		URN iam.PolicyURN `json:"urn"`
	}

	return response.URN, pc.parseResponse(res, &response)
}

// Delete deletes a policy.
func (pc *PolicyClient) Delete(ctx context.Context, urn iam.PolicyURN) error {
	name := urn.PolicyName()
	if name == "" {
		return errors.New("Invalid policy")
	}

	req, err := pc.newRequest(ctx, "DELETE", "/v1/policies/"+name, nil)
	if err != nil {
		return err
	}

	res, err := pc.cli.Do(req)
	if err != nil {
		return err
	}

	return pc.parseResponse(res, nil)
}

// Load loads the policy with the given URN.
func (pc *PolicyClient) Load(ctx context.Context, urn iam.PolicyURN) (iam.Policy, error) {
	name := urn.PolicyName()
	if name == "" {
		return iam.Policy{}, errors.New("Invalid policy")
	}

	req, err := pc.newRequest(ctx, "GET", "/v1/policies/"+name, nil)
	if err != nil {
		return iam.Policy{}, err
	}

	res, err := pc.cli.Do(req)
	if err != nil {
		return iam.Policy{}, err
	}

	var response iam.Policy

	return response, pc.parseResponse(res, &response)
}

// Update updates an existing policy.
func (pc *PolicyClient) Update(ctx context.Context, urn iam.PolicyURN, p iam.Policy) error {
	name := urn.PolicyName()
	if name == "" {
		return errors.New("Invalid policy")
	}

	req, err := pc.newRequest(ctx, "GET", "/v1/policies/"+name, p)
	if err != nil {
		return err
	}

	res, err := pc.cli.Do(req)
	if err != nil {
		return err
	}

	return pc.parseResponse(res, nil)
}

// List returns a list of all available policies.
func (pc *PolicyClient) List(ctx context.Context) ([]iam.Policy, error) {
	req, err := pc.newRequest(ctx, "GET", "/v1/policies/", nil)
	if err != nil {
		return nil, err
	}

	res, err := pc.cli.Do(req)
	if err != nil {
		return nil, err
	}

	var response struct {
		Policies []iam.Policy `json:"policies"`
	}

	return response.Policies, pc.parseResponse(res, &response)
}
