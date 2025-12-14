package buffer

import (
	"io"
	"slices"
)

type Buffer struct {
	data []byte
	pool *shardedPool
}

func (b *Buffer) ReadAll(r io.ReadCloser) error {
	defer r.Close()
	_, err := b.ReadFrom(r)
	return err
}

func (b *Buffer) ReadFrom(r io.Reader) (int64, error) {
	b.Resize(0)
	for {
		n, err := r.Read(b.rest())
		b.Resize(len(b.data) + n)

		if err != nil {
			if err == io.EOF {
				err = nil
			}
			return int64(len(b.data)), err
		}

		b.growIfNeeded()
	}
}

func (b *Buffer) Data() []byte {
	return b.data
}

func (b *Buffer) Resize(n int) {
	if n > cap(b.data) {
		b.data = slices.Grow(b.data, n-len(b.data))
	}
	b.data = b.data[:n]
}

func (b *Buffer) Free() {
	if b.pool != nil && cap(b.data) > 0 {
		b.pool.getPoolForSize(cap(b.data)).put(b.data)
	}
	b.data = nil
}

func (b *Buffer) rest() []byte {
	return b.data[len(b.data):cap(b.data)]
}

func (b *Buffer) growIfNeeded() {
	b.data = slices.Grow(b.data, 1)
}
