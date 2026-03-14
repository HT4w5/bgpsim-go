package rset

import (
	"maps"

	"github.com/HT4w5/bgpsim-go/pkg/bgp/route"
)

const (
	maxPaths = 8
)

type RouteSet struct {
	bestPath     *route.BgpRoute
	multipathSet map[uint32]*route.BgpRoute

	// Config
	multipath bool
}

// Build a RouteSet from multiple BgpRoutes
// Assume provided routes have the same prefix (and len != 0)
func New(routes []*route.BgpRoute, opts ...func(*RouteSet)) *RouteSet {
	rs := &RouteSet{}

	for _, opt := range opts {
		opt(rs)
	}

	if !rs.multipath {
		// Get best
		compFunc := func(a, b *route.BgpRoute) int {
			comp := route.CompareMultipath(a, b)
			if comp == 0 {
				comp = route.CompareTieBreak(a, b)
			}
			return comp
		}
		best := routes[0] // Panics if len(routes) == 0

		for _, r := range routes {
			if compFunc(best, r) > 0 {
				best = r
			}
		}

		rs.bestPath = best
	} else {
		// Get best multipath

		best := routes[0] // Panics if len(routes) == 0

		for _, r := range routes {
			if route.CompareMultipath(best, r) > 0 {
				best = r
			}
		}

		rs.multipathSet = make(map[uint32]*route.BgpRoute)

		multipathCount := 0 // Track multipath count
		for _, r := range routes {
			if route.CompareMultipath(best, r) == 0 {
				rs.multipathSet[r.Hash()] = r
				multipathCount++
			}
			if multipathCount >= maxPaths {
				break
			}
		}

		for _, best = range rs.multipathSet { // Get first multipath route
			break
		}

		for _, r := range rs.multipathSet {
			if route.CompareTieBreak(best, r) > 0 {
				best = r
			}
		}

		rs.bestPath = best
	}

	return rs
}

// Options

func WithMultipath(m bool) func(*RouteSet) {
	return func(rs *RouteSet) {
		rs.multipath = m
	}
}

// Compare equality of two RouteSets
func (rs *RouteSet) Eq(other *RouteSet) bool {
	if rs.multipath != other.multipath {
		return false // Shouldn't happen in reality since multipath is explicitly set to the same on the same process
	}

	if !rs.multipath {
		return rs.bestPath.Hash() == other.bestPath.Hash()
	} else {
		equalMap := maps.EqualFunc(rs.multipathSet, other.multipathSet, func(a, b *route.BgpRoute) bool { return true })
		if !equalMap {
			return false
		}
		return rs.bestPath.Hash() == other.bestPath.Hash()
	}
}

// Getters

func (rs *RouteSet) BestPath() *route.BgpRoute {
	return rs.bestPath
}

func (rs *RouteSet) MultipathSet() []*route.BgpRoute {
	set := make([]*route.BgpRoute, 0, len(rs.multipathSet))
	for _, r := range rs.multipathSet {
		set = append(set, r)
	}
	return set
}
