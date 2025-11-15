package model

import (
	"net"
	"net/netip"
	"time"
)

type BgpOriginType int

const (
	BO_LOCAL BgpOriginType = iota
	BO_IP
	BO_INTERFACE
)

type BgpRoute struct {
	attrs   BgpAttrs
	prefix  netip.Prefix
	nextHop NextHop
	origin  BgpOrigin
	arrival time.Time
}

type BgpAttrs struct {
	LocalPref int
	Weight    int
	Tag       int
	Metric    int
	AD        int
	AsPath    *AsPath
}

type BgpOrigin struct {
	Iface       string
	LinkLocalIp net.IP
	T           BgpOriginType
}
