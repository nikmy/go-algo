package static

import (
	"testing"

	"github.com/stretchr/testify/require"
)

func TestMap(t *testing.T) {
	t.Run("int keys", func(t *testing.T) {
		type kv = KV[int, int]
		squares := make([]kv, 0, 200)
		for i := range cap(squares) {
			squares = append(squares, kv{i, i * i})
		}
		m := NewMap(squares...)
		for k, v := range m.Entries() {
			require.Equal(t, k*k, v)
		}
	})
}
