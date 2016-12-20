package grx

import "time"

// Subscription is normally returned from any subscription with the time
// the subscription happens and the expected time to unsubscribe (if available)
type Subscription struct {
	_observable   Observable
	SubscribeAt   time.Time
	UnsubscribeAt time.Time
}

// Dispose clean up an observable and notify the unsubscribe channel, then return a Subscriptor.
func (s *Subscription) Dispose() Subscriptor {
	go func() {
		s._observable.unsubscribe()
	}()
	s.UnsubscribeAt = time.Now()
	return s
}

// Unsubscribe is an alias for Dispose
func (s *Subscription) Unsubscribe() Subscriptor {
	return s.Dispose()
}

// UnscribeIn notify the unscribe channel in d duration, then return the Subscriptor.
func (s *Subscription) UnsubscribeIn(d time.Duration) Subscriptor {
	<-time.After(d)
	return s.Dispose()
}
