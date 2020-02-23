package mocks

import (
	"context"

	"github.com/stretchr/testify/mock"
	"github.com/tierklinik-dobersberg/identity-server/iam"
)

type MembershipRepository struct {
	mock.Mock
}

func (m *MembershipRepository) AddMember(ctx context.Context, user iam.UserURN, group iam.GroupURN) error {
	return m.Called(user, group).Error(0)
}

func (m *MembershipRepository) DeleteMember(ctx context.Context, user iam.UserURN, group iam.GroupURN) error {
	return m.Called(user, group).Error(0)
}

func (m *MembershipRepository) Memberships(ctx context.Context, user iam.UserURN) ([]iam.GroupURN, error) {
	args := m.Called(user)
	return args.Get(0).([]iam.GroupURN), args.Error(1)
}

func (m *MembershipRepository) Members(ctx context.Context, group iam.GroupURN) ([]iam.UserURN, error) {
	args := m.Called(group)
	return args.Get(0).([]iam.UserURN), args.Error(1)
}

func NewMembershipRepository() *MembershipRepository {
	return &MembershipRepository{}
}
