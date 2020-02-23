package inmem

import (
	"context"
	"sync"

	"github.com/tierklinik-dobersberg/identity-server/iam"
	"github.com/tierklinik-dobersberg/identity-server/pkg/common"
)

type groupRepo struct {
	rw     sync.RWMutex
	groups map[iam.GroupURN]iam.Group
}

// NewGroupRepository creates a new in-memory group repository
func NewGroupRepository() iam.GroupRepository {
	return &groupRepo{
		groups: make(map[iam.GroupURN]iam.Group),
	}
}

func (repo *groupRepo) Store(ctx context.Context, group iam.Group) error {
	repo.rw.Lock()
	defer repo.rw.Unlock()

	repo.groups[group.ID] = group

	return ctx.Err()
}

func (repo *groupRepo) Delete(ctx context.Context, urn iam.GroupURN) error {
	repo.rw.Lock()
	defer repo.rw.Unlock()

	delete(repo.groups, urn)

	return ctx.Err()
}

func (repo *groupRepo) Load(ctx context.Context, urn iam.GroupURN) (iam.Group, error) {
	repo.rw.RLock()
	defer repo.rw.RUnlock()

	if g, ok := repo.groups[urn]; ok {
		return g, ctx.Err()
	}

	return iam.Group{}, common.NewNotFoundError("group")
}

func (repo *groupRepo) Get(ctx context.Context) ([]iam.Group, error) {
	repo.rw.RLock()
	defer repo.rw.RUnlock()

	groups := make([]iam.Group, 0, len(repo.groups))

	for _, g := range repo.groups {
		groups = append(groups, g)
	}

	return groups, ctx.Err()
}
