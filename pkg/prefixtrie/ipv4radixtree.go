package prefixtrie

import "net/netip"

type IPv4RadixTree[V any] struct {
	root *IPv4RadixTreeNode[V]
	len  int
}

type IPv4RadixTreeNode[V any] struct {
	zero *IPv4RadixTreeNode[V]
	one  *IPv4RadixTreeNode[V]

	sequence uint32
	length   int

	value *V
}

func NewIPv4RadixTree[V any]() *IPv4RadixTree[V] {
	return &IPv4RadixTree[V]{
		len: 0,
		root: &IPv4RadixTreeNode[V]{
			zero:     nil,
			one:      nil,
			sequence: 0,
			length:   0,
			value:    nil,
		},
	}
}

func (t *IPv4RadixTree[V]) Insert(p netip.Prefix, v V) error {
	return nil
}

func (t *IPv4RadixTree[V]) Query(addr netip.Addr) Match[V] {
	return Match[V]{}
}

func (t *IPv4RadixTree[V]) Delete(prefix netip.Prefix) bool {
	return false
}

func (t *IPv4RadixTree[V]) Len() int {
	return 0
}
