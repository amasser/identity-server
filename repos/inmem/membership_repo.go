package inmem

import (
	"context"
	"sync"

	"github.com/tierklinik-dobersberg/identity-server/iam"
	"github.com/tierklinik-dobersberg/identity-server/pkg/common"
)

type membershipRepo struct {
	rw      sync.RWMutex
	members map[iam.GroupURN][]iam.UserURN
	users   map[iam.UserURN][]iam.GroupURN
}

// NewMembershipRepository creates a new in-memory membership
// repository
func NewMembershipRepository() iam.MembershipRepository {
	return &membershipRepo{
		members: make(map[iam.GroupURN][]iam.UserURN),
		users:   make(map[iam.UserURN][]iam.GroupURN),
	}
}

// AddMember adds a user as a new member to grp. If user is already a member of grp AddMember
// is a no-op.
func (repo *membershipRepo) AddMember(ctx context.Context, user iam.UserURN, grp iam.GroupURN) error {
	repo.rw.Lock()
	defer repo.rw.Unlock()

	members := repo.members[grp]
	for _, m := range members {
		if m == user {
			return nil
		}
	}
	members = append(members, user)
	groups := append(repo.users[user], grp)

	repo.members[grp] = members
	repo.users[user] = groups

	return nil
}

func (repo *membershipRepo) DeleteMember(ctx context.Context, user iam.UserURN, grp iam.GroupURN) error {
	repo.rw.Lock()
	defer repo.rw.RUnlock()

	members := repo.members[grp]
	foundMember := false
	for i, m := range members {
		if m == user {
			members = append(members[:i], members[i+1:]...)
			foundMember = true
			break
		}
	}

	if foundMember {
		groups := repo.users[user]
		for i, g := range groups {
			if g == grp {
				groups = append(groups[:i], groups[i+1:]...)
				break
			}
		}

		repo.members[grp] = members
		repo.users[user] = groups

		return nil
	}

	return common.NewNotFoundError("user membership")
}

func (repo *membershipRepo) Members(ctx context.Context, grp iam.GroupURN) ([]iam.UserURN, error) {
	repo.rw.RLock()
	defer repo.rw.RUnlock()

	members := repo.members[grp]
	users := make([]iam.UserURN, len(members))

	for i, u := range members {
		users[i] = u
	}

	return users, nil
}

func (repo *membershipRepo) Memberships(ctx context.Context, user iam.UserURN) ([]iam.GroupURN, error) {
	repo.rw.RLock()
	defer repo.rw.RUnlock()

	groups := repo.users[user]
	result := make([]iam.GroupURN, len(groups))

	for i, g := range groups {
		result[i] = g
	}

	return result, nil
}
