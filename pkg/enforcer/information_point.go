package enforcer

import (
	"context"
)

// InfoPoint implements a Policy Information Point (PIP) that is queried
// by Policy Enforcement Points (PEP) whenever additional information about a
// subject or resource - that is part of an authorization request -
// is required.
type InfoPoint interface {
	// GetResourceContext returns additional context for the resource in question.
	// Note that sensitive data should not be exposed as additional context.
	GetResourceContext(ctx context.Context, resource string) (Context, error)
}
