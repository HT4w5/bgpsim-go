package prefixtrie

import "net/netip"

type Match[V any] struct {
	Prefix netip.Prefix
	Value  V
	Found  bool
}

type PrefixTrie[V any] interface {
	Insert(prefix netip.Prefix, value V) error
	Query(addr netip.Addr) Match[V]
	Delete(prefix netip.Prefix) bool
	Len() int
}
