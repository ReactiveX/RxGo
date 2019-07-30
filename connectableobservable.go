package rxgo

import (
	"context"
	"sync"

	"github.com/pkg/errors"
)

// ConnectableObservable is an observable that send items once an observer is connected
type ConnectableObservable interface {
	Observable
	Connect() Observer
}

type connectableObservable struct {
	observable     Observable
	observersMutex sync.Mutex
	observers      []Observer
}

func newConnectableObservableFromObservable(observable Observable) ConnectableObservable {
	return &connectableObservable{
		observable: observable,
	}
}

func (c *connectableObservable) Iterator(ctx context.Context) Iterator {
	return c.observable.Iterator(context.Background())
}

func (c *connectableObservable) All(predicate Predicate) Single {
	return c.observable.All(predicate)
}

func (c *connectableObservable) AverageFloat32() Single {
	return c.observable.AverageFloat32()
}

func (c *connectableObservable) AverageFloat64() Single {
	return c.observable.AverageFloat64()
}

func (c *connectableObservable) AverageInt() Single {
	return c.observable.AverageInt()
}

func (c *connectableObservable) AverageInt8() Single {
	return c.observable.AverageInt8()
}

func (c *connectableObservable) AverageInt16() Single {
	return c.observable.AverageInt16()
}

func (c *connectableObservable) AverageInt32() Single {
	return c.observable.AverageInt32()
}

func (c *connectableObservable) AverageInt64() Single {
	return c.observable.AverageInt64()
}

func (c *connectableObservable) BufferWithCount(count, skip int) Observable {
	return c.observable.BufferWithCount(count, skip)
}

func (c *connectableObservable) BufferWithTime(timespan, timeshift Duration) Observable {
	return c.observable.BufferWithTime(timespan, timeshift)
}

func (c *connectableObservable) BufferWithTimeOrCount(timespan Duration, count int) Observable {
	return c.observable.BufferWithTimeOrCount(timespan, count)
}

func (c *connectableObservable) Connect() Observer {
	out := NewObserver()
	go func() {
		it := c.observable.Iterator(context.Background())
		for {
			if item, err := it.Next(context.Background()); err == nil {
				c.observersMutex.Lock()
				for _, observer := range c.observers {
					c.observersMutex.Unlock()
					select {
					case observer.getItemChannel() <- item:
					default:
					}
					c.observersMutex.Lock()
				}
				c.observersMutex.Unlock()
			} else {
				break
			}
		}
	}()
	return out
}

func (c *connectableObservable) Contains(equal Predicate) Single {
	return c.observable.Contains(equal)
}

func (c *connectableObservable) Count() Single {
	return c.observable.Count()
}

func (c *connectableObservable) DefaultIfEmpty(defaultValue interface{}) Observable {
	return c.observable.DefaultIfEmpty(defaultValue)
}

func (c *connectableObservable) Distinct(apply Function) Observable {
	return c.observable.Distinct(apply)
}

func (c *connectableObservable) DistinctUntilChanged(apply Function) Observable {
	return c.observable.DistinctUntilChanged(apply)
}

func (c *connectableObservable) DoOnEach(onNotification Consumer) Observable {
	return c.observable.DoOnEach(onNotification)
}

func (c *connectableObservable) ElementAt(index uint) Single {
	return c.observable.ElementAt(index)
}

func (c *connectableObservable) Filter(apply Predicate) Observable {
	return c.observable.Filter(apply)
}

func (c *connectableObservable) First() Observable {
	return c.observable.First()
}

func (c *connectableObservable) FirstOrDefault(defaultValue interface{}) Single {
	return c.observable.FirstOrDefault(defaultValue)
}

func (c *connectableObservable) FlatMap(apply func(interface{}) Observable, maxInParallel uint) Observable {
	return c.observable.FlatMap(apply, maxInParallel)
}

func (c *connectableObservable) ForEach(nextFunc NextFunc, errFunc ErrFunc,
	doneFunc DoneFunc, opts ...Option) Observer {
	return c.observable.ForEach(nextFunc, errFunc, doneFunc, opts...)
}

func (c *connectableObservable) IgnoreElements() Observable {
	return c.observable.IgnoreElements()
}

func (c *connectableObservable) Last() Observable {
	return c.observable.Last()
}

func (c *connectableObservable) LastOrDefault(defaultValue interface{}) Single {
	return c.observable.LastOrDefault(defaultValue)
}

func (c *connectableObservable) Map(apply Function) Observable {
	return c.observable.Map(apply)
}

func (c *connectableObservable) Max(comparator Comparator) OptionalSingle {
	return c.observable.Max(comparator)
}

func (c *connectableObservable) Min(comparator Comparator) OptionalSingle {
	return c.observable.Min(comparator)
}

func (c *connectableObservable) OnErrorResumeNext(resumeSequence ErrorToObservableFunction) Observable {
	return c.observable.OnErrorResumeNext(resumeSequence)
}

func (c *connectableObservable) OnErrorReturn(resumeFunc ErrorFunction) Observable {
	return c.observable.OnErrorReturn(resumeFunc)
}

func (c *connectableObservable) OnErrorReturnItem(item interface{}) Observable {
	return c.observable.OnErrorReturnItem(item)
}

func (c *connectableObservable) Publish() ConnectableObservable {
	return c.observable.Publish()
}

func (c *connectableObservable) Reduce(apply Function2) OptionalSingle {
	return c.observable.Reduce(apply)
}

func (c *connectableObservable) Repeat(count int64, frequency Duration) Observable {
	return c.observable.Repeat(count, frequency)
}

func (c *connectableObservable) Sample(obs Observable) Observable {
	return c.observable.Sample(obs)
}

func (c *connectableObservable) Scan(apply Function2) Observable {
	return c.observable.Scan(apply)
}

func (c *connectableObservable) SequenceEqual(obs Observable) Single {
	return c.observable.SequenceEqual(obs)
}

func (c *connectableObservable) Skip(nth uint) Observable {
	return c.observable.Skip(nth)
}

func (c *connectableObservable) SkipLast(nth uint) Observable {
	return c.observable.SkipLast(nth)
}

func (c *connectableObservable) SkipWhile(apply Predicate) Observable {
	return c.observable.SkipWhile(apply)
}

func (c *connectableObservable) StartWithItems(item interface{}, items ...interface{}) Observable {
	return c.observable.StartWithItems(item, items...)
}

func (c *connectableObservable) StartWithIterable(iterable Iterable) Observable {
	return c.observable.StartWithIterable(iterable)
}

func (c *connectableObservable) StartWithObservable(observable Observable) Observable {
	return c.observable.StartWithObservable(observable)
}

func (c *connectableObservable) Subscribe(handler EventHandler, opts ...Option) Observer {
	observableOptions := ParseOptions(opts...)

	ob := NewObserver(handler)
	var ch chan interface{}
	if observableOptions.BackpressureStrategy() == Buffer {
		ch = make(chan interface{}, observableOptions.Buffer())
	} else {
		ch = make(chan interface{})
	}
	ob.setItemChannel(ch)
	c.observersMutex.Lock()
	c.observers = append(c.observers, ob)
	c.observersMutex.Unlock()

	go func() {
		for item := range ch {
			switch item := item.(type) {
			case error:
				err := ob.OnError(item)
				if err != nil {
					panic(errors.Wrap(err, "error while sending error item from connectable observable"))
				}
				return
			default:
				err := ob.OnNext(item)
				if err != nil {
					panic(errors.Wrap(err, "error while sending next item from connectable observable"))
				}
			}
		}
	}()

	return ob
}

func (c *connectableObservable) SumFloat32() Single {
	return c.observable.SumFloat32()
}

func (c *connectableObservable) SumFloat64() Single {
	return c.observable.SumFloat64()
}

func (c *connectableObservable) SumInt64() Single {
	return c.observable.SumInt64()
}

func (c *connectableObservable) Take(nth uint) Observable {
	return c.observable.Take(nth)
}

func (c *connectableObservable) TakeLast(nth uint) Observable {
	return c.observable.TakeLast(nth)
}

func (c *connectableObservable) TakeUntil(apply Predicate) Observable {
	return c.observable.TakeUntil(apply)
}

func (c *connectableObservable) TakeWhile(apply Predicate) Observable {
	return c.observable.TakeWhile(apply)
}

func (c *connectableObservable) Timeout(ctx context.Context) Observable {
	return c.observable.Timeout(ctx)
}

func (c *connectableObservable) ToChannel(opts ...Option) Channel {
	return c.observable.ToChannel(opts...)
}

func (c *connectableObservable) ToMap(keySelector Function) Single {
	return c.observable.ToMap(keySelector)
}

func (c *connectableObservable) ToSlice() Single {
	return c.observable.ToSlice()
}

func (c *connectableObservable) ToMapWithValueSelector(keySelector, valueSelector Function) Single {
	return c.observable.ToMapWithValueSelector(keySelector, valueSelector)
}

func (c *connectableObservable) ZipFromObservable(publisher Observable, zipper Function2) Observable {
	return c.observable.ZipFromObservable(publisher, zipper)
}

func (c *connectableObservable) getCustomErrorStrategy() func(Observable, Observer, error) error {
	return c.observable.getCustomErrorStrategy()
}

func (c *connectableObservable) getNextStrategy() func(Observer, interface{}) error {
	return c.observable.getNextStrategy()
}
