package syncx

import (
	"sync"

	"github.com/nikmy/algo/syncx/atomx"
)

// NewBarrier creates new Barrier primitive
// for synchronizing n goroutines.
func NewBarrier(n int) *Barrier {
	b := &Barrier{
		allDone: NewCond(new(sync.Mutex)),
	}
	b.counter.Store(int64(n))
	return b
}

// Barrier is synchronization primitive
// to make sure that all workers hit
// checkpoint.
//
// Typical use case:
// 	ready := NewBarrier(n)
// 	for i := range n {
// 	    go func() {
//          prepare()
//	        ready.Pass()
//          run()
//      }
// 	}
type Barrier struct {
	counter atomx.Int64
	allDone *Cond
}

// Pass decrements the counter and suspends till all
// Pass will be called. All operations after Pass calls
// will be synchronized with happens-before relationship
// with all operations before Pass calls.
func (b *Barrier) Pass() {
	if b.counter.Load() == 0 {
		panic("invalid use of barrier")
	}

	if b.counter.Add(-1) == 0 {
		b.allDone.Broadcast()
		return
	}

	b.allDone.L.Lock()
	defer b.allDone.L.Unlock()
	for b.counter.Load() > 0 {
		b.allDone.Wait()
	}
}

func (b *Barrier) Done() bool {
	return b.counter.Load() == 0
}
