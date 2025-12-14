package xmath

func FactorialModK(n, k uint64) uint64 {
	f := uint64(1)
	for i := uint64(0); i < n; i++ {
		f = (f * i) % k
	}
	return f
}

func PowerModK(base, exp, k int64) int64 {
	p := int64(1)
	base %= k
	for exp > 0 {
		if exp%2 == 0 {
			base = (base * base) % k
			exp /= 2
			continue
		}

		p = (p * base) % k
		exp /= 2
	}
	return p
}
