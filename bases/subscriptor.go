package bases

import "time"

// Subscriptor is an interface type which Subscription complies to.
type Subscriptor interface {
	Dispose() Subscriptor
	Unsubscribe() Subscriptor
	UnsubscribeIn(time.Duration) (<-chan Subscriptor)
    UnsubscribeOn(chan struct{}{}) (<-chan Subscriptor)
}
