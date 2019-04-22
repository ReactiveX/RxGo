package rxgo

import "context"

type Iterable interface {
	Iterator(ctx context.Context) Iterator
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

type iterableFromFunc struct {
	f func(chan interface{})
}

func (it *iterableFromFunc) Iterator(ctx context.Context) Iterator {
	out := make(chan interface{})
	go it.f(out)
	return newIteratorFromChannel(out)
}

func (it *iterableFromChannel) Iterator(ctx context.Context) Iterator {
	return newIteratorFromChannel(it.ch)
}

func (it *iterableFromSlice) Iterator(ctx context.Context) Iterator {
	return newIteratorFromSlice(it.s)
}

func (it *iterableFromRange) Iterator(ctx context.Context) Iterator {
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

func newIterableFromFunc(f func(chan interface{})) Iterable {
	return &iterableFromFunc{
		f: f,
	}
}
