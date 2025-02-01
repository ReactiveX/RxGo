# teivah/onecontext

[![Go Report Card](https://goreportcard.com/badge/github.com/teivah/onecontext)](https://goreportcard.com/report/github.com/teivah/onecontext)

`teivah/onecontext` is a set of context utilities.

## One Context

Have you ever faced the situation where you have to merge multiple existing contexts?
If not, then you might, eventually.

For example, we can face the situation where we are building an application using a library that gives us a global context (for example, `urfave/cli`).
This context expires once the application is stopped.

Meanwhile, we are exposing a gRPC service like this:

```go
func (f *Foo) Get(ctx context.Context, bar *Bar) (*Baz, error) {
	
}
```

Here, we receive another context provided by gRPC.

Then, in the `Get` implementation, we want to query a database and provide a context for that.

Ideally, we would like to provide a merged context that would expire either:
- When the application is stopped
- Or when the received gRPC context expires

It's is precisely the purpose of this library.

In our case, we can now merge the two contexts in a single one like this:

```go
ctx, cancel := onecontext.Merge(ctx1, ctx2)
```

It returns a merged context that we can now propagate.

## Detach

`onecontext.Detach` returns a context detached from the original cancellation signal. It can be helpful, for example, if we implement an HTTP handler and that we have to pass a context to another goroutine that may expire after we send back a response.

## Reset Values

`onecontext.ResetValues` reset the values of an existing context.
