# Tooth
**tooth** abstracts channel operations in Go by means of a pubsub-like methods.

## Usage
`Tooth` instance is used to exchange data within the code.

## Examples
```go

frodo := tooth.New("frodo")
sam := tooth.New("sam")
gollum := tooth.New("gollum")

sam.Subscribe(frodo)
frodo.Subscribe(gollum)
sam.Subscribe(gollum)

frodo.Publish("I'm so tired")
gollum.Publish("My preciousssss")

msg1 := sam.FetchAll()      // [I'm so tired My preciousssss]
msg2 := frodo.Fetch(gollum) // My preciousssss

```

Right now it's limited to only `string` but I'm working on it.
