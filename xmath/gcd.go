package xmath

import "golang.org/x/exp/constraints"

// ExtendedGCD returns x, y, g such that a*x + b*y = g, g = gcd(a, b)
func ExtendedGCD[T constraints.Integer](a, b T) (x T, y T, g T) {
	if b == 0 {
		return a, 1, 0
	}
	g, x, y = ExtendedGCD(b, a%b)
	x, y = y, x-(a/b)*y
	return
}

// GCD return the greatest common divisor of a and b
func GCD[T constraints.Integer](a, b T) T {
	return binaryGcd(a, b)
}

// LCM returns the least common multiple of a and b
func LCM[T constraints.Integer](a, b T) T {
	return a / binaryGcd(a, b) * b
}

func binaryGcd[T constraints.Integer](a, b T) T {
	switch {
	case a == b || b == 0:
		return a
	case a == 0:
		return b
	case a%2 == 0:
		if b%2 == 0 {
			return binaryGcd(a/2, b/2) * 2
		}
		return binaryGcd(a/2, b)
	case b%2 == 0:
		return binaryGcd(a, b/2)
	}

	if a > b {
		return binaryGcd((a-b)/2, b)
	}
	return binaryGcd((b-a)/2, a)
}
