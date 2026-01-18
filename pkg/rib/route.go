package rib

import "net/netip"

type Route struct {
	Prefix        netip.Prefix
	NextHop       netip.Addr
	Protocol      Protocol
	AdminCost     int
	Metric        uint64
	NonForwarding bool
}
