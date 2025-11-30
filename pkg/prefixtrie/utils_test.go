package prefixtrie

import (
	"fmt"
	"testing"
)

func TestUint32LCPL(t *testing.T) {
	tests := []struct {
		a    uint32
		b    uint32
		want int
	}{
		{0, 0, 32},
		{0, ^uint32(0), 0},
		{^uint32(0), ^uint32(0), 32},
		{255, 511, 23},
		{2147483648, 3221225472, 1},
	}

	for _, tt := range tests {
		t.Run(fmt.Sprintf("%d,%d", tt.a, tt.b), func(t *testing.T) {
			ans := uint32LCPL(tt.a, tt.b)
			if ans != tt.want {
				t.Errorf("got %d, want %d", ans, tt.want)
			}
		})
	}
}
