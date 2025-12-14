package deque

import "unsafe"

type Deque[T any] struct {
	bufHead *buffer[T]
	bufTail *buffer[T]
}

func (d *Deque[T]) Clear() {
	d.bufTail = d.bufHead
	d.bufHead.head = 0
	d.bufHead.tail = 0
}

func (d *Deque[T]) Len() int {
	size := 0
	for i := d.bufHead; i != nil; i = i.next {
		size += i.size
	}

	return size
}

func (d *Deque[T]) IsEmpty() bool {
	return !d.bufHead.isEmpty()
}

func (d *Deque[T]) PushFront(x T) {
	if d.bufHead == nil {
		d.bufHead = newBuffer[T]()
		d.bufTail = d.bufHead
	}

	if d.bufHead.isFull() {
		old := d.bufHead
		d.bufHead = newBuffer[T]()
		d.bufHead.next = old
		old.prev = d.bufHead
	}

	d.bufHead.pushFront(x)
}

func (d *Deque[T]) PushBack(x T) {
	if d.bufHead == nil {
		d.bufHead = newBuffer[T]()
		d.bufTail = d.bufHead
	}

	if d.bufTail.isFull() {
		old := d.bufTail
		d.bufTail = newBuffer[T]()
		old.next = d.bufTail
		d.bufTail.prev = old
	}

	d.bufTail.pushBack(x)
}

func (d *Deque[T]) PopFront() T {
	if d.bufHead.isEmpty() {
		panic("pop of empty deque")
	}

	pop := d.bufHead.popFront()
	if d.bufHead != d.bufTail && d.bufHead.isEmpty() {
		d.bufHead = d.bufHead.next
		d.bufHead.prev = nil

	}

	return pop
}

func (d *Deque[T]) PopBack() T {
	if d.bufTail.isEmpty() {
		panic("pop of empty deque")
	}

	pop := d.bufTail.popBack()
	if d.bufTail != d.bufHead && d.bufTail.isEmpty() {
		d.bufTail = d.bufHead.prev
		d.bufTail.next = nil
	}

	return pop
}

const bufferBytesLen = 4096

func newBuffer[T any]() *buffer[T] {
	tSize, tAlign := unsafe.Sizeof(*new(T)), unsafe.Alignof(*new(T))
	if r := tSize % tAlign; r > 0 {
		tSize += tAlign - r
	}

	targetLen := bufferBytesLen / tSize
	return &buffer[T]{constBuffer: constBuffer[T]{data: make([]T, targetLen)}}
}

type buffer[T any] struct {
	constBuffer[T]
	next *buffer[T]
	prev *buffer[T]
}

func (b *buffer[T]) pushBack(x T) {
	b.data[b.tail] = x
	b.tail = b.idx(b.tail + 1)
	b.size++
}

func (b *buffer[T]) pushFront(x T) {
	b.head = b.idx(b.head - 1)
	b.data[b.head] = x
	b.size++
}

func (b *buffer[T]) popFront() T {
	pop := b.front()
	b.head = b.idx(b.head + 1)
	b.size--
	return pop
}

func (b *buffer[T]) popBack() T {
	pop := b.back()
	b.tail = b.idx(b.tail - 1)
	b.size--
	return pop
}

type constBuffer[T any] struct {
	data []T
	head int
	tail int
	size int
}

func (b constBuffer[T]) front() T {
	return b.data[b.head]
}

func (b constBuffer[T]) back() T {
	return b.data[b.idx(b.tail-1)]
}

func (b constBuffer[T]) isFull() bool {
	return b.size == len(b.data)
}

func (b constBuffer[T]) isEmpty() bool {
	return b.size == 0
}

func (b constBuffer[T]) idx(i int) int {
	return i % len(b.data)
}
