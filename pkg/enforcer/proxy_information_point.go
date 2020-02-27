package enforcer

import (
	"context"
	"sync"

	"github.com/tierklinik-dobersberg/identity-server/pkg/common"
)

// ProxyInfoPoint proxies policy information requests
// to one or more Policy Information Points (PIP) based
// on resource matching.
type ProxyInfoPoint struct {
	rw     sync.RWMutex
	points map[ResourceMatcher]InfoPoint
}

// NewProxyInfoPoint returns a new proxy information point.
func NewProxyInfoPoint() *ProxyInfoPoint {
	return &ProxyInfoPoint{
		points: make(map[ResourceMatcher]InfoPoint),
	}
}

// Add adds a new Policy Information Point (PIP) to be used
// when matcher applies.
func (pip *ProxyInfoPoint) Add(matcher ResourceMatcher, point InfoPoint) {
	pip.rw.Lock()
	defer pip.rw.Unlock()

	pip.points[matcher] = point
}

// GetResourceContext implements InfoPoint and calls the matching
// info point. Note that ProxyInfoPoint does not handle multiple matches.
func (pip *ProxyInfoPoint) GetResourceContext(ctx context.Context, resource string) (Context, error) {
	pip.rw.RLock()
	defer pip.rw.RUnlock()

	for matcher, infoPoint := range pip.points {
		if matcher.Match(resource) {
			return infoPoint.GetResourceContext(ctx, resource)
		}
	}

	return nil, common.NewNotFoundError("resource type")
}
