package buffer

import (
	"fmt"
	"runtime"
	"testing"
)

func BenchmarkSizedPool(b *testing.B) {
	for _, kind := range allSizes {
		b.Run(fmt.Sprintf("%s buffers from pool", kind.Name), func(b *testing.B) {
			b.Cleanup(runtime.GC)
			p := sizedPool{
				bufferSize:  kind.Size,
				maxPoolSize: 1,
			}

			b.ResetTimer()
			for i := 0; i < b.N; i++ {
				buf := p.get(kind.Size)
				p.put(buf)
			}
			b.StopTimer()
		})

		b.Run(fmt.Sprintf("%s buffers from malloc", kind.Name), func(b *testing.B) {
			b.Cleanup(runtime.GC)
			for i := 0; i < b.N; i++ {
				_ = make([]byte, kind.Size)
			}
		})
	}
}

func BenchmarkSizedPool_Parallel(b *testing.B) {
	for _, parallel := range []int{1, 4, 16} {
		b.Logf("parallel=%d", parallel)
		for _, kind := range allSizes {
			b.Run(
				fmt.Sprintf("[parallel=%d] %s buffers from pool", parallel, kind.Name),
				func(b *testing.B) {
					p := sizedPool{
						bufferSize:  kind.Size,
						maxPoolSize: parallel,
					}

					b.Cleanup(runtime.GC)
					b.SetParallelism(parallel)

					b.ResetTimer()
					b.RunParallel(func(pb *testing.PB) {
						b.ReportAllocs()
						for pb.Next() {
							buf := p.get(kind.Size)
							p.put(buf)
						}
					})
					b.StopTimer()
				},
			)

			b.Run(
				fmt.Sprintf("[parallel=%d] %s buffers from malloc", parallel, kind.Name),
				func(b *testing.B) {
					b.Cleanup(runtime.GC)
					b.SetParallelism(parallel)

					b.ResetTimer()
					b.RunParallel(func(pb *testing.PB) {
						b.ReportAllocs()
						for pb.Next() {
							_ = make([]byte, kind.Size)
						}
					})
					b.StopTimer()
				},
			)
		}
	}
}
