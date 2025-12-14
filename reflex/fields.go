package reflex

import (
	"iter"
	"reflect"
)

func ValueFields(v reflect.Value) iter.Seq[reflect.Value] {
	if v.Kind() != reflect.Struct {
		panic("obtaining fields of non-struct type")
	}

	return func(yield func(reflect.Value) bool) {
		for i := 0; i < v.NumField(); i++ {
			if !yield(v.Field(i)) {
				return
			}
		}
	}
}

func TypeFields(t reflect.Type) iter.Seq[reflect.StructField] {
	if t.Kind() != reflect.Struct {
		panic("obtaining fields of non-struct type")
	}

	return func(yield func(reflect.StructField) bool) {
		for i := 0; i < t.NumField(); i++ {
			if !yield(t.Field(i)) {
				return
			}
		}
	}
}

func EqualFields(a, b []reflect.StructField) bool {
	if len(a) != len(b) {
		return false
	}

	for i := range a {
		if a[i].Type != b[i].Type || a[i].Name != b[i].Name {
			return false
		}
	}

	return true
}
