package rxgo

import "context"

type (
	// Item is a wrapper having either a value or an error.
	Item struct {
		V interface{}
		E error
	}
)

// Of creates an item from a value.
func Of(i interface{}) Item {
	return Item{V: i}
}

// Error creates an item from an error.
func Error(err error) Item {
	return Item{E: err}
}

// Error checks if an item is an error.
func (i Item) Error() bool {
	return i.E != nil
}

// SendBlocking sends an item and blocks until it is sent.
func (i Item) SendBlocking(ch chan<- Item) {
	ch <- i
}

// SendContext sends an item and blocks until it is sent or a context canceled.
// It returns a boolean to indicate whether the item was sent.
func (i Item) SendContext(ctx context.Context, ch chan<- Item) bool {
	select {
	case <-ctx.Done():
		return false
	case ch <- i:
		return true
	}
}

// SendNonBlocking sends an item without blocking.
// It returns a boolean to indicate whether the item was sent.
func (i Item) SendNonBlocking(ch chan<- Item) bool {
	select {
	default:
		return false
	case ch <- i:
		return true
	}
}
