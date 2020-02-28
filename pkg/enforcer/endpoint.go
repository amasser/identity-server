package enforcer

import (
	"context"

	"github.com/go-kit/kit/endpoint"
)

type contextKey string

const (
	// EnforcerActionContextKey is used to store the action that is being performed
	// in a request context.
	EnforcerActionContextKey contextKey = "enforcer:action"

	// EnforcerSubjectContextKey is used to store the subject that is going to perform
	// an action in a request context.
	EnforcerSubjectContextKey contextKey = "enforcer:subject"

	// EnforcerResourceContextKey is used to store the resource that is being operated
	// on in a request context.
	EnforcerResourceContextKey contextKey = "enforcer:resource"

	// EnforcerContextContextKey is used to store additional context values for an
	// operation.
	EnforcerContextContextKey contextKey = "enforcer:context"
)

// WithSubject adds subject to the request context.
func WithSubject(ctx context.Context, subject string) context.Context {
	return context.WithValue(ctx, EnforcerSubjectContextKey, subject)
}

// Subject returns the subject associated with ctx.
func Subject(ctx context.Context) (string, bool) {
	val := ctx.Value(EnforcerSubjectContextKey)
	if val == nil {
		return "", false
	}
	return val.(string), true
}

// WithAction adds action to the request context.
func WithAction(ctx context.Context, action string) context.Context {
	return context.WithValue(ctx, EnforcerActionContextKey, action)
}

// Action returns the action associated with ctx.
func Action(ctx context.Context) (string, bool) {
	val := ctx.Value(EnforcerActionContextKey)
	if val == nil {
		return "", false
	}
	return val.(string), true
}

// WithResource adds resource to the request context.
func WithResource(ctx context.Context, resource string) context.Context {
	return context.WithValue(ctx, EnforcerResourceContextKey, resource)
}

// Resource returns the resource associated with ctx.
func Resource(ctx context.Context) (string, bool) {
	val := ctx.Value(EnforcerResourceContextKey)
	if val == "" {
		return "", false
	}
	return val.(string), true
}

// WithPolicyContext adds values to the request context.
func WithPolicyContext(ctx context.Context, values Context) context.Context {
	return context.WithValue(ctx, EnforcerContextContextKey, values)
}

// PolicyContext returns the policy context associated with ctx.
func PolicyContext(ctx context.Context) (Context, bool) {
	val := ctx.Value(EnforcerContextContextKey)
	if val == nil {
		return nil, false
	}
	return val.(Context), true
}

// Enforce returns an endpoint.Middleware that uses enforcer to ensure
// the request subject is allowed to perform the given action on a resource.
func Enforce(enforcer Enforcer) endpoint.Middleware {
	return func(next endpoint.Endpoint) endpoint.Endpoint {
		return func(ctx context.Context, request interface{}) (interface{}, error) {
			action, ok := Action(ctx)
			if !ok {
				return nil, &PermissionDeniedError{"No action defined"}
			}

			subject, ok := Subject(ctx)
			if !ok {
				return nil, &PermissionDeniedError{"No subject defined"}
			}

			resource, ok := Resource(ctx)
			if !ok {
				return nil, &PermissionDeniedError{"No resource defined"}
			}

			context, _ := PolicyContext(ctx)

			if err := enforcer.Enforce(ctx, subject, action, resource, context); err != nil {
				return nil, err
			}

			return next(ctx, request)
		}
	}
}
