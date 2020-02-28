package enforcer

import (
	"context"
	"regexp"
	"strings"
	"sync"
)

// MultiResourceInfoPoint proxies policy information requests
// to one or more Policy Information Points (PIP) based
// on resource matching.
type MultiResourceInfoPoint struct {
	rw     sync.RWMutex
	points map[ResourceMatcher]InfoPoint
}

// NewMultiResourceInfoPoint returns a new proxy policy information point.
func NewMultiResourceInfoPoint() *MultiResourceInfoPoint {
	return &MultiResourceInfoPoint{
		points: make(map[ResourceMatcher]InfoPoint),
	}
}

// Add adds a new Policy Information Point (PIP) to be used
// when matcher applies.
func (pip *MultiResourceInfoPoint) Add(matcher ResourceMatcher, point InfoPoint) {
	pip.rw.Lock()
	defer pip.rw.Unlock()

	pip.points[matcher] = point
}

// GetResourceContext implements InfoPoint and calls the matching
// info point. Note that MultiResourceInfoPoint does not handle multiple matches.
func (pip *MultiResourceInfoPoint) GetResourceContext(ctx context.Context, resource string) (Context, error) {
	pip.rw.RLock()
	defer pip.rw.RUnlock()

	for matcher, infoPoint := range pip.points {
		if matcher.Match(resource) {
			return infoPoint.GetResourceContext(ctx, resource)
		}
	}

	return nil, nil
}

// ResourceMatcher decided if it matches for the given resource or not.
type ResourceMatcher interface {
	Match(resource string) bool
}

// ResourceMatcherFunc is a convenience wrapper for implementing
// ResourceMatcher using a simple function.
type ResourceMatcherFunc func(r string) bool

// Match implements ResourceMatcher and simply calls rmf.
func (rmf ResourceMatcherFunc) Match(r string) bool {
	return rmf(r)
}

// RegexpResourceMatcher implements ResourceMatcher by appliying a
// regular expression against the resource.
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

// Match implements ResourceMatcher.
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
