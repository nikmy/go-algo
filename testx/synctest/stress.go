package synctest

import (
	"runtime"
	"sync"
	"testing"

	"github.com/stretchr/testify/require"
)

type Operation struct {
	Runner func()
	Actors int
}

func Stress(t testing.TB, f faultInjector, iters int, operations ...Operation) {
	m.Lock()
	defer m.Unlock()
	F = f
	require.NotPanics(t, func() {
		defer runtime.GOMAXPROCS(runtime.GOMAXPROCS(1))
		stress(t, iters, operations...)
	})
}

func stress(t testing.TB, iters int, operations ...Operation) {
	runtime.LockOSThread()
	defer runtime.UnlockOSThread()

	var (
		group = sync.WaitGroup{}
		start = make(chan struct{})
	)

	for _, op := range operations {
		f := op.Runner
		for i := 0; i < op.Actors; i++ {
			group.Add(1)
			go func() {
				runtime.LockOSThread()
				defer runtime.UnlockOSThread()

				<-start

				defer group.Done()

				defer func() {
					if r := recover(); r != nil {
						t.Helper()
						t.Fatalf("panic: %v", r)
					}
				}()

				for i := 0; i < iters; i++ {
					f()
				}
			}()
		}
	}

	close(start)
	group.Wait()
}
