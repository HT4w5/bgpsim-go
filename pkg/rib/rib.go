package rib

import (
	"net/netip"

	"github.com/HT4w5/bgpsim-go/pkg/rpool"
	"github.com/gaissmai/bart"
)

type Rib struct {
	ipt *bart.Table[*rpool.RoutePool]
}

func MakeRib() *Rib {
	return &Rib{
		ipt: &bart.Table[*rpool.RoutePool]{},
	}
}

// AddRoute adds a route to the RIB
func (rib *Rib) AddRoute(route Route) bool {
	changed := false
	p, ok := rib.ipt.LookupPrefix(route.Prefix)
	if !ok {
		changed = true
		p = rpool.MakeRoutePool()
		rib.ipt.Insert(route.Prefix, p)
	}
	changed = p.Insert(route) || changed
	return changed
}

// RemoveRoute removes a route from the RIB
func (rib *Rib) RemoveRoute(route Route) bool {
	p, ok := rib.ipt.LookupPrefix(route.Prefix)
	if !ok {
		return false
	}
	return p.Remove(route)
}

// GetRoutes returns all routes in the RIB
func (rib *Rib) GetRoutes() []Route {
	routes := make([]Route, rib.ipt.Size())
	for _, m := range rib.ipt.All() {
		for r := range m.All() {
			routes = append(routes, r.(Route))
		}
	}
	return routes
}

// GetRoutesForPrefix returns routes for a specific prefix
func (rib *Rib) GetRoutesForPrefix(prefix netip.Prefix) []Route {
	routes := make([]Route, 0)
	for _, m := range rib.ipt.Subnets(prefix) {
		for r := range m.All() {
			routes = append(routes, r.(Route))
		}
	}
	return routes
}

// LongestPrefixMatch finds the most specific routes matching an IP
func (rib *Rib) LongestPrefixMatch(ip netip.Addr) []Route {
	routes := make([]Route, 0)
	m, ok := rib.ipt.Lookup(ip)
	if !ok {
		return routes
	}
	for r := range m.All() {
		routes = append(routes, r.(Route))
	}
	return routes
}
