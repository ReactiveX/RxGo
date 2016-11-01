package event

// Event either emits a value, an error, or notify as completed.
type Event struct {
	Value     interface{}
	Error     error
	Completed bool
}

// New returns an instance of Event. 
func New() *Event {
	return &Event{}
}


