package lockfree

import (
	"testing"

	"github.com/nikmy/algo/testx/faulty"
	"github.com/nikmy/algo/testx/synctest"
)

func TestPipe_Safety(t *testing.T) {
	pipe := NewPipe[int](256)

	consume := synctest.Operation{
		Runner: func() { pipe.TryConsume(new(int)) },
		Actors: 1,
	}

	produce := synctest.Operation{
		Runner: func() { pipe.TryProduce(42) },
		Actors: 1,
	}

	c := faulty.NewController(t, 42)
	c.SetFaultProbability(0.3)

	synctest.Stress(t, c.FaultInjector, 1_000_000, produce, consume)
}

func BenchmarkPipe_CompareWithChannel(b *testing.B) {
	compareParSeq(b, func(size uint32) bufferImpl { return NewPipe[int](size) })
}
