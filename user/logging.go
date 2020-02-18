package user

import (
	"context"
	"time"

	"github.com/go-kit/kit/log"
	"github.com/tierklinik-dobersberg/iam/v2/iam"
)

type loggingService struct {
	logger log.Logger
	Service
}

// NewLoggerService returns wraps s into a logging service.
func NewLoggerService(logger log.Logger, s Service) Service {
	return &loggingService{
		logger:  logger,
		Service: s,
	}
}

func (s *loggingService) CreateUser(ctx context.Context, accountID int, username string, attrs map[string]interface{}) (id iam.UserURN, err error) {
	defer func(begin time.Time) {
		s.logger.Log(
			"method", "create_user",
			"accountID", accountID,
			"username", username,
			"took", time.Since(begin),
			"err", err,
		)
	}(time.Now())
	return s.Service.CreateUser(ctx, accountID, username, attrs)
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
