package user

import (
	"context"
	"errors"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
	"github.com/tierklinik-dobersberg/identity-server/pkg/iam"
)

func Test_userEndpoint(t *testing.T) {
	t.Run("Successful", func(t *testing.T) {
		s := &serviceMock{}
		ep := makeCreateUserEndpoint(s)
		attrs := map[string]interface{}{"job": "developer"}

		expectedUser := iam.User{
			AccountID:  10,
			Username:   "admin",
			ID:         "urn:iam::user/10",
			Attributes: attrs,
		}

		req := createUserRequest{
			Username:   "admin",
			Attributes: attrs,
			Password:   "password",
		}

		s.On("CreateUser", "admin", "password", attrs).Once().Return(
			iam.UserURN("urn:iam::user/10"),
			nil,
		)
		s.On("LoadUser", iam.UserURN("urn:iam::user/10")).Once().Return(
			expectedUser,
			nil,
		)

		res, err := ep(bg, req)
		assert.NoError(t, err)
		assert.Equal(t, expectedUser, res.(createUserResponse).User)
		assert.NoError(t, res.(createUserResponse).Err)
	})

	t.Run("CreateUser failed", func(t *testing.T) {
		s := &serviceMock{}
		ep := makeCreateUserEndpoint(s)
		req := createUserRequest{
			Username: "admin",
			Password: "password",
		}

		s.On("CreateUser", "admin", "password", map[string]interface{}(nil)).Once().Return(
			iam.UserURN(""),
			errors.New("some-error"),
		)

		res, err := ep(bg, req)
		assert.NoError(t, err)
		assert.Error(t, res.(createUserResponse).Err)
	})

	t.Run("LoadUser_failed", func(t *testing.T) {
		s := &serviceMock{}
		ep := makeCreateUserEndpoint(s)

		req := createUserRequest{
			Username: "admin",
			Password: "password",
		}

		s.On("CreateUser", "admin", "password", map[string]interface{}(nil)).Once().Return(
			iam.UserURN("urn:iam::user/10"),
			nil,
		)
		s.On("LoadUser", iam.UserURN("urn:iam::user/10")).Once().Return(
			iam.User{},
			errors.New("some-error"),
		)

		res, err := ep(bg, req)
		assert.NoError(t, err)
		assert.Equal(t, iam.User{}, res.(createUserResponse).User)
		assert.Error(t, res.(createUserResponse).Err)
	})
}

func Test_loadUserEndpoint(t *testing.T) {
	s := &serviceMock{}
	ep := makeLoadUserEndpoint(s)

	t.Run("LoadUser_Success", func(t *testing.T) {
		expectedUser := iam.User{
			Username: "admin",
		}
		s.On("LoadUser", iam.UserURN("urn:iam::user/10")).Once().Return(
			expectedUser,
			nil,
		)

		res, err := ep(bg, loadUserRequest{URN: "urn:iam::user/10"})
		assert.NoError(t, err)
		assert.Equal(t, &expectedUser, res.(loadUserResponse).User)
		assert.NoError(t, res.(loadUserResponse).Err)
	})

	t.Run("LoadUser_Error", func(t *testing.T) {
		s.On("LoadUser", iam.UserURN("urn:iam::user/10")).Once().Return(
			iam.User{},
			errors.New("some-error"),
		)

		res, err := ep(bg, loadUserRequest{URN: "urn:iam::user/10"})
		assert.NoError(t, err)
		assert.Nil(t, res.(loadUserResponse).User)
		assert.Error(t, res.(loadUserResponse).Err)
	})
}

func Test_deleteUserEndpoint(t *testing.T) {
	s := &serviceMock{}
	ep := makeDeleteUserEndpoint(s)

	t.Run("DeleteUser_Success", func(t *testing.T) {
		s.On("DeleteUser", iam.UserURN("urn:iam::user/10")).Once().Return(nil)
		res, err := ep(bg, deleteUserRequest{URN: "urn:iam::user/10"})
		assert.NoError(t, err)
		assert.NoError(t, res.(deleteUserResponse).Err)
	})

	t.Run("DeleteUser_Failure", func(t *testing.T) {
		s.On("DeleteUser", iam.UserURN("urn:iam::user/10")).Once().Return(errors.New("simulated"))
		res, err := ep(bg, deleteUserRequest{URN: "urn:iam::user/10"})
		assert.NoError(t, err)
		assert.Error(t, res.(deleteUserResponse).Err)
		assert.Equal(t, "simulated", res.(deleteUserResponse).Err.Error())
	})
}

func Test_lockUserEndpoint(t *testing.T) {
	s := &serviceMock{}
	ep := makeLockUserEndpoint(s)

	t.Run("DeleteUser_Success", func(t *testing.T) {
		s.On("LockUser", iam.UserURN("urn:iam::user/10"), true).Once().Return(nil)
		res, err := ep(bg, lockUserRequest{URN: "urn:iam::user/10", Locked: true})
		assert.NoError(t, err)
		assert.NoError(t, res.(lockUserResponse).Err)
	})

	t.Run("DeleteUser_Failure", func(t *testing.T) {
		s.On("LockUser", iam.UserURN("urn:iam::user/10"), false).Once().Return(errors.New("simulated"))
		res, err := ep(bg, lockUserRequest{URN: "urn:iam::user/10", Locked: false})
		assert.NoError(t, err)
		assert.Error(t, res.(lockUserResponse).Err)
		assert.Equal(t, "simulated", res.(lockUserResponse).Err.Error())
	})
}

func Test_listUsersEndpoint(t *testing.T) {
	s := &serviceMock{}
	ep := makeListUsersEndpoint(s)
	expectedUser := iam.User{
		Username: "admin",
	}

	s.On("Users").Once().Return([]iam.User(nil), errors.New("some-error"))
	res, err := ep(bg, listUsersRequest{})
	assert.NoError(t, err)
	assert.NotNil(t, res)
	assert.Nil(t, res.(listUsersResponse).Users)
	assert.Error(t, res.(listUsersResponse).Err)

	s.On("Users").Once().Return([]iam.User{expectedUser}, nil)
	res, err = ep(bg, listUsersRequest{})
	assert.NoError(t, err)
	assert.NotNil(t, res)
	assert.NotNil(t, res.(listUsersResponse).Users)
	assert.Equal(t, []iam.User{expectedUser}, res.(listUsersResponse).Users)
	assert.NoError(t, res.(listUsersResponse).Err)
}

func Test_updateAttributesEndpoint(t *testing.T) {
	s := &serviceMock{}
	ep := makeUpdateAttrsEndpoint(s)
	attr := map[string]interface{}{
		"job": "developer",
	}

	s.On("UpdateAttrs", iam.UserURN("urn:iam::user/10"), attr).Once().Return(nil)

	res, err := ep(bg, updateAttrsRequest{URN: "urn:iam::user/10", Attributes: attr})
	assert.NoError(t, err)
	assert.NoError(t, res.(updateAttrsResponse).Err)

	s.On("UpdateAttrs", iam.UserURN("urn:iam::user/10"), attr).Once().Return(errors.New("some-error"))

	res, err = ep(bg, updateAttrsRequest{URN: "urn:iam::user/10", Attributes: attr})
	assert.NoError(t, err)
	assert.Error(t, res.(updateAttrsResponse).Err)
}

func Test_setAttributeEndpoint(t *testing.T) {
	s := &serviceMock{}
	ep := makeSetAttrEndpoint(s)

	s.On("SetAttr", iam.UserURN("urn:iam::user/10"), "key", "value").Once().Return(nil)

	res, err := ep(bg, setAttrRequest{URN: "urn:iam::user/10", Key: "key", Value: "value"})
	assert.NoError(t, err)
	assert.NoError(t, res.(setAttrResponse).Err)

	s.On("SetAttr", iam.UserURN("urn:iam::user/10"), "key", "value").Once().Return(errors.New("some-error"))

	res, err = ep(bg, setAttrRequest{URN: "urn:iam::user/10", Key: "key", Value: "value"})
	assert.NoError(t, err)
	assert.Error(t, res.(setAttrResponse).Err)
}

func Test_deleteAttributeEndpoint(t *testing.T) {
	s := &serviceMock{}
	ep := makeDeleteAttrRequest(s)

	s.On("DeleteAttr", iam.UserURN("urn:iam::user/10"), "key").Once().Return(nil)

	res, err := ep(bg, deleteAttrRequest{URN: "urn:iam::user/10", Key: "key"})
	assert.NoError(t, err)
	assert.NoError(t, res.(deleteAttrResponse).Err)

	s.On("DeleteAttr", iam.UserURN("urn:iam::user/10"), "key").Once().Return(errors.New("some-error"))

	res, err = ep(bg, deleteAttrRequest{URN: "urn:iam::user/10", Key: "key"})
	assert.NoError(t, err)
	assert.Error(t, res.(deleteAttrResponse).Err)
}

type serviceMock struct {
	mock.Mock
}

func (s *serviceMock) CreateUser(_ context.Context, username, password string, attrs map[string]interface{}) (iam.UserURN, error) {
	args := s.Called(username, password, attrs)
	return args.Get(0).(iam.UserURN), args.Error(1)
}

func (s *serviceMock) LoadUser(_ context.Context, urn iam.UserURN) (iam.User, error) {
	args := s.Called(urn)
	return args.Get(0).(iam.User), args.Error(1)
}

func (s *serviceMock) DeleteUser(_ context.Context, urn iam.UserURN) error {
	return s.Called(urn).Error(0)
}

func (s *serviceMock) LockUser(_ context.Context, urn iam.UserURN, locked bool) error {
	return s.Called(urn, locked).Error(0)
}

func (s *serviceMock) Users(_ context.Context) ([]iam.User, error) {
	args := s.Called()
	return args.Get(0).([]iam.User), args.Error(1)
}

func (s *serviceMock) UpdateAttrs(_ context.Context, id iam.UserURN, attrs map[string]interface{}) error {
	return s.Called(id, attrs).Error(0)
}

func (s *serviceMock) SetAttr(_ context.Context, urn iam.UserURN, key string, value interface{}) error {
	return s.Called(urn, key, value).Error(0)
}

func (s *serviceMock) DeleteAttr(_ context.Context, urn iam.UserURN, key string) error {
	return s.Called(urn, key).Error(0)
}

func (s *serviceMock) OnDelete(_ context.Context, fn OnDeleteFunc) {
	s.Called(fn)
}
