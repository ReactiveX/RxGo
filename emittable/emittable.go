package emittable

import (
	"github.com/jochasinga/grx/bases"
)

// Emittable is a high-level alias for empty interface which may any underlying type including error
type Emittable struct {
	value interface{}
}

func (e *Emittable) Emit() (bases.Item, error) {
	switch e := e.value.(type) {
	case error:
		return nil, e
	default:
		break
	}
	return bases.Item(e), nil
}

func New(any interface{}) *Emittable {
	return &Emittable{any}
}
