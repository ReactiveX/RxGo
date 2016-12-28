package subscription

import (
	"time"

	"github.com/jochasinga/grx/bases"
	"github.com/jochasinga/grx/subject"
)

// Subscription is usually returned from any subscription
type Subscription struct {
	observable    *subject.Subject
	subscribeAt   time.Time
	unsubscribeAt time.Time
}

// DefaultSubscription is a default Subscription
var DefaultSubscription = &Subscription{}

func New(fs ...func(*Subscription)) *Subscription {
	s := DefaultSubscription
	if len(fs) > 0 {
		for _, f := range fs {
			f(s)
		}
	}
	return s
}

func (s *Subscription) Stream() bases.Stream {
	return s.observable.Stream
}

// Dispose cleans up an observable and notify its unsubscribe channel and return a Subscriptor
func (s *Subscription) Dispose() {
	go func() {
		s.observable.Unsubscribe()
	}()
}

func (s *Subscription) Subscribe() bases.Subscriptor {
	s.subscribeAt = time.Now()
	return s
}

func (s *Subscription) SubscribeAt() time.Time {
	return s.subscribeAt
}

// Unsubscribe is an alias for Dispose
func (s *Subscription) Unsubscribe() bases.Subscriptor {
	s.unsubscribeAt = time.Now()
	return s
}

func (s *Subscription) UnsubscribeAt() time.Time {
	return s.unsubscribeAt
}

// UnscribeIn notify the unsubscribe channel in d duration, then return the Subscriptor
func (s *Subscription) UnsubscribeIn(d time.Duration) <-chan bases.Subscriptor {
	out := make(chan bases.Subscriptor)
	go func() {
		<-time.After(d)
		out <- s.Unsubscribe()
		close(out)
	}()
	return out
}

// UnscribeOn unsubscribes an Observable
func (s *Subscription) UnsubscribeOn(sig chan struct{}, timeout time.Duration) <-chan bases.Subscriptor {
	out := make(chan bases.Subscriptor)

	if timeout >= 0 {
		go func() {
			select {
			case <-time.After(timeout):
				return
			case <-sig:
				out <- s.Unsubscribe()
				return
			}
		}()
	}
	return out
}
