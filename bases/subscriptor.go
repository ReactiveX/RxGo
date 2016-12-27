package bases

import "time"

// Subscriptor is an interface type which Subscription implements.
type Subscriptor interface {
	Dispose() Subscriptor
	Subscribe() Subscriptor
	Unsubscribe() Subscriptor
	UnsubscribeIn(time.Duration) <-chan Subscriptor
	UnsubscribeOn(chan struct{}, time.Duration) <-chan Subscriptor
}
