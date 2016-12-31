package observable

import (
	"errors"
	"time"
	//"math/rand"
	//"net/http"
	"testing"
	//"time"

	//"github.com/bmizerany/assert"
	"github.com/jochasinga/grx/bases"
	"github.com/jochasinga/grx/emittable"
	//"github.com/jochasinga/grx/handlers"
	//"github.com/jochasinga/grx/iterable"
	"github.com/jochasinga/grx/observer"
	"github.com/stretchr/testify/assert"
)

type Fixture struct {
	num                      int
	text                     string
	char                     rune
	err                      error
	fin                      bool
	errch                    chan error
	eint, estr, echar, echan bases.Emitter
}

func (f *Fixture) Setup(conf func(*Fixture)) *Fixture {
	conf(f)
	return f
}

func setupDefaultFixture() *Fixture {
	return (&Fixture{}).Setup(func(f *Fixture) {
		f.errch = make(chan error, 1)
		f.eint = emittable.From(10)
		f.estr = emittable.From("hello")
		f.echar = emittable.From('a')
		f.echan = emittable.From(f.errch)
	})
}

func TestBasicObservableConstructor(t *testing.T) {
	assert := assert.New(t)
	basic := NewBasic(0)

	assert.IsType((Basic)(nil), basic)
	assert.Equal(0, cap(basic))

	basic = NewBasic(3)
	assert.Equal(3, cap(basic))
}

func TestConnectableObservableConstructor(t *testing.T) {
	assert := assert.New(t)
	text := "hello"
	connectable := NewConnectable(0)

	if assert.IsType(Connectable{}, connectable) {
		assert.Equal(0, cap(connectable.emitters))
	}

	connectable = NewConnectable(3)
	assert.Equal(3, cap(connectable.emitters))
	ob := observer.Observer{
		NextHandler: func(item bases.Item) {
			text += item.(string)
		},
	}

	connectable = NewConnectable(6, ob)
	assert.Equal(6, cap(connectable.emitters))
	connectable.observers[0].NextHandler(bases.Item(" world"))
	assert.Equal("hello world", text)
}

func TestBasicSubscription(t *testing.T) {
	fixture := setupDefaultFixture()

	// Send an error over to errch
	go func() {
		fixture.errch <- errors.New("yike")
		return
	}()

	basic := BasicFrom([]bases.Emitter{
		fixture.eint,
		fixture.estr,
		fixture.echar,
		fixture.echan,
	})

	ob := observer.Observer{
		NextHandler: func(it bases.Item) {
			switch it := it.(type) {
			case int:
				t.Logf("Item is an integer: %d\n", it)
				fixture.num += it
			case string:
				t.Logf("Item is a string: %q\n", it)
				fixture.text += it
			case rune:
				t.Logf("Item is a rune: %v\n", it)
				fixture.char += it
			case chan error:
				if e, ok := <-it; ok {
					t.Logf("Item is an emitted error: %v", e)
					fixture.err = e
				}
			}
		},

		DoneHandler: func() {
			t.Log("done")
			fixture.fin = !fixture.fin
		},
	}
	done := basic.Subscribe(ob)
	<-done

	subtests := []struct {
		n, expected interface{}
	}{
		{fixture.num, 10},
		{fixture.text, "hello"},
		{fixture.char, 'a'},
		{fixture.err, errors.New("yike")},
		{fixture.fin, true},
	}

	for _, tt := range subtests {
		assert.Equal(t, tt.expected, tt.n)
	}
}

func TestBasicMap(t *testing.T) {
	fixture := setupDefaultFixture()

	sourceSlice := []bases.Emitter{
		fixture.eint,
		fixture.estr,
		fixture.echar,
		fixture.echan,
	}

	basic := BasicFrom(sourceSlice)

	multiplyAllIntBy := func(n interface{}) func(bases.Emitter) bases.Emitter {
		return func(e bases.Emitter) bases.Emitter {
			if item, err := e.Emit(); err == nil {
				if val, ok := item.(int); ok {
					return emittable.From(val * n.(int))
				}
			}
			return e
		}
	}

	basic = basic.Map(multiplyAllIntBy(100))

	subtests := []bases.Emitter{
		emittable.From(1000),
		fixture.estr,
		fixture.echar,
		fixture.echan,
	}

	i := 0
	for e := range basic {
		assert.Equal(t, subtests[i], e)
		i++
	}
}

func TestBasicFilter(t *testing.T) {
	fixture := setupDefaultFixture()

	sourceSlice := []bases.Emitter{
		fixture.eint,
		fixture.estr,
		fixture.echar,
		fixture.echan,
	}

	basic := BasicFrom(sourceSlice)

	isIntOrString := func(e bases.Emitter) bool {
		if item, err := e.Emit(); err == nil {
			switch item.(type) {
			case int, string:
				return true
			}
		}
		return false
	}

	basic = basic.Filter(isIntOrString)

	assert.Equal(t, fixture.eint, <-basic)
	assert.Equal(t, fixture.estr, <-basic)
}

func TestBasicEmpty(t *testing.T) {

	basic := Empty()
	isDone := false

	done := basic.Subscribe(observer.Observer{
		DoneHandler: func() {
			isDone = !isDone
		},
	})

	<-done
	assert.True(t, isDone)
}

func TestBasicInterval(t *testing.T) {

	fin := make(chan struct{})
	basic := Interval(500*time.Millisecond, fin)

	_ = basic.Subscribe(observer.Observer{
		NextHandler: func(item bases.Item) {
			t.Log(item)
		},
		DoneHandler: func() {
			t.Log("done")
		},
	})

	<-time.After(2 * time.Second)
	fin <- struct{}{}
}

func TestConnectableMap(t *testing.T) {
	fixture := setupDefaultFixture()

	sourceSlice := []bases.Emitter{
		fixture.eint,
		fixture.estr,
		fixture.echar,
		fixture.echan,
	}

	connectable := ConnectableFrom(sourceSlice)

	// multiplyAllIntBy is a CurryableFunc
	multiplyAllIntBy := func(n interface{}) func(bases.Emitter) bases.Emitter {
		return func(e bases.Emitter) bases.Emitter {
			if item, err := e.Emit(); err == nil {
				if val, ok := item.(int); ok {
					return emittable.From(val * n.(int))
				}
			}
			return e
		}
	}

	connectable = connectable.Map(multiplyAllIntBy(100))

	subtests := []bases.Emitter{
		emittable.From(1000),
		fixture.estr,
		fixture.echar,
		fixture.echan,
	}

	i := 0
	for e := range connectable.emitters {
		assert.Equal(t, subtests[i], e)
		i++
	}
}

func TestConnectableSubscription(t *testing.T) {

	fixture := (&Fixture{}).Setup(func(f *Fixture) {
		f.errch = make(chan error, 1)
		f.eint = emittable.From(10)
		f.estr = emittable.From("hello")
		f.echar = emittable.From('a')
		f.echan = emittable.From(f.errch)
	})

	// Send an error over to errch
	go func() {
		fixture.errch <- errors.New("yike")
		return
	}()

	connectable := ConnectableFrom([]bases.Emitter{
		fixture.eint,
		fixture.estr,
		fixture.echar,
		fixture.echan,
	})

	ob1 := observer.Observer{
		NextHandler: func(it bases.Item) {
			switch it := it.(type) {
			case int:
				t.Logf("Item is an integer: %d\n", it)
				fixture.num += it
			case string:
				t.Logf("Item is a string: %q\n", it)
				fixture.text += it
			case rune:
				t.Logf("Item is a rune: %v\n", it)
				fixture.char += it
			case chan error:
				if e, ok := <-it; ok {
					t.Logf("Item is an emitted error: %v", e)
					fixture.err = e
				}
			}
		},
		DoneHandler: func() {
			t.Log("done")
			fixture.fin = !fixture.fin
		},
	}

	ob2 := observer.Observer{
		NextHandler: func(it bases.Item) {
			switch it := it.(type) {
			case int:
				t.Logf("Item is indeed an integer: %d\n", it)
			case string:
				t.Logf("Item is indeed a string: %q\n", it)
			case rune:
				t.Logf("Item is indeed a rune: %v\n", it)
			case chan error:
				if e, ok := <-it; ok {
					t.Logf("Item is an indeed emitted error: %v", e)
					fixture.err = e
				}
			}
		},
		DoneHandler: func() {
			t.Log("Indeed it's done")
		},
	}

	beforetests := []struct {
		n, expected interface{}
	}{
		{fixture.num, 0},
		{fixture.text, ""},
		{fixture.char, rune(0)},
		{fixture.err, error(nil)},
		{fixture.fin, false},
	}

	for _, tt := range beforetests {
		assert.Equal(t, tt.expected, tt.n)
	}

	connectable = connectable.Subscribe(ob1).Subscribe(ob2)

	done := connectable.Connect()
	<-done

	subtests := []struct {
		n, expected interface{}
	}{
		{fixture.num, 10},
		{fixture.text, "hello"},
		{fixture.char, 'a'},
		{fixture.err, errors.New("yike")},
		{fixture.fin, true},
	}

	for _, tt := range subtests {
		assert.Equal(t, tt.expected, tt.n)
	}
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

func TestStartMethodWithFakeExternalCalls(t *testing.T) {
	fakeHttpResponses := []*http.Response{}

	// NOTE: HTTP Response errors such as status 500 does not return an error
	fakeHttpErrors := []error{}

	// Fake directives that returns an Event containing an HTTP response.
	d1 := func() bases.Emitter {
		res := &http.Response{
			Status:     "404 NOT FOUND",
			StatusCode: 404,
			Proto:      "HTTP/1.0",
			ProtoMajor: 1,
		}

		// Simulating an I/O block
		time.Sleep(20 * time.Millisecond)
		return emittable.From(res)
	}

	d2 := func() bases.Emitter {
		res := &http.Response{
			Status:     "200 OK",
			StatusCode: 200,
			Proto:      "HTTP/1.0",
			ProtoMajor: 1,
		}
		time.Sleep(10 * time.Millisecond)
		return emittable.From(res)
	}

	d3 := func() bases.Emitter {
		res := &http.Response{
			Status:     "500 SERVER ERROR",
			StatusCode: 500,
			Proto:      "HTTP/1.0",
			ProtoMajor: 1,
		}
		time.Sleep(30 * time.Millisecond)
		return emittable.From(res)
	}

	d4 := func() bases.Emitter {
		err := errors.New("Some kind of error")
		time.Sleep(50 * time.Millisecond)
		return emittable.From(err)
	}

	watcher := &observer.Observer{
		NextHandler: handlers.NextFunc(func(it bases.Item) {
			if res, ok := it.(*http.Response); ok {
				fakeHttpResponses = append(fakeHttpResponses, res)
			}
		}),
		ErrHandler: handlers.ErrFunc(func(err error) {
			fakeHttpErrors = append(fakeHttpErrors, err)
		}),
		DoneHandler: handlers.DoneFunc(func() {
			fakeHttpResponses = append(fakeHttpResponses, &http.Response{
				Status:     "999 End",
				StatusCode: 999,
			})
		}),
	}

	source := Start(d1, d2, d3, d4)
	_, err := source.Subscribe(watcher)

	assert := assert.New(t)
	assert.Nil(err)

	<-time.After(100 * time.Millisecond)

	assert.IsType((*Observable)(nil), source)
	assert.Equal(4, len(fakeHttpResponses))
	assert.Equal(1, len(fakeHttpErrors))
	assert.Equal(200, fakeHttpResponses[0].StatusCode)
	assert.Equal(404, fakeHttpResponses[1].StatusCode)
	assert.Equal(500, fakeHttpResponses[2].StatusCode)
	assert.Equal(999, fakeHttpResponses[3].StatusCode)
	assert.Equal("Some kind of error", fakeHttpErrors[0])
}

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
