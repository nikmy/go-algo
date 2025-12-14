package it

import (
	"reflect"

	"github.com/stretchr/testify/assert"
)

func ShouldNotEqual(t testingT, want, got interface{}, msgAndArgs ...any) bool {
	wv, gv := reflect.ValueOf(want), reflect.ValueOf(got)
	if isFloat(wv.Type()) && isFloat(gv.Type()) {
		return equalFloats(t, true, wv.Float(), wv.Float(), 1e-10, msgAndArgs...)
	}

	return assert.Equal(t, want, got, msgAndArgs...)
}

func MustEqual(t testingT, x, y any, msgAndArgs ...any) {
	if h, ok := t.(tHelper); ok {
		h.Helper()
	}
	if !ShouldEqual(t, x, y, msgAndArgs...) {
		t.FailNow()
	}
}

func ShouldEqual(t testingT, want, got interface{}, msgAndArgs ...any) bool {
	wv, gv := reflect.ValueOf(want), reflect.ValueOf(got)
	if isFloat(wv.Type()) && isFloat(gv.Type()) {
		return equalFloats(t, false, wv.Float(), wv.Float(), 1e-10, msgAndArgs...)
	}

	return assert.Equal(t, want, got, msgAndArgs...)
}

func isFloat(t reflect.Type) bool {
	return t.Kind() == reflect.Float32 || t.Kind() == reflect.Float64
}

func equalFloats(t testingT, flip bool, want, got float64, accuracy float64, msgAndArgs ...any) bool {
	if want == got {
		return !flip
	}

	if want == 0 {
		if flip {
			return assert.NotZero(t, got, msgAndArgs...)
		}
		return assert.Zero(t, got, msgAndArgs...)
	}

	// TODO
	return assert.InEpsilon(t, want, got, accuracy, msgAndArgs)
}
