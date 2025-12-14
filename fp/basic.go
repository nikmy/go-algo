package fp

func Id[T any](x T) T { return x }

func Empty[T any]() T {
	var empty T
	return empty
}