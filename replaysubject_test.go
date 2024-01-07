package rxgo

import (
	"fmt"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
)

// TestReplaySubject verifies that a new subscriber receives the entire history
func TestReplaySubject(t *testing.T) {
	subject := NewReplaySubject(10)

	// load buffer
	for i := 0; i < 3; i++ {
		subject.Next(i)
	}

	_, obs := subject.Subscribe()

	values := make([]int, 0)
	obs.DoOnNext(func(i interface{}) {
		values = append(values, i.(int))
	})

	// add more
	for i := 3; i < 5; i++ {
		subject.Next(i)
		// slow down to let subscriber read from buffer
		time.Sleep(10 * time.Millisecond)
	}

	assert.Equal(t, []int{0, 1, 2, 3, 4}, values)
	fmt.Printf("values: %v", values)
}

// TestMaxItemsReplay verifies only the last n elements are kept in replay buffer
func TestMaxItemsReplay(t *testing.T) {
	subject := NewReplaySubject(2)

	// load buffer, expect to keep 2,3 in buffer
	for i := 0; i < 4; i++ {
		subject.Next(i)
	}

	_, obs := subject.Subscribe()

	values := make([]int, 0)
	obs.DoOnNext(func(i interface{}) {
		values = append(values, i.(int))
	})

	// add more
	for i := 4; i < 6; i++ {
		subject.Next(i)
		// slow down to let subscriber read from buffer
		time.Sleep(10 * time.Millisecond)
	}

	assert.Equal(t, []int{2, 3, 4, 5}, values)
	fmt.Printf("values: %v", values)
}
