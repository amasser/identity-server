package inmem

import (
	"context"
	"sync"

	"github.com/tierklinik-dobersberg/identity-server/iam"
	"github.com/tierklinik-dobersberg/identity-server/pkg/common"
)

type policyRepo struct {
	l sync.RWMutex
	m map[iam.PolicyURN]iam.Policy
}

func (r *policyRepo) Store(ctx context.Context, policy iam.Policy) error {
	r.l.Lock()
	defer r.l.Unlock()

	r.m[iam.PolicyURN(policy.ID)] = policy
	return nil
}

func (r *policyRepo) Delete(ctx context.Context, urn iam.PolicyURN) error {
	r.l.Lock()
	defer r.l.Unlock()

	if _, ok := r.m[urn]; !ok {
		return common.NewNotFoundError("policy")
	}

	delete(r.m, urn)

	return nil
}

func (r *policyRepo) Load(ctx context.Context, urn iam.PolicyURN) (iam.Policy, error) {
	r.l.RLock()
	defer r.l.RUnlock()

	p, ok := r.m[urn]
	if !ok {
		return iam.Policy{}, common.NewNotFoundError("policy")
	}

	return p, nil
}

func (r *policyRepo) Get(ctx context.Context) ([]iam.Policy, error) {
	r.l.RLock()
	defer r.l.RUnlock()

	policies := make([]iam.Policy, 0, len(r.m))
	for _, p := range r.m {
		policies = append(policies, p)
	}

	return policies, nil
}

// NewPolicyRepository returns a new in-memory policy repository.
func NewPolicyRepository() iam.PolicyRepository {
	return &policyRepo{
		m: make(map[iam.PolicyURN]iam.Policy),
	}
}
