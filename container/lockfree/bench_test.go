package lockfree

import (
	"runtime"
	"testing"
	"time"

	"github.com/nikmy/algo/syncx"
	"github.com/nikmy/algo/testx/faulty"
)

type chanBuf chan int

func (c chanBuf) TryProduce(x int) bool {
	select {
	case c <- x:
		return true
	default:
		return false
	}
}

func (c chanBuf) TryConsume(x *int) bool {
	select {
	case *x = <-c:
		return true
	default:
		return false
	}
}

func compareParSeq(t testing.TB, init func(uint32) bufferImpl) {
	c := faulty.NewController(t, 42)

	c.Parallel(false)
	t.Logf("---- Sequential Access ----")
	compareWithChannel(t, init)

	t.Logf("\n\n")

	c.Parallel(true)
	t.Logf("---- Parallel Access ----")
	compareWithChannel(t, init)
}

func compareWithChannel(t testing.TB, init func(uint32) bufferImpl) {
	for _, bSize := range []uint32{1, 16, 64, 256, 1024, 4096} {
		ch, pipe := make(chanBuf, bSize), init(bSize)
		chP, chC, chNP, chnC := benchBuffer(ch)
		implP, implC, implNP, implNC := benchBuffer(pipe)

		pDiff := implP - chP
		cDiff := implC - chC

		npDiff := (implNP - chNP) * 1000 / chNP
		ncDiff := (implNC - chnC) * 1000 / chnC

		t.Logf("Buffer Size:  %d\n", bSize)
		t.Logf("\tProduced:   %+d\n", npDiff)
		t.Logf("\tConsumed:   %+d\n", ncDiff)
		t.Logf("\tTryProduce: %+.1f %%\n", float64(pDiff)/float64(chP)*100)
		t.Logf("\tTryConsume: %+.1f %%\n", float64(cDiff)/float64(chC)*100)
	}
}

type bufferImpl interface {
	TryProduce(int) bool
	TryConsume(*int) bool
}

func benchBuffer(buf bufferImpl) (time.Duration, time.Duration, int, int) {
	var (
		elapsedC time.Duration
		elapsedP time.Duration

		produced int
		consumed int
	)

	var wg syncx.WaitGroup

	const iterations = 1_000_000

	nProc := runtime.GOMAXPROCS(0)

	wg.Add(1)
	go func() {
		defer wg.Done()
		start := time.Now()
		for i := 0; i < iterations; i++ {
			if buf.TryProduce(i) {
				produced++
			} else if nProc == 1 {
				runtime.Gosched()
			}
		}
		elapsedP = time.Since(start)
	}()

	wg.Add(1)
	go func() {
		defer wg.Done()
		start := time.Now()
		for i := 0; i < iterations; i++ {
			var x int
			if buf.TryConsume(&x) {
				consumed++
			} else if nProc == 1 {
				runtime.Gosched()
			}
			_ = x
		}
		elapsedC = time.Since(start)
	}()

	wg.Wait()

	return elapsedP, elapsedC, produced, consumed
}
