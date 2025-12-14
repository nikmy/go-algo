package lockfree

import (
	"sync/atomic"

	"github.com/nikmy/algo/syncx/atomx"
)

type queueNode[T any] struct {
	elem T
	next atomic.Pointer[queueNode[T]]
}

func NewQueue[T any]() *Queue[T] {
	q := &Queue[T]{sent: new(queueNode[T])}
	q.sent.next.Store(q.sent)
	q.head.Store(q.sent)
	q.tail.Store(q.sent)

	return q
}

// Queue is lock-free queue implementation
type Queue[T any] struct {
	sent *queueNode[T]
	head atomx.Pointer[queueNode[T]]
	tail atomx.Pointer[queueNode[T]]
}

func (q *Queue[T]) TryProduce(x T) bool {
	q.PushBack(x)
	return true
}

func (q *Queue[T]) TryConsume(x *T) bool {
	if pop := q.PopFront(); pop != nil {
		*x = *pop
		return true
	}

	return false
}

func (q *Queue[T]) PushBack(x T) {
	n := queueNode[T]{elem: x}
	n.next.Store(q.sent)

	for {
		tail := q.tail.Load()

		next := tail.next.Load()
		if next != q.sent {
			// helping
			q.tail.CompareAndSwap(tail, next)
			continue
		}

		if !tail.next.CompareAndSwap(next, &n) {
			continue
		}

		q.tail.Store(&n)
		break
	}
}

func (q *Queue[T]) PopFront() *T {
	for {
		head := q.head.Load()
		if head == q.sent {
			return nil
		}

		if q.head.CompareAndSwap(head, head.next.Load()) {
			return &head.elem
		}
	}
}
