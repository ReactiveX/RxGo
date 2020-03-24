# FromChannel Operator

## Overview

Create a cold observable from a channel.

The items are consumed when an Observer subscribes.

## Example

```go
ch := make(chan rxgo.Item)
observable := rxgo.FromChannel(ch)
```

* [WithPublishStrategy](options.md#withpublishstrategy)