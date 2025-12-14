package buffer

import "io"

type Reader interface {
	io.ReadCloser
}
