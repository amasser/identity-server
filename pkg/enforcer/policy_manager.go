package enforcer

import (
	"context"

	"github.com/ory/ladon"
	"github.com/tierklinik-dobersberg/identity-server/pkg/common"
	"github.com/tierklinik-dobersberg/identity-server/pkg/iam"
)

// PolicyManager implements the ladon.Manager interface.
type PolicyManager struct {
	repo iam.PolicyRepository
}

// NewPolicyManager returns a new policy manager that can be used
// as the ladon.Manger in NewLadonEnforcer()
func NewPolicyManager(repo iam.PolicyRepository) *PolicyManager {
	return &PolicyManager{
		repo: repo,
	}
}

// Create implements ladon.Manager but does nothing.
func (*PolicyManager) Create(ladon.Policy) error {
	return common.ErrNotImplemented
}

// Update implements ladon.Manager but does nothing.
func (*PolicyManager) Update(ladon.Policy) error {
	return common.ErrNotImplemented
}

// Delete implements ladon.Manager but does nothing.
func (*PolicyManager) Delete(string) error {
	return common.ErrNotImplemented
}

// Get implements ladon.Manager and returns the policy associated
// with id.
func (p *PolicyManager) Get(id string) (ladon.Policy, error) {
	policy, err := p.repo.Load(context.Background(), iam.PolicyURN(id))
	if err != nil {
		return nil, err
	}

	policy.DefaultPolicy.ID = string(policy.ID)
	return &policy, nil
}

// GetAll implements ladon.Manager and returns all policies. Note that
// the limit and offset parameters are not yet supported.
func (p *PolicyManager) GetAll(limit, offset int64) (ladon.Policies, error) {
	if limit != 0 || offset != 0 {
		return nil, common.ErrNotImplemented
	}

	all, err := p.repo.Get(context.Background())
	if err != nil {
		return nil, err
	}

	policies := make(ladon.Policies, len(all))
	for i, p := range all {
		p.DefaultPolicy.ID = string(p.ID)
		policies[i] = &p
	}

	return policies, nil
}

// FindRequestCandidates returns all available policies.
func (p *PolicyManager) FindRequestCandidates(r *ladon.Request) (ladon.Policies, error) {
	return p.GetAll(0, 0)
}

// FindPoliciesForSubject return all available policies.
func (p *PolicyManager) FindPoliciesForSubject(subject string) (ladon.Policies, error) {
	return p.GetAll(0, 0)
}

// FindPoliciesForResource return all available policies.
func (p *PolicyManager) FindPoliciesForResource(resource string) (ladon.Policies, error) {
	return p.GetAll(0, 0)
}
