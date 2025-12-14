package hashmap

import (
	"math/rand"
	"testing"

	"github.com/stretchr/testify/require"
)

func TestCuckoo(t *testing.T) {
	m := NewCuckoo[int, int]()

	rnd := rand.New(rand.NewSource(42))

	const n = 8

	et := make(map[int]int, n)
	for i := 0; i < n; i++ {
		et[rnd.Int()] = rnd.Int()
	}

	for k, v := range et {
		m.Insert(k, v)
	}

	for k, v := range et {
		got, ok := m.Lookup(k)
		require.Equal(t, true, ok)
		require.Equal(t, v, got)
	}

	require.Less(t, m.Len()/n, 8, "overhead over 100x")

}
