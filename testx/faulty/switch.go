package faulty

import (
	"runtime"
	"sync"
)

type Async struct {
	g sync.WaitGroup
}

func (c *Async) Parallel(enableParallelism bool) {
	if enableParallelism {
		runtime.GOMAXPROCS(runtime.NumCPU())
	} else {
		runtime.GOMAXPROCS(1)
	}
}

func (c *Async) Yield() {
	runtime.Gosched()
}

func (c *Async) Go(f func()) {
	c.g.Add(1)
	go func() {
		defer c.g.Done()
		f()
	}()
}

func (c *Async) Wait() {
	c.g.Wait()
}
