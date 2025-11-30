package prefixtrie

import "net/netip"

type IPv4PrefixBuilder struct {
	seq uint32
	len uint8
}

func NewIPv4PrefixBuilder() *IPv4PrefixBuilder {
	return &IPv4PrefixBuilder{
		seq: 0,
		len: 0,
	}
}

func (b *IPv4PrefixBuilder) PushSeq(seq uint32, len uint8) {
	b.seq |= seq >> b.len
	b.len += len
}

func (b *IPv4PrefixBuilder) PopSeq(len uint8) {
	b.len -= len
	b.seq &= uint32PrefixMask(b.len)
}

func (b *IPv4PrefixBuilder) Build() netip.Prefix {
	return netip.PrefixFrom(Uint32ToAddr(b.seq), int(b.len))
}

func Uint32ToAddr(ipUint32 uint32) netip.Addr {
	var ipBytes [4]byte
	ipBytes[0] = byte(ipUint32 >> 24)
	ipBytes[1] = byte(ipUint32 >> 16)
	ipBytes[2] = byte(ipUint32 >> 8)
	ipBytes[3] = byte(ipUint32)

	return netip.AddrFrom4(ipBytes)
}
