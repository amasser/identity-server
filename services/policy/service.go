package policy

import (
	"context"
	"fmt"

	"github.com/tierklinik-dobersberg/identity-server/iam"
	"github.com/tierklinik-dobersberg/identity-server/pkg/mutex"
)

// Service implements polciy management functionallity.
type Service interface {
	// Create creates a new access policy under name.
	Create(ctx context.Context, name string, policy iam.Policy) (iam.PolicyURN, error)

	// Delete deletes a policy.
	Delete(ctx context.Context, urn iam.PolicyURN) error

	// Load loads the policy with the given URN.
	Load(ctx context.Context, urn iam.PolicyURN) (iam.Policy, error)

	// Update updates an existing policy.
	Update(ctx context.Context, urn iam.PolicyURN, p iam.Policy) error

	// List returns a list of all available policies.
	List(ctx context.Context) ([]iam.Policy, error)
}

type service struct {
	m    *mutex.Mutex
	repo iam.PolicyRepository
}

func (s *service) Create(ctx context.Context, name string, policy iam.Policy) (iam.PolicyURN, error) {
	if !s.m.TryLock(ctx) {
		return "", ctx.Err()
	}
	defer s.m.Unlock()

	policy.ID = iam.PolicyURN(fmt.Sprintf("urn:iam::policy/%s", name))

	if err := s.repo.Store(ctx, policy); err != nil {
		return "", nil
	}

	return policy.ID, nil
}

func (s *service) Delete(ctx context.Context, urn iam.PolicyURN) error {
	if !s.m.TryLock(ctx) {
		return ctx.Err()
	}
	defer s.m.Unlock()

	return s.repo.Delete(ctx, urn)
}

func (s *service) Load(ctx context.Context, urn iam.PolicyURN) (iam.Policy, error) {
	if !s.m.TryLock(ctx) {
		return iam.Policy{}, ctx.Err()
	}
	defer s.m.Unlock()

	return s.repo.Load(ctx, urn)
}

func (s *service) Update(ctx context.Context, urn iam.PolicyURN, p iam.Policy) error {
	if !s.m.TryLock(ctx) {
		return ctx.Err()
	}
	defer s.m.Unlock()

	p.ID = urn
	if err := s.repo.Store(ctx, p); err != nil {
		return err
	}

	return nil
}

func (s *service) List(ctx context.Context) ([]iam.Policy, error) {
	if !s.m.TryLock(ctx) {
		return nil, ctx.Err()
	}
	defer s.m.Unlock()

	return s.repo.Get(ctx)
}

// NewService returns a new policy management service.
func NewService(repo iam.PolicyRepository) Service {
	return &service{
		m:    mutex.New(),
		repo: repo,
	}
}
