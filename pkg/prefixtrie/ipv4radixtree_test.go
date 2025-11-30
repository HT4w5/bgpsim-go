package prefixtrie_test

import (
	"fmt"
	"net/netip"
	"testing"

	"github.com/HT4w5/bgpsim-go/pkg/prefixtrie"
)

func TestIPv4RadixTreeBasic(t *testing.T) {
	rt := prefixtrie.NewIPv4RadixTree[int]()

	rt.Insert(netip.MustParsePrefix("114.51.4.0/24"), 1919810)
	m := rt.Query(netip.MustParseAddr("114.51.4.1"))

	if !m.Found() || m.GetValue() != 1919810 || m.GetPrefix().String() != "114.51.4.0/24" {
		t.Errorf("Basic lookup failed: got %v, %v, %v", m.Found(), m.GetValue(), m.GetPrefix())
	}
}

func TestIPv4RadixTreeUpdate(t *testing.T) {
	rt := prefixtrie.NewIPv4RadixTree[int]()
	p := netip.MustParsePrefix("10.0.0.0/8")

	// Initial Insert
	rt.Insert(p, 100)
	m1 := rt.Query(netip.MustParseAddr("10.0.0.1"))
	if m1.GetValue() != 100 {
		t.Errorf("Expected 100, got %d", m1.GetValue())
	}

	// Overwrite
	rt.Insert(p, 200)
	m2 := rt.Query(netip.MustParseAddr("10.0.0.1"))
	if m2.GetValue() != 200 {
		t.Errorf("Expected 200 after update, got %d", m2.GetValue())
	}
}

func TestIPv4RadixTreeNoMatch(t *testing.T) {
	rt := prefixtrie.NewIPv4RadixTree[int]()

	// Insert something that doesn't cover everything
	rt.Insert(netip.MustParsePrefix("192.168.0.0/16"), 1)

	// Query something outside
	m := rt.Query(netip.MustParseAddr("10.0.0.1"))
	if m.Found() {
		t.Error("Expected no match for 10.0.0.1, but found one")
	}
}

func TestIPv4RadixTreeEmptyQuery(t *testing.T) {
	rt := prefixtrie.NewIPv4RadixTree[int]()
	m := rt.Query(netip.MustParseAddr("114.51.4.0"))
	if m.Found() {
		t.Error("Query() returned found on empty tree")
	}
}

func TestIPv4RadixTreeInsertQuery(t *testing.T) {
	rt := prefixtrie.NewIPv4RadixTree[string]()

	// Define prefix table
	inserts := []struct {
		prefix string
		val    string
	}{
		{"0.0.0.0/0", "default"}, // Default route
		{"10.0.0.0/8", "private-10"},
		{"10.1.0.0/16", "private-10-1"},
		{"10.1.5.0/24", "private-10-1-5"},
		{"192.168.1.0/24", "home-network"},
		{"192.168.1.1/32", "router-interface"},
	}

	for _, ins := range inserts {
		p := netip.MustParsePrefix(ins.prefix)
		rt.Insert(p, ins.val)
	}

	fmt.Println(rt.GetTable())

	// Query test cases
	tests := []struct {
		name       string
		queryAddr  string
		wantFound  bool
		wantVal    string
		wantPrefix string
	}{
		{
			name:       "Exact match /32",
			queryAddr:  "192.168.1.1",
			wantFound:  true,
			wantVal:    "router-interface",
			wantPrefix: "192.168.1.1/32",
		},
		{
			name:       "Fallback to /24 inside /32 range",
			queryAddr:  "192.168.1.2",
			wantFound:  true,
			wantVal:    "home-network",
			wantPrefix: "192.168.1.0/24",
		},
		{
			name:       "Longest prefix match (deepest child)",
			queryAddr:  "10.1.5.50",
			wantFound:  true,
			wantVal:    "private-10-1-5",
			wantPrefix: "10.1.5.0/24",
		},
		{
			name:       "Match parent of deepest child",
			queryAddr:  "10.1.4.99",
			wantFound:  true,
			wantVal:    "private-10-1",
			wantPrefix: "10.1.0.0/16",
		},
		{
			name:       "Match grandparent",
			queryAddr:  "10.2.2.2",
			wantFound:  true,
			wantVal:    "private-10",
			wantPrefix: "10.0.0.0/8",
		},
		{
			name:       "Match default route",
			queryAddr:  "8.8.8.8",
			wantFound:  true,
			wantVal:    "default",
			wantPrefix: "0.0.0.0/0",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			addr := netip.MustParseAddr(tt.queryAddr)
			m := rt.Query(addr)

			if m.Found() != tt.wantFound {
				t.Fatalf("Query(%q) found=%v, want %v", tt.queryAddr, m.Found(), tt.wantFound)
			}

			if !tt.wantFound {
				return
			}

			if m.GetValue() != tt.wantVal {
				t.Errorf("Query(%q) value=%v, want %v", tt.queryAddr, m.GetValue(), tt.wantVal)
			}

			if m.GetPrefix().String() != tt.wantPrefix {
				t.Errorf("Query(%q) prefix=%v, want %v", tt.queryAddr, m.GetPrefix(), tt.wantPrefix)
			}
		})
	}
}

func TestIPv4RadixTreeInsertForkCase1(t *testing.T) {
	rt := prefixtrie.NewIPv4RadixTree[string]()

	// Define prefix table
	inserts := []struct {
		prefix string
		val    string
	}{
		{"223.5.5.5/32", "p1"},
		{"223.5.5.0/24", "p2"},
		{"223.5.0.0/16", "p3"}, // Should cause fork case 2
	}

	for _, ins := range inserts {
		p := netip.MustParsePrefix(ins.prefix)
		rt.Insert(p, ins.val)
	}

	fmt.Println(rt.GetTable())

	// Query test cases
	tests := []struct {
		name       string
		queryAddr  string
		wantFound  bool
		wantVal    string
		wantPrefix string
	}{
		{
			name:       "p1",
			queryAddr:  "223.5.5.5",
			wantFound:  true,
			wantVal:    "p1",
			wantPrefix: "223.5.5.5/32",
		},
		{
			name:       "p2",
			queryAddr:  "223.5.5.1",
			wantFound:  true,
			wantVal:    "p2",
			wantPrefix: "223.5.5.0/24",
		},
		{
			name:       "p3",
			queryAddr:  "223.5.1.1",
			wantFound:  true,
			wantVal:    "p3",
			wantPrefix: "223.5.0.0/16",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			addr := netip.MustParseAddr(tt.queryAddr)
			m := rt.Query(addr)

			if m.Found() != tt.wantFound {
				t.Fatalf("Query(%q) found=%v, want %v", tt.queryAddr, m.Found(), tt.wantFound)
			}

			if !tt.wantFound {
				return
			}

			if m.GetValue() != tt.wantVal {
				t.Errorf("Query(%q) value=%v, want %v", tt.queryAddr, m.GetValue(), tt.wantVal)
			}

			if m.GetPrefix().String() != tt.wantPrefix {
				t.Errorf("Query(%q) prefix=%v, want %v", tt.queryAddr, m.GetPrefix(), tt.wantPrefix)
			}
		})
	}
}

func TestIPv4RadixTreeInsertForkCase2(t *testing.T) {
	rt := prefixtrie.NewIPv4RadixTree[string]()

	// Define prefix table
	inserts := []struct {
		prefix string
		val    string
	}{
		{"0.0.0.0/24", "p1"},
		{"0.0.1.0/24", "p2"},
		{"0.1.0.0/16", "p3"}, // Should cause fork case 2
	}

	for _, ins := range inserts {
		p := netip.MustParsePrefix(ins.prefix)
		rt.Insert(p, ins.val)
	}

	fmt.Println(rt.GetTable())

	// Query test cases
	tests := []struct {
		name       string
		queryAddr  string
		wantFound  bool
		wantVal    string
		wantPrefix string
	}{
		{
			name:       "p1-1",
			queryAddr:  "0.0.0.1",
			wantFound:  true,
			wantVal:    "p1",
			wantPrefix: "0.0.0.0/24",
		},
		{
			name:       "p1-2",
			queryAddr:  "0.0.0.128",
			wantFound:  true,
			wantVal:    "p1",
			wantPrefix: "0.0.0.0/24",
		},
		{
			name:       "p1-3",
			queryAddr:  "0.0.0.254",
			wantFound:  true,
			wantVal:    "p1",
			wantPrefix: "0.0.0.0/24",
		},
		{
			name:       "p2-1",
			queryAddr:  "0.0.1.1",
			wantFound:  true,
			wantVal:    "p2",
			wantPrefix: "0.0.1.0/24",
		},
		{
			name:       "p2-2",
			queryAddr:  "0.0.1.32",
			wantFound:  true,
			wantVal:    "p2",
			wantPrefix: "0.0.1.0/24",
		},
		{
			name:       "p2-3",
			queryAddr:  "0.0.1.64",
			wantFound:  true,
			wantVal:    "p2",
			wantPrefix: "0.0.1.0/24",
		},
		{
			name:       "p3-1",
			queryAddr:  "0.1.0.1",
			wantFound:  true,
			wantVal:    "p3",
			wantPrefix: "0.1.0.0/16",
		},
		{
			name:       "p3-2",
			queryAddr:  "0.1.254.1",
			wantFound:  true,
			wantVal:    "p3",
			wantPrefix: "0.1.0.0/16",
		},
		{
			name:       "p3-3",
			queryAddr:  "0.1.18.48",
			wantFound:  true,
			wantVal:    "p3",
			wantPrefix: "0.1.0.0/16",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			addr := netip.MustParseAddr(tt.queryAddr)
			m := rt.Query(addr)

			if m.Found() != tt.wantFound {
				t.Fatalf("Query(%q) found=%v, want %v", tt.queryAddr, m.Found(), tt.wantFound)
			}

			if !tt.wantFound {
				return
			}

			if m.GetValue() != tt.wantVal {
				t.Errorf("Query(%q) value=%v, want %v", tt.queryAddr, m.GetValue(), tt.wantVal)
			}

			if m.GetPrefix().String() != tt.wantPrefix {
				t.Errorf("Query(%q) prefix=%v, want %v", tt.queryAddr, m.GetPrefix(), tt.wantPrefix)
			}
		})
	}
}

func TestIPv4RadixTreeDelete(t *testing.T) {
	rt := prefixtrie.NewIPv4RadixTree[int]()

	rt.Insert(netip.MustParsePrefix("114.51.4.0/24"), 1919810)

	if rt.Delete(netip.MustParsePrefix("114.51.0.0/16")) {
		t.Error("Deleting non-existent prefix shouldn't succeed")
	}

	if !rt.Delete(netip.MustParsePrefix("114.51.4.0/24")) {
		t.Error("Deleting existing prefix should succeed")
	}

	if rt.Len() != 0 {
		t.Errorf("Incorrect Len(), expected 0, got %d", rt.Len())
	}
}

// Warning: LLM generated

func TestIPv4RadixTree_HappyPath(t *testing.T) {
	tree := prefixtrie.NewIPv4RadixTree[string]()
	prefix1, _ := netip.ParsePrefix("192.168.1.0/24")
	prefix2, _ := netip.ParsePrefix("192.168.1.0/28")
	addr1, _ := netip.ParseAddr("192.168.1.10")
	addr2, _ := netip.ParseAddr("192.168.1.5")

	tree.Insert(prefix1, "Network 1")
	tree.Insert(prefix2, "Network 2")

	if match := tree.Query(addr1); !match.Found() || match.GetValue() != "Network 2" {
		t.Errorf("Expected to find 'Network 2', got %v", match.GetValue())
	}

	if match := tree.Query(addr2); !match.Found() || match.GetValue() != "Network 2" {
		t.Errorf("Expected to find 'Network 2', got %v", match.GetValue())
	}
}

func TestIPv4RadixTree_EmptyTree(t *testing.T) {
	tree := prefixtrie.NewIPv4RadixTree[string]()
	addr, _ := netip.ParseAddr("192.168.1.1")

	if match := tree.Query(addr); match.Found() {
		t.Errorf("Expected not to find any value, got %v", match.GetValue())
	}
}

func TestIPv4RadixTree_DeleteRoot(t *testing.T) {
	tree := prefixtrie.NewIPv4RadixTree[string]()
	prefix, _ := netip.ParsePrefix("0.0.0.0/0")
	addr, _ := netip.ParseAddr("192.168.1.1")

	tree.Insert(prefix, "Root Network")
	if !tree.Delete(prefix) {
		t.Errorf("Expected to delete root node, but failed")
	}

	if match := tree.Query(addr); match.Found() {
		t.Errorf("Expected not to find any value after deletion, got %v", match.GetValue())
	}
}

func TestIPv4RadixTree_DeleteLeaf(t *testing.T) {
	tree := prefixtrie.NewIPv4RadixTree[string]()
	prefix1, _ := netip.ParsePrefix("192.168.1.0/24")
	prefix2, _ := netip.ParsePrefix("192.168.1.0/28")
	addr, _ := netip.ParseAddr("192.168.1.10")

	tree.Insert(prefix1, "Network 1")
	tree.Insert(prefix2, "Network 2")

	if !tree.Delete(prefix2) {
		t.Errorf("Expected to delete leaf node, but failed")
	}

	if match := tree.Query(addr); !match.Found() || match.GetValue() != "Network 1" {
		t.Errorf("Expected to find 'Network 1' after deletion of 'Network 2', got %v", match.GetValue())
	}
}

func TestIPv4RadixTree_DeleteWithOneChild(t *testing.T) {
	tree := prefixtrie.NewIPv4RadixTree[string]()
	prefix1, _ := netip.ParsePrefix("192.168.1.0/24")
	prefix2, _ := netip.ParsePrefix("192.168.1.0/28")
	prefix3, _ := netip.ParsePrefix("192.168.1.16/28")
	addr, _ := netip.ParseAddr("192.168.1.20")

	tree.Insert(prefix1, "Network 1")
	tree.Insert(prefix2, "Network 2")
	tree.Insert(prefix3, "Network 3")

	fmt.Println(tree.GetTable())

	if !tree.Delete(prefix2) {
		t.Errorf("Expected to delete node with one child, but failed")
	}

	fmt.Println(tree.GetTable())

	if match := tree.Query(addr); !match.Found() || match.GetValue() != "Network 3" {
		t.Errorf("Expected to find 'Network 3' after deletion of 'Network 2', got %v", match.GetValue())
	}
}

func TestIPv4RadixTree_DeleteWithTwoChildren(t *testing.T) {
	tree := prefixtrie.NewIPv4RadixTree[string]()
	prefix1, _ := netip.ParsePrefix("192.168.1.0/24")
	prefix2, _ := netip.ParsePrefix("192.168.1.0/28")
	prefix3, _ := netip.ParsePrefix("192.168.1.16/28")
	addr1, _ := netip.ParseAddr("192.168.1.10")
	addr2, _ := netip.ParseAddr("192.168.1.20")

	tree.Insert(prefix1, "Network 1")
	tree.Insert(prefix2, "Network 2")
	tree.Insert(prefix3, "Network 3")

	if !tree.Delete(prefix1) {
		t.Errorf("Expected to delete node with two children, but failed")
	}

	if match := tree.Query(addr1); !match.Found() || match.GetValue() != "Network 2" {
		t.Errorf("Expected to find 'Network 2' after deletion of 'Network 1', got %v", match.GetValue())
	}

	if match := tree.Query(addr2); !match.Found() || match.GetValue() != "Network 3" {
		t.Errorf("Expected to find 'Network 3' after deletion of 'Network 1', got %v", match.GetValue())
	}
}

func TestIPv4RadixTree_InsertDuplicate(t *testing.T) {
	tree := prefixtrie.NewIPv4RadixTree[string]()
	prefix, _ := netip.ParsePrefix("192.168.1.0/24")
	addr, _ := netip.ParseAddr("192.168.1.10")

	tree.Insert(prefix, "Network 1")
	tree.Insert(prefix, "Network 2")

	if match := tree.Query(addr); !match.Found() || match.GetValue() != "Network 2" {
		t.Errorf("Expected to find updated 'Network 2', got %v", match.GetValue())
	}
}

func TestIPv4RadixTree_QueryExactMatch(t *testing.T) {
	tree := prefixtrie.NewIPv4RadixTree[string]()
	prefix, _ := netip.ParsePrefix("192.168.1.0/24")
	addr, _ := netip.ParseAddr("192.168.1.10")

	tree.Insert(prefix, "Network 1")

	if match := tree.Query(addr); !match.Found() || match.GetValue() != "Network 1" {
		t.Errorf("Expected to find 'Network 1', got %v", match.GetValue())
	}
}

func TestIPv4RadixTree_QueryLongestPrefix(t *testing.T) {
	tree := prefixtrie.NewIPv4RadixTree[string]()
	prefix1, _ := netip.ParsePrefix("192.168.0.0/16")
	prefix2, _ := netip.ParsePrefix("192.168.1.0/24")
	addr, _ := netip.ParseAddr("192.168.1.10")

	tree.Insert(prefix1, "Network 1")
	tree.Insert(prefix2, "Network 2")

	if match := tree.Query(addr); !match.Found() || match.GetValue() != "Network 2" {
		t.Errorf("Expected to find 'Network 2', got %v", match.GetValue())
	}
}

func TestIPv4RadixTree_QueryNoMatch(t *testing.T) {
	tree := prefixtrie.NewIPv4RadixTree[string]()
	prefix, _ := netip.ParsePrefix("192.168.1.0/24")
	addr, _ := netip.ParseAddr("10.0.0.1")

	tree.Insert(prefix, "Network 1")

	if match := tree.Query(addr); match.Found() {
		t.Errorf("Expected not to find any value, got %v", match.GetValue())
	}
}

func TestIPv4RadixTree_Len(t *testing.T) {
	tree := prefixtrie.NewIPv4RadixTree[string]()
	prefix1, _ := netip.ParsePrefix("192.168.0.0/16")
	prefix2, _ := netip.ParsePrefix("192.168.1.0/24")

	tree.Insert(prefix1, "Network 1")
	tree.Insert(prefix2, "Network 2")

	if length := tree.Len(); length != 2 {
		t.Errorf("Expected length 2, got %d", length)
	}

	tree.Delete(prefix1)

	if length := tree.Len(); length != 1 {
		t.Errorf("Expected length 1, got %d", length)
	}
}
