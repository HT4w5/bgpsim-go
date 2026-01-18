package rpool

import "iter"

type Route interface {
	Hash() uint32 // Stable, non-colliding hash
}

type RoutePool struct {
	routes []Route
	idxMap map[uint32]int // route.Hash() -> index
	len    int
}

func MakeRoutePool() *RoutePool {
	return &RoutePool{
		routes: make([]Route, 0),
		idxMap: make(map[uint32]int),
		len:    0,
	}
}

// Insert a route, return true if changed
func (rp *RoutePool) Insert(route Route) bool {
	h := route.Hash()
	_, ok := rp.idxMap[h]
	if !ok {
		rp.idxMap[h] = len(rp.idxMap)
		rp.routes = append(rp.routes, route)
		rp.len++
		rp.rehash()
		return true
	}
	return false
}

// Remove a route, return true if changed
func (rp *RoutePool) Remove(route Route) bool {
	h := route.Hash()
	_, ok := rp.idxMap[h]
	if !ok {
		return false
	}
	delete(rp.idxMap, h)
	rp.len--
	rp.rehash()
	return true
}

func (rp *RoutePool) All() iter.Seq[Route] {
	return func(yield func(Route) bool) {
		for _, i := range rp.idxMap {
			if !yield(rp.routes[i]) {
				return
			}
		}
	}
}

func (rp *RoutePool) Len() int {
	return rp.len
}

func (rp *RoutePool) rehash() {
	// Rehash on condition
	if len(rp.routes)/rp.len < 2 {
		return
	}

	newRoutes := make([]Route, 0, cap(rp.routes))
	newIdxMap := make(map[uint32]int)

	for k, v := range rp.idxMap {
		newIdxMap[k] = len(newRoutes)
		newRoutes = append(newRoutes, rp.routes[v])
	}

	rp.routes = newRoutes
	rp.idxMap = newIdxMap
}
