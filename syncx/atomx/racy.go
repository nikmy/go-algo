//go:build race

package atomx

import (
	"sync/atomic"

	"github.com/nikmy/algo/fp"
	"github.com/nikmy/algo/syncx/internal/fuzz"
)

type (
	Bool    atomic.Bool
	Int32   atomic.Int32
	Int64   atomic.Int64
	Uint32  atomic.Uint32
	Uint64  atomic.Uint64
	Uintptr atomic.Uintptr
	Value   atomic.Value
)

func (x *Bool) impl() *atomic.Bool {
	return fp.UnsafeCast[Bool, atomic.Bool](x)
}

func (x *Bool) Load() bool {
	fuzz.MaybeYield()
	return x.impl().Load()
}

func (x *Bool) Store(val bool) {
	fuzz.MaybeYield()
	x.impl().Store(val)
}

func (x *Bool) Swap(new bool) (old bool) {
	fuzz.MaybeYield()
	return x.impl().Swap(new)
}

func (x *Bool) CompareAndSwap(old, new bool) (swapped bool) {
	fuzz.MaybeYield()
	return x.impl().CompareAndSwap(old, new)
}

func (x *Int32) impl() *atomic.Int32 {
	return fp.UnsafeCast[Int32, atomic.Int32](x)
}

func (x *Int32) Add(delta int32) int32 {
	fuzz.MaybeYield()
	return x.impl().Add(delta)
}

func (x *Int32) Load() int32 {
	fuzz.MaybeYield()
	return x.impl().Load()
}

func (x *Int32) Store(val int32) {
	fuzz.MaybeYield()
	x.impl().Store(val)
}

func (x *Int32) Swap(new int32) (old int32) {
	fuzz.MaybeYield()
	return x.impl().Swap(new)
}

func (x *Int32) CompareAndSwap(old, new int32) (swapped bool) {
	fuzz.MaybeYield()
	return x.impl().CompareAndSwap(old, new)
}

func (x *Uint32) impl() *atomic.Uint32 {
	return fp.UnsafeCast[Uint32, atomic.Uint32](x)
}

func (x *Uint32) Add(delta uint32) uint32 {
	fuzz.MaybeYield()
	return x.impl().Add(delta)
}

func (x *Uint32) Load() uint32 {
	fuzz.MaybeYield()
	return x.impl().Load()
}

func (x *Uint32) Store(val uint32) {
	fuzz.MaybeYield()
	x.impl().Store(val)
}

func (x *Uint32) Swap(new uint32) (old uint32) {
	fuzz.MaybeYield()
	return x.impl().Swap(new)
}

func (x *Uint32) CompareAndSwap(old, new uint32) (swapped bool) {
	fuzz.MaybeYield()
	return x.impl().CompareAndSwap(old, new)
}

func (x *Int64) impl() *atomic.Int64 {
	return fp.UnsafeCast[Int64, atomic.Int64](x)
}

func (x *Int64) Add(delta int64) int64 {
	fuzz.MaybeYield()
	return x.impl().Add(delta)
}

func (x *Int64) Load() int64 {
	fuzz.MaybeYield()
	return x.impl().Load()
}

func (x *Int64) Store(val int64) {
	fuzz.MaybeYield()
	x.impl().Store(val)
}

func (x *Int64) Swap(new int64) (old int64) {
	fuzz.MaybeYield()
	return x.impl().Swap(new)
}

func (x *Int64) CompareAndSwap(old, new int64) (swapped bool) {
	fuzz.MaybeYield()
	return x.impl().CompareAndSwap(old, new)
}

func (x *Uint64) impl() *atomic.Uint64 {
	return fp.UnsafeCast[Uint64, atomic.Uint64](x)
}

func (x *Uint64) Add(delta uint64) uint64 {
	fuzz.MaybeYield()
	return x.impl().Add(delta)
}

func (x *Uint64) Load() uint64 {
	fuzz.MaybeYield()
	return x.impl().Load()
}

func (x *Uint64) Store(val uint64) {
	fuzz.MaybeYield()
	x.impl().Store(val)
}

func (x *Uint64) Swap(new uint64) (old uint64) {
	fuzz.MaybeYield()
	return x.impl().Swap(new)
}

func (x *Uint64) CompareAndSwap(old, new uint64) (swapped bool) {
	fuzz.MaybeYield()
	return x.impl().CompareAndSwap(old, new)
}

func (x *Uintptr) impl() *atomic.Uintptr {
	return fp.UnsafeCast[Uintptr, atomic.Uintptr](x)
}

func (x *Uintptr) Add(delta uintptr) uintptr {
	fuzz.MaybeYield()
	return x.impl().Add(delta)
}

func (x *Uintptr) Load() uintptr {
	fuzz.MaybeYield()
	return x.impl().Load()
}

func (x *Uintptr) Store(val uintptr) {
	fuzz.MaybeYield()
	x.impl().Store(val)
}

func (x *Uintptr) Swap(new uintptr) (old uintptr) {
	fuzz.MaybeYield()
	return x.impl().Swap(new)
}

func (x *Uintptr) CompareAndSwap(old, new uintptr) (swapped bool) {
	fuzz.MaybeYield()
	return x.impl().CompareAndSwap(old, new)
}

func (v *Value) impl() *atomic.Value {
	return fp.UnsafeCast[Value, atomic.Value](v)
}

func (v *Value) Load() any {
	fuzz.MaybeYield()
	return v.impl().Load()
}

func (v *Value) Store(val any) {
	fuzz.MaybeYield()
	v.impl().Store(val)
}

func (v *Value) Swap(new any) any {
	fuzz.MaybeYield()
	return v.impl().Swap(new)
}

func (v *Value) CompareAndSwap(old, new any) (swapped bool) {
	fuzz.MaybeYield()
	return v.impl().CompareAndSwap(old, new)
}
