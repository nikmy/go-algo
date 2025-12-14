package heap

type PriorityQueue[T any] interface {
	Min() T
	Size() int
	Empty() bool

	Push(T)
	Pop() T
}

func NewPriorityQueue[T any](compareFunc func(T, T) int) *pq[T] {
	q := new(pq[T])
	q.init(compareFunc, nil, q.swapElems)
	return q
}

func Prioritize[T any](compareFunc func(T, T) int, elems []T) *pq[T] {
	q := new(pq[T])
	q.init(compareFunc, elems, q.swapElems)
	return q
}

var _ PriorityQueue[struct{}] = new(pq[struct{}])

type pq[T any] struct {
	comp func(T, T) int
	swap func(int, int)
	data []T
}

func (q *pq[T]) init(compareFunc func(T, T) int, elems []T, swap func(i, j int)) *pq[T] {
	q.comp = compareFunc
	q.data = elems
	q.swap = swap

	if len(elems) > 1 {
		q.heapify()
	}

	return q
}

func (q *pq[T]) heapify() {
	for i := len(q.data) / 2; i >= 0; i-- {
		q.siftDown(i)
	}
}

func (q *pq[T]) Size() int {
	return len(q.data)
}

func (q *pq[T]) Empty() bool {
	return q.Size() == 0
}

func (q *pq[T]) Min() T {
	return q.data[0]
}

func (q *pq[T]) Push(x T) {
	q.data = append(q.data, x)
	q.siftUp(len(q.data) - 1)
}

func (q *pq[T]) Pop() T {
	m := q.data[0]
	last := len(q.data) - 1
	q.swap(0, last)
	q.data = q.data[:last]
	q.siftDown(0)
	return m
}

func (q *pq[T]) siftUp(i int) {
	for q.less(i, (i-1)/2) {
		q.swap(i, (i-1)/2)
		i = (i - 1) / 2
	}
}

func (q *pq[T]) siftDown(i int) {
	for 2*i+1 < q.Size() {
		left := 2*i + 1
		right := 2*i + 2
		j := left
		if right < q.Size() && q.less(right, left) {
			j = right
		}
		if q.lessOrEqual(i, j) {
			break
		}
		q.swap(i, j)
		i = j
	}
}

func (q *pq[T]) swapElems(i, j int) {
	q.data[i], q.data[j] = q.data[j], q.data[i]
}

func (q *pq[T]) less(i, j int) bool {
	return q.comp(q.data[i], q.data[j]) < 0
}

func (q *pq[T]) lessOrEqual(i, j int) bool {
	return q.comp(q.data[i], q.data[j]) <= 0
}

