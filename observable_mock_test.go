package rxgo

import (
	"bufio"
	"context"
	"github.com/reactivex/rxgo/handlers"
	"github.com/reactivex/rxgo/options"
	"github.com/stretchr/testify/mock"
	"strconv"
	"strings"
	"testing"
)

func newSyncObservable(iterator Iterator) Observable {
	return &observable{
		observableType: cold,
		iterable: &syncIterable{
			iterator: iterator,
		},
	}
}

type syncIterable struct {
	iterator Iterator
}

func (s *syncIterable) Iterator(ctx context.Context) Iterator {
	return s.iterator
}

type MockObservable struct {
	mock.Mock
}

type MockIterator struct {
	mock.Mock
}

type task struct {
	observable int
	item       int
	close      bool
}

func TestMockIterators(t *testing.T) {
	mockIterators("1:1,1:2,2:0,1:3,1:4,2:0,1:5,1:6,2:0,1:close,2:close")
}

func countTab(line string) int {
	i := 0
	for _, runeValue := range line {
		if runeValue == '\t' {
			i++
		} else {
			break
		}
	}
	return i
}

func mockIterators(in string) ([]*MockIterator, error) {
	scanner := bufio.NewScanner(strings.NewReader(in))
	m := make(map[int]int)
	tasks := make([]task, 0)
	count := 0
	for scanner.Scan() {
		s := scanner.Text()
		if s == "" {
			continue
		}
		observable := countTab(s)
		v := strings.TrimSpace(s)
		if v == "x" {
			tasks = append(tasks, task{
				observable: observable,
				close:      true,
			})
		} else {
			n, err := strconv.Atoi(v)
			if err != nil {
				return nil, err
			}
			tasks = append(tasks, task{
				observable: observable,
				item:       n,
			})
		}
		if _, contains := m[observable]; !contains {
			m[observable] = count
			count++
		}
	}

	observables := make([]*MockIterator, 0, len(m))
	calls := make([]*mock.Call, len(m))
	for i := 0; i < len(m); i++ {
		observables = append(observables, new(MockIterator))
	}

	item, err := args(tasks[0])
	call := observables[0].On("Next", mock.Anything).Once().Return(item, err)
	calls[0] = call

	var lastCh chan struct{}
	lastObs := tasks[0].observable
	for i := 1; i < len(tasks); i++ {
		t := tasks[i]
		index := m[t.observable]
		obs := observables[index]
		item, err := args(t)
		if lastObs == t.observable {
			if calls[index] == nil {
				calls[index] = obs.On("Next", mock.Anything).Once().Return(item, err)
			} else {
				calls[index].On("Next", mock.Anything).Once().Return(item, err)
			}
		} else {
			lastObs = t.observable
			if lastCh == nil {
				ch := make(chan struct{})
				lastCh = ch
				if calls[index] == nil {
					calls[index] = obs.On("Next", mock.Anything).Once().Return(item, err).
						Run(func(args mock.Arguments) {
							run(ch, nil)
						})
				} else {
					calls[index].On("Next", mock.Anything).Once().Return(item, err).
						Run(func(args mock.Arguments) {
							run(ch, nil)
						})
				}
			} else {
				ch := make(chan struct{})
				previous := lastCh
				if calls[index] == nil {
					calls[index] = obs.On("Next", mock.Anything).Once().Return(item, err).
						Run(func(args mock.Arguments) {
							run(ch, previous)
						})
				} else {
					calls[index].On("Next", mock.Anything).Once().Return(item, err).
						Run(func(args mock.Arguments) {
							run(ch, previous)
						})
				}
				lastCh = ch
			}
		}
	}
	return observables, nil
}

func args(t task) (interface{}, error) {
	if t.close {
		return nil, &NoSuchElementError{}
	}
	return t.item, nil
}

func run(wait chan struct{}, send chan struct{}) {
	if send != nil {
		send <- struct{}{}
	}
	if wait != nil {
		<-wait
	}
}

func (m *MockIterator) Next(ctx context.Context) (interface{}, error) {
	args := m.Called(ctx)
	return args.Get(0), args.Error(1)
}

func (m *MockObservable) Iterator(ctx context.Context) Iterator {
	args := m.Called(ctx)
	return args.Get(0).(Iterator)
}

func (m *MockObservable) All(predicate Predicate) Single {
	args := m.Called(predicate)
	return args.Get(0).(Single)
}

func (m *MockObservable) AverageFloat32() Single {
	panic("implement me")
}

func (m *MockObservable) AverageFloat64() Single {
	panic("implement me")
}

func (m *MockObservable) AverageInt() Single {
	panic("implement me")
}

func (m *MockObservable) AverageInt8() Single {
	panic("implement me")
}

func (m *MockObservable) AverageInt16() Single {
	panic("implement me")
}

func (m *MockObservable) AverageInt32() Single {
	panic("implement me")
}

func (m *MockObservable) AverageInt64() Single {
	panic("implement me")
}

func (m *MockObservable) BufferWithCount(count, skip int) Observable {
	panic("implement me")
}

func (m *MockObservable) BufferWithTime(timespan, timeshift Duration) Observable {
	panic("implement me")
}

func (m *MockObservable) BufferWithTimeOrCount(timespan Duration, count int) Observable {
	panic("implement me")
}

func (m *MockObservable) Contains(equal Predicate) Single {
	panic("implement me")
}

func (m *MockObservable) Count() Single {
	panic("implement me")
}

func (m *MockObservable) DefaultIfEmpty(defaultValue interface{}) Observable {
	panic("implement me")
}

func (m *MockObservable) Distinct(apply Function) Observable {
	panic("implement me")
}

func (m *MockObservable) DistinctUntilChanged(apply Function) Observable {
	panic("implement me")
}

func (m *MockObservable) DoOnEach(onNotification Consumer) Observable {
	panic("implement me")
}

func (m *MockObservable) ElementAt(index uint) Single {
	panic("implement me")
}

func (m *MockObservable) Filter(apply Predicate) Observable {
	panic("implement me")
}

func (m *MockObservable) First() Observable {
	panic("implement me")
}

func (m *MockObservable) FirstOrDefault(defaultValue interface{}) Single {
	panic("implement me")
}

func (m *MockObservable) FlatMap(apply func(interface{}) Observable, maxInParallel uint) Observable {
	panic("implement me")
}

func (m *MockObservable) ForEach(nextFunc handlers.NextFunc, errFunc handlers.ErrFunc,
	doneFunc handlers.DoneFunc, opts ...options.Option) Observer {
	panic("implement me")
}

func (m *MockObservable) IgnoreElements() Observable {
	panic("implement me")
}

func (m *MockObservable) Last() Observable {
	panic("implement me")
}

func (m *MockObservable) LastOrDefault(defaultValue interface{}) Single {
	panic("implement me")
}

func (m *MockObservable) Map(apply Function) Observable {
	panic("implement me")
}

func (m *MockObservable) Max(comparator Comparator) OptionalSingle {
	panic("implement me")
}

func (m *MockObservable) Min(comparator Comparator) OptionalSingle {
	panic("implement me")
}

func (m *MockObservable) OnErrorResumeNext(resumeSequence ErrorToObservableFunction) Observable {
	panic("implement me")
}

func (m *MockObservable) OnErrorReturn(resumeFunc ErrorFunction) Observable {
	panic("implement me")
}

func (m *MockObservable) OnErrorReturnItem(item interface{}) Observable {
	panic("implement me")
}

func (m *MockObservable) Publish() ConnectableObservable {
	panic("implement me")
}

func (m *MockObservable) Reduce(apply Function2) OptionalSingle {
	panic("implement me")
}

func (m *MockObservable) Repeat(count int64, frequency Duration) Observable {
	panic("implement me")
}

func (m *MockObservable) Sample(obs Observable) Observable {
	panic("implement me")
}

func (m *MockObservable) Scan(apply Function2) Observable {
	panic("implement me")
}

func (m *MockObservable) SequenceEqual(obs Observable) Single {
	panic("implement me")
}

func (m *MockObservable) Skip(nth uint) Observable {
	panic("implement me")
}

func (m *MockObservable) SkipLast(nth uint) Observable {
	panic("implement me")
}

func (m *MockObservable) SkipWhile(apply Predicate) Observable {
	panic("implement me")
}

func (m *MockObservable) StartWithItems(item interface{}, items ...interface{}) Observable {
	panic("implement me")
}

func (m *MockObservable) StartWithIterable(iterable Iterable) Observable {
	panic("implement me")
}

func (m *MockObservable) StartWithObservable(observable Observable) Observable {
	panic("implement me")
}

func (m *MockObservable) Subscribe(handler handlers.EventHandler, opts ...options.Option) Observer {
	panic("implement me")
}

func (m *MockObservable) SumFloat32() Single {
	panic("implement me")
}

func (m *MockObservable) SumFloat64() Single {
	panic("implement me")
}

func (m *MockObservable) SumInt64() Single {
	panic("implement me")
}

func (m *MockObservable) Take(nth uint) Observable {
	panic("implement me")
}

func (m *MockObservable) TakeLast(nth uint) Observable {
	panic("implement me")
}

func (m *MockObservable) TakeUntil(apply Predicate) Observable {
	panic("implement me")
}

func (m *MockObservable) TakeWhile(apply Predicate) Observable {
	panic("implement me")
}

func (m *MockObservable) Timeout(duration Duration) Observable {
	panic("implement me")
}

func (m *MockObservable) ToChannel(opts ...options.Option) Channel {
	panic("implement me")
}

func (m *MockObservable) ToMap(keySelector Function) Single {
	panic("implement me")
}

func (m *MockObservable) ToMapWithValueSelector(keySelector, valueSelector Function) Single {
	panic("implement me")
}

func (m *MockObservable) ToSlice() Single {
	panic("implement me")
}

func (m *MockObservable) ZipFromObservable(publisher Observable, zipper Function2) Observable {
	panic("implement me")
}

func (m *MockObservable) getIgnoreElements() bool {
	panic("implement me")
}

func (m *MockObservable) getOnErrorResumeNext() ErrorToObservableFunction {
	panic("implement me")
}

func (m *MockObservable) getOnErrorReturn() ErrorFunction {
	panic("implement me")
}

func (m *MockObservable) getOnErrorReturnItem() interface{} {
	panic("implement me")
}
