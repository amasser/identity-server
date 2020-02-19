package user

import (
	"github.com/stretchr/testify/mock"
	"github.com/tierklinik-dobersberg/iam/v2/pkg/authn"
)

type authnMock struct {
	mock.Mock
}

func (a *authnMock) ImportAccount(username, password string, locked bool) (int, error) {
	args := a.Called(username, password, locked)
	return args.Int(0), args.Error(1)
}

func (a *authnMock) GetAccount(id int) (authn.Account, error) {
	args := a.Called(id)
	return args.Get(0).(authn.Account), args.Error(1)
}

func (a *authnMock) LockAccount(id int) error {
	return a.Called(id).Error(0)
}

func (a *authnMock) UnlockAccount(id int) error {
	return a.Called(id).Error(0)
}

func (a *authnMock) ArchiveAccount(id int) error {
	return a.Called(id).Error(0)
}
