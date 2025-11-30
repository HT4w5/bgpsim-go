package ipv4map

import "net/netip"

// Simple implementation for cross-testing

type prefixEntry struct {
	prefix  netip.Prefix
	value   string
	deleted bool
}

type IPv4Map struct {
	table []prefixEntry
	len   int
}

func NewIPv4Map() *IPv4Map {
	return &IPv4Map{
		table: make([]prefixEntry, 0),
	}
}

func (m *IPv4Map) Insert(p netip.Prefix, s string) {
	prefix := p.Masked()

	found := false
	for i, v := range m.table {
		if v.deleted {
			continue
		}
		if v.prefix == prefix {
			m.table[i].value = s
			found = true
			break
		}
	}

	if !found {
		m.len++
		m.table = append(m.table, prefixEntry{
			prefix:  prefix,
			value:   s,
			deleted: false,
		})
	}
}

func (m *IPv4Map) Query(addr netip.Addr) (netip.Prefix, string, bool) {
	found := false
	bestMatch := netip.MustParsePrefix("0.0.0.0/0")
	bestValue := ""
	for _, v := range m.table {
		if v.deleted {
			continue
		}
		if v.prefix.Contains(addr) {
			found = true
			if bestMatch.Bits() <= v.prefix.Bits() {
				bestMatch = v.prefix
				bestValue = v.value
			}
		}
	}

	return bestMatch, bestValue, found
}

func (m *IPv4Map) Delete(p netip.Prefix) bool {
	prefix := p.Masked()

	found := false
	for i, v := range m.table {
		if v.deleted {
			continue
		}
		if v.prefix == prefix {
			m.table[i].deleted = true
			found = true
			m.len--
			break
		}
	}

	return found
}

func (m *IPv4Map) Len() int {
	return m.len
}
