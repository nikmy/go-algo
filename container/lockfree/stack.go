package lockfree

import (
	"github.com/nikmy/algo/syncx/atomx"
)

type stackNode[T any] struct {
	elem T
	next *stackNode[T]
}

type Stack[T any] struct {
	top atomx.Pointer[stackNode[T]]
}

func (s *Stack[T]) Empty() bool {
	return s.top.Load() == nil
}

func (s *Stack[T]) TryPush(x T) bool {
	n := stackNode[T]{
		elem: x,
		next: s.top.Load(),
	}

	return s.top.CompareAndSwap(n.next, &n)
}

func (s *Stack[T]) TryPop(x *T) bool {
	top := s.top.Load()
	if top == nil {
		return false
	}

	if !s.top.CompareAndSwap(top, top.next) {
		return false
	}

	*x = top.elem
	return true
}
