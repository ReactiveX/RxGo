package grx

import "errors"

var (
	EndOfIteratorError       = errors.New("End of Iterator")
	NilObservableError       = errors.New("Observable is nil")
	NilObservableCError      = errors.New("Observable's internal channel is nil")
	NilObserverError         = errors.New("Observer is nil")
	UndefinedObservableError = errors.New("Undefined error")
)
