package user

import (
	"net/http/httptest"
	"strings"
	"testing"

	"github.com/go-kit/kit/log"
	"github.com/gorilla/mux"
	"github.com/stretchr/testify/assert"
	"github.com/tierklinik-dobersberg/iam/v2/iam"
)

func Test_decodeCreateUserRequest(t *testing.T) {
	payload := `
	{
		"username": "admin",
		"attrs": {
			"job": "developer"
		},
		"accountID": 10,
		"password": "password"
	}
	`
	r := httptest.NewRequest("POST", "/v1/users/", strings.NewReader(payload))

	expectedRequest := createUserRequest{
		User: iam.User{
			AccountID: 10,
			Attributes: map[string]interface{}{
				"job": "developer",
			},
			Username: "admin",
		},
		Password: "password",
	}

	req, err := decodeCreateUserRequest(nil, r)
	assert.NoError(t, err)
	assert.Equal(t, expectedRequest, req)

	r = httptest.NewRequest("POST", "/v1/users/", strings.NewReader("199"))
	req, err = decodeCreateUserRequest(nil, r)
	assert.Error(t, err)
}

func Test_decodeLoadUserRequest(t *testing.T) {
	r := httptest.NewRequest("GET", "/v1/users/10", nil)
	r = mux.SetURLVars(r, map[string]string{"id": "10"})

	expectedRequest := loadUserRequest{
		URN: "urn:iam::user/10",
	}

	req, err := decodeLoadUserRequest(nil, r)
	assert.NoError(t, err)
	assert.Equal(t, expectedRequest, req)

	r = mux.SetURLVars(r, nil)
	req, err = decodeLoadUserRequest(nil, r)
	assert.Error(t, err)
}

func Test_decodeListUserRequest(t *testing.T) {
	r := httptest.NewRequest("GET", "/v1/users", nil)
	req, err := decodeListUserRequest(nil, r)
	assert.NoError(t, err)
	assert.Equal(t, listUsersRequest{}, req)
}

func Test_updateAttrRequest(t *testing.T) {
	payload := `
	{
		"job": "developer"
	}
	`
	r := httptest.NewRequest("PUT", "/v1/users/10/attrs", strings.NewReader(payload))
	r = mux.SetURLVars(r, map[string]string{"id": "10"})

	expectedRequest := updateAttrsRequest{
		URN: "urn:iam::user/10",
		Attributes: map[string]interface{}{
			"job": "developer",
		},
	}

	req, err := decodeUpdateAttrRequest(nil, r)
	assert.NoError(t, err)
	assert.Equal(t, expectedRequest, req)

	r = httptest.NewRequest("POST", "/v1/users/", strings.NewReader(`"invalid-body"`))
	r = mux.SetURLVars(r, map[string]string{"id": "10"})
	req, err = decodeUpdateAttrRequest(nil, r)
	assert.Error(t, err)

	r = httptest.NewRequest("POST", "/v1/users/", strings.NewReader(`{}`))
	req, err = decodeUpdateAttrRequest(nil, r)
	assert.Error(t, err) // no user id in mux.Vars
}

func Test_setAttrRequest(t *testing.T) {
	payload := `
	"developer"
	`
	r := httptest.NewRequest("PUT", "/v1/users/10/attrs/job", strings.NewReader(payload))
	r = mux.SetURLVars(r, map[string]string{"id": "10", "key": "job"})

	expectedRequest := setAttrRequest{
		URN:   "urn:iam::user/10",
		Key:   "job",
		Value: "developer",
	}

	req, err := decodeSetAttrRequest(nil, r)
	assert.NoError(t, err)
	assert.Equal(t, expectedRequest, req)

	r = httptest.NewRequest("PUT", "/v1/users/10/attrs/job", strings.NewReader(`invalid-json`))
	r = mux.SetURLVars(r, map[string]string{"id": "10", "key": "job"})
	req, err = decodeSetAttrRequest(nil, r)
	assert.Error(t, err)

	r = httptest.NewRequest("PUT", "/v1/users/10/attrs/job", strings.NewReader(`{}`))
	r = mux.SetURLVars(r, map[string]string{"key": "job"})
	req, err = decodeSetAttrRequest(nil, r)
	assert.Error(t, err) // no user id in mux.Vars

	r = httptest.NewRequest("PUT", "/v1/users/10/attrs/job", strings.NewReader(`{}`))
	r = mux.SetURLVars(r, map[string]string{"id": "10"})
	req, err = decodeSetAttrRequest(nil, r)
	assert.Error(t, err) // no user key in mux.Vars
}

func Test_delAttrRequest(t *testing.T) {
	r := httptest.NewRequest("DELETE", "/v1/users/10/attrs/job", nil)
	r = mux.SetURLVars(r, map[string]string{"id": "10", "key": "job"})

	expectedRequest := deleteAttrRequest{
		URN: "urn:iam::user/10",
		Key: "job",
	}

	req, err := decodeDeleteAttrRequest(nil, r)
	assert.NoError(t, err)
	assert.Equal(t, expectedRequest, req)

	r = httptest.NewRequest("DELETE", "/v1/users/10/attrs/job", nil)
	r = mux.SetURLVars(r, map[string]string{"key": "job"})
	req, err = decodeDeleteAttrRequest(nil, r)
	assert.Error(t, err) // no user id in mux.Vars

	r = httptest.NewRequest("DELETE", "/v1/users/10/attrs/job", nil)
	r = mux.SetURLVars(r, map[string]string{"id": "10"})
	req, err = decodeDeleteAttrRequest(nil, r)
	assert.Error(t, err) // no user key in mux.Vars
}

func Test_MakeHandler(t *testing.T) {
	r := MakeHandler(NewService(&userRepoMock{}), log.NewNopLogger())
	assert.NotNil(t, r)
}
