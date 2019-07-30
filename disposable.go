package rxgo

// Disposable allows to dispose a subscription
type Disposable interface {
	Dispose()
	IsDisposed() bool
	Notify(chan<- struct{})
}
