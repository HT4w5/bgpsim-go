package prefixtrie_test

import (
	"testing"

	"github.com/HT4w5/bgpsim-go/pkg/prefixtrie"
)

func TestStackPushPop(t *testing.T) {
	input := [...]int{1, 2, 3}
	expected := [...]int{3, 2, 1}

	s := prefixtrie.NewStack[int](3)
	for _, v := range input {
		s.Push(v)
	}

	for _, v := range expected {
		if p, ok := s.Pop(); !ok {
			t.Error("failed to Pop()")
		} else if p != v {
			t.Errorf("Pop(): got %d, expected %d", p, v)
		}
	}
}

func TestStackLen(t *testing.T) {
	s := prefixtrie.NewStack[int](100)
	for i := range 100 {
		s.Push(i)
		l := s.Len()
		if l != i+1 {
			t.Errorf("Len(): got %d, expected %d", l, i+1)
		}
	}
}

func TestStack1k(t *testing.T) {
	s := prefixtrie.NewStack[int](1000)
	for i := range 1000 {
		s.Push(i)
	}

	for i := 999; i >= 0; i-- {
		if v, ok := s.Pop(); !ok {
			t.Error("failed to Pop()")
		} else if i != v {
			t.Errorf("Pop(): got %d, expected %d", v, i)
		}
	}
}

func TestStack10k(t *testing.T) {
	s := prefixtrie.NewStack[int](10000)
	for i := range 10000 {
		s.Push(i)
	}

	for i := 9999; i >= 0; i-- {
		if v, ok := s.Pop(); !ok {
			t.Error("failed to Pop()")
		} else if i != v {
			t.Errorf("Pop(): got %d, expected %d", v, i)
		}
	}
}

func BenchmarkStack1k(b *testing.B) {
	for b.Loop() {
		s := prefixtrie.NewStack[int](1000)
		for i := range 1000 {
			s.Push(i)
		}
		for range 1000 {
			s.Pop()
		}
	}
}
