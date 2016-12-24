package bang

// Notifier consists of a pair of channels used in notifying an Observable
type Notifier struct {
	unsubscribe chan struct{}
	done        chan struct{}
}

var Bang = func() struct{}{} {
    return struct{}{}
}()

func (notif *Notifier) Done() {
	go func() {
		notif.done <- Bang
		close(notif.done)
	}()
}

func (notif *Notifier) Unsubscribe() {
	go func() {
		notif.unsubscribe <- Bang
		close(notif.unsubscribe)
	}()
}

func New() *Notifier {
	return &Notifier{
		unsubscribe: make(chan struct{}, 1),
		done:        make(chan struct{}, 1),
	}
}
