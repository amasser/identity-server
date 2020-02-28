package policy

import (
	"context"
	"encoding/json"
	"fmt"
	"net/http"

	"github.com/go-kit/kit/endpoint"
	"github.com/go-kit/kit/log"
	"github.com/go-kit/kit/transport"
	kithttp "github.com/go-kit/kit/transport/http"
	"github.com/gorilla/mux"
	"github.com/tierklinik-dobersberg/identity-server/iam"
	"github.com/tierklinik-dobersberg/identity-server/pkg/authn"
	"github.com/tierklinik-dobersberg/identity-server/pkg/common"
)

// MakeHandler returns a http.Handler for the policy management service.
func MakeHandler(s Service, extractor authn.SubjectExtractorFunc, logger log.Logger) http.Handler {
	opts := []kithttp.ServerOption{
		kithttp.ServerErrorHandler(transport.NewLogErrorHandler(logger)),
		kithttp.ServerErrorEncoder(kithttp.DefaultErrorEncoder),
		kithttp.ServerBefore(kithttp.PopulateRequestContext),
	}

	auth := authn.NewAuthenticator(extractor)
	makeEndpoint := func(factory func(s Service) endpoint.Endpoint) endpoint.Endpoint {
		return endpoint.Chain(auth)(factory(s))
	}

	createPolicyHandler := kithttp.NewServer(
		makeEndpoint(makeCreatePolicyEndpoint),
		decodeCreatePolicyRequest,
		kithttp.EncodeJSONResponse,
		opts...,
	)

	deletePolicyHandler := kithttp.NewServer(
		makeEndpoint(makeDeletePolicyEndpoint),
		decodeDeletePolicyRequest,
		kithttp.EncodeJSONResponse,
		opts...,
	)

	loadPolicyHandler := kithttp.NewServer(
		makeEndpoint(makeLoadPolicyEndpoint),
		decodeLoadPolicyRequest,
		kithttp.EncodeJSONResponse,
		opts...,
	)

	updatePolicyHandler := kithttp.NewServer(
		makeEndpoint(makeUpdatePolicyEndpoint),
		decodeUpdatePolicyRequest,
		kithttp.EncodeJSONResponse,
		opts...,
	)

	listPoliciesHandler := kithttp.NewServer(
		makeEndpoint(makeListPoliciesEndpoint),
		decodeListPoliciesRequest,
		kithttp.EncodeJSONResponse,
		opts...,
	)

	r := mux.NewRouter()

	r.Handle("/v1/policies/", listPoliciesHandler).Methods("GET")
	r.Handle("/v1/policies/", createPolicyHandler).Methods("POST")
	r.Handle("/v1/policies/{id}", loadPolicyHandler).Methods("GET")
	r.Handle("/v1/policies/{id}", updatePolicyHandler).Methods("PUT")
	r.Handle("/v1/policies/{id}", deletePolicyHandler).Methods("DELETE")

	return r
}

func decodeCreatePolicyRequest(ctx context.Context, r *http.Request) (interface{}, error) {
	var req createPolicyRequest

	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		return nil, err
	}

	return req, nil
}

func decodeDeletePolicyRequest(ctx context.Context, r *http.Request) (interface{}, error) {
	urn, err := getPolicyURN(r, "id")
	if err != nil {
		return nil, err
	}

	return deletePolicyRequest{urn}, nil
}

func decodeLoadPolicyRequest(ctx context.Context, r *http.Request) (interface{}, error) {
	urn, err := getPolicyURN(r, "id")
	if err != nil {
		return nil, err
	}

	return loadPolicyRequest{urn}, nil
}

func decodeUpdatePolicyRequest(ctx context.Context, r *http.Request) (interface{}, error) {
	var req updatePolicyRequest
	if err := json.NewDecoder(r.Body).Decode(&req.Policy); err != nil {
		return nil, err
	}

	var err error
	req.URN, err = getPolicyURN(r, "id")
	if err != nil {
		return nil, err
	}

	return req, nil
}

func decodeListPoliciesRequest(ctx context.Context, r *http.Request) (interface{}, error) {
	return listPoliciesRequest{}, nil
}

func getPolicyURN(r *http.Request, key string) (iam.PolicyURN, error) {
	vars := mux.Vars(r)
	id, ok := vars[key]
	if !ok {
		return "", common.NewInvalidArgumentError("bad route")
	}

	urn := iam.PolicyURN(fmt.Sprintf("urn:iam::policy/%s", id))
	return urn, nil
}
