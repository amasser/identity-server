// Package user provides the use-case of creating and managing users. Used by
// views facing an adimistrator.
package user

import (
	"context"
	"fmt"
	"os"
	"sync"
	"sync/atomic"

	"github.com/tierklinik-dobersberg/identity-server/iam"
	"github.com/tierklinik-dobersberg/identity-server/pkg/authn"
	"github.com/tierklinik-dobersberg/identity-server/pkg/common"
	"github.com/tierklinik-dobersberg/identity-server/pkg/mutex"
)

// ErrInvalidArgument is returned when an invalid argument is passed to
// a Service method
var ErrInvalidArgument = common.NewInvalidArgumentError("invalid argument")

// OnDeleteFunc is a callback function that is invoked when a user is deleted.
// See Service.OnDelete() for more information.
type OnDeleteFunc func(urn iam.UserURN)

// Service is the interface that provides user management methods.
type Service interface {
	// CreateUser creates a new user account in the user management system and returns
	// the new unique user URN.
	CreateUser(ctx context.Context, username, password string, attrs map[string]interface{}) (iam.UserURN, error)

	// LoadUser returns the read model of a user.
	LoadUser(ctx context.Context, urn iam.UserURN) (iam.User, error)

	// DeleteUser deletes the user account from IAM and archives it
	// on authn-server
	DeleteUser(ctx context.Context, urn iam.UserURN) error

	// LockUser locks or unlocks a user account
	LockUser(ctx context.Context, urn iam.UserURN, locked bool) error

	// Users returns the read model of all available users.
	Users(ctx context.Context) ([]iam.User, error)

	// UpdateAttrs replaces all user attributes from user `id` with `attrs`.
	UpdateAttrs(ctx context.Context, id iam.UserURN, attrs map[string]interface{}) error

	// SetAttr updates the attribute key with value of the user identified by id.
	SetAttr(ctx context.Context, id iam.UserURN, key string, value interface{}) error

	// DeleteAttr deletes the attr key from the user identified by `id`
	DeleteAttr(ctx context.Context, id iam.UserURN, key string) error

	// OnDelete registers a callback function that is invoked whenever a user
	// is deleted/archived. Note that the callback function *may* be unregistered
	// when the provided context is cancelled.
	OnDelete(ctx context.Context, fn OnDeleteFunc)
}

type service struct {
	authn authn.Service
	m     *mutex.Mutex
	repo  iam.UserRepository

	deleteFnsLock sync.RWMutex
	deleteFns     map[int64]OnDeleteFunc
}

func (s *service) CreateUser(ctx context.Context, username, password string, attrs map[string]interface{}) (urn iam.UserURN, err error) {
	if !s.m.TryLock(ctx) {
		return "", ctx.Err()
	}
	defer s.m.Unlock()

	accountID, err := s.authn.ImportAccount(username, password, false)
	if err != nil {
		return "", err
	}
	defer func() {
		if err != nil {
			// TODO(ppacher): log that!
			s.authn.ArchiveAccount(accountID)
		}
	}()

	urn = iam.UserURN(fmt.Sprintf("urn:iam::user/%d", accountID))

	_, err = s.repo.Load(ctx, urn)
	if err == nil {
		return "", common.NewConflictError("username")
	}
	if !common.IsNotFound(err) && !os.IsNotExist(err) {
		return "", err
	}

	user := iam.User{
		AccountID:  accountID,
		Username:   username,
		ID:         urn,
		Attributes: attrs,
	}

	return urn, s.repo.Store(ctx, user)
}

func (s *service) LoadUser(ctx context.Context, urn iam.UserURN) (iam.User, error) {
	if urn == "" {
		return iam.User{}, ErrInvalidArgument
	}
	if !s.m.TryLock(ctx) {
		return iam.User{}, ctx.Err()
	}
	defer s.m.Unlock()

	return s.repo.Load(ctx, urn)
}

func (s *service) Users(ctx context.Context) ([]iam.User, error) {
	if !s.m.TryLock(ctx) {
		return nil, ctx.Err()
	}
	defer s.m.Unlock()

	return s.repo.Get(ctx)
}

func (s *service) DeleteUser(ctx context.Context, urn iam.UserURN) error {
	if urn == "" {
		return ErrInvalidArgument
	}

	if !s.m.TryLock(ctx) {
		return ctx.Err()
	}
	defer s.m.Unlock()

	user, err := s.repo.Load(ctx, urn)
	if err != nil {
		return err
	}

	if err := s.authn.ArchiveAccount(user.AccountID); err != nil {
		return err
	}

	// notify all on-delete subscribes even
	// if the actual delete operation fails
	s.deleteFnsLock.RLock()
	for _, fn := range s.deleteFns {
		go fn(urn)
	}
	s.deleteFnsLock.RUnlock()

	return s.repo.Delete(ctx, urn)
}

func (s *service) LockUser(ctx context.Context, urn iam.UserURN, locked bool) error {
	if urn == "" {
		return ErrInvalidArgument
	}

	if !s.m.TryLock(ctx) {
		return ctx.Err()
	}
	defer s.m.Unlock()

	user, err := s.repo.Load(ctx, urn)
	if err != nil {
		return err
	}

	if locked {
		err = s.authn.LockAccount(user.AccountID)
	} else {
		err = s.authn.UnlockAccount(user.AccountID)
	}

	return err
}

func (s *service) UpdateAttrs(ctx context.Context, urn iam.UserURN, attr map[string]interface{}) error {
	if urn == "" {
		return ErrInvalidArgument
	}

	if !s.m.TryLock(ctx) {
		return ctx.Err()
	}
	defer s.m.Unlock()

	user, err := s.repo.Load(ctx, urn)
	if err != nil {
		return err
	}

	user.Attributes = attr
	return s.repo.Store(ctx, user)
}

func (s *service) SetAttr(ctx context.Context, urn iam.UserURN, key string, value interface{}) error {
	if urn == "" {
		return ErrInvalidArgument
	}

	if !s.m.TryLock(ctx) {
		return ctx.Err()
	}
	defer s.m.Unlock()

	user, err := s.repo.Load(ctx, urn)
	if err != nil {
		return err
	}

	if user.Attributes == nil {
		user.Attributes = make(map[string]interface{})
	}

	user.Attributes[key] = value
	return s.repo.Store(ctx, user)
}

func (s *service) DeleteAttr(ctx context.Context, urn iam.UserURN, key string) error {
	if urn == "" {
		return ErrInvalidArgument
	}

	if !s.m.TryLock(ctx) {
		return ctx.Err()
	}
	defer s.m.Unlock()

	user, err := s.repo.Load(ctx, urn)
	if err != nil {
		return err
	}

	if user.Attributes == nil {
		return nil
	}

	delete(user.Attributes, key)
	return s.repo.Store(ctx, user)
}

func (s *service) OnDelete(ctx context.Context, fn OnDeleteFunc) {
	id := getUniqueSubscriberID()

	s.deleteFnsLock.Lock()
	s.deleteFns[id] = fn
	s.deleteFnsLock.Unlock()
}

// NewService creates a new user management services
func NewService(repo iam.UserRepository, authn authn.Service) Service {
	return &service{
		authn:     authn,
		m:         mutex.New(),
		repo:      repo,
		deleteFns: make(map[int64]OnDeleteFunc, 10),
	}
}

var nextUniqueSubscriber int64

func getUniqueSubscriberID() int64 {
	return atomic.AddInt64(&nextUniqueSubscriber, 1)
}
