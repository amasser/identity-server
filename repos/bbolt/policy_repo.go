package bbolt

import (
	"context"
	"encoding/json"

	"github.com/tierklinik-dobersberg/identity-server/iam"
	"github.com/tierklinik-dobersberg/identity-server/pkg/common"
	"go.etcd.io/bbolt"
)

var errPolicyNotFound = common.NewNotFoundError("policy")

type policyRepo struct {
	*Database
}

func (db *policyRepo) Store(ctx context.Context, policy iam.Policy) error {
	blob, err := json.Marshal(policy)
	if err != nil {
		return err
	}

	return db.db.Update(func(tx *bbolt.Tx) error {
		b, err := tx.CreateBucketIfNotExists(policyBucketKey)
		if err != nil {
			return err
		}

		return b.Put([]byte(policy.ID), blob)
	})
}

func (db *policyRepo) Delete(ctx context.Context, urn iam.PolicyURN) error {
	return db.db.Update(func(tx *bbolt.Tx) error {
		b := tx.Bucket(policyBucketKey)
		if b == nil {
			return errPolicyNotFound
		}

		if b.Get([]byte(urn)) == nil {
			return errPolicyNotFound
		}

		return b.Delete([]byte(urn))
	})
}

func (db *policyRepo) Load(ctx context.Context, urn iam.PolicyURN) (iam.Policy, error) {
	var blob []byte
	err := db.db.View(func(tx *bbolt.Tx) error {
		b := tx.Bucket(policyBucketKey)
		if b == nil {
			return errPolicyNotFound
		}

		blob = b.Get([]byte(urn))
		if blob == nil {
			return errPolicyNotFound
		}
		return nil
	})

	if err != nil {
		return iam.Policy{}, err
	}

	var p iam.Policy
	if err := json.Unmarshal(blob, &p); err != nil {
		return p, nil
	}

	return p, nil
}

func (db *policyRepo) Get(ctx context.Context) ([]iam.Policy, error) {
	var blobs [][]byte
	err := db.db.View(func(tx *bbolt.Tx) error {
		b := tx.Bucket(policyBucketKey)
		if b == nil {
			return nil
		}

		cursor := b.Cursor()
		key, value := cursor.First()

		for key != nil {
			blobs = append(blobs, value)
			key, value = cursor.Next()
		}

		return nil
	})
	if err != nil {
		return nil, err
	}

	policies := make([]iam.Policy, len(blobs))
	for i, b := range blobs {
		var p iam.Policy
		if err := json.Unmarshal(b, &p); err != nil {
			return nil, err
		}

		policies[i] = p
	}

	return policies, nil
}
