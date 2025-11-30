// Warning: LLM generated
package ipv4map_test

import (
	"net/netip"
	"testing"

	"github.com/HT4w5/bgpsim-go/pkg/prefixtrie/ipv4map"
)

// Helper to parse address/prefix strings safely in tests
func mustAddr(s string) netip.Addr {
	return netip.MustParseAddr(s)
}

func mustPrefix(s string) netip.Prefix {
	return netip.MustParsePrefix(s)
}

// Test basic Insert and overwrite functionality
func TestInsert(t *testing.T) {
	m := ipv4map.NewIPv4Map()

	// 1. Basic Insertion
	p1 := mustPrefix("192.168.1.0/24")
	m.Insert(p1, "LAN Segment 1")

	_, v1, ok1 := m.Query(mustAddr("192.168.1.10"))
	if !ok1 || v1 != "LAN Segment 1" {
		t.Errorf("Insert failed: Expected 'LAN Segment 1', got %s (ok: %v)", v1, ok1)
	}

	// 2. Insertion with unmasked prefix (should mask internally)
	p2Unmasked := mustPrefix("10.0.0.10/8") // Will be stored as 10.0.0.0/8
	m.Insert(p2Unmasked, "Global Network")

	_, v2, ok2 := m.Query(mustAddr("10.1.2.3"))
	if !ok2 || v2 != "Global Network" {
		t.Errorf("Insert masking failed: Expected 'Global Network', got %s (ok: %v)", v2, ok2)
	}

	// 3. Overwrite existing prefix
	m.Insert(p1, "LAN Segment A (New)")
	_, v3, ok3 := m.Query(mustAddr("192.168.1.10"))
	if !ok3 || v3 != "LAN Segment A (New)" {
		t.Errorf("Insert overwrite failed: Expected 'LAN Segment A (New)', got %s (ok: %v)", v3, ok3)
	}
}

// Test Longest Prefix Match (LPM) logic, which was the primary bug area.
func TestQueryLPM(t *testing.T) {
	m := ipv4map.NewIPv4Map()

	// Insert prefixes from widest to narrowest
	m.Insert(mustPrefix("0.0.0.0/0"), "Default")
	m.Insert(mustPrefix("10.0.0.0/8"), "Internal")
	m.Insert(mustPrefix("10.10.0.0/16"), "Region A")
	m.Insert(mustPrefix("10.10.10.0/24"), "Host Group")

	tests := []struct {
		name           string
		queryAddr      string
		expectedPrefix string
		expectedValue  string
		expectedFound  bool
	}{
		{
			name:           "Perfect /24 Match (Longest)",
			queryAddr:      "10.10.10.5",
			expectedPrefix: "10.10.10.0/24",
			expectedValue:  "Host Group",
			expectedFound:  true,
		},
		{
			name:           "Intermediate /16 Match",
			queryAddr:      "10.10.20.5", // Falls outside /24, inside /16
			expectedPrefix: "10.10.0.0/16",
			expectedValue:  "Region A",
			expectedFound:  true,
		},
		{
			name:           "Widest /8 Match",
			queryAddr:      "10.20.30.40", // Falls outside /16, inside /8
			expectedPrefix: "10.0.0.0/8",
			expectedValue:  "Internal",
			expectedFound:  true,
		},
		{
			name:           "Default /0 Match",
			queryAddr:      "11.11.11.11", // Falls outside /8, inside /0
			expectedPrefix: "0.0.0.0/0",
			expectedValue:  "Default",
			expectedFound:  true,
		},
		{
			name:           "/32 Host Match",
			queryAddr:      "10.10.10.10",
			expectedPrefix: "10.10.10.0/24", // The /24 is still the best match
			expectedValue:  "Host Group",
			expectedFound:  true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			prefix, value, found := m.Query(mustAddr(tt.queryAddr))

			if found != tt.expectedFound {
				t.Fatalf("Found mismatch: Expected %v, got %v", tt.expectedFound, found)
			}
			if found && prefix.String() != tt.expectedPrefix {
				t.Errorf("Prefix mismatch: Expected %s, got %s", tt.expectedPrefix, prefix.String())
			}
			if found && value != tt.expectedValue {
				t.Errorf("Value mismatch: Expected %s, got %s", tt.expectedValue, value)
			}
		})
	}
}

// Test various non-LPM edge cases for Query
func TestQueryEdgeCases(t *testing.T) {
	m := ipv4map.NewIPv4Map()
	m.Insert(mustPrefix("192.168.0.0/16"), "Office")

	// 1. Query an empty map
	mEmpty := ipv4map.NewIPv4Map()
	_, _, okEmpty := mEmpty.Query(mustAddr("8.8.8.8"))
	if okEmpty {
		t.Error("Expected no match on empty map, but found one")
	}

	// 2. Query an address that is not in any prefix
	_, _, okMiss := m.Query(mustAddr("1.2.3.4"))
	if okMiss {
		t.Error("Expected no match for address 1.2.3.4, but found one")
	}

	// 3. Query the network address (first address in the prefix)
	p, _, ok := m.Query(mustAddr("192.168.0.0"))
	if !ok || p.String() != "192.168.0.0/16" {
		t.Errorf("Expected match for 192.168.0.0/16, got %s", p.String())
	}
}

// Test Delete functionality
func TestDelete(t *testing.T) {
	m := ipv4map.NewIPv4Map()
	p1 := mustPrefix("172.16.1.0/24")
	p2 := mustPrefix("172.16.2.0/24")

	m.Insert(p1, "Active 1")
	m.Insert(p2, "Active 2")

	// 1. Delete an existing prefix
	deleted := m.Delete(p1)
	if !deleted {
		t.Error("Expected Delete to return true for existing prefix")
	}

	// 2. Query the deleted prefix (should fail)
	_, _, okDeleted := m.Query(mustAddr("172.16.1.5"))
	if okDeleted {
		t.Error("Query found a match for a deleted prefix")
	}

	// 3. Delete a non-existent prefix
	pNonExist := mustPrefix("172.16.3.0/24")
	notDeleted := m.Delete(pNonExist)
	if notDeleted {
		t.Error("Expected Delete to return false for non-existent prefix")
	}

	// 4. Ensure soft deletion works (other prefixes remain)
	_, v2, ok2 := m.Query(mustAddr("172.16.2.5"))
	if !ok2 || v2 != "Active 2" {
		t.Errorf("Active prefix was unexpectedly deleted: %s", v2)
	}

	// 5. Insert after soft deletion (should insert, not reuse deleted slot in this simple implementation)
	m.Insert(p1, "Re-inserted")
	_, v3, ok3 := m.Query(mustAddr("172.16.1.5"))
	if !ok3 || v3 != "Re-inserted" {
		t.Errorf("Re-insertion failed, got %s", v3)
	}
}
