package user

import (
	"context"
	"time"

	"github.com/go-kit/kit/log"
	"github.com/tierklinik-dobersberg/identity-server/iam"
)

type loggingService struct {
	logger log.Logger
	Service
}

// NewLoggingService returns wraps s into a logging service.
func NewLoggingService(logger log.Logger, s Service) Service {
	return &loggingService{
		logger:  logger,
		Service: s,
	}
}

func (s *loggingService) CreateUser(ctx context.Context, username, password string, attrs map[string]interface{}) (id iam.UserURN, err error) {
	defer func(begin time.Time) {
		s.logger.Log(
			"method", "create_user",
			"username", username,
			"took", time.Since(begin),
			"err", err,
		)
	}(time.Now())
	return s.Service.CreateUser(ctx, username, password, attrs)
}

func (s *loggingService) LoadUser(ctx context.Context, urn iam.UserURN) (user iam.User, err error) {
	defer func(begin time.Time) {
		s.logger.Log(
			"method", "load_user",
			"urn", urn,
			"took", time.Since(begin),
			"err", err,
		)
	}(time.Now())
	return s.Service.LoadUser(ctx, urn)
}

func (s *loggingService) DeleteUser(ctx context.Context, urn iam.UserURN) (err error) {
	defer func(begin time.Time) {
		s.logger.Log(
			"method", "delete_user",
			"urn", urn,
			"took", time.Since(begin),
			"err", err,
		)
	}(time.Now())
	return s.Service.DeleteUser(ctx, urn)
}

func (s *loggingService) LockUser(ctx context.Context, urn iam.UserURN, locked bool) (err error) {
	defer func(begin time.Time) {
		s.logger.Log(
			"method", "delete_user",
			"urn", urn,
			"locked", locked,
			"took", time.Since(begin),
			"err", err,
		)
	}(time.Now())
	return s.Service.LockUser(ctx, urn, locked)
}

func (s *loggingService) Users(ctx context.Context) (users []iam.User, err error) {
	defer func(begin time.Time) {
		s.logger.Log(
			"method", "list_users",
			"took", time.Since(begin),
			"err", err,
		)
	}(time.Now())
	return s.Service.Users(ctx)
}

func (s *loggingService) UpdateAttrs(ctx context.Context, urn iam.UserURN, attrs map[string]interface{}) (err error) {
	defer func(begin time.Time) {
		s.logger.Log(
			"method", "update_attrs",
			"urn", urn,
			"took", time.Since(begin),
			"err", err,
		)
	}(time.Now())
	return s.Service.UpdateAttrs(ctx, urn, attrs)
}

func (s *loggingService) SetAttr(ctx context.Context, urn iam.UserURN, key string, value interface{}) (err error) {
	defer func(begin time.Time) {
		s.logger.Log(
			"method", "set_attr",
			"urn", urn,
			"attr_key", key,
			"took", time.Since(begin),
			"err", err,
		)
	}(time.Now())
	return s.Service.SetAttr(ctx, urn, key, value)
}

func (s *loggingService) DeleteAttr(ctx context.Context, urn iam.UserURN, key string) (err error) {
	defer func(begin time.Time) {
		s.logger.Log(
			"method", "delete_attr",
			"urn", urn,
			"attr_key", key,
			"took", time.Since(begin),
			"err", err,
		)
	}(time.Now())
	return s.Service.DeleteAttr(ctx, urn, key)
}

func (s *loggingService) OnDelete(ctx context.Context, fn OnDeleteFunc) {
	defer func(begin time.Time) {
		s.logger.Log(
			"method", "on_delete",
			"took", time.Since(begin),
		)
	}(time.Now())

	wrapped := func(urn iam.UserURN) {
		defer func(begin time.Time) {
			s.logger.Log(
				"method", "on_delete_callback",
				"took", time.Since(begin),
				"urn", urn,
			)
		}(time.Now())

		fn(urn)
	}

	s.Service.OnDelete(ctx, wrapped)
}
