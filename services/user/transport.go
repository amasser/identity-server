package user

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"net/http"
	"os"

	"github.com/go-kit/kit/endpoint"
	"github.com/go-kit/kit/log"
	"github.com/go-kit/kit/transport"
	kithttp "github.com/go-kit/kit/transport/http"
	"github.com/gorilla/mux"
	"github.com/tierklinik-dobersberg/identity-server/iam"
	"github.com/tierklinik-dobersberg/identity-server/pkg/authn"
	"github.com/tierklinik-dobersberg/identity-server/pkg/enforcer"
)

const (
	// ActionWriteUser allows the subject to create new users.
	ActionWriteUser = "iam:user:write"

	// ActionLoadUser allows the subject to load a user.
	ActionLoadUser = "iam:user:load"

	// ActionListUsers allows the subject ot list users.
	ActionListUsers = "iam:user:list"

	// ActionDeleteUser allows the subject to delete a user.
	ActionDeleteUser = "iam:user:delete"

	// ActionLockUnlockUser allows the subject to lock or unlock a user account.
	ActionLockUnlockUser = "iam:user:lock-unlock"

	// ActionUpdateUserAttr allows the subject to update a users attributes.
	ActionUpdateUserAttr = "iam:user:write-attr"
)

// MakeHandler returns a http.Handler for the user management service
func MakeHandler(s Service, extractor authn.SubjectExtractorFunc, authz enforcer.Enforcer, logger log.Logger) http.Handler {
	opts := []kithttp.ServerOption{
		kithttp.ServerErrorHandler(transport.NewLogErrorHandler(logger)),
		kithttp.ServerErrorEncoder(kithttp.DefaultErrorEncoder),
		kithttp.ServerBefore(kithttp.PopulateRequestContext),
	}

	makeEndpoint := func(action string, factory func(s Service) endpoint.Endpoint) endpoint.Endpoint {
		return endpoint.Chain(
			authn.NewAuthenticator(extractor),
			enforcer.NewActionEndpoint(action),
			enforcer.NewEnforcedEndpoint(authz),
		)(factory(s))
	}

	createUserHandler := kithttp.NewServer(
		makeEndpoint(ActionWriteUser, makeCreateUserEndpoint),
		decodeCreateUserRequest,
		encodeResponse,
		opts...,
	)

	loadUserHandler := kithttp.NewServer(
		makeEndpoint(ActionLoadUser, makeLoadUserEndpoint),
		decodeLoadUserRequest,
		encodeResponse,
		opts...,
	)

	deleteUserHandler := kithttp.NewServer(
		makeEndpoint(ActionDeleteUser, makeDeleteUserEndpoint),
		decodeDeleteUserRequest,
		encodeStatusOnlyResponse,
		opts...,
	)

	lockUserHandler := kithttp.NewServer(
		makeEndpoint(ActionLockUnlockUser, makeLockUserEndpoint),
		decodeLockUserRequest,
		encodeStatusOnlyResponse,
		opts...,
	)

	listUsersHandler := kithttp.NewServer(
		makeEndpoint(ActionListUsers, makeListUsersEndpoint),
		decodeListUserRequest,
		encodeResponse,
		opts...,
	)

	updateAttrHandler := kithttp.NewServer(
		makeEndpoint(ActionUpdateUserAttr, makeUpdateAttrsEndpoint),
		decodeUpdateAttrRequest,
		encodeStatusOnlyResponse,
		opts...,
	)

	setAttrHandler := kithttp.NewServer(
		makeEndpoint(ActionUpdateUserAttr, makeSetAttrEndpoint),
		decodeSetAttrRequest,
		encodeStatusOnlyResponse,
		opts...,
	)

	deleteAttrHandler := kithttp.NewServer(
		makeEndpoint(ActionUpdateUserAttr, makeDeleteAttrRequest),
		decodeDeleteAttrRequest,
		encodeStatusOnlyResponse,
		opts...,
	)

	r := mux.NewRouter()

	// swagger:route GET /v1/users/ users listUsers
	//
	// List all users accounts stored in IAM.
	//
	//     Produces:
	//     - application/json
	//
	//     Schemes: http, https
	//
	//     Responses:
	//       default: body:genericError
	//       200: userList
	r.Handle("/v1/users/", listUsersHandler).Methods("GET")

	// swagger:route POST /v1/users/ users createUser
	//
	// Creates a new user account on IAM and authn-server.
	//
	//     Consumes:
	//     - application/json
	//
	//     Produces:
	//     - application/json
	//
	//     Schemes: http, https
	//
	//     Parameters:
	//     + in: body
	//       type: createUserBody
	//
	//     Responses:
	//       default: body:genericError
	//       200: User
	r.Handle("/v1/users/", createUserHandler).Methods("POST")

	// swagger:route GET /v1/users/{id} user getUser
	//
	// Returns a user account identified by it's ID
	//
	//     Produces:
	//	   - application/json
	//
	//	   Schemes: http, https
	//
	//	   Parameters:
	//     + name: id
	//       in: path
	//       type: number
	//	     required: true
	//
	//     Responses:
	//       default: body:genericError
	//       200: User
	r.Handle("/v1/users/{id}", loadUserHandler).Methods("GET")

	// swagger:route DELETE /v1/users/{id} user deleteUser
	//
	// Deletes a user account from IAM and archives it on authn-server.
	//
	//     Schemes: http, https
	//
	//     Parameters:
	//     + name: id
	//       in: path
	//       type: number
	//       required: true
	//
	//     Responses:
	//       default: body:genericError
	//       202: description:User has been deleted successfully
	r.Handle("/v1/users/{id}", deleteUserHandler).Methods("DELETE")

	// swagger:route PUT /v1/users/{id}/locked user lockUser
	//
	// Locks a user account.
	//
	//     Schemes: http, https
	//
	//     Parameters:
	//     + name: id
	//       in: path
	//       type: number
	//       required: true
	//
	//     Responses:
	//       default: body:genericError
	//       202: description:User has been locked successfully
	//

	// swagger:route DELETE /v1/users/{id}/locked user lockUser
	//
	// Unlock a user account.
	//
	//     Schemes: http, https
	//
	//     Parameters:
	//     + name: id
	//       in: path
	//       type: number
	//       required: true
	//
	//     Responses:
	//       default: body:genericError
	//       202: description:User has been unlocked successfully
	r.Handle("/v1/users/{id}/locked", lockUserHandler).Methods("PUT", "DELETE")

	// swagger:route PUT /v1/users/{id}/attrs/ user attributes updateAttributes
	//
	// Replaces all attributes of a user account.
	//
	//     Schemes: http, https
	//
	//     Consumes:
	//     - application/json
	//
	//     Parameters:
	//     + name: id
	//       in: path
	//       type: number
	//       required: true
	//
	//     Responses:
	//       default: body:genericError
	//       202: description:User attributes have been successfully replaced
	r.Handle("/v1/users/{id}/attrs/", updateAttrHandler).Methods("PUT")

	// swagger:route PUT /v1/users/{id}/attrs/{key} user attributes setAttribute
	//
	// Updates a single user attribute.
	//
	//     Schemes: http, https
	//
	//     Parameters:
	//     + name: id
	//       in: path
	//       type: number
	//       required: true
	//
	//     Responses:
	//       default: body:genericError
	//       202: description:The user attribute has been successfully stored
	r.Handle("/v1/users/{id}/attrs/{key}", setAttrHandler).Methods("PUT")

	// swagger:route DELETE /v1/users/{id}/attrs/{key} user attributes deleteAttribute
	//
	// Updates a single user attribute.
	//
	//     Schemes: http, https
	//
	//     Parameters:
	//     + name: id
	//       in: path
	//       type: number
	//       required: true
	//
	//     Responses:
	//       default: body:genericError
	//       202: description:The user attribute has been successfully stored
	r.Handle("/v1/users/{id}/attrs/{key}", deleteAttrHandler).Methods("DELETE")

	return r
}

var errBadRoute = errors.New("bad route")

func decodeCreateUserRequest(_ context.Context, r *http.Request) (interface{}, error) {
	var req createUserRequest

	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		return nil, err
	}

	return req, nil
}

func decodeLoadUserRequest(_ context.Context, r *http.Request) (interface{}, error) {
	urn, err := getURNFromVars(r, "id")
	return loadUserRequest{URN: urn}, err
}

func decodeDeleteUserRequest(_ context.Context, r *http.Request) (interface{}, error) {
	urn, err := getURNFromVars(r, "id")
	return deleteUserRequest{URN: urn}, err
}

func decodeLockUserRequest(_ context.Context, r *http.Request) (interface{}, error) {
	urn, err := getURNFromVars(r, "id")
	if err != nil {
		return nil, err
	}

	locked := r.Method == "PUT"
	return lockUserRequest{URN: urn, Locked: locked}, nil
}

func decodeListUserRequest(_ context.Context, r *http.Request) (interface{}, error) {
	return listUsersRequest{}, nil
}

func decodeUpdateAttrRequest(_ context.Context, r *http.Request) (interface{}, error) {
	var (
		req updateAttrsRequest
		err error
	)

	req.URN, err = getURNFromVars(r, "id")
	if err != nil {
		return nil, err
	}

	if err := json.NewDecoder(r.Body).Decode(&req.Attributes); err != nil {
		return nil, err
	}

	return req, nil
}

func decodeSetAttrRequest(_ context.Context, r *http.Request) (interface{}, error) {
	var (
		req setAttrRequest
		err error
		ok  bool
	)

	req.URN, err = getURNFromVars(r, "id")
	if err != nil {
		return nil, err
	}

	vars := mux.Vars(r)
	req.Key, ok = vars["key"]
	if !ok {
		return nil, errBadRoute
	}

	if err := json.NewDecoder(r.Body).Decode(&req.Value); err != nil {
		return nil, err
	}

	return req, nil
}

func decodeDeleteAttrRequest(_ context.Context, r *http.Request) (interface{}, error) {
	var (
		req deleteAttrRequest
		err error
		ok  bool
	)

	req.URN, err = getURNFromVars(r, "id")
	if err != nil {
		return nil, err
	}

	vars := mux.Vars(r)
	req.Key, ok = vars["key"]
	if !ok {
		return nil, errBadRoute
	}

	return req, nil
}

func encodeStatusOnlyResponse(ctx context.Context, w http.ResponseWriter, response interface{}) error {
	if e, ok := response.(errorer); ok && e.error() != nil {
		encodeError(ctx, e.error(), w)
		return nil
	}

	w.WriteHeader(http.StatusAccepted)

	return nil
}

func encodeResponse(ctx context.Context, w http.ResponseWriter, response interface{}) error {
	if e, ok := response.(errorer); ok && e.error() != nil {
		encodeError(ctx, e.error(), w)
		return nil
	}

	w.Header().Set("Content-Type", "application/json; charset=utf-8")
	return json.NewEncoder(w).Encode(response)
}

type errorer interface {
	error() error
}

// A generic error message
// swagger:model genericError
type genericError struct {
	// The error message that occured
	//
	// Required: true
	// Example: Invalid parameter
	Error string `json:"error"`
}

// encode errors from business-logic
func encodeError(_ context.Context, err error, w http.ResponseWriter) {
	w.Header().Set("Content-Type", "application/json; charset=utf-8")
	switch err {
	case os.ErrNotExist:
		w.WriteHeader(http.StatusNotFound)
		err = errors.New("resource not found")
	case os.ErrExist:
		w.WriteHeader(http.StatusConflict)
		err = errors.New("resource exists")
	case ErrInvalidArgument:
		w.WriteHeader(http.StatusBadRequest)
	default:
		w.WriteHeader(http.StatusInternalServerError)
	}
	json.NewEncoder(w).Encode(map[string]interface{}{
		"error": err.Error(),
	})
}

func getURNFromVars(r *http.Request, key string) (iam.UserURN, error) {
	vars := mux.Vars(r)
	id, ok := vars[key]
	if !ok {
		return "", errBadRoute
	}

	// TODO(ppacher): convert to number for input validation
	urn := iam.UserURN(fmt.Sprintf("urn:iam::user/%s", id))
	return urn, nil
}
