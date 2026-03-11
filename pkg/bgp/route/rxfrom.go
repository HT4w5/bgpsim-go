package route

import (
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
	hash        uint32
	iface       string
	linkLocalIP netip.Addr
	t           RxFromType
}

func NewRxFrom(opts ...func(*RxFrom)) RxFrom {
	rxf := RxFrom{}
	for _, opt := range opts {
		opt(&rxf)
	}

	rxf.computeHash()
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
	return rxf.hash
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

func (rxf *RxFrom) computeHash() {
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

	rxf.hash = h.Sum32()
}
