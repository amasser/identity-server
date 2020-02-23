package bbolt

import (
	"github.com/tierklinik-dobersberg/identity-server/iam"
	"go.etcd.io/bbolt"
)

var (
	userBucketKey       = []byte("iam-v1-users")
	groupBucketKey      = []byte("iam-v1-groups")
	membershipBucketKey = []byte("iam-v1-memberships")
)

// Database provides persistence for users, groups and policies
// using a bbolt database. It implements various iam.*Repository
// interfaces.
type Database struct {
	db *bbolt.DB
}

// UserRepo returns a iam.UserRepository backed by db
func (db *Database) UserRepo() iam.UserRepository {
	return &userRepo{db}
}

// GroupRepo returns a iam.GroupRepository backed by db
func (db *Database) GroupRepo() iam.GroupRepository {
	return &groupRepo{db}
}

// Open opes the database file at path and returns
// a new Database instance
func Open(path string) (*Database, error) {
	db, err := bbolt.Open(path, 0600, nil)
	if err != nil {
		return nil, err
	}

	return &Database{db: db}, nil
}
