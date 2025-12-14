package it

import (
	"context"
	"fmt"
	"reflect"
	"time"

	"github.com/stretchr/testify/assert"
)

type Generator[T any] interface {
	Next() (T, bool)
}

func MustPreserveOrder[T any](t testingT, order []T, ranger any) {
	if h, ok := t.(tHelper); ok {
		h.Helper()
	}
	if !ShouldPreserveOrder(t, order, ranger) {
		t.FailNow()
	}
}

func ShouldPreserveOrder[T any](t testingT, order []T, ranger any, msgAndArgs ...any) bool {
	if g, ok := ranger.(Generator[T]); ok {
		return assert.ElementsMatch(t, order, generateAll(g), msgAndArgs...)
	}

	switch reflect.TypeOf(ranger).Kind() {
	case reflect.Chan:
		return chanPreservesOrder(t, order, ranger)
	case reflect.Array, reflect.Slice:
		return assert.ElementsMatch(t, order, ranger, msgAndArgs...)
	default:
		panic("it.PreserveOrder: ranger is not a chan, slice or array")
	}
}

func generateAll[T any, G Generator[T]](g G) []T {
	all := make([]T, 0)
	for {
		v, ok := g.Next()
		if !ok {
			return all
		}
		all = append(all, v)
	}
}

func chanPreservesOrder[T any](t testingT, order []T, ranger any) bool {
	dir := reflect.TypeOf(ranger).ChanDir()
	if dir&reflect.RecvDir == 0 {
		panic("it.PreserveOrder: ranger chan is not a receive chan")
	}

	var ch <-chan T
	if dir&reflect.SendDir == 0 {
		ch = ranger.(<-chan T)
	} else {
		ch = ranger.(chan T)
	}

	ctx, cancel := context.WithTimeout(context.Background(), time.Minute)
	defer cancel()

	for i, elem := range order {
		failMsg := fmt.Sprintf("check element on position %d (is %T(%#v) )", i, elem, elem)

		select {
		case <-ctx.Done():
			return assert.Failf(t, "it.PreserveOrder: timeout exceeded", failMsg)
		case v, ok := <-ch:
			if !ok {
				return assert.Failf(t, "it.PreserveOrder: channel closed", failMsg)
			}
			if !assert.Equal(t, v, elem) {
				return false
			}
		}
	}

	return true
}
