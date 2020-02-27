package policy

import (
	"context"
	"net/http"

	"github.com/go-kit/kit/endpoint"
	"github.com/tierklinik-dobersberg/identity-server/iam"
)

type createPolicyRequest struct {
	Name   string     `json:"name"`
	Policy iam.Policy `json:"policy"`
}

type createPolicyResponse struct {
	URN iam.PolicyURN `json:"urn"`
}

func makeCreatePolicyEndpoint(s Service) endpoint.Endpoint {
	return func(ctx context.Context, request interface{}) (interface{}, error) {
		req := request.(createPolicyRequest)
		urn, err := s.Create(ctx, req.Name, req.Policy)
		if err != nil {
			return nil, err
		}

		return createPolicyResponse{urn}, nil
	}
}

type deletePolicyRequest struct {
	URN iam.PolicyURN
}
type deletePolicyResponse struct{}

func (deletePolicyResponse) StatusCode() int { return http.StatusNoContent }

func makeDeletePolicyEndpoint(s Service) endpoint.Endpoint {
	return func(ctx context.Context, request interface{}) (interface{}, error) {
		req := request.(deletePolicyRequest)
		err := s.Delete(ctx, req.URN)
		if err != nil {
			return nil, err
		}

		return deletePolicyResponse{}, nil
	}
}

type updatePolicyRequest struct {
	URN    iam.PolicyURN
	Policy iam.Policy
}

type updatePolicyResponse struct {
}

func (updatePolicyResponse) StatusCode() int { return http.StatusNoContent }

func makeUpdatePolicyEndpoint(s Service) endpoint.Endpoint {
	return func(ctx context.Context, request interface{}) (interface{}, error) {
		req := request.(updatePolicyRequest)
		err := s.Update(ctx, req.URN, req.Policy)
		if err != nil {
			return nil, err
		}

		return updatePolicyResponse{}, nil
	}
}

type listPoliciesRequest struct{}
type listPoliciesResponse struct {
	Policies []iam.Policy `json:"policies"`
}

func makeListPoliciesEndpoint(s Service) endpoint.Endpoint {
	return func(ctx context.Context, request interface{}) (interface{}, error) {
		_ = request.(listPoliciesRequest)
		policies, err := s.List(ctx)
		if err != nil {
			return nil, err
		}

		return listPoliciesResponse{policies}, nil
	}
}
