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


