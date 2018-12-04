package rxgo

import (
	"errors"
	"net/http"
	"strconv"
	"sync/atomic"
	"testing"
	"time"

	"github.com/reactivex/rxgo/handlers"
	"github.com/reactivex/rxgo/iterable"
	"github.com/reactivex/rxgo/optional"
	"github.com/reactivex/rxgo/options"
	"github.com/stretchr/testify/assert"
)

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

func TestRange(t *testing.T) {
	got := []interface{}{}
	r, err := Range(1, 5)
	if err != nil {
		t.Fail()
	}
	r.Subscribe(handlers.NextFunc(func(i interface{}) {
		got = append(got, i)
	})).Block()
	assert.Equal(t, []interface{}{1, 2, 3, 4, 5}, got)
}

func TestRangeWithNegativeCount(t *testing.T) {
	r, err := Range(1, -5)
	assert.NotNil(t, err)
	assert.Nil(t, r)
}

func TestRangeWithMaximumExceeded(t *testing.T) {
	r, err := Range(1<<31, 1)
	assert.NotNil(t, err)
	assert.Nil(t, r)
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

	d1 := Supplier(func() interface{} {
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

	d2 := Supplier(func() interface{} {
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

	d3 := Supplier(func() interface{} {
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

	e1 := Supplier(func() interface{} {
		err := errors.New("Bad URL")
		res, err := fakeGet("badurl.err", 100*time.Millisecond, err)
		if err != nil {
			return err
		}
		return res
	})

	d4 := Supplier(func() interface{} {
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
	stream := Just(1, 2, 3).Map(func(i interface{}) interface{} {
		return i.(int) * 10
	})

	AssertThatObservable(t, stream, HasItems(10, 20, 30))
}

func TestObservableTake(t *testing.T) {
	stream := Just(1, 2, 3, 4, 5).Take(3)
	AssertThatObservable(t, stream, HasItems(1, 2, 3))
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

	lt := func(target interface{}) Predicate {
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

	myStream.Subscribe(ob, options.WithParallelism(2)).Block()

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

	myStream.Subscribe(ob, options.WithParallelism(2)).Block()

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
	g, err := got.Get()
	assert.Nil(t, err)
	assert.Nil(t, g)
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

func TestObservableSkipWhile(t *testing.T) {
	got := []interface{}{}
	Just(1, 2, 3, 4, 5).SkipWhile(func(i interface{}) bool {
		switch i := i.(type) {
		case int:
			return i != 3
		default:
			return true
		}
	}).Subscribe(handlers.NextFunc(func(i interface{}) {
		got = append(got, i)
	})).Block()

	assert.Equal(t, []interface{}{3, 4, 5}, got)
}

func TestObservableSkipWhileWithEmpty(t *testing.T) {
	got := []interface{}{}
	Empty().SkipWhile(func(i interface{}) bool {
		switch i := i.(type) {
		case int:
			return i != 3
		default:
			return false
		}
	}).Subscribe(handlers.NextFunc(func(i interface{}) {
		got = append(got, i)
	})).Block()

	assert.Equal(t, []interface{}{}, got)
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

func TestObservableToMap(t *testing.T) {
	items := []interface{}{3, 4, 5, true, false}
	it, err := iterable.New(items)
	if err != nil {
		t.Fail()
	}
	stream1 := From(it)
	stream2 := stream1.ToMap(func(i interface{}) interface{} {
		switch v := i.(type) {
		case int:
			return v
		case bool:
			if v {
				return 0
			} else {
				return 1
			}
		default:
			return i
		}
	})
	var got interface{}
	stream2.Subscribe(handlers.NextFunc(func(i interface{}) {
		got = i
	})).Block()
	expected := map[interface{}]interface{}{
		3: 3,
		4: 4,
		5: 5,
		0: true,
		1: false,
	}
	assert.Exactly(t, expected, got)
}

func TestObservableToMapWithEmpty(t *testing.T) {
	stream1 := Empty()
	stream2 := stream1.ToMap(func(i interface{}) interface{} {
		return i
	})
	var got interface{}
	stream2.Subscribe(handlers.NextFunc(func(i interface{}) {
		got = i
	})).Block()
	assert.Exactly(t, map[interface{}]interface{}{}, got)
}

func TestObservableToMapWithValueSelector(t *testing.T) {
	items := []interface{}{3, 4, 5, true, false}
	it, err := iterable.New(items)
	if err != nil {
		t.Fail()
	}
	stream1 := From(it)
	keySelector := func(i interface{}) interface{} {
		switch v := i.(type) {
		case int:
			return v
		case bool:
			if v {
				return 0
			} else {
				return 1
			}
		default:
			return i
		}
	}
	valueSelector := func(i interface{}) interface{} {
		switch v := i.(type) {
		case int:
			return v * 10
		case bool:
			return v
		default:
			return i
		}
	}
	stream2 := stream1.ToMapWithValueSelector(keySelector, valueSelector)
	var got interface{}
	stream2.Subscribe(handlers.NextFunc(func(i interface{}) {
		got = i
	})).Block()
	expected := map[interface{}]interface{}{
		3: 30,
		4: 40,
		5: 50,
		0: true,
		1: false,
	}
	assert.Exactly(t, expected, got)
}

func TestObservableToMapWithValueSelectorWithEmpty(t *testing.T) {
	stream1 := Empty()
	f := func(i interface{}) interface{} {
		return i
	}
	stream2 := stream1.ToMapWithValueSelector(f, f)
	var got interface{}
	stream2.Subscribe(handlers.NextFunc(func(i interface{}) {
		got = i
	})).Block()
	assert.Exactly(t, map[interface{}]interface{}{}, got)
}

func TestObservableZip(t *testing.T) {
	items := []interface{}{1, 2, 3}
	it, err := iterable.New(items)
	if err != nil {
		t.Fail()
	}
	stream1 := From(it)
	items2 := []interface{}{10, 20, 30}
	it2, err := iterable.New(items2)
	if err != nil {
		t.Fail()
	}
	stream2 := From(it2)
	zipper := func(elem1 interface{}, elem2 interface{}) interface{} {
		switch v1 := elem1.(type) {
		case int:
			switch v2 := elem2.(type) {
			case int:
				return v1 + v2
			}
		}
		return 0
	}
	zip := stream1.ZipFromObservable(stream2, zipper)
	nums := []int{}
	onNext := handlers.NextFunc(func(item interface{}) {
		if num, ok := item.(int); ok {
			nums = append(nums, num)
		}
	})
	zip.Subscribe(onNext).Block()
	assert.Exactly(t, []int{11, 22, 33}, nums)
}

func TestObservableZipWithDifferentLength1(t *testing.T) {
	items := []interface{}{1, 2, 3}
	it, err := iterable.New(items)
	if err != nil {
		t.Fail()
	}
	stream1 := From(it)
	items2 := []interface{}{10, 20}
	it2, err := iterable.New(items2)
	if err != nil {
		t.Fail()
	}
	stream2 := From(it2)
	zipper := func(elem1 interface{}, elem2 interface{}) interface{} {
		switch v1 := elem1.(type) {
		case int:
			switch v2 := elem2.(type) {
			case int:
				return v1 + v2
			}
		}
		return 0
	}
	zip := stream1.ZipFromObservable(stream2, zipper)
	nums := []int{}
	onNext := handlers.NextFunc(func(item interface{}) {
		if num, ok := item.(int); ok {
			nums = append(nums, num)
		}
	})
	zip.Subscribe(onNext).Block()
	assert.Exactly(t, []int{11, 22}, nums)
}

func TestObservableZipWithDifferentLength2(t *testing.T) {
	items := []interface{}{1, 2}
	it, err := iterable.New(items)
	if err != nil {
		t.Fail()
	}
	stream1 := From(it)
	items2 := []interface{}{10, 20, 30}
	it2, err := iterable.New(items2)
	if err != nil {
		t.Fail()
	}
	stream2 := From(it2)
	zipper := func(elem1 interface{}, elem2 interface{}) interface{} {
		switch v1 := elem1.(type) {
		case int:
			switch v2 := elem2.(type) {
			case int:
				return v1 + v2
			}
		}
		return 0
	}
	zip := stream1.ZipFromObservable(stream2, zipper)
	nums := []int{}
	onNext := handlers.NextFunc(func(item interface{}) {
		if num, ok := item.(int); ok {
			nums = append(nums, num)
		}
	})
	zip.Subscribe(onNext).Block()
	assert.Exactly(t, []int{11, 22}, nums)
}

func TestObservableZipWithEmpty(t *testing.T) {
	stream1 := Empty()
	stream2 := Empty()
	zipper := func(elem1 interface{}, elem2 interface{}) interface{} {
		return 0
	}
	zip := stream1.ZipFromObservable(stream2, zipper)
	nums := []int{}
	onNext := handlers.NextFunc(func(item interface{}) {
		if num, ok := item.(int); ok {
			nums = append(nums, num)
		}
	})
	zip.Subscribe(onNext).Block()
	assert.Exactly(t, []int{}, nums)
}

func TestCheckEventHandlers(t *testing.T) {
	i := 0
	nf := handlers.NextFunc(func(interface{}) {
		i = i + 2
	})
	df := handlers.DoneFunc(func() {
		i = i + 5
	})
	ob1 := CheckEventHandlers(nf, df)
	ob1.OnNext("")
	ob1.OnDone()
	assert.Equal(t, 7, i)
}

func TestObservableForEach(t *testing.T) {
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
	sub := myStream.ForEach(onNext, onError, onDone).Block()
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

func TestOnErrorReturn(t *testing.T) {
	got := make([]int, 0)

	obs := Just(1, 2, 3, errors.New("7"), 4).
		OnErrorReturn(func(e error) interface{} {
			i, err := strconv.Atoi(e.Error())
			if err != nil {
				t.Fail()
			}
			return i
		})

	obs.Subscribe(handlers.NextFunc(func(i interface{}) {
		got = append(got, i.(int))
	})).Block()

	assert.Equal(t, []int{1, 2, 3, 7}, got)
}

func TestOnErrorResumeNext(t *testing.T) {
	got := make([]int, 0)

	obs := Just(1, 2, 3, errors.New("7"), 4).
		OnErrorResumeNext(func(e error) Observable {
			return Just(5, 6, errors.New("8"), 9)
		})

	err := obs.Subscribe(handlers.NextFunc(func(i interface{}) {
		got = append(got, i.(int))
	})).Block()

	assert.Equal(t, []int{1, 2, 3, 5, 6}, got)
	assert.NotNil(t, err)
}

func TestAll(t *testing.T) {
	predicateAllInt := func(i interface{}) bool {
		switch i.(type) {
		case int:
			return true
		default:
			return false
		}
	}

	got1, err := Just(1, 2, 3).All(predicateAllInt).
		Subscribe(nil).Block()
	assert.Nil(t, err)
	assert.Equal(t, true, got1)

	got2, err := Just(1, "x", 3).All(predicateAllInt).
		Subscribe(nil).Block()
	assert.Nil(t, err)
	assert.Equal(t, false, got2)
}

func TestContain(t *testing.T) {
	predicate := func(i interface{}) bool {
		switch i := i.(type) {
		case int:
			return i == 2
		default:
			return false
		}
	}

	var got1, got2 bool

	Just(1, 2, 3).Contains(predicate).
		Subscribe(handlers.NextFunc(func(i interface{}) {
			got1 = i.(bool)
		})).Block()
	assert.True(t, got1)

	Just(1, 5, 3).Contains(predicate).
		Subscribe(handlers.NextFunc(func(i interface{}) {
			got2 = i.(bool)
		})).Block()
	assert.False(t, got2)
}

func TestDefaultIfEmpty(t *testing.T) {
	got1 := 0
	Empty().DefaultIfEmpty(3).Subscribe(handlers.NextFunc(func(i interface{}) {
		got1 = i.(int)
	})).Block()
	assert.Equal(t, 3, got1)
}

func TestDefaultIfEmptyWithNonEmpty(t *testing.T) {
	got1 := 0
	Just(1).DefaultIfEmpty(3).Subscribe(handlers.NextFunc(func(i interface{}) {
		got1 = i.(int)
	})).Block()
	assert.Equal(t, 1, got1)
}

func TestDoOnEach(t *testing.T) {
	sum := 0
	stream := Just(1, 2, 3).DoOnEach(func(i interface{}) {
		sum += i.(int)
	})

	AssertThatObservable(t, stream, HasItems(1, 2, 3))
	assert.Equal(t, 6, sum)
}

func TestDoOnEachWithEmpty(t *testing.T) {
	sum := 0
	stream := Empty().DoOnEach(func(i interface{}) {
		sum += i.(int)
	})

	AssertThatObservable(t, stream, HasSize(0))
	assert.Equal(t, 0, sum)
}

func TestRepeat(t *testing.T) {
	repeat := Just(1, 2, 3).Repeat(1, nil)
	AssertThatObservable(t, repeat, HasItems(1, 2, 3, 1, 2, 3))
}

func TestRepeatZeroTimes(t *testing.T) {
	repeat := Just(1, 2, 3).Repeat(0, nil)
	AssertThatObservable(t, repeat, HasItems(1, 2, 3))
}

func TestRepeatWithNegativeCount(t *testing.T) {
	repeat := Just(1, 2, 3).Repeat(-2, nil)
	AssertThatObservable(t, repeat, HasItems(1, 2, 3))
}

func TestRepeatWithFrequency(t *testing.T) {
	frequency := new(mockDuration)
	frequency.On("duration").Return(time.Millisecond)

	repeat := Just(1, 2, 3).Repeat(1, frequency)
	AssertThatObservable(t, repeat, HasItems(1, 2, 3, 1, 2, 3))
	frequency.AssertNumberOfCalls(t, "duration", 1)
	frequency.AssertExpectations(t)
}

func TestAverageInt(t *testing.T) {
	AssertThatSingle(t, Just(1, 2, 3).AverageInt(), HasValue(2))
	AssertThatSingle(t, Just(1, 20).AverageInt(), HasValue(10))
	AssertThatSingle(t, Empty().AverageInt(), HasValue(0))
	AssertThatSingle(t, Just(1.1, 2.2, 3.3).AverageInt(), HasRaisedAnError())
}

func TestAverageInt8(t *testing.T) {
	AssertThatSingle(t, Just(int8(1), int8(2), int8(3)).AverageInt8(), HasValue(int8(2)))
	AssertThatSingle(t, Just(int8(1), int8(20)).AverageInt8(), HasValue(int8(10)))
	AssertThatSingle(t, Empty().AverageInt8(), HasValue(0))
	AssertThatSingle(t, Just(1.1, 2.2, 3.3).AverageInt8(), HasRaisedAnError())
}

func TestAverageInt16(t *testing.T) {
	AssertThatSingle(t, Just(int16(1), int16(2), int16(3)).AverageInt16(), HasValue(int16(2)))
	AssertThatSingle(t, Just(int16(1), int16(20)).AverageInt16(), HasValue(int16(10)))
	AssertThatSingle(t, Empty().AverageInt16(), HasValue(0))
	AssertThatSingle(t, Just(1.1, 2.2, 3.3).AverageInt16(), HasRaisedAnError())
}

func TestAverageInt32(t *testing.T) {
	AssertThatSingle(t, Just(int32(1), int32(2), int32(3)).AverageInt32(), HasValue(int32(2)))
	AssertThatSingle(t, Just(int32(1), int32(20)).AverageInt32(), HasValue(int32(10)))
	AssertThatSingle(t, Empty().AverageInt32(), HasValue(0))
	AssertThatSingle(t, Just(1.1, 2.2, 3.3).AverageInt32(), HasRaisedAnError())
}

func TestAverageInt64(t *testing.T) {
	AssertThatSingle(t, Just(int64(1), int64(2), int64(3)).AverageInt64(), HasValue(int64(2)))
	AssertThatSingle(t, Just(int64(1), int64(20)).AverageInt64(), HasValue(int64(10)))
	AssertThatSingle(t, Empty().AverageInt64(), HasValue(0))
	AssertThatSingle(t, Just(1.1, 2.2, 3.3).AverageInt64(), HasRaisedAnError())
}

func TestAverageFloat32(t *testing.T) {
	AssertThatSingle(t, Just(float32(1), float32(2), float32(3)).AverageFloat32(), HasValue(float32(2)))
	AssertThatSingle(t, Just(float32(1), float32(20)).AverageFloat32(), HasValue(float32(10.5)))
	AssertThatSingle(t, Empty().AverageFloat32(), HasValue(0))
	AssertThatSingle(t, Just("x").AverageFloat32(), HasRaisedAnError())
}

func TestAverageFloat64(t *testing.T) {
	AssertThatSingle(t, Just(float64(1), float64(2), float64(3)).AverageFloat64(), HasValue(float64(2)))
	AssertThatSingle(t, Just(float64(1), float64(20)).AverageFloat64(), HasValue(float64(10.5)))
	AssertThatSingle(t, Empty().AverageFloat64(), HasValue(0))
	AssertThatSingle(t, Just("x").AverageFloat64(), HasRaisedAnError())
}

func TestMax(t *testing.T) {
	comparator := func(e1, e2 interface{}) Comparison {
		i1 := e1.(int)
		i2 := e2.(int)
		if i1 > i2 {
			return Greater
		} else if i1 < i2 {
			return Smaller
		} else {
			return Equals
		}
	}

	optionalSingle := Just(1, 5, 1).Max(comparator)
	AssertThatOptionalSingle(t, optionalSingle, HasValue(5))
}

func TestMaxWithEmpty(t *testing.T) {
	comparator := func(interface{}, interface{}) Comparison {
		return Equals
	}

	optionalSingle := Empty().Max(comparator)
	AssertThatOptionalSingle(t, optionalSingle, IsEmpty())
}

func TestMin(t *testing.T) {
	comparator := func(e1, e2 interface{}) Comparison {
		i1 := e1.(int)
		i2 := e2.(int)
		if i1 > i2 {
			return Greater
		} else if i1 < i2 {
			return Smaller
		} else {
			return Equals
		}
	}

	optionalSingle := Just(5, 1, 5).Min(comparator)
	AssertThatOptionalSingle(t, optionalSingle, HasValue(1))
}

func TestMinWithEmpty(t *testing.T) {
	comparator := func(interface{}, interface{}) Comparison {
		return Equals
	}

	optionalSingle := Empty().Min(comparator)
	AssertThatOptionalSingle(t, optionalSingle, IsEmpty())
}

func TestSumInt64(t *testing.T) {
	AssertThatSingle(t, Just(1, 2, 3).SumInt64(), HasValue(int64(6)))
	AssertThatSingle(t, Just(int8(1), int(2), int16(3), int32(4), int64(5)).SumInt64(),
		HasValue(int64(15)))
	AssertThatSingle(t, Just(1.1, 2.2, 3.3).SumInt64(), HasRaisedAnError())
	AssertThatSingle(t, Empty().SumInt64(), HasValue(int64(0)))
}

func TestSumFloat32(t *testing.T) {
	AssertThatSingle(t, Just(float32(1.0), float32(2.0), float32(3.0)).SumFloat32(),
		HasValue(float32(6.)))
	AssertThatSingle(t, Just(float32(1.1), 2, int8(3), int16(1), int32(1), int64(1)).SumFloat32(),
		HasValue(float32(9.1)))
	AssertThatSingle(t, Just(1.1, 2.2, 3.3).SumFloat32(), HasRaisedAnError())
	AssertThatSingle(t, Empty().SumFloat32(), HasValue(float32(0)))
}

func TestSumFloat64(t *testing.T) {
	AssertThatSingle(t, Just(1.1, 2.2, 3.3).SumFloat64(),
		HasValue(6.6))
	AssertThatSingle(t, Just(float32(1.0), 2, int8(3), 4., int16(1), int32(1), int64(1)).SumFloat64(),
		HasValue(float64(13.)))
	AssertThatSingle(t, Just("x").SumFloat64(), HasRaisedAnError())
	AssertThatSingle(t, Empty().SumFloat64(), HasValue(float64(0)))
}
