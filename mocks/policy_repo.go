package mocks

import (
	"context"

	"github.com/stretchr/testify/mock"
	"github.com/tierklinik-dobersberg/identity-server/pkg/iam"
)

type PolicyRepository struct {
	mock.Mock
}

func (p *PolicyRepository) Store(ctx context.Context, policy iam.Policy) error {
	return p.Called(policy).Error(0)
}

func (p *PolicyRepository) Delete(ctx context.Context, urn iam.PolicyURN) error {
	return p.Called(urn).Error(0)
}

func (p *PolicyRepository) Load(ctx context.Context, urn iam.PolicyURN) (iam.Policy, error) {
	args := p.Called(urn)
	return args.Get(0).(iam.Policy), args.Error(1)
}

func (p *PolicyRepository) Get(ctx context.Context) ([]iam.Policy, error) {
	args := p.Called()
	return args.Get(0).([]iam.Policy), args.Error(1)
}
