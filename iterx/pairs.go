package iterx

import (
	"iter"
)

func Left[L, R any](seq iter.Seq2[L, R]) iter.Seq[L] {
	return func(yield func(L) bool) {
		for l := range seq {
			if !yield(l) {
				return
			}
		}
	}
}

func Right[L, R any](seq iter.Seq2[L, R]) iter.Seq[R] {
	return func(yield func(R) bool) {
		for _, r := range seq {
			if !yield(r) {
				return
			}
		}
	}
}

func Zip[L, R any](left iter.Seq[L], right iter.Seq[R]) iter.Seq2[L, R] {
	return Pairs(left, right)
}

func Pairs[L, R any](left iter.Seq[L], right iter.Seq[R]) iter.Seq2[L, R] {
	return func(yield func(L, R) bool) {
		nextRight, stop := iter.Pull(right)
		defer stop()

		for l := range left {
			r, ok := nextRight()
			if !ok {
				return
			}
			yield(l, r)
		}
	}
}
