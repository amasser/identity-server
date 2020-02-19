package group

import (
	"context"
	"errors"

	"github.com/tierklinik-dobersberg/identity-server/iam"
	"github.com/tierklinik-dobersberg/identity-server/pkg/mutex"
	"github.com/tierklinik-dobersberg/identity-server/services/user"
)

// Service is the interface providing group management methods
type Service interface {
	// Create a new group with the given name and comment. Note that
	// groups will always be created without any initial members.
	// See AddMember for more information
	Create(ctx context.Context, groupName string, groupComment string) (iam.GroupURN, error)

	// Delete an existing account group and cancel the membership of all users
	Delete(ctx context.Context, urn iam.GroupURN) error

	// Load loads an existing account from the repository optionally including
	// a list of members URNs.
	Load(ctx context.Context, urn iam.GroupURN, withMembers bool) (iam.Group, error)

	// UpdateComment updates the comment of an account group
	UpdateComment(ctx context.Context, urn iam.GroupURN, comment string) error

	// AddMember adds a new memeber to the group
	AddMember(ctx context.Context, grp iam.GroupURN, memeber iam.UserURN) error

	// DeleteMember deletes a member from the group
	DeleteMember(ctx context.Context, grp iam.GroupURN, member iam.UserURN) error
}

type service struct {
	users user.Service
	l     *mutex.Mutex
	repo  iam.GroupRepository
}

// NewService returns a new service for managing account group memberships.
// It depends on having access to the user management service as well as
// a group repository for persisting changes.
func NewService(us user.Service, repo iam.GroupRepository) Service {
	return &service{
		users: us,
		l:     mutex.New(),
		repo:  repo,
	}
}

var errNotImplemented = errors.New("not yet implemented")

func (s *service) Create(ctx context.Context, groupName string, groupComment string) (iam.GroupURN, error) {
	return "", errNotImplemented
}

func (s *service) Delete(ctx context.Context, urn iam.GroupURN) error {
	return errNotImplemented
}

func (s *service) Load(ctx context.Context, urn iam.GroupURN, withMembers bool) (iam.Group, error) {
	return iam.Group{}, errNotImplemented
}

func (s *service) UpdateComment(ctx context.Context, urn iam.GroupURN, comment string) error {
	return errNotImplemented
}

func (s *service) AddMember(ctx context.Context, grp iam.GroupURN, memeber iam.UserURN) error {
	return errNotImplemented
}

func (s *service) DeleteMember(ctx context.Context, grp iam.GroupURN, member iam.UserURN) error {
	return errNotImplemented
}
