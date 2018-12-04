package rxgo

import (
	"time"

	"github.com/stretchr/testify/mock"
)

var Indefinitely int64 = -1

type Duration interface {
	duration() time.Duration
}

type duration struct {
	d time.Duration
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

func WithDuration(d time.Duration) Duration {
	return &duration{
		d: d,
	}
}
