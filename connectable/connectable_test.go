package connectable

import (
	"errors"
	"github.com/reactivex/rxgo"
	"github.com/reactivex/rxgo/fx"
	"github.com/reactivex/rxgo/handlers"
	"github.com/reactivex/rxgo/iterable"
	"github.com/stretchr/testify/assert"
	"sync"
	"testing"
	"time"
)

func TestNew(t *testing.T) {
	observers := make([]rx.Observer, 5, 5)
	connectable := New(0, observers...)

	switch v := connectable.(type) {
	case *connector:
		assert.Exactly(t, observers, v.observers)
	default:
		t.Fail()
	}
}

func TestDoOperator(t *testing.T) {
	co := Just(1, 2, 3)
	num := 0

	nextf := func(item interface{}) {
		num += item.(int)
	}

	co = co.Do(nextf)
	sub := co.Connect()
	<-sub

	assert.Equal(t, 6, num)
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

	it, err := iterable.New([]interface{}{1, 2, 3, errors.New("bang"), 9})
	if err != nil {
		t.Fail()
	}
	co := From(it)

	onNext := handlers.NextFunc(func(item interface{}) {
		num += item.(int)
	})

	onError := handlers.ErrFunc(func(err error) {
		myErr = err
	})

	onDone := handlers.DoneFunc(func() {
		done = "done"
	})

	ob := rx.NewObserver(onError, onDone, onNext)

	sub := co.Subscribe(ob).Connect()

	for c := range sub {
		s := c.Block()
		assert.Equal("bang", s.Error())
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

	it, err := iterable.New([]interface{}{1, 2, 3, errors.New("bang"), 9})
	if err != nil {
		t.Fail()
	}

	co := From(it)

	var mutex = &sync.Mutex{}

	ob1 := rx.NewObserver(
		handlers.NextFunc(func(item interface{}) {
			<-time.After(100 * time.Millisecond)
			mutex.Lock()
			nums = append(nums, item.(int))
			mutex.Unlock()
		}), handlers.ErrFunc(func(err error) {
			mutex.Lock()
			errs = append(errs, err)
			mutex.Unlock()
		}), handlers.DoneFunc(func() {
			mutex.Lock()
			dones = append(dones, "D1")
			mutex.Unlock()
		}))

	ob2 := rx.NewObserver(
		handlers.NextFunc(func(item interface{}) {
			mutex.Lock()
			nums = append(nums, item.(int)*2)
			mutex.Unlock()
		}), handlers.ErrFunc(func(err error) {
			mutex.Lock()
			errs = append(errs, err)
			mutex.Unlock()
		}), handlers.DoneFunc(func() {
			mutex.Lock()
			dones = append(dones, "D2")
			mutex.Unlock()
		}))

	ob3 := handlers.NextFunc(func(item interface{}) {
		<-time.After(200 * time.Millisecond)
		mutex.Lock()
		nums = append(nums, item.(int)*10)
		mutex.Unlock()
	})

	co = co.Subscribe(ob1).Subscribe(ob3).Subscribe(ob2)
	subs := co.Connect()

	for sub := range subs {
		s := sub.Block()
		assert.Equal("bang", s.Error())
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
	it, err := iterable.New(items)
	if err != nil {
		t.Fail()
	}

	stream := From(it)

	multiplyAllIntBy := func(factor interface{}) fx.Function {
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
	it, err := iterable.New(items)
	if err != nil {
		t.Fail()
	}

	stream := From(it)

	lt := func(target interface{}) fx.Predicate {
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
	it, err := iterable.New(items)
	if err != nil {
		t.Fail()
	}
	stream := From(it)

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
	it, err := iterable.New(items)
	if err != nil {
		t.Fail()
	}

	stream := From(it)

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

func TestConnectableFirst(t *testing.T) {
	items := []interface{}{0, 1, 3}
	it, err := iterable.New(items)
	if err != nil {
		t.Fail()
	}

	co1 := From(it)
	co2 := co1.First()

	nums := []int{}
	onNext := handlers.NextFunc(func(item interface{}) {
		if num, ok := item.(int); ok {
			nums = append(nums, num)
		}
	})

	subs := co2.Subscribe(onNext).Connect()
	<-subs

	assert.Exactly(t, []int{0}, nums)
}

func TestConnectableFirstWithEmpty(t *testing.T) {
	co1 := Empty()

	co2 := co1.First()

	nums := []int{}
	onNext := handlers.NextFunc(func(item interface{}) {
		if num, ok := item.(int); ok {
			nums = append(nums, num)
		}
	})

	subs := co2.Subscribe(onNext).Connect()
	<-subs

	assert.Exactly(t, []int{}, nums)
}

func TestObservableLast(t *testing.T) {
	items := []interface{}{0, 1, 3}
	it, err := iterable.New(items)
	if err != nil {
		t.Fail()
	}

	co := From(it)

	co = co.Last()

	nums := []int{}
	onNext := handlers.NextFunc(func(item interface{}) {
		if num, ok := item.(int); ok {
			nums = append(nums, num)
		}
	})

	subs := co.Subscribe(onNext).Connect()
	<-subs

	assert.Exactly(t, []int{3}, nums)
}

func TestObservableLastWithEmpty(t *testing.T) {
	co := Empty()

	co = co.Last()

	nums := []int{}
	onNext := handlers.NextFunc(func(item interface{}) {
		if num, ok := item.(int); ok {
			nums = append(nums, num)
		}
	})

	subs := co.Subscribe(onNext).Connect()
	<-subs

	assert.Exactly(t, []int{}, nums)
}

func TestConnectableDistinct(t *testing.T) {
	items := []interface{}{1, 2, 2, 1, 3}
	it, err := iterable.New(items)
	if err != nil {
		t.Fail()
	}

	co := From(it)

	id := func(item interface{}) interface{} {
		return item
	}

	co = co.Distinct(id)

	nums := []int{}
	onNext := handlers.NextFunc(func(item interface{}) {
		if num, ok := item.(int); ok {
			nums = append(nums, num)
		}
	})

	subs := co.Subscribe(onNext).Connect()
	<-subs

	assert.Exactly(t, []int{1, 2, 3}, nums)
}

func TestConnectableDistinctUntilChanged(t *testing.T) {
	items := []interface{}{1, 2, 2, 1, 3}
	it, err := iterable.New(items)
	if err != nil {
		t.Fail()
	}

	co := From(it)

	id := func(item interface{}) interface{} {
		return item
	}

	co = co.DistinctUntilChanged(id)

	nums := []int{}
	onNext := handlers.NextFunc(func(item interface{}) {
		if num, ok := item.(int); ok {
			nums = append(nums, num)
		}
	})

	subs := co.Subscribe(onNext).Connect()
	<-subs

	assert.Exactly(t, []int{1, 2, 1, 3}, nums)
}
