package rxgo

import (
	"bufio"
	"context"
	"github.com/stretchr/testify/mock"
	"strconv"
	"strings"
)

type mockIterable struct {
	iterator Iterator
}

type mockIterator struct {
	mock.Mock
}

type task struct {
	observable int
	item       int
	close      bool
}

func (s *mockIterable) Iterator(ctx context.Context) Iterator {
	return s.iterator
}

func mockObservable(iterator Iterator) Observable {
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

func mockIterators(in string) ([]*mockIterator, error) {
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

	observables := make([]*mockIterator, 0, len(m))
	calls := make([]*mock.Call, len(m))
	for i := 0; i < len(m); i++ {
		observables = append(observables, new(mockIterator))
	}

	item, err := args(tasks[0])
	call := observables[0].On("Next", mock.Anything).Once().Return(item, err)
	calls[0] = call

	var lastCh chan struct{}
	lastObservableType := tasks[0].observable
	for i := 1; i < len(tasks); i++ {
		t := tasks[i]
		index := m[t.observable]
		obs := observables[index]
		item, err := args(t)
		if lastObservableType == t.observable {
			if calls[index] == nil {
				calls[index] = obs.On("Next", mock.Anything).Once().Return(item, err)
			} else {
				calls[index].On("Next", mock.Anything).Once().Return(item, err)
			}
		} else {
			lastObservableType = t.observable
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

func (m *mockIterator) Next(ctx context.Context) (interface{}, error) {
	args := m.Called(ctx)
	return args.Get(0), args.Error(1)
}
