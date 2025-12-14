package it

import (
	"errors"
	"fmt"
	"reflect"
	"regexp"
	"testing"

	"github.com/stretchr/testify/assert"
)

func _() {
	re := regexp.MustCompile("^[a-zA-Z0-9]+$")

	op := func() (int, error) {
		return 0, errors.New("test error")
	}

	t := new(testing.T)
	Check(t).
		ThatValue(17).IsEqualTo(17).
		ThatValue(3).IsNotEqualTo(4).
		ThatResult(op()).IsNotZero().
		ThatValue("a115").Matches(re).
		ThatValue([]int{1, 2, 3}).HasLen(3).
		ThatValue(make([]int, 0, 10)).HasCap(10).
		ThatValue(true).IsTrue().AndExitIfNot().
		ThatValue(false).HasTypeOf(true)
}

type testRunner[T any] interface {
	testingT
	Run(string, func(T)) bool
}

func NewT[TT testRunner[TT]](tt TT) *T[TT] {
	return &T[TT]{
		t: tt,
	}
}

type T[TT testRunner[TT]] struct {
	t TT
	v any
}

func (t *T[TT]) Run(name string, f func(t *T[TT])) {
	t.t.Run(name, func(t TT) {
		f(NewT[TT](t))
	})
}

func (t *T[TT]) T() TT {
	return t.t
}

func Check[TT testRunner[TT]](t TT) *chain[TT] {
	return &chain[TT]{
		t: t,
	}
}

type chain[TT testRunner[TT]] struct {
	t TT
	v any
	e error
	r bool
}

func (c *chain[TT]) ThatValue(v any) *chain[TT] {
	c.r = true
	c.v = v
	c.e = nil
	return c
}

func (c *chain[TT]) ThatResult(v any, err error) *chain[TT] {
	c.v, c.e = v, err
	return c
}

func (c *chain[TT]) AndExitIfNot() *chain[TT] {
	if h, ok := any(c.t).(tHelper); ok {
		h.Helper()
	}
	if !c.r {
		c.t.FailNow()
	}
	return c
}

func (c *chain[TT]) IsTrue() *chain[TT] {
	if !assert.NoError(c.t, c.e) {
		c.r = false
		return c
	}
	c.r = assert.True(c.t, c.v.(bool))
	return c
}

func (c *chain[TT]) IsFalse() *chain[TT] {
	if !assert.NoError(c.t, c.e) {
		c.r = false
		return c
	}
	c.r = assert.False(c.t, c.v.(bool))
	return c
}

func (c *chain[TT]) IsEqualTo(want any) *chain[TT] {
	if !assert.NoError(c.t, c.e) {
		c.r = false
		return c
	}
	c.r = ShouldEqual(c.t, want, c.v)
	return c
}

func (c *chain[TT]) IsNotEqualTo(want any) *chain[TT] {
	if !assert.NoError(c.t, c.e) {
		c.r = false
		return c
	}
	c.r = ShouldEqual(c.t, want, c.v)
	return c
}

func (c *chain[TT]) HasLen(want int) *chain[TT] {
	if !assert.NoError(c.t, c.e) {
		c.r = false
		return c
	}
	c.r = assert.Len(c.t, c.v, want)
	return c
}

func (c *chain[TT]) HasCap(want int) *chain[TT] {
	if !assert.NoError(c.t, c.e) {
		c.r = false
		return c
	}
	got, ok := getCap(c.v)
	if !ok {
		c.r = assert.Fail(c.t, fmt.Sprintf("\"%T\" could not be applied builtin cap()", c.v))
	}
	c.r = assert.Equal(c.t, want, got)
	return c
}

func getCap(v any) (c int, ok bool) {
	defer func() { ok = recover() == nil }()
	return reflect.ValueOf(v).Cap(), true
}

func (c *chain[TT]) IsZero() *chain[TT] {
	if !assert.NoError(c.t, c.e) {
		c.r = false
		return c
	}

	c.r = assert.Zero(c.t, c.v)
	return c
}

func (c *chain[TT]) IsNotZero() *chain[TT] {
	if !assert.NoError(c.t, c.e) {
		c.r = false
		return c
	}

	c.r = assert.NotZero(c.t, c.v)
	return c
}

func (c *chain[TT]) IsEmpty() *chain[TT] {
	if !assert.NoError(c.t, c.e) {
		c.r = false
		return c
	}
	c.r = assert.Empty(c.t, c.v)
	return c
}

func (c *chain[TT]) IsNotEmpty() *chain[TT] {
	if !assert.NoError(c.t, c.e) {
		c.r = false
		return c
	}
	c.r = assert.NotEmpty(c.t, c.v)
	return c
}

func (c *chain[TT]) IsNil() *chain[TT] {
	if !assert.NoError(c.t, c.e) {
		c.r = false
		return c
	}
	c.r = assert.Nil(c.t, c.v)
	return c
}

func (c *chain[TT]) IsNotNil() *chain[TT] {
	if !assert.NoError(c.t, c.e) {
		c.r = false
		return c
	}
	c.r = assert.NotNil(c.t, c.v)
	return c
}

func (c *chain[TT]) HasTypeOf(v any) *chain[TT] {
	if !assert.NoError(c.t, c.e) {
		c.r = false
		return c
	}
	c.r = assert.IsType(c.t, v, c.v)
	return c
}

func (c *chain[TT]) Matches(re *regexp.Regexp) *chain[TT] {
	if !assert.NoError(c.t, c.e) {
		c.r = false
		return c
	}

	stringType := reflect.TypeFor[string]()
	if !reflect.TypeOf(c.v).ConvertibleTo(stringType) {
		panic(fmt.Sprintf("cannot match %T", c.v))
	}

	s := reflect.ValueOf(c.v).Convert(stringType).Interface().(string)
	m := re.MatchString(s)
	if !m {
		c.r = assert.Failf(c.t, "regexp match failed", "%s does not match '%s'", s, re.String())
	}

	return c
}
