package connectable

import (
	"testing"

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

	/*
		co4 := co3.Subscribe(ob)

		assert.Equal(0, cap(co4.Observable))
		assert.Equal("hello", text)

		sub := co4.Connect()
		<-sub

		assert.Equal("helloworld", text)
	*/
}

/*
func (suite *ConnectableSuite) TestSubscription() {

	// Send an error over to errch
	go func() {
		suite.fixture.errchan <- errors.New("yike")
		return
	}()

	co1 := From(suite.fixture.emitters)

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

	beforetests := []struct {
		n, expected interface{}
	}{
		{suite.fixture.num, 0},
		{suite.fixture.text, ""},
		{suite.fixture.char, rune(0)},
		{suite.fixture.err, error(nil)},
		{suite.fixture.isdone, false},
	}

	for _, tt := range beforetests {
		assert.Equal(suite.T(), tt.expected, tt.n)
	}

	co2 := co1.Subscribe(ob)

	done := co2.Connect()
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

func (suite *ConnectableSuite) TestConnectableMap() {

	co1 := From(suite.fixture.emitters)

	// multiplyAllIntBy is a CurryableFunc
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

	co2 := co1.Map(multiplyAllIntBy(100))

	cotests := []bases.Emitter{
		emittable.From(1000),
		suite.fixture.etext,
		suite.fixture.echar,
		suite.fixture.echan,
	}

	i := 0
	for e := range co2.Basic {
		assert.Equal(suite.T(), cotests[i], e)
		i++
	}
}

func (suite *ConnectableSuite) TestFilter() {

	co1 := From(suite.fixture.emitters)

	isIntOrString := func(e bases.Emitter) bool {
		if item, err := e.Emit(); err == nil {
			switch item.(type) {
			case int, string:
				return true
			}
		}
		return false
	}

	co2 := co1.Filter(isIntOrString)

	assert.Equal(suite.T(), suite.fixture.eint, <-co2.Basic)
	assert.Equal(suite.T(), suite.fixture.etext, <-co2.Basic)
	assert.Equal(suite.T(), nil, <-co2.Basic)
}

func TestConnectableSuite(t *testing.T) {
	suite.Run(t, new(ConnectableSuite))
}
*/
