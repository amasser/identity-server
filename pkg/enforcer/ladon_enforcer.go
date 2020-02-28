package enforcer

import (
	"context"
	"fmt"

	"github.com/ory/ladon"
)

// LadonEnforcer implements the Enforcer interface based on
// awesome ory/ladon package.
type LadonEnforcer struct {
	infoPoint InfoPoint

	warden ladon.Warden
}

// NewLadonEnforcer returns a new ory/ladon based enforcer.
func NewLadonEnforcer(manager ladon.Manager, infoPoint InfoPoint) *LadonEnforcer {
	return &LadonEnforcer{
		infoPoint: infoPoint,
		warden: &ladon.Ladon{
			Manager: manager,
		},
	}
}

// Enforce checks if subject is allowed to perform action on resource. It implements the Enforcer interface.
func (e *LadonEnforcer) Enforce(ctx context.Context, subject, action, resource string, context Context) error {
	// TODO(ppacher): get subject and resource context in parallel.

	subjectContext, err := e.infoPoint.GetResourceContext(ctx, subject)
	if err != nil {
		return nil
	}

	resourceContext, err := e.infoPoint.GetResourceContext(ctx, resource)
	if err != nil {
		return nil
	}

	resultCtx := make(Context)

	for k, v := range context {
		resultCtx[k] = v
	}

	for k, v := range subjectContext {
		if _, ok := resultCtx[k]; ok {
			return fmt.Errorf("subject-context duplicates context key %q", k)
		}
		resultCtx[k] = v
	}

	for k, v := range resourceContext {
		if _, ok := resultCtx[k]; ok {
			return fmt.Errorf("resource-context duplicates context key %q", k)
		}
		resultCtx[k] = v
	}

	request := &ladon.Request{
		Action:   action,
		Subject:  subject,
		Resource: resource,
		Context:  ladon.Context(resultCtx),
	}

	if err := e.warden.IsAllowed(request); err != nil {
		return &PermissionDeniedError{Reason: err.Error()}
	}

	return nil
}
