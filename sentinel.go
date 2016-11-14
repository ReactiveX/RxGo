package grx

// Sentinel is an interface for any object that implements the following method set.
type Sentinel interface {
	OnNext(interface{})
	OnError(error)
	OnDone()
}
