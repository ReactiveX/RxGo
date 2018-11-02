package subscription

import (
	"context"
	"time"
)

// Subscription is usually returned from any subscription
type Subscription struct {
	ctx context.Context

	SubscribeTime   time.Time
	UnsubscribeTime time.Time
	Error           error
	//term          chan struct{}
}

// DefaultSubscription is a default Subscription.
var DefaultSubscription = Subscription{
	ctx: context.Background(),
}

// New creates a DefaultSubscription.
func New() Subscription {
	return DefaultSubscription
}

//Deadline return when the emitter is completed, and whether or not a deadline is set.
func (s Subscription) Deadline() (deadline time.Time, ok bool) {
	return s.ctx.Deadline()
}

// Done returns a channel that signals when the emitter has completed.
func (s Subscription) Done() <-chan struct{} {
	return s.ctx.Done()
}

// Err returns the error that has occured for the emitter.
func (s Subscription) Err() error {
	if s.Error != nil {
		return s.Error
	}
	return s.ctx.Err()
}

// Value assigns a key to the emitter.
func (s Subscription) Value(key interface{}) interface{} {
	return s.ctx.Value(key)
}

// SubscribeAt returns the time at which the subscription was subscribed.
func (s Subscription) SubscribeAt() time.Time {
	return s.SubscribeTime
}

// Subscribe records the time of subscription.
func (s Subscription) Subscribe() Subscription {
	s.SubscribeTime = time.Now()
	return s
}

// Unsubscribe records the time of unsubscription.
func (s Subscription) Unsubscribe() Subscription {
	s.UnsubscribeTime = time.Now()
	return s
}

//WithContext constructs a subscription with the specified context.
func WithContext(ctx context.Context) Subscription {
	return Subscription{
		ctx: ctx,
	}
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
