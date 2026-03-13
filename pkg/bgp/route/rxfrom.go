package route

import (
	"cmp"
	"encoding/binary"
	"hash/fnv"
	"net/netip"
)

type RxFromType int

const (
	Local RxFromType = iota
	IP
	Interface
)

type RxFrom struct {
	iface       string
	linkLocalIP netip.Addr
	t           RxFromType
}

func NewRxFrom(opts ...func(*RxFrom)) RxFrom {
	rxf := RxFrom{}
	for _, opt := range opts {
		opt(&rxf)
	}

	return rxf
}

// Option functions

func WithLocal() func(*RxFrom) {
	return func(rxf *RxFrom) {
		rxf.t = Local
	}
}

func WithIP(ip netip.Addr) func(*RxFrom) {
	return func(rxf *RxFrom) {
		rxf.linkLocalIP = ip
		rxf.t = IP
	}
}

func WithInterface(iface string) func(*RxFrom) {
	return func(rxf *RxFrom) {
		rxf.iface = iface
		rxf.t = Interface
	}
}

func WithType(t RxFromType) func(*RxFrom) {
	return func(rxf *RxFrom) {
		rxf.t = t
	}
}

// Getter methods

func (rxf *RxFrom) Hash() uint32 {
	h := fnv.New32a()
	binary.Write(h, binary.BigEndian, rxf.t)

	switch rxf.t {
	case Local:
	case IP:
		ipBytes := rxf.linkLocalIP.As16()
		h.Write(ipBytes[:])
	case Interface:
		h.Write([]byte(rxf.iface))
	}

	return h.Sum32()
}

func (rxf *RxFrom) Iface() string {
	return rxf.iface
}

func (rxf *RxFrom) LinkLocalIP() netip.Addr {
	return rxf.linkLocalIP
}

func (rxf *RxFrom) Type() RxFromType {
	return rxf.t
}

func compareRxFrom(a, b RxFrom) int {
	// Lowest type
	// Local < IP < Interface
	if a.t != b.t {
		return cmp.Compare(a.t, b.t) // Prefer lowest
	}

	switch a.t {
	case Local:
		fallthrough
	case IP:
		return a.linkLocalIP.Compare(b.linkLocalIP)
	case Interface:
		return cmp.Compare(a.iface, b.iface)
	}

	return 0
}
