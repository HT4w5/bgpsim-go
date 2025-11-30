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

func (n *IPv4RadixTreeNode[V]) forkAt(len uint8) bool {
	child := &IPv4RadixTreeNode[V]{
		zero: n.zero,
		one:  n.one,
		val:  n.val,
		seq:  n.seq << len,
		len:  n.len - len,
	}

	n.val = nil
	childDirection := uint32GetBit(n.seq, len) // Direction in which child is attached
	if childDirection {
		n.one = child
		n.zero = nil
	} else {
		n.one = nil
		n.zero = child
	}
	n.len = len
	n.seq &= uint32PrefixMask(len)

	return childDirection
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

	// Empty tree
	if t.len == 0 {
		t.root = &IPv4RadixTreeNode[V]{
			zero: nil,
			one:  nil,
			val:  &v,
			seq:  seq,
			len:  len,
		}
		t.len++
		return
	}
	var parent *IPv4RadixTreeNode[V]
	parent = nil
	direction := false
	node := t.root

	for true {
		if node == nil {
			break
		}

		lcpl := min(uint8(uint32LCPL(seq, node.seq)), len, node.len)

		if lcpl == node.len {
			// Go to next iteration
			direction = uint32GetBit(seq, node.len)
			parent = node
			seq <<= node.len
			len -= node.len
			if direction {
				node = parent.one
			} else {
				node = parent.zero
			}
			// Break on next iteration if no length left
			if len == 0 {
				node = nil
			}
			continue
		}

		// Fork current node (N)
		//  ...      ...
		//   |        |
		//   N   ->   N
		//  / \      / \
		// X   Y   nil  B
		//             / \
		//            X   Y
		direction = !node.forkAt(lcpl)
		parent = node
		seq <<= lcpl
		len -= lcpl
		node = nil
	}

	// Insert new value
	if len == 0 {
		// Just set value of parent
		if parent.val == nil {
			t.len++
		}
		parent.val = &v
		return
	}

	// Attach new node to parent
	t.len++
	newNode := &IPv4RadixTreeNode[V]{
		zero: nil,
		one:  nil,
		val:  &v,
		seq:  seq,
		len:  len,
	}

	if direction {
		parent.one = newNode
	} else {
		parent.zero = newNode
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
		if node.zero != nil && node.one != nil {
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
			if parent == nil {
				// Default route
				node.val = nil
				return true
			}
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
