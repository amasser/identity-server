package enforcer

import "context"

// NoOpEnforcer implements the Enforcer interface but allows every operation.
type NoOpEnforcer struct{}

// NewNoOpEnforcer returns a new NoOpEnforcer.
func NewNoOpEnforcer() *NoOpEnforcer {
	return &NoOpEnforcer{}
}

// Enforce implements the Enforcer interface but does nothing.
func (NoOpEnforcer) Enforce(ctx context.Context, subject, action, resource string, context Context) error {
	return nil
}
