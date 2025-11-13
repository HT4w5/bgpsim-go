package model

import (
	"fmt"
	"net"
	"net/netip"
)

type StaticRouteType int

const (
	SR_LOCAL StaticRouteType = iota
	SR_CONNECTED
	SR_BGP
	SR_STATIC
)

var staticRouteTypeMap = map[string]StaticRouteType{
	"LOCAL":     SR_LOCAL,
	"CONNECTED": SR_CONNECTED,
	"BGP":       SR_BGP,
	"STATIC":    SR_STATIC,
}

type StaticRoute struct {
	prefix             netip.Prefix
	nextHop            NextHop
	administrativeCost int
	tag                int
	weight             int
	metric             int
	t                  StaticRouteType
}

type StaticRouteConfig struct {
	Network            string `json:"network"`
	NextHopIp          string `json:"nextHopIp"`
	NextHopInterface   string `json:"nextHopInterface"`
	AdministrativeCost int    `json:"administrativeCost"`
	Tag                int    `json:"tag"`
	Weight             int    `json:"weight"`
	Metric             int    `json:"metric"`
	Type               string `json:"type"`
}

func NewStaticRoute(cfg *StaticRouteConfig) (*StaticRoute, error) {
	sr := &StaticRoute{
		administrativeCost: cfg.AdministrativeCost,
		tag:                cfg.Tag,
		weight:             cfg.Weight,
		metric:             cfg.Metric,
	}

	// Determine static route type
	if t, ok := staticRouteTypeMap[cfg.Type]; !ok {
		return nil, fmt.Errorf("invalid type: %s", cfg.Type)
	} else {
		sr.t = t
	}

	// Parse prefix
	if prefix, err := netip.ParsePrefix(cfg.Network); err != nil {
		return nil, fmt.Errorf("invalid prefix %s: %w", cfg.Network, err)
	} else {
		sr.prefix = prefix
	}

	// Parse next hop
	if nextHop := NewNextHop(net.ParseIP(cfg.NextHopIp), cfg.NextHopInterface); nextHop == nil {
		return nil, fmt.Errorf("invalid next hop: %s %s", cfg.NextHopIp, cfg.NextHopInterface)
	} else {
		sr.nextHop = *nextHop
	}

	return sr, nil
}
