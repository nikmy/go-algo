package reflex

import (
	"reflect"
	"slices"
)

func Override(orig, overrider reflect.Value) {
	if orig.Kind() == reflect.Interface {
		Override(orig.Elem(), overrider)
	}

	if orig.Kind() != reflect.Pointer {
		panic("overriden object must be of pointer type")
	}

	overrideCopy(orig, overrider)
}

func overrideCopy(w, r reflect.Value) {
	assertCanOverride(w, r)

	if r.IsZero() {
		return
	}

	if w.IsZero() {
		w.Set(r)
		return
	}

	if w.Kind() == reflect.Struct {
		rFields := slices.Collect(ValueFields(r))
		i := 0
		for wf := range ValueFields(w) {
			overrideCopy(wf, rFields[i])
			i++
		}
		return
	}

	if w.Kind() == reflect.Pointer {
		// override value by pointer recursively
		overrideCopy(w.Elem(), r.Elem())
		return
	}

	w.Set(r)
}

func assertCanOverride(w, r reflect.Value) {
	if w.Kind() != r.Kind() {
		panic("kind mismatch")
	}

	if w.Type() == r.Type() {
		return
	}

	if w.Kind() == reflect.Struct && EqualFields(slices.Collect(TypeFields(w.Type())), slices.Collect(TypeFields(r.Type()))) {
		return
	}

	panic("type mismatch")
}
