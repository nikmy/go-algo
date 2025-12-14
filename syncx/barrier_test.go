package syncx

import (
	"testing"

	"github.com/nikmy/algo/testx/faulty"
	"github.com/nikmy/algo/testx/synctest"
)

func TestBarrier_Pass(t *testing.T) {
	const n = 100

	barrier := NewBarrier(n)
	synctest.Stress(t, faulty.NewController(t, 42), 1,
		synctest.Operation{
			Actors: n,
			Runner: barrier.Pass,
		},
	)
}
