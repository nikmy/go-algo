package reflex

import "reflect"

func TryAssign(lhs, rhs reflect.Value) bool {
	dstType := lhs.Type()
	if rhs.Type().ConvertibleTo(dstType) && dstType.AssignableTo(dstType) {
		lhs.Set(rhs.Convert(dstType))
		return true
	}
	return false
}
