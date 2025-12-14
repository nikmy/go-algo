package fp

func Not[I any, P ~func(I) bool](condition P) P {
	return func(i I) bool {
		return !condition(i)
	}
}

func Or[I any, P ~func(I) bool](p, q P) P {
	return func(i I) bool {
		return p(i) || q(i)
	}
}

func And[I any, P ~func(I) bool](p, q P) P {
	return func(i I) bool {
		return p(i) && q(i)
	}
}
