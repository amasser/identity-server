package user

import (
	"context"

	"github.com/go-kit/kit/endpoint"
	"github.com/tierklinik-dobersberg/iam/v2/iam"
)

type createUserRequest struct {
	iam.User
	Password string `json:"password"`
}
type createUserResponse struct {
	iam.User
	Err error `json:"error,omitempty"`
}

func (r createUserResponse) error() error { return r.Err }

func makeCreateUserEndpoint(s Service) endpoint.Endpoint {
	return func(ctx context.Context, request interface{}) (interface{}, error) {
		req := request.(createUserRequest)

		urn, err := s.CreateUser(ctx, req.User.AccountID, req.User.Username, req.User.Attributes)
		if err != nil {
			return createUserResponse{Err: err}, nil
		}

		user, err := s.LoadUser(ctx, urn)
		return createUserResponse{User: user, Err: err}, nil
	}
}

type loadUserRequest struct {
	URN iam.UserURN
}
type loadUserResponse struct {
	iam.User
	Err error `json:"error,omitempty"`
}

func (r loadUserResponse) error() error { return r.Err }

func makeLoadUserEndpoint(s Service) endpoint.Endpoint {
	return func(ctx context.Context, request interface{}) (interface{}, error) {
		req := request.(loadUserRequest)

		user, err := s.LoadUser(ctx, req.URN)
		return loadUserResponse{User: user, Err: err}, nil
	}
}

type listUsersRequest struct{}
type listUsersResponse struct {
	Users []iam.User `json:"users,omitempty"`
	Err   error      `json:"error,omitempty"`
}

func (r listUsersResponse) error() error { return r.Err }

func makeListUsersEndpoint(s Service) endpoint.Endpoint {
	return func(ctx context.Context, request interface{}) (interface{}, error) {
		_ = request.(listUsersRequest)
		users, err := s.Users(ctx)
		return listUsersResponse{Users: users, Err: err}, nil
	}
}

type updateAttrsRequest struct {
	URN        iam.UserURN
	Attributes map[string]interface{}
}
type updateAttrsResponse struct {
	Err error `json:"error,omitempty"`
}

func makeUpdateAttrsEndpoint(s Service) endpoint.Endpoint {
	return func(ctx context.Context, request interface{}) (interface{}, error) {
		req := request.(updateAttrsRequest)
		err := s.UpdateAttrs(ctx, req.URN, req.Attributes)
		return updateAttrsResponse{Err: err}, nil
	}
}

type setAttrRequest struct {
	URN   iam.UserURN
	Key   string
	Value interface{}
}
type setAttrResponse struct {
	Err error `json:"error,omitempty"`
}

func (r setAttrResponse) error() error { return r.Err }

func makeSetAttrEndpoint(s Service) endpoint.Endpoint {
	return func(ctx context.Context, request interface{}) (interface{}, error) {
		req := request.(setAttrRequest)
		err := s.SetAttr(ctx, req.URN, req.Key, req.Value)
		return setAttrResponse{Err: err}, nil
	}
}

type deleteAttrRequest struct {
	URN iam.UserURN
	Key string
}
type deleteAttrResponse struct {
	Err error `json:"error,omitempty"`
}

func (r deleteAttrResponse) error() error { return r.Err }

func makeDeleteAttrRequest(s Service) endpoint.Endpoint {
	return func(ctx context.Context, request interface{}) (interface{}, error) {
		req := request.(deleteAttrRequest)
		err := s.DeleteAttr(ctx, req.URN, req.Key)
		return deleteAttrResponse{Err: err}, nil
	}
}
