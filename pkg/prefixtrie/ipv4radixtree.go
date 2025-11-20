package prefixtrie

import (
	"encoding/binary"
	"math/bits"
	"net/netip"
)

const ipv4AddrLen = 32

type IPv4RadixTree[V any] struct {
	root *IPv4RadixTreeNode[V]
	len  int
}

type IPv4RadixTreeNode[V any] struct {
	zero *IPv4RadixTreeNode[V]
	one  *IPv4RadixTreeNode[V]

	seq uint32
	len uint8

	val *V
}

func NewIPv4RadixTree[V any]() *IPv4RadixTree[V] {
	return &IPv4RadixTree[V]{
		len: 0,
		root: &IPv4RadixTreeNode[V]{
			zero: nil,
			one:  nil,
			seq:  0,
			len:  0,
			val:  nil,
		},
	}
}

func (t *IPv4RadixTree[V]) Insert(p netip.Prefix, v V) error {
	seqBytes := p.Addr().As4()
	seq := binary.BigEndian.Uint32(seqBytes[:])
	len := uint8(p.Bits())
	node := t.root

	for true {
		if len == node.len {
			if seq == node.seq {
				// Override existing
				node.val = &v
				return nil
			}
		}
	}
}

func (t *IPv4RadixTree[V]) Query(addr netip.Addr) Match[V] {
	if t.len == 0 {
		return NewNotFoundMatch[V]()
	}

	seqBytes := addr.As4()
	seq := binary.BigEndian.Uint32(seqBytes[:])
	len := uint8(ipv4AddrLen)
	node := t.root
	builder := NewIPv4PrefixBuilder()

	for true {
		builder.PushSeq(node.seq, node.len)
		// Special case where prefix reaches length of address
		if len == node.len {
			if seq == node.seq {
				return NewFoundMatch(builder.Build(), *node.val)
			}
			break
		}

		if (seq & uint32PrefixMask(node.len)) == node.seq {
			var nextNode *IPv4RadixTreeNode[V]
			if uint32GetBit(seq, node.len) {
				// Next bit is 1
				if node.one != nil {
					nextNode = node.one
				} else {
					// No further match
					if node.val != nil {
						return NewFoundMatch(builder.Build(), *node.val)
					} else {
						// No value
						break
					}
				}
			} else {
				// Next bit is 0
				if node.zero != nil {
					nextNode = node.zero
				} else {
					// No further match
					if node.val != nil {
						return NewFoundMatch(builder.Build(), *node.val)
					} else {
						// No value
						break
					}
				}
			}
			seq <<= node.len
			len -= node.len
			node = nextNode
		}
	}

	return NewNotFoundMatch[V]()
}

func (t *IPv4RadixTree[V]) Delete(prefix netip.Prefix) bool {
	return false
}

func (t *IPv4RadixTree[V]) Len() int {
	return t.len
}

func uint32PrefixMask(len uint8) uint32 {
	return ^uint32(0) << (32 - len)
}

func uint32GetBit(seq uint32, idx uint8) bool {
	mask := uint32(1 << (31 - idx))
	return (seq & mask) != 0
}

// Get longest common prefix length of two uint32 values
func uint32LCPL(a uint32, b uint32) int {
	diff := a ^ b
	return bits.LeadingZeros32(diff)
}
