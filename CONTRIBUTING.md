# Contributions Guidelines

Contributions are always welcome. However, to make this a smooth collaboration experience for everyone and to maintain the quality of the code, here is a few things to consider before and after making a pull request:

## Consistency

There are already +80 operators and +250 unit tests. Please don't try to reinvent the wheel and make sure to check first how the current implementation solve the most common problems.

## Edge Case

When we develop a new operator, there are a lot of edge cases to handle (eager/lazy observation, sequential vs parallel, on error strategy, etc.). The utility functions `observable()`, `single()` and `optionalSingle()` are there to help. Yet, it is not always possible to use them (observation of multiple Observables, etc.). In this case, you may want to take a look at existing operators like `WindowWithTime` to see an exhaustive implementation.

## Unit Tests

Make sure to include unit tests. Again, consistency is key. In most of the unit tests we use the RxGo assertion API.

## Duration

If an operator input contains a duration, we should use `rxgo.Duration`. It allows us to mock it and to implement deterministic tests whenever possible using `timeCausality()`.

## Code Formatting

Before to create a pull request, make sure to format your code using:

* [gofumpt](https://github.com/mvdan/gofumpt)
    * Install: `go get mvdan.cc/gofumpt`
    * Execution: `gofumpt -s -w .`

* [goimports](https://godoc.org/golang.org/x/tools/cmd/goimports)
    * Install: `go get golang.o`
    * Execution: `goimports -w .`