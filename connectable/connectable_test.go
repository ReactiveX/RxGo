package connectable

import (
	"errors"
	"testing"
	"time"

	"github.com/jochasinga/grx/fx"
	"github.com/jochasinga/grx/handlers"
	"github.com/jochasinga/grx/observer"

	"github.com/stretchr/testify/assert"
)

func TestCreateConnectableWithConstructor(t *testing.T) {
	assert := assert.New(t)
	text := "hello"
	co1 := New(0)
	co2 := New(3)
	co3 := Just("world")

	cotests := []struct {
		expect, suspect int
	}{
		{0, cap(co1.Observable)},
		{3, cap(co2.Observable)},
		{0, cap(co3.Observable)},
	}

	if assert.IsType(Connectable{}, co1) &&
		assert.IsType(Connectable{}, co2) &&
		assert.IsType(Connectable{}, co3) {

		for _, tt := range cotests {
			assert.Equal(tt.suspect, tt.expect)
		}
	}

	ob := observer.New(handlers.NextFunc(func(item interface{}) {
		text += item.(string)
	}),
	)

	co4 := New(0, ob)
	assert.Equal(0, cap(co4.Observable))

	co4.observers[0].OnNext("world")
	assert.Equal("helloworld", text)
}

func TestSubscribeToNextFunc(t *testing.T) {
	co := Just(1, 2, 3)
	num := 0

	onNext := handlers.NextFunc(func(item interface{}) {
		num += item.(int)
	})

	co = co.Subscribe(onNext)
	sub := co.Connect()
	<-sub

	assert.Equal(t, 6, num)
}

func TestSubscribeToErrFunc(t *testing.T) {
	co := Just(errors.New("bang"))

	var myError error

	onError := handlers.ErrFunc(func(err error) {
		myError = err
	})

	co = co.Subscribe(onError)
	sub := co.Connect()
	<-sub

	if assert.NotNil(t, myError) {
		assert.Equal(t, "bang", myError.Error())
	}
}

func TestSubscribeToDoneFunc(t *testing.T) {
	co := Empty()

	text := ""

	onDone := handlers.DoneFunc(func() {
		text += "done"
	})

	sub := co.Subscribe(onDone).Connect()
	<-sub

	if assert.NotEmpty(t, text) {
		assert.Equal(t, "done", text)
	}

}

func TestSubscribeToObserver(t *testing.T) {
	assert := assert.New(t)

	var (
		num   int
		myErr error
		done  string
	)

	co := From([]interface{}{1, 2, 3, errors.New("bang"), 9})

	onNext := handlers.NextFunc(func(item interface{}) {
		num += item.(int)
	})

	onError := handlers.ErrFunc(func(err error) {
		myErr = err
	})

	onDone := handlers.DoneFunc(func() {
		done = "done"
	})

	ob := observer.New(onError, onDone, onNext)

	sub := co.Subscribe(ob).Connect()
	//<-sub

	for c := range sub {
		for s := range c {
			assert.Equal("bang", s.Error.Error())
		}
	}

	assert.Equal(6, num)
	assert.Equal("bang", myErr.Error())
	assert.Empty(done)
}

func TestSubscribeToManyObservers(t *testing.T) {
	assert := assert.New(t)

	var (
		nums  []int
		errs  []error
		dones []string
	)

	co := From([]interface{}{1, 2, 3, errors.New("bang"), 9})

	ob1 := observer.Observer{
		NextHandler: func(item interface{}) {
			<-time.After(100 * time.Millisecond)
			nums = append(nums, item.(int))
		},
		ErrHandler: func(err error) {
			errs = append(errs, err)
		},
		DoneHandler: func() {
			dones = append(dones, "D1")
		},
	}

	ob2 := observer.Observer{
		NextHandler: func(item interface{}) {
			nums = append(nums, item.(int)*2)
		},
		ErrHandler: func(err error) {
			errs = append(errs, err)
		},
		DoneHandler: func() {
			dones = append(dones, "D2")
		},
	}

	ob3 := handlers.NextFunc(func(item interface{}) {
		<-time.After(200 * time.Millisecond)
		nums = append(nums, item.(int)*10)
	})

	co = co.Subscribe(ob1).Subscribe(ob3).Subscribe(ob2)
	subs := co.Connect()

	for sub := range subs {
		for s := range sub {
			assert.Equal("bang", s.Error.Error())
		}
	}

	expectedNums := []int{2, 4, 6, 1, 10, 2, 3, 20, 30}
	for _, num := range expectedNums {
		assert.Contains(nums, num)
	}

	expectedErr := errors.New("bang")
	assert.Exactly([]error{expectedErr, expectedErr}, errs)

	assert.Empty(dones)
}

func TestConnectableMap(t *testing.T) {
	items := []interface{}{1, 2, 3, "foo", "bar", []byte("baz")}
	stream := From(items)

	multiplyAllIntBy := func(factor interface{}) fx.MappableFunc {
		return func(item interface{}) interface{} {
			if num, ok := item.(int); ok {
				return num * factor.(int)
			}
			return item
		}
	}
	stream = stream.Map(multiplyAllIntBy(10))

	nums := []int{}
	onNext := handlers.NextFunc(func(item interface{}) {
		if num, ok := item.(int); ok {
			nums = append(nums, num)
		}
	})

	subs := stream.Subscribe(onNext).Connect()
	<-subs

	assert.Exactly(t, []int{10, 20, 30}, nums)
}

func TestConnectableFilter(t *testing.T) {
	items := []interface{}{1, 2, 3, 120, []byte("baz"), 7, 10, 13}
	stream := From(items)

	lt := func(target interface{}) fx.FilterableFunc {
		return func(item interface{}) bool {
			if num, ok := item.(int); ok {
				if num < 9 {
					return true
				}
			}
			return false
		}
	}

	stream = stream.Filter(lt(9))

	nums := []int{}
	onNext := handlers.NextFunc(func(item interface{}) {
		if num, ok := item.(int); ok {
			nums = append(nums, num)
		}
	})

	subs := stream.Subscribe(onNext).Connect()
	<-subs

	assert.Exactly(t, []int{1, 2, 3, 7}, nums)
}

func TestConnectableScanWithIntegers(t *testing.T) {
	items := []interface{}{0, 1, 3, 5, 1, 8}
	stream := From(items)

	stream = stream.Scan(func(x, y interface{}) interface{} {
		var v1, v2 int

		if x, ok := x.(int); ok {
			v1 = x
		}

		if y, ok := y.(int); ok {
			v2 = y
		}

		return v1 + v2
	})

	nums := []int{}
	onNext := handlers.NextFunc(func(item interface{}) {
		if num, ok := item.(int); ok {
			nums = append(nums, num)
		}
	})

	subs := stream.Subscribe(onNext).Connect()
	<-subs

	assert.Exactly(t, []int{0, 1, 4, 9, 10, 18}, nums)
}

func TestConnectableScanWithStrings(t *testing.T) {
	items := []interface{}{"hello", "world", "this", "is", "foo"}
	stream := From(items)

	stream = stream.Scan(func(x, y interface{}) interface{} {
		var w1, w2 string

		if x, ok := x.(string); ok {
			w1 = x
		}

		if y, ok := y.(string); ok {
			w2 = y
		}

		return w1 + w2
	})

	words := []string{}
	onNext := handlers.NextFunc(func(item interface{}) {
		if word, ok := item.(string); ok {
			words = append(words, word)
		}
	})

	subs := stream.Subscribe(onNext).Connect()
	<-subs

	expected := []string{
		"hello",
		"helloworld",
		"helloworldthis",
		"helloworldthisis",
		"helloworldthisisfoo",
	}

	assert.Exactly(t, expected, words)
}
