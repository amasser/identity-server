package user

import (
	"errors"
	"net/http"
	"net/http/httptest"
	"os"
	"strings"
	"testing"

	"github.com/go-kit/kit/log"
	"github.com/gorilla/mux"
	"github.com/stretchr/testify/assert"
	"github.com/tierklinik-dobersberg/identity-server/iam"
	"github.com/tierklinik-dobersberg/identity-server/pkg/enforcer"
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
		Attributes: map[string]interface{}{
			"job": "developer",
		},
		Username: "admin",
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

func Test_decodeDeleteUserRequest(t *testing.T) {
	r := httptest.NewRequest("DELETE", "/v1/users/10", nil)
	r = mux.SetURLVars(r, map[string]string{"id": "10"})

	expectedRequest := deleteUserRequest{
		URN: "urn:iam::user/10",
	}

	req, err := decodeDeleteUserRequest(nil, r)
	assert.NoError(t, err)
	assert.Equal(t, expectedRequest, req)

	r = mux.SetURLVars(r, nil)
	req, err = decodeDeleteUserRequest(nil, r)
	assert.Error(t, err)
}

func Test_decodeLockUserRequest_Locked(t *testing.T) {
	r := httptest.NewRequest("PUT", "/v1/users/10", nil)
	r = mux.SetURLVars(r, map[string]string{"id": "10"})

	expectedRequest := lockUserRequest{
		URN:    "urn:iam::user/10",
		Locked: true,
	}

	req, err := decodeLockUserRequest(nil, r)
	assert.NoError(t, err)
	assert.Equal(t, expectedRequest, req)

	r = mux.SetURLVars(r, nil)
	req, err = decodeLockUserRequest(nil, r)
	assert.Error(t, err)
}

func Test_decodeLockUserRequest_Unlocked(t *testing.T) {
	r := httptest.NewRequest("DELETE", "/v1/users/10", nil)
	r = mux.SetURLVars(r, map[string]string{"id": "10"})

	expectedRequest := lockUserRequest{
		URN:    "urn:iam::user/10",
		Locked: false,
	}

	req, err := decodeLockUserRequest(nil, r)
	assert.NoError(t, err)
	assert.Equal(t, expectedRequest, req)

	r = mux.SetURLVars(r, nil)
	req, err = decodeLockUserRequest(nil, r)
	assert.Error(t, err)
}

func Test_decodeListUserRequest(t *testing.T) {
	r := httptest.NewRequest("GET", "/v1/users", nil)
	req, err := decodeListUserRequest(nil, r)
	assert.NoError(t, err)
	assert.Equal(t, listUsersRequest{}, req)
}

func Test_decodeUpdateAttrRequest(t *testing.T) {
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

func Test_decodeSetAttrRequest(t *testing.T) {
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

func Test_decodeDelAttrRequest(t *testing.T) {
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
	svc, _, _ := setupServiceTestBed()
	jwtTokenExtractor := func(string) (string, error) { return "", nil }
	r := MakeHandler(svc, jwtTokenExtractor, enforcer.NewNoOpEnforcer(), log.NewNopLogger())
	assert.NotNil(t, r)
}

func Test_encodeError(t *testing.T) {
	w := httptest.NewRecorder()
	err := errors.New("internal-server-error")
	encodeError(bg, err, w)
	assert.Equal(t, http.StatusInternalServerError, w.Code)
	assert.Equal(t, "{\"error\":\"internal-server-error\"}\n", string(w.Body.Bytes()))
	assert.Equal(t, []string{"application/json; charset=utf-8"}, w.HeaderMap["Content-Type"])

	w = httptest.NewRecorder()
	encodeError(bg, os.ErrNotExist, w)
	assert.Equal(t, http.StatusNotFound, w.Code)
	assert.Equal(t, "{\"error\":\"resource not found\"}\n", string(w.Body.Bytes()))
	assert.Equal(t, []string{"application/json; charset=utf-8"}, w.HeaderMap["Content-Type"])

	w = httptest.NewRecorder()
	encodeError(bg, ErrInvalidArgument, w)
	assert.Equal(t, http.StatusBadRequest, w.Code)
	assert.Equal(t, "{\"error\":\"invalid argument\"}\n", string(w.Body.Bytes()))
	assert.Equal(t, []string{"application/json; charset=utf-8"}, w.HeaderMap["Content-Type"])
}

func Test_encodeStatusOnlyResponse(t *testing.T) {
	w := httptest.NewRecorder()

	err := encodeStatusOnlyResponse(bg, w, deleteAttrResponse{})
	assert.NoError(t, err)
	assert.Equal(t, http.StatusAccepted, w.Code)
	assert.Len(t, w.Body.Bytes(), 0)

	w = httptest.NewRecorder()
	err = encodeStatusOnlyResponse(bg, w, deleteAttrResponse{Err: errors.New("some-error")})
	assert.NoError(t, err)
	assert.Equal(t, http.StatusInternalServerError, w.Code)
}

func Test_encodeResponse(t *testing.T) {
	w := httptest.NewRecorder()
	err := encodeResponse(bg, w, createUserResponse{User: iam.User{ID: "urn:iam::user/10"}})
	assert.NoError(t, err)
	assert.Equal(t, http.StatusOK, w.Code)
	assert.NotEqual(t, 0, len(w.Body.Bytes()))
	assert.Equal(t, []string{"application/json; charset=utf-8"}, w.HeaderMap["Content-Type"])

	w = httptest.NewRecorder()
	err = encodeResponse(bg, w, createUserResponse{Err: errors.New("some-error")})
	assert.NoError(t, err)
	assert.Equal(t, http.StatusInternalServerError, w.Code)
	assert.NotEqual(t, 0, len(w.Body.Bytes()))
}
