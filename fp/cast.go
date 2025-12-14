package fp

import "unsafe"

type NumLike interface {
	~int | ~int8 | ~int16 | ~int32 | ~int64 |
		~uint | ~uint8 | ~uint16 | ~uint32 | ~uint64 |
		~float32 | ~float64
}

func StaticCastNumber[From, To NumLike](x From) To {
	return (To)(x)
}

type StringLike interface {
	~[]byte | ~string
}

func StaticCastString[From, To StringLike](x From) To {
	return (To)(x)
}

func DynamicCast[From, To any](x From) *To {
	to, ok := any(x).(To)
	if !ok {
		return nil
	}
	return &to
}

func DynamicCastResult[To, Arg, From any, F ~func(Arg) From](f F) func(Arg) To {
	return func(a Arg) To { return any(f(a)).(To) }
}

func DynamicCastFunc[
	FromArg, ToArg any,
	FromResult, ToResult any,
	From ~func(FromArg) FromResult,
](f From) func(ToArg) ToResult {
	return func(arg ToArg) ToResult {
		return any(f(any(arg).(FromArg))).(ToResult)
	}
}

func UnsafeCast[From, To any](arg *From) *To {
	return (*To)(unsafe.Pointer(arg))
}

func UnsafeAnyCast[To any](v any) *To {
	type iface struct {
		typ unsafe.Pointer
		val unsafe.Pointer
	}
	return (*To)((*iface)(unsafe.Pointer(&v)).val)
}
