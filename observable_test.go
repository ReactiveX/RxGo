package rxgo

import (
	"fmt"
	"testing"

	"github.com/stretchr/testify/require"
	"go.uber.org/goleak"
)

func TestMain(m *testing.M) {
	goleak.VerifyTestMain(m)
}

func TestObservable(t *testing.T) {
	// obs := &observableWrapper[string]{}
	// obs.subscribeOn(func(s string) {}, func(err error) {}, func() {}, func() {})
}

func TestNever(t *testing.T) {
	// ctx, cancel := context.WithTimeout(context.Background(), time.Second)
	// defer cancel()
	// NEVER[any]().SubscribeSync(func(a any) {}, func(err error) {}, func() {
	// 	log.Println("Completed")
	// })
}

func TestEmpty(t *testing.T) {
	checkObservableResults(t, EMPTY[any](), []any{}, nil, true)
}

func TestThrownError(t *testing.T) {
	var v = fmt.Errorf("uncaught error")
	checkObservableResults(t, ThrownError[string](func() error {
		return v
	}), []string{}, v, true)
}

func TestDefer(t *testing.T) {
	values := []string{"a", "b", "c"}
	obs := Defer(func() IObservable[string] {
		return newObservable(func(subscriber Subscriber[string]) {
			for _, v := range values {
				subscriber.Next(v)
			}
			subscriber.Complete()
		})
	})
	checkObservableResults(t, obs, values, nil, true)
}

func TestRange(t *testing.T) {
	t.Run("Range from 1 to 10", func(t *testing.T) {
		checkObservableResults(t, Range[uint](1, 10), []uint{1, 2, 3, 4, 5, 6, 7, 8, 9, 10}, nil, true)
	})
	t.Run("Range from 0 to 3", func(t *testing.T) {
		checkObservableResults(t, Range[uint](0, 3), []uint{0, 1, 2}, nil, true)
	})
}

func TestInterval(t *testing.T) {
}

func TestScheduled(t *testing.T) {

}

func TestTimer(t *testing.T) {

}

func checkObservableResult[T any](t *testing.T, obs IObservable[T], result T, err error, isCompleted bool) {
	var (
		hasCompleted  bool
		collectedErr  error
		collectedData T
	)
	obs.SubscribeSync(func(v T) {
		collectedData = v
	}, func(err error) {
		collectedErr = err
	}, func() {
		hasCompleted = true
	})
	require.Equal(t, collectedData, result)
	require.Equal(t, hasCompleted, isCompleted)
	require.Equal(t, collectedErr, err)
}

func checkObservableResults[T any](t *testing.T, obs IObservable[T], result []T, err error, isCompleted bool) {
	var (
		hasCompleted  bool
		collectedErr  error
		collectedData = make([]T, 0, len(result))
	)
	obs.SubscribeSync(func(v T) {
		collectedData = append(collectedData, v)
	}, func(err error) {
		collectedErr = err
	}, func() {
		hasCompleted = true
	})
	require.ElementsMatch(t, collectedData, result)
	require.Equal(t, hasCompleted, isCompleted)
	require.Equal(t, collectedErr, err)
}
