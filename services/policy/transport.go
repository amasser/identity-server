package policy

import (
	"context"
	"encoding/json"
	"fmt"
	"net/http"

	"github.com/gorilla/mux"
	"github.com/tierklinik-dobersberg/identity-server/iam"
	"github.com/tierklinik-dobersberg/identity-server/pkg/common"
)

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
