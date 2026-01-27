package optional

type Optional[T any] struct {
	value T
	ok    bool
}

func Of[T any](v T) Optional[T] {
	return Optional[T]{
		value: v,
		ok:    true,
	}
}

func Empty[T any]() Optional[T] {
	return Optional[T]{
		ok: false,
	}
}

func (o Optional[T]) IsValid() bool {
	return o.ok
}

func (o Optional[T]) Get() (T, bool) {
	return o.value, o.ok
}

func (o Optional[T]) OrElse(fallback T) T {
	if !o.ok {
		return fallback
	}
	return o.value
}
