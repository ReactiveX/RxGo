package observable

import (
	"errors"
	"net/http"
	"testing"
	"time"

	"github.com/reactivex/rxgo/fx"
	"github.com/reactivex/rxgo/handlers"
	"github.com/reactivex/rxgo/iterable"
	"github.com/reactivex/rxgo/observer"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
)

func TestDefaultObservable(t *testing.T) {
	assert.Equal(t, 0, cap(DefaultObservable))
}

func TestCreateObservableWithConstructor(t *testing.T) {
	assert := assert.New(t)

	stream1 := New(0)
	stream2 := New(3)

	if assert.IsType(Observable(nil), stream1) && assert.IsType(Observable(nil), stream2) {
		assert.Equal(0, cap(stream1))
		assert.Equal(3, cap(stream2))
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

	myObserver := observer.New(df)

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
	sub := myStream.Subscribe(onDone)
	<-sub

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

	sub := myStream.Subscribe(onNext)
	<-sub

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

	sub := myStream.Subscribe(observer.New(onNext, onDone))
	<-sub

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

	sub := myStream.Subscribe(onNext)
	<-sub

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

	sub := myStream.Subscribe(onNext)
	<-sub

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

	d1 := fx.EmittableFunc(func() interface{} {
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

	d2 := fx.EmittableFunc(func() interface{} {
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

	d3 := fx.EmittableFunc(func() interface{} {
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

	e1 := fx.EmittableFunc(func() interface{} {
		err := errors.New("Bad URL")
		res, err := fakeGet("badurl.err", 100*time.Millisecond, err)
		if err != nil {
			return err
		}
		return res
	})

	d4 := fx.EmittableFunc(func() interface{} {
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

	myObserver := observer.New(onNext, onError, onDone)

	myStream := Start(d1, d3, d4, e1, d2)

	sub := myStream.Subscribe(myObserver)
	s := <-sub

	assert.Exactly(t, []int{200, 301, 500}, responseCodes)
	assert.Equal(t, "Bad URL", s.Err().Error())

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

	done := myStream.Subscribe(nf)
	<-done

	assert.Equal(t, 6, mynum)
}

func TestSubscribeToErrFunc(t *testing.T) {
	myStream := Just(1, "hello", errors.New("bang"), 43.5)

	var myerr error

	ef := handlers.ErrFunc(func(err error) {
		myerr = err
	})

	done := myStream.Subscribe(ef)
	sub := <-done

	assert.Equal(t, "bang", myerr.Error())
	assert.Equal(t, "bang", sub.Error.Error())
}

func TestSubscribeToDoneFunc(t *testing.T) {
	myStream := Just(nil)

	donetext := ""

	df := handlers.DoneFunc(func() {
		donetext = "done"
	})

	done := myStream.Subscribe(df)
	<-done

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

	ob := observer.New(onNext, onError, onDone)

	done := myStream.Subscribe(ob)
	sub := <-done

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

	assert.Equal("bang", sub.Err().Error())
}

func TestObservableMap(t *testing.T) {
	items := []interface{}{1, 2, 3, "foo", "bar", []byte("baz")}
	it, err := iterable.New(items)
	if err != nil {
		t.Fail()
	}

	stream1 := From(it)

	multiplyAllIntBy := func(factor interface{}) fx.MappableFunc {
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

	sub := stream2.Subscribe(onNext)
	<-sub

	assert.Exactly(t, []int{10, 20, 30}, nums)
}

func TestObservableCount(t *testing.T) {
	items := []interface{}{1, 2, 3, "foo", "bar", errors.New("error")}
	it, err := iterable.New(items)
	if err != nil {
		t.Fail()
	}

	stream := From(it)
	count := stream.Count()
	total := <-count

	assert.Exactly(t, int64(6), total)
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

	sub := stream2.Subscribe(onNext)
	<-sub

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

	sub := stream2.Subscribe(onNext)
	<-sub

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

	sub := stream2.Subscribe(onNext)
	<-sub

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

	stream2 := stream1.Filter(lt(9))

	nums := []int{}
	onNext := handlers.NextFunc(func(item interface{}) {
		if num, ok := item.(int); ok {
			nums = append(nums, num)
		}
	})

	sub := stream2.Subscribe(onNext)
	<-sub

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

	sub := stream2.Subscribe(onNext)
	<-sub

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

	sub := stream2.Subscribe(onNext)
	<-sub

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

	sub := stream2.Subscribe(onNext)
	<-sub

	assert.Exactly(t, []int{3}, nums)
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

	sub := stream2.Subscribe(onNext)
	<-sub

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

	sub := stream2.Subscribe(onNext)
	<-sub

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

	sub := stream2.Subscribe(onNext)
	<-sub

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

	sub := stream2.Subscribe(onNext)
	<-sub

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

	sub := stream2.Subscribe(onNext)
	<-sub

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

	sub := stream2.Subscribe(onNext)
	<-sub

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

	sub := stream2.Subscribe(onNext)
	<-sub

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

	sub := stream2.Subscribe(onNext)
	<-sub

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

	sub := stream2.Subscribe(onNext)
	<-sub

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

	sub := myStream.Subscribe(observer.New(onNext, onDone))
	<-sub

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

	sub := myStream.Subscribe(observer.New(onNext, onDone))
	<-sub

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

	sub := myStream.Subscribe(observer.New(onNext, onDone))
	<-sub

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

	sub := myStream.Subscribe(observer.New(onNext, onDone))
	<-sub

	assert.Exactly(t, []string{"end"}, stringarray)
}

func TestEmptyCompletesSequence(t *testing.T) {
	// given
	emissionObserver := observer.NewObserverMock()

	// and empty sequence
	sequence := Empty()

	// when subscribes to the sequence
	<-sequence.Subscribe(emissionObserver.Capture())

	// then completes without any emission
	emissionObserver.AssertNotCalled(t, "OnNext", mock.Anything)
	emissionObserver.AssertNotCalled(t, "OnError", mock.Anything)
	emissionObserver.AssertCalled(t, "OnDone")
}
