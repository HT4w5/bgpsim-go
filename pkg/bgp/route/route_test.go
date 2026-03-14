package route

import (
	"net/netip"
	"testing"
)

func TestCompareMultipath(t *testing.T) {
	tests := []struct {
		name     string
		a        *BgpRoute
		b        *BgpRoute
		expected int
	}{
		{
			name:     "both nil",
			a:        nil,
			b:        nil,
			expected: 0,
		},
		{
			name:     "a is nil",
			a:        nil,
			b:        New(),
			expected: 1,
		},
		{
			name:     "b is nil",
			a:        New(),
			b:        nil,
			expected: -1,
		},
		{
			name:     "equal routes",
			a:        New(),
			b:        New(),
			expected: 0,
		},
		{
			name:     "prefer higher weight - a wins",
			a:        New(WithWeight(100)),
			b:        New(WithWeight(50)),
			expected: -1,
		},
		{
			name:     "prefer higher weight - b wins",
			a:        New(WithWeight(50)),
			b:        New(WithWeight(100)),
			expected: 1,
		},
		{
			name:     "prefer higher localPreference - a wins",
			a:        New(WithLocalPreference(200)),
			b:        New(WithLocalPreference(100)),
			expected: -1,
		},
		{
			name:     "prefer higher localPreference - b wins",
			a:        New(WithLocalPreference(100)),
			b:        New(WithLocalPreference(200)),
			expected: 1,
		},
		{
			name:     "weight takes precedence over localPreference",
			a:        New(WithWeight(100), WithLocalPreference(50)),
			b:        New(WithWeight(50), WithLocalPreference(200)),
			expected: -1,
		},
		{
			name:     "prefer lower adminCost - a wins",
			a:        New(WithAdminCost(10)),
			b:        New(WithAdminCost(20)),
			expected: -1,
		},
		{
			name:     "prefer lower adminCost - b wins",
			a:        New(WithAdminCost(20)),
			b:        New(WithAdminCost(10)),
			expected: 1,
		},
		{
			name:     "localPreference takes precedence over adminCost",
			a:        New(WithLocalPreference(100), WithAdminCost(50)),
			b:        New(WithLocalPreference(200), WithAdminCost(10)),
			expected: 1,
		},
		{
			name:     "prefer shorter asPath - a wins",
			a:        New(WithAsPath([]uint32{1, 2})),
			b:        New(WithAsPath([]uint32{1, 2, 3})),
			expected: -1,
		},
		{
			name:     "prefer shorter asPath - b wins",
			a:        New(WithAsPath([]uint32{1, 2, 3, 4})),
			b:        New(WithAsPath([]uint32{1, 2})),
			expected: 1,
		},
		{
			name:     "adminCost takes precedence over asPath",
			a:        New(WithAdminCost(10), WithAsPath([]uint32{1, 2, 3, 4})),
			b:        New(WithAdminCost(20), WithAsPath([]uint32{1})),
			expected: -1,
		},
		{
			name:     "prefer lower metric - a wins",
			a:        New(WithMetric(10)),
			b:        New(WithMetric(20)),
			expected: -1,
		},
		{
			name:     "prefer lower metric - b wins",
			a:        New(WithMetric(30)),
			b:        New(WithMetric(15)),
			expected: 1,
		},
		{
			name:     "asPath length takes precedence over metric",
			a:        New(WithAsPath([]uint32{1}), WithMetric(100)),
			b:        New(WithAsPath([]uint32{1, 2, 3}), WithMetric(5)),
			expected: -1,
		},
		{
			name:     "equal metric results in tie",
			a:        New(WithMetric(50)),
			b:        New(WithMetric(50)),
			expected: 0,
		},
		{
			name: "full comparison - all equal",
			a: New(
				WithWeight(100),
				WithLocalPreference(200),
				WithAdminCost(10),
				WithAsPath([]uint32{1, 2}),
				WithMetric(50),
			),
			b: New(
				WithWeight(100),
				WithLocalPreference(200),
				WithAdminCost(10),
				WithAsPath([]uint32{3, 4}),
				WithMetric(50),
			),
			expected: 0,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := CompareMultipath(tt.a, tt.b)
			if result != tt.expected {
				t.Errorf("CompareMultipath() = %d, expected %d", result, tt.expected)
			}
		})
	}
}

func TestCompareTieBreak(t *testing.T) {
	tests := []struct {
		name     string
		a        *BgpRoute
		b        *BgpRoute
		expected int
	}{
		{
			name:     "both nil",
			a:        nil,
			b:        nil,
			expected: 0,
		},
		{
			name:     "a is nil",
			a:        nil,
			b:        New(),
			expected: 1,
		},
		{
			name:     "b is nil",
			a:        New(),
			b:        nil,
			expected: -1,
		},
		{
			name: "prefer earlier arrival - a wins",
			a: func() *BgpRoute {
				r := New()
				r.SetArrival(1000)
				return r
			}(),
			b: func() *BgpRoute {
				r := New()
				r.SetArrival(2000)
				return r
			}(),
			expected: -1,
		},
		{
			name: "prefer earlier arrival - b wins",
			a: func() *BgpRoute {
				r := New()
				r.SetArrival(3000)
				return r
			}(),
			b: func() *BgpRoute {
				r := New()
				r.SetArrival(1000)
				return r
			}(),
			expected: 1,
		},
		{
			name: "same arrival - compare RxFrom - Local < IP",
			a: func() *BgpRoute {
				r := New(WithReceivedFrom(NewRxFrom(WithLocal())))
				r.SetArrival(1000)
				return r
			}(),
			b: func() *BgpRoute {
				r := New(WithReceivedFrom(NewRxFrom(WithIP(netip.MustParseAddr("10.0.0.1")))))
				r.SetArrival(1000)
				return r
			}(),
			expected: -1,
		},
		{
			name: "same arrival - compare RxFrom - IP < Interface",
			a: func() *BgpRoute {
				r := New(WithReceivedFrom(NewRxFrom(WithIP(netip.MustParseAddr("10.0.0.1")))))
				r.SetArrival(1000)
				return r
			}(),
			b: func() *BgpRoute {
				r := New(WithReceivedFrom(NewRxFrom(WithInterface("eth0"))))
				r.SetArrival(1000)
				return r
			}(),
			expected: -1,
		},
		{
			name: "same arrival - compare RxFrom - Local < Interface",
			a: func() *BgpRoute {
				r := New(WithReceivedFrom(NewRxFrom(WithLocal())))
				r.SetArrival(1000)
				return r
			}(),
			b: func() *BgpRoute {
				r := New(WithReceivedFrom(NewRxFrom(WithInterface("eth0"))))
				r.SetArrival(1000)
				return r
			}(),
			expected: -1,
		},
		{
			name: "same arrival and type - compare IP addresses",
			a: func() *BgpRoute {
				r := New(WithReceivedFrom(NewRxFrom(WithIP(netip.MustParseAddr("10.0.0.1")))))
				r.SetArrival(1000)
				return r
			}(),
			b: func() *BgpRoute {
				r := New(WithReceivedFrom(NewRxFrom(WithIP(netip.MustParseAddr("10.0.0.2")))))
				r.SetArrival(1000)
				return r
			}(),
			expected: -1,
		},
		{
			name: "same arrival and type - compare IP addresses - b wins",
			a: func() *BgpRoute {
				r := New(WithReceivedFrom(NewRxFrom(WithIP(netip.MustParseAddr("192.168.1.1")))))
				r.SetArrival(1000)
				return r
			}(),
			b: func() *BgpRoute {
				r := New(WithReceivedFrom(NewRxFrom(WithIP(netip.MustParseAddr("10.0.0.1")))))
				r.SetArrival(1000)
				return r
			}(),
			expected: 1,
		},
		{
			name: "same arrival and type - compare interface names",
			a: func() *BgpRoute {
				r := New(WithReceivedFrom(NewRxFrom(WithInterface("eth0"))))
				r.SetArrival(1000)
				return r
			}(),
			b: func() *BgpRoute {
				r := New(WithReceivedFrom(NewRxFrom(WithInterface("eth1"))))
				r.SetArrival(1000)
				return r
			}(),
			expected: -1,
		},
		{
			name: "same arrival and type - compare interface names - b wins",
			a: func() *BgpRoute {
				r := New(WithReceivedFrom(NewRxFrom(WithInterface("eth2"))))
				r.SetArrival(1000)
				return r
			}(),
			b: func() *BgpRoute {
				r := New(WithReceivedFrom(NewRxFrom(WithInterface("eth0"))))
				r.SetArrival(1000)
				return r
			}(),
			expected: 1,
		},
		{
			name: "same arrival and same RxFrom - tie",
			a: func() *BgpRoute {
				r := New(WithReceivedFrom(NewRxFrom(WithLocal())))
				r.SetArrival(1000)
				return r
			}(),
			b: func() *BgpRoute {
				r := New(WithReceivedFrom(NewRxFrom(WithLocal())))
				r.SetArrival(1000)
				return r
			}(),
			expected: 0,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := CompareTieBreak(tt.a, tt.b)
			if result != tt.expected {
				t.Errorf("CompareTieBreak() = %d, expected %d", result, tt.expected)
			}
		})
	}
}
