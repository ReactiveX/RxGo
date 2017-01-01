package fx

import (
	"github.com/jochasinga/grx/bases"
)

type (

	// DirectiveFunc defines a func that should be passed to the observable.Start method,
	// and represents a simple func that takes no arguments and return a bases.Emitter type.
	DirectiveFunc func() bases.Emitter

	// MappableFunc defines a func that should be passed to the Map operator.
	MappableFunc func(bases.Emitter) bases.Emitter

	// FilterableFunc defines a func that should be passed to the Filter operator.
	FilterableFunc func(bases.Emitter) bool
)
