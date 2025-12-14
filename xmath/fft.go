package xmath

import (
	"math"
	"math/cmplx"
)

type Point struct {
	X complex128
	Y float64
}

// FFT performs Fast Fourier Transform.
func FFT(coefs []complex128) (values []complex128) {
	return fftRecursive(coefs)
}

const twopi complex128 = 2 * math.Pi * 1i

func fftRecursive(p []complex128) []complex128 {
	n := len(p)
	if n < 3 {
		return p
	}

	omega := cmplx.Exp(twopi / complex(float64(n), 0))

	if len(p)%2 == 1 {
		// x^0, x^1, ..., x^2k

		pe, po := make([]complex128, n/2+1), make([]complex128, n/2)
		for i := 0; i < n-1; i += 2 {
			pe[i/2] = p[i]
			po[i/2] = p[i+1]
		}
		pe[n/2] = p[n-1]

		fe, fo := fftRecursive(pe), fftRecursive(po)
		f := make([]complex128, len(p))
		for i := range len(f) / 2 {
			wi := cmplx.Pow(omega, complex(float64(i), 0))
			f[i] = fe[i] + wi*fo[i]
			f[i+n/2] = fe[i] - wi*fo[i]
		}

		return f
	}

	panic("not implemented")
}
