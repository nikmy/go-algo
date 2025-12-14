package iterx

import (
	"iter"
)

func Map[From, To any](mapper func(From) To, s iter.Seq[From]) iter.Seq[To] {
	return func(yield func(To) bool) {
		for v := range s {
			if !yield(mapper(v)) {
				return
			}
		}
	}
}

func Filter[T any](filter func(T) bool, s iter.Seq[T]) iter.Seq[T] {
	return func(yield func(T) bool) {
		for v := range s {
			if !filter(v) {
				continue
			}
			if !yield(v) {
				return
			}
		}
	}
}

func Reduce[T any](reducer func(accum, val T) T, init T, s iter.Seq[T]) T {
	for v := range s {
		init = reducer(init, v)
	}
	return init
}

func Batches[T any](maxSize int, s iter.Seq[T]) iter.Seq[[]T] {
	return func(yield func([]T) bool) {
		batch := make([]T, 0, maxSize)
		for v := range s {
			batch = append(batch, v)
			if len(batch) == cap(batch) {
				if !yield(batch) {
					break
				}
			}
		}
		if len(batch) > 0 {
			yield(batch)
		}
	}
}

func Flatten[T any, S ~[]T](batches iter.Seq[S]) iter.Seq[T] {
	return func(yield func(T) bool) {
		for b := range batches {
			for _, v := range b {
				if !yield(v) {
					return
				}
			}
		}
	}
}
