package user

import (
	"context"

	"github.com/go-kit/kit/endpoint"
	"github.com/tierklinik-dobersberg/identity-server/iam"
)

// Request body used to create a new user account
// swagger:model createUserBody
type createUserRequest struct {
	Username   string                 `json:"username"`
	Attributes map[string]interface{} `json:"attrs"`
	Password   string                 `json:"password"`
}
type createUserResponse struct {
	iam.User
	Err error `json:"error,omitempty"`
}

func (r createUserResponse) error() error { return r.Err }

func makeCreateUserEndpoint(s Service) endpoint.Endpoint {
	return func(ctx context.Context, request interface{}) (interface{}, error) {
		req := request.(createUserRequest)

		urn, err := s.CreateUser(ctx, req.Username, req.Password, req.Attributes)
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
	*iam.User
	Err error `json:"error,omitempty"`
}

func (r loadUserResponse) error() error { return r.Err }

func makeLoadUserEndpoint(s Service) endpoint.Endpoint {
	return func(ctx context.Context, request interface{}) (interface{}, error) {
		req := request.(loadUserRequest)

		user, err := s.LoadUser(ctx, req.URN)
		if err != nil {
			return loadUserResponse{Err: err}, nil
		}
		return loadUserResponse{User: &user, Err: err}, nil
	}
}

type deleteUserRequest struct {
	URN iam.UserURN
}
type deleteUserResponse struct {
	Err error `json:"error,omitempty"`
}

func (r deleteUserResponse) error() error { return r.Err }

func makeDeleteUserEndpoint(s Service) endpoint.Endpoint {
	return func(ctx context.Context, request interface{}) (interface{}, error) {
		req := request.(deleteUserRequest)
		return deleteUserResponse{Err: s.DeleteUser(ctx, req.URN)}, nil
	}
}

type lockUserRequest struct {
	URN    iam.UserURN
	Locked bool
}
type lockUserResponse struct {
	Err error `json:"error,omitempty"`
}

func (r lockUserResponse) error() error { return r.Err }

func makeLockUserEndpoint(s Service) endpoint.Endpoint {
	return func(ctx context.Context, request interface{}) (interface{}, error) {
		req := request.(lockUserRequest)
		return lockUserResponse{Err: s.LockUser(ctx, req.URN, req.Locked)}, nil
	}
}

type listUsersRequest struct{}

// A list of users accounts
// swagger:model userList
type listUsersResponse struct {
	// All users accounts stored in IAM
	Users []iam.User `json:"users,omitempty"`

	// swagger:ignore
	Err error `json:"error,omitempty"`
}

func (r listUsersResponse) error() error { return r.Err }

func makeListUsersEndpoint(s Service) endpoint.Endpoint {
	return func(ctx context.Context, request interface{}) (interface{}, error) {
		_ = request.(listUsersRequest)
		users, err := s.Users(ctx)
		return listUsersResponse{Users: users, Err: err}, nil
	}
}

// Updates and replaces all attributes of a user account.
// swagger:parameters updateAttributes
type updateAttrsRequest struct {
	// swagger:ignore
	URN iam.UserURN

	// Attributes for the user
	// in: body
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

// Sets a user attribute to a specific value
// swagger:parameters setAttribute
type setAttrRequest struct {
	// swagger:ignore
	URN iam.UserURN

	// The name of the attribute
	// in: path
	// required: true
	Key string `json:"key"`

	// The value for the new attribute
	// in: body
	// required: true
	Value interface{} `json:"value"`
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

// Sets a user attribute to a specific value
// swagger:parameters deleteAttribute
type deleteAttrRequest struct {
	// swagger:ignore
	URN iam.UserURN

	// The name of the attribute
	// in: path
	// required: true
	Key string `json:"key"`
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
