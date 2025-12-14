package lockfree

import (
	"fmt"
	"math/rand/v2"
	"slices"
	"sync"
	"testing"
	"unsafe"
)

// Alternatives to compare with
// ///////////////////////////////////////////////////////////////////////////////////////////////

func Benchmark_BinarySearch(b *testing.B) {
	sizes := []int64{10, 100, 1_000, 10_000, 100_000, 1_000_000}
	parallel := []int{1, 4, 16, 32}
	for _, p := range parallel {
		for _, size := range sizes {
			b.Run(fmt.Sprintf("%d_parallel_%d_elements", p, size), func(b *testing.B) {
				b.ReportAllocs()
				b.SetParallelism(16)

				list := generateInts(0, int(size), 1)

				b.ResetTimer()
				b.RunParallel(func(pb *testing.PB) {
					i := rand.Int64() % size
					for pb.Next() {
						if _, ok := slices.BinarySearch(list, i); !ok {
							panic(fmt.Sprintf("%d must be found", i))
						}
						i += 7
						if i >= size {
							i -= size
						}
					}
				})
			})
		}
	}
}

func BenchmarkLock(b *testing.B) {
	b.SetParallelism(32)

	var lock sync.Mutex
	b.RunParallel(func(pb *testing.PB) {
		for pb.Next() {
			lock.Lock()
			lock.Unlock()
		}
	})
}

func BenchmarkSharedLock(b *testing.B) {
	b.SetParallelism(32)

	var lock sync.RWMutex
	b.RunParallel(func(pb *testing.PB) {
		seed := uint64(uintptr(unsafe.Pointer(pb)))
		mask := uint64(1<<32 - 1)
		rnd := rand.New(rand.NewPCG(seed&mask, seed^mask))
		for pb.Next() {
			if rnd.Int()%3 == 0 {
				lock.RLock()
				lock.RUnlock()
			} else {
				lock.Lock()
				lock.Unlock()
			}
		}
	})
}

/////////////////////////////////////////////////////////////////////////////////////////////////

func BenchmarkMakeSkipList(b *testing.B) {
	values := generateInts(-1_000, 1_000, 1)
	rand.Shuffle(len(values), func(i, j int) { values[i], values[j] = values[j], values[i] })

	b.Run("20_elements", func(b *testing.B) {
		b.ReportAllocs()
		b.SetParallelism(16)
		b.RunParallel(func(pb *testing.PB) {
			for pb.Next() {
				_ = MakeSkipList(values[:len(values)/100]...)
			}
		})
	})

	b.Run("200_elements", func(b *testing.B) {
		b.ReportAllocs()
		b.SetParallelism(16)
		b.RunParallel(func(pb *testing.PB) {
			for pb.Next() {
				_ = MakeSkipList(values[:len(values)/10]...)
			}
		})
	})

	b.Run("2000_elements", func(b *testing.B) {
		b.ReportAllocs()
		b.SetParallelism(16)
		b.RunParallel(func(pb *testing.PB) {
			for pb.Next() {
				_ = MakeSkipList(values...)
			}
		})
	})
}

func BenchmarkSkipList_Lookup(b *testing.B) {
	sizes := []int64{100, 10_000, 1_000_000}
	parallel := []int{1, 16, 32}
	for _, p := range parallel {
		for _, size := range sizes {
			b.Run(fmt.Sprintf("%d_parallel_%d_elements", p, size), func(b *testing.B) {
				b.ReportAllocs()
				b.SetParallelism(16)

				list := MakeSkipList(generateInts(0, int(size), 1)...)

				b.ResetTimer()
				b.RunParallel(func(pb *testing.PB) {
					i := rand.Int64() % size
					for pb.Next() {
						if !list.Lookup(i) {
							panic(fmt.Sprintf("%d must be found", i))
						}
						i += 7
						if i >= size {
							i -= size
						}
					}
				})
			})
		}
	}
}

func BenchmarkSkipList_Insert(b *testing.B) {
	const mod = int64(1<<31 - 1)

	list := NewSkipList[int64]()
	b.ReportAllocs()
	b.SetParallelism(16)
	b.RunParallel(func(pb *testing.PB) {
		step := rand.Int64() % 61
		i := int64(0)
		for pb.Next() {
			list.Insert(i)
			i = (i + step) % mod
		}
	})
}

func BenchmarkSkipList_Delete(b *testing.B) {
	list := MakeSkipList(generateInts(-1_000_000, 1_000_000, 1)...)
	b.ReportAllocs()
	b.SetParallelism(16)
	b.RunParallel(func(pb *testing.PB) {
		i := int64(0)
		for pb.Next() {
			if list.IsEmpty() {
				break
			}
			list.Delete(i)
			i += 7
			if i > 1_000_000 {
				i -= 2_000_000
			}
		}
	})
}

func BenchmarkSkipList_Stress(b *testing.B) {
	list := MakeSkipList(generateInts(-1_000, 1_000, 4)...)

	b.ReportAllocs()
	b.SetParallelism(32)
	b.RunParallel(func(pb *testing.PB) {
		seed := uint64(uintptr(unsafe.Pointer(pb)))
		mask := uint64(1<<32 - 1)
		rnd := rand.New(rand.NewPCG(seed&mask, seed^mask))
		for pb.Next() {
			x := rnd.Int64()
			switch rnd.Int() % 3 {
			case 0:
				list.Lookup(x)
			case 1:
				list.Insert(x)
			case 2:
				list.Delete(x)
			}
		}
	})
}
