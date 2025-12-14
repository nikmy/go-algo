package heap

type Heap[T comparable] interface {
	PriorityQueue[T]
	DecreaseKey(old, new T) (found bool)
}

var _ Heap[struct{}] = new(heap[struct{}])

func Make[T comparable](compareFunc func(T, T) int, elems []T) *heap[T] {
	h := New(compareFunc)
	h.q.heapify()
	return h
}

func New[T comparable](compareFunc func(T, T) int) *heap[T] {
	h := &heap[T]{
		q: new(pq[T]),
		m: make(map[T]int),
	}

	swap := func(i int, j int) {
		h.m[h.q.data[i]] = j
		h.m[h.q.data[j]] = i
		h.q.swapElems(i, j)
	}

	h.q.init(compareFunc, nil, swap)

	return h
}

type heap[T comparable] struct {
	q *pq[T]
	m map[T]int
}

func (h *heap[T]) Min() T {
	return h.q.Min()
}

func (h *heap[T]) Size() int {
	return h.q.Size()
}

func (h *heap[T]) Empty() bool {
	return h.q.Empty()
}

func (h *heap[T]) Push(x T) {
	h.q.Push(x)
}

func (h *heap[T]) Pop() T {
	v := h.q.Pop()
	delete(h.m, v)
	return v
}

func (h *heap[T]) DecreaseKey(old, new T) bool {
	idx, found := h.m[old]
	if !found {
		return false
	}

	h.q.data[idx] = new

	c := h.q.comp(old, new)
	switch {
	case c < 0:
		h.q.siftUp(idx)
	case c > 0:
		h.q.siftDown(idx)
	}

	return true
}
