package enforcer

import (
	"context"
	"fmt"
	"net/http"
)

// PermissionDeniedError is returned when a requested action is
// denied.
type PermissionDeniedError struct {
	// Reason holds a human readable reason why the request was denied.
	Reason string `json:"reason"`
}

// Error implements the built-in error interface.
func (pde *PermissionDeniedError) Error() string {
	return fmt.Sprintf("Access denied: %s", pde.Reason)
}

// StatusCode returns http.StatusForbidden and implements
// the StatusCoder interface of go-kit's HTTP transport.
func (pde *PermissionDeniedError) StatusCode() int {
	return http.StatusForbidden
}

// Context represents environmental context of a permission/access
// request.
type Context map[string]interface{}

// Enforcer is responsible of enforcing permission checks.
type Enforcer interface {
	// Enforce checks all available policies and rules and decided whenter subject is
	// allowed to perform action on resource (taking additional environmental context
	// into account). Implementations return nil if the action is allowed.
	// Any error returned should be treated as "permission denied". This ensures that
	// even in case of unreachable remote policy enforcement points (PEP) or policy
	// information points (PIP) no action is allowed by accident.
	Enforce(ctx context.Context, subject, action, resource string, context Context) error
}
