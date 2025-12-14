package lockfree

import (
	"testing"

	"github.com/nikmy/algo/syncx"
	"github.com/nikmy/algo/testx/faulty"
	"github.com/nikmy/algo/testx/synctest"
)

func TestStack_Safety(t *testing.T) {
	var stack Stack[int]

	produce := synctest.Operation{
		Runner: func() { stack.TryPush(42) },
		Actors: 2,
	}

	consume := synctest.Operation{
		Runner: func() { stack.TryPop(new(int)) },
		Actors: 2,
	}

	c := faulty.NewController(t, 42)
	c.SetFaultProbability(0.2)

	synctest.Stress(t, c, 1_000_000, produce, consume)
}

func TestStack_LIFO(t *testing.T) {
	c := faulty.NewController(t, 42)

	c.Parallel(false)

	order := c.Fuzzed().Perm(1_000_000)

	var (
		stack Stack[int]
		valid []int
		mutex syncx.Mutex
	)

	c.SetFaultProbability(1)

	c.Go(func() {
		for x := range order {
			full := !stack.TryPush(x)
			if !full {
				mutex.Lock()
				valid = append(valid, x)
				mutex.Unlock()
			}
			if full || c.Fault() {
				c.Yield()
			}
		}
	})

	c.Go(func() {
		var pop int

		for range order {
			empty := !stack.TryPop(&pop)
			if !empty {
				mutex.Lock()
				if len(valid) == 0 {
					mutex.Unlock()
					c.Fail("missed items %v", valid)
					return
				}

				if valid[len(valid)-1] != pop {
					mutex.Unlock()
					c.Fail("LIFO property violation", "got item %d, want %d", pop, valid[len(valid)-1])
					return
				}

				valid = valid[:len(valid)-1]
				mutex.Unlock()
			}

			if empty || c.Fault() {
				c.Yield()
			}
		}
	})

	c.Wait()
}
