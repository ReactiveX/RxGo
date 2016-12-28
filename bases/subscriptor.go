package bases

import "time"

// Subscriptor is an interface type which Subscription implements.
type Subscriptor interface {
	Dispose()
	Subscribe() Subscriptor
	SubscribeAt() time.Time
	Unsubscribe() Subscriptor
	UnsubscribeAt() time.Time
	UnsubscribeIn(time.Duration) <-chan Subscriptor
	UnsubscribeOn(chan struct{}, time.Duration) <-chan Subscriptor
}
