package group

import (
	"context"
	"errors"
	"testing"

	"github.com/go-kit/kit/log"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
	"github.com/tierklinik-dobersberg/identity-server/iam"
	"github.com/tierklinik-dobersberg/identity-server/mocks"
	"github.com/tierklinik-dobersberg/identity-server/pkg/common"
	"github.com/tierklinik-dobersberg/identity-server/services/user"
)

var testCtx = context.Background()

func TestService_Create(t *testing.T) {
	t.Run("Invalid name", func(t *testing.T) {
		t.Parallel()
		s := setupTestBed()

		urn, err := s.Create(testCtx, "", "some-comment")
		assert.Equal(t, iam.GroupURN(""), urn)
		assert.Error(t, err)
		s.AssertExpectations(t)
	})

	t.Run("Store failed", func(t *testing.T) {
		t.Parallel()
		s := setupTestBed()
		expectedGroup := iam.Group{
			ID:      "urn:iam::group/devs",
			Name:    "devs",
			Comment: "some-comment",
		}

		s.groups.On("Store", expectedGroup).Once().Return(errors.New("simulated error"))

		urn, err := s.Create(testCtx, "devs", "some-comment")
		assert.Equal(t, iam.GroupURN(""), urn)
		assert.Error(t, err)
		s.AssertExpectations(t)
	})

	t.Run("Success", func(t *testing.T) {
		t.Parallel()
		s := setupTestBed()
		expectedGroup := iam.Group{
			ID:      "urn:iam::group/devs",
			Name:    "devs",
			Comment: "some-comment",
		}

		s.groups.On("Store", expectedGroup).Once().Return(nil)

		urn, err := s.Create(testCtx, "devs", "some-comment")
		assert.Equal(t, iam.GroupURN("urn:iam::group/devs"), urn)
		assert.NoError(t, err)
		s.AssertExpectations(t)
	})
}

func TestService_Delete(t *testing.T) {
	t.Run("Invalid group name", func(t *testing.T) {
		t.Parallel()
		s := setupTestBed()

		assert.Equal(t, ErrInvalidParameter, s.Delete(testCtx, ""))
	})

	t.Run("Members() failed", func(t *testing.T) {
		t.Parallel()
		s := setupTestBed()

		s.members.On("Members", iam.GroupURN("urn:iam::group/devs")).Once().Return([]iam.UserURN(nil), errors.New("simulated"))
		err := s.Delete(testCtx, "urn:iam::group/devs")
		assert.Error(t, err)
		assert.Equal(t, "simulated", err.Error())
	})

	t.Run("Delete all members", func(t *testing.T) {
		t.Parallel()

		t.Run("DeleteMember fails", func(t *testing.T) {
			t.Parallel()
			s := setupTestBed()

			s.members.On("Members", iam.GroupURN("urn:iam::group/devs")).Once().Return([]iam.UserURN{
				"urn:iam::user/10",
				"urn:iam::user/22",
			}, nil)

			s.members.On("DeleteMember", iam.UserURN("urn:iam::user/10"), iam.GroupURN("urn:iam::group/devs")).Once().Return(nil)
			s.members.On("DeleteMember", iam.UserURN("urn:iam::user/22"), iam.GroupURN("urn:iam::group/devs")).Once().Return(errors.New("simulated"))

			err := s.Delete(testCtx, "urn:iam::group/devs")
			assert.Error(t, err)
			assert.Equal(t, "simulated", err.Error())
			s.AssertExpectations(t)
		})

		t.Run("Success", func(t *testing.T) {
			t.Parallel()
			s := setupTestBed()

			s.members.On("Members", iam.GroupURN("urn:iam::group/devs")).Once().Return([]iam.UserURN{
				"urn:iam::user/10",
				"urn:iam::user/22",
			}, nil)

			s.members.On("DeleteMember", iam.UserURN("urn:iam::user/10"), iam.GroupURN("urn:iam::group/devs")).Once().Return(nil)
			s.members.On("DeleteMember", iam.UserURN("urn:iam::user/22"), iam.GroupURN("urn:iam::group/devs")).Once().Return(nil)
			s.groups.On("Delete", iam.GroupURN("urn:iam::group/devs")).Once().Return(nil)

			err := s.Delete(testCtx, "urn:iam::group/devs")
			assert.NoError(t, err)
			s.AssertExpectations(t)
		})
	})
}

func TestService_Load(t *testing.T) {
	t.Run("invalid name", func(t *testing.T) {
		t.Parallel()
		s := setupTestBed()

		_, err := s.Load(testCtx, "")
		assert.Equal(t, ErrInvalidParameter, err)
		s.AssertExpectations(t)
	})

	t.Run("Load failes", func(t *testing.T) {
		t.Parallel()
		s := setupTestBed()

		s.groups.On("Load", iam.GroupURN("urn:iam::group/devs")).Return(iam.Group{}, errors.New("simulated"))

		_, err := s.Load(testCtx, "urn:iam::group/devs")
		assert.Error(t, err)
		s.AssertExpectations(t)
	})

	t.Run("Success", func(t *testing.T) {
		t.Parallel()
		s := setupTestBed()

		expectedGroup := iam.Group{
			ID:      "urn:iam::group/devs",
			Name:    "devs",
			Comment: "A useful comment",
		}

		s.groups.On("Load", iam.GroupURN("urn:iam::group/devs")).Return(expectedGroup, nil)

		g, err := s.Load(testCtx, "urn:iam::group/devs")
		assert.NoError(t, err)
		assert.Equal(t, expectedGroup, g)
		s.AssertExpectations(t)
	})
}

func TestService_UpdateComment(t *testing.T) {
	t.Run("Invalid name", func(t *testing.T) {
		t.Parallel()
		s := setupTestBed()

		err := s.UpdateComment(testCtx, "", "new-comment")
		assert.Error(t, err)
		assert.Equal(t, ErrInvalidParameter, err)
		s.AssertExpectations(t)
	})

	t.Run("Load fails", func(t *testing.T) {
		t.Parallel()
		s := setupTestBed()

		s.groups.On("Load", iam.GroupURN("urn:iam::group/devs")).Once().Return(iam.Group{}, errors.New("simulated"))

		err := s.UpdateComment(testCtx, "urn:iam::group/devs", "new comment")
		assert.Error(t, err)
		s.AssertExpectations(t)
	})

	t.Run("Store fails", func(t *testing.T) {
		t.Parallel()
		s := setupTestBed()

		existingGroup := iam.Group{
			ID:      "urn:iam::group/devs",
			Name:    "devs",
			Comment: "old comment",
		}
		expectedGroup := existingGroup
		expectedGroup.Comment = "new comment"

		s.groups.On("Load", iam.GroupURN("urn:iam::group/devs")).Once().Return(existingGroup, nil)
		s.groups.On("Store", expectedGroup).Once().Return(errors.New("simulated"))

		err := s.UpdateComment(testCtx, "urn:iam::group/devs", "new comment")
		assert.Error(t, err)
		s.AssertExpectations(t)
	})

	t.Run("Success", func(t *testing.T) {
		t.Parallel()
		s := setupTestBed()

		existingGroup := iam.Group{
			ID:      "urn:iam::group/devs",
			Name:    "devs",
			Comment: "old comment",
		}
		expectedGroup := existingGroup
		expectedGroup.Comment = "new comment"

		s.groups.On("Load", iam.GroupURN("urn:iam::group/devs")).Once().Return(existingGroup, nil)
		s.groups.On("Store", expectedGroup).Once().Return(nil)

		err := s.UpdateComment(testCtx, "urn:iam::group/devs", "new comment")
		assert.NoError(t, err)
		s.AssertExpectations(t)
	})
}

func TestService_AddMember(t *testing.T) {
	t.Run("Load failes", func(t *testing.T) {
		t.Parallel()
		s := setupTestBed()

		s.groups.On("Load", iam.GroupURN("urn:iam::group/devs")).Once().Return(iam.Group{}, errors.New("simulated"))
		err := s.AddMember(testCtx, "urn:iam::group/devs", "urn:iam::user/10")
		assert.Error(t, err)
		s.AssertExpectations(t)
	})

	t.Run("User does not exist", func(t *testing.T) {
		t.Parallel()
		s := setupTestBed()

		s.groups.On("Load", iam.GroupURN("urn:iam::group/devs")).Once().Return(iam.Group{}, nil)
		s.users.On("LoadUser", iam.UserURN("urn:iam::user/10")).Once().Return(iam.User{}, &common.NotFoundError{})

		err := s.AddMember(testCtx, "urn:iam::group/devs", "urn:iam::user/10")
		assert.True(t, common.IsNotFound(err))
		s.AssertExpectations(t)
	})

	t.Run("AddMember fails", func(t *testing.T) {
		t.Parallel()
		s := setupTestBed()

		s.groups.On("Load", iam.GroupURN("urn:iam::group/devs")).Once().Return(iam.Group{}, nil)
		s.users.On("LoadUser", iam.UserURN("urn:iam::user/10")).Once().Return(iam.User{}, nil)
		s.members.On("AddMember", iam.UserURN("urn:iam::user/10"), iam.GroupURN("urn:iam::group/devs")).Return(errors.New("simulated"))

		err := s.AddMember(testCtx, "urn:iam::group/devs", "urn:iam::user/10")
		assert.Error(t, err)
		s.AssertExpectations(t)
	})

	t.Run("Success", func(t *testing.T) {
		t.Parallel()
		s := setupTestBed()

		s.groups.On("Load", iam.GroupURN("urn:iam::group/devs")).Once().Return(iam.Group{}, nil)
		s.users.On("LoadUser", iam.UserURN("urn:iam::user/10")).Once().Return(iam.User{}, nil)
		s.members.On("AddMember", iam.UserURN("urn:iam::user/10"), iam.GroupURN("urn:iam::group/devs")).Return(nil)

		err := s.AddMember(testCtx, "urn:iam::group/devs", "urn:iam::user/10")
		s.AssertExpectations(t)
		assert.NoError(t, err)
	})
}

func TestService_DeleteMember(t *testing.T) {
	t.Run("Load failes", func(t *testing.T) {
		t.Parallel()
		s := setupTestBed()

		s.groups.On("Load", iam.GroupURN("urn:iam::group/devs")).Return(iam.Group{}, common.NewNotFoundError(""))
		err := s.DeleteMember(testCtx, "urn:iam::group/devs", "urn:iam::user/10")
		assert.True(t, common.IsNotFound(err))
		s.AssertExpectations(t)
	})

	t.Run("DeleteMember failes", func(t *testing.T) {
		t.Parallel()
		s := setupTestBed()

		s.groups.On("Load", iam.GroupURN("urn:iam::group/devs")).Return(iam.Group{}, nil)
		s.members.On("DeleteMember", iam.UserURN("urn:iam::user/10"), iam.GroupURN("urn:iam::group/devs")).Return(errors.New("simulated"))

		err := s.DeleteMember(testCtx, "urn:iam::group/devs", "urn:iam::user/10")
		assert.Error(t, err)
		s.AssertExpectations(t)
	})

	t.Run("Success", func(t *testing.T) {
		t.Parallel()
		s := setupTestBed()

		s.groups.On("Load", iam.GroupURN("urn:iam::group/devs")).Return(iam.Group{}, nil)
		s.members.On("DeleteMember", iam.UserURN("urn:iam::user/10"), iam.GroupURN("urn:iam::group/devs")).Return(nil)

		err := s.DeleteMember(testCtx, "urn:iam::group/devs", "urn:iam::user/10")
		assert.NoError(t, err)
		s.AssertExpectations(t)
	})
}

func TestService_onUserDeleted(t *testing.T) {
	t.Run("Remove a deleted user form all groups", func(t *testing.T) {
		t.Parallel()
		s := setupTestBed()

		s.members.On("Memberships", iam.UserURN("urn:iam::user/10")).Once().Return([]iam.GroupURN{"urn:iam::group/devs"}, nil)
		s.members.On("DeleteMember", iam.UserURN("urn:iam::user/10"), iam.GroupURN("urn:iam::group/devs")).Once().Return(nil)

		s.onDelete("urn:iam::user/10")
		s.AssertExpectations(t)
	})

	t.Run("Memberships failed to load", func(t *testing.T) {
		t.Parallel()
		s := setupTestBed()

		s.members.On("Memberships", iam.UserURN("urn:iam::user/10")).Once().Return([]iam.GroupURN(nil), errors.New("simulated"))

		s.onDelete("urn:iam::user/10")
		s.AssertExpectations(t)
	})

	t.Run("Remove a deleted user form all groups", func(t *testing.T) {
		t.Parallel()
		s := setupTestBed()

		s.members.On("Memberships", iam.UserURN("urn:iam::user/10")).Once().Return([]iam.GroupURN{
			"urn:iam::group/devs",
			"urn:iam::group/admins",
		}, nil)
		s.members.On("DeleteMember", iam.UserURN("urn:iam::user/10"), iam.GroupURN("urn:iam::group/devs")).Once().Return(errors.New("simulated"))
		s.members.On("DeleteMember", iam.UserURN("urn:iam::user/10"), iam.GroupURN("urn:iam::group/admins")).Once().Return(nil)

		s.onDelete("urn:iam::user/10")
		s.AssertExpectations(t)
	})
}

type testBed struct {
	Service

	onDelete user.OnDeleteFunc
	groups   *mocks.GroupRepository
	members  *mocks.MembershipRepository
	users    *userServiceMock
}

func setupTestBed() *testBed {
	var fn user.OnDeleteFunc
	us := &userServiceMock{}
	gr := mocks.NewGroupRepository()
	mr := mocks.NewMembershipRepository()

	us.On("OnDelete", mock.Anything).Once().Run(func(args mock.Arguments) {
		fn = args[0].(user.OnDeleteFunc)
	})

	s := NewService(us, gr, mr, log.NewNopLogger())
	s = NewLoggingService(s, log.NewNopLogger())

	return &testBed{
		Service:  s,
		onDelete: fn,
		groups:   gr,
		members:  mr,
		users:    us,
	}
}

func (bed *testBed) AssertExpectations(t *testing.T) {
	bed.groups.AssertExpectations(t)
	bed.members.AssertExpectations(t)
}

type userServiceMock struct {
	mock.Mock
}

func (s *userServiceMock) CreateUser(_ context.Context, username, password string, attrs map[string]interface{}) (iam.UserURN, error) {
	args := s.Called(username, password, attrs)
	return args.Get(0).(iam.UserURN), args.Error(1)
}

func (s *userServiceMock) LoadUser(_ context.Context, urn iam.UserURN) (iam.User, error) {
	args := s.Called(urn)
	return args.Get(0).(iam.User), args.Error(1)
}

func (s *userServiceMock) DeleteUser(_ context.Context, urn iam.UserURN) error {
	return s.Called(urn).Error(0)
}

func (s *userServiceMock) LockUser(_ context.Context, urn iam.UserURN, locked bool) error {
	return s.Called(urn, locked).Error(0)
}

func (s *userServiceMock) Users(_ context.Context) ([]iam.User, error) {
	args := s.Called()
	return args.Get(0).([]iam.User), args.Error(1)
}

func (s *userServiceMock) UpdateAttrs(_ context.Context, id iam.UserURN, attrs map[string]interface{}) error {
	return s.Called(id, attrs).Error(0)
}

func (s *userServiceMock) SetAttr(_ context.Context, urn iam.UserURN, key string, value interface{}) error {
	return s.Called(urn, key, value).Error(0)
}

func (s *userServiceMock) DeleteAttr(_ context.Context, urn iam.UserURN, key string) error {
	return s.Called(urn, key).Error(0)
}

func (s *userServiceMock) OnDelete(_ context.Context, fn user.OnDeleteFunc) {
	s.Called(fn)
}
