package group

import (
	"context"
	"errors"
	"net/http"
	"testing"

	kithttp "github.com/go-kit/kit/transport/http"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
	"github.com/tierklinik-dobersberg/identity-server/pkg/iam"
)

func Test_CreateGroupEndpoint(t *testing.T) {
	s := &mockService{}
	ep := makeCreateGroupEndpoint(s)

	s.On("Create", "admins", "All IT administrators").Once().Return(
		iam.GroupURN("urn:iam::group/admins"),
		nil,
	)
	res, err := ep(testCtx, createGroupRequest{Name: "admins", Comment: "All IT administrators"})
	assert.NoError(t, err)
	assert.Equal(t, createGroupResponse{URN: "urn:iam::group/admins"}, res)

	s.On("Create", "devs", "some comment").Once().Return(iam.GroupURN(""), errors.New("simualted"))
	res, err = ep(testCtx, createGroupRequest{Name: "devs", Comment: "some comment"})
	assert.Error(t, err)
	assert.Nil(t, res)

	s.AssertExpectations(t)
}

func Test_DeleteGroupEndpoint(t *testing.T) {
	s := &mockService{}
	ep := makeDeleteGroupEndpoint(s)

	s.On("Delete", iam.GroupURN("urn:iam::group/admins")).Once().Return(nil)
	res, err := ep(testCtx, deleteGroupRequest{URN: "urn:iam::group/admins"})
	assert.NoError(t, err)
	assert.Equal(t, deleteGroupResponse{}, res)

	c, ok := res.(kithttp.StatusCoder)
	assert.True(t, ok)
	assert.Equal(t, http.StatusNoContent, c.StatusCode())

	s.On("Delete", iam.GroupURN("urn:iam::group/devs")).Once().Return(errors.New("simulated"))
	res, err = ep(testCtx, deleteGroupRequest{URN: "urn:iam::group/devs"})
	assert.Nil(t, res)
	assert.Error(t, err)
}

func Test_LoadGroupEndpoint(t *testing.T) {
	s := &mockService{}
	ep := makeLoadGroupEndpoint(s)

	s.On("Load", iam.GroupURN("urn:iam::group/devs")).Once().Return(iam.Group{}, errors.New("simulated"))
	res, err := ep(testCtx, loadGroupRequest{URN: "urn:iam::group/devs"})
	assert.Error(t, err)
	assert.Nil(t, res)

	grp := iam.Group{
		ID:   "urn:iam::group/devs",
		Name: "devs",
	}
	s.On("Load", iam.GroupURN("urn:iam::group/devs")).Once().Return(grp, nil)
	res, err = ep(testCtx, loadGroupRequest{URN: "urn:iam::group/devs"})
	assert.NoError(t, err)
	assert.Equal(t, loadGroupResponse{Group: grp}, res)
}

func Test_GetGroupsEndpoint(t *testing.T) {
	s := &mockService{}
	ep := makeGetGroupsEndpoint(s)
	grp := iam.Group{
		ID:   "urn:iam::group/devs",
		Name: "devs",
	}
	s.On("Get").Once().Return([]iam.Group{grp}, nil)

	grps, err := ep(testCtx, getGroupsRequest{})
	assert.NoError(t, err)
	assert.Equal(t, getGroupsResponse{Groups: []iam.Group{grp}}, grps)

	s.On("Get").Once().Return([]iam.Group(nil), errors.New("simulated"))
	grps, err = ep(testCtx, getGroupsRequest{})
	assert.Error(t, err)
	assert.Nil(t, grps)
}

func Test_UpdateGroupCommentEndpoint(t *testing.T) {
	s := &mockService{}
	ep := makeUpdateGroupCommentEndpoint(s)

	s.On("UpdateComment", iam.GroupURN("urn:iam::group/admins"), "new comment").Once().Return(nil)

	res, err := ep(testCtx, updateGroupCommentRequest{URN: "urn:iam::group/admins", NewComment: "new comment"})
	assert.NoError(t, err)
	assert.Equal(t, updateGroupCommentResponse{}, res)

	assert.Equal(t, http.StatusNoContent, res.(kithttp.StatusCoder).StatusCode())

	s.On("UpdateComment", iam.GroupURN("urn:iam::group/admins"), "new comment").Once().Return(errors.New("error"))

	res, err = ep(testCtx, updateGroupCommentRequest{URN: "urn:iam::group/admins", NewComment: "new comment"})
	assert.Error(t, err)
	assert.Nil(t, res)
}

func Test_AddMemberEndpoint(t *testing.T) {
	s := &mockService{}
	ep := makeAddMemberEndpoint(s)

	s.On("AddMember", iam.GroupURN("urn:iam::group/devs"), iam.UserURN("urn:iam::user/10")).Once().Return(nil)
	res, err := ep(testCtx, addMemberRequest{
		Group: "urn:iam::group/devs",
		User:  "urn:iam::user/10",
	})
	assert.NoError(t, err)
	assert.Equal(t, addMemberResponse{}, res)

	assert.Equal(t, http.StatusNoContent, res.(kithttp.StatusCoder).StatusCode())

	s.On("AddMember", iam.GroupURN("urn:iam::group/devs"), iam.UserURN("urn:iam::user/10")).Once().Return(errors.New("error"))
	res, err = ep(testCtx, addMemberRequest{
		Group: "urn:iam::group/devs",
		User:  "urn:iam::user/10",
	})
	assert.Error(t, err)
	assert.Nil(t, res)
}

func Test_DeleteMemberEndpoint(t *testing.T) {
	s := &mockService{}
	ep := makeDeleteMemberEndpoint(s)

	s.On("DeleteMember", iam.GroupURN("urn:iam::group/devs"), iam.UserURN("urn:iam::user/10")).Once().Return(nil)
	res, err := ep(testCtx, deleteMemberRequest{
		Group: "urn:iam::group/devs",
		User:  "urn:iam::user/10",
	})
	assert.NoError(t, err)
	assert.Equal(t, deleteMemberResponse{}, res)

	assert.Equal(t, http.StatusNoContent, res.(kithttp.StatusCoder).StatusCode())

	s.On("DeleteMember", iam.GroupURN("urn:iam::group/devs"), iam.UserURN("urn:iam::user/10")).Once().Return(errors.New("error"))
	res, err = ep(testCtx, deleteMemberRequest{
		Group: "urn:iam::group/devs",
		User:  "urn:iam::user/10",
	})
	assert.Error(t, err)
	assert.Nil(t, res)
}

type mockService struct {
	mock.Mock
}

func (m *mockService) Get(ctx context.Context) ([]iam.Group, error) {
	args := m.Called()
	return args.Get(0).([]iam.Group), args.Error(1)
}

func (m *mockService) Create(ctx context.Context, groupName string, groupComment string) (iam.GroupURN, error) {
	args := m.Called(groupName, groupComment)
	return args.Get(0).(iam.GroupURN), args.Error(1)
}

func (m *mockService) Delete(ctx context.Context, urn iam.GroupURN) error {
	args := m.Called(urn)
	return args.Error(0)
}

func (m *mockService) Load(ctx context.Context, urn iam.GroupURN) (iam.Group, error) {
	args := m.Called(urn)
	return args.Get(0).(iam.Group), args.Error(1)
}

func (m *mockService) UpdateComment(ctx context.Context, urn iam.GroupURN, comment string) error {
	args := m.Called(urn, comment)
	return args.Error(0)
}

func (m *mockService) AddMember(ctx context.Context, grp iam.GroupURN, member iam.UserURN) error {
	args := m.Called(grp, member)
	return args.Error(0)
}

func (m *mockService) DeleteMember(ctx context.Context, grp iam.GroupURN, member iam.UserURN) error {
	args := m.Called(grp, member)
	return args.Error(0)
}
