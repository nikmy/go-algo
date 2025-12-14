package static

import (
	"hash/maphash"
	"iter"
	"math/rand/v2"
)

// UNSTABLE
type Map[K comparable, V any] = *hmap[K, V]

type KV[K comparable, V any] struct {
	Key K
	Val V
}

// UNSTABLE
func NewMap[K comparable, V any](entries ...KV[K, V]) Map[K, V] {
	size := 4 * len(entries)
	seed := maphash.MakeSeed()
	for {
		a, b := rand.Uint64(), rand.Uint64N(uint64(size))
		hash := func(k K) int {
			p := maphash.Comparable(seed, k)
			h := int((a*p + b) % (1<<31 - 1))
			h = h % size
			if h < 0 {
				h += size
			}
			return h
		}

		if tryHash(hash, entries) {
			data := make([]KV[K, V], size)
			for _, kv := range entries {
				data[hash(kv.Key)] = kv
			}
			return &hmap[K, V]{
				hash: hash,
				data: data,
			}
		}
	}
}

func tryHash[K comparable, V any](hash func(K) int, entries []KV[K, V]) bool {
	taken := map[int]struct{}{}
	for _, kv := range entries {
		i := hash(kv.Key)
		if _, ok := taken[i]; ok {
			return false
		}
		taken[i] = struct{}{}
	}
	return true
}

type hmap[K comparable, V any] struct {
	hash func(K) int
	data []KV[K, V]
}

func (h *hmap[K, V]) HasKey(key K) bool {
	return h.data[h.hash(key)].Key == key
}

func (h *hmap[K, V]) Lookup(key K) (V, bool) {
	kv := h.data[h.hash(key)]
	if kv.Key == key {
		return kv.Val, true
	}
	return *new(V), false
}

func (h *hmap[K, V]) Get(key K) V {
	v, _ := h.Lookup(key)
	return v
}

func (h *hmap[K, V]) Entries() iter.Seq2[K, V] {
	return func(yield func(K, V) bool) {
		for _, entry := range h.data {
			if !yield(entry.Key, entry.Val) {
				return
			}
		}
	}
}
