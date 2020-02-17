package user

import (
	"context"
	"errors"
	"fmt"
	"os"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
	"github.com/tierklinik-dobersberg/iam/v2/iam"
)

var bg = context.Background()

func TestService_CreateUser(t *testing.T) {
	t.Run("Create_Successful", func(t *testing.T) {
		t.Parallel()

		r := &userRepoMock{}
		svc := NewService(r)

		expectedUser := iam.User{
			AccountID: 1,
			Username:  "admin",
			ID:        "urn:iam::user/1",
			Attributes: map[string]interface{}{
				"role":       "admin",
				"job":        "IT Administrator",
				"department": "IT",
			},
		}

		r.On("Load", iam.UserURN("urn:iam::user/1")).Once().Return(iam.User{}, os.ErrNotExist)
		r.On("Store", expectedUser).Return(nil)

		userUrn, err := svc.CreateUser(bg, 1, "admin", map[string]interface{}{
			"role":       "admin",
			"job":        "IT Administrator",
			"department": "IT",
		})

		assert.NoError(t, err)
		assert.Equal(t, iam.UserURN("urn:iam::user/1"), userUrn)
		r.AssertExpectations(t)
	})

	t.Run("Create_UserExists", func(t *testing.T) {
		t.Parallel()

		r := &userRepoMock{}
		svc := NewService(r)

		r.On("Load", iam.UserURN("urn:iam::user/2")).Once().Return(iam.User{}, nil)
		_, err := svc.CreateUser(bg, 2, "admin", nil)
		assert.Error(t, err)
		assert.True(t, os.IsExist(err))
		r.AssertExpectations(t)
	})

	t.Run("Create_UnknownError", func(t *testing.T) {
		t.Parallel()

		r := &userRepoMock{}
		svc := NewService(r)

		r.On("Load", iam.UserURN("urn:iam::user/3")).Once().Return(iam.User{}, errors.New("simulated"))
		_, err := svc.CreateUser(bg, 3, "admin", nil)
		assert.Error(t, err)
		assert.Equal(t, "simulated", err.Error())

		r.On("Load", iam.UserURN("urn:iam::user/3")).Once().Return(iam.User{}, os.ErrNotExist)
		r.On("Store", mock.Anything).Once().Return(errors.New("simulated 2"))
		_, err = svc.CreateUser(bg, 3, "admin", nil)
		assert.Error(t, err)
		assert.Equal(t, "simulated 2", err.Error())
	})
}

func TestService_LoadUser(t *testing.T) {

	t.Run("Load_InvalidArgs", func(t *testing.T) {
		t.Parallel()

		r := &userRepoMock{}
		svc := NewService(r)

		_, err := svc.LoadUser(bg, "")
		assert.Error(t, err)
		assert.Equal(t, ErrInvalidArgument, err)

		r.AssertExpectations(t)
	})

	t.Run("Load_Success", func(t *testing.T) {
		t.Parallel()

		r := &userRepoMock{}
		svc := NewService(r)

		expectedUser := expectedUser(10)

		r.On("Load", iam.UserURN("urn:iam::user/10")).Once().Return(expectedUser, nil)

		u, err := svc.LoadUser(bg, iam.UserURN("urn:iam::user/10"))
		assert.NoError(t, err)
		assert.Equal(t, expectedUser, u)
	})
}

func TestService_Users(t *testing.T) {
	t.Run("Users_Success", func(t *testing.T) {
		t.Parallel()

		r := &userRepoMock{}
		svc := NewService(r)

		expectedUser := expectedUser(10)

		r.On("Get").Return([]iam.User{
			expectedUser,
		}, nil)

		users, err := svc.Users(bg)
		assert.NoError(t, err)
		assert.Len(t, users, 1)
		assert.Equal(t, expectedUser, users[0])
	})

	t.Run("Users_Failure", func(t *testing.T) {
		t.Parallel()

		r := &userRepoMock{}
		svc := NewService(r)

		r.On("Get").Return(([]iam.User)(nil), errors.New("simulated"))
		_, err := svc.Users(bg)
		assert.Equal(t, "simulated", err.Error())
	})
}

func TestService_UpdateAttrs(t *testing.T) {
	t.Run("UpdateAttr_InvalidArg", func(t *testing.T) {
		t.Parallel()

		r := &userRepoMock{}
		svc := NewService(r)
		err := svc.UpdateAttrs(bg, "", map[string]interface{}{"some": "key"})
		assert.Error(t, err)
		assert.Equal(t, ErrInvalidArgument, err)
	})

	t.Run("UpdateAttr_NotExist", func(t *testing.T) {
		t.Parallel()

		r := &userRepoMock{}
		svc := NewService(r)
		urn := iam.UserURN("urn:iam::user/11")

		r.On("Load", urn).Once().Return(iam.User{}, os.ErrNotExist)
		err := svc.UpdateAttrs(bg, urn, map[string]interface{}{"new": "value"})
		assert.Error(t, err)
		assert.True(t, os.IsNotExist(err))
		r.AssertExpectations(t)
	})

	t.Run("UpdateAttr_Success", func(t *testing.T) {
		t.Parallel()

		r := &userRepoMock{}
		svc := NewService(r)
		urn := iam.UserURN("urn:iam::user/11")

		inputUser := expectedUser(10)
		expectedAttrs := map[string]interface{}{
			"new": "value",
		}
		expectedUser := inputUser
		expectedUser.Attributes = expectedAttrs

		r.On("Load", urn).Once().Return(inputUser, nil)
		r.On("Store", expectedUser).Once().Return(nil)

		err := svc.UpdateAttrs(bg, urn, expectedAttrs)
		assert.NoError(t, err)
		r.AssertExpectations(t)
	})
}

func TestService_SetAttr(t *testing.T) {
	t.Run("SetAttr_InvalidArg", func(t *testing.T) {
		t.Parallel()

		r := &userRepoMock{}
		svc := NewService(r)
		err := svc.SetAttr(bg, "", "some", "key")
		assert.Error(t, err)
		assert.Equal(t, ErrInvalidArgument, err)
	})

	t.Run("SetAttr_NotExist", func(t *testing.T) {
		t.Parallel()

		r := &userRepoMock{}
		svc := NewService(r)
		urn := iam.UserURN("urn:iam::user/11")

		r.On("Load", urn).Once().Return(iam.User{}, os.ErrNotExist)
		err := svc.SetAttr(bg, urn, "key", "value")
		assert.Error(t, err)
		assert.True(t, os.IsNotExist(err))
		r.AssertExpectations(t)
	})

	t.Run("SetAttr_Success", func(t *testing.T) {
		t.Parallel()

		r := &userRepoMock{}
		svc := NewService(r)
		urn := iam.UserURN("urn:iam::user/11")

		inputUser := expectedUser(10)
		expectedAttrs := map[string]interface{}{
			"new": "value",
			"job": "Developer",
		}
		expectedUser := inputUser
		expectedUser.Attributes = expectedAttrs

		r.On("Load", urn).Once().Return(inputUser, nil)
		r.On("Store", expectedUser).Once().Return(nil)

		err := svc.SetAttr(bg, urn, "new", "value")
		assert.NoError(t, err)
		r.AssertExpectations(t)
	})

	t.Run("SetAttr_Success_nil_Attributes", func(t *testing.T) {
		t.Parallel()

		r := &userRepoMock{}
		svc := NewService(r)
		urn := iam.UserURN("urn:iam::user/11")

		inputUser := expectedUser(10)
		inputUser.Attributes = nil
		expectedAttrs := map[string]interface{}{
			"new": "value",
		}
		expectedUser := inputUser
		expectedUser.Attributes = expectedAttrs

		r.On("Load", urn).Once().Return(inputUser, nil)
		r.On("Store", expectedUser).Once().Return(nil)

		err := svc.SetAttr(bg, urn, "new", "value")
		assert.NoError(t, err)
		r.AssertExpectations(t)
	})
}

func TestService_DeleteAttr(t *testing.T) {
	t.Run("DeleteAttr_InvalidArg", func(t *testing.T) {
		t.Parallel()

		r := &userRepoMock{}
		svc := NewService(r)
		err := svc.DeleteAttr(bg, "", "job")
		assert.Error(t, err)
		assert.Equal(t, ErrInvalidArgument, err)
	})

	t.Run("DeleteAttr_NotExist", func(t *testing.T) {
		t.Parallel()

		r := &userRepoMock{}
		svc := NewService(r)
		urn := iam.UserURN("urn:iam::user/11")

		r.On("Load", urn).Once().Return(iam.User{}, os.ErrNotExist)
		err := svc.DeleteAttr(bg, urn, "job")
		assert.Error(t, err)
		assert.True(t, os.IsNotExist(err))
		r.AssertExpectations(t)
	})

	t.Run("DeleteAttr_Success", func(t *testing.T) {
		t.Parallel()

		r := &userRepoMock{}
		svc := NewService(r)
		urn := iam.UserURN("urn:iam::user/11")

		inputUser := expectedUser(10)
		expectedAttrs := map[string]interface{}{}
		expectedUser := inputUser
		expectedUser.Attributes = expectedAttrs

		r.On("Load", urn).Once().Return(inputUser, nil)
		r.On("Store", expectedUser).Once().Return(nil)

		err := svc.DeleteAttr(bg, urn, "job")
		assert.NoError(t, err)
		r.AssertExpectations(t)
	})

	t.Run("DeleteAttr_Success_niL_Attributes", func(t *testing.T) {
		t.Parallel()

		r := &userRepoMock{}
		svc := NewService(r)
		urn := iam.UserURN("urn:iam::user/11")

		inputUser := expectedUser(10)
		inputUser.Attributes = nil
		expectedAttrs := map[string]interface{}{}
		expectedUser := inputUser
		expectedUser.Attributes = expectedAttrs

		r.On("Load", urn).Once().Return(inputUser, nil)

		err := svc.DeleteAttr(bg, urn, "job")
		assert.NoError(t, err)
		r.AssertExpectations(t)
	})

}

type userRepoMock struct {
	mock.Mock
}

func (rm *userRepoMock) Store(_ context.Context, user iam.User) error {
	return rm.Called(user).Error(0)
}

func (rm *userRepoMock) Delete(_ context.Context, urn iam.UserURN) error {
	return rm.Called(urn).Error(0)
}

func (rm *userRepoMock) Load(_ context.Context, urn iam.UserURN) (iam.User, error) {
	args := rm.Called(urn)
	return args.Get(0).(iam.User), args.Error(1)
}

func (rm *userRepoMock) Get(_ context.Context) ([]iam.User, error) {
	args := rm.Called()
	return args.Get(0).([]iam.User), args.Error(1)
}

func expectedUser(id int) iam.User {
	return iam.User{
		AccountID: id,
		ID:        iam.UserURN(fmt.Sprintf("urn:iam::user/%d", id)),
		Attributes: map[string]interface{}{
			"job": "Developer",
		},
	}

}
