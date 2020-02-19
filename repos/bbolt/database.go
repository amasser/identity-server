package bbolt

import "go.etcd.io/bbolt"

var (
	userBucketKey = []byte("iam-v1-users")
)

// Database provides persistence for users, groups and policies
// using a bbolt database. It implements various iam.*Repository
// interfaces.
type Database struct {
	db *bbolt.DB
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
