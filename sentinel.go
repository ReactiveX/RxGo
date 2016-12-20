package grx

type Sentinel interface {
	OnNext(interface{})
	OnError(error)
	OnDone()
}

// Sentinel is an interface for any object that implements the following method set.
//type Sentinel interface {
type Observer interface {
	Sentinel

	/******* Unexported *******/
	setInnerObservableTo(Observable) Observable
	getInnerObservable() Observable
}
