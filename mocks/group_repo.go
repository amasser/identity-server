package mocks

import (
	"context"

	"github.com/stretchr/testify/mock"
	"github.com/tierklinik-dobersberg/identity-server/pkg/iam"
)

type GroupRepository struct {
	mock.Mock
}

func (m *GroupRepository) Store(_ context.Context, group iam.Group) error {
	return m.Called(group).Error(0)
}

func (m *GroupRepository) Delete(_ context.Context, urn iam.GroupURN) error {
	return m.Called(urn).Error(0)
}

func (m *GroupRepository) Load(_ context.Context, urn iam.GroupURN) (iam.Group, error) {
	args := m.Called(urn)
	return args.Get(0).(iam.Group), args.Error(1)
}

func (m *GroupRepository) Get(_ context.Context) ([]iam.Group, error) {
	args := m.Called()
	return args.Get(0).([]iam.Group), args.Error(1)
}

func NewGroupRepository() *GroupRepository {
	return &GroupRepository{}
}
