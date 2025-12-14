package streams

import (
	"slices"
	"unsafe"

	"github.com/nikmy/algo/fp"
)

func ReadAll[T any](r Reader[T], to []T) ([]T, error) {
	noReadAttempts := 0
	for {
		n, err := r.Read(to[len(to):cap(to)])
		to = slices.Grow(to, n)
		if err == EOS {
			return to, nil
		}
		if err != nil {
			return to, err
		}

		if n == 0 {
			noReadAttempts++
			if noReadAttempts > maxAttemptsWithoutProgress {
				return to, ErrNoReadProgress
			}
			continue
		}
		noReadAttempts = 0

		to = slices.Grow(to, 1)
	}
}

func Copy[T any](r Reader[T], w Writer[T]) error {
	return CopyBuffer(defaultLenBuffer[T](), r, w)
}

const maxAttemptsWithoutProgress = 10

func CopyBuffer[T any](buf []T, r Reader[T], w Writer[T]) error {
	noReadAttempts := 0
	for {
		read, err := r.Read(buf)
		if err == EOS {
			return nil
		}
		if err != nil {
			return err
		}

		if read == 0 {
			noReadAttempts++
			if noReadAttempts > maxAttemptsWithoutProgress {
				return ErrNoReadProgress
			}
			continue
		}
		noReadAttempts = 0

		err = writeAll(w, buf[:read])
		if err != nil {
			return err
		}
	}
}

func writeAll[T any](w Writer[T], buf []T) error {
	noProgressAttempts := 0
	for {
		written, err := w.Write(buf)
		if err != nil {
			return err
		}

		if written == 0 {
			noProgressAttempts++
			if noProgressAttempts > maxAttemptsWithoutProgress {
				return ErrNoWriteProgress
			}
			continue
		}
		noProgressAttempts = 0

		if written < len(buf) {
			buf = buf[written:]
		}
	}
}

type readerFunc[T any] func([]T) (int, error)

func (f readerFunc[T]) Read(p []T) (int, error) { return f(p) }

// Map returns a reader, that in its Read method returns
// number of processed objects and error from original Reader.
func Map[From, To any](mapFunc func(From) To, r Reader[From]) Reader[To] {
	return readerFunc[To](func(p []To) (int, error) {
		if len(p) == 0 {
			return 0, nil
		}

		buf := make([]From, len(p))

		n, err := r.Read(buf)
		for i, v := range buf[:n] {
			p[i] = mapFunc(v)
		}

		return n, err
	})
}

func Filter[T any](filterFunc func(T) bool, r Reader[T]) Reader[T] {
	return readerFunc[T](func(p []T) (int, error) {
		if len(p) == 0 {
			return 0, nil
		}

		buf := make([]T, len(p))

		n, err := r.Read(buf)
		filtered := 0
		for _, v := range buf[:n] {
			if filterFunc(v) {
				p[filtered] = v
				filtered++
			}
		}

		return filtered, err
	})
}

// Reduce returns a writer that writes reduced value into p[0].
func Reduce[E, R any](reduceFunc func(R, E) R, init R, r Reader[E]) Reader[R] {
	return readerFunc[R](func(p []R) (int, error) {
		if len(p) == 0 {
			return 0, nil
		}

		buf := make([]E, len(p))

		n, err := r.Read(buf)
		if err != nil {
			return 0, err
		}

		p[0] = init
		for _, e := range buf[:n] {
			reduceFunc(p[0], e)
		}

		return min(1, n), nil
	})
}

func Uniq[T comparable](r Reader[T]) Reader[T] {
	uniq := map[T]struct{}{}

	return readerFunc[T](func(p []T) (int, error) {
		if len(p) == 0 {
			return 0, nil
		}

		buf := make([]T, len(p))

		n, err := r.Read(buf)
		if err != nil {
			return 0, err
		}

		uniques := 0
		for _, v := range buf[:n] {
			if _, ok := uniq[v]; ok {
				continue
			}
			p[uniques] = v
			uniques++
		}

		return uniques, nil
	})
}

func UniqFunc[T any](equalFunc func(T, T) bool, r Reader[T]) Reader[T] {
	uniq := defaultLenBuffer[T]()

	return readerFunc[T](func(p []T) (int, error) {
		if len(p) == 0 {
			return 0, nil
		}

		buf := make([]T, len(p))

		n, err := r.Read(buf)

		uniques := 0
		for _, v := range buf[:n] {
			if slices.ContainsFunc(uniq, fp.Carry2(equalFunc, v)) {
				continue
			}
			p[uniques] = v
			uniques++
		}

		return uniques, err
	})
}

func defaultLenBuffer[T any]() []T {
	const (
		minSize     = 4
		maxBytesLen = 4095
	)

	bufSize := max(minSize, maxBytesLen/unsafe.Sizeof(fp.Empty[T]()))
	return make([]T, bufSize)
}
