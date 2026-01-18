package rib

import (
	"encoding/binary"
	"hash/fnv"
	"net/netip"
)

type Protocol int

const (
	BGP Protocol = iota
)

type Route struct {
	Prefix        netip.Prefix
	NextHop       netip.Addr
	Protocol      Protocol
	AdminCost     int
	Metric        uint64
	NonForwarding bool
}

func (r Route) Hash() uint32 {
	h := fnv.New32a()

	prefixBytes, _ := r.Prefix.MarshalBinary()
	h.Write(prefixBytes)

	nextHopBytes, _ := r.NextHop.MarshalBinary()
	h.Write(nextHopBytes)

	binary.Write(h, binary.BigEndian, r.Protocol)
	binary.Write(h, binary.BigEndian, r.AdminCost)
	binary.Write(h, binary.BigEndian, r.Metric)
	binary.Write(h, binary.BigEndian, r.NonForwarding)

	return h.Sum32()
}
