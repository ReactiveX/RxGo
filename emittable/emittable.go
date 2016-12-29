package emittable

import (
	"github.com/jochasinga/grx/bases"
)

// Emittable is a high-level alias for empty interface which may any underlying type including error
type Emittable struct {
	value interface{}
}

var DefaultEmittable = &Emittable{}

func (e *Emittable) Emit() (bases.Item, error) {
	switch e := e.value.(type) {
	case error:
		return nil, e
	case interface{}:
		return bases.Item(e), nil
	}
	return bases.Item(e), nil
}

func From(any interface{}) *Emittable {
	return &Emittable{any}
}
