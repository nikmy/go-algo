package ptest

import (
	"context"
	"runtime"
	"sync"
	"sync/atomic"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
)

func NewController(t *testing.T) *pController {
	return &pController{
		t:  t,
		do: make(chan struct{}),
	}
}

type pController struct {
	t  *testing.T
	wg sync.WaitGroup
	do chan struct{}
	n  atomic.Int64
}

func (c *pController) Run(timeout time.Duration) {
	ctx, cancel := context.WithTimeout(context.Background(), timeout)
	defer cancel()

	close(c.do)
	done := make(chan struct{})
	go func() {
		defer close(done)
		c.wg.Wait()
	}()
	select {
	case <-done:
	case <-ctx.Done():
		stack := make([]byte, 100_000)
		n := runtime.Stack(stack, true)
		assert.Failf(c.t, "timed out", "%d goroutines stuck:\n--- stacktrace ---\n%s", c.n.Load(), stack[:n])
	}
}

func (c *pController) Spawn(n int, g func()) {
	c.wg.Add(n)
	c.n.Add(int64(n))
	for i := 0; i < n; i++ {
		go func() {
			defer c.wg.Done()
			<-c.do
			g()
			c.n.Add(-1)
		}()
	}
}
