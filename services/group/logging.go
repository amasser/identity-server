package group

import (
	"context"
	"time"

	"github.com/go-kit/kit/log"
	"github.com/tierklinik-dobersberg/identity-server/pkg/iam"
)

type loggingService struct {
	l log.Logger
	Service
}

// NewLoggingService returns a service that logs method calls of Service
func NewLoggingService(s Service, logger log.Logger) Service {
	return &loggingService{
		l:       logger,
		Service: s,
	}
}

func (s *loggingService) Get(ctx context.Context) (grps []iam.Group, err error) {
	defer func(begin time.Time) {
		s.l.Log(
			"method", "list_groups",
			"groups", len(grps),
			"took", time.Since(begin),
			"err", err,
		)
	}(time.Now())
	return s.Service.Get(ctx)
}

func (s *loggingService) Create(ctx context.Context, groupName, groupComment string) (urn iam.GroupURN, err error) {
	defer func(begin time.Time) {
		s.l.Log(
			"method", "create_group",
			"name", groupName,
			"comment", groupComment,
			"took", time.Since(begin),
			"urn", urn,
			"err", err,
		)
	}(time.Now())
	return s.Service.Create(ctx, groupName, groupComment)
}

func (s *loggingService) Delete(ctx context.Context, urn iam.GroupURN) (err error) {
	defer func(begin time.Time) {
		s.l.Log(
			"method", "delete_group",
			"urn", urn,
			"took", time.Since(begin),
			"err", err,
		)
	}(time.Now())
	return s.Service.Delete(ctx, urn)
}

func (s *loggingService) Load(ctx context.Context, urn iam.GroupURN) (grp iam.Group, err error) {
	defer func(begin time.Time) {
		s.l.Log(
			"method", "load_group",
			"urn", "urn",
			"took", time.Since(begin),
			err, "err",
		)
	}(time.Now())
	return s.Service.Load(ctx, urn)
}

func (s *loggingService) GetMembers(ctx context.Context, urn iam.GroupURN) (members []iam.UserURN, err error) {
	defer func(begin time.Time) {
		s.l.Log(
			"method", "get_group_members",
			"urn", "urn",
			"took", time.Since(begin),
			"count_members", len(members),
			err, "err",
		)
	}(time.Now())
	return s.Service.GetMembers(ctx, urn)
}

func (s *loggingService) UdpateComment(ctx context.Context, urn iam.GroupURN, comment string) (err error) {
	defer func(begin time.Time) {
		s.l.Log(
			"method", "update_comment",
			"urn", urn,
			"comment", comment,
			"took", time.Since(begin),
			"err", err,
		)
	}(time.Now())
	return s.Service.UpdateComment(ctx, urn, comment)
}

func (s *loggingService) AddMember(ctx context.Context, grp iam.GroupURN, member iam.UserURN) (err error) {
	defer func(begin time.Time) {
		s.l.Log(
			"method", "add_member",
			"grp", grp,
			"member", member,
			"took", time.Since(begin),
			"err", err,
		)
	}(time.Now())
	return s.Service.AddMember(ctx, grp, member)
}

func (s *loggingService) DeleteMember(ctx context.Context, grp iam.GroupURN, member iam.UserURN) (err error) {
	defer func(begin time.Time) {
		s.l.Log(
			"method", "delete_member",
			"grp", grp,
			"member", member,
			"took", time.Since(begin),
			"err", err,
		)
	}(time.Now())
	return s.Service.DeleteMember(ctx, grp, member)
}
