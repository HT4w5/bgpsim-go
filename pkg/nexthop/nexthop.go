package nexthop

import (
	"encoding/binary"
	"hash/fnv"
	"net/netip"
)

type NextHopType int

const (
	IP NextHopType = iota
	Interface
	Discard
	Invalid
)

type NextHop struct {
	hash  uint32
	ip    netip.Addr
	iface string
	t     NextHopType
}

func New(opts ...func(*NextHop)) NextHop {
	n := NextHop{}
	for _, opt := range opts {
		opt(&n)
	}

	n.computeHash()
	return n
}

func WithIP(ip netip.Addr) func(*NextHop) {
	return func(n *NextHop) {
		n.ip = ip
		n.t = IP
	}
}

func WithInterface(iface string) func(*NextHop) {
	return func(n *NextHop) {
		n.iface = iface
		n.t = Interface
	}
}

func WithDiscard() func(*NextHop) {
	return func(n *NextHop) {
		n.t = Discard
	}
}

func WithType(t NextHopType) func(*NextHop) {
	return func(n *NextHop) {
		n.t = t
	}
}

func (nexthop *NextHop) Hash() uint32 {
	return nexthop.hash
}

func (nexthop *NextHop) IP() netip.Addr {
	return nexthop.ip
}

func (nexthop *NextHop) Iface() string {
	return nexthop.iface
}

func (nexthop *NextHop) Type() NextHopType {
	return nexthop.t
}

func (nexthop *NextHop) computeHash() {
	h := fnv.New32a()
	switch nexthop.t {
	case IP:
		binary.Write(h, binary.BigEndian, nexthop.t)
		ipBytes := nexthop.ip.As16()
		h.Write(ipBytes[:])
	case Interface:
		binary.Write(h, binary.BigEndian, nexthop.t)
		h.Write([]byte(nexthop.iface))
	case Discard:
		binary.Write(h, binary.BigEndian, nexthop.t)
	}

	nexthop.hash = h.Sum32()
}
