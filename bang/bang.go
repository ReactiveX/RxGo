package bang

// Notifier consists of a pair of channels used in notifying an Observable
type Notifier struct {
	Unsubscribed chan struct{}
	IsDone       chan struct{}
}

var Bang = func() struct{} {
	return struct{}{}
}()

func (notif *Notifier) Done() {
	go func() {
		notif.IsDone <- Bang
		close(notif.IsDone)
	}()
}

func (notif *Notifier) Unsubscribe() {
	go func() {
		notif.Unsubscribed <- Bang
		close(notif.Unsubscribed)
	}()
}

func New() *Notifier {
	return &Notifier{
		Unsubscribed: make(chan struct{}, 1),
		IsDone:       make(chan struct{}, 1),
	}
}
