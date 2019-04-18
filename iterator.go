package rxgo

import (
	"context"
	"github.com/reactivex/rxgo/errors"
)

type Iterator interface {
	cancel()
	Next() (interface{}, error)
}

type iteratorFromChannel struct {
	ch         chan interface{}
	ctx        context.Context
	cancelFunc context.CancelFunc
}

type iteratorFromRange struct {
	current int
	end     int // Included
}

type iteratorFromSlice struct {
	index int
	s     []interface{}
}

func (it *iteratorFromChannel) cancel() {
	it.cancelFunc()
}

func (it *iteratorFromChannel) Next() (interface{}, error) {
	select {
	case <-it.ctx.Done():
		return nil, errors.New(errors.CancelledIteratorError)
	case next, ok := <-it.ch:
		if ok {
			return next, nil
		}
		return nil, errors.New(errors.EndOfIteratorError)
	}
}

func (it *iteratorFromRange) cancel() {
	// TODO
}

func (it *iteratorFromRange) Next() (interface{}, error) {
	it.current++
	if it.current <= it.end {
		return it.current, nil
	}
	return nil, errors.New(errors.EndOfIteratorError)
}

func (it *iteratorFromSlice) cancel() {
	// TODO
}

func (it *iteratorFromSlice) Next() (interface{}, error) {
	it.index++
	if it.index < len(it.s) {
		return it.s[it.index], nil
	}
	return nil, errors.New(errors.EndOfIteratorError)
}

func newIteratorFromChannel(ch chan interface{}) Iterator {
	ctx, cancel := context.WithCancel(context.Background())
	return &iteratorFromChannel{
		ch:         ch,
		ctx:        ctx,
		cancelFunc: cancel,
	}
}

func newIteratorFromRange(start, end int) Iterator {
	return &iteratorFromRange{
		current: start,
		end:     end,
	}
}

func newIteratorFromSlice(s []interface{}) Iterator {
	return &iteratorFromSlice{
		index: -1,
		s:     s,
	}
}
