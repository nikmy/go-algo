package faulty

import (
	"testing"

	"github.com/stretchr/testify/require"
)

func TestDeterministic(t *testing.T) {
	c1 := NewController(t, 0)
	c2 := NewController(t, 0)

	require.Equal(t, c1.Fuzzed().Perm(1_000), c2.Fuzzed().Perm(1_000))

	c1.SetFaultProbability(0.35)
	c2.SetFaultProbability(0.35)
	for i := 0; i < 100; i++ {
		require.Equal(t, c1.Fault(), c2.Fault())
	}
}

func TestAsync(t *testing.T) {
	c := NewController(t, 0)

	c.Parallel(false)

	ping, pong, end := false, false, false
	c.Go(func() {
		defer func() { end = true }()

		require.True(t, ping)
		pong = true
		c.Yield()
	})

	require.False(t, pong)

	ping = true
	c.Yield()

	require.True(t, pong)
	require.False(t, end)

	c.Wait()

	require.True(t, end)
}
