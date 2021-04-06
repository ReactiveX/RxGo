// +build !all

package rxgo

import (
	"context"
	"errors"
	"fmt"
	"math/rand"
	"testing"
	"time"

	"go.uber.org/goleak"
)

const maxSleepNs = 10_000_000 // 10 ms

// TODO Keep enriching tests
func TestLeak(t *testing.T) {
	var (
		count  = 100
		fooErr = errors.New("")
	)

	observables := map[string]func(context.Context) Observable{
		"Amb": func(ctx context.Context) Observable {
			obs := FromChannel(make(chan Item), WithContext(ctx))
			return Amb([]Observable{obs}, WithContext(ctx))
		},
		"CombineLatest": func(ctx context.Context) Observable {
			return CombineLatest(func(i ...interface{}) interface{} {
				sum := 0
				for _, v := range i {
					if v == nil {
						continue
					}
					sum += v.(int)
				}
				return sum
			}, []Observable{
				Just(1, 2)(),
				Just(10, 11)(),
			})
		},
		"Concat": func(ctx context.Context) Observable {
			return Concat([]Observable{
				Just(1, 2, 3)(),
				Just(4, 5, 6)(),
			})
		},
		"FromChannel": func(ctx context.Context) Observable {
			return FromChannel(getChannel(ctx), WithContext(ctx))
		},
		"FromEventSource": func(ctx context.Context) Observable {
			return FromEventSource(getChannel(ctx), WithContext(ctx))
		},
	}

	actions := map[string]func(context.Context, Observable){
		"All": func(ctx context.Context, obs Observable) {
			obs.All(func(_ interface{}) bool {
				return true
			}, WithContext(ctx))
		},
		"Average": func(ctx context.Context, obs Observable) {
			obs.AverageInt(WithContext(ctx))
		},
		"BufferWithTime": func(ctx context.Context, obs Observable) {
			obs.BufferWithTime(WithDuration(time.Millisecond), WithContext(ctx))
		},
		"Connect": func(ctx context.Context, obs Observable) {
			obs.Connect(ctx)
		},
		"Contains": func(ctx context.Context, obs Observable) {
			obs.Contains(func(i interface{}) bool {
				return i == 2
			}, WithContext(ctx))
		},
		"For each": func(_ context.Context, obs Observable) {
			obs.ForEach(func(_ interface{}) {}, func(_ error) {}, func() {})
		},
	}

	defer goleak.VerifyNone(t)
	for testObservable, factory := range observables {
		for testAction, action := range actions {
			for i := 0; i < count; i++ {
				waitTime := randomTime()
				factory := factory
				action := action
				t.Run(fmt.Sprintf("%s - %s - %v - single", testObservable, testAction, waitTime), func(t *testing.T) {
					t.Parallel()
					ctx, cancel := context.WithTimeout(context.Background(), waitTime)
					defer cancel()
					action(ctx, factory(ctx))
				})
				t.Run(fmt.Sprintf("%s - %s - %v - composed", testObservable, testAction, waitTime), func(t *testing.T) {
					t.Parallel()
					ctx, cancel := context.WithTimeout(context.Background(), waitTime)
					defer cancel()
					action(ctx, factory(ctx).Map(func(_ context.Context, i interface{}) (interface{}, error) {
						return i, nil
					}))
				})
				t.Run(fmt.Sprintf("%s - %s - %v - erritem", testObservable, testAction, waitTime), func(t *testing.T) {
					t.Parallel()
					ctx, cancel := context.WithTimeout(context.Background(), waitTime)
					defer cancel()
					action(ctx, factory(ctx).Map(func(_ context.Context, i interface{}) (interface{}, error) {
						return nil, fooErr
					}))
				})
			}
			t.Run(fmt.Sprintf("%s - %s - already cancelled", testObservable, testAction), func(t *testing.T) {
				t.Parallel()
				ctx, cancel := context.WithCancel(context.Background())
				cancel()
				action(ctx, factory(ctx))
			})
		}
	}
}

func getChannel(ctx context.Context) chan Item {
	ch := make(chan Item, 3)
	go func() {
		time.Sleep(randomTime())
		Of(1).SendContext(ctx, ch)
		time.Sleep(randomTime())
		Of(2).SendContext(ctx, ch)
		time.Sleep(randomTime())
		Of(3).SendContext(ctx, ch)
	}()
	return ch
}

func randomTime() time.Duration {
	return time.Duration(rand.Intn(maxSleepNs)) * time.Nanosecond
}
