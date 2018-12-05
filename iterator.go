package rxgo

type Iterator interface {
	Next() bool
	Value() interface{}
}

type iteratorFromChannel struct {
	item interface{}
	ch   chan interface{}
}

func (s *iteratorFromChannel) Next() bool {
	if v, ok := <-s.ch; ok {
		s.item = v
		return true
	}

	return false
}

func (s *iteratorFromChannel) Value() interface{} {
	return s.item
}

func NewIteratorFromChannel(ch chan interface{}) Iterator {
	return &iteratorFromChannel{
		ch: ch,
	}
}
