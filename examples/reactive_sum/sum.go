/*
 * Copyright (c) 2016-2017 Joe Chasinga
 *
 * Permission is hereby granted, free of charge, to any person obtaining a copy
 * of this software and associated documentation files (the "Software"), to deal
 * in the Software without restriction, including without limitation the rights
 * to use, copy, modify, merge, publish, distribute, sublicense, and/or sell
 * copies of the Software, and to permit persons to whom the Software is
 * furnished to do so, subject to the following conditions:
 *
 * The above copyright notice and this permission notice shall be included in all
 * copies or substantial portions of the Software.
 *
 * THE SOFTWARE IS PROVIDED "AS IS", WITHOUT WARRANTY OF ANY KIND, EXPRESS OR
 * IMPLIED, INCLUDING BUT NOT LIMITED TO THE WARRANTIES OF MERCHANTABILITY,
 * FITNESS FOR A PARTICULAR PURPOSE AND NONINFRINGEMENT. IN NO EVENT SHALL THE
 * AUTHORS OR COPYRIGHT HOLDERS BE LIABLE FOR ANY CLAIM, DAMAGES OR OTHER
 * LIABILITY, WHETHER IN AN ACTION OF CONTRACT, TORT OR OTHERWISE, ARISING FROM,
 * OUT OF OR IN CONNECTION WITH THE SOFTWARE OR THE USE OR OTHER DEALINGS IN THE
 * SOFTWARE.
 *
 */

package main

import (
	"bufio"
	"fmt"
	"os"
	"strconv"
	"strings"

	"github.com/reactivex/rxgo/handlers"
	"github.com/reactivex/rxgo/observable"
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
