package lockfree

import (
	"math/bits"

	"github.com/nikmy/algo/syncx/atomx"
)

func NewPipe[T any](size uint32) *Pipe[T] {
	mask := bits.Len32(size) - 1
	if bits.OnesCount32(size) > 1 {
		mask++
	}

	size = max(uint32(2), 1<<uint32(mask))
	return &Pipe[T]{
		data: make([]T, size),
		mask: size - 1,
	}
}

// Pipe is optimized channel
// for single producer - single consumer
// case. It is about four times faster than
// builtin channel for parallel access, and
// double time faster in single thread case.
// TryProduce and TryConsume are wait free
// for both sides.
type Pipe[T any] struct {
	data []T

	head uint32
	_    [64]byte

	tail uint32
	_    [64]byte

	headCache uint32
	tailCache uint32

	mask uint32
}

func (p *Pipe[T]) TryConsume(x *T) bool {
	currHead := atomx.LoadUint32(&p.head)
	currTail := p.tailCache

	if currHead == currTail {
		currTail = atomx.LoadUint32(&p.tail)
		p.tailCache = currTail
	}

	if p.isEmpty(currHead, currTail) {
		return false
	}

	*x = p.data[currHead&p.mask]
	atomx.StoreUint32(&p.head, p.next(currHead))

	return true
}

func (p *Pipe[T]) TryProduce(x T) bool {
	currTail := atomx.LoadUint32(&p.tail)
	currHead := p.headCache

	if p.next(currTail) == currHead {
		currHead = atomx.LoadUint32(&p.head)
		p.headCache = currHead
	}

	if p.isFull(currHead, currTail) {
		return false
	}

	p.data[currTail&p.mask] = x
	atomx.StoreUint32(&p.tail, p.next(currTail))

	return true
}

func (p *Pipe[T]) isEmpty(head, tail uint32) bool {
	return tail == head
}

func (p *Pipe[T]) isFull(head, tail uint32) bool {
	return p.next(tail) == head
}

func (p *Pipe[T]) next(idx uint32) uint32 {
	return (idx + 1) & p.mask
}
