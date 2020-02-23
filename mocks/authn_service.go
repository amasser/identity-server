package mocks

import (
	"github.com/stretchr/testify/mock"
	"github.com/tierklinik-dobersberg/identity-server/pkg/authn"
)

type AuthnService struct {
	mock.Mock
}

func (a *AuthnService) ImportAccount(username, password string, locked bool) (int, error) {
	args := a.Called(username, password, locked)
	return args.Int(0), args.Error(1)
}

func (a *AuthnService) GetAccount(id int) (authn.Account, error) {
	args := a.Called(id)
	return args.Get(0).(authn.Account), args.Error(1)
}

func (a *AuthnService) LockAccount(id int) error {
	return a.Called(id).Error(0)
}

func (a *AuthnService) UnlockAccount(id int) error {
	return a.Called(id).Error(0)
}

func (a *AuthnService) ArchiveAccount(id int) error {
	return a.Called(id).Error(0)
}

func NewAuthnService() *AuthnService {
	return &AuthnService{}
}
