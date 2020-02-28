package group

import (
	"context"
	"net/http"

	"github.com/go-kit/kit/endpoint"
	"github.com/tierklinik-dobersberg/identity-server/pkg/iam"
)

// Creates a new group in IAM:
// swagger:model createGroupRequest
type createGroupRequest struct {
	// Name is the name of the new group.
	// Required: true
	Name string `json:"name"`

	// Comment is an optional comment for the new group.
	Comment string `json:"comment"`
}

// Response used in a successful call to createGroup
// swagger:model createGroupResponse
type createGroupResponse struct {
	// URN is the URN of the newly created group.
	URN iam.GroupURN `json:"urn"`
}

func makeCreateGroupEndpoint(s Service) endpoint.Endpoint {
	return func(ctx context.Context, request interface{}) (interface{}, error) {
		req := request.(createGroupRequest)
		urn, err := s.Create(ctx, req.Name, req.Comment)
		if err != nil {
			return nil, err
		}
		return createGroupResponse{URN: urn}, nil
	}
}

type deleteGroupRequest struct {
	URN iam.GroupURN
}

type deleteGroupResponse struct{}

func (d deleteGroupResponse) StatusCode() int {
	return http.StatusNoContent
}

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

// A list of all groups stored and managed by IAM.
// swagger:model groupList
type getGroupsResponse struct {
	// All groups stored and managed by IAM.
	Groups []iam.Group `json:"groups,omitempty"`
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

// Request body when updating a group comment.
// swagger:model updateGroupCommentRequest
type updateGroupCommentRequest struct {
	// URN is the URN of the group to update.
	// swagger:ignore
	URN iam.GroupURN `json:"-"`

	// The new comment for the group.
	NewComment string `json:"comment"`
}
type updateGroupCommentResponse struct{}

func (d updateGroupCommentResponse) StatusCode() int {
	return http.StatusNoContent
}

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

func (d addMemberResponse) StatusCode() int {
	return http.StatusNoContent
}

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

func (d deleteMemberResponse) StatusCode() int {
	return http.StatusNoContent
}

func makeDeleteMemberEndpoint(s Service) endpoint.Endpoint {
	return func(ctx context.Context, request interface{}) (interface{}, error) {
		req := request.(deleteMemberRequest)
		if err := s.DeleteMember(ctx, req.Group, req.User); err != nil {
			return nil, err
		}
		return deleteMemberResponse{}, nil
	}
}
