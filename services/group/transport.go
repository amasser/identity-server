package group

import (
	"context"
	"encoding/json"
	"fmt"
	"net/http"

	"github.com/go-kit/kit/log"
	"github.com/go-kit/kit/transport"
	kithttp "github.com/go-kit/kit/transport/http"
	"github.com/gorilla/mux"
	"github.com/tierklinik-dobersberg/identity-server/iam"
	"github.com/tierklinik-dobersberg/identity-server/pkg/common"
)

// MakeHandler returns a http.Handler for the group management service
func MakeHandler(s Service, logger log.Logger) http.Handler {
	opts := []kithttp.ServerOption{
		kithttp.ServerErrorHandler(transport.NewLogErrorHandler(logger)),
		kithttp.ServerErrorEncoder(kithttp.DefaultErrorEncoder),
	}

	listGroupsHandler := kithttp.NewServer(
		makeGetGroupsEndpoint(s),
		decodeGetGroupsRequest,
		kithttp.EncodeJSONResponse,
		opts...,
	)

	createGroupHandler := kithttp.NewServer(
		makeCreateGroupEndpoint(s),
		decodeCreateGroupRequest,
		kithttp.EncodeJSONResponse,
		opts...,
	)

	deleteGroupHandler := kithttp.NewServer(
		makeDeleteGroupEndpoint(s),
		decodeDeleteGroupRequest,
		kithttp.EncodeJSONResponse,
		opts...,
	)

	loadGroupHandler := kithttp.NewServer(
		makeLoadGroupEndpoint(s),
		decodeLoadGroupRequest,
		kithttp.EncodeJSONResponse,
		opts...,
	)

	updateCommentHandler := kithttp.NewServer(
		makeUpdateGroupCommentEndpoint(s),
		decodeUpdateCommentRequest,
		kithttp.EncodeJSONResponse,
		opts...,
	)

	addMemberHandler := kithttp.NewServer(
		makeAddMemberEndpoint(s),
		decodeAddMemberRequest,
		kithttp.EncodeJSONResponse,
		opts...,
	)

	deleteMemberHandler := kithttp.NewServer(
		makeDeleteMemberEndpoint(s),
		decodeDeleteMemberRequest,
		kithttp.EncodeJSONResponse,
		opts...,
	)

	r := mux.NewRouter()

	r.Handle("/v1/groups/", listGroupsHandler).Methods("GET")
	r.Handle("/v1/groups/", createGroupHandler).Methods("POST")
	r.Handle("/v1/groups/{id}", loadGroupHandler).Methods("GET")
	r.Handle("/v1/groups/{id}", updateCommentHandler).Methods("PUT", "PATCH")
	r.Handle("/v1/groups/{id}", deleteGroupHandler).Methods("DELETE")
	r.Handle("/v1/groups/{id}/members/{user}", addMemberHandler).Methods("PUT")
	r.Handle("/v1/groups/{id}/members/{user}", deleteMemberHandler).Methods("DELETE")

	return r
}

func decodeGetGroupsRequest(ctx context.Context, req *http.Request) (interface{}, error) {
	return getGroupsRequest{}, nil
}

func decodeCreateGroupRequest(ctx context.Context, req *http.Request) (interface{}, error) {
	var request createGroupRequest
	if err := json.NewDecoder(req.Body).Decode(&request); err != nil {
		return nil, err
	}

	return request, nil
}

func decodeDeleteGroupRequest(ctx context.Context, req *http.Request) (interface{}, error) {
	urn, err := getGroupURN(req, "id")
	if err != nil {
		return nil, err
	}

	return deleteGroupRequest{URN: urn}, nil
}

func decodeLoadGroupRequest(ctx context.Context, req *http.Request) (interface{}, error) {
	urn, err := getGroupURN(req, "id")
	if err != nil {
		return nil, err
	}

	return loadGroupRequest{URN: urn}, nil
}

func decodeUpdateCommentRequest(ctx context.Context, req *http.Request) (interface{}, error) {
	urn, err := getGroupURN(req, "id")
	if err != nil {
		return nil, err
	}

	update := updateGroupCommentRequest{URN: urn}
	if err := json.NewDecoder(req.Body).Decode(&update); err != nil {
		return nil, err
	}
	return update, nil
}

func decodeAddMemberRequest(ctx context.Context, req *http.Request) (interface{}, error) {
	grp, err := getGroupURN(req, "id")
	if err != nil {
		return nil, err
	}
	user, err := getUserURN(req, "user")
	if err != nil {
		return nil, err
	}

	return addMemberRequest{Group: grp, User: user}, nil
}

func decodeDeleteMemberRequest(ctx context.Context, req *http.Request) (interface{}, error) {
	grp, err := getGroupURN(req, "id")
	if err != nil {
		return nil, err
	}
	user, err := getUserURN(req, "user")
	if err != nil {
		return nil, err
	}

	return deleteMemberRequest{Group: grp, User: user}, nil
}

func getGroupURN(r *http.Request, key string) (iam.GroupURN, error) {
	vars := mux.Vars(r)
	id, ok := vars[key]
	if !ok {
		return "", common.NewInvalidArgumentError("bad route")
	}

	urn := iam.GroupURN(fmt.Sprintf("urn:iam::group/%s", id))
	return urn, nil
}

func getUserURN(r *http.Request, key string) (iam.UserURN, error) {
	vars := mux.Vars(r)
	id, ok := vars[key]
	if !ok {
		return "", common.NewInvalidArgumentError("bad route")
	}

	urn := iam.UserURN(fmt.Sprintf("urn:iam::user/%s", id))
	return urn, nil
}
