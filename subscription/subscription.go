package subscription

import "time"

// Subscription is usually returned from any subscription
type Subscription struct {
	SubscribeAt   time.Time
	UnsubscribeAt time.Time
	Error         error
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
