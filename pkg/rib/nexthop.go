package rib

import (
	"encoding/binary"
	"hash/fnv"
	"net/netip"

	"github.com/HT4w5/bgpsim-go/pkg/optional"
)

type NextHopType int

const (
	IP NextHopType = iota
	Interface
	Discard
	Invalid
)

type NextHop struct {
	ip    optional.Optional[netip.Addr]
	iface optional.Optional[string]
	t     NextHopType
}

func MakeNextHop() *NextHop {
	return &NextHop{
		t: Invalid,
	}
}

func (nh *NextHop) WithIP(ip netip.Addr) *NextHop {
	nh.t = IP
	nh.ip = optional.Of(ip)
	nh.iface = optional.Empty[string]()
	return nh
}

func (nh *NextHop) WithInterface(iface string) *NextHop {
	nh.t = Interface
	nh.iface = optional.Of(iface)
	nh.ip = optional.Empty[netip.Addr]()
	return nh
}

func (nh *NextHop) WithDiscard() *NextHop {
	nh.t = Discard
	return nh
}

func (nh *NextHop) Hash() uint32 {
	h := fnv.New32a()
	ipBytes := nh.ip.OrElse(netip.Addr{}).As16()
	h.Write(ipBytes[:])
	h.Write([]byte(nh.iface.OrElse("")))
	binary.Write(h, binary.BigEndian, nh.t)
	return h.Sum32()
}
