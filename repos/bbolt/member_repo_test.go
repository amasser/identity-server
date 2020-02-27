package bbolt

import (
	"testing"

	"go.etcd.io/bbolt"
)

func testOpenMemberRepo(t *testing.T, f string) *memberRepo {
	db, err := Open(f)
	if err != nil {
		t.FailNow()
		return nil
	}

	return db.MembershipRepo().(*memberRepo)
}

func getTempMemberRepoWithData(t *testing.T) (*memberRepo, func()) {
	s, c := getTempDb()
	db := testOpenMemberRepo(t, s)

	err := db.db.Update(func(tx *bbolt.Tx) error {
		mu, err := tx.CreateBucket(membershipUserBucketKey)
		if err != nil {
			return err
		}

		mg, err := tx.CreateBucket(membershipGroupBucketKey)
		if err != nil {
			return err
		}

		if err := mu.Put([]byte("urn:iam::user/10"), []byte(`[
			"urn:iam::group/admins"	
		]`)); err != nil {
			return err
		}

		if err := mg.Put([]byte("urn:iam::group/admins"), []byte(`[
			"urn:iam::user/10"
		]`)); err != nil {
			return err
		}

		return nil
	})

	if err != nil {
		t.FailNow()
	}

	return db, c
}
