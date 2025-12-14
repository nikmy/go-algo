package faulty

import (
	"math/rand"
	"sync"
)

type FaultInjector struct {
	rnd *rand.Rand
	thr float64
	mu  sync.Mutex
}

func (r *FaultInjector) SetFaultProbability(p float64) {
	r.mu.Lock()
	defer r.mu.Unlock()
	r.thr = p
}

func (r *FaultInjector) Fault() bool {
	r.mu.Lock()
	defer r.mu.Unlock()
	return r.rnd.Float64() < r.thr
}
