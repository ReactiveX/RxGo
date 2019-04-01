package rxgo

import (
	"github.com/reactivex/rxgo/handlers"
	"github.com/reactivex/rxgo/options"
	"sync"
)

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

func (c *connectableObservable) Iterator() Iterator {
	return c.Iterator()
}

func (c *connectableObservable) All(predicate Predicate) Single {
	return c.All(predicate)
}

func (c *connectableObservable) AverageFloat32() Single {
	return c.AverageFloat32()
}

func (c *connectableObservable) AverageFloat64() Single {
	return c.AverageFloat64()
}

func (c *connectableObservable) AverageInt() Single {
	return c.AverageInt()
}

func (c *connectableObservable) AverageInt8() Single {
	return c.AverageInt8()
}

func (c *connectableObservable) AverageInt16() Single {
	return c.AverageInt16()
}

func (c *connectableObservable) AverageInt32() Single {
	return c.AverageInt32()
}

func (c *connectableObservable) AverageInt64() Single {
	return c.AverageInt64()
}

func (c *connectableObservable) BufferWithCount(count, skip int) Observable {
	return c.BufferWithCount(count, skip)
}

func (c *connectableObservable) BufferWithTime(timespan, timeshift Duration) Observable {
	return c.BufferWithTime(timespan, timeshift)
}

func (c *connectableObservable) BufferWithTimeOrCount(timespan Duration, count int) Observable {
	return c.BufferWithTimeOrCount(timespan, count)
}

func (c *connectableObservable) Connect() Observer {
	out := NewObserver()
	go func() {
		it := c.observable.Iterator()
		for {
			if item, err := it.Next(); err == nil {
				c.observersMutex.Lock()
				for _, observer := range c.observers {
					c.observersMutex.Unlock()
					select {
					case observer.getChannel() <- item:
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
	return c.Contains(equal)
}

func (c *connectableObservable) Count() Single {
	return c.Count()
}

func (c *connectableObservable) DefaultIfEmpty(defaultValue interface{}) Observable {
	return c.DefaultIfEmpty(defaultValue)
}

func (c *connectableObservable) Distinct(apply Function) Observable {
	return c.Distinct(apply)
}

func (c *connectableObservable) DistinctUntilChanged(apply Function) Observable {
	return c.DistinctUntilChanged(apply)
}

func (c *connectableObservable) DoOnEach(onNotification Consumer) Observable {
	return c.DoOnEach(onNotification)
}

func (c *connectableObservable) ElementAt(index uint) Single {
	return c.ElementAt(index)
}

func (c *connectableObservable) Filter(apply Predicate) Observable {
	return c.Filter(apply)
}

func (c *connectableObservable) First() Observable {
	return c.First()
}

func (c *connectableObservable) FirstOrDefault(defaultValue interface{}) Single {
	return c.FirstOrDefault(defaultValue)
}

func (c *connectableObservable) FlatMap(apply func(interface{}) Observable, maxInParallel uint) Observable {
	return c.FlatMap(apply, maxInParallel)
}

func (c *connectableObservable) ForEach(nextFunc handlers.NextFunc, errFunc handlers.ErrFunc,
	doneFunc handlers.DoneFunc, opts ...options.Option) Observer {
	return c.ForEach(nextFunc, errFunc, doneFunc, opts...)
}

func (c *connectableObservable) IgnoreElements() Observable {
	return c.IgnoreElements()
}

func (c *connectableObservable) Last() Observable {
	return c.Last()
}

func (c *connectableObservable) LastOrDefault(defaultValue interface{}) Single {
	return c.LastOrDefault(defaultValue)
}

func (c *connectableObservable) Map(apply Function) Observable {
	return c.Map(apply)
}

func (c *connectableObservable) Max(comparator Comparator) OptionalSingle {
	return c.Max(comparator)
}

func (c *connectableObservable) Min(comparator Comparator) OptionalSingle {
	return c.Min(comparator)
}

func (c *connectableObservable) OnErrorResumeNext(resumeSequence ErrorToObservableFunction) Observable {
	return c.OnErrorResumeNext(resumeSequence)
}

func (c *connectableObservable) OnErrorReturn(resumeFunc ErrorFunction) Observable {
	return c.OnErrorReturn(resumeFunc)
}

func (c *connectableObservable) OnErrorReturnItem(item interface{}) Observable {
	return c.OnErrorReturnItem(item)
}

func (c *connectableObservable) Publish() ConnectableObservable {
	return c.Publish()
}

func (c *connectableObservable) Reduce(apply Function2) OptionalSingle {
	return c.Reduce(apply)
}

func (c *connectableObservable) Repeat(count int64, frequency Duration) Observable {
	return c.Repeat(count, frequency)
}

func (c *connectableObservable) Scan(apply Function2) Observable {
	return c.Scan(apply)
}

func (c *connectableObservable) SequenceEqual(obs Observable) Single {
	return c.SequenceEqual(obs)
}

func (c *connectableObservable) Skip(nth uint) Observable {
	return c.Skip(nth)
}

func (c *connectableObservable) SkipLast(nth uint) Observable {
	return c.SkipLast(nth)
}

func (c *connectableObservable) SkipWhile(apply Predicate) Observable {
	return c.SkipWhile(apply)
}

func (c *connectableObservable) StartWithItems(items ...interface{}) Observable {
	return c.StartWithItems(items...)
}

func (c *connectableObservable) StartWithIterable(iterable Iterable) Observable {
	return c.StartWithIterable(iterable)
}

func (c *connectableObservable) StartWithObservable(observable Observable) Observable {
	return c.StartWithObservable(observable)
}

func (o *connectableObservable) Subscribe(handler handlers.EventHandler, opts ...options.Option) Observer {
	observableOptions := options.ParseOptions(opts...)

	ob := CheckEventHandler(handler)
	ob.setBackpressureStrategy(observableOptions.BackpressureStrategy())
	var ch chan interface{}
	if observableOptions.BackpressureStrategy() == options.Buffer {
		ch = make(chan interface{}, observableOptions.Buffer())
	} else {
		ch = make(chan interface{})
	}
	ob.setChannel(ch)
	o.observersMutex.Lock()
	o.observers = append(o.observers, ob)
	o.observersMutex.Unlock()

	go func() {
		for item := range ch {
			switch item := item.(type) {
			case error:
				ob.OnError(item)
				return
			default:
				ob.OnNext(item)
			}
		}
	}()

	return ob
}

func (c *connectableObservable) SumFloat32() Single {
	return c.SumFloat32()
}

func (c *connectableObservable) SumFloat64() Single {
	return c.SumFloat64()
}

func (c *connectableObservable) SumInt64() Single {
	return c.SumInt64()
}

func (c *connectableObservable) Take(nth uint) Observable {
	return c.Take(nth)
}

func (c *connectableObservable) TakeLast(nth uint) Observable {
	return c.TakeLast(nth)
}

func (c *connectableObservable) TakeWhile(apply Predicate) Observable {
	return c.TakeWhile(apply)
}

func (c *connectableObservable) ToList() Observable {
	return c.ToList()
}

func (c *connectableObservable) ToMap(keySelector Function) Observable {
	return c.ToMap(keySelector)
}

func (c *connectableObservable) ToMapWithValueSelector(keySelector Function, valueSelector Function) Observable {
	return c.ToMapWithValueSelector(keySelector, valueSelector)
}

func (c *connectableObservable) ZipFromObservable(publisher Observable, zipper Function2) Observable {
	return c.ZipFromObservable(publisher, zipper)
}

func (c *connectableObservable) getIgnoreElements() bool {
	return c.getIgnoreElements()
}

func (c *connectableObservable) getOnErrorResumeNext() ErrorToObservableFunction {
	return c.getOnErrorResumeNext()
}

func (c *connectableObservable) getOnErrorReturn() ErrorFunction {
	return c.getOnErrorReturn()
}

func (c *connectableObservable) getOnErrorReturnItem() interface{} {
	return c.getOnErrorReturnItem()
}
