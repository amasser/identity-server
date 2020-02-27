package policy

import (
	"context"
	"strings"
	"time"

	"github.com/go-kit/kit/log"
	"github.com/tierklinik-dobersberg/identity-server/iam"
)

type loggingService struct {
	Service
	l log.Logger
}

// NewLoggingService returns a new service that logs every request to
// the logging service.
func NewLoggingService(l log.Logger, s Service) Service {
	return &loggingService{
		Service: s,
		l:       l,
	}
}

func (l *loggingService) Create(ctx context.Context, name string, policy iam.Policy) (urn iam.PolicyURN, err error) {
	defer func(begin time.Time) {
		l.l.Log(
			"method", "create_policy",
			"name", name,
			"subjects", strings.Join(policy.Subjects, ","),
			"resources", strings.Join(policy.Resources, ", "),
			"took", time.Since(begin),
			"urn", urn,
			"err", err,
		)
	}(time.Now())

	return l.Service.Create(ctx, name, policy)
}

func (l *loggingService) Delete(ctx context.Context, urn iam.PolicyURN) (err error) {
	defer func(begin time.Time) {
		l.l.Log(
			"method", "delete_policy",
			"took", time.Since(begin),
			"urn", urn,
			"err", err,
		)
	}(time.Now())

	return l.Service.Delete(ctx, urn)
}

func (l *loggingService) Update(ctx context.Context, urn iam.PolicyURN, p iam.Policy) (err error) {
	defer func(begin time.Time) {
		l.l.Log(
			"method", "update_policy",
			"took", time.Since(begin),
			"urn", urn,
			"subjects", strings.Join(p.Subjects, ", "),
			"resources", strings.Join(p.Resources, ", "),
			"err", err,
		)
	}(time.Now())

	return l.Service.Update(ctx, urn, p)
}

func (l *loggingService) List(ctx context.Context) (policies []iam.Policy, err error) {
	defer func(begin time.Time) {
		l.l.Log(
			"method", "list_policies",
			"took", time.Since(begin),
			"policies", len(policies),
			"err", err,
		)
	}(time.Now())

	return l.Service.List(ctx)
}
