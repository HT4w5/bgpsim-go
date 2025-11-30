package prefixtrie

import (
	"encoding/binary"
	"fmt"
	"net/netip"
	"strings"
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

func (n *IPv4RadixTreeNode[V]) mergeWith(o *IPv4RadixTreeNode[V]) {
	n.zero = o.zero
	n.one = o.one
	n.val = o.val
	n.seq |= o.seq >> n.len
	n.len += o.len
}

func NewIPv4RadixTree[V any]() *IPv4RadixTree[V] {
	return &IPv4RadixTree[V]{
		len:  0,
		root: nil,
	}
}

func (t *IPv4RadixTree[V]) Insert(p netip.Prefix, v V) {
	p = p.Masked()

	seqBytes := p.Addr().As4()
	seq := binary.BigEndian.Uint32(seqBytes[:])
	len := uint8(p.Bits())
	node := t.root

	for true {
		// Equally-lengthed current and node sequence
		lenCmp := uint8Cmp(len, node.len)
		// Longest common prefix length
		lcpl := uint32LCPL(seq, node.seq)

		if lenCmp == 0 {
			if lcpl >= int(len) {
				// Override value of current node
				if node.val == nil {
					t.len++
				}
				node.val = &v
				return
			}
		}

		// Current sequence is longer than node sequence
		// AND lcpl not shorter than node sequence
		if lenCmp > 0 && lcpl >= int(node.len) {
			// The first bit where the current seqence is longer than the node sequence
			nextBit := uint32GetBit(seq, node.len)
			// Go to next iteration
			var nextNode *IPv4RadixTreeNode[V]
			if nextBit {
				// Next bit is 1
				if node.one != nil {
					nextNode = node.one
				} else {
					// No further match
					// Create new node
					node.one = &IPv4RadixTreeNode[V]{
						zero: nil,
						one:  nil,
						seq:  seq << node.len,
						len:  len - node.len,
						val:  &v,
					}
					// Increment size counter
					t.len++
					return
				}
			} else {
				// Next bit is 0
				if node.zero != nil {
					nextNode = node.zero
				} else {
					// No further match
					// Create new node
					node.zero = &IPv4RadixTreeNode[V]{
						zero: nil,
						one:  nil,
						seq:  seq << node.len,
						len:  len - node.len,
						val:  &v,
					}
					// Increment size counter
					t.len++
					return
				}
			}
			seq <<= node.len
			len -= node.len
			node = nextNode
			continue
		}

		if lcpl >= int(len) {
			// The first bit where the node seqence is longer than the current length
			nextBit := uint32GetBit(node.seq, len)
			// Fork case 1
			//     ...
			//      |
			//      A <-- (node)
			//     / \
			//   nil  C
			//        |
			//       ...
			//
			// Original node is divided into A and C,
			// while current insertion is bound to A.val

			nodeC := &IPv4RadixTreeNode[V]{
				zero: node.zero,
				one:  node.one,
				seq:  node.seq << len,
				len:  node.len - len,
				val:  node.val,
			}

			// Use original node as A
			if nextBit {
				// C is attached to A.one
				node.zero = nil
				node.one = nodeC
			} else {
				// The opposite
				node.one = nil
				node.zero = nodeC
			}

			node.seq &= uint32PrefixMask(len)
			node.len = len
			node.val = &v
			// Increment size counter
			t.len++
			return
		}

		// Fork case 2
		//     ...
		//      |
		//      A <-- (node)
		//     / \
		//    B   C
		//        |
		//       ...
		//
		// Original node is divided into A and C,
		// while current insertion is bound to B.val,
		// which is newly created

		// The first bit where the current seqence is longer than the LCPL
		nextBit := uint32GetBit(seq, uint8(lcpl))

		nodeB := &IPv4RadixTreeNode[V]{
			zero: nil,
			one:  nil,
			seq:  seq << lcpl,
			len:  len - uint8(lcpl), // If previous conditions are met, len is guaranteed to be larger than lcpl
			val:  &v,
		}

		nodeC := &IPv4RadixTreeNode[V]{
			zero: node.zero,
			one:  node.one,
			seq:  node.seq << lcpl,
			len:  node.len - uint8(lcpl),
			val:  node.val,
		}

		// Use original node as A
		if nextBit {
			// B is attached to A.one
			// C is attached to A.zero
			node.one = nodeB
			node.zero = nodeC
		} else {
			// The opposite
			node.zero = nodeB
			node.one = nodeC
		}

		node.seq &= uint32PrefixMask(uint8(lcpl))
		node.len = uint8(lcpl)
		node.val = nil
		// Increment size counter
		t.len++
		return
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
	bestMatch := NewNotFoundMatch[V]()

	for true {
		builder.PushSeq(node.seq, node.len)
		// Special case where prefix reaches length of address
		if len == node.len {
			if seq == node.seq {
				bestMatch = NewFoundMatch(builder.Build(), node.val)
			}
			break
		}

		if (seq & uint32PrefixMask(node.len)) == node.seq {
			// Update best match
			if node.val != nil {
				bestMatch = NewFoundMatch(builder.Build(), node.val)
			}
			// Go to next iteration
			var nextNode *IPv4RadixTreeNode[V]
			if uint32GetBit(seq, node.len) {
				// Next bit is 1
				if node.one != nil {
					nextNode = node.one
				} else {
					// No further match
					if node.val != nil {
						bestMatch = NewFoundMatch(builder.Build(), node.val)
					}
					break
				}
			} else {
				// Next bit is 0
				if node.zero != nil {
					nextNode = node.zero
				} else {
					// No further match
					if node.val != nil {
						bestMatch = NewFoundMatch(builder.Build(), node.val)
					}
					break
				}
			}
			seq <<= node.len
			len -= node.len
			node = nextNode
			continue
		}
		break
	}

	return bestMatch
}

type visitState int

const (
	pre visitState = iota
	afterZero
	afterOne
)

type stackEntry[V any] struct {
	node  *IPv4RadixTreeNode[V]
	state visitState
}

func (t *IPv4RadixTree[V]) GetTable() string {
	st := NewStack[*stackEntry[V]](t.len)
	st.Push(&stackEntry[V]{
		t.root,
		pre,
	})

	builder := NewIPv4PrefixBuilder()
	out := &strings.Builder{}

	for st.Len() > 0 {
		top, _ := st.Peek()
		switch top.state {
		case pre:
			// Enter node
			builder.PushSeq(top.node.seq, top.node.len)

			if top.node.val != nil {
				fmt.Fprintf(out, "%s: %v\n", builder.Build().String(), *top.node.val)
			}

			// Goto zero
			top.state = afterZero
			if top.node.zero != nil {
				st.Push(&stackEntry[V]{
					top.node.zero,
					pre,
				})
			}
		case afterZero:
			// Goto one
			top.state = afterOne
			if top.node.one != nil {
				st.Push(&stackEntry[V]{
					top.node.one,
					pre,
				})
			}
		case afterOne:
			// Exit node
			builder.PopSeq(top.node.len)
			st.Pop()
		}
	}

	return out.String()
}

func (t *IPv4RadixTree[V]) Delete(prefix netip.Prefix) bool {
	prefix = prefix.Masked()
	if t.len == 0 {
		return false
	}

	seqBytes := prefix.Addr().As4()
	seq := binary.BigEndian.Uint32(seqBytes[:])
	len := uint8(prefix.Bits())
	node := t.root
	var parent *IPv4RadixTreeNode[V] // Track parent for merging
	parent = nil
	direction := false // Direction of last descent

	found := false
	for true {
		// Special case where prefix reaches length of address
		// Exact match
		if len == node.len {
			if seq == node.seq {
				found = true
			}
			break
		}

		if (seq & uint32PrefixMask(node.len)) == node.seq {
			// Go to next iteration
			var nextNode *IPv4RadixTreeNode[V]
			if uint32GetBit(seq, node.len) {
				// Next bit is 1
				if node.one != nil {
					direction = true
					nextNode = node.one
				} else {
					// No further match
					break
				}
			} else {
				// Next bit is 0
				if node.zero != nil {
					direction = false
					nextNode = node.zero
				} else {
					// No further match
					break
				}
			}
			seq <<= node.len
			len -= node.len
			parent = node
			node = nextNode
			continue
		}
		break
	}

	if found {
		// Cases where a value potentially doesn't exist
		// Default route or node with two children
		if node.zero != nil && node.one != nil || parent == nil {
			// 2 children
			// Only delete value
			if node.val == nil {
				return false
			}
			node.val = nil
			t.len--
			return true
		}
		// From now on, there must be a value to be deleted
		t.len--
		if node.zero == nil && node.one == nil {
			// No children
			// Delete node
			var child *IPv4RadixTreeNode[V]
			if direction {
				parent.one = nil
				child = parent.zero
			} else {
				parent.zero = nil
				child = parent.one
			}
			if parent.val == nil && child != nil {
				parent.mergeWith(child)
			}
		} else {
			// One child
			// Merge with child
			child := node.zero
			if child == nil {
				child = node.one
			}
			node.mergeWith(child)
		}
	}

	return found
}

func (t *IPv4RadixTree[V]) Len() int {
	return t.len
}
