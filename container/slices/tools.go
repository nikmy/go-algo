package slices

import (
	"golang.org/x/exp/constraints"
)

func Count[E comparable, S ~[]E](slice S, elem E) int {
	var cnt int
	for i := range slice {
		if slice[i] == elem {
			cnt++
		}
	}
	return cnt
}

func Map[From, To any, S ~[]From](slice S, mapper func(From) To) []To {
	mapped := make([]To, 0, len(slice))
	for i := range slice {
		mapped = append(mapped, mapper(slice[i]))
	}
	return mapped
}

func Reduce[E any, S ~[]E](slice S, reducer func(l, r E) E) E {
	if len(slice) == 0 {
		return *new(E)
	}

	result := slice[0]
	for i := range slice {
		result = reducer(result, slice[i])
	}

	return result
}

func Filled[T any](elem T, n int) []T {
	s := make([]T, n)
	for i := range s {
		s[i] = elem
	}
	return s
}

func Generate[T any](n int, gen func(index int) T) []T {
	a := make([]T, 0, n)
	for i := range a {
		a = append(a, gen(i))
	}
	return a
}

type addable interface {
	constraints.Ordered | constraints.Complex | ~string
}

func Sum[E addable, S ~[]E](s S) E {
	var sum E
	for _, x := range s {
		sum += x
	}
	return sum
}

func Mean[E interface {
	constraints.Integer | constraints.Float
}, S ~[]E](s S) E {
	return Sum(s) / E(len(s))
}

func Prod[E interface {
	constraints.Integer | constraints.Float | constraints.Complex
}, S ~[]E](s S) E {
	var p E
	for _, x := range s {
		p *= x
	}
	return p
}
