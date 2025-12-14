package buffer

import (
	"math/rand"
	stdslices "slices"
	"testing"
	"time"

	"github.com/stretchr/testify/require"
	"go.uber.org/goleak"

	"github.com/nikmy/algo/syncx"
)

type bufferKind struct {
	Size int
	Name string
}

var allSizes = []bufferKind{
	{Size: tinySize, Name: "512B"},
	{Size: smallSize, Name: "1KB"},
	{Size: mediumSize, Name: "4KB"},
	{Size: largeSize, Name: "1MB"},
	{Size: hugeSize, Name: "8MB"},
}

func TestShardedPool_SizeMatch(t *testing.T) {
	p := NewPool(mediumSize, FileReadingPreset)

	t.Run("resize to bigger", func(t *testing.T) {
		defer p.Release()

		small := p.Get(smallSize)
		small.Resize(mediumSize)
		small.Free()

		require.Empty(t, p.getPoolForSize(smallSize).buffers)

		mediumPoolBuffers := p.getPoolForSize(mediumSize).buffers
		require.Len(t, mediumPoolBuffers, 1)
		require.Equal(t, mediumSize, cap(mediumPoolBuffers[0]))
	})

	t.Run("resize to smaller", func(t *testing.T) {
		defer p.Release()

		buf := p.Get(mediumSize)
		buf.Resize(smallSize)
		buf.Free()

		require.Empty(t, p.getPoolForSize(smallSize).buffers)

		mediumPoolBuffers := p.getPoolForSize(mediumSize).buffers
		require.Len(t, mediumPoolBuffers, 1)
		require.Equal(t, mediumSize, cap(mediumPoolBuffers[0]))
	})

	t.Run("second free do nothing", func(t *testing.T) {
		defer p.Release()

		buf := p.Get(tinySize)
		buf.Free()
		buf.Free()

		for _, kind := range allSizes {
			requiredLen := 0
			if kind.Size == tinySize {
				requiredLen = 1
			}

			pool := p.getPoolForSize(kind.Size)

			require.Len(t, pool.buffers, requiredLen, "mismatch for %s size", kind.Name)
		}
	})
}

func TestPool(t *testing.T) {
	defer goleak.VerifyNone(t)

	var exactSizes []int
	for _, kind := range allSizes {
		exactSizes = append(exactSizes, kind.Size)
	}
	exactSizes = append(exactSizes, largeSize*2, largeSize*4)

	sizes := []int{0, 1}
	for _, size := range exactSizes {
		sizes = append(sizes, size-1, size, size+1)
	}

	p := NewPool(largeSize*2, FileReadingPreset)

	start := syncx.NewBarrier(33)

	var wg syncx.WaitGroup
	for i := int64(0); i < 32; i++ {
		src := rand.NewSource(i)

		order := stdslices.Clone(sizes)
		rand.New(src).Shuffle(len(order), func(i, j int) {
			order[i], order[j] = order[j], order[i]
		})

		wg.Add(1)
		go func() {
			start.Pass()
			defer wg.Done()

			for _, s := range order {
				buf := p.Get(s)
				require.Equal(t, s, len(buf.Data()))
				buf.Free()
			}
		}()
	}

	validateSizedPool := func(pool *sizedPool) {
		pool.Lock()
		defer pool.Unlock()

		require.LessOrEqual(t, len(pool.buffers), pool.maxPoolSize)
		for _, buf := range pool.buffers {
			require.Equal(t, pool.bufferSize, cap(buf))
		}
	}

	wg.Add(1)
	go func() {
		start.Pass()
		defer wg.Done()

		for i := 0; i < 10000; i++ {
			time.Sleep(time.Nanosecond * 10)
			for i := range *p {
				validateSizedPool(&(*p)[i])
			}
		}
	}()

	wg.Wait()
}
