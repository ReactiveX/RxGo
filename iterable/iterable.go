// Package iterable provides an Iterable type that is capable of converting
// sequences of empty interface such as slice and channel to an Iterator.
package iterable

import (
	"reflect"

	"github.com/reactivex/rxgo/errors"
)

// Iterable converts channel and slice into an Iterator.
type Iterable <-chan interface{}

// Next returns the next element in an Iterable sequence and an
// error when it reaches the end. Next registers Iterable to Iterator.
func (it Iterable) Next() (interface{}, error) {
	if next, ok := <-it; ok {
		return next, nil
	}
	return nil, errors.New(errors.EndOfIteratorError)
}

// New creates a new Iterable from a slice or a channel or a map of empty interface.
func New(any interface{}) (Iterable, error) {
	refAny := reflect.ValueOf(any)
	switch refAny.Type().Kind() {
	case reflect.Slice:
		c := make(chan interface{}, refAny.Len())
		go func() {
			for i := 0; i < refAny.Len(); i++ {
				c <- refAny.Index(i).Interface()
			}
			close(c)
		}()
		return Iterable(c), nil
	case reflect.Map:
		keys := refAny.MapKeys()
		c := make(chan interface{}, len(keys))
		go func() {
			for _, key := range keys {
				c <- []interface{}{key.Interface(), refAny.MapIndex(key).Interface()}
			}
			close(c)
		}()
		return Iterable(c), nil
	case reflect.Chan:
		chanDir := refAny.Type().ChanDir()
		if chanDir == reflect.SendDir {
			return nil, errors.New(errors.IterableError)
		}
		if refAny.Type().Elem().Kind() == reflect.Interface {
			if chanDir == reflect.RecvDir {
				return Iterable(any.(<-chan interface{})), nil
			}
			return Iterable(any.(chan interface{})), nil
		}
		c := make(chan interface{}, refAny.Cap())
		go func() {
			for v, ok := refAny.Recv(); ok; v, ok = refAny.Recv() {
				c <- v.Interface()
			}
			close(c)
		}()
		return Iterable(c), nil
	default:
		return nil, errors.New(errors.IterableError)
	}
}
