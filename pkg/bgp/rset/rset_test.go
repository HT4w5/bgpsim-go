package rset

import (
	"net/netip"
	"testing"

	"github.com/HT4w5/bgpsim-go/pkg/bgp/route"
	"github.com/HT4w5/bgpsim-go/pkg/nexthop"
)

// Helper to create a test route
func makeRoute(localPref uint32, asPathLen int, metric int, weight int) *route.BgpRoute {
	asPath := make([]uint32, asPathLen)
	for i := 0; i < asPathLen; i++ {
		asPath[i] = uint32(i + 1)
	}
	return route.New(
		route.WithLocalPreference(localPref),
		route.WithAsPath(asPath),
		route.WithMetric(metric),
		route.WithWeight(weight),
		route.WithPrefix(netip.MustParsePrefix("10.0.0.0/24")),
	)
}

// Helper to create a route with a unique nexthop (for unique hash)
func makeRouteWithNexthop(localPref uint32, asPathLen int, metric int, weight int, nhIP string) *route.BgpRoute {
	asPath := make([]uint32, asPathLen)
	for i := 0; i < asPathLen; i++ {
		asPath[i] = uint32(i + 1)
	}
	nh := nexthop.New(nexthop.WithIP(netip.MustParseAddr(nhIP)))
	return route.New(
		route.WithLocalPreference(localPref),
		route.WithAsPath(asPath),
		route.WithMetric(metric),
		route.WithWeight(weight),
		route.WithPrefix(netip.MustParsePrefix("10.0.0.0/24")),
		route.WithNextHop(nh),
	)
}

// RouteSet creation panics when provided route slice is empty or nil
func TestNew_PanicOnEmptySlice(t *testing.T) {
	defer func() {
		if r := recover(); r == nil {
			t.Error("Expected panic when creating RouteSet with empty slice, but no panic occurred")
		}
	}()

	New([]*route.BgpRoute{})
}

func TestNew_PanicOnNilSlice(t *testing.T) {
	defer func() {
		if r := recover(); r == nil {
			t.Error("Expected panic when creating RouteSet with nil slice, but no panic occurred")
		}
	}()

	New(nil)
}

// RouteSet creation without multipath
func TestNew_WithoutMultipath(t *testing.T) {
	routes := []*route.BgpRoute{
		makeRoute(100, 3, 10, 0), // localPref=100, asPath=3, metric=10
		makeRoute(200, 2, 5, 0),  // localPref=200, asPath=2, metric=5 (best - highest localPref)
		makeRoute(150, 1, 20, 0), // localPref=150, asPath=1, metric=20
	}

	rs := New(routes)

	if rs.BestPath() == nil {
		t.Fatal("BestPath should not be nil")
	}

	if rs.BestPath().LocalPreference() != 200 {
		t.Errorf("Expected best path with LocalPreference 200, got %d", rs.BestPath().LocalPreference())
	}

	if len(rs.MultipathSet()) != 0 {
		t.Errorf("Expected empty MultipathSet when multipath is disabled, got %d routes", len(rs.MultipathSet()))
	}
}

func TestNew_WithoutMultipath_TieBreak(t *testing.T) {
	r1 := makeRouteWithNexthop(100, 2, 10, 0, "192.168.1.1")
	r1.SetArrival(1000) // Earlier arrival

	r2 := makeRouteWithNexthop(100, 2, 10, 0, "192.168.1.2")
	r2.SetArrival(2000) // Later arrival

	routes := []*route.BgpRoute{r2, r1}

	rs := New(routes)

	if rs.BestPath().Arrival() != 1000 {
		t.Errorf("Expected best path with arrival 1000, got %d", rs.BestPath().Arrival())
	}
}

// RouteSet creation with multipath
func TestNew_WithMultipath(t *testing.T) {
	r1 := makeRouteWithNexthop(100, 2, 10, 0, "192.168.1.1")
	r1.SetArrival(1000)

	r2 := makeRouteWithNexthop(100, 2, 10, 0, "192.168.1.2")
	r2.SetArrival(2000)

	r3 := makeRouteWithNexthop(50, 2, 10, 0, "192.168.1.3") // Different localPref - not equal

	routes := []*route.BgpRoute{r1, r2, r3}

	rs := New(routes, WithMultipath(true))

	if rs.BestPath() == nil {
		t.Fatal("BestPath should not be nil")
	}

	if rs.BestPath().LocalPreference() != 100 {
		t.Errorf("Expected best path with LocalPreference 100, got %d", rs.BestPath().LocalPreference())
	}

	multipathSet := rs.MultipathSet()
	if len(multipathSet) != 2 {
		t.Errorf("Expected MultipathSet with 2 routes, got %d", len(multipathSet))
	}

	for _, r := range multipathSet {
		if r.LocalPreference() != 100 {
			t.Errorf("Expected all multipath routes to have LocalPreference 100, got %d", r.LocalPreference())
		}
	}
}

func TestNew_WithMultipath_BestPathSelection(t *testing.T) {
	r1 := makeRouteWithNexthop(100, 2, 10, 0, "192.168.1.1")
	r1.SetArrival(3000) // Latest

	r2 := makeRouteWithNexthop(100, 2, 10, 0, "192.168.1.2")
	r2.SetArrival(1000) // Earliest - should be best

	r3 := makeRouteWithNexthop(100, 2, 10, 0, "192.168.1.3")
	r3.SetArrival(2000)

	routes := []*route.BgpRoute{r1, r2, r3}

	rs := New(routes, WithMultipath(true))

	if rs.BestPath().Arrival() != 1000 {
		t.Errorf("Expected best path with arrival 1000, got %d", rs.BestPath().Arrival())
	}

	if len(rs.MultipathSet()) != 3 {
		t.Errorf("Expected MultipathSet with 3 routes, got %d", len(rs.MultipathSet()))
	}
}

func TestNew_WithMultipath_OnlyOne(t *testing.T) {
	r1 := makeRouteWithNexthop(100, 2, 10, 0, "192.168.1.1")

	routes := []*route.BgpRoute{r1}

	rs := New(routes, WithMultipath(true))

	if rs.BestPath() != r1 {
		t.Errorf("Expected best path %v, got %v", r1, rs.BestPath())
	}
}

func TestNew_WithMultipath_MultipathLimit(t *testing.T) {
	r1 := makeRouteWithNexthop(100, 2, 10, 0, "192.168.1.1")
	r2 := makeRouteWithNexthop(100, 2, 10, 0, "192.168.1.2")
	r3 := makeRouteWithNexthop(100, 2, 10, 0, "192.168.1.3")
	r4 := makeRouteWithNexthop(100, 2, 10, 0, "192.168.1.4")
	r5 := makeRouteWithNexthop(100, 2, 10, 0, "192.168.1.5")
	r6 := makeRouteWithNexthop(100, 2, 10, 0, "192.168.1.6")
	r7 := makeRouteWithNexthop(100, 2, 10, 0, "192.168.1.7")
	r8 := makeRouteWithNexthop(100, 2, 10, 0, "192.168.1.8")
	r9 := makeRouteWithNexthop(100, 2, 10, 0, "192.168.1.9")
	r10 := makeRouteWithNexthop(100, 2, 10, 0, "192.168.1.10")

	routes := []*route.BgpRoute{r1, r2, r3, r4, r5, r6, r7, r8, r9, r10}

	rs := New(routes, WithMultipath(true))

	multipathCount := len(rs.MultipathSet())
	if multipathCount != maxPaths {
		t.Errorf("Expected len(multipathSet) %d, got %d", maxPaths, multipathCount)
	}
}

// RouteSet Eq without multipath
func TestEq_WithoutMultipath_Equal(t *testing.T) {
	r1 := makeRoute(100, 2, 10, 0)
	r2 := makeRoute(50, 3, 5, 0)

	rs1 := New([]*route.BgpRoute{r1, r2})
	rs2 := New([]*route.BgpRoute{r1, r2})

	if !rs1.Eq(rs2) {
		t.Error("Expected RouteSets with same best path to be equal")
	}
}

func TestEq_WithoutMultipath_NotEqual(t *testing.T) {
	r1 := makeRoute(100, 2, 10, 0)
	r2 := makeRoute(200, 3, 5, 0)

	rs1 := New([]*route.BgpRoute{r1})
	rs2 := New([]*route.BgpRoute{r2})

	if rs1.Eq(rs2) {
		t.Error("Expected RouteSets with different best paths to be not equal")
	}
}

// RouteSet Eq with multipath
func TestEq_WithMultipath_Equal(t *testing.T) {
	r1 := makeRouteWithNexthop(100, 2, 10, 0, "192.168.1.1")
	r1.SetArrival(1000)

	r2 := makeRouteWithNexthop(100, 2, 10, 0, "192.168.1.2")
	r2.SetArrival(2000)

	rs1 := New([]*route.BgpRoute{r1, r2}, WithMultipath(true))
	rs2 := New([]*route.BgpRoute{r1, r2}, WithMultipath(true))

	if !rs1.Eq(rs2) {
		t.Error("Expected RouteSets with same multipath set to be equal")
	}
}

func TestEq_WithMultipath_NotEqual_DifferentSet(t *testing.T) {
	r1 := makeRouteWithNexthop(100, 2, 10, 0, "192.168.1.1")
	r1.SetArrival(1000)

	r2 := makeRouteWithNexthop(100, 2, 10, 0, "192.168.1.2")
	r2.SetArrival(2000)

	r3 := makeRouteWithNexthop(100, 2, 10, 0, "192.168.1.3")
	r3.SetArrival(3000)

	rs1 := New([]*route.BgpRoute{r1, r2}, WithMultipath(true))
	rs2 := New([]*route.BgpRoute{r1, r3}, WithMultipath(true))

	if rs1.Eq(rs2) {
		t.Error("Expected RouteSets with different multipath sets to be not equal")
	}
}

func TestEq_WithMultipath_NotEqual_DifferentBestPath(t *testing.T) {
	r1 := makeRouteWithNexthop(100, 2, 10, 0, "192.168.1.1")
	r1.SetArrival(1000)

	r2 := makeRouteWithNexthop(100, 2, 10, 0, "192.168.1.2")
	r2.SetArrival(2000)

	rs1 := New([]*route.BgpRoute{r1, r2}, WithMultipath(true))

	r3 := makeRouteWithNexthop(100, 2, 10, 0, "192.168.1.3")
	r3.SetArrival(500)

	rs2 := New([]*route.BgpRoute{r1, r3}, WithMultipath(true))

	if rs1.Eq(rs2) {
		t.Error("Expected RouteSets with different multipath sets to be not equal")
	}
}

func TestEq_DifferentMultipathConfig(t *testing.T) {
	r1 := makeRoute(100, 2, 10, 0)

	rs1 := New([]*route.BgpRoute{r1}, WithMultipath(false))
	rs2 := New([]*route.BgpRoute{r1}, WithMultipath(true))

	if rs1.Eq(rs2) {
		t.Error("Expected RouteSets with different multipath config to be not equal")
	}
}
