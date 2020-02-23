package group

import (
	"context"
	"fmt"

	"github.com/go-kit/kit/log"
	"github.com/tierklinik-dobersberg/identity-server/iam"
	"github.com/tierklinik-dobersberg/identity-server/pkg/common"
	"github.com/tierklinik-dobersberg/identity-server/pkg/mutex"
	"github.com/tierklinik-dobersberg/identity-server/services/user"
)

// Service is the interface providing group management methods
type Service interface {
	// Get returns a list of all groups.
	Get(ctx context.Context) ([]iam.Group, error)

	// Create a new group with the given name and comment. Note that
	// groups will always be created without any initial members.
	// See AddMember for more information.
	Create(ctx context.Context, groupName string, groupComment string) (iam.GroupURN, error)

	// Delete an existing account group and cancel the membership of all users.
	Delete(ctx context.Context, urn iam.GroupURN) error

	// Load loads an existing account from the repository optionally including
	// a list of members URNs.
	Load(ctx context.Context, urn iam.GroupURN) (iam.Group, error)

	// UpdateComment updates the comment of an account group.
	UpdateComment(ctx context.Context, urn iam.GroupURN, comment string) error

	// AddMember adds a new memeber to the group.
	AddMember(ctx context.Context, grp iam.GroupURN, memeber iam.UserURN) error

	// DeleteMember deletes a member from the group.
	DeleteMember(ctx context.Context, grp iam.GroupURN, member iam.UserURN) error
}

type service struct {
	users   user.Service
	l       *mutex.Mutex
	groups  iam.GroupRepository
	members iam.MembershipRepository
	log     log.Logger
}

// NewService returns a new service for managing account group memberships.
// It depends on having access to the user management service as well as
// a group repository for persisting changes.
func NewService(us user.Service, groups iam.GroupRepository, members iam.MembershipRepository, logger log.Logger) Service {
	svc := &service{
		users:   us,
		l:       mutex.New(),
		groups:  groups,
		members: members,
		log:     logger,
	}

	us.OnDelete(context.Background(), svc.userDeleted)

	return svc
}

// ErrInvalidParameter is returned from the group management service if
// invalid parameters are supplied
var ErrInvalidParameter = common.NewInvalidArgumentError("invalid parameter")

func (s *service) Get(ctx context.Context) ([]iam.Group, error) {
	if !s.l.TryLock(ctx) {
		return nil, ctx.Err()
	}
	defer s.l.Unlock()

	return s.groups.Get(ctx)
}

func (s *service) Create(ctx context.Context, groupName string, groupComment string) (iam.GroupURN, error) {
	if groupName == "" {
		return "", ErrInvalidParameter
	}

	if !s.l.TryLock(ctx) {
		return "", ctx.Err()
	}
	defer s.l.Unlock()

	grp := iam.Group{
		ID:      iam.GroupURN(fmt.Sprintf("urn:iam::group/%s", groupName)),
		Name:    groupName,
		Comment: groupComment,
	}

	if err := s.groups.Store(ctx, grp); err != nil {
		return "", err
	}

	return grp.ID, nil
}

func (s *service) Delete(ctx context.Context, urn iam.GroupURN) error {
	if urn == "" {
		return ErrInvalidParameter
	}

	if !s.l.TryLock(ctx) {
		return ctx.Err()
	}
	defer s.l.Unlock()

	members, err := s.members.Members(ctx, urn)
	if err != nil {
		return err
	}

	for _, m := range members {
		if err := s.members.DeleteMember(ctx, m, urn); err != nil {
			return err
		}
	}

	return s.groups.Delete(ctx, urn)
}

func (s *service) Load(ctx context.Context, urn iam.GroupURN) (iam.Group, error) {
	if urn == "" {
		return iam.Group{}, ErrInvalidParameter
	}
	if !s.l.TryLock(ctx) {
		return iam.Group{}, ctx.Err()
	}
	defer s.l.Unlock()

	grp, err := s.groups.Load(ctx, urn)
	return grp, err
}

func (s *service) UpdateComment(ctx context.Context, urn iam.GroupURN, comment string) error {
	if urn == "" {
		return ErrInvalidParameter
	}
	if !s.l.TryLock(ctx) {
		return ctx.Err()
	}
	defer s.l.Unlock()

	grp, err := s.groups.Load(ctx, urn)
	if err != nil {
		return err
	}

	grp.Comment = comment

	return s.groups.Store(ctx, grp)
}

func (s *service) AddMember(ctx context.Context, grp iam.GroupURN, member iam.UserURN) error {
	if !s.l.TryLock(ctx) {
		return ctx.Err()
	}
	defer s.l.Unlock()

	// Ensure the target group exists.
	if _, err := s.groups.Load(ctx, grp); err != nil {
		return err
	}

	// Ensure the user actually exists.
	// We still might race with a Delete() call for the user,
	// however, this isn't a problem as the callback for OnDelete()
	// will be blocked until we finished adding the user to the
	// group. Once it is unblocked, the user will be removed
	// again.
	if _, err := s.users.LoadUser(ctx, member); err != nil {
		return err
	}

	// MembershipRepository already handles the case of a user
	// already being part of the group. It's a no-op then.
	if err := s.members.AddMember(ctx, member, grp); err != nil {
		return err
	}

	return nil
}

func (s *service) DeleteMember(ctx context.Context, grp iam.GroupURN, member iam.UserURN) error {
	if !s.l.TryLock(ctx) {
		return ctx.Err()
	}
	defer s.l.Unlock()

	if _, err := s.groups.Load(ctx, grp); err != nil {
		return err
	}

	if err := s.members.DeleteMember(ctx, member, grp); err != nil {
		return err
	}

	return nil
}

func (s *service) userDeleted(user iam.UserURN) {
	s.l.Lock()
	defer s.l.Unlock()

	groups, err := s.members.Memberships(context.Background(), user)
	if err != nil {
		// TODO: we might have zombie users in that group now
		// maybe retry?
		s.log.Log(
			"method", "userDeleted",
			"user", user,
			"err", err,
			"msg", "Failed to get group memberships",
		)
		return
	}

	for _, g := range groups {
		if err := s.members.DeleteMember(context.Background(), user, g); err != nil {
			// TODO: we might have zombie users in that group now
			// maybe retry?
			s.log.Log(
				"method", "userDeleted",
				"user", user,
				"group", g,
				"err", err,
				"msg", "Failed to remove user from group",
			)
		}
	}
}
