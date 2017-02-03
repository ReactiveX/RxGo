package subscription

import "time"

// Subscription is usually returned from any subscription
type Subscription struct {
	SubscribeAt   time.Time
	UnsubscribeAt time.Time
	Error         error
	//term          chan struct{}
}

// DefaultSubscription is a default Subscription.
var DefaultSubscription = Subscription{}

// New creates a DefaultSubscription.
func New() Subscription {
	return DefaultSubscription
}

// Err returns an error recorded from a stream.
func (s Subscription) Err() error {
	return s.Error
}

// Subscribe records the time of subscription.
func (s Subscription) Subscribe() Subscription {
	s.SubscribeAt = time.Now()
	return s
}

// Unsubscribe records the time of unsubscription.
func (s Subscription) Unsubscribe() Subscription {
	s.UnsubscribeAt = time.Now()
	return s
}

/* TODO:
func (s Subscription) Dispose() {
	s.term <- struct{}{}
	close(s.term)
}

func (s Subscription) Terminated() chan struct{} {
	return s.term
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
*/
