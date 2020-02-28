package bbolt

import (
	"context"
	"encoding/json"

	"github.com/go-kit/kit/log/level"
	"github.com/tierklinik-dobersberg/identity-server/pkg/common"
	"github.com/tierklinik-dobersberg/identity-server/pkg/iam"
	"go.etcd.io/bbolt"
)

var errNoMember = common.NewNotFoundError("user membership")

type memberRepo struct {
	*Database
}

func (db *memberRepo) AddMember(ctx context.Context, user iam.UserURN, group iam.GroupURN) error {
	userKey := []byte(user)
	groupKey := []byte(group)

	return db.db.Update(func(tx *bbolt.Tx) error {
		userBucket, err := tx.CreateBucketIfNotExists(membershipUserBucketKey)
		if err != nil {
			return err
		}

		var existingGroups groupList
		if userBucket.Get(userKey) != nil {
			err = json.Unmarshal(userBucket.Get(userKey), &existingGroups)
			if err != nil {
				return err
			}
		}

		groupBucket, err := tx.CreateBucketIfNotExists(membershipGroupBucketKey)
		if err != nil {
			return err
		}

		var existingMembers userList
		if groupBucket.Get(groupKey) != nil {
			err = json.Unmarshal(groupBucket.Get(groupKey), &existingMembers)
			if err != nil {
				return err
			}
		}

		if !existingGroups.Has(group) {
			existingGroups = append(existingGroups, group)
			blob, err := json.Marshal(existingGroups)
			if err != nil {
				return err
			}

			if err := userBucket.Put(userKey, blob); err != nil {
				return err
			}
		}

		if !existingMembers.Has(user) {
			existingMembers = append(existingMembers, user)
			blob, err := json.Marshal(existingMembers)
			if err != nil {
				return err
			}

			if err := groupBucket.Put(groupKey, blob); err != nil {
				return err
			}
		}

		return nil
	})
}

func (db *memberRepo) DeleteMember(ctx context.Context, user iam.UserURN, group iam.GroupURN) error {
	userKey := []byte(user)
	groupKey := []byte(group)

	return db.db.Update(func(tx *bbolt.Tx) error {
		userBucket := tx.Bucket(membershipUserBucketKey)
		groupBucket := tx.Bucket(membershipGroupBucketKey)

		if userBucket == nil || groupBucket == nil {
			level.Debug(db.l).Log(
				"method", "DeleteMember",
				"user", user,
				"group", group,
				"msg", "Found empty buckets",
			)
			return errNoMember
		}

		groupBlob := groupBucket.Get(groupKey)
		userBlob := userBucket.Get(userKey)

		if groupBlob == nil || userBlob == nil {
			level.Debug(db.l).Log(
				"method", "DeleteMember",
				"user", user,
				"group", group,
				"msg", "either group or user entry not found",
			)
			return errNoMember
		}

		{
			var allMembers userList
			if err := json.Unmarshal(groupBlob, &allMembers); err != nil {
				return err
			}

			if !allMembers.Delete(user) {
				level.Debug(db.l).Log(
					"method", "DeleteMember",
					"user", user,
					"group", group,
					"msg", "user not part of member list",
				)
				return errNoMember
			}
			blob, err := json.Marshal(allMembers)
			if err != nil {
				return err
			}
			if err := groupBucket.Put(groupKey, blob); err != nil {
				return err
			}

			level.Debug(db.l).Log(
				"method", "DeleteMember",
				"user", user,
				"group", group,
				"members", len(allMembers),
				"msg", "updated member list",
			)
		}

		{
			var allGroups groupList
			if err := json.Unmarshal(userBlob, &allGroups); err != nil {
				return err
			}

			if !allGroups.Delete(group) {
				// TODO(ppacher): this acutally means we have an inconsistent data set
				level.Debug(db.l).Log(
					"method", "DeleteMember",
					"user", user,
					"group", group,
					"msg", "group not found in user groups",
				)
				return errNoMember
			}
			blob, err := json.Marshal(allGroups)
			if err != nil {
				return err
			}
			if err := userBucket.Put(userKey, blob); err != nil {
				return err
			}

			level.Debug(db.l).Log(
				"method", "DeleteMember",
				"user", user,
				"group", group,
				"groups", len(allGroups),
				"msg", "updated group list",
			)
		}

		level.Debug(db.l).Log(
			"method", "DeleteMember",
			"user", user,
			"group", group,
			"msg", "deleted user from group",
		)
		return nil
	})
}

func (db *memberRepo) Memberships(ctx context.Context, user iam.UserURN) ([]iam.GroupURN, error) {
	var groups groupList
	var blob []byte
	err := db.db.View(func(tx *bbolt.Tx) error {
		bucket := tx.Bucket(membershipUserBucketKey)
		if bucket == nil {
			return nil
		}

		blob = bucket.Get([]byte(user))
		return nil
	})
	if err != nil {
		return nil, err
	}

	if blob != nil {
		err = json.Unmarshal(blob, &groups)
	}
	return groups, err
}

func (db *memberRepo) Members(ctx context.Context, group iam.GroupURN) ([]iam.UserURN, error) {
	var members userList
	var blob []byte
	err := db.db.View(func(tx *bbolt.Tx) error {
		bucket := tx.Bucket(membershipGroupBucketKey)
		if bucket == nil {
			return nil
		}

		blob = bucket.Get([]byte(group))
		return nil
	})
	if err != nil {
		return nil, err
	}

	if blob != nil {
		err = json.Unmarshal(blob, &members)
	}
	return members, err
}

type userList []iam.UserURN

func (ul userList) Has(urn iam.UserURN) bool {
	for _, u := range ul {
		if u == urn {
			return true
		}
	}
	return false
}

func (ul *userList) Delete(urn iam.UserURN) bool {
	for i, u := range *ul {
		if u == urn {
			*ul = append((*ul)[:i], (*ul)[i+1:]...)
			return true
		}
	}

	return false
}

type groupList []iam.GroupURN

func (gl groupList) Has(urn iam.GroupURN) bool {
	for _, g := range gl {
		if g == urn {
			return true
		}
	}
	return false
}

func (gl *groupList) Delete(urn iam.GroupURN) bool {
	for i, g := range *gl {
		if g == urn {
			*gl = append((*gl)[:i], (*gl)[i+1:]...)
			return true
		}
	}

	return false
}
