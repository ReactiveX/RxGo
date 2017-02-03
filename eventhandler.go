package rx

// EventHandler type is implemented by all handlers and Observer.
type EventHandler interface {
	Handle(interface{})
}
