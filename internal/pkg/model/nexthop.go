package model

import (
	"net"
)

type NextHopType int

const (
	IP NextHopType = iota
	INTERFACE
	DISCARD
	INVALID
)

type NextHop struct {
	ip    net.IP
	iface *string
	t     NextHopType
}

func NewInvalidNextHop() *NextHop {
	return &NextHop{
		ip:    nil,
		iface: nil,
		t:     INVALID,
	}
}

func NewNextHop(ip net.IP, iface string) *NextHop {
	if ip == nil && len(iface) == 0 {
		return nil
	}
	if ip == nil {
		return &NextHop{
			ip:    nil,
			iface: &iface,
			t:     INTERFACE,
		}
	} else {
		return &NextHop{
			ip:    ip,
			iface: nil,
			t:     IP,
		}
	}
}
