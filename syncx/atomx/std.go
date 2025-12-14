//go:build !race

package atomx

import "sync/atomic"

func maybeYield() {}

type (
	Bool    = atomic.Bool
	Int32   = atomic.Int32
	Int64   = atomic.Int64
	Uint32  = atomic.Uint32
	Uint64  = atomic.Uint64
	Uintptr = atomic.Uintptr
	Value   = atomic.Value
)
