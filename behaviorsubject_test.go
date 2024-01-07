package rxgo

import (
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
)

// TestBehaviorSubject verifies that two subscribers started at different times receive the same result
func TestBehaviorSubject(t *testing.T) {
	subject := NewBehaviorSubject()

	// obs1 no last value
	_, obs1 := subject.Subscribe()
	values1 := make([]int, 0)
	obs1.DoOnNext(func(i interface{}) {
		values1 = append(values1, i.(int))
	})

	subject.Next(1)

	// obs2 receives last value and new
	_, obs2 := subject.Subscribe()
	values2 := make([]int, 0)
	obs2.DoOnNext(func(i interface{}) {
		values2 = append(values2, i.(int))
	})

	subject.Next(2)

	time.Sleep(10 * time.Millisecond)

	assert.Equal(t, []int{1, 2}, values1)
	assert.Equal(t, []int{1, 2}, values2)
}
