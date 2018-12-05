package rxgo

type Iterator interface {
	Next() bool
	Value() interface{}
}

type iteratorFromChannel struct {
	item interface{}
	ch   chan interface{}
}

type iteratorFromSlice struct {
	index int
	s     []interface{}
}

func (it *iteratorFromChannel) Next() bool {
	if v, ok := <-it.ch; ok {
		it.item = v
		return true
	}

	return false
}

func (it *iteratorFromChannel) Value() interface{} {
	return it.item
}

func (it *iteratorFromSlice) Next() bool {
	it.index = it.index + 1
	if it.index >= len(it.s) {
		return false
	} else {
		return true
	}
}

func (it *iteratorFromSlice) Value() interface{} {
	return it.s[it.index]
}

func NewIteratorFromChannel(ch chan interface{}) Iterator {
	return &iteratorFromChannel{
		ch: ch,
	}
}

func NewIteratorFromSlice(s []interface{}) Iterator {
	return &iteratorFromSlice{
		index: -1,
		s:     s,
	}
}
