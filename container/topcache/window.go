package topcache

import "sort"

type window struct {
	data []item
	head int
	size int
	gran int64
}

type item struct {
	cnt int64
	utc int64
}

func (w *window) Update(delta, utc int64) {
	w.cleanup(utc)

	at := w.search(utc)
	if len(w.data) > 0 && at == len(w.data) {
		w.popFront()
	}

	if at < w.size {
		w.data[w.idx(at)].cnt += delta
		return
	}

	w.data[w.idx(at)] = item{cnt: delta, utc: w.round(utc)}
	w.size++
}

func (w *window) cleanup(now int64) {
	firstValid := w.idx(w.search(now - w.gran * int64(len(w.data))))
	for w.head != firstValid {
		w.popFront()
	}
}

func (w *window) search(t int64) int {
	i := sort.Search(w.size-1, func(i int) bool {
		return w.data[w.idx(i)].utc >= t
	})

	if (t - w.data[w.idx(i)].utc) >= w.gran {
		i++
	}

	return i
}

func (w *window) popFront() {
	w.head = (w.head + 1) % len(w.data)
	w.size--
}

func (w *window) round(t int64) int64 {
	return t / w.gran * w.gran
}

func (w *window) idx(i int) int {
	return (w.head + i) % len(w.data)
}
