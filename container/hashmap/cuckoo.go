// UNSTABLE
package hashmap

import (
	"math"
	"math/rand/v2"

	"github.com/dolthub/maphash"
)

type Hash[T any] func(T) uintptr // TODO: universal hash for all types

const (
	notFound uint64 = 0xFFFFFFFFFFFFFFF
	baseXor  uint64 = 0x101010101
)

func NewCuckoo[K comparable, V any]() *Cuckoo[K, V] {
	return &Cuckoo[K, V]{
		mainHash: maphash.NewHasher[K](),
		xorValue: rand.N[uint64](math.MaxUint64) & baseXor,
		logB:     2,
		buckets:  make([]cuckooBucket[K, V], 4),
	}
}

// Cuckoo is hashmap optimized for reads. It has fair O(1)
// asymptotic for lookups and removes, and amortized O*(1)
// asymptotic for inserts. It is guaranteed that memory is
// used for storing elements is O(4N).
type Cuckoo[K comparable, V any] struct {
	mainHash maphash.Hasher[K]
	xorValue uint64
	logB     uint8
	buckets  []cuckooBucket[K, V]
	busy     int
}

func (c *Cuckoo[K, V]) Len() int {
	return 1 << c.logB
}

func (c *Cuckoo[K, V]) Lookup(key K) (V, bool) {
	kv, _ := c.lookup(key)
	return kv.v, kv.ok
}

func (c *Cuckoo[K, V]) Remove(key K) bool {
	_, b := c.lookup(key)
	if b == notFound {
		return false
	}

	c.buckets[b] = cuckooBucket[K, V]{}
	c.busy--
	return true
}

func (c *Cuckoo[K, V]) Insert(key K, value V) {
	for !c.tryInsert(key, value) {
		c.evacuate()
	}
	c.busy++
}

func (c *Cuckoo[K, V]) lookup(key K) (cuckooBucket[K, V], uint64) {
	b := c.hash(key)

	first := c.buckets[b]
	if first.ok && first.k == key {
		return first, b
	}

	b = c.flipHash(b)
	second := c.buckets[b]
	if second.ok && second.k == key {
		return second, b
	}

	return cuckooBucket[K, V]{}, notFound
}

func (c *Cuckoo[K, V]) tryInsert(key K, value V) bool {
	b := c.hash(key)

	first := c.buckets[b]
	if !first.ok {
		c.buckets[b] = cuckooBucket[K, V]{true, key, value}
		return true
	}
	if first.k == key {
		return true
	}

	b = c.flipHash(b)
	second := c.buckets[b]
	if !second.ok {
		c.buckets[b] = cuckooBucket[K, V]{true, key, value}
		return true
	}
	if second.k == key {
		return true
	}

	if c.free(b) {
		c.buckets[b] = cuckooBucket[K, V]{true, key, value}
		return true
	}

	return false
}

func (c *Cuckoo[K, V]) flipHash(h uint64) uint64 {
	return h ^ (c.xorValue & uint64(c.Len()-1))
}

func (c *Cuckoo[K, V]) hash(key K) uint64 {

	return c.mainHash.Hash(key) & uint64(c.Len()-1)
}

func (c *Cuckoo[K, V]) free(b uint64) bool {
	swapPath := make([]uint64, 0, 64)
	for len(swapPath) < c.Len() && c.buckets[b].ok {
		swapPath = append(swapPath, b)
		b = c.flipHash(b)
	}

	if len(swapPath) >= c.Len() {
		return false
	}

	for len(swapPath) > 0 {
		s := len(swapPath) - 1
		c.buckets[b], c.buckets[s] = c.buckets[s], c.buckets[b]
		b = swapPath[s]
	}

	return true
}

func (c *Cuckoo[K, V]) overflow() bool {
	return c.busy > c.Len()/4
}

func (c *Cuckoo[K, V]) evacuate() {
	oldBuckets := c.buckets
	if c.overflow() {
		c.logB++
	}
	c.buckets = make([]cuckooBucket[K, V], c.Len())

tryHash:
	for _, kv := range oldBuckets {
		if !kv.ok {
			continue
		}
		if !c.tryInsert(kv.k, kv.v) {
			c.mainHash = maphash.NewSeed[K](c.mainHash)
			goto tryHash
		}
	}
}

type cuckooBucket[K comparable, V any] struct {
	ok bool
	k  K
	v  V
}
