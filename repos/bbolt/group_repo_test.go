package bbolt

import (
	"context"
	"encoding/json"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/tierklinik-dobersberg/identity-server/pkg/iam"
	"go.etcd.io/bbolt"
)

func TestGroup_Store(t *testing.T) {
	db, cleanup := getTempGroupRepoWithData(t)
	defer cleanup()

	err := db.Store(context.Background(), testGroup)
	assert.NoError(t, err)

	db.db.View(func(tx *bbolt.Tx) error {
		b := tx.Bucket(groupBucketKey)

		assert.NotNil(t, b.Get([]byte("urn:iam::group/admins")))
		return nil
	})
}

func TestGroup_Delete(t *testing.T) {
	db, cleanup := getTempGroupRepoWithData(t)
	defer cleanup()

	err := db.Delete(context.Background(), "urn:iam::group/admins")
	assert.NoError(t, err)

	err = db.Delete(context.Background(), "urn:iam::group/admins")
	assert.Error(t, err)
}

func TestGroup_Load(t *testing.T) {
	db, cleanup := getTempGroupRepoWithData(t)
	defer cleanup()

	grp, err := db.Load(context.Background(), "urn:iam::group/admins")
	assert.NoError(t, err)
	assert.Equal(t, testGroup, grp)

	grp, err = db.Load(context.Background(), "urn:iam::group/devs")
	assert.Error(t, err)
	assert.Equal(t, iam.Group{}, grp)
}

func TestGroup_Get(t *testing.T) {
	db, cleanup := getTempGroupRepoWithData(t)
	defer cleanup()

	grps, err := db.Get(context.Background())
	assert.NoError(t, err)
	assert.Equal(t, []iam.Group{testGroup}, grps)
}

var testGroup = iam.Group{
	ID:      "urn:iam::group/admins",
	Name:    "admins",
	Comment: "A group used for unit tests",
}

func testOpenGroupRepo(t *testing.T, f string) *groupRepo {
	db, err := Open(f)
	if err != nil {
		t.FailNow()
		return nil
	}

	return db.GroupRepo().(*groupRepo)
}

func getTempGroupRepoWithData(t *testing.T) (*groupRepo, func()) {
	s, c := getTempDb()
	db := testOpenGroupRepo(t, s)
	blob, _ := json.Marshal(testGroup)
	err := db.db.Update(func(tx *bbolt.Tx) error {
		b, err := tx.CreateBucket(groupBucketKey)
		if err != nil {
			return err
		}

		return b.Put([]byte(testGroup.ID), blob)
	})

	if err != nil {
		t.FailNow()
		return nil, nil
	}

	return db, c
}
