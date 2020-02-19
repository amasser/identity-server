package iam

import "context"

// UserRepository provides persistent storage for user account information.
type UserRepository interface {
	// Store stores a user, overwriting and existing one if necassary.
	Store(ctx context.Context, user User) error

	// Delete deletes the user with the given urn. Implementations
	// do not need to check if the user actually exists. If such a
	// check exists, implementations are encouraged to return
	// os.ErrNotExist so transports can convert it to an appropriate
	// error code (like 404 Not Found).
	Delete(ctx context.Context, urn UserURN) error

	// Load returns a read model of the user identified by urn.
	Load(ctx context.Context, urn UserURN) (User, error)

	// Get returns a list of read models of all users.
	Get(ctx context.Context) ([]User, error)
}

// GroupRepository provides persistent storage for group information.
type GroupRepository interface {
	// Store stores a group, overwritting and existing one if necassary.
	// The groups member slice should only be persisted if storeMembers
	// is set to true. Otherwise any previously set member slice must be
	// copied for the new group.
	Store(ctx context.Context, group Group, storeMembers bool) error

	// Delete deletes an existing account group. Implementations tracking
	// user-group assignments should clean them up as well. It is not required
	// for the group to be deleted to exist. If it does not exist, implementations
	// should either return nil or os.ErrNotExist. Any other error will be treated
	// as something else and returned back to the user.
	Delete(ctx context.Context, urn GroupURN) error

	// Load loads the account group from the persistent storage and returns it.
	// If the group does not exist os.ErrNotFound should be returned.
	Load(ctx context.Context, urn GroupURN) (Group, error)

	// Get should return all groups from the persistent storage.
	Get(ctx context.Context) ([]Group, error)
}
