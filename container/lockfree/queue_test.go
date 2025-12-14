package lockfree

import (
	"testing"
)

func BenchmarkQueue_CompareWithChannel(b *testing.B) {
	compareParSeq(b, func(size uint32) bufferImpl { return NewQueue[int]() })
}
