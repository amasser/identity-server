package bbolt

import (
	"context"
	"encoding/json"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"github.com/tierklinik-dobersberg/identity-server/pkg/iam"
	"go.etcd.io/bbolt"
)

func testOpenUserRepo(f string) (*userRepo, error) {
	db, err := Open(f)
	if err != nil {
		return nil, err
	}

	return db.UserRepo().(*userRepo), nil
}

func Test_Store(t *testing.T) {
	f, cleanup := getTempDb()
	defer cleanup()
	db, err := testOpenUserRepo(f)
	require.NoError(t, err)

	user := iam.User{
		AccountID: 10,
		Username:  "admin",
		ID:        "urn:iam::user/10",
		Attributes: map[string]interface{}{
			"job": "developer",
		},
	}

	err = db.Store(context.Background(), user)
	assert.NoError(t, err)

	db.db.View(func(tx *bbolt.Tx) error {
		bucket := tx.Bucket(userBucketKey)
		require.NotNil(t, bucket)
		key := []byte("urn:iam::user/10")
		value := bucket.Get(key)
		require.NotNil(t, value)

		var u iam.User
		require.NoError(t, json.Unmarshal(value, &u))
		require.Equal(t, user, u)
		return nil
	})

	// store must overwrite
	user = iam.User{
		AccountID: 10,
		Username:  "admin2",
		ID:        "urn:iam::user/10",
		Attributes: map[string]interface{}{
			"job": "developer",
		},
	}
	err = db.Store(context.Background(), user)
	assert.NoError(t, err)

	db.db.View(func(tx *bbolt.Tx) error {
		bucket := tx.Bucket(userBucketKey)
		require.NotNil(t, bucket)
		key := []byte("urn:iam::user/10")
		value := bucket.Get(key)
		require.NotNil(t, value)

		var u iam.User
		require.NoError(t, json.Unmarshal(value, &u))
		require.Equal(t, user, u)
		return nil
	})
}

func Test_Delete(t *testing.T) {
	db, cleanup := getTempUserRepoWithData(t)
	defer cleanup()

	assert.NoError(t, db.Delete(context.Background(), "urn:iam::user/10"))
	db.db.View(func(tx *bbolt.Tx) error {
		b := tx.Bucket(userBucketKey)
		key := []byte("urn:iam::user/10")
		assert.Nil(t, b.Get(key))

		return nil
	})
}

func Test_Load(t *testing.T) {
	db, cleanup := getTempUserRepoWithData(t)
	defer cleanup()

	u, err := db.Load(context.Background(), "urn:iam::user/10")
	assert.NoError(t, err)
	assert.Equal(t, testExistingUser, u)

	u, err = db.Load(context.Background(), "urn:iam::user/100")
	assert.Error(t, err)
	assert.Equal(t, iam.User{}, u)
}

func Test_Get(t *testing.T) {
	db, cleanup := getTempUserRepoWithData(t)
	defer cleanup()

	list, err := db.Get(context.Background())
	assert.NoError(t, err)
	assert.Equal(t, []iam.User{testExistingUser}, list)
}
