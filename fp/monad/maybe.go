package monad

func Just[T any](value T) Maybe[T] { return Maybe[T]{ptr: &value} }
func Nothing[T any]() Maybe[T]     { return Maybe[T]{} }

// Maybe represent the object that can be absent.
// It must always be checked with Ok method before
// calling Get.
type Maybe[T any] struct {
	ptr *T
}

func (m Maybe[T]) Ok() bool { return m.ptr != nil }
func (m Maybe[T]) Get() *T  { return m.ptr }
