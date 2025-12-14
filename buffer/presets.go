package buffer

const (
	Byte     = 1
	Kilobyte = 10 << Byte
	Megabyte = 10 << Kilobyte

	// number of shards in pool
	nShards = 6

	// files size distribution
	tinySize   = 512 * Byte
	smallSize  = 1 * Kilobyte
	mediumSize = 4 * Kilobyte
	largeSize  = 1 * Megabyte
	hugeSize   = 8 * Megabyte
)

type Preset [nShards - 1]int

var (
	// FileReadingPreset is good for caching buffers for
	// rereading runtime configuration files.
	FileReadingPreset = Preset{tinySize, smallSize, mediumSize, largeSize, hugeSize}

	// SmallBuffersPreset contains sizes of small pieces. It includes
	// standard 64B cacheline size, 4KB page size, 32KB entire L1 size
	// and 64KB — maximum size that does not escape to heap in Go.
	SmallBuffersPreset = Preset{64 * Byte, Kilobyte, 4 * Kilobyte, 32 * Kilobyte, 64 * Kilobyte}

	// LargeBuffersPreset contains sizes of large pieces. It includes
	// 256KB, 1MB, 2MB (standard L2 cache size), 4MB and 8MB (L3).
	LargeBuffersPreset = Preset{256 * Kilobyte, Megabyte, 2 * Megabyte, 4 * Megabyte, 8 * Megabyte}
)
