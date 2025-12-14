package faulty

import (
	"math/rand"
	"testing"

	"github.com/stretchr/testify/require"
)

func NewController(t testing.TB, randomSeed int64) *Controller {
	c := &Controller{
		t:             t,
		rnd:           rand.New(rand.NewSource(randomSeed)),
		Async:         new(Async),
		FaultInjector: new(FaultInjector),
	}

	c.FaultInjector.rnd = c.rnd

	return c
}

type Controller struct {
	t testing.TB

	rnd *rand.Rand

	*Async
	*FaultInjector
}

func (c *Controller) Fuzzed() *rand.Rand {
	return c.rnd
}

func (c *Controller) Seed(seed int64) {
	c.rnd.Seed(seed)
}

func (c *Controller) Fail(msg string, args ...any) {
	require.Fail(c.t, msg, args...)
}
