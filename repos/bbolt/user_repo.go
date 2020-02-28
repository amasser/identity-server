package bbolt

import (
	"context"
	"encoding/json"

	"github.com/tierklinik-dobersberg/identity-server/pkg/common"
	"github.com/tierklinik-dobersberg/identity-server/pkg/iam"
	"go.etcd.io/bbolt"
)

var _ iam.UserRepository = &userRepo{}
var errUserNotFound = common.NewNotFoundError("user")

type userRepo struct {
	*Database
}

// Store impelements iam.UserRepository
func (db *userRepo) Store(ctx context.Context, user iam.User) error {
	blob, err := json.Marshal(user)
	if err != nil {
		return err
	}

	return db.db.Update(func(tx *bbolt.Tx) error {
		bucket, err := tx.CreateBucketIfNotExists(userBucketKey)
		if err != nil {
			return err
		}

		return bucket.Put([]byte(user.ID), blob)
	})
}

// Delete implements iam.UserRepository
func (db *userRepo) Delete(ctx context.Context, urn iam.UserURN) error {
	return db.db.Update(func(tx *bbolt.Tx) error {
		bucket := tx.Bucket(userBucketKey)
		if bucket == nil {
			return nil
		}

		return bucket.Delete([]byte(urn))
	})
}

// Load implements iam.UserRepository
func (db *userRepo) Load(ctx context.Context, urn iam.UserURN) (user iam.User, err error) {
	var blob []byte
	err = db.db.View(func(tx *bbolt.Tx) error {
		bucket := tx.Bucket(userBucketKey)
		if bucket == nil {
			return errUserNotFound
		}

		blob = bucket.Get([]byte(urn))
		return nil
	})
	if err != nil {
		return
	}

	if blob == nil {
		err = errUserNotFound
	} else {
		err = json.Unmarshal(blob, &user)
	}

	return
}

// Get implements iam.UserRepository
func (db *userRepo) Get(ctx context.Context) (users []iam.User, err error) {
	var blobs [][]byte
	err = db.db.View(func(tx *bbolt.Tx) error {
		bucket := tx.Bucket(userBucketKey)
		if bucket == nil {
			return nil
		}

		cursor := bucket.Cursor()
		key, blob := cursor.First()
		for key != nil {
			blobs = append(blobs, blob)
			key, blob = cursor.Next()
		}

		return nil
	})
	if err != nil {
		return
	}

	users = make([]iam.User, len(blobs))
	for i, b := range blobs {
		var u iam.User
		if err = json.Unmarshal(b, &u); err != nil {
			return
		}

		users[i] = u
	}

	return
}
