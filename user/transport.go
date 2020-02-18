package user

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"net/http"
	"os"

	"github.com/go-kit/kit/log"
	"github.com/go-kit/kit/transport"
	kithttp "github.com/go-kit/kit/transport/http"
	"github.com/gorilla/mux"
	"github.com/tierklinik-dobersberg/iam/v2/iam"
)

// MakeHandler returns a http.Handler for the user management service
func MakeHandler(s Service, logger log.Logger) http.Handler {
	opts := []kithttp.ServerOption{
		kithttp.ServerErrorHandler(transport.NewLogErrorHandler(logger)),
	}

	createUserHandler := kithttp.NewServer(
		makeCreateUserEndpoint(s),
		decodeCreateUserRequest,
		encodeResponse,
		opts...,
	)

	loadUserHandler := kithttp.NewServer(
		makeLoadUserEndpoint(s),
		decodeLoadUserRequest,
		encodeResponse,
		opts...,
	)

	listUsersHandler := kithttp.NewServer(
		makeListUsersEndpoint(s),
		decodeListUserRequest,
		encodeResponse,
		opts...,
	)

	updateAttrHandler := kithttp.NewServer(
		makeUpdateAttrsEndpoint(s),
		decodeUpdateAttrRequest,
		encodeResponse,
		opts...,
	)

	setAttrHandler := kithttp.NewServer(
		makeSetAttrEndpoint(s),
		decodeSetAttrRequest,
		encodeResponse,
		opts...,
	)

	deleteAttrHandler := kithttp.NewServer(
		makeDeleteAttrRequest(s),
		decodeDeleteAttrRequest,
		encodeResponse,
		opts...,
	)

	r := mux.NewRouter()

	r.Handle("/v1/users", createUserHandler).Methods("POST")
	r.Handle("/v1/users", listUsersHandler).Methods("GET")
	r.Handle("/v1/users/{id}", loadUserHandler).Methods("GET")
	r.Handle("/v1/users/{id}/attrs", updateAttrHandler).Methods("PUT")
	r.Handle("/v1/users/{id}/attrs/{key}", setAttrHandler).Methods("PUT")
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

func statusResponseEncoder(status http.Status)

func encodeResponse(ctx context.Context, w http.ResponseWriter, response interface{}) error {
	if r, ok := response.(httpResponseWriter); ok {
		return r.writeHTTP(ctx, w)
	}

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

type httpResponseWriter interface {
	writeHTTP(context.Context, http.ResponseWriter) error
}

// encode errors from business-logic
func encodeError(_ context.Context, err error, w http.ResponseWriter) {
	w.Header().Set("Content-Type", "application/json; charset=utf-8")
	switch err {
	case os.ErrNotExist:
		w.WriteHeader(http.StatusNotFound)
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
	id, ok := vars["id"]
	if !ok {
		return "", errBadRoute
	}

	// TODO(ppacher): convert to number for input validation
	urn := iam.UserURN(fmt.Sprintf("urn:iam::user/%s", id))
	return urn, nil
}
