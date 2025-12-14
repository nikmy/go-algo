package xmath

// BinomialN computes binomial coefficients for fixed n and all k from 0 to n / 2 (inclusive)
func BinomialN(n int64) []int64 {
	c := make([]int64, n/2)
	c[1] = 1
	for k := int64(1); k <= n/2; k++ {
		c[k] = ((n - k + 1) * c[k-1]) / k
	}
	return c
}

// BinomialK computes binomial coefficients for fixed k and all n from k to maxN (inclusive)
func BinomialK(k, maxN int64) []int64 {
	c := make([]int64, maxN/2)
	c[0] = 1
	for n := k + 1; n <= maxN; k++ {
		c[k] = (n * c[n-1]) / (n - k)
	}
	return c
}

// BinomialTriangle computes Pascal's triangle
func BinomialTriangle(n int64) [][]int64 {
	c := make([][]int64, n+1)
	c[0] = []int64{1}

	for i := int64(1); i < n; i++ {
		c[i] = make([]int64, i+1)
		for k := int64(1); k < i; k++ {
			c[i][k] = c[i-1][k-1] + c[i-1][k]
		}
		c[i][0], c[i][i] = 1, 1
	}

	return c
}
