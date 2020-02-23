package group

import (
	"context"

	"github.com/go-kit/kit/endpoint"
	"github.com/tierklinik-dobersberg/identity-server/iam"
)

type createUserRequest struct {
	Name    string `json:"name"`
	Comment string `json:"comment"`
}

type createUserResponse struct {
	URN iam.GroupURN `json:"urn"`
}

func makeCreateGroupEndpoint(s Service) endpoint.Endpoint {
	return func(ctx context.Context, request interface{}) (interface{}, error) {
		req := request.(createUserRequest)
		urn, err := s.Create(ctx, req.Name, req.Comment)
		if err != nil {
			return nil, err
		}
		return createUserResponse{URN: urn}, nil
	}
}

type deleteGroupRequest struct {
	URN iam.GroupURN
}

type deleteGroupResponse struct{}

func makeDeleteGroupEndpoint(s Service) endpoint.Endpoint {
	return func(ctx context.Context, request interface{}) (interface{}, error) {
		req := request.(deleteGroupRequest)
		if err := s.Delete(ctx, req.URN); err != nil {
			return nil, err
		}
		return deleteGroupResponse{}, nil
	}
}

type loadGroupRequest struct {
	URN iam.GroupURN
}
type loadGroupResponse struct {
	iam.Group
}

func makeLoadGroupEndpoint(s Service) endpoint.Endpoint {
	return func(ctx context.Context, request interface{}) (interface{}, error) {
		req := request.(loadGroupRequest)
		grp, err := s.Load(ctx, req.URN)
		if err != nil {
			return nil, err
		}

		return loadGroupResponse{Group: grp}, nil
	}
}

type getGroupsRequest struct{}
type getGroupsResponse struct {
	Groups []iam.Group `json:"groups"`
}

func makeGetGroupsEndpoint(s Service) endpoint.Endpoint {
	return func(ctx context.Context, request interface{}) (interface{}, error) {
		_ = request.(getGroupsRequest)
		grps, err := s.Get(ctx)
		if err != nil {
			return nil, err
		}
		return getGroupsResponse{Groups: grps}, nil
	}
}

type updateGroupCommentRequest struct {
	URN        iam.GroupURN `json:"-"`
	NewComment string       `json:"comment"`
}
type updateGroupCommentResponse struct{}

func makeUpdateGroupCommentEndpoint(s Service) endpoint.Endpoint {
	return func(ctx context.Context, request interface{}) (interface{}, error) {
		req := request.(updateGroupCommentRequest)
		if err := s.UpdateComment(ctx, req.URN, req.NewComment); err != nil {
			return nil, err
		}
		return updateGroupCommentResponse{}, nil
	}
}

type addMemberRequest struct {
	Group iam.GroupURN
	User  iam.UserURN
}

type addMemberResponse struct{}

func makeAddMemberEndpoint(s Service) endpoint.Endpoint {
	return func(ctx context.Context, request interface{}) (interface{}, error) {
		req := request.(addMemberRequest)
		if err := s.AddMember(ctx, req.Group, req.User); err != nil {
			return nil, err
		}
		return addMemberResponse{}, nil
	}
}

type deleteMemberRequest struct {
	Group iam.GroupURN
	User  iam.UserURN
}

type deleteMemberResponse struct{}

func makeDeleteMemberEndpoint(s Service) endpoint.Endpoint {
	return func(ctx context.Context, request interface{}) (interface{}, error) {
		req := request.(deleteMemberRequest)
		if err := s.DeleteMember(ctx, req.Group, req.User); err != nil {
			return nil, err
		}
		return deleteMemberResponse{}, nil
	}
}
