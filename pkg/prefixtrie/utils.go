package prefixtrie

import (
	"math/bits"
	"net/netip"
)

// Set left len bits to 1 in uint32
func uint32PrefixMask(len uint8) uint32 {
	return ^uint32(0) << (32 - len)
}

// Get bit at idx (starting from 0) in seq
func uint32GetBit(seq uint32, idx uint8) bool {
	mask := uint32(1 << (31 - idx))
	return (seq & mask) != 0
}

// Get longest common prefix length of two uint32 values
func uint32LCPL(a uint32, b uint32) int {
	diff := a ^ b
	return bits.LeadingZeros32(diff)
}

// Compare two uint8s
// Return value of a - b (can be negative)
func uint8Cmp(a uint8, b uint8) int {
	if a == b {
		return 0
	} else if a > b {
		return int(a - b)
	} else {
		return -int(b - a)
	}
}

// Return value type for prefix trie query
type Match[V any] struct {
	prefix netip.Prefix
	value  *V
	found  bool
}

func (m Match[V]) GetValue() V {
	if !m.found {
		var zero V
		return zero
	}
	return *m.value
}

func (m Match[V]) GetPrefix() netip.Prefix {
	if !m.found {
		return netip.Prefix{}
	}
	return m.prefix
}

func (m Match[V]) Found() bool {
	return m.found
}

func NewFoundMatch[V any](prefix netip.Prefix, value *V) Match[V] {
	return Match[V]{
		prefix: prefix,
		value:  value,
		found:  true,
	}
}

func NewNotFoundMatch[V any]() Match[V] {
	return Match[V]{
		found: false,
	}
}
