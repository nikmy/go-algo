package streams

import (
	"errors"
	"io"
)

var (
	EOS = errors.New("end of stream")

	ErrNoReadProgress  = errors.New("multiple reads return zero")
	ErrNoWriteProgress = errors.New("multiple writes return zero")
)

type Reader[T any] interface {
	Read(p []T) (n int, err error)
}

type Writer[T any] interface {
	Write(p []T) (n int, err error)
}

type ReadWriter[T any] interface {
	Reader[T]
	Writer[T]
}

type Closer = io.Closer

type ReadCloser[T any] interface {
	Reader[T]
	Closer
}

type WriteCloser[T any] interface {
	Writer[T]
	Closer
}

type ReadWriteCloser[T any] interface {
	Reader[T]
	Writer[T]
	Closer
}

func NopCloser[T any](r Reader[T]) ReadCloser[T] {
	return nopReadCloser[T]{Reader: r}
}

type nopReadCloser[T any] struct{ Reader[T] }

func (nopReadCloser[T]) Close() error { return nil }
