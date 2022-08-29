package rxgo

import (
	"log"
	"testing"
)

func TestNever(t *testing.T) {
	NEVER[any]().SubscribeSync(func(a any) {}, func(err error) {}, func() {
		log.Println("Completed")
	})
}
