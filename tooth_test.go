package tooth

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

const eq = "Should be equal"

func TestTooth(t *testing.T) {
	grumpy := New("grumpy")
	happy := New("happy")
	messy := New("messy")

	happy.Subscribe(grumpy)
	messy.Subscribe(grumpy)
	messy.Subscribe(happy)
	grumpy.Subscribe(happy)

	grumpy.Publish("Hello everyone")
	happy.Publish("I hate life!")

	msg1 := messy.Fetch(grumpy)
	msg2 := happy.Fetch(grumpy)
	msg3 := grumpy.Fetch(happy)

	messages := messy.FetchAll()
	expected := []string{"Hello everyone", "I hate life!"}

	assert.Equal(t, msg1, "Hello everyone", eq)
	assert.Equal(t, msg2, "Hello everyone", eq)
	assert.Equal(t, msg3, "I hate life!", eq)
	assert.NotEmpty(t, messages)
	assert.ObjectsAreEqual(expected, messages)
}
