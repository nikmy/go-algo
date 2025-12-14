package atomx

import (
	"sync/atomic"

	"github.com/nikmy/algo/fp"
	"github.com/nikmy/algo/syncx/internal/fuzz"
)

/*
	Since pointer is generic, it cannot be aliased, so we need to
	use unsafe cast always. Hope it's not so expensive.
*/

type Pointer[T any] atomic.Pointer[T]

func (x *Pointer[T]) impl() *atomic.Pointer[T] {
	return fp.UnsafeCast[Pointer[T], atomic.Pointer[T]](x)
}

func (x *Pointer[T]) Load() *T {
	fuzz.MaybeYield()
	return x.impl().Load()
}

func (x *Pointer[T]) Store(val *T) {
	fuzz.MaybeYield()
	x.impl().Store(val)
}

func (x *Pointer[T]) Swap(new *T) (old *T) {
	fuzz.MaybeYield()
	return x.impl().Swap(new)
}

func (x *Pointer[T]) CompareAndSwap(old, new *T) (swapped bool) {
	fuzz.MaybeYield()
	return x.impl().CompareAndSwap(old, new)
}
