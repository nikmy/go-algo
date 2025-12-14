package atomx

import (
	"sync/atomic"
	"unsafe"

	"github.com/nikmy/algo/syncx/internal/fuzz"
)

func SwapInt32(addr *int32, new int32) (old int32) {
	fuzz.MaybeYield()
	return atomic.SwapInt32(addr, new)
}

func SwapInt64(addr *int64, new int64) (old int64) {
	fuzz.MaybeYield()
	return atomic.SwapInt64(addr, new)
}

func SwapUint32(addr *uint32, new uint32) (old uint32) {
	fuzz.MaybeYield()
	return atomic.SwapUint32(addr, new)
}

func SwapUint64(addr *uint64, new uint64) (old uint64) {
	fuzz.MaybeYield()
	return atomic.SwapUint64(addr, new)
}

func SwapUintptr(addr *uintptr, new uintptr) (old uintptr) {
	fuzz.MaybeYield()
	return atomic.SwapUintptr(addr, new)
}

func SwapPointer(addr *unsafe.Pointer, new unsafe.Pointer) (old unsafe.Pointer) {
	fuzz.MaybeYield()
	return atomic.SwapPointer(addr, new)
}

func CompareAndSwapInt32(addr *int32, old, new int32) (swapped bool) {
	fuzz.MaybeYield()
	return atomic.CompareAndSwapInt32(addr, old, new)
}

func CompareAndSwapInt64(addr *int64, old, new int64) (swapped bool) {
	fuzz.MaybeYield()
	return atomic.CompareAndSwapInt64(addr, old, new)
}

func CompareAndSwapUint32(addr *uint32, old, new uint32) (swapped bool) {
	fuzz.MaybeYield()
	return atomic.CompareAndSwapUint32(addr, old, new)
}

// CompareAndSwapUint64 executes the compare-and-swap operation for a uint64 value.
// Consider using the more ergonomic and less error-prone [Uint64.CompareAndSwap] instead
// (particularly if you target 32-bit platforms; see the bugs section).
func CompareAndSwapUint64(addr *uint64, old, new uint64) (swapped bool) {
	fuzz.MaybeYield()
	return atomic.CompareAndSwapUint64(addr, old, new)
}

func CompareAndSwapUintptr(addr *uintptr, old, new uintptr) (swapped bool) {
	fuzz.MaybeYield()
	return atomic.CompareAndSwapUintptr(addr, old, new)
}

func CompareAndSwapPointer(addr *unsafe.Pointer, old, new unsafe.Pointer) (swapped bool) {
	fuzz.MaybeYield()
	return atomic.CompareAndSwapPointer(addr, old, new)
}

func AddInt32(addr *int32, delta int32) (new int32) {
	fuzz.MaybeYield()
	return atomic.AddInt32(addr, delta)
}

func AddUint32(addr *uint32, delta uint32) (new uint32) {
	fuzz.MaybeYield()
	return atomic.AddUint32(addr, delta)
}

func AddInt64(addr *int64, delta int64) (new int64) {
	fuzz.MaybeYield()
	return atomic.AddInt64(addr, delta)
}

func AddUint64(addr *uint64, delta uint64) (new uint64) {
	fuzz.MaybeYield()
	return atomic.AddUint64(addr, delta)
}

func AddUintptr(addr *uintptr, delta uintptr) (new uintptr) {
	fuzz.MaybeYield()
	return atomic.AddUintptr(addr, delta)
}

func LoadInt32(addr *int32) (val int32) {
	fuzz.MaybeYield()
	return atomic.LoadInt32(addr)
}

func LoadInt64(addr *int64) (val int64) {
	fuzz.MaybeYield()
	return atomic.LoadInt64(addr)
}

func LoadUint32(addr *uint32) (val uint32) {
	fuzz.MaybeYield()
	return atomic.LoadUint32(addr)
}

func LoadUint64(addr *uint64) (val uint64) {
	fuzz.MaybeYield()
	return atomic.LoadUint64(addr)
}

func LoadUintptr(addr *uintptr) (val uintptr) {
	fuzz.MaybeYield()
	return atomic.LoadUintptr(addr)
}

func LoadPointer(addr *unsafe.Pointer) (val unsafe.Pointer) {
	fuzz.MaybeYield()
	return atomic.LoadPointer(addr)
}

func StoreInt32(addr *int32, val int32) {
	fuzz.MaybeYield()
	atomic.StoreInt32(addr, val)
}

func StoreInt64(addr *int64, val int64) {
	fuzz.MaybeYield()
	atomic.StoreInt64(addr, val)
}

func StoreUint32(addr *uint32, val uint32) {
	fuzz.MaybeYield()
	atomic.StoreUint32(addr, val)
}

func StoreUint64(addr *uint64, val uint64) {
	fuzz.MaybeYield()
	atomic.StoreUint64(addr, val)
}

func StoreUintptr(addr *uintptr, val uintptr) {
	fuzz.MaybeYield()
	atomic.StoreUintptr(addr, val)
}

func StorePointer(addr *unsafe.Pointer, val unsafe.Pointer) {
	fuzz.MaybeYield()
	atomic.StorePointer(addr, val)
}
