package rxgo

import (
	"context"
	"errors"
	"net/http"
	"strconv"
	"sync"
	"testing"
	"time"

	"github.com/reactivex/rxgo/options"

	"github.com/reactivex/rxgo/handlers"
	"github.com/stretchr/testify/assert"
)

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
	myStream := Just(1, 3.1416, &struct{ foo string }{"bar"})
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

		res, err := fakeGet("somehost.com", 0*time.Millisecond, result)
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

		res, err := fakeGet("somehost.com", 50*time.Millisecond, result)
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

		res, err := fakeGet("badserver.co", 100*time.Millisecond, result)
		if err != nil {
			return err
		}
		return res
	})

	e1 := Supplier(func() interface{} {
		err := errors.New("Bad URL")
		res, err := fakeGet("badurl.err", 150*time.Millisecond, err)
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

// FIXME
//func TestSubscribeToErrFunc(t *testing.T) {
//	myStream := Just(1, "hello", errors.New("bang"), 43.5)
//
//	var myerr error
//
//	ef := handlers.ErrFunc(func(err error) {
//		myerr = err
//	})
//
//	sub := myStream.Subscribe(ef).Block()
//
//	assert.Equal(t, "bang", myerr.Error())
//	assert.Equal(t, "bang", sub.Error())
//}

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

	myStream := Just("foo", "bar", "baz", 'a', 'b', errors.New("bang"), 99)

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

	AssertObservable(t, stream, HasItems(10, 20, 30))
}

func TestObservableTake(t *testing.T) {
	stream := Just(1, 2, 3, 4, 5).Take(3)
	AssertObservable(t, stream, HasItems(1, 2, 3))
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
	stream1 := Just(1, 2, 3, 4, 5)
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

func TestObservableTakeLastLessThanNth(t *testing.T) {
	stream1 := Just(4, 5)
	stream2 := stream1.TakeLast(3)

	nums := []int{}
	onNext := handlers.NextFunc(func(item interface{}) {
		if num, ok := item.(int); ok {
			nums = append(nums, num)
		}
	})

	stream2.Subscribe(onNext).Block()

	assert.Exactly(t, []int{4, 5}, nums)
}

func TestObservableTakeLastLessThanNth2(t *testing.T) {
	stream1 := Just(4, 5)
	stream2 := stream1.TakeLast(100000)

	nums := []int{}
	onNext := handlers.NextFunc(func(item interface{}) {
		if num, ok := item.(int); ok {
			nums = append(nums, num)
		}
	})

	stream2.Subscribe(onNext).Block()

	assert.Exactly(t, []int{4, 5}, nums)
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
	stream1 := Just(1, 2, 3, 120, []byte("baz"), 7, 10, 13)

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
	stream1 := Just(0, 1, 3)
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
	stream1 := Just(0, 1, 3)

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

// FIXME Data race
//func TestParallelSubscribeToObserver(t *testing.T) {
//	assert := assert.New(t)
//	myStream := Just("foo", "bar", "baz", 'a', 'b', 99)
//
//	var wordsCount uint64
//	var charsCount uint64
//	var integersCount uint64
//	finished := false
//
//	onNext := handlers.NextFunc(func(item interface{}) {
//		switch item.(type) {
//		case string:
//			atomic.AddUint64(&wordsCount, 1)
//		case rune:
//			atomic.AddUint64(&charsCount, 1)
//		case int:
//			atomic.AddUint64(&integersCount, 1)
//		}
//	})
//
//	onError := handlers.ErrFunc(func(err error) {
//		t.Logf("Error emitted in the stream: %v\n", err)
//	})
//
//	onDone := handlers.DoneFunc(func() {
//		finished = true
//	})
//
//	ob := NewObserver(onNext, onError, onDone)
//
//	myStream.Subscribe(ob, options.WithParallelism(2)).Block()
//
//	assert.True(finished)
//
//	assert.Equal(integersCount, uint64(0x1))
//	assert.Equal(wordsCount, uint64(0x3))
//	assert.Equal(charsCount, uint64(0x2))
//}

// FIXME Data race
//func TestParallelSubscribeToObserverWithError(t *testing.T) {
//	assert := assert.New(t)
//
//	myStream := Just("foo", "bar", "baz", 'a', 'b', 99, errors.New("error"))
//
//	finished := false
//
//	onNext := handlers.NextFunc(func(item interface{}) {
//	})
//
//	onError := handlers.ErrFunc(func(err error) {
//		t.Logf("Error emitted in the stream: %v\n", err)
//	})
//
//	onDone := handlers.DoneFunc(func() {
//		finished = true
//	})
//
//	ob := NewObserver(onNext, onError, onDone)
//
//	myStream.Subscribe(ob, options.WithParallelism(2)).Block()
//
//	assert.False(finished)
//}

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
	stream1 := Just(0, 1, 3, 5, 1, 8)

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
	stream1 := Just(0, 1, 3, 5, 1, 8)

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
	stream1 := Just(1, 2, 2, 1, 3)

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
	stream1 := Just(1, 2, 2, 1, 3)

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

func TestObservableSequenceEqualWithCorrectSequence(t *testing.T) {
	sequence := Just(2, 5, 12, 43, 98, 100, 213)
	result := Just(2, 5, 12, 43, 98, 100, 213).SequenceEqual(sequence)
	AssertSingle(t, result, HasValue(true))
}

func TestObservableSequenceEqualWithIncorrectSequence(t *testing.T) {
	sequence := Just(2, 5, 12, 43, 98, 100, 213)
	result := Just(2, 5, 12, 43, 15, 100, 213).SequenceEqual(sequence)
	AssertSingle(t, result, HasValue(false))
}

func TestObservableSequenceEqualWithDifferentLengthSequence(t *testing.T) {
	sequenceShorter := Just(2, 5, 12, 43, 98, 100)
	sequenceLonger := Just(2, 5, 12, 43, 98, 100, 213, 512)

	resultForShorter := Just(2, 5, 12, 43, 98, 100, 213).SequenceEqual(sequenceShorter)
	AssertSingle(t, resultForShorter, HasValue(false))

	resultForLonger := Just(2, 5, 12, 43, 98, 100, 213).SequenceEqual(sequenceLonger)
	AssertSingle(t, resultForLonger, HasValue(false))
}

func TestObservableSequenceEqualWithEmpty(t *testing.T) {
	result := Empty().SequenceEqual(Empty())
	AssertSingle(t, result, HasValue(true))
}

func TestObservableScanWithIntegers(t *testing.T) {
	stream1 := Just(0, 1, 3, 5, 1, 8)

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
	stream1 := Just("hello", "world", "this", "is", "foo")

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
	stream1 := Just(1, 2, 3, 4, 5)
	add := func(acc interface{}, elem interface{}) interface{} {
		if a, ok := acc.(int); ok {
			if b, ok := elem.(int); ok {
				return a + b
			}
		}
		return 0
	}

	var got Optional
	_, err := stream1.Reduce(add).Subscribe(handlers.NextFunc(func(i interface{}) {
		got = i.(Optional)
	})).Block()
	if err != nil {
		t.Fail()
	}
	assert.False(t, got.IsEmpty())
	assert.Exactly(t, Of(14), got)
}

func TestObservableReduceEmpty(t *testing.T) {
	add := func(acc interface{}, elem interface{}) interface{} {
		if a, ok := acc.(int); ok {
			if b, ok := elem.(int); ok {
				return a + b
			}
		}
		return 0
	}
	stream := Empty()

	var got Optional
	_, err := stream.Reduce(add).Subscribe(handlers.NextFunc(func(i interface{}) {
		got = i.(Optional)
	})).Block()
	if err != nil {
		t.Fail()
	}
	assert.True(t, got.IsEmpty())
}

func TestObservableReduceNil(t *testing.T) {
	stream := Just(1, 2, 3, 4, 5)
	nilReduce := func(acc interface{}, elem interface{}) interface{} {
		return nil
	}
	var got Optional
	_, err := stream.Reduce(nilReduce).Subscribe(handlers.NextFunc(func(i interface{}) {
		got = i.(Optional)
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
	stream := Just(1, 2, 3, "foo", "bar", errors.New("error"))
	count, err := stream.Count().Subscribe(nil).Block()
	if err != nil {
		t.Fail()
	}
	assert.Exactly(t, int64(6), count)
}

func TestObservableFirstOrDefault(t *testing.T) {
	v, err := Empty().FirstOrDefault(7).Subscribe(nil).Block()
	if err != nil {
		t.Fail()
	}
	assert.Exactly(t, 7, v)
}

func TestObservableFirstOrDefaultWithValue(t *testing.T) {
	v, err := Just(0, 1, 2).FirstOrDefault(7).Subscribe(nil).Block()
	if err != nil {
		t.Fail()
	}
	assert.Exactly(t, 0, v)
}

func TestObservableLastOrDefault(t *testing.T) {
	v, err := Empty().LastOrDefault(7).Subscribe(nil).Block()
	if err != nil {
		t.Fail()
	}
	assert.Exactly(t, 7, v)
}

func TestObservableLastOrDefaultWithValue(t *testing.T) {
	v, err := Just(0, 1, 3).LastOrDefault(7).Subscribe(nil).Block()
	if err != nil {
		t.Fail()
	}
	assert.Exactly(t, 3, v)
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

func TestObservableZip(t *testing.T) {
	stream1 := Just(1, 2, 3)
	stream2 := Just(10, 20, 30)
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
	AssertObservable(t, zip, HasItems(11, 22, 33))
}

func TestObservableZipWithDifferentLength1(t *testing.T) {
	stream1 := Just(1, 2, 3)
	stream2 := Just(10, 20)
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
	AssertObservable(t, zip, HasItems(11, 22))
}

func TestObservableZipWithDifferentLength2(t *testing.T) {
	stream1 := Just(1, 2)
	stream2 := Just(10, 20, 30)
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
	AssertObservable(t, zip, HasItems(11, 22))
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
		i += 2
	})
	df := handlers.DoneFunc(func() {
		i += 5
	})
	ob1 := CheckEventHandlers(nf, df)
	ob1.OnNext("")
	ob1.OnDone()
	assert.Equal(t, 7, i)
}

func TestObservableForEach(t *testing.T) {
	assert := assert.New(t)
	myStream := Just("foo", "bar", "baz", 'a', 'b', errors.New("bang"), 99)
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

func TestAll(t *testing.T) {
	predicateAllInt := func(i interface{}) bool {
		switch i.(type) {
		case int:
			return true
		default:
			return false
		}
	}

	AssertSingle(t, Just(1, 2, 3).All(predicateAllInt),
		HasValue(true), HasNotRaisedAnyError())

	AssertSingle(t, Just(1, "x", 3).All(predicateAllInt),
		HasValue(false), HasNotRaisedAnyError())
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

	AssertObservable(t, stream, HasItems(1, 2, 3))
	assert.Equal(t, 6, sum)
}

func TestDoOnEachWithEmpty(t *testing.T) {
	sum := 0
	stream := Empty().DoOnEach(func(i interface{}) {
		sum += i.(int)
	})

	AssertObservable(t, stream, HasSize(0))
	assert.Equal(t, 0, sum)
}

func TestRepeat(t *testing.T) {
	repeat := Just(1, 2, 3).Repeat(1, nil)
	AssertObservable(t, repeat, HasItems(1, 2, 3, 1, 2, 3))
}

func TestRepeatZeroTimes(t *testing.T) {
	repeat := Just(1, 2, 3).Repeat(0, nil)
	AssertObservable(t, repeat, HasItems(1, 2, 3))
}

func TestRepeatWithNegativeCount(t *testing.T) {
	repeat := Just(1, 2, 3).Repeat(-2, nil)
	AssertObservable(t, repeat, HasItems(1, 2, 3))
}

func TestRepeatWithFrequency(t *testing.T) {
	frequency := new(mockDuration)
	frequency.On("duration").Return(time.Millisecond)

	repeat := Just(1, 2, 3).Repeat(1, frequency)
	AssertObservable(t, repeat, HasItems(1, 2, 3, 1, 2, 3))
	frequency.AssertNumberOfCalls(t, "duration", 1)
	frequency.AssertExpectations(t)
}

func TestAverageInt(t *testing.T) {
	AssertSingle(t, Just(1, 2, 3).AverageInt(), HasValue(2))
	AssertSingle(t, Just(1, 20).AverageInt(), HasValue(10))
	AssertSingle(t, Empty().AverageInt(), HasValue(0))
	AssertSingle(t, Just(1.1, 2.2, 3.3).AverageInt(), HasRaisedAnError())
}

func TestAverageInt8(t *testing.T) {
	AssertSingle(t, Just(int8(1), int8(2), int8(3)).AverageInt8(), HasValue(int8(2)))
	AssertSingle(t, Just(int8(1), int8(20)).AverageInt8(), HasValue(int8(10)))
	AssertSingle(t, Empty().AverageInt8(), HasValue(0))
	AssertSingle(t, Just(1.1, 2.2, 3.3).AverageInt8(), HasRaisedAnError())
}

func TestAverageInt16(t *testing.T) {
	AssertSingle(t, Just(int16(1), int16(2), int16(3)).AverageInt16(), HasValue(int16(2)))
	AssertSingle(t, Just(int16(1), int16(20)).AverageInt16(), HasValue(int16(10)))
	AssertSingle(t, Empty().AverageInt16(), HasValue(0))
	AssertSingle(t, Just(1.1, 2.2, 3.3).AverageInt16(), HasRaisedAnError())
}

func TestAverageInt32(t *testing.T) {
	AssertSingle(t, Just(int32(1), int32(2), int32(3)).AverageInt32(), HasValue(int32(2)))
	AssertSingle(t, Just(int32(1), int32(20)).AverageInt32(), HasValue(int32(10)))
	AssertSingle(t, Empty().AverageInt32(), HasValue(0))
	AssertSingle(t, Just(1.1, 2.2, 3.3).AverageInt32(), HasRaisedAnError())
}

func TestAverageInt64(t *testing.T) {
	AssertSingle(t, Just(int64(1), int64(2), int64(3)).AverageInt64(), HasValue(int64(2)))
	AssertSingle(t, Just(int64(1), int64(20)).AverageInt64(), HasValue(int64(10)))
	AssertSingle(t, Empty().AverageInt64(), HasValue(0))
	AssertSingle(t, Just(1.1, 2.2, 3.3).AverageInt64(), HasRaisedAnError())
}

func TestAverageFloat32(t *testing.T) {
	AssertSingle(t, Just(float32(1), float32(2), float32(3)).AverageFloat32(), HasValue(float32(2)))
	AssertSingle(t, Just(float32(1), float32(20)).AverageFloat32(), HasValue(float32(10.5)))
	AssertSingle(t, Empty().AverageFloat32(), HasValue(0))
	AssertSingle(t, Just("x").AverageFloat32(), HasRaisedAnError())
}

func TestAverageFloat64(t *testing.T) {
	AssertSingle(t, Just(float64(1), float64(2), float64(3)).AverageFloat64(), HasValue(float64(2)))
	AssertSingle(t, Just(float64(1), float64(20)).AverageFloat64(), HasValue(float64(10.5)))
	AssertSingle(t, Empty().AverageFloat64(), HasValue(0))
	AssertSingle(t, Just("x").AverageFloat64(), HasRaisedAnError())
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
	AssertOptionalSingle(t, optionalSingle, HasValue(5))
}

func TestMaxWithEmpty(t *testing.T) {
	comparator := func(interface{}, interface{}) Comparison {
		return Equals
	}

	optionalSingle := Empty().Max(comparator)
	AssertOptionalSingle(t, optionalSingle, IsEmpty())
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
	AssertOptionalSingle(t, optionalSingle, HasValue(1))
}

func TestMinWithEmpty(t *testing.T) {
	comparator := func(interface{}, interface{}) Comparison {
		return Equals
	}

	optionalSingle := Empty().Min(comparator)
	AssertOptionalSingle(t, optionalSingle, IsEmpty())
}

func TestBufferWithCountWithCountAndSkipEqual(t *testing.T) {
	obs := Just(1, 2, 3, 4, 5, 6).BufferWithCount(3, 3)
	AssertObservable(t, obs, HasItems([]interface{}{1, 2, 3}, []interface{}{4, 5, 6}))
}

func TestBufferWithCountWithCountAndSkipNotEqual(t *testing.T) {
	obs := Just(1, 2, 3, 4, 5, 6).BufferWithCount(2, 3)
	AssertObservable(t, obs, HasItems([]interface{}{1, 2}, []interface{}{4, 5}))
}

func TestBufferWithCountWithEmpty(t *testing.T) {
	obs := Empty().BufferWithCount(2, 3)
	AssertObservable(t, obs, IsEmpty())
}

func TestBufferWithCountWithIncompleteLastItem(t *testing.T) {
	obs := Just(1, 2, 3, 4).BufferWithCount(2, 3)
	AssertObservable(t, obs, HasItems([]interface{}{1, 2}, []interface{}{4}))
}

func TestBufferWithCountWithError(t *testing.T) {
	obs := Just(1, 2, 3, 4, errors.New("")).BufferWithCount(3, 3)
	AssertObservable(t, obs, HasItems([]interface{}{1, 2, 3}, []interface{}{4}))
	AssertObservable(t, obs, HasRaisedError(errors.New("")))
}

func TestBufferWithInvalidInputs(t *testing.T) {
	// TODO HasRaisedAnErrorWithCause
	obs := Just(1, 2, 3, 4).BufferWithCount(0, 5)
	AssertObservable(t, obs, HasRaisedAnError())

	obs = Just(1, 2, 3, 4).BufferWithCount(5, 0)
	AssertObservable(t, obs, HasRaisedAnError())
}

func TestBufferWithTimeWithMockedTime(t *testing.T) {
	just := Just(1, 2, 3)

	timespan := new(mockDuration)
	timespan.On("duration").Return(10 * time.Second)

	timeshift := new(mockDuration)
	timeshift.On("duration").Return(10 * time.Second)

	obs := just.BufferWithTime(timespan, timeshift)

	AssertObservable(t, obs, HasItems([]interface{}{1, 2, 3}))
	timespan.AssertCalled(t, "duration")
	timeshift.AssertNotCalled(t, "duration")
}

func TestBufferWithTimeWithMinorMockedTime(t *testing.T) {
	ch := make(chan interface{})
	from := FromIterator(newIteratorFromChannel(ch))

	timespan := new(mockDuration)
	timespan.On("duration").Return(1 * time.Millisecond)

	timeshift := new(mockDuration)
	timeshift.On("duration").Return(1 * time.Millisecond)

	obs := from.BufferWithTime(timespan, timeshift)

	time.Sleep(10 * time.Millisecond)
	ch <- 1
	close(ch)

	obs.Subscribe(nil).Block()

	timespan.AssertCalled(t, "duration")
}

func TestBufferWithTimeWithIllegalInput(t *testing.T) {
	AssertObservable(t, Empty().BufferWithTime(nil, nil), HasRaisedAnError())
	AssertObservable(t, Empty().BufferWithTime(WithDuration(0*time.Second), nil), HasRaisedAnError())
}

func TestBufferWithTimeWithNilTimeshift(t *testing.T) {
	just := Just(1, 2, 3)
	obs := just.BufferWithTime(WithDuration(1*time.Second), nil)
	AssertObservable(t, obs, IsNotEmpty())
}

func TestBufferWithTimeWithError(t *testing.T) {
	just := Just(1, 2, 3, errors.New(""))
	obs := just.BufferWithTime(WithDuration(1*time.Second), nil)
	AssertObservable(t, obs, HasItems([]interface{}{1, 2, 3}), HasRaisedAnError())
}

func TestBufferWithTimeWithEmpty(t *testing.T) {
	obs := Empty().BufferWithTime(WithDuration(1*time.Second), WithDuration(1*time.Second))
	AssertObservable(t, obs, IsEmpty())
}

func TestBufferWithTimeOrCountWithInvalidInputs(t *testing.T) {
	obs := Empty().BufferWithTimeOrCount(nil, 5)
	AssertObservable(t, obs, HasRaisedAnError())

	obs = Empty().BufferWithTimeOrCount(WithDuration(0), 5)
	AssertObservable(t, obs, HasRaisedAnError())

	obs = Empty().BufferWithTimeOrCount(WithDuration(time.Millisecond*5), 0)
	AssertObservable(t, obs, HasRaisedAnError())
}

func TestBufferWithTimeOrCountWithCount(t *testing.T) {
	just := Just(1, 2, 3)
	obs := just.BufferWithTimeOrCount(WithDuration(1*time.Second), 2)
	AssertObservable(t, obs, HasItems([]interface{}{1, 2}, []interface{}{3}))
}

func TestBufferWithTimeOrCountWithTime(t *testing.T) {
	ch := make(chan interface{})
	from := FromIterator(newIteratorFromChannel(ch))

	got := make([]interface{}, 0)

	obs := from.BufferWithTimeOrCount(WithDuration(1*time.Millisecond), 10).
		Subscribe(handlers.NextFunc(func(i interface{}) {
			got = append(got, i)
		}))
	ch <- 1
	time.Sleep(50 * time.Millisecond)
	ch <- 2
	close(ch)
	obs.Block()

	// Check items are included
	got1 := false
	got2 := false
	for _, v := range got {
		switch v := v.(type) {
		case []interface{}:
			if len(v) == 1 {
				if v[0] == 1 {
					got1 = true
				} else if v[0] == 2 {
					got2 = true
				}
			}
		}
	}

	assert.True(t, got1)
	assert.True(t, got2)
}

func TestBufferWithTimeOrCountWithMockedTime(t *testing.T) {
	ch := make(chan interface{})
	from := FromIterator(newIteratorFromChannel(ch))

	timespan := new(mockDuration)
	timespan.On("duration").Return(1 * time.Millisecond)

	obs := from.BufferWithTimeOrCount(timespan, 5)

	time.Sleep(50 * time.Millisecond)
	ch <- 1
	close(ch)

	obs.Subscribe(nil).Block()

	timespan.AssertCalled(t, "duration")
}

func TestBufferWithTimeOrCountWithError(t *testing.T) {
	just := Just(1, 2, 3, errors.New(""), 4)
	obs := just.BufferWithTimeOrCount(WithDuration(10*time.Second), 2)
	AssertObservable(t, obs, HasItems([]interface{}{1, 2}, []interface{}{3}),
		HasRaisedAnError())
}

func TestSumInt64(t *testing.T) {
	AssertSingle(t, Just(1, 2, 3).SumInt64(), HasValue(int64(6)))
	AssertSingle(t, Just(int8(1), int(2), int16(3), int32(4), int64(5)).SumInt64(),
		HasValue(int64(15)))
	AssertSingle(t, Just(1.1, 2.2, 3.3).SumInt64(), HasRaisedAnError())
	AssertSingle(t, Empty().SumInt64(), HasValue(int64(0)))
}

func TestSumFloat32(t *testing.T) {
	AssertSingle(t, Just(float32(1.0), float32(2.0), float32(3.0)).SumFloat32(),
		HasValue(float32(6.)))
	AssertSingle(t, Just(float32(1.1), 2, int8(3), int16(1), int32(1), int64(1)).SumFloat32(),
		HasValue(float32(9.1)))
	AssertSingle(t, Just(1.1, 2.2, 3.3).SumFloat32(), HasRaisedAnError())
	AssertSingle(t, Empty().SumFloat32(), HasValue(float32(0)))
}

func TestSumFloat64(t *testing.T) {
	AssertSingle(t, Just(1.1, 2.2, 3.3).SumFloat64(),
		HasValue(6.6))
	AssertSingle(t, Just(float32(1.0), 2, int8(3), 4., int16(1), int32(1), int64(1)).SumFloat64(),
		HasValue(float64(13.)))
	AssertSingle(t, Just("x").SumFloat64(), HasRaisedAnError())
	AssertSingle(t, Empty().SumFloat64(), HasValue(float64(0)))
}

func TestMapWithTwoSubscription(t *testing.T) {
	just := Just(1).Map(func(i interface{}) interface{} {
		return 1 + i.(int)
	}).Map(func(i interface{}) interface{} {
		return 1 + i.(int)
	})

	AssertObservable(t, just, HasItems(3))
	AssertObservable(t, just, HasItems(3))
}

func TestMapWithConcurrentSubscriptions(t *testing.T) {
	just := Just(1).Map(func(i interface{}) interface{} {
		return 1 + i.(int)
	}).Map(func(i interface{}) interface{} {
		return 1 + i.(int)
	})

	wg := sync.WaitGroup{}
	for i := 0; i < 100; i++ {
		wg.Add(1)
		go func() {
			defer wg.Done()
			AssertObservable(t, just, HasItems(3))
		}()
	}

	wg.Wait()
}

func TestStartWithItems(t *testing.T) {
	obs := Just(1, 2, 3).StartWithItems(10, 20)
	AssertObservable(t, obs, HasItems(10, 20, 1, 2, 3))
}

func TestIgnoreElements(t *testing.T) {
	obs := Just(1, 2, 3).IgnoreElements()
	AssertObservable(t, obs, IsEmpty())
}

func TestIgnoreElements_WithError(t *testing.T) {
	err := errors.New("")
	obs := Just(1, err, 3).IgnoreElements()
	AssertObservable(t, obs, HasRaisedError(err))
}

func TestOnErrorReturn(t *testing.T) {
	obs := Just(1, 2, errors.New("3"), 4).OnErrorReturn(func(e error) interface{} {
		i, _ := strconv.Atoi(e.Error())
		return i
	})
	AssertObservable(t, obs, HasItems(1, 2, 3, 4))
}

func TestOnErrorReturnItem(t *testing.T) {
	obs := Just(1, 2, errors.New("3"), 4).OnErrorReturnItem(3)
	AssertObservable(t, obs, HasItems(1, 2, 3, 4))
}

func TestOnErrorResumeNext(t *testing.T) {
	err := errors.New("8")
	obs := Just(1, 2, errors.New("3"), 4).OnErrorResumeNext(func(e error) Observable {
		return Just(5, 6, err, 9)
	})
	AssertObservable(t, obs, HasItems(1, 2, 5, 6), HasRaisedError(err))
}

func TestToSlice(t *testing.T) {
	single := Just(1, 2, 3).ToSlice()
	AssertSingle(t, single, HasValue([]interface{}{1, 2, 3}))
}

func TestToSlice_Empty(t *testing.T) {
	single := Empty().ToSlice()
	AssertSingle(t, single, HasValue([]interface{}{}))
}

func TestToMap(t *testing.T) {
	single := Just(3, 4, 5, true, false).ToMap(func(i interface{}) interface{} {
		switch v := i.(type) {
		case int:
			return v
		case bool:
			if v {
				return 0
			}
			return 1
		default:
			return i
		}
	})
	AssertSingle(t, single, HasValue(map[interface{}]interface{}{
		3: 3,
		4: 4,
		5: 5,
		0: true,
		1: false,
	}))
}

func TestToMap_Empty(t *testing.T) {
	single := Empty().ToMap(func(i interface{}) interface{} {
		return i
	})
	AssertSingle(t, single, HasValue(map[interface{}]interface{}{}))
}

func TestToMapWithValueSelector(t *testing.T) {
	keySelector := func(i interface{}) interface{} {
		switch v := i.(type) {
		case int:
			return v
		case bool:
			if v {
				return 0
			}
			return 1
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
	single := Just(3, 4, 5, true, false).ToMapWithValueSelector(keySelector, valueSelector)
	AssertSingle(t, single, HasValue(map[interface{}]interface{}{
		3: 30,
		4: 40,
		5: 50,
		0: true,
		1: false,
	}))
}

func TestToMapWithValueSelector_Empty(t *testing.T) {
	single := Empty().ToMapWithValueSelector(func(i interface{}) interface{} {
		return i
	}, func(i interface{}) interface{} {
		return i
	})
	AssertSingle(t, single, HasValue(map[interface{}]interface{}{}))
}

func channel(ch chan interface{}, t time.Duration) (item interface{}, closed bool, cancelled bool) {
	ctx, cancel := context.WithTimeout(context.Background(), t)
	defer cancel()

	select {
	case i, ok := <-ch:
		if ok {
			return i, false, false
		} else {
			return nil, true, false
		}
	case <-ctx.Done():
		return nil, false, true
	}
}

func TestToChannel(t *testing.T) {
	ch := Just(1, 2, 3).ToChannel()
	item, _, _ := channel(ch, wait)
	assert.Equal(t, 1, item)
	item, _, _ = channel(ch, wait)
	assert.Equal(t, 2, item)
	item, _, _ = channel(ch, wait)
	assert.Equal(t, 3, item)
	_, closed, _ := channel(ch, wait)
	assert.True(t, closed)
}

func TestToChannel_BufferedChannel(t *testing.T) {
	ch := Just(1, 2, 3).ToChannel(options.WithBufferedChannel(3))
	item, _, _ := channel(ch, wait)
	assert.Equal(t, 1, item)
	item, _, _ = channel(ch, wait)
	assert.Equal(t, 2, item)
	item, _, _ = channel(ch, wait)
	assert.Equal(t, 3, item)
	_, closed, _ := channel(ch, wait)
	assert.True(t, closed)
}

func TestTakeWhile(t *testing.T) {
	obs := Just(1, 2, 3, 4, 5).TakeWhile(func(item interface{}) bool {
		return item != 3
	})
	AssertObservable(t, obs, HasItems(1, 2))
}

func TestTakeWhile_Empty(t *testing.T) {
	obs := Empty().TakeWhile(func(item interface{}) bool {
		return item != 3
	})
	AssertObservable(t, obs, IsEmpty())
}

func TestTakeUntil(t *testing.T) {
	obs := Just(1, 2, 3, 4, 5).TakeUntil(func(item interface{}) bool {
		return item == 3
	})
	AssertObservable(t, obs, HasItems(1, 2, 3))
}

func TestTakeUntil_Empty(t *testing.T) {
	obs := Empty().TakeUntil(func(item interface{}) bool {
		return item == 3
	})
	AssertObservable(t, obs, IsEmpty())
}

func TestStartWith(t *testing.T) {
	obs := Just(1, 2, 3).StartWithItems(10, 20)
	AssertObservable(t, obs, HasItems(10, 20, 1, 2, 3))
}

func TestStartWith_WithError(t *testing.T) {
	err := errors.New("")
	obs := Just(1, 2, 3).StartWithItems(10, err)
	AssertObservable(t, obs, HasItems(10), HasRaisedError(err))
}

func TestStartWith_Empty(t *testing.T) {
	obs := Empty().StartWithItems(10, 20)
	AssertObservable(t, obs, HasItems(10, 20))
}

func TestStartWithIterable(t *testing.T) {
	ch := make(chan interface{}, 2)
	obs := Just(1, 2, 3).StartWithIterable(newIterableFromChannel(ch))
	ch <- 10
	ch <- 20
	close(ch)
	AssertObservable(t, obs, HasItems(10, 20, 1, 2, 3))
}

func TestStartWithIterable_WithError(t *testing.T) {
	err := errors.New("")
	ch := make(chan interface{}, 2)
	obs := Just(1, 2, 3).StartWithIterable(newIterableFromChannel(ch))
	ch <- 10
	ch <- err
	close(ch)
	AssertObservable(t, obs, HasItems(10), HasRaisedError(err))
}

func TestStartWithIterable_Empty(t *testing.T) {
	ch := make(chan interface{}, 2)
	obs := Empty().StartWithIterable(newIterableFromChannel(ch))
	ch <- 10
	ch <- 20
	close(ch)
	AssertObservable(t, obs, HasItems(10, 20))
}

func TestStartWithIterable_WithoutItems(t *testing.T) {
	ch := make(chan interface{})
	obs := Just(1, 2, 3).StartWithIterable(newIterableFromChannel(ch))
	close(ch)
	AssertObservable(t, obs, HasItems(1, 2, 3))
}

func TestStartWithObservable(t *testing.T) {
	obs := Just(1, 2, 3).StartWithObservable(Just(10, 20))
	AssertObservable(t, obs, HasItems(10, 20, 1, 2, 3))
}

func TestStartWithObservable_WithError(t *testing.T) {
	err := errors.New("")
	obs := Just(1, 2, 3).StartWithObservable(Just(10, err))
	AssertObservable(t, obs, HasItems(10), HasRaisedError(err))
}

func TestStartWithObservable_Empty1(t *testing.T) {
	obs := Empty().StartWithObservable(Just(10, 20))
	AssertObservable(t, obs, HasItems(10, 20))
}

func TestStartWithObservable_Empty2(t *testing.T) {
	obs := Just(1, 2, 3).StartWithObservable(Empty())
	AssertObservable(t, obs, HasItems(1, 2, 3))
}

//var _ = Describe("Timeout operator", func() {
// FIXME
//Context("when creating an observable with timeout operator", func() {
//	ch := make(chan interface{}, 10)
//	duration := WithDuration(pollingInterval)
//	o := FromChannel(ch).Timeout(duration)
//	Context("after a given period without items", func() {
//		outNext, outErr, _ := subscribe(o)
//
//		ch <- 1
//		ch <- 2
//		ch <- 3
//		time.Sleep(time.Second)
//		ch <- 4
//		It("should receive the elements before the timeout", func() {
//			Expect(pollItems(outNext, timeout)).Should(Equal([]interface{}{1, 2, 3}))
//		})
//		It("should receive a TimeoutError", func() {
//			Expect(pollItem(outErr, timeout)).Should(Equal(&TimeoutError{}))
//		})
//	})
//})
//})
