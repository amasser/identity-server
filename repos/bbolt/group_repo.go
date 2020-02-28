package bbolt

import (
	"context"
	"encoding/json"

	"github.com/tierklinik-dobersberg/identity-server/pkg/common"
	"github.com/tierklinik-dobersberg/identity-server/pkg/iam"
	"go.etcd.io/bbolt"
)

var errGroupNotFound = common.NewNotFoundError("group")

type groupRepo struct {
	*Database
}

func (db *groupRepo) Store(ctx context.Context, group iam.Group) error {
	blob, err := json.Marshal(group)
	if err != nil {
		return err
	}

	return db.db.Update(func(tx *bbolt.Tx) error {
		bucket, err := tx.CreateBucketIfNotExists(groupBucketKey)
		if err != nil {
			return err
		}

		return bucket.Put([]byte(group.ID), blob)
	})
}

func (db *groupRepo) Delete(ctx context.Context, urn iam.GroupURN) error {
	return db.db.Update(func(tx *bbolt.Tx) error {
		bucket := tx.Bucket(groupBucketKey)
		if bucket == nil {
			return errGroupNotFound
		}

		if bucket.Get([]byte(urn)) == nil {
			return errGroupNotFound
		}

		return bucket.Delete([]byte(urn))
	})
}

func (db *groupRepo) Load(ctx context.Context, urn iam.GroupURN) (iam.Group, error) {
	var grp iam.Group
	var blob []byte

	err := db.db.View(func(tx *bbolt.Tx) error {
		bucket := tx.Bucket(groupBucketKey)
		if bucket == nil {
			return errGroupNotFound
		}

		blob = bucket.Get([]byte(urn))
		if blob == nil {
			return errGroupNotFound
		}
		return nil
	})

	if err == nil {
		err = json.Unmarshal(blob, &grp)
	}

	return grp, err
}

func (db *groupRepo) Get(ctx context.Context) (groups []iam.Group, err error) {
	var blobs [][]byte

	err = db.db.View(func(tx *bbolt.Tx) error {
		bucket := tx.Bucket(groupBucketKey)
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

	groups = make([]iam.Group, len(blobs))
	for i, b := range blobs {
		var g iam.Group
		if err = json.Unmarshal(b, &g); err != nil {
			return
		}

		groups[i] = g
	}
	return
}
