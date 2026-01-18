package rib

import (
	"net/netip"

	"github.com/gaissmai/bart"
)

type Protocol int

const (
	BGP Protocol = iota
)

type Rib struct {
	ipt *bart.Table[map[Route]struct{}]
}

func MakeRib() *Rib {
	return &Rib{
		ipt: &bart.Table[map[Route]struct{}]{},
	}
}

// AddRoute adds a route to the RIB
func (rib *Rib) AddRoute(route Route) {
	m, ok := rib.ipt.LookupPrefix(route.Prefix)
	if !ok {
		m = make(map[Route]struct{})
		rib.ipt.Insert(route.Prefix, m)
	}
	m[route] = struct{}{}
}

// RemoveRoute removes a route from the RIB
func (rib *Rib) RemoveRoute(route Route) {
	m, ok := rib.ipt.LookupPrefix(route.Prefix)
	if !ok {
		return
	}
	delete(m, route)
}

// GetRoutes returns all routes in the RIB
func (rib *Rib) GetRoutes() []Route {
	routes := make([]Route, rib.ipt.Size())
	for _, m := range rib.ipt.All() {
		for r, _ := range m {
			routes = append(routes, r)
		}
	}
	return routes
}

// GetRoutesForPrefix returns routes for a specific prefix
func (rib *Rib) GetRoutesForPrefix(prefix netip.Prefix) []Route {
	routes := make([]Route, 0)
	for _, m := range rib.ipt.Subnets(prefix) {
		for r, _ := range m {
			routes = append(routes, r)
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
	for r, _ := range m {
		routes = append(routes, r)
	}
	return routes
}
