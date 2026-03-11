package route

import (
	"encoding/binary"
	"hash/fnv"
	"net/netip"

	"github.com/HT4w5/bgpsim-go/pkg/nexthop"
)

// BgpRoute represents a BGP route
// Inmutable once created
type BgpRoute struct {
	adminCost       int
	asPath          []uint32
	bgpAdminCost    int
	hash            uint32
	localPreference uint32
	metric          int
	nextHop         nexthop.NextHop
	nonForwarding   bool
	nonRouting      bool
	pathID          int
	prefix          netip.Prefix
	receivedFrom    RxFrom
	srcPrefixLength int
	tag             int
	weight          int
}

// Create new BgpRoute
func New(opts ...func(*BgpRoute)) *BgpRoute {
	r := &BgpRoute{}
	for _, opt := range opts {
		opt(r)
	}
	r.computeHash()
	return r
}

func (r *BgpRoute) Clone(opts ...func(*BgpRoute)) *BgpRoute {
	clone := &BgpRoute{
		adminCost:       r.adminCost,
		asPath:          make([]uint32, len(r.asPath)),
		bgpAdminCost:    r.bgpAdminCost,
		hash:            r.hash,
		localPreference: r.localPreference,
		metric:          r.metric,
		nextHop:         r.nextHop,
		nonForwarding:   r.nonForwarding,
		nonRouting:      r.nonRouting,
		pathID:          r.pathID,
		prefix:          r.prefix,
		receivedFrom:    r.receivedFrom,
		srcPrefixLength: r.srcPrefixLength,
		tag:             r.tag,
		weight:          r.weight,
	}

	// Copy AsPath
	copy(clone.asPath, r.asPath)

	// Apply options
	for _, opt := range opts {
		opt(clone)
	}

	// Recompute hash
	clone.computeHash()
	return clone
}

func (r *BgpRoute) computeHash() {
	h := fnv.New32a()

	binary.Write(h, binary.BigEndian, r.adminCost)
	binary.Write(h, binary.BigEndian, r.asPath)
	binary.Write(h, binary.BigEndian, r.bgpAdminCost)
	binary.Write(h, binary.BigEndian, r.localPreference)
	binary.Write(h, binary.BigEndian, r.metric)
	binary.Write(h, binary.BigEndian, r.nextHop.Hash())
	binary.Write(h, binary.BigEndian, r.nonForwarding)
	binary.Write(h, binary.BigEndian, r.nonRouting)
	binary.Write(h, binary.BigEndian, r.pathID)
	prefixBytes, _ := r.prefix.MarshalBinary()
	h.Write(prefixBytes)
	binary.Write(h, binary.BigEndian, r.receivedFrom.Hash())
	binary.Write(h, binary.BigEndian, r.srcPrefixLength)
	binary.Write(h, binary.BigEndian, r.tag)
	binary.Write(h, binary.BigEndian, r.weight)

	r.hash = h.Sum32()
}

// Getter methods

func (r *BgpRoute) AdminCost() int {
	return r.adminCost
}

func (r *BgpRoute) AsPath() []uint32 {
	return r.asPath
}

func (r *BgpRoute) BgpAdminCost() int {
	return r.bgpAdminCost
}

func (r *BgpRoute) Hash() uint32 {
	return r.hash
}

func (r *BgpRoute) LocalPreference() uint32 {
	return r.localPreference
}

func (r *BgpRoute) Metric() int {
	return r.metric
}

func (r *BgpRoute) NextHop() nexthop.NextHop {
	return r.nextHop
}

func (r *BgpRoute) NonForwarding() bool {
	return r.nonForwarding
}

func (r *BgpRoute) NonRouting() bool {
	return r.nonRouting
}

func (r *BgpRoute) PathID() int {
	return r.pathID
}

func (r *BgpRoute) Prefix() netip.Prefix {
	return r.prefix
}

func (r *BgpRoute) ReceivedFrom() RxFrom {
	return r.receivedFrom
}

func (r *BgpRoute) SrcPrefixLength() int {
	return r.srcPrefixLength
}

func (r *BgpRoute) Tag() int {
	return r.tag
}

func (r *BgpRoute) Weight() int {
	return r.weight
}

// Options

func WithAdminCost(adminCost int) func(*BgpRoute) {
	return func(r *BgpRoute) {
		r.adminCost = adminCost
	}
}

func WithAsPath(asPath []uint32) func(*BgpRoute) {
	return func(r *BgpRoute) {
		r.asPath = asPath
	}
}

func WithBgpAdminCost(bgpAdminCost int) func(*BgpRoute) {
	return func(r *BgpRoute) {
		r.bgpAdminCost = bgpAdminCost
	}
}

func WithLocalPreference(localPreference uint32) func(*BgpRoute) {
	return func(r *BgpRoute) {
		r.localPreference = localPreference
	}
}

func WithMetric(metric int) func(*BgpRoute) {
	return func(r *BgpRoute) {
		r.metric = metric
	}
}

func WithNextHop(nextHop nexthop.NextHop) func(*BgpRoute) {
	return func(r *BgpRoute) {
		r.nextHop = nextHop
	}
}

func WithNonForwarding(nonForwarding bool) func(*BgpRoute) {
	return func(r *BgpRoute) {
		r.nonForwarding = nonForwarding
	}
}

func WithNonRouting(nonRouting bool) func(*BgpRoute) {
	return func(r *BgpRoute) {
		r.nonRouting = nonRouting
	}
}

func WithPathID(pathID int) func(*BgpRoute) {
	return func(r *BgpRoute) {
		r.pathID = pathID
	}
}

func WithPrefix(prefix netip.Prefix) func(*BgpRoute) {
	return func(r *BgpRoute) {
		r.prefix = prefix
	}
}

func WithReceivedFrom(receivedFrom RxFrom) func(*BgpRoute) {
	return func(r *BgpRoute) {
		r.receivedFrom = receivedFrom
	}
}

func WithSrcPrefixLength(srcPrefixLength int) func(*BgpRoute) {
	return func(r *BgpRoute) {
		r.srcPrefixLength = srcPrefixLength
	}
}

func WithTag(tag int) func(*BgpRoute) {
	return func(r *BgpRoute) {
		r.tag = tag
	}
}

func WithWeight(weight int) func(*BgpRoute) {
	return func(r *BgpRoute) {
		r.weight = weight
	}
}
