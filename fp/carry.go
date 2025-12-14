package fp

func Carry2[Arg1, Arg2, Result any](f func (Arg1, Arg2) Result, arg1 Arg1) func(Arg2) Result {
	return func(arg2 Arg2) Result {
		return f(arg1, arg2)
	}
}
