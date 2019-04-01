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

func (connectableObservable) Iterator() Iterator {
	panic("implement me")
}

func (connectableObservable) All(predicate Predicate) Single {
	panic("implement me")
}

func (connectableObservable) AverageFloat32() Single {
	panic("implement me")
}

func (connectableObservable) AverageFloat64() Single {
	panic("implement me")
}

func (connectableObservable) AverageInt() Single {
	panic("implement me")
}

func (connectableObservable) AverageInt8() Single {
	panic("implement me")
}

func (connectableObservable) AverageInt16() Single {
	panic("implement me")
}

func (connectableObservable) AverageInt32() Single {
	panic("implement me")
}

func (connectableObservable) AverageInt64() Single {
	panic("implement me")
}

func (connectableObservable) BufferWithCount(count, skip int) Observable {
	panic("implement me")
}

func (connectableObservable) BufferWithTime(timespan, timeshift Duration) Observable {
	panic("implement me")
}

func (connectableObservable) BufferWithTimeOrCount(timespan Duration, count int) Observable {
	panic("implement me")
}

func (o *connectableObservable) Connect() Observer {
	out := NewObserver()
	go func() {
		it := o.observable.Iterator()
		for {
			if item, err := it.Next(); err == nil {
				o.observersMutex.Lock()
				for _, observer := range o.observers {
					o.observersMutex.Unlock()
					select {
					case observer.getChannel() <- item:
					default:
					}
					o.observersMutex.Lock()
				}
				o.observersMutex.Unlock()
			} else {
				break
			}
		}
	}()
	return out
}

func (connectableObservable) Contains(equal Predicate) Single {
	panic("implement me")
}

func (connectableObservable) Count() Single {
	panic("implement me")
}

func (connectableObservable) DefaultIfEmpty(defaultValue interface{}) Observable {
	panic("implement me")
}

func (connectableObservable) Distinct(apply Function) Observable {
	panic("implement me")
}

func (connectableObservable) DistinctUntilChanged(apply Function) Observable {
	panic("implement me")
}

func (connectableObservable) DoOnEach(onNotification Consumer) Observable {
	panic("implement me")
}

func (connectableObservable) ElementAt(index uint) Single {
	panic("implement me")
}

func (connectableObservable) Filter(apply Predicate) Observable {
	panic("implement me")
}

func (connectableObservable) First() Observable {
	panic("implement me")
}

func (connectableObservable) FirstOrDefault(defaultValue interface{}) Single {
	panic("implement me")
}

func (connectableObservable) FlatMap(apply func(interface{}) Observable, maxInParallel uint) Observable {
	panic("implement me")
}

func (connectableObservable) ForEach(nextFunc handlers.NextFunc, errFunc handlers.ErrFunc,
	doneFunc handlers.DoneFunc, opts ...options.Option) Observer {
	panic("implement me")
}

func (connectableObservable) Last() Observable {
	panic("implement me")
}

func (connectableObservable) LastOrDefault(defaultValue interface{}) Single {
	panic("implement me")
}

func (connectableObservable) Map(apply Function) Observable {
	panic("implement me")
}

func (connectableObservable) Max(comparator Comparator) OptionalSingle {
	panic("implement me")
}

func (connectableObservable) Min(comparator Comparator) OptionalSingle {
	panic("implement me")
}

func (connectableObservable) OnErrorResumeNext(resumeSequence ErrorToObservableFunction) Observable {
	panic("implement me")
}

func (connectableObservable) OnErrorReturn(resumeFunc ErrorFunction) Observable {
	panic("implement me")
}

func (connectableObservable) Publish() ConnectableObservable {
	panic("implement me")
}

func (connectableObservable) Reduce(apply Function2) OptionalSingle {
	panic("implement me")
}

func (connectableObservable) Repeat(count int64, frequency Duration) Observable {
	panic("implement me")
}

func (connectableObservable) Scan(apply Function2) Observable {
	panic("implement me")
}

func (connectableObservable) SequenceEqual(obs Observable) Single {
	panic("implement me")
}

func (connectableObservable) Skip(nth uint) Observable {
	panic("implement me")
}

func (connectableObservable) SkipLast(nth uint) Observable {
	panic("implement me")
}

func (connectableObservable) SkipWhile(apply Predicate) Observable {
	panic("implement me")
}

func (connectableObservable) StartWithItems(items ...interface{}) Observable {
	panic("implement me")
}

func (connectableObservable) StartWithIterable(iterable Iterable) Observable {
	panic("implement me")
}

func (connectableObservable) StartWithObservable(observable Observable) Observable {
	panic("implement me")
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

func (connectableObservable) SumFloat32() Single {
	panic("implement me")
}

func (connectableObservable) SumFloat64() Single {
	panic("implement me")
}

func (connectableObservable) SumInt64() Single {
	panic("implement me")
}

func (connectableObservable) Take(nth uint) Observable {
	panic("implement me")
}

func (connectableObservable) TakeLast(nth uint) Observable {
	panic("implement me")
}

func (connectableObservable) TakeWhile(apply Predicate) Observable {
	panic("implement me")
}

func (connectableObservable) ToList() Observable {
	panic("implement me")
}

func (connectableObservable) ToMap(keySelector Function) Observable {
	panic("implement me")
}

func (connectableObservable) ToMapWithValueSelector(keySelector Function, valueSelector Function) Observable {
	panic("implement me")
}

func (connectableObservable) ZipFromObservable(publisher Observable, zipper Function2) Observable {
	panic("implement me")
}

func (connectableObservable) getOnErrorResumeNext() ErrorToObservableFunction {
	panic("implement me")
}

func (connectableObservable) getOnErrorReturn() ErrorFunction {
	panic("implement me")
}
