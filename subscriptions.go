package grx

import "time"

// Subscription is normally returned from any subscription with the time
// the subscription happens and the expected time to unsubscribe, if available.
type Subscription struct {
	observable *Observable
	Subscribe time.Time
	Unsubscribe time.Time
}

func (s *Subscription) Close() {
	go func() {
		s.observable.unsubscribed <- struct{}{}
	}()
	s.Unsubscribe = time.Now()
}
