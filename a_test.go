package rxgo_test

import (
	"fmt"
	"github.com/reactivex/rxgo/v2"
	"testing"
	"time"
)

func TestMap(t *testing.T) {
	// Create a regular Observable
	ch := make(chan rxgo.Item)
	go func() {
		ch <- rxgo.Of(1)
		ch <- rxgo.Of(2)
		ch <- rxgo.Of(3)
		close(ch)
	}()
	observable := rxgo.FromChannel(ch, rxgo.WithPublishStrategy())

	// Create the first Observer
	observable.DoOnNext(func(i interface{}) {
		fmt.Printf("First observer: %d\n", i)
	})

	// Create the second Observer
	observable.DoOnNext(func(i interface{}) {
		fmt.Printf("Second observer: %d\n", i)
	})

	observable.Connect()

	time.Sleep(500 * time.Millisecond)

}
