package prefixtrie

type Stack[T any] struct {
	elements []T
}

func NewStack[T any](alloc int) *Stack[T] {
	return &Stack[T]{
		elements: make([]T, 0, alloc),
	}
}

func (s *Stack[T]) Push(e T) {
	s.elements = append(s.elements, e)
}

func (s *Stack[T]) Pop() (T, bool) {
	size := len(s.elements)
	if size == 0 {
		var zero T
		return zero, false
	}

	e := s.elements[size-1]

	s.elements = s.elements[:size-1]

	return e, true
}

func (s *Stack[T]) Peek() (T, bool) {
	size := len(s.elements)
	if size == 0 {
		var zero T
		return zero, false
	}

	return s.elements[size-1], true
}

func (s *Stack[T]) Len() int {
	return len(s.elements)
}
