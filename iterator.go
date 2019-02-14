package rxgo

import "github.com/reactivex/rxgo/errors"

type Iterator interface {
	Next() (interface{}, error)
}

type iteratorFromChannel struct {
	ch chan interface{}
}

type iteratorFromSlice struct {
	index int
	s     []interface{}
}

type iteratorFromRange struct {
	current int
	end     int // Included
}

func (it *iteratorFromChannel) Next() (interface{}, error) {
	if next, ok := <-it.ch; ok {
		return next, nil
	}

	return nil, errors.New(errors.EndOfIteratorError)
}

func (it *iteratorFromSlice) Next() (interface{}, error) {
	it.index = it.index + 1
	if it.index < len(it.s) {
		return it.s[it.index], nil
	} else {
		return nil, errors.New(errors.EndOfIteratorError)
	}
}

func (it *iteratorFromRange) Next() (interface{}, error) {
	it.current = it.current + 1
	if it.current <= it.end {
		return it.current, nil
	} else {
		return nil, errors.New(errors.EndOfIteratorError)
	}
}

func newIteratorFromChannel(ch chan interface{}) Iterator {
	return &iteratorFromChannel{
		ch: ch,
	}
}

func newIteratorFromSlice(s []interface{}) Iterator {
	return &iteratorFromSlice{
		index: -1,
		s:     s,
	}
}

func newIteratorFromRange(start, end int) Iterator {
	return &iteratorFromRange{
		current: start,
		end:     end,
	}
}
