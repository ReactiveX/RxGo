package rxgo

import (
	"context"
	"reflect"
	"time"
)

type (
	// Item is a wrapper having either a value or an error.
	Item struct {
		V interface{}
		E error
	}

	// TimestampItem attach a timestamp to an item.
	TimestampItem struct {
		Timestamp time.Time
		V         interface{}
	}

	// CloseChannelStrategy indicates a strategy on whether to close a channel.
	CloseChannelStrategy uint32
)

const (
	// LeaveChannelOpen indicates to leave the channel open after completion.
	LeaveChannelOpen CloseChannelStrategy = iota
	// CloseChannel indicates to close the channel open after completion.
	CloseChannel
)

// Of creates an item from a value.
func Of(i interface{}) Item {
	return Item{V: i}
}

// Error creates an item from an error.
func Error(err error) Item {
	return Item{E: err}
}

// SendItems is an utility function that send a list of interface{} and indicate a strategy on whether to close
// the channel once the function completes.
func SendItems(ch chan<- Item, strategy CloseChannelStrategy, items ...interface{}) {
	if strategy == CloseChannel {
		defer close(ch)
	}
	for _, currentItem := range items {
		switch item := currentItem.(type) {
		default:
			rt := reflect.TypeOf(item)
			switch rt.Kind() {
			default:
				ch <- Of(item)
			case reflect.Slice:
				s := reflect.ValueOf(currentItem)
				for i := 0; i < s.Len(); i++ {
					currentItem := s.Index(i).Interface()
					switch item := currentItem.(type) {
					default:
						ch <- Of(item)
					case error:
						ch <- Error(item)
					}
				}
			}
		case error:
			ch <- Error(item)
		}
	}
}

// Error checks if an item is an error.
func (i Item) Error() bool {
	return i.E != nil
}

// SendBlocking sends an item and blocks until it is sent.
func (i Item) SendBlocking(ch chan<- Item) {
	ch <- i
}

// SendCtx sends an item and blocks until it is sent or a context canceled.
// It returns a boolean to indicate whether the item was sent.
func (i Item) SendCtx(ctx context.Context, ch chan<- Item) bool {
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
