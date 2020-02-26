package bbolt

import (
	"encoding/json"
	"io/ioutil"
	"log"
	"os"
	"testing"

	"github.com/stretchr/testify/require"
	"github.com/tierklinik-dobersberg/identity-server/iam"
	"go.etcd.io/bbolt"
)

func getTempDb() (string, func()) {
	file, err := ioutil.TempFile("", "unit-test-db-*.db")
	if err != nil {
		log.Fatal(err)
	}

	return file.Name(), func() {
		os.Remove(file.Name())
	}
}

var testExistingUser = iam.User{
	AccountID: 10,
	Username:  "admin",
	ID:        "urn:iam::user/10",
	Attributes: map[string]interface{}{
		"job": "developer",
	},
}

func getTempUserRepoWithData(t *testing.T) (*userRepo, func()) {
	s, c := getTempDb()

	db, err := testOpenUserRepo(s)
	require.NoError(t, err)

	blob, _ := json.Marshal(testExistingUser)

	err = db.db.Update(func(tx *bbolt.Tx) error {
		b, err := tx.CreateBucket(userBucketKey)
		if err != nil {
			return err
		}

		return b.Put([]byte("urn:iam::user/10"), blob)
	})
	require.NoError(t, err)

	return db, c
}
