package rxgo

type Iterable interface {
	Next() <-chan Item
}

type Source struct {
	next <-chan Item
}

func newIterable(next <-chan Item) Iterable {
	return &Source{next: next}
}

func (s *Source) Next() <-chan Item {
	return s.next
}
