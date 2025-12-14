//go:build !race

package syncx

import "sync"

type (
	WaitGroup = sync.WaitGroup
	Locker    = sync.Locker
	Once = sync.Once

	Mutex   = sync.Mutex
	RWMutex = sync.RWMutex
	Cond    = sync.Cond
)

var NewCond = sync.NewCond
