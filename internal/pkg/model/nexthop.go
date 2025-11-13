package model

import (
	"net"
)

type NextHopType int

const (
	NH_IP NextHopType = iota
	NH_INTERFACE
	NH_DISCARD
	NH_INVALID
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
		t:     NH_INVALID,
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
			t:     NH_INTERFACE,
		}
	} else {
		return &NextHop{
			ip:    ip,
			iface: nil,
			t:     NH_IP,
		}
	}
}
