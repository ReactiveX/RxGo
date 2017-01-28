package observable

import (

	//"fmt"

	//"math/rand"
	//"net/http"
	"errors"
	"testing"
	"time"

	"github.com/jochasinga/grx/bases"
	//"github.com/jochasinga/grx/emittable"

	"github.com/jochasinga/grx/handlers"
	//"github.com/jochasinga/grx/iterable"

	"github.com/jochasinga/grx/observer"
	"github.com/stretchr/testify/assert"
	//"github.com/stretchr/testify/suite"
)

func TestObservableImplementsBaseObservable(t *testing.T) {
	t.Skip("Skipping implementation test for now")
	assert.Implements(t, (*bases.Observable)(nil), Observable(nil))
}

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

	ob1 := checkEventHandler(myObserver)
	ob2 := checkEventHandler(df)

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
	numItems := 5
	yes := make(chan struct{})

	n := 0
	go func() {
		for sig := range yes {
			assert.Equal(t, struct{}{}, sig)
			n++
		}
	}()

	onNext := handlers.NextFunc(func(item interface{}) {
		yes <- struct{}{}
	})

	sub := myStream.Subscribe(onNext)
	<-sub

	assert.Equal(t, numItems, n)
}

func TestFromOperator(t *testing.T) {
	items := []interface{}{1, 3.1416, &struct{ foo string }{"bar"}}
	myStream := From(items)
	lenItems := len(items)
	yes := make(chan struct{})

	n := 0
	go func() {
		for sig := range yes {
			assert.Equal(t, struct{}{}, sig)
			n++
		}
	}()

	onNext := handlers.NextFunc(func(item interface{}) {
		yes <- struct{}{}
	})

	sub := myStream.Subscribe(onNext)
	<-sub

	assert.Equal(t, lenItems, n)
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

	myStream := From([]interface{}{
		"foo", "bar", "baz", 'a', 'b', errors.New("bang"), 99,
	})

	strings := make(chan string)
	chars := make(chan rune)
	integers := make(chan int)
	fin := make(chan struct{})

	onnext := handlers.NextFunc(func(item interface{}) {
		switch item := item.(type) {
		case string:
			strings <- item
		case rune:
			chars <- item
		case int:
			integers <- item
		}
	})

	onerr := handlers.ErrFunc(func(err error) {
		t.Logf("Error emitted in the stream: %v\n", err)
	})

	ondone := handlers.DoneFunc(func() {
		fin <- struct{}{}
	})

	ob := observer.New(onnext, onerr, ondone)

	go func() {
		expected := []string{"foo", "bar", "baz"}
		n := 0
		for string := range strings {
			assert.Equal(expected[n], string)
			n++
		}
	}()

	go func() {
		expected := []rune{'a', 'b', 'c'}
		n := 0
		for char := range chars {
			assert.Equal(expected[n], char)
			n++
		}
	}()

	go func() {
		num := <-integers
		assert.Equal(0, num, "integers should not receive anything.")
	}()

	go func() {
		sig := <-fin
		assert.Equal(nil, sig, "fin should not receive anything.")
	}()

	done := myStream.Subscribe(ob)
	sub := <-done

	assert.Equal("bang", sub.Error.Error())
}

/*
func (suite *BasicSuite) TestSubscription() {

	// Send an error over to errch
	go func() {
		suite.fixture.errchan <- errors.New("yike")
		return
	}()

	bs := From([]bases.Emitter{
		suite.fixture.eint,
		suite.fixture.etext,
		suite.fixture.echar,
		suite.fixture.echan,
	})

	ob := observer.Observer{
		NextHandler: func(it bases.Item) {
			switch it := it.(type) {
			case int:
				suite.fixture.num += it
			case string:
				suite.fixture.text += it
			case rune:
				suite.fixture.char += it
			case chan error:
				if e, ok := <-it; ok {
					suite.fixture.err = e
				}
			}
		},

		DoneHandler: func() {
			suite.fixture.isdone = !suite.fixture.isdone
		},
	}
	done := bs.Subscribe(ob)
	<-done

	subtests := []struct {
		n, expected interface{}
	}{
		{suite.fixture.num, 10},
		{suite.fixture.text, "hello"},
		{suite.fixture.char, 'a'},
		{suite.fixture.err, errors.New("yike")},
		{suite.fixture.isdone, true},
	}

	for _, tt := range subtests {
		assert.Equal(suite.T(), tt.expected, tt.n)
	}
}

func (suite *BasicSuite) TestBasicMap() {

	bs1 := From(suite.fixture.emitters)

	multiplyAllIntBy := func(n interface{}) fx.MappableFunc {
		return func(e bases.Emitter) bases.Emitter {
			if item, err := e.Emit(); err == nil {
				if val, ok := item.(int); ok {
					return emittable.From(val * n.(int))
				}
			}
			return e
		}
	}

	bs2 := bs1.Map(multiplyAllIntBy(100))

	subtests := []bases.Emitter{
		emittable.From(1000),
		suite.fixture.etext,
		suite.fixture.echar,
		suite.fixture.echan,
	}

	i := 0
	for e := range bs2 {
		assert.Equal(suite.T(), subtests[i], e)
		i++
	}
}

func (suite *BasicSuite) TestBasicFilter() {

	bs1 := From(suite.fixture.emitters)

	isIntOrString := func(e bases.Emitter) bool {
		if item, err := e.Emit(); err == nil {
			switch item.(type) {
			case int, string:
				return true
			}
		}
		return false
	}

	bs2 := bs1.Filter(isIntOrString)

	assert.Equal(suite.T(), suite.fixture.eint, <-bs2)
	assert.Equal(suite.T(), suite.fixture.etext, <-bs2)
	assert.Nil(suite.T(), <-bs2)
}

func (suite *BasicSuite) TestEmpty() {

	bs := Empty()
	finished := false

	done := bs.Subscribe(observer.Observer{
		DoneHandler: func() {
			finished = !finished
		},
	})

	<-done
	assert.True(suite.T(), finished)
}

func (suite *BasicSuite) TestInterval() {

	numch := make(chan int, 1)
	term := make(chan struct{}, 1)
	source := Interval(term, 1*time.Second)
	assert.IsType(suite.T(), (Basic)(nil), source)

	_ = source.Subscribe(observer.Observer{
		NextHandler: func(it bases.Item) {
			if num, ok := it.(int); ok {
				numch <- num
			}
		},
	})

	i := 0
	go func() {
		for {
			select {
			case num := <-numch:
				assert.Equal(suite.T(), i, num)
			}
			i++
		}
	}()
	<-time.After(5 * time.Second)
	term <- struct{}{}
}

type FakeHttp struct {
	responses []bases.Item
	errors    []error
}

func (fhttp *FakeHttp) Get(delay time.Duration, fn func() (*http.Response, error)) (*http.Response, error) {
	time.Sleep(delay)
	return fn()
}

func (suite *BasicSuite) TestFakeBlockingExternalCalls() {

	fhttp := new(FakeHttp)

	// Fake directives that returns an Event containing an HTTP response.
	d1 := func() &http.Response {
		return &http.Response{
			Status:     "404 NOT FOUND",
			StatusCode: 404,
			Proto:      "HTTP/1.0",
			ProtoMajor: 1,
		}
	}

	d2 := func() &http.Response {
		return &http.Response{
			Status:     "200 OK",
			StatusCode: 200,
			Proto:      "HTTP/1.0",
			ProtoMajor: 1,
		}
	}

	d3 := func() &http.Response {
		return &http.Response{
			Status:     "500 SERVER ERROR",
			StatusCode: 500,
			Proto:      "HTTP/1.0",
			ProtoMajor: 1,
		}
	}

	d4 := func() &http.Response {
		return errors.New("Some kind of error")
	}

	e1 := emittable.From(func() interface{} {
		res, _ := fhttp.Get(10 * time.Millisecond, d1)
		return res
	})
	e2 := emitable.From(func() interface{} {
		res, _ := fhttp.Get(30 * time.Millisecond, d2)
		return res
	})
	e3 := emittable.From(func() interface{} {
		res, _ := fhttp.Get(20 * time.Millisecond, d3)
		return res
	})
	e4 := emittable.From(func() interface{} {
		_, err := fhttp.Get(50 * time.Millisecond, d4)
		return err
	})

	basic := Just(e1, e2, e3, e4)

	watcher := &observer.Observer{
		NextHandler: func(it bases.Item) {
			fhttp.responses = append(fhttp.responses, it)
		}),
		ErrHandler: func(err error) {
			fhttp.errors = append(fhttp.errors, err)
		}),
		DoneHandler: handlers.DoneFunc(func() {
			fhttp.responses = append(fhttp.responses, bases.Item("Oho end")
		}),
	}

	done := basic.Subscribe(watcher)
	<-done

	if assert.NotEmpty(suite.T(), fhttp.responses) && assert.NotEmpty(suite.T(), fhttp.errors) {

		assert.Len(suite.T(), 3, fhttp.responses)
		assert.Len(suite.T(), 1, fhttp.errors)
	}
}

func TestBasicSuite(t *testing.T) {
	suite.Run(t, new(BasicSuite))
}

/*
func TestObservableImplementStream(t *testing.T) {
	assert.Implements(t, (*bases.Stream)(nil), DefaultObservable)
}

func TestObservableImplementIterator(t *testing.T) {
	assert.Implements(t, (*bases.Iterator)(nil), DefaultObservable)
}

func TestCreateObservableWithConstructor(t *testing.T) {
	source := New()
	assert.IsType(t, (*Observable)(nil), source)
}

func TestCreateOperator(t *testing.T) {
	source := Create(func(ob *observer.Observer) {
		ob.OnNext("Hello")
	})

	empty := ""

	eventHandler := handlers.NextFunc(func(it bases.Item) {
		if text, ok := it.(string); ok {
			empty += text
		} else {
			panic("Item is not a string")
		}
	})

	_, _ = source.Subscribe(eventHandler)
	<-time.After(100 * time.Millisecond)
	assert.Equal(t, "Hello", empty)

	source = Create(func(ob *observer.Observer) {
		ob.OnError(errors.New("OMG this is an error"))
	})

	errText := ""
	errHandler := handlers.ErrFunc(func(err error) {
		errText += err.Error()
	})
	_, _ = source.Subscribe(errHandler)
	<-time.After(100 * time.Millisecond)
	assert.Equal(t, "OMG this is an error", errText)
}

func TestEmptyOperator(t *testing.T) {
	msg := "Sumpin's"
	source := Empty()

	watcher := &observer.Observer{
		NextHandler: handlers.NextFunc(func(i bases.Item) {
			panic("NextHandler shouldn't be called")
		}),
		DoneHandler: handlers.DoneFunc(func() {
			msg += " brewin'"
		}),
	}
	_, err := source.Subscribe(watcher)
	assert.Nil(t, err)

	<-time.After(100 * time.Millisecond)
	assert.Equal(t, "Sumpin's brewin'", msg)
}

func TestJustOperator(t *testing.T) {
	assert := assert.New(t)
	url := "http://api.com/api/v1.0/user"
	source := Just(url)

	assert.IsType((*Observable)(nil), source)

	urlWithQueryString := ""
	queryString := "?id=999"
	expected := url + queryString

	watcher := &observer.Observer{
		NextHandler: handlers.NextFunc(func(it bases.Item) {
			if url, ok := it.(string); ok {
				urlWithQueryString += url
			}
		}),
		DoneHandler: handlers.DoneFunc(func() {
			urlWithQueryString += queryString
		}),
	}

	sub, err := source.Subscribe(watcher)
	assert.Nil(err)
	assert.NotNil(sub)
	<-time.After(10 * time.Millisecond)
	assert.Equal(expected, urlWithQueryString)

	source = Just('R', 'x', 'G', 'o')
	e, err := source.Next()
	assert.IsType((*Observable)(nil), source)
	assert.Nil(err)
	assert.IsType((*emittable.Emittable)(nil), e)
	assert.Implements((*bases.Emitter)(nil), e)
}

func TestFromOperator(t *testing.T) {
	assert := assert.New(t)

	iterableUrls := iterable.From([]interface{}{
		"http://api.com/api/v1.0/user",
		"https://dropbox.com/api/v2.1/get/storage",
		"http://googleapi.com/map",
	})

	responses := []string{}

	request := func(url string) string {
		randomNum := rand.Intn(100)
		time.Sleep((200 - time.Duration(randomNum)) * time.Millisecond)
		return fmt.Sprintf("{\"url\":%q}", url)
	}

	urlStream := From(iterableUrls)
	urlObserver := &observer.Observer{
		NextHandler: handlers.NextFunc(func(it bases.Item) {
			if url, ok := it.(string); ok {
				res := request(url)
				responses = append(responses, res)
			} else {
				assert.Fail("Item is not a string as expected")
			}
		}),
		DoneHandler: handlers.DoneFunc(func() {
			responses = append(responses, "END")
		}),
	}

	sub, err := urlStream.Subscribe(urlObserver)
	assert.NotNil(sub)
	assert.Nil(err)

	<-time.After(100 * time.Millisecond)
	expectedStrings := []string{
		"{\"url\":\"http://api.com/api/v1.0/user\"}",
		"{\"url\":\"https://dropbox.com/api/v2.1/get/storage\"}",
		"{\"url\":\"http://googleapi.com/map\"}",
		"END",
	}
	assert.Exactly(expectedStrings, responses)

	iterableNums := iterable.From([]interface{}{1, 2, 3, 4, 5, 6})
	numCopy := []int{}

	numStream := From(iterableNums)
	numObserver := &observer.Observer{
		NextHandler: handlers.NextFunc(func(it bases.Item) {
			if num, ok := it.(int); ok {
				numCopy = append(numCopy, num+1)
			} else {
				assert.Fail("Item is not an integer as expected")
			}
		}),
		DoneHandler: handlers.DoneFunc(func() {
			numCopy = append(numCopy, 0)
		}),
	}

	sub, err = numStream.Subscribe(numObserver)
	assert.NotNil(sub)
	assert.Nil(err)

	<-time.After(100 * time.Millisecond)
	expectedNums := []int{2, 3, 4, 5, 6, 7, 0}
	assert.Exactly(expectedNums, numCopy)
}

func TestStartOperator(t *testing.T) {
	assert := assert.New(t)
	d1 := func() bases.Emitter {
		return emittable.From(333)
	}
	d2 := func() bases.Emitter {
		return emittable.From(666)
	}
	d3 := func() bases.Emitter {
		return emittable.From(999)
	}

	source := Start(d1, d2, d3)
	e, err := source.Next()

	assert.IsType((*Observable)(nil), source)
	assert.Nil(err)
	assert.IsType((*emittable.Emittable)(nil), e)
	assert.Implements((*bases.Emitter)(nil), e)

	nums := []int{}
	watcher := &observer.Observer{
		NextHandler: handlers.NextFunc(func(it bases.Item) {
			if num, ok := it.(int); ok {
				nums = append(nums, num+111)
			} else {
				assert.Fail("Item is not an integer as expected")
			}
		}),
		DoneHandler: handlers.DoneFunc(func() {
			nums = append(nums, 0)
		}),
	}

	_, _ = source.Subscribe(watcher)
	<-time.After(100 * time.Millisecond)
	expected := []int{444, 777, 1110, 0}
	assert.Exactly(expected, nums)
}
*/

/*
func TestIntervalOperator(t *testing.T) {
	assert := assert.New(t)
	numch := make(chan int, 1)
	source := Interval(1 * time.Millisecond)
	assert.IsType((*Observable)(nil), source)

	go func() {
		_, _ = source.Subscribe(&observer.Observer{
			NextHandler: handlers.NextFunc(func(it bases.Item) {
				if num, ok := it.(int); ok {
					numch <- num
				}
			}),
		})
	}()

	i := 0

	select {
	case <-time.After(1 * time.Millisecond):
		if i >= 10 {
			return
		}
		i++
	case num := <-numch:
		assert.Equal(i, num)
	}

	//	for i := 0; i <= 10; i++ {
	//		<-time.After(1 * time.Millisecond)
	//		assert.Equal(i, <-numch)
	//	}

}

func TestRangeOperator(t *testing.T) {
	assert := assert.New(t)
	nums := []int{}
	watcher := &observer.Observer{
		NextHandler: handlers.NextFunc(func(it bases.Item) {
			if num, ok := it.(int); ok {
				nums = append(nums, num)
			} else {
				assert.Fail("Item is not an integer as expected")
			}
		}),
	}
	source := Range(1, 10)
	assert.IsType((*Observable)(nil), source)

	_, err := source.Subscribe(watcher)
	assert.Nil(err)

	<-time.After(100 * time.Millisecond)
	expected := []int{1, 2, 3, 4, 5, 6, 7, 8, 9}
	assert.Exactly(expected, nums)
}

func TestSubscriptionIsNonBlocking(t *testing.T) {
	var (
		s1 = Just("Hello", "world", "foo", 'a', 1.2, -3111.02, []rune{}, struct{}{})
		s2 = From(iterable.From([]interface{}{1, 2, "Hi", 'd', 2.10, -54, []byte{}}))
		s3 = Range(1, 100)
		s4 = Interval(1 * time.Second)
		s5 = Empty()
	)

	sources := []*Observable{s1, s2, s3, s4, s5}

	watcher := &observer.Observer{
		NextHandler: handlers.NextFunc(func(it bases.Item) {
			time.Sleep(1 * time.Second)
			t.Log(it)
			return
		}),
		DoneHandler: handlers.DoneFunc(func() {
			t.Log("DONE")
		}),
	}

	first := time.Now()

	for _, source := range sources {
		_, err := source.Subscribe(watcher)
		assert.Nil(t, err)
	}

	elapsed := time.Since(first)

	comp := assert.Comparison(func() bool {
		return elapsed < 1*time.Second
	})
	assert.Condition(t, comp)
}

func TestObservableDoneChannel(t *testing.T) {
	assert := assert.New(t)
	o := Range(1, 10)
	_, err := o.Subscribe(&observer.Observer{
		NextHandler: handlers.NextFunc(func(it bases.Item) {
			t.Logf("Test value: %v\n", it)
		}),
		DoneHandler: handlers.DoneFunc(func() {
			o.notifier.Done()
		}),
	})

	assert.Nil(err)

	<-time.After(100 * time.Millisecond)
	assert.Equal(struct{}{}, <-o.notifier.IsDone)
}

func TestObservableGetDisposedViaSubscription(t *testing.T) {
	nums := []int{}
	source := Interval(100 * time.Millisecond)
	sub, err := source.Subscribe(&observer.Observer{
		NextHandler: handlers.NextFunc(func(it bases.Item) {
			if num, ok := it.(int); ok {
				nums = append(nums, num)
			} else {
				assert.Fail(t, "Item is not an integer as expected")
			}
		}),
	})

	assert.Nil(t, err)
	assert.NotNil(t, sub)

	<-time.After(300 * time.Millisecond)
	sub.Dispose()

	comp := assert.Comparison(func() bool {
		return len(nums) <= 4
	})
	assert.Condition(t, comp)
}
*/
