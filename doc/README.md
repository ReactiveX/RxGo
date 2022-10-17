## Categories of operators

There are operators for different purposes, and they may be categorized as: creation, transformation, filtering, joining, multicasting, error handling, utility, etc. In the following list you will find all the operators organized in categories.

## Creation Operators

<!-- - from -->
<!-- - fromEventPattern -->
<!-- - generate -->

- Of ✅
- [Defer](./defer.md) ✅ 📝
- [Empty](./empty.md) ✅ 📝
- [Interval](./interval.md) ✅ 📝
- [Never](./never.md) ✅ 📝
- [Range](./range.md) ✅ 📝
- [Throw](./throw.md) ✅ 📝
- [Timer](./timer.md) ✅ 📝
- [Iif](./iif.md) ✅ 📝

## Join Creation Operators

> These are Observable creation operators that also have join functionality -- emitting values of multiple source Observables.

<!-- - Partition -->

- [ConcatAll](./concat-all.md) ✅
- [ConcatWith](./concat-with.md) ✅
- CombineLatestAll ✅
- CombineLatestWith ✅
- ExhaustAll
- [ForkJoin](./fork-join.md) ✅ 📝
- MergeAll 🚧
- MergeWith 🚧
- RaceWith 🚧
- [ZipAll] ✅
- [ZipWith] ✅
- SwitchAll
- StartWith
- WithLatestFrom

## Transformation Operators

- [Buffer](./buffer.md) ✅
- [BufferCount](./buffer-count.md) ✅ 📝
- [BufferTime](./buffer-time.md) ✅
- BufferToggle ✅
- BufferWhen ✅
- [ConcatMap](./concat-map.md) ✅
- ExhaustMap ✅
- Expand
- GroupBy 🚧
- [Map](./map.md) ✅ 📝
- MergeMap 🚧
- MergeScan
- Pairwise ✅
- [Scan](./scan.md) ✅
- SwitchScan
- SwitchMap
- Window
- WindowCount
- WindowTime
- WindowToggle
- WindowWhen

## Filtering Operators

- Audit ✅
- AuditTime ✅
- Debounce ✅
- DebounceTime ✅
- Distinct ✅
- [DistinctUntilChanged](./distinct-until-changed.md) ✅ 📝
- [ElementAt](./element-at.md) ✅ 📝
- [Filter](./filter.md) ✅ 📝
- [First](./first.md) ✅ 📝
- [IgnoreElements](./ignore-elements.md) ✅ 📝
- [Last](./last.md) ✅ 📝
- Sample ✅
- SampleTime ✅
- [Single](./single.md) ✅ 📝
- [Skip](./skip.md) ✅ 📝
- [SkipLast](./skiplast.md) ✅ 📝
- SkipUntil ✅
- [SkipWhile](./skip-while.md) ✅ 📝
- [Take](./take.md) ✅ 📝
- [TakeLast](./takelast.md) ✅ 📝
- TakeUntil ✅
- TakeWhile ✅
- Throttle 🚧
- ThrottleTime 🚧

## Multicasting Operators

- Multicast
- Publish
- PublishBehavior
- PublishLast
- PublishReplay
- Share

## Error Handling Operators

- Catch ✅
- Retry ✅
- ~~RetryWhen~~

## Utility Operators

- [Do](./do.md) ✅
- [Delay](./delay.md) ✅
- [DelayWhen](./delay-when.md) 🚧
- [Dematerialize](./dematerialize.md) ✅ 📝
- [Materialize](./materialize.md) ✅ 📝
- ObserveOn
- SubscribeOn
- [Repeat](./repeat.md) ✅
- ~~RepeatWhen~~
- [TimeInterval](./time-interval.md) ✅
- [Timestamp](./timestamp.md) ✅ 📝
- [Timeout](./timeout.md) ✅
- ~~TimeoutWith~~
- [ToSlice](./to-slice.md) ✅ 📝

## Conditional and Boolean Operators

- [DefaultIfEmpty](./default-if-empty.md) ✅ 📝
- [Every](./every.md) ✅ 📝
- [Find](./find.md) ✅
- [FindIndex](./find-index.md) ✅
- [IsEmpty](./is-empty.md) ✅ 📝
- [SequenceEqual](./sequence-equal.md) ✅ 📝
- [ThrowIfEmpty] ✅ 📝

## Mathematical and Aggregate Operators

- [Count](./count.md) ✅ 📝
- [Max](./max.md) ✅ 📝
- [Min](./min.md) ✅ 📝
- [Reduce](./reduce.md) ✅ 📝
