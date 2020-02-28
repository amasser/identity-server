package enforcer

import (
	"context"

	"github.com/go-kit/kit/endpoint"
)

type contextKey string

const (
	// ContextKeyAction is used to store the action that is being performed
	// in a request context.
	ContextKeyAction contextKey = "enforcer:action"

	// ContextKeySubject is used to store the subject that is going to perform
	// an action in a request context.
	ContextKeySubject contextKey = "enforcer:subject"

	// ContextKeyResource is used to store the resource that is being operated
	// on in a request context.
	ContextKeyResource contextKey = "enforcer:resource"

	// ContextKeyContext is used to store additional context values for an
	// operation.
	ContextKeyContext contextKey = "enforcer:context"
)

// WithSubject adds subject to the request context.
func WithSubject(ctx context.Context, subject string) context.Context {
	return context.WithValue(ctx, ContextKeySubject, subject)
}

// Subject returns the subject associated with ctx.
func Subject(ctx context.Context) (string, bool) {
	val := ctx.Value(ContextKeySubject)
	if val == nil {
		return "", false
	}
	return val.(string), true
}

// WithAction adds action to the request context.
func WithAction(ctx context.Context, action string) context.Context {
	return context.WithValue(ctx, ContextKeyAction, action)
}

// Action returns the action associated with ctx.
func Action(ctx context.Context) (string, bool) {
	val := ctx.Value(ContextKeyAction)
	if val == nil {
		return "", false
	}
	return val.(string), true
}

// WithResource adds resource to the request context.
func WithResource(ctx context.Context, resource string) context.Context {
	return context.WithValue(ctx, ContextKeyResource, resource)
}

// Resource returns the resource associated with ctx.
func Resource(ctx context.Context) (string, bool) {
	val := ctx.Value(ContextKeyResource)
	if val == nil {
		return "", false
	}
	return val.(string), true
}

// WithPolicyContext adds values to the request context.
func WithPolicyContext(ctx context.Context, values Context) context.Context {
	return context.WithValue(ctx, ContextKeyContext, values)
}

// PolicyContext returns the policy context associated with ctx.
func PolicyContext(ctx context.Context) (Context, bool) {
	val := ctx.Value(ContextKeyContext)
	if val == nil {
		return nil, false
	}
	return val.(Context), true
}

// NewActionEndpoint adds the specified endpoint the each request context of
// the wrapped endpoint.
func NewActionEndpoint(action string) endpoint.Middleware {
	return func(next endpoint.Endpoint) endpoint.Endpoint {
		return func(ctx context.Context, request interface{}) (interface{}, error) {
			ctx = WithAction(ctx, action)
			return next(ctx, request)
		}
	}
}

// NewResourceEndpoint returns an endpoint.Middleware that extracts a the resource name from a request
// and adds it to the request context of the wrapped endpoint.
func NewResourceEndpoint(resourcer func(ctx context.Context, request interface{}) (string, error)) endpoint.Middleware {
	return func(next endpoint.Endpoint) endpoint.Endpoint {
		return func(ctx context.Context, request interface{}) (interface{}, error) {
			resource, err := resourcer(ctx, request)
			if err != nil {
				return nil, err
			}

			ctx = WithResource(ctx, resource)
			return next(ctx, request)
		}
	}
}

// NewEnforcedEndpoint returns an endpoint.Middleware that uses enforcer to ensure
// the request subject is allowed to perform the given action on a resource.
func NewEnforcedEndpoint(enforcer Enforcer) endpoint.Middleware {
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

			resource, _ := Resource(ctx)
			/*
				if !ok {
					return nil, &PermissionDeniedError{"No resource defined"}
				}
			*/

			context, _ := PolicyContext(ctx)

			if err := enforcer.Enforce(ctx, subject, action, resource, context); err != nil {
				return nil, err
			}

			return next(ctx, request)
		}
	}
}
