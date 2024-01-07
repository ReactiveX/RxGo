## Subjects
A detailed description of Subjects can be found here [Subjects](https://reactivex.io/documentation/subject.html).

Available Subject types:
* Subject - a simple fan-out with the ability to subscribe and unsubscribe any time
* BehaviorSubject - a subject which replays the last published item to every new subscriber
* ReplaySubject - a subject which replays the last n published items to every new subscriber

### Design
Subjects are created with a set of Observable options. Every subject subscriber receives a Subscription and an Observable Object. The Subscription can be used to unsubscribe from the Subject. The Observable is used to receive items from the Subject. Each Observable is a cold Observable with its own event source channel.

> [!NOTE]  
> Even though Subjects accept all options to create new Subscriber Observables, not all combinations make sense. For example, because each observer has its own Observable there is no point using connectable Observer options.

### Simple Subject Example
The code below shows a simple Subject example:
```go
// new subject with default options
subject := NewSubject()

sub, obs := subject.Subscribe()
obs.DoOnNext(func(i interface{}) {
    // handle items
})

// call unsubscribe when done
sub.Unsubscribe()
```

### Subject with BackPressure Strategy
By default a slow Subscriber would block all other Subscribers. This can be changed by creating Subscribers with BackPressure Strategy Drop:
```go
subject := NewSubject(WithBackPressureStrategy(Drop))
```

### Behavior and Replay Subject Design
Both Behavior and Replay Subjects publish one or more stored items to new subscribers before publishing new items. To achieve this, these subjects use buffered Go channels to create hot Observables. The buffer size equals to max number of replay items. This is to ensure that new Observers can create Subscriptions without a deadlock or blocking out other Subscribers. The entire Subject will be locked until a new Subscriber consumed all replay items.

> [!NOTE]  
> Even though Behavior and Replay Subjects accept all options to create new Subscriber Observables, not all combinations make sense. BackPressure strategy Drop should only used with care.

### Replay Subject Construction
The ReplaySubject constructor has an additional parameter "maxReplayItems". This parameter controls how many items are held in buffer for new subscribers.

