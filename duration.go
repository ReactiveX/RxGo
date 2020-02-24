package rxgo

import (
	"context"
	"time"

	"github.com/stretchr/testify/mock"
)

// Infinite represents an infinite wait time
var Infinite int64 = -1

// Duration represents a duration
type Duration interface {
	duration() time.Duration
}

type duration struct {
	d time.Duration
}

func (d *duration) duration() time.Duration {
	return d.d
}

// WithDuration is a duration option
func WithDuration(d time.Duration) Duration {
	return &duration{
		d: d,
	}
}

var tick = struct{}{}

type causalityDuration struct {
	fs []func()
}

func timeCausality(elems ...interface{}) (context.Context, Observable, Duration) {
	ch := make(chan Item, 1)
	fs := make([]func(), len(elems)+1)
	ctx, cancel := context.WithCancel(context.Background())
	for i, elem := range elems {
		i := i
		elem := elem
		if elem == tick {
			fs[i] = func() {}
		} else {
			switch elem := elem.(type) {
			default:
				fs[i] = func() {
					ch <- Of(elem)
				}
			case error:
				fs[i] = func() {
					ch <- Error(elem)
				}
			}
		}
	}
	fs[len(elems)] = func() {
		cancel()
	}
	return ctx, FromChannel(ch), &causalityDuration{fs: fs}
}

func (d *causalityDuration) duration() time.Duration {
	d.fs[0]()
	d.fs = d.fs[1:]
	return time.Nanosecond
}

type mockDuration struct {
	mock.Mock
}

func (m *mockDuration) duration() time.Duration {
	args := m.Called()
	return args.Get(0).(time.Duration)
}
