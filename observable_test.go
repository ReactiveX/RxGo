package rx

import (
	"errors"
	"net/http"
	"testing"
	"time"

	"github.com/reactivex/rxgo/fx"
	"github.com/reactivex/rxgo/handlers"
	"github.com/reactivex/rxgo/iterable"
	"github.com/reactivex/rxgo/optional"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
	"sync/atomic"
)

// TODO
//func TestNewFromChannel(t *testing.T) {
//	ch := make(chan interface{}, 5)
//
//	observable := NewObservableFromChannel(ch)
//	switch v := observable.(type) {
//	case *observable:
//		assert.Exactly(t, ch, v.ch)
//	default:
//		t.Fail()
//	}
//}

func TestCreateObservableWithConstructor(t *testing.T) {
	assert := assert.New(t)

	stream1 := NewObservable(0)
	stream2 := NewObservable(3)

	switch v := stream1.(type) {
	case *observable:
		assert.Equal(0, cap(v.ch))
	default:
		t.Fail()
	}

	switch v := stream2.(type) {
	case *observable:
		assert.Equal(3, cap(v.ch))
	default:
		t.Fail()
	}
}

func TestCheckEventHandler(t *testing.T) {
	if testing.Short() {
		t.Skip("Skip testing of unexported testCheckEventHandler")
	}

	testtext := ""

	df := handlers.DoneFunc(func() {
		testtext += "done"
	})

	myObserver := NewObserver(df)

	ob1 := CheckEventHandler(myObserver)
	ob2 := CheckEventHandler(df)

	ob1.OnDone()
	assert.Equal(t, "done", testtext)

	ob2.OnDone()
	assert.Equal(t, "donedone", testtext)
}

func TestEmptyOperator(t *testing.T) {
	myStream := Empty()
	text := ""

	onDone := handlers.DoneFunc(func() {
		text += "done"
	})
	myStream.Subscribe(onDone).Block()

	assert.Equal(t, "done", text)
}

func TestIntervalOperator(t *testing.T) {
	fin := make(chan struct{})
	myStream := Interval(fin, 10*time.Millisecond)
	nums := []int{}

	onNext := handlers.NextFunc(func(item interface{}) {
		if num, ok := item.(int); ok {
			if num >= 5 {
				fin <- struct{}{}
				close(fin)
			}
			nums = append(nums, num)
		}
	})

	myStream.Subscribe(onNext).Block()

	assert.Exactly(t, []int{0, 1, 2, 3, 4, 5}, nums)
}

func TestRangeOperator(t *testing.T) {
	myStream := Range(2, 6)
	nums := []int{}

	onNext := handlers.NextFunc(func(item interface{}) {
		if num, ok := item.(int); ok {
			nums = append(nums, num)
		}
	})

	onDone := handlers.DoneFunc(func() {
		nums = append(nums, 1000)
	})

	myStream.Subscribe(NewObserver(onNext, onDone)).Block()

	assert.Exactly(t, []int{2, 3, 4, 5, 1000}, nums)
}

func TestJustOperator(t *testing.T) {
	myStream := Just(1, 2.01, "foo", map[string]string{"bar": "baz"}, 'a')
	//numItems := 5
	//yes := make(chan struct{})
	stuff := []interface{}{}

	onNext := handlers.NextFunc(func(item interface{}) {
		stuff = append(stuff, item)
	})

	myStream.Subscribe(onNext).Block()

	expected := []interface{}{
		1, 2.01, "foo", map[string]string{"bar": "baz"}, 'a',
	}

	assert.Exactly(t, expected, stuff)
}

func TestFromOperator(t *testing.T) {
	items := []interface{}{1, 3.1416, &struct{ foo string }{"bar"}}
	it, err := iterable.New(items)
	if err != nil {
		t.Fail()
	}

	myStream := From(it)
	nums := []interface{}{}

	onNext := handlers.NextFunc(func(item interface{}) {
		switch item := item.(type) {
		case int, float64:
			nums = append(nums, item)
		}
	})

	myStream.Subscribe(onNext).Block()

	expected := []interface{}{1, 3.1416}
	assert.Len(t, nums, 2)
	for n, num := range nums {
		assert.Equal(t, expected[n], num)
	}
}

func fakeGet(url string, delay time.Duration, result interface{}) (interface{}, error) {
	<-time.After(delay)
	if err, isErr := result.(error); isErr {
		return nil, err
	}

	return result, nil
}

func TestStartOperator(t *testing.T) {

	responseCodes := []int{}
	done := false

	d1 := fx.Supplier(func() interface{} {
		result := &http.Response{
			Status:     "200 OK",
			StatusCode: 200,
		}

		res, err := fakeGet("somehost.com", 10*time.Millisecond, result)
		if err != nil {
			return err
		}
		return res
	})

	d2 := fx.Supplier(func() interface{} {
		result := &http.Response{
			Status:     "301 Moved Permanently",
			StatusCode: 301,
		}

		res, err := fakeGet("somehost.com", 15*time.Millisecond, result)
		if err != nil {
			return err
		}
		return res
	})

	d3 := fx.Supplier(func() interface{} {
		result := &http.Response{
			Status:     "500 Server Error",
			StatusCode: 500,
		}

		res, err := fakeGet("badserver.co", 50*time.Millisecond, result)
		if err != nil {
			return err
		}
		return res
	})

	e1 := fx.Supplier(func() interface{} {
		err := errors.New("Bad URL")
		res, err := fakeGet("badurl.err", 100*time.Millisecond, err)
		if err != nil {
			return err
		}
		return res
	})

	d4 := fx.Supplier(func() interface{} {
		result := &http.Response{
			Status:     "404 Not Found",
			StatusCode: 400,
		}
		res, err := fakeGet("notfound.org", 200*time.Millisecond, result)
		if err != nil {
			return err
		}
		return res
	})

	onNext := handlers.NextFunc(func(item interface{}) {
		if res, ok := item.(*http.Response); ok {
			responseCodes = append(responseCodes, res.StatusCode)
		}
	})

	onError := handlers.ErrFunc(func(err error) {
		t.Log(err)
	})

	onDone := handlers.DoneFunc(func() {
		done = true
	})

	myObserver := NewObserver(onNext, onError, onDone)

	myStream := Start(d1, d3, d4, e1, d2)

	err := myStream.Subscribe(myObserver).Block()

	assert.Exactly(t, []int{200, 301, 500}, responseCodes)
	assert.Equal(t, "Bad URL", err.Error())

	// Error should have prevented OnDone from being called
	assert.False(t, done)
}

func TestSubscribeToNextFunc(t *testing.T) {
	myStream := Just(1, 2, 3, errors.New("4"), 5)
	mynum := 0

	nf := handlers.NextFunc(func(item interface{}) {
		if num, ok := item.(int); ok {
			mynum += num
		}
	})

	myStream.Subscribe(nf).Block()

	assert.Equal(t, 6, mynum)
}

func TestSubscribeToErrFunc(t *testing.T) {
	myStream := Just(1, "hello", errors.New("bang"), 43.5)

	var myerr error

	ef := handlers.ErrFunc(func(err error) {
		myerr = err
	})

	sub := myStream.Subscribe(ef).Block()

	assert.Equal(t, "bang", myerr.Error())
	assert.Equal(t, "bang", sub.Error())
}

func TestSubscribeToDoneFunc(t *testing.T) {
	myStream := Just(nil)

	donetext := ""

	df := handlers.DoneFunc(func() {
		donetext = "done"
	})

	myStream.Subscribe(df).Block()

	assert.Equal(t, "done", donetext)
}

func TestSubscribeToObserver(t *testing.T) {
	assert := assert.New(t)

	it, err := iterable.New([]interface{}{
		"foo", "bar", "baz", 'a', 'b', errors.New("bang"), 99,
	})
	if err != nil {
		t.Fail()
	}
	myStream := From(it)

	words := []string{}
	chars := []rune{}
	integers := []int{}
	finished := false

	onNext := handlers.NextFunc(func(item interface{}) {
		switch item := item.(type) {
		case string:
			words = append(words, item)
		case rune:
			chars = append(chars, item)
		case int:
			integers = append(integers, item)
		}
	})

	onError := handlers.ErrFunc(func(err error) {
		t.Logf("Error emitted in the stream: %v\n", err)
	})

	onDone := handlers.DoneFunc(func() {
		finished = true
	})

	ob := NewObserver(onNext, onError, onDone)

	done := myStream.Subscribe(ob)
	sub := done.Block()

	assert.Empty(integers)
	assert.False(finished)

	expectedWords := []string{"foo", "bar", "baz"}
	assert.Len(words, len(expectedWords))
	for n, word := range words {
		assert.Equal(expectedWords[n], word)
	}

	expectedChars := []rune{'a', 'b'}
	assert.Len(chars, len(expectedChars))
	for n, char := range chars {
		assert.Equal(expectedChars[n], char)
	}

	assert.Equal("bang", sub.Error())
}

func TestObservableMap(t *testing.T) {
	items := []interface{}{1, 2, 3, "foo", "bar", []byte("baz")}
	it, err := iterable.New(items)
	if err != nil {
		t.Fail()
	}

	stream1 := From(it)

	multiplyAllIntBy := func(factor interface{}) fx.Function {
		return func(item interface{}) interface{} {
			if num, ok := item.(int); ok {
				return num * factor.(int)
			}
			return item
		}
	}
	stream2 := stream1.Map(multiplyAllIntBy(10))

	nums := []int{}
	onNext := handlers.NextFunc(func(item interface{}) {
		if num, ok := item.(int); ok {
			nums = append(nums, num)
		}
	})

	stream2.Subscribe(onNext).Block()

	assert.Exactly(t, []int{10, 20, 30}, nums)
}

func TestObservableTake(t *testing.T) {
	items := []interface{}{1, 2, 3, 4, 5}
	it, err := iterable.New(items)
	if err != nil {
		t.Fail()
	}

	stream1 := From(it)
	stream2 := stream1.Take(3)

	nums := []int{}
	onNext := handlers.NextFunc(func(item interface{}) {
		if num, ok := item.(int); ok {
			nums = append(nums, num)
		}
	})

	stream2.Subscribe(onNext).Block()

	assert.Exactly(t, []int{1, 2, 3}, nums)
}

func TestObservableTakeWithEmpty(t *testing.T) {
	stream1 := Empty()
	stream2 := stream1.Take(3)

	nums := []int{}
	onNext := handlers.NextFunc(func(item interface{}) {
		if num, ok := item.(int); ok {
			nums = append(nums, num)
		}
	})

	stream2.Subscribe(onNext).Block()

	assert.Exactly(t, []int{}, nums)
}

func TestObservableTakeLast(t *testing.T) {
	items := []interface{}{1, 2, 3, 4, 5}
	it, err := iterable.New(items)
	if err != nil {
		t.Fail()
	}

	stream1 := From(it)
	stream2 := stream1.TakeLast(3)

	nums := []int{}
	onNext := handlers.NextFunc(func(item interface{}) {
		if num, ok := item.(int); ok {
			nums = append(nums, num)
		}
	})

	stream2.Subscribe(onNext).Block()

	assert.Exactly(t, []int{3, 4, 5}, nums)
}

/*
func TestObservableTakeLastWithEmpty(t *testing.T) {
	stream1 := Empty()
	stream2 := stream1.TakeLast(3)

	nums := []int{}
	onNext := handlers.NextFunc(func(item interface{}) {
		if num, ok := item.(int); ok {
			nums = append(nums, num)
		}
	})

	sub := stream2.Subscribe(onNext)
	<-sub

	assert.Exactly(t, []int{}, nums)
}*/

func TestObservableFilter(t *testing.T) {
	items := []interface{}{1, 2, 3, 120, []byte("baz"), 7, 10, 13}
	it, err := iterable.New(items)
	if err != nil {
		t.Fail()
	}

	stream1 := From(it)

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

	stream2 := stream1.Filter(lt(9))

	nums := []int{}
	onNext := handlers.NextFunc(func(item interface{}) {
		if num, ok := item.(int); ok {
			nums = append(nums, num)
		}
	})

	stream2.Subscribe(onNext).Block()

	assert.Exactly(t, []int{1, 2, 3, 7}, nums)
}

func TestObservableFirst(t *testing.T) {
	items := []interface{}{0, 1, 3}
	it, err := iterable.New(items)
	if err != nil {
		t.Fail()
	}

	stream1 := From(it)
	stream2 := stream1.First()

	nums := []int{}
	onNext := handlers.NextFunc(func(item interface{}) {
		if num, ok := item.(int); ok {
			nums = append(nums, num)
		}
	})

	stream2.Subscribe(onNext).Block()

	assert.Exactly(t, []int{0}, nums)
}

func TestObservableFirstWithEmpty(t *testing.T) {
	stream1 := Empty()

	stream2 := stream1.First()

	nums := []int{}
	onNext := handlers.NextFunc(func(item interface{}) {
		if num, ok := item.(int); ok {
			nums = append(nums, num)
		}
	})

	stream2.Subscribe(onNext).Block()

	assert.Exactly(t, []int{}, nums)
}

func TestObservableLast(t *testing.T) {
	items := []interface{}{0, 1, 3}
	it, err := iterable.New(items)
	if err != nil {
		t.Fail()
	}

	stream1 := From(it)

	stream2 := stream1.Last()

	nums := []int{}
	onNext := handlers.NextFunc(func(item interface{}) {
		if num, ok := item.(int); ok {
			nums = append(nums, num)
		}
	})

	stream2.Subscribe(onNext).Block()

	assert.Exactly(t, []int{3}, nums)
}

func TestParallelSubscribeToObserver(t *testing.T) {
	assert := assert.New(t)

	it, err := iterable.New([]interface{}{
		"foo", "bar", "baz", 'a', 'b', 99,
	})
	if err != nil {
		t.Fail()
	}
	myStream := From(it)

	var wordsCount uint64
	var charsCount uint64
	var integersCount uint64
	finished := false

	onNext := handlers.NextFunc(func(item interface{}) {
		switch item.(type) {
		case string:
			atomic.AddUint64(&wordsCount, 1)
		case rune:
			atomic.AddUint64(&charsCount, 1)
		case int:
			atomic.AddUint64(&integersCount, 1)
		}
	})

	onError := handlers.ErrFunc(func(err error) {
		t.Logf("Error emitted in the stream: %v\n", err)
	})

	onDone := handlers.DoneFunc(func() {
		finished = true
	})

	ob := NewObserver(onNext, onError, onDone)

	myStream.Subscribe(ob, WithParallelism(2)).Block()

	assert.True(finished)

	assert.Equal(integersCount, uint64(0x1))
	assert.Equal(wordsCount, uint64(0x3))
	assert.Equal(charsCount, uint64(0x2))
}

func TestParallelSubscribeToObserverWithError(t *testing.T) {
	assert := assert.New(t)

	it, err := iterable.New([]interface{}{
		"foo", "bar", "baz", 'a', 'b', 99, errors.New("error"),
	})
	if err != nil {
		t.Fail()
	}
	myStream := From(it)

	finished := false

	onNext := handlers.NextFunc(func(item interface{}) {
	})

	onError := handlers.ErrFunc(func(err error) {
		t.Logf("Error emitted in the stream: %v\n", err)
	})

	onDone := handlers.DoneFunc(func() {
		finished = true
	})

	ob := NewObserver(onNext, onError, onDone)

	myStream.Subscribe(ob, WithParallelism(2)).Block()

	assert.False(finished)
}

func TestObservableLastWithEmpty(t *testing.T) {
	stream1 := Empty()

	stream2 := stream1.Last()

	nums := []int{}
	onNext := handlers.NextFunc(func(item interface{}) {
		if num, ok := item.(int); ok {
			nums = append(nums, num)
		}
	})

	stream2.Subscribe(onNext).Block()

	assert.Exactly(t, []int{}, nums)
}

func TestObservableSkip(t *testing.T) {
	items := []interface{}{0, 1, 3, 5, 1, 8}
	it, err := iterable.New(items)
	if err != nil {
		t.Fail()
	}

	stream1 := From(it)

	stream2 := stream1.Skip(3)

	nums := []int{}
	onNext := handlers.NextFunc(func(item interface{}) {
		if num, ok := item.(int); ok {
			nums = append(nums, num)
		}
	})

	stream2.Subscribe(onNext).Block()

	assert.Exactly(t, []int{5, 1, 8}, nums)
}

func TestObservableSkipWithEmpty(t *testing.T) {
	stream1 := Empty()

	stream2 := stream1.Skip(3)

	nums := []int{}
	onNext := handlers.NextFunc(func(item interface{}) {
		if num, ok := item.(int); ok {
			nums = append(nums, num)
		}
	})

	stream2.Subscribe(onNext).Block()

	assert.Exactly(t, []int{}, nums)
}

func TestObservableSkipLast(t *testing.T) {
	items := []interface{}{0, 1, 3, 5, 1, 8}
	it, err := iterable.New(items)
	if err != nil {
		t.Fail()
	}

	stream1 := From(it)

	stream2 := stream1.SkipLast(3)

	nums := []int{}
	onNext := handlers.NextFunc(func(item interface{}) {
		if num, ok := item.(int); ok {
			nums = append(nums, num)
		}
	})

	stream2.Subscribe(onNext).Block()

	assert.Exactly(t, []int{0, 1, 3}, nums)
}

func TestObservableSkipLastWithEmpty(t *testing.T) {
	stream1 := Empty()

	stream2 := stream1.SkipLast(3)

	nums := []int{}
	onNext := handlers.NextFunc(func(item interface{}) {
		if num, ok := item.(int); ok {
			nums = append(nums, num)
		}
	})

	stream2.Subscribe(onNext).Block()

	assert.Exactly(t, []int{}, nums)
}

func TestObservableDistinct(t *testing.T) {
	items := []interface{}{1, 2, 2, 1, 3}
	it, err := iterable.New(items)
	if err != nil {
		t.Fail()
	}

	stream1 := From(it)

	id := func(item interface{}) interface{} {
		return item
	}

	stream2 := stream1.Distinct(id)

	nums := []int{}
	onNext := handlers.NextFunc(func(item interface{}) {
		if num, ok := item.(int); ok {
			nums = append(nums, num)
		}
	})

	stream2.Subscribe(onNext).Block()

	assert.Exactly(t, []int{1, 2, 3}, nums)
}

func TestObservableDistinctUntilChanged(t *testing.T) {
	items := []interface{}{1, 2, 2, 1, 3}
	it, err := iterable.New(items)
	if err != nil {
		t.Fail()
	}

	stream1 := From(it)

	id := func(item interface{}) interface{} {
		return item
	}

	stream2 := stream1.DistinctUntilChanged(id)

	nums := []int{}
	onNext := handlers.NextFunc(func(item interface{}) {
		if num, ok := item.(int); ok {
			nums = append(nums, num)
		}
	})

	stream2.Subscribe(onNext).Block()

	assert.Exactly(t, []int{1, 2, 1, 3}, nums)
}

func TestObservableScanWithIntegers(t *testing.T) {
	items := []interface{}{0, 1, 3, 5, 1, 8}
	it, err := iterable.New(items)
	if err != nil {
		t.Fail()
	}

	stream1 := From(it)

	stream2 := stream1.Scan(func(x, y interface{}) interface{} {
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

	stream2.Subscribe(onNext).Block()

	assert.Exactly(t, []int{0, 1, 4, 9, 10, 18}, nums)
}

func TestObservableScanWithString(t *testing.T) {
	items := []interface{}{"hello", "world", "this", "is", "foo"}
	it, err := iterable.New(items)
	if err != nil {
		t.Fail()
	}

	stream1 := From(it)

	stream2 := stream1.Scan(func(x, y interface{}) interface{} {
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

	stream2.Subscribe(onNext).Block()

	expected := []string{
		"hello",
		"helloworld",
		"helloworldthis",
		"helloworldthisis",
		"helloworldthisisfoo",
	}

	assert.Exactly(t, expected, words)
}

func TestRepeatInfinityOperator(t *testing.T) {
	myStream := Repeat("mystring")

	item, err := myStream.Next()

	if err != nil {
		assert.Fail(t, "fail to emit next item", err)
	}

	if value, ok := item.(string); ok {
		assert.Equal(t, value, "mystring")
	} else {
		assert.Fail(t, "fail to emit next item", err)
	}
}

func TestRepeatNtimeOperator(t *testing.T) {
	myStream := Repeat("mystring", 2)
	stringarray := []string{}

	onNext := handlers.NextFunc(func(item interface{}) {
		if value, ok := item.(string); ok {
			stringarray = append(stringarray, value)
		}
	})

	onDone := handlers.DoneFunc(func() {
		stringarray = append(stringarray, "end")
	})

	myStream.Subscribe(NewObserver(onNext, onDone)).Block()

	assert.Exactly(t, []string{"mystring", "mystring", "end"}, stringarray)
}

func TestRepeatNtimeMultiVariadicOperator(t *testing.T) {
	myStream := Repeat("mystring", 2, 2, 3, 4, 5, 6, 7)
	stringarray := []string{}

	onNext := handlers.NextFunc(func(item interface{}) {
		if value, ok := item.(string); ok {
			stringarray = append(stringarray, value)
		}
	})

	onDone := handlers.DoneFunc(func() {
		stringarray = append(stringarray, "end")
	})

	myStream.Subscribe(NewObserver(onNext, onDone)).Block()

	assert.Exactly(t, []string{"mystring", "mystring", "end"}, stringarray)
}

func TestRepeatWithZeroNtimeOperator(t *testing.T) {
	myStream := Repeat("mystring", 0)
	stringarray := []string{}

	onNext := handlers.NextFunc(func(item interface{}) {
		if value, ok := item.(string); ok {
			stringarray = append(stringarray, value)
		}
	})

	onDone := handlers.DoneFunc(func() {
		stringarray = append(stringarray, "end")
	})

	myStream.Subscribe(NewObserver(onNext, onDone)).Block()

	assert.Exactly(t, []string{"end"}, stringarray)
}

func TestRepeatWithNegativeTimesOperator(t *testing.T) {
	myStream := Repeat("mystring", -10)
	stringarray := []string{}

	onNext := handlers.NextFunc(func(item interface{}) {
		if value, ok := item.(string); ok {
			stringarray = append(stringarray, value)
		}
	})

	onDone := handlers.DoneFunc(func() {
		stringarray = append(stringarray, "end")
	})

	myStream.Subscribe(NewObserver(onNext, onDone)).Block()

	assert.Exactly(t, []string{"end"}, stringarray)
}

func TestEmptyCompletesSequence(t *testing.T) {
	// given
	emissionObserver := NewObserverMock()

	// and empty sequence
	sequence := Empty()

	// when subscribes to the sequence
	sequence.Subscribe(emissionObserver.Capture()).Block()

	// then completes without any emission
	emissionObserver.AssertNotCalled(t, "OnNext", mock.Anything)
	emissionObserver.AssertNotCalled(t, "OnError", mock.Anything)
	emissionObserver.AssertCalled(t, "OnDone")
}

func TestError(t *testing.T) {
	var got error
	err := errors.New("foo")
	stream := Error(err)
	stream.Subscribe(handlers.ErrFunc(func(e error) {
		got = e
	})).Block()

	assert.Equal(t, err, got)
}

func TestElementAt(t *testing.T) {
	got := 0
	just := Just(0, 1, 2, 3, 4)
	single := just.ElementAt(2)
	single.Subscribe(handlers.NextFunc(func(i interface{}) {
		switch i := i.(type) {
		case int:
			got = i
		}
	})).Block()

	assert.Equal(t, 2, got)
}

func TestElementAtWithError(t *testing.T) {
	got := 0
	just := Just(0, 1, 2, 3, 4)
	single := just.ElementAt(10)
	single.Subscribe(handlers.ErrFunc(func(error) {
		got = 10
	})).Block()

	assert.Equal(t, 10, got)
}

func TestObservableReduce(t *testing.T) {
	items := []interface{}{1, 2, 3, 4, 5}
	it, err := iterable.New(items)
	if err != nil {
		t.Fail()
	}
	stream1 := From(it)
	add := func(acc interface{}, elem interface{}) interface{} {
		if a, ok := acc.(int); ok {
			if b, ok := elem.(int); ok {
				return a + b
			}
		}
		return 0
	}

	var got optional.Optional
	_, err = stream1.Reduce(add).Subscribe(handlers.NextFunc(func(i interface{}) {
		got = i.(optional.Optional)
	})).Block()
	if err != nil {
		t.Fail()
	}
	assert.False(t, got.IsEmpty())
	assert.Exactly(t, optional.Of(14), got)
}

func TestObservableReduceEmpty(t *testing.T) {
	it, err := iterable.New([]interface{}{})
	if err != nil {
		t.Fail()
	}
	add := func(acc interface{}, elem interface{}) interface{} {
		if a, ok := acc.(int); ok {
			if b, ok := elem.(int); ok {
				return a + b
			}
		}
		return 0
	}
	stream := From(it)

	var got optional.Optional
	_, err = stream.Reduce(add).Subscribe(handlers.NextFunc(func(i interface{}) {
		got = i.(optional.Optional)
	})).Block()
	if err != nil {
		t.Fail()
	}
	assert.True(t, got.IsEmpty())
}

func TestObservableReduceNil(t *testing.T) {
	items := []interface{}{1, 2, 3, 4, 5}
	it, err := iterable.New(items)
	if err != nil {
		t.Fail()
	}
	stream := From(it)
	nilReduce := func(acc interface{}, elem interface{}) interface{} {
		return nil
	}
	var got optional.Optional
	_, err = stream.Reduce(nilReduce).Subscribe(handlers.NextFunc(func(i interface{}) {
		got = i.(optional.Optional)
	})).Block()
	if err != nil {
		t.Fail()
	}
	assert.False(t, got.IsEmpty())
	assert.Nil(t, got.Get())
}

func TestObservableCount(t *testing.T) {
	items := []interface{}{1, 2, 3, "foo", "bar", errors.New("error")}
	it, err := iterable.New(items)
	if err != nil {
		t.Fail()
	}
	stream := From(it)
	count, err := stream.Count().Subscribe(nil).Block()
	if err != nil {
		t.Fail()
	}
	assert.Exactly(t, int64(6), count)
}

func TestObservableFirstOrDefault(t *testing.T) {
	var items []interface{}
	it, err := iterable.New(items)
	if err != nil {
		t.Fail()
	}
	v, err := From(it).FirstOrDefault(7).Subscribe(nil).Block()
	if err != nil {
		t.Fail()
	}
	assert.Exactly(t, 7, v)
}

func TestObservableFirstOrDefaultWithValue(t *testing.T) {
	items := []interface{}{0, 1, 2}
	it, err := iterable.New(items)
	if err != nil {
		t.Fail()
	}
	v, err := From(it).FirstOrDefault(7).Subscribe(nil).Block()
	if err != nil {
		t.Fail()
	}
	assert.Exactly(t, 0, v)
}

func TestObservableLastOrDefault(t *testing.T) {
	var items []interface{}
	it, err := iterable.New(items)
	if err != nil {
		t.Fail()
	}
	v, err := From(it).LastOrDefault(7).Subscribe(nil).Block()
	if err != nil {
		t.Fail()
	}
	assert.Exactly(t, 7, v)
}

func TestObservableLastOrDefaultWithValue(t *testing.T) {
	items := []interface{}{0, 1, 3}
	it, err := iterable.New(items)
	if err != nil {
		t.Fail()
	}
	v, err := From(it).LastOrDefault(7).Subscribe(nil).Block()
	if err != nil {
		t.Fail()
	}
	assert.Exactly(t, 3, v)
}

func TestObservableTakeWhile(t *testing.T) {
	items := []interface{}{1, 2, 3, 4, 5}
	it, err := iterable.New(items)
	if err != nil {
		t.Fail()
	}
	stream1 := From(it)
	stream2 := stream1.TakeWhile(func(item interface{}) bool {
		return item != 3
	})
	nums := []int{}
	onNext := handlers.NextFunc(func(item interface{}) {
		if num, ok := item.(int); ok {
			nums = append(nums, num)
		}
	})
	stream2.Subscribe(onNext).Block()
	assert.Exactly(t, []int{1, 2}, nums)
}

func TestObservableTakeWhileWithEmpty(t *testing.T) {
	stream1 := Empty()
	stream2 := stream1.TakeWhile(func(item interface{}) bool {
		return item != 3
	})
	nums := []int{}
	onNext := handlers.NextFunc(func(item interface{}) {
		if num, ok := item.(int); ok {
			nums = append(nums, num)
		}
	})
	stream2.Subscribe(onNext).Block()
	assert.Exactly(t, []int{}, nums)
}

func TestObservableToList(t *testing.T) {
	items := []interface{}{1, "hello", false, .0}
	it, err := iterable.New(items)
	if err != nil {
		t.Fail()
	}
	var got interface{}
	stream1 := From(it)
	stream1.ToList().Subscribe(handlers.NextFunc(func(i interface{}) {
		got = i
	})).Block()
	assert.Exactly(t, []interface{}{1, "hello", false, .0}, got)
}

func TestObservableToListWithEmpty(t *testing.T) {
	stream1 := Empty()
	var got interface{}
	stream1.ToList().Subscribe(handlers.NextFunc(func(i interface{}) {
		got = i
	})).Block()
	assert.Exactly(t, []interface{}{}, got)
}
