package rxgo

type Iterable interface {
	Iterator() Iterator
}

type iterableFromChannel struct {
	ch chan interface{}
}

type iterableFromSlice struct {
	s []interface{}
}

type iterableFromRange struct {
	start int
	count int
}

func (it *iterableFromChannel) Iterator() Iterator {
	return newIteratorFromChannel(it.ch)
}

func (it *iterableFromSlice) Iterator() Iterator {
	return newIteratorFromSlice(it.s)
}

func (it *iterableFromRange) Iterator() Iterator {
	return newIteratorFromRange(it.start-1, it.start+it.count)
}

func newIterableFromChannel(ch chan interface{}) Iterable {
	return &iterableFromChannel{
		ch: ch,
	}
}

func newIterableFromSlice(s []interface{}) Iterable {
	return &iterableFromSlice{
		s: s,
	}
}

func newIterableFromRange(start, count int) Iterable {
	return &iterableFromRange{
		start: start,
		count: count,
	}
}
