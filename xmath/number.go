package xmath

import (
	"math/cmplx"

	"golang.org/x/exp/constraints"
)

type Real interface {
	constraints.Integer | constraints.Float
}

type Complex interface {
	Real | constraints.Complex
}

func Conjugate[T constraints.Complex](x T) T {
	return T(cmplx.Conj(complex128(x)))
}

func Distance[T Real](a, b T) T {
	if a > b {
		return a - b
	}
	return b - a
}

func Abs[T Real](x T) T {
	if x < 0 {
		return -x
	}
	return x
}
