package rxgo

type Iterable interface {
	Iterator() Iterator
}

type iterableFromChannel struct {
	ch chan interface{}
}

func (it *iterableFromChannel) Iterator() Iterator {
	return newIteratorFromChannel(it.ch)
}

func newIterableFromChannel(ch chan interface{}) Iterable {
	return &iterableFromChannel{
		ch: ch,
	}
}

type iterableFromSlice struct {
	s []interface{}
}

func (it *iterableFromSlice) Iterator() Iterator {
	return newIteratorFromSlice(it.s)
}

func newIterableFromSlice(s []interface{}) Iterable {
	return &iterableFromSlice{
		s: s,
	}
}
