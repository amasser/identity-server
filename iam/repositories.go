package iam

import "context"

// UserRepository provides persistent storage for user account information
type UserRepository interface {
	// Store stores a user, overwriting and existing one if necassary
	Store(ctx context.Context, user User) error

	// Delete deletes the user with the given urn. Implementations
	// do not need to check if the user actually exists. If such a
	// check exists, implementations are encouraged to return
	// os.ErrNotExist so transports can convert it to an appropriate
	// error code (like 404 Not Found).
	Delete(ctx context.Context, urn UserURN) error

	// Load returns a read model of the user identified by urn
	Load(ctx context.Context, urn UserURN) (User, error)

	// Get returns a list of read models of all users
	Get(ctx context.Context) ([]User, error)
}
