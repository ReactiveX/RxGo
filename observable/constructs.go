package observable

import (
	"github.com/jochasinga/grx/event"
)

// Observable is a stream of events
type Observable struct {
	Stream    chan *event.Event
}

