package enforcer

import (
	"context"
	"regexp"
	"strings"
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

// ResourceMatcher decided if it matches for the given resource or not.
type ResourceMatcher interface {
	Match(resource string) bool
}

// ResourceMatcherFunc is a convenience wrapper for implementing
// ResourceMatcher using a simple function.
type ResourceMatcherFunc func(r string) bool

// Match implements ResourceMatcher and simply calls rmf
func (rmf ResourceMatcherFunc) Match(r string) bool {
	return rmf(r)
}

// RegexpResourceMatcher implements ResourceMatcher by appliying a
// regular expression against the resource
type RegexpResourceMatcher struct {
	r *regexp.Regexp
}

// NewRegexpResourceMatcher returns a new resource matcher that
// matches r against the resource name.
func NewRegexpResourceMatcher(r string) (*RegexpResourceMatcher, error) {
	reg, err := regexp.CompilePOSIX(r)
	if err != nil {
		return nil, err
	}

	return &RegexpResourceMatcher{reg}, nil
}

// Match implements ResourceMatcher
func (matcher *RegexpResourceMatcher) Match(resource string) bool {
	return matcher.r.Match([]byte(resource))
}

// PrefixMatcher returns a ResourceMatcherFunc that matches on a
// prefix string.
func PrefixMatcher(prefix string) ResourceMatcherFunc {
	return func(resource string) bool {
		return strings.HasPrefix(resource, prefix)
	}
}
