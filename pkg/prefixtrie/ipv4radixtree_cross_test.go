package prefixtrie_test

import (
	"encoding/binary"
	"fmt"
	"math/rand"
	"net/netip"
	"testing"

	"github.com/HT4w5/bgpsim-go/pkg/prefixtrie"
	"github.com/HT4w5/bgpsim-go/pkg/prefixtrie/ipv4map"
)

const ipv4AddrLen = 32

type ptOperation int

const (
	ptInsert ptOperation = iota
	ptQuery
	ptDelete
)

var ptOps = []ptOperation{
	ptInsert,
	ptQuery,
	ptDelete,
}

type ptCrossTester struct {
	rt         *prefixtrie.IPv4RadixTree[string]
	st         *ipv4map.IPv4Map
	ran        *rand.Rand
	prefixPool []netip.Prefix
}

func newCrossTester(seed int64) *ptCrossTester {
	return &ptCrossTester{
		rt:         prefixtrie.NewIPv4RadixTree[string](),
		st:         ipv4map.NewIPv4Map(),
		ran:        rand.New(rand.NewSource(seed)),
		prefixPool: make([]netip.Prefix, 0),
	}
}

func (t *ptCrossTester) tick() (string, error) {
	op := ptOps[t.ran.Intn(len(ptOps))]

	var msg string
	var err error
	err = nil

	switch op {
	case ptInsert:
		prefix := randomIPv4Prefix(t.ran)
		msg = fmt.Sprintf("Insert(%s)", prefix.String())
		value := fmt.Sprintf("Network %s", prefix.String())
		t.rt.Insert(prefix, value)
		t.st.Insert(prefix, value)
		t.prefixPool = append(t.prefixPool, prefix)

		if t.rt.Len() != t.st.Len() {
			err = fmt.Errorf("Len() mismatch after Insert(%s). rt: %d, st: %d", prefix.String(), t.rt.Len(), t.st.Len())
		}
	case ptQuery:
		addr := randomIPv4Addr(t.ran)
		msg = fmt.Sprintf("Query(%s)", addr.String())
		rm := t.rt.Query(addr)
		sp, sv, ok := t.st.Query(addr)

		mismatch := false
		// Both found
		if rm.Found() == ok && ok {
			if rm.GetPrefix() != sp || rm.GetValue() != sv {
				mismatch = true
			}
		} else {
			if rm.Found() != ok {
				mismatch = true
			}
		}

		if mismatch {
			err = fmt.Errorf("mismatch on Query(%s). rt: %v,%s, %s; st: %v,%s, %s", addr.String(), rm.Found(), rm.GetPrefix().String(), rm.GetValue(), ok, sp.String(), sv)
		}
	case ptDelete:
		if len(t.prefixPool) == 0 {
			msg = "skipping Delete()"
			break
		}

		t.ran.Shuffle(len(t.prefixPool), func(i, j int) {
			tmp := t.prefixPool[i]
			t.prefixPool[i] = t.prefixPool[j]
			t.prefixPool[j] = tmp
		})

		prefix := t.prefixPool[len(t.prefixPool)-1]
		t.prefixPool = t.prefixPool[:len(t.prefixPool)-1]

		msg = fmt.Sprintf("Delete(%s)", prefix.String())

		rOk := t.rt.Delete(prefix)
		sOk := t.st.Delete(prefix)

		if rOk != sOk {
			err = fmt.Errorf("mismatch on Delete(%s). rt: %v; st: %v", prefix.String(), rOk, sOk)
			break
		}

		if t.rt.Len() != t.st.Len() {
			err = fmt.Errorf("Len() mismatch after Delete(%s). rt: %d, st: %d", prefix.String(), t.rt.Len(), t.st.Len())
		}
	}

	return msg, err
}

func randomIPv4Prefix(ran *rand.Rand) netip.Prefix {
	addrUint32 := ran.Uint32()
	var addrBytes [4]byte
	binary.BigEndian.PutUint32(addrBytes[:], addrUint32)

	ipAddr := netip.AddrFrom4(addrBytes)
	prefixLength := ran.Intn(ipv4AddrLen + 1)
	prefix := netip.PrefixFrom(ipAddr, prefixLength)
	return prefix.Masked()
}

func randomIPv4Addr(ran *rand.Rand) netip.Addr {
	addrUint32 := ran.Uint32()
	var addrBytes [4]byte
	binary.BigEndian.PutUint32(addrBytes[:], addrUint32)

	return netip.AddrFrom4(addrBytes)
}

func TestFixedSeed(t *testing.T) {
	ct := newCrossTester(1)
	for i := range 10 {
		msg, err := ct.tick()
		fmt.Printf("Iteration %d: %s    %v\n", i, msg, err)
		if err != nil {
			t.Errorf("Error on iteration %d: %v", i, err)
		}
	}
}

func TestMultipleSeeds(t *testing.T) {
	for seed := range 1000 {
		ct := newCrossTester(int64(seed))
		for range 1000 {
			_, err := ct.tick()
			if err != nil {
				t.Errorf("error on seed %d", seed)
				break
			}
		}
	}
}
