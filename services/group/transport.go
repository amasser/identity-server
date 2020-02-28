package group

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
	"github.com/tierklinik-dobersberg/identity-server/pkg/authn"
	"github.com/tierklinik-dobersberg/identity-server/pkg/common"
	"github.com/tierklinik-dobersberg/identity-server/pkg/enforcer"
	"github.com/tierklinik-dobersberg/identity-server/pkg/iam"
)

const (
	// ActionGroupRead represents the action to read one or more groups.
	ActionGroupRead = "iam:groups.read"

	// ActionGroupWrite is the action to create, delete or update groups.
	ActionGroupWrite = "iam:groups.write"
)

// MakeHandler returns a http.Handler for the group management service
func MakeHandler(s Service, extractor authn.SubjectExtractorFunc, authz enforcer.Enforcer, logger log.Logger) http.Handler {
	opts := []kithttp.ServerOption{
		kithttp.ServerErrorHandler(transport.NewLogErrorHandler(logger)),
		kithttp.ServerErrorEncoder(kithttp.DefaultErrorEncoder),
		kithttp.ServerBefore(kithttp.PopulateRequestContext),
	}

	makeEndpoint := func(action string, factory func(Service) endpoint.Endpoint) endpoint.Endpoint {
		return endpoint.Chain(
			authn.NewAuthenticator(extractor),
			enforcer.NewActionEndpoint(action),
			enforcer.NewEnforcedEndpoint(authz),
		)(factory(s))
	}

	listGroupsHandler := kithttp.NewServer(
		makeEndpoint(ActionGroupRead, makeGetGroupsEndpoint),
		decodeGetGroupsRequest,
		kithttp.EncodeJSONResponse,
		opts...,
	)

	createGroupHandler := kithttp.NewServer(
		makeEndpoint(ActionGroupWrite, makeCreateGroupEndpoint),
		decodeCreateGroupRequest,
		kithttp.EncodeJSONResponse,
		opts...,
	)

	deleteGroupHandler := kithttp.NewServer(
		makeEndpoint(ActionGroupWrite, makeDeleteGroupEndpoint),
		decodeDeleteGroupRequest,
		kithttp.EncodeJSONResponse,
		opts...,
	)

	loadGroupHandler := kithttp.NewServer(
		makeEndpoint(ActionGroupRead, makeLoadGroupEndpoint),
		decodeLoadGroupRequest,
		kithttp.EncodeJSONResponse,
		opts...,
	)

	updateCommentHandler := kithttp.NewServer(
		makeEndpoint(ActionGroupWrite, makeUpdateGroupCommentEndpoint),
		decodeUpdateCommentRequest,
		kithttp.EncodeJSONResponse,
		opts...,
	)

	addMemberHandler := kithttp.NewServer(
		makeEndpoint(ActionGroupWrite, makeAddMemberEndpoint),
		decodeAddMemberRequest,
		kithttp.EncodeJSONResponse,
		opts...,
	)

	deleteMemberHandler := kithttp.NewServer(
		makeEndpoint(ActionGroupWrite, makeDeleteMemberEndpoint),
		decodeDeleteMemberRequest,
		kithttp.EncodeJSONResponse,
		opts...,
	)

	r := mux.NewRouter()

	// swagger:route GET /v1/groups/ groups listGroups
	//
	// List all groups stored and managed by IAM.
	//
	//	Produces:
	//	- application/json
	//
	//	Schemes: http, https
	//
	//	Responses:
	//		default: body:genericError
	//		200: groupList
	r.Handle("/v1/groups/", listGroupsHandler).Methods("GET")

	// swagger:route POST /v1/groups/ groups createGroup
	//
	// Create a new group.
	//
	// 	Produces:
	//	- application/json
	//
	//	Consumes:
	//	-application/json
	//
	//	Schemes: http, https
	//
	//	Parameters:
	//	+	in: body
	//		type: createGroupRequest
	//
	//	Responses:
	//		default: body:genericError
	//		200: createGroupResponse
	r.Handle("/v1/groups/", createGroupHandler).Methods("POST")

	// swagger:route GET /v1/groups/{id} groups getGroup
	//
	// Get a specific group.
	//
	// 	Produces:
	//	- application/json
	//
	//	Schemes: http, https
	//
	//	Parameters:
	//	+	in: path
	//		name: id
	//		description: The ID of the group.
	//
	//	Responses:
	//		default: body:genericError
	//		200: Group
	r.Handle("/v1/groups/{id}", loadGroupHandler).Methods("GET")

	// swagger:route PUT /v1/groups/{id} groups updateGroupComment
	//
	// Update the comment of a specific group.
	//
	//	Produces:
	//	- application/json
	//
	//	Schemes: http, https
	//
	//	Parameters:
	//	+	in: path
	//		name: id
	//		description: The ID of the group.
	//	+	in: body
	//		type: updateGroupCommentRequest
	//
	//	Responses:
	//		default: body:genericError
	//		201: description: Group updated successfully.
	r.Handle("/v1/groups/{id}", updateCommentHandler).Methods("PUT", "PATCH")

	// swagger:route DELETE /v1/groups/{id} groups deleteGroup
	//
	// Delete an existing group and remove all members of that group.
	//
	//	Produces:
	//	- application/json
	//
	//	Schemes: http, https
	//
	//	Parameters:
	//	+	in: path
	//		name: id
	//		description: The ID of the group to delete.
	//
	//	Responses:
	//		default: body:genericError
	//		201: description: Group deleted successfully.
	r.Handle("/v1/groups/{id}", deleteGroupHandler).Methods("DELETE")

	// swagger:route PUT /v1/groups/{id}/members/{user} groups addMemberToGroup
	//
	// Add a new member to a group.
	//
	//	Produces:
	//	- application/json
	//
	//	Schemes: http, https
	//
	//	Parameters:
	//	+	in: path
	//		name: id
	//		description: The ID of the group.
	//	+	in:path
	//		name: user
	//		description: The ID of the user.
	//
	//	Responses:
	//		default: body:genericError
	//		201: description: User added successfully.
	r.Handle("/v1/groups/{id}/members/{user}", addMemberHandler).Methods("PUT")

	// swagger:route DELETE /v1/groups/{id}/members/{user} groups deleteMemberFromGroup
	//
	// Delete a user from a group.
	//
	//	Produces:
	//	- application/json
	//
	//	Schemes: http, https
	//
	// 	Parameters:
	//	+	in: path
	//		name: id
	//		description: The ID of the group.
	//	+	in: path
	//		name: user
	//		description: The ID of the user to remove from the group.
	//
	//	Responses:
	//		default: body:genericError
	//		201: description: User successfully removed from group.
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
