package bitset

import (
	"fmt"
	"math/bits"
)

func New(size int) *Bitset {
	lastBucketSize, size := size&63, size>>6
	if lastBucketSize > 0 {
		size++
	}
	data := make([]uint64, size)
	return &Bitset{data, lastBucketSize}
}

func Equals(left, right *Bitset) bool {
	if left.Size() != right.Size() {
		return false
	}
	return !Xor(left, right).Any()
}

func Xor(bs, mask *Bitset) *Bitset {
	if bs.Size() < mask.Size() {
		bs, mask = mask, bs
	}
	xor := New(bs.Size())
	for b := 0; b < mask.Size(); b++ {
		xor.data[b] = bs.data[b] ^ mask.data[b]
	}
	copy(xor.data[mask.Size():], bs.data[mask.Size():])
	return xor
}

type Bitset struct {
	data           []uint64
	lastBucketSize int
}

func (bs *Bitset) Size() int {
	b := len(bs.data)
	if b == 0 {
		return 0
	}
	return (b-1)<<6 + bs.lastBucketSize
}

func (bs *Bitset) FlipBit(idx int) {
	bucket, ind := idx>>6, uint64(idx&63)
	mask := uint64(1) << ind
	bs.data[bucket] ^= mask
}

func (bs *Bitset) Fix(idx int) {
	bucket, ind := idx>>6, uint64(idx&63)
	mask := uint64(1) << ind
	bs.data[bucket] |= mask
}

func (bs *Bitset) Unfix(idx int) {
	bucket, ind := idx>>6, uint64(idx&63)
	mask := ^(uint64(1) << ind)
	bs.data[bucket] &= mask
}

func (bs *Bitset) Get(idx int) bool {
	bucket, ind := idx>>6, uint64(idx&63)
	mask := uint64(1) << ind
	return bs.data[bucket]&mask != 0
}

func (bs *Bitset) Flip() {
	mask := ^uint64(0)
	for i := range bs.data {
		if i == len(bs.data)-1 && bs.lastBucketSize > 0 {
			mask = (uint64(1) << bs.lastBucketSize) - 1
		}
		bs.data[i] ^= mask
	}
}

func (bs *Bitset) All() bool {
	allOnes := ^uint64(0)
	for i, b := range bs.data {
		if b != allOnes {
			if bs.lastBucketSize == 0 || i != len(bs.data)-1 {
				return false
			}
			if b != (uint64(1)<<bs.lastBucketSize)-1 {
				return false
			}
		}
	}
	return true
}

func (bs *Bitset) Any() bool {
	allOnes := ^uint64(0)
	for _, b := range bs.data {
		if (b & allOnes) != 0 {
			return true
		}
	}
	return false
}

func (bs *Bitset) Count() int {
	cnt := 0
	for _, b := range bs.data {
		cnt += bits.OnesCount64(b)
	}
	return cnt
}

func (bs *Bitset) ToString() string {
	str := ""

	special := 0
	if bs.lastBucketSize > 0 {
		special = 1
	}

	for i := 0; i < len(bs.data)-special; i++ {
		str += fmt.Sprintf("%064b", bits.Reverse64(bs.data[i]))

	}

	if special != 0 {
		b := bs.data[len(bs.data)-1]
		str += fmt.Sprintf("%064b", bits.Reverse64(b))[:bs.lastBucketSize]
	}

	return str
}
