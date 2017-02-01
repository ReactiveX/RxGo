package main

import (
	"bufio"
	"fmt"
	"os"
	"strconv"
	"strings"

	"github.com/jochasinga/grx/handlers"
	"github.com/jochasinga/grx/observable"
)

func main() {

	var num1, num2 int
	reader := bufio.NewReader(os.Stdin)

	processText := func(text, prefix string, numPtr *int) {
		text = strings.Trim(strings.TrimPrefix(text, prefix), " \n")
		*numPtr, _ = strconv.Atoi(text)
	}

	// All side effects are consolidated into this handler.
	onNext := handlers.NextFunc(func(item interface{}) {
		if text, ok := item.(string); ok {
			switch {
			case strings.HasPrefix(text, "a:"):
				processText(text, "a:", &num1)
			case strings.HasPrefix(text, "b:"):
				processText(text, "b:", &num2)
			default:
				fmt.Println("Input does not start with prefix \"a:\" or \"b:\"!")
				return
			}
		}

		fmt.Printf("Running sum: %d\n", num1+num2)

	})

	for {
		fmt.Print("type> ")

		sub := observable.Start(func() interface{} {
			text, err := reader.ReadString('\n')
			if err != nil {
				return err
			}
			return text
		}).Subscribe(onNext)

		<-sub
	}
}
