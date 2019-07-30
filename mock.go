package rxgo

import (
	"bufio"
	"context"
	"strconv"
	"strings"
	"time"

	"github.com/pkg/errors"
	"github.com/stretchr/testify/mock"
)

const signalCh = byte(0)

var mockError = errors.New("")

type mockContext struct {
	mock.Mock
}

type mockIterable struct {
	iterator Iterator
}

type mockIterator struct {
	mock.Mock
}

type task struct {
	index   int
	item    int
	error   error
	close   bool
	context bool
}

type mockType struct {
	observable bool
	context    bool
	index      int
}

func (m *mockContext) Deadline() (deadline time.Time, ok bool) {
	panic("implement me")
}

func (m *mockContext) Done() <-chan struct{} {
	outputs := m.Called()
	return outputs.Get(0).(chan struct{})
}

func (m *mockContext) Err() error {
	panic("implement me")
}

func (m *mockContext) Value(key interface{}) interface{} {
	panic("implement me")
}

func (s *mockIterable) Iterator(ctx context.Context) Iterator {
	return s.iterator
}

func newMockObservable(iterator Iterator) Observable {
	return &observable{
		observableType: cold,
		iterable: &mockIterable{
			iterator: iterator,
		},
	}
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

// TODO Causality with more than two observables
func causality(in string) ([]Observable, []context.Context) {
	scanner := bufio.NewScanner(strings.NewReader(in))
	types := make([]mockType, 0)
	tasks := make([]task, 0)
	// Search header
	countObservables := 0
	countContexts := 0
	for scanner.Scan() {
		s := scanner.Text()
		if s == "" {
			continue
		}
		splits := strings.Split(s, "\t")
		for _, split := range splits {
			switch split {
			case "o":
				types = append(types, mockType{
					observable: true,
					index:      countObservables,
				})
				countObservables++
			case "c":
				types = append(types, mockType{
					context: true,
					index:   countContexts,
				})
				countContexts++
			default:
				panic("unknown type: " + split)
			}
		}
		break
	}
	for scanner.Scan() {
		s := scanner.Text()
		if s == "" {
			continue
		}
		index := countTab(s)
		v := strings.TrimSpace(s)
		if types[index].observable {
			switch v {
			case "x":
				tasks = append(tasks, task{
					index: types[index].index,
					close: true,
				})
			case "e":
				tasks = append(tasks, task{
					index: types[index].index,
					error: mockError,
				})
			default:
				n, err := strconv.Atoi(v)
				if err != nil {
					panic(err)
				}
				tasks = append(tasks, task{
					index: types[index].index,
					item:  n,
				})
			}
		} else if types[index].context {
			tasks = append(tasks, task{
				index:   types[index].index,
				context: true,
			})
		} else {
			panic("unknwon type")
		}
	}

	iterators := make([]*mockIterator, 0, countObservables)
	for i := 0; i < countObservables; i++ {
		iterators = append(iterators, new(mockIterator))
	}
	iteratorCalls := make([]*mock.Call, countObservables)

	contexts := make([]*mockContext, 0, countContexts)
	for i := 0; i < countContexts; i++ {
		contexts = append(contexts, new(mockContext))
	}
	iteratorContexts := make([]*mock.Call, countContexts)

	lastObservableType := -1
	if tasks[0].context {
		notif := make(chan struct{}, 1)
		call := contexts[0].On("Done").Once().Return(notif)
		iteratorContexts[0] = call
	} else {
		item, err := args(tasks[0])
		call := iterators[0].On("Next", mock.Anything).Once().Return(item, err)
		iteratorCalls[0] = call
		lastObservableType = tasks[0].index
	}

	var lastCh chan struct{}
	for i := 1; i < len(tasks); i++ {
		t := tasks[i]
		index := t.index

		if t.context {
			ctx := contexts[index]
			notif := make(chan struct{}, 1)
			lastObservableType = -1
			if lastCh == nil {
				ch := make(chan struct{}, 1)
				lastCh = ch
				if iteratorContexts[index] == nil {
					iteratorContexts[index] = ctx.On("Done").Once().Return(notif).
						Run(func(args mock.Arguments) {
							preDone(notif, ch, nil)
						})
				} else {
					iteratorContexts[index].On("Done").Once().Return(notif).
						Run(func(args mock.Arguments) {
							preDone(notif, ch, nil)
						})
				}
			} else {
				var ch chan struct{}
				// If this is the latest task we do not set any wait channel
				if i != len(tasks)-1 {
					ch = make(chan struct{}, 1)
				}
				previous := lastCh
				if iteratorContexts[index] == nil {
					iteratorContexts[index] = ctx.On("Done").Once().Return(notif).
						Run(func(args mock.Arguments) {
							preDone(notif, ch, previous)
						})
				} else {
					iteratorContexts[index].On("Done").Once().Return(notif).
						Run(func(args mock.Arguments) {
							preDone(notif, ch, previous)
						})
				}
				lastCh = ch
			}
		} else {
			obs := iterators[index]
			item, err := args(t)
			if lastObservableType == t.index {
				if iteratorCalls[index] == nil {
					iteratorCalls[index] = obs.On("Next", mock.Anything).Once().Return(item, err)
				} else {
					iteratorCalls[index].On("Next", mock.Anything).Once().Return(item, err)
				}
			} else {
				lastObservableType = t.index
				if lastCh == nil {
					ch := make(chan struct{})
					lastCh = ch
					if iteratorCalls[index] == nil {
						iteratorCalls[index] = obs.On("Next", mock.Anything).Once().Return(item, err).
							Run(func(args mock.Arguments) {
								preNext(args, ch, nil)
							})
					} else {
						iteratorCalls[index].On("Next", mock.Anything).Once().Return(item, err).
							Run(func(args mock.Arguments) {
								preNext(args, ch, nil)
							})
					}
				} else {
					var ch chan struct{}
					// If this is the latest task we do not set any wait channel
					if i != len(tasks)-1 {
						ch = make(chan struct{})
					}
					previous := lastCh
					if iteratorCalls[index] == nil {
						iteratorCalls[index] = obs.On("Next", mock.Anything).Once().Return(item, err).
							Run(func(args mock.Arguments) {
								preNext(args, ch, previous)
							})
					} else {
						iteratorCalls[index].On("Next", mock.Anything).Once().Return(item, err).
							Run(func(args mock.Arguments) {
								preNext(args, ch, previous)
							})
					}
					lastCh = ch
				}
			}
		}

	}

	observables := make([]Observable, 0, len(iterators))
	for _, iterator := range iterators {
		observables = append(observables, newMockObservable(iterator))
	}
	ctxs := make([]context.Context, 0, len(iterators))
	for _, c := range contexts {
		ctxs = append(ctxs, c)
	}
	return observables, ctxs
}

func args(t task) (interface{}, error) {
	if t.close {
		return nil, &NoSuchElementError{}
	}
	if t.error != nil {
		return t.error, nil
	}
	return t.item, nil
}

func preDone(notif chan struct{}, wait chan struct{}, send chan struct{}) {
	if send != nil {
		send <- struct{}{}
	}
	if wait != nil {
		<-wait
	}
	notif <- struct{}{}
}

func preNext(args mock.Arguments, wait chan struct{}, send chan struct{}) {
	if send != nil {
		send <- struct{}{}
	}
	if wait == nil {
		return
	}
	if len(args) == 1 {
		if ctx, ok := args[0].(context.Context); ok {
			if sig, ok := ctx.Value(signalCh).(chan struct{}); ok {
				select {
				case <-wait:
				case <-ctx.Done():
					sig <- struct{}{}
				}
				return
			}
		}
	}
	<-wait
}

func (m *mockIterator) Next(ctx context.Context) (interface{}, error) {
	sig := make(chan struct{}, 1)
	defer close(sig)
	outputs := m.Called(context.WithValue(ctx, signalCh, sig))
	select {
	case <-sig:
		return nil, &CancelledIteratorError{}
	default:
		return outputs.Get(0), outputs.Error(1)
	}
}
