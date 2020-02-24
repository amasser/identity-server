package bbolt

import (
	"github.com/go-kit/kit/log"
	"github.com/tierklinik-dobersberg/identity-server/iam"
	"go.etcd.io/bbolt"
)

var (
	userBucketKey            = []byte("iam-v1-users")
	groupBucketKey           = []byte("iam-v1-groups")
	membershipGroupBucketKey = []byte("iam-v1-memberships-group")
	membershipUserBucketKey  = []byte("iam-v1-memberships-user")
)

// Database provides persistence for users, groups and policies
// using a bbolt database. It implements various iam.*Repository
// interfaces.
type Database struct {
	db *bbolt.DB
	l  log.Logger
}

// UserRepo returns a iam.UserRepository backed by db.
func (db *Database) UserRepo() iam.UserRepository {
	return &userRepo{db}
}

// GroupRepo returns a iam.GroupRepository backed by db.
func (db *Database) GroupRepo() iam.GroupRepository {
	return &groupRepo{db}
}

// MembershipRepo returns a iam.MembershipRepository backed by db.
func (db *Database) MembershipRepo() iam.MembershipRepository {
	return &memberRepo{db}
}

// Open opes the database file at path and returns
// a new Database instance
func Open(path string) (*Database, error) {
	return OpenWithLogger(path, log.NewNopLogger())
}

// OpenWithLogger is like Open but allows specifying a logger to use
func OpenWithLogger(path string, l log.Logger) (*Database, error) {
	db, err := bbolt.Open(path, 0600, nil)
	if err != nil {
		return nil, err
	}

	return &Database{
		db: db,
		l:  l,
	}, nil
}
