//go:build race

package syncx

import (
	"sync"

	"github.com/nikmy/algo/fp"
	"github.com/nikmy/algo/syncx/internal/fuzz"
)

type (
	WaitGroup = sync.WaitGroup
	Locker    = sync.Locker
	Once      = sync.Once

	Mutex   sync.Mutex
	RWMutex sync.RWMutex

	privateCond = *sync.Cond
	Cond        struct {
		privateCond
	}
)

func (m *Mutex) impl() *sync.Mutex {
	return fp.UnsafeCast[Mutex, sync.Mutex](m)
}

func (m *Mutex) Lock() {
	fuzz.MaybeYield()
	m.impl().Lock()
}

func (m *Mutex) Unlock() {
	fuzz.MaybeYield()
	m.impl().Unlock()
}

func (m *Mutex) TryLock() bool {
	fuzz.MaybeYield()
	return m.impl().TryLock()
}

func (m *RWMutex) impl() *sync.RWMutex {
	return fp.UnsafeCast[RWMutex, sync.RWMutex](m)
}

func (m *RWMutex) Lock() {
	fuzz.MaybeYield()
	m.impl().Lock()
}

func (m *RWMutex) Unlock() {
	fuzz.MaybeYield()
	m.impl().Unlock()
}

func (m *RWMutex) TryLock() bool {
	fuzz.MaybeYield()
	return m.impl().TryLock()
}

func (m *RWMutex) RLock() {
	fuzz.MaybeYield()
	m.impl().RLock()
}

func (m *RWMutex) RUnlock() {
	fuzz.MaybeYield()
	m.impl().RUnlock()
}

func (m *RWMutex) TryRLock() bool {
	fuzz.MaybeYield()
	return m.impl().TryRLock()
}

func (m *RWMutex) RLocker() Locker {
	return m.impl().RLocker()
}

func NewCond(l Locker) *Cond {
	return &Cond{privateCond: sync.NewCond(l)}
}

func (c *Cond) Signal() {
	fuzz.MaybeYield()
	c.privateCond.Signal()
}

func (c *Cond) Broadcast() {
	fuzz.MaybeYield()
	c.privateCond.Broadcast()
}

func (c *Cond) Wait() {
	fuzz.MaybeYield()
	c.privateCond.Wait()
}
