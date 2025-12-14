package syncx

import (
	"github.com/nikmy/algo/syncx/atomx"
)

func NewSemaphore(limit int64) *Semaphore {
	return &Semaphore{owners: limit}
}

type Semaphore struct{
	owners int64
}

func (s *Semaphore) TryAcquire(n int64) bool {
	current := atomx.LoadInt64(&s.owners)
	if current < n {
		return false
	}

	return atomx.CompareAndSwapInt64(&s.owners, current, current-n)
}

func (s *Semaphore) Release(n int64) {
	atomx.AddInt64(&s.owners, n)
}
