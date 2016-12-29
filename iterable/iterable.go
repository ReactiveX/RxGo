package iterable

import (
	"github.com/jochasinga/grx/bases"
	"github.com/jochasinga/grx/emittable"
	"github.com/jochasinga/grx/errors"
)

// Emittable is a high-level alias for empty interface which may any underlying type including error
type Iterable chan bases.Emitter

func (it Iterable) Next() (bases.Emitter, error) {
	if next, ok := <-(chan bases.Emitter)(it); ok {
		return next, nil
	}
	return nil, errors.New(errors.EndOfIteratorError)
}

func From(any []interface{}) Iterable {
	it := make(chan bases.Emitter, len(any))
	go func() {
		for _, val := range any {
			it <- emittable.From(val)
		}
		close(it)
	}()
	return Iterable(it)
}
