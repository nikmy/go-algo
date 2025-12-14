package topcache

type entry[K comparable] struct {
	Key K
	Cnt int64

	w window
}

func (e *entry[K]) update(delta, utc int64) {
	e.w.Update(delta, utc)
	e.Cnt += delta
}
