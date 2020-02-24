package rxgo

import (
	"context"
	"fmt"
	"github.com/stretchr/testify/assert"
	"golang.org/x/sync/errgroup"
	"reflect"
	"sync"
	"testing"
	"time"
)

func testConnectableSingle(t *testing.T, obs Observable) {
	ctx, cancel := context.WithTimeout(context.Background(), time.Second)
	defer cancel()
	eg, _ := errgroup.WithContext(ctx)

	expected := []interface{}{1, 2, 3}

	nbConsumers := 3
	wg := sync.WaitGroup{}
	wg.Add(nbConsumers)
	// Before Connect() is called we create multiple observers
	// We check all observers receive the same items
	for i := 0; i < nbConsumers; i++ {
		eg.Go(func() error {
			observer := obs.Observe(WithContext(ctx))
			wg.Done()
			got, err := collect(ctx, observer)
			if err != nil {
				return err
			}
			if !reflect.DeepEqual(got, expected) {
				return fmt.Errorf("expected: %v, got: %v", expected, got)
			}
			return nil
		})
	}

	wg.Wait()
	obs.Connect()
	assert.NoError(t, eg.Wait())
}

func testConnectableComposed(t *testing.T, obs Observable) {
	obs = obs.Map(func(_ context.Context, i interface{}) (interface{}, error) {
		return i.(int) + 1, nil
	}, WithPublishStrategy())

	ctx, cancel := context.WithTimeout(context.Background(), time.Second)
	defer cancel()
	eg, _ := errgroup.WithContext(ctx)

	expected := []interface{}{2, 3, 4}

	nbConsumers := 3
	wg := sync.WaitGroup{}
	wg.Add(nbConsumers)
	// Before Connect() is called we create multiple observers
	// We check all observers receive the same items
	for i := 0; i < nbConsumers; i++ {
		eg.Go(func() error {
			observer := obs.Observe(WithContext(ctx))
			wg.Done()

			got, err := collect(ctx, observer)
			if err != nil {
				return err
			}
			if !reflect.DeepEqual(got, expected) {
				return fmt.Errorf("expected: %v, got: %v", expected, got)
			}
			return nil
		})
	}

	wg.Wait()
	obs.Connect()
	assert.NoError(t, eg.Wait())
}

func testConnectableWithoutConnect(t *testing.T, obs Observable) {
	ctx, cancel := context.WithTimeout(context.Background(), 50*time.Millisecond)
	defer cancel()
	Assert(ctx, t, obs, IsEmpty())
}

func Test_Connectable_IterableChannel_Single(t *testing.T) {
	ch := make(chan Item, 10)
	go func() {
		ch <- Of(1)
		ch <- Of(2)
		ch <- Of(3)
		close(ch)
	}()
	testConnectableSingle(t, FromChannel(ch, WithPublishStrategy()))
}

func Test_Connectable_IterableChannel_Composed(t *testing.T) {
	ch := make(chan Item, 10)
	go func() {
		ch <- Of(1)
		ch <- Of(2)
		ch <- Of(3)
		close(ch)
	}()
	testConnectableComposed(t, FromChannel(ch, WithPublishStrategy()))
}

func Test_Connectable_IterableChannel_WithoutConnect(t *testing.T) {
	ch := make(chan Item, 10)
	go func() {
		ch <- Of(1)
		ch <- Of(2)
		ch <- Of(3)
		close(ch)
	}()
	obs := FromChannel(ch, WithPublishStrategy(), WithBufferedChannel(10))
	testConnectableWithoutConnect(t, obs)
}
