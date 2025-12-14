package topcache

import "sync"

func New[K comparable](topSize int, windowLen int, gran int64) *cache[K] {
	return &cache[K]{
		ordered: newTopList[K](topSize),
		heap:    make(map[K]*entry[K]),
		gran:    gran,
		wLen:    windowLen,
	}
}

type cache[K comparable] struct {
	lock sync.RWMutex

	ordered *topList[K]
	heap    map[K]*entry[K]

	gran int64
	wLen int
}

func (c *cache[K]) Top(k int) []K {
	c.lock.RLock()
	defer c.lock.RUnlock()
	return c.ordered.Top(k)
}

func (c *cache[K]) Update(key K, delta int64, utc int64) {
	c.lock.Lock()
	defer c.lock.Unlock()

	e, isSmall := c.heap[key]

	if !isSmall && !c.ordered.Has(key) {
		e = c.newEntry(key, delta, utc)

		if e.Cnt <= c.ordered.Min() {
			c.heap[key] = e
			return
		}

		c.pushTop(e)
		return
	}

	if !isSmall {
		c.ordered.Update(key, delta, utc)
		// maybe there is an entry in heap
		// having greater count, but we ignore
		// this weird case for optimization
		return
	}

	e.Cnt += delta
	if e.Cnt < c.ordered.Min() {
		return
	}

	delete(c.heap, key)
	c.pushTop(e)
}

func (c *cache[K]) newEntry(key K, cnt int64, utc int64) *entry[K] {
	w := window{
		gran: c.gran,
		data: make([]item, c.wLen),
	}
	w.Update(cnt, utc)

	return &entry[K]{
		Key: key,
		Cnt: cnt,
		w:   w,
	}
}

func (c *cache[K]) pushTop(e *entry[K]) {
	evicted := c.ordered.Push(e)
	if evicted != nil {
		c.heap[evicted.Key] = evicted
	}
}
