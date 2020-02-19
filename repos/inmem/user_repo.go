package inmem

import (
	"context"
	"os"
	"sync"

	"github.com/tierklinik-dobersberg/iam/v2/iam"
)

type userRepository struct {
	l     sync.RWMutex
	users map[iam.UserURN]iam.User
}

func (r *userRepository) Store(ctx context.Context, user iam.User) error {
	r.l.Lock()
	defer r.l.Unlock()

	r.users[user.ID] = user

	return ctx.Err()
}

func (r *userRepository) Delete(ctx context.Context, urn iam.UserURN) error {
	r.l.Lock()
	defer r.l.Unlock()

	if _, ok := r.users[urn]; !ok {
		return os.ErrNotExist
	}

	delete(r.users, urn)

	return ctx.Err()
}

func (r *userRepository) Load(ctx context.Context, urn iam.UserURN) (iam.User, error) {
	r.l.RLock()
	defer r.l.RUnlock()

	if u, ok := r.users[urn]; ok {
		return u, nil
	}

	return iam.User{}, os.ErrNotExist
}

func (r *userRepository) Get(ctx context.Context) ([]iam.User, error) {
	r.l.RLock()
	defer r.l.RUnlock()

	users := make([]iam.User, 0, len(r.users))
	for _, u := range r.users {
		users = append(users, u)
	}

	return users, ctx.Err()
}

// NewUserRepository returns a new in-memory user repository
func NewUserRepository() iam.UserRepository {
	return &userRepository{
		users: make(map[iam.UserURN]iam.User),
	}
}
