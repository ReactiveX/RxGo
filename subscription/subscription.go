package subscription

import (
	"time"

	"github.com/jochasinga/bases"
	"github.com/jochasinga/subject"
)

// Subscription is usually returned from any subscription
type Subscription struct {
	observable    *subject.Subject
	SubscribeAt   time.Time
	UnsubscribeAt time.Time
}

// DefaultSubscription is a default Subscription
var DefaultSubscription = &Subscription{}

func New(fs ...func(*Subscription)) {
	s := DefaultSubscription
	if len(fs) > 0 {
		for _, f := range fs {
			f(s)
		}
	}
	return s
}

// Dispose cleans up an observable and notify its unsubscribe channel and return a Subscriptor
func (s *Subscription) Dispose() bases.Subscriptor {
	go func() {
		s.observable.Unsubscribe()
		s.UnsubscribeAt = time.Now()
	}()
	return s
}

// Unsubscribe is an alias for Dispose
func (s *Subscription) Unsubscribe() bases.Subscriptor {
	return s.Dispose()
}

// UnscribeIn notify the unsubscribe channel in d duration, then return the Subscriptor
func (s *Subscription) UnsubscribeIn(d time.Duration) <-chan bases.Subscriptor {
	out := make(chan bases.Subscriptor)
	go func() {
		<-time.After(d)
		out <- s.Dispose()
		close(out)
	}()
	return out
}

// UnscribeOn unsubscribes an Observable
func (s *Subscription) UnsubscribeOn(sig chan struct{}, timeout ...time.Duration) <-chan Subscriptor {
	out := make(chan Subscriptor)
	go func() {
		select {
		case <-time.After(timeout):
		case <-sig:
			out <- s.Dispose()
		}
	}()
	return out

}
