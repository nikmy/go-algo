package buffer

import (
	"runtime"
	"sync/atomic"
)

type Pool = *shardedPool

func NewPool(maxBufferSize int, preset Preset) Pool {
	p := shardedPool{
		{bufferSize: preset[0]},
		{bufferSize: preset[1]},
		{bufferSize: preset[2]},
		{bufferSize: preset[3]},
		{bufferSize: preset[4]},
		{bufferSize: maxBufferSize},
	}

	for i := range p {
		p[i].maxPoolSize = maxBufferSize / p[i].bufferSize
	}

	return &p
}

type shardedPool [nShards]sizedPool

func (p *shardedPool) Get(n int) *Buffer {
	buf := &Buffer{
		pool: p,
		data: p.getPoolForSize(n).get(n),
	}
	buf.Resize(n)
	return buf
}

func (p *shardedPool) Release() {
	for i := range p {
		p[i].release()
	}
}

func (p *shardedPool) AggregateStats() PoolStats {
	if p == nil {
		return PoolStats{}
	}

	var stats [nShards]PoolStats
	for i := 0; i < nShards; i++ {
		p[i].readStats(&stats[i])
	}

	aggr := PoolStats{}
	for _, stat := range &stats {
		aggr.TotalAlloc += stat.TotalAlloc
		aggr.TotalCrops += stat.TotalCrops
		aggr.TotalOverflows += stat.TotalOverflows
		aggr.InUse += stat.InUse
	}

	return aggr
}

func (p *shardedPool) getPoolForSize(n int) *sizedPool {
	return &p[p.category(n)]
}

func (p *shardedPool) category(n int) int {
	switch {
	case n <= p[0].bufferSize:
		return 0
	case n <= p[1].bufferSize:
		return 1
	case n <= p[2].bufferSize:
		return 2
	case n <= p[3].bufferSize:
		return 3
	case n <= p[4].bufferSize:
		return 4
	default:
		return 5
	}
}

type sizedPool struct {
	spinlock

	buffers     [][]byte
	bufferSize  int
	maxPoolSize int

	stats PoolStats
}

func (p *sizedPool) get(n int) []byte {
	p.Lock()
	if len(p.buffers) == 0 {
		p.Unlock()
		return p.new(n)
	}

	last := len(p.buffers) - 1

	buf := p.buffers[last]
	p.buffers[last] = nil
	p.buffers = p.buffers[:last]

	p.Unlock()

	return buf
}

func (p *sizedPool) put(buf []byte) {
	if p.bufferSize < cap(buf) {
		buf = buf[:0:p.bufferSize]
	}

	p.Lock()
	defer p.Unlock()
	if len(p.buffers) < p.maxPoolSize {
		p.buffers = append(p.buffers, buf[:0])
	}
}

func (p *sizedPool) release() {
	p.Lock()
	defer p.Unlock()
	p.buffers = nil
}

func (p *sizedPool) new(n int) []byte {
	return make([]byte, 0, max(p.bufferSize, n))
}

func (p *sizedPool) readStats(stats *PoolStats) {
	stats.TotalAlloc = atomic.LoadInt64(&p.stats.TotalAlloc)
	stats.TotalCrops = atomic.LoadInt64(&p.stats.TotalCrops)
	stats.TotalOverflows = atomic.LoadInt64(&p.stats.TotalOverflows)
	stats.InUse = atomic.LoadInt64(&p.stats.InUse)
}

type spinlock struct {
	locked uint32
}

func (s *spinlock) Lock() {
	for {
		if atomic.CompareAndSwapUint32(&s.locked, 0, 1) {
			break
		}

		for atomic.LoadUint32(&s.locked) == 1 {
			// spin
			runtime.Gosched()
			continue
		}
	}
}

func (s *spinlock) Unlock() {
	atomic.StoreUint32(&s.locked, 0)
}
