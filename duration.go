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

type testDuration struct {
	fs []func()
}

var tick = struct{}{}

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
			fs[i] = func() {
				ch <- Of(elem)
			}
		}
	}
	fs[len(elems)] = func() {
		cancel()
	}
	return ctx, FromChannel(ch), &testDuration{fs: fs}
}

func (d *testDuration) append(fs ...func()) {
	if d.fs == nil {
		d.fs = make([]func(), 0)
	}
	for _, f := range fs {
		d.fs = append(d.fs, f)
	}
}

func (d *testDuration) duration() time.Duration {
	d.fs[0]()
	d.fs = d.fs[1:]
	return 0
}

type mockDuration struct {
	mock.Mock
}

func (m *mockDuration) duration() time.Duration {
	args := m.Called()
	return args.Get(0).(time.Duration)
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
