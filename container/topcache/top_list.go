package topcache

func newTopList[K comparable](size int) *topList[K] {
	return &topList[K]{
		data: make([]*entry[K], 0, size),
		fidx: make(map[K]int),
	}
}

type topList[K comparable] struct {
	fidx map[K]int
	data []*entry[K]
}

func (l *topList[K]) Top(k int) []K {
	k = min(k, len(l.data))

	r := make([]K, 0, k)
	for _, entry := range l.data[:k] {
		r = append(r, entry.Key)
	}

	return r
}

func (l *topList[K]) Min() int64 {
	return l.data[len(l.data)-1].Cnt
}

func (l *topList[K]) Push(entry *entry[K]) (evicted *entry[K]) {
	if len(l.data) == cap(l.data) && entry.Cnt <= l.Min() {
		return entry
	}

	var i int
	if len(l.data) == cap(l.data) {
		evicted, l.data[cap(l.data)-1] = l.data[cap(l.data)-1], entry
		delete(l.fidx, evicted.Key)
		i = cap(l.data) - 1
	} else {
		l.data = append(l.data, entry)
		i = len(l.data) - 1
	}

	l.fidx[entry.Key] = i
	l.sift(i)
	return
}

func (l *topList[K]) Update(key K, delta int64, utc int64) int64 {
	i, ok := l.fidx[key]
	if !ok {
		panic("update: entry not found")
	}

	l.data[i].update(delta, utc)
	l.sift(i)
	return l.data[i].Cnt
}

func (l *topList[K]) Has(key K) bool {
	_, ok := l.fidx[key]
	return ok
}

func (l *topList[K]) sift(i int) {
	p := i - 1
	for p >= 0 && l.data[i].Cnt > l.data[p].Cnt {
		l.swap(p, i)
		i--
		p--
	}

	n := i + 1
	for n < len(l.data) && l.data[n].Cnt > l.data[i].Cnt {
		l.swap(i, n)
		i++
		n++
	}
}

func (l *topList[K]) swap(i, j int) {
	l.data[i], l.data[j] = l.data[j], l.data[i]
	l.fidx[l.data[i].Key], l.fidx[l.data[j].Key] = i, j
}
