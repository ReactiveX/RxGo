package observer

import (
	"github.com/jochasinga/grx/event"
)

// Observer is a "sentinel" object consisting of three methods to handle event stream.
type Observer struct {
	OnNext      func(e *event.Event)
	OnError     func(e *event.Event)
	OnCompleted func(e *event.Event)
}

// New constructs a new Observer instance with default handlers.
func New() *Observer {
	return &Observer{
		OnNext: func(e *event.Event) {},
		OnError: func(e *event.Event) { panic(e.Error) },
		OnCompleted: func(e *event.Event) {},
	}
}


