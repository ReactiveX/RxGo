package rxgo

// Subscription allows subscriber to unsubscribe from the subject
type Subscription interface {
	GetId() int
	Unsubscribe()
}

// subscription implementation of the Subscription interface
type subscription struct {
	id      int
	subject ISubject
}

func NewSubscription(id int, subject ISubject) *subscription {
	res := subscription{
		id:      id,
		subject: subject,
	}

	return &res
}

func (s *subscription) GetId() int {
	return s.id
}

// Unsubscribe called by a subscriber to remove itself
func (s *subscription) Unsubscribe() {
	s.subject.Unsubscribe(s.id)
}
