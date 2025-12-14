package monad

func Error[T any](err error) Result[T] {
	return Result[T]{err: err}
}

func Value[T any](value T) Result[T] {
	return Result[T]{val: &value}
}

// Result is a monad for (value, error) pair.
// It can have only error or only value, and
// cannot be completely empty
type Result[T any] struct {
	val *T
	err error
}

func (r Result[T]) Ok() bool {
	return r.err == nil
}

func (r Result[T]) Err() error {
	return r.err
}

func (r Result[T]) Get() *T {
	return r.val
}
