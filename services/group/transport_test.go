package group

import (
	"net/http/httptest"
	"strings"
	"testing"

	"github.com/go-kit/kit/log"
	"github.com/gorilla/mux"
	"github.com/stretchr/testify/assert"
)

func Test_decodeGetGroupsRequest(t *testing.T) {
	r := httptest.NewRequest("GET", "/v1/groups/", nil)
	res, err := decodeGetGroupsRequest(testCtx, r)
	assert.NoError(t, err)
	assert.Equal(t, getGroupsRequest{}, res)
}

func Test_decodeDecodeCreateGroupRequest(t *testing.T) {
	r := httptest.NewRequest("POST", "/v1/groups/", strings.NewReader(`
	{
		"name": "group-name",
		"comment": "and comment"	
	}`))
	res, err := decodeCreateGroupRequest(testCtx, r)
	assert.NoError(t, err)
	assert.Equal(t, createGroupRequest{
		Name:    "group-name",
		Comment: "and comment",
	}, res)

	r = httptest.NewRequest("POST", "/v1/groups/", strings.NewReader(`invalid-json`))
	res, err = decodeCreateGroupRequest(testCtx, r)
	assert.Nil(t, res)
	assert.Error(t, err)
}

func Test_decodeDeleteGroupRequest(t *testing.T) {
	r := httptest.NewRequest("DELETE", "/v1/groups/admins", nil)
	r = mux.SetURLVars(r, map[string]string{"id": "admins"})

	res, err := decodeDeleteGroupRequest(testCtx, r)
	assert.NoError(t, err)
	assert.Equal(t, deleteGroupRequest{
		URN: "urn:iam::group/admins",
	}, res)

	r = mux.SetURLVars(r, map[string]string{})
	res, err = decodeDeleteGroupRequest(testCtx, r)
	assert.Nil(t, res)
	assert.Error(t, err)
}

func Test_decodeLoadGroupRequest(t *testing.T) {
	r := mux.SetURLVars(
		httptest.NewRequest("GET", "/v1/groups/admins", nil),
		map[string]string{"id": "admins"},
	)
	res, err := decodeLoadGroupRequest(testCtx, r)
	assert.NoError(t, err)
	assert.Equal(t, loadGroupRequest{URN: "urn:iam::group/admins"}, res)

	r = mux.SetURLVars(r, nil)
	res, err = decodeLoadGroupRequest(testCtx, r)
	assert.Nil(t, res)
	assert.Error(t, err)
}

func Test_decodeAddMemberRequest(t *testing.T) {
	r := mux.SetURLVars(
		httptest.NewRequest("PUT", "/v1/groups/admins/member/10", nil),
		map[string]string{
			"id":   "admins",
			"user": "10",
		},
	)

	res, err := decodeAddMemberRequest(testCtx, r)
	assert.NoError(t, err)
	assert.Equal(t, addMemberRequest{
		Group: "urn:iam::group/admins",
		User:  "urn:iam::user/10",
	}, res)

	r1 := mux.SetURLVars(r, map[string]string{
		"id": "admins",
	})
	res, err = decodeAddMemberRequest(testCtx, r1)
	assert.Nil(t, res)
	assert.Error(t, err)

	r2 := mux.SetURLVars(r, map[string]string{
		"user": "10",
	})
	res, err = decodeAddMemberRequest(testCtx, r2)
	assert.Nil(t, res)
	assert.Error(t, err)
}

func Test_decodeDeleteMemberRequest(t *testing.T) {
	r := mux.SetURLVars(
		httptest.NewRequest("PUT", "/v1/groups/admins/member/10", nil),
		map[string]string{
			"id":   "admins",
			"user": "10",
		},
	)

	res, err := decodeDeleteMemberRequest(testCtx, r)
	assert.NoError(t, err)
	assert.Equal(t, deleteMemberRequest{
		Group: "urn:iam::group/admins",
		User:  "urn:iam::user/10",
	}, res)

	r1 := mux.SetURLVars(r, map[string]string{
		"id": "admins",
	})
	res, err = decodeDeleteMemberRequest(testCtx, r1)
	assert.Nil(t, res)
	assert.Error(t, err)

	r2 := mux.SetURLVars(r, map[string]string{
		"user": "10",
	})
	res, err = decodeDeleteMemberRequest(testCtx, r2)
	assert.Nil(t, res)
	assert.Error(t, err)
}

func Test_MakeHandler(t *testing.T) {
	s := &mockService{}
	_ = MakeHandler(s, log.NewNopLogger())
}

func Test_decodeUpdateCommentRequest(t *testing.T) {
	r := mux.SetURLVars(
		httptest.NewRequest("PUT", "/v1/groups/admins", strings.NewReader(`
			{
				"comment": "new comment"
			}
		`)),
		map[string]string{"id": "admins"},
	)
	res, err := decodeUpdateCommentRequest(testCtx, r)
	assert.NoError(t, err)
	assert.Equal(t, updateGroupCommentRequest{
		URN:        "urn:iam::group/admins",
		NewComment: "new comment",
	}, res)

	r = mux.SetURLVars(r, nil)
	res, err = decodeUpdateCommentRequest(testCtx, r)
	assert.Error(t, err)
	assert.Nil(t, res)
}
