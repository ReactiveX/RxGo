## Categories of operators

There are operators for different purposes, and they may be categorized as: creation, transformation, filtering, joining, multicasting, error handling, utility, etc. In the following list you will find all the operators organized in categories.

## Creation Operators

<!-- - bindCallback -->
<!-- - bindNodeCallback -->
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
- Timer ✅
- [Iif](./iif.md) ✅

## Join Creation Operators

> These are Observable creation operators that also have join functionality -- emitting values of multiple source Observables.

<!-- - Partition -->

- ConcatAll ✅
- ConcatWith ✅
- CombineLatestAll ✅
- CombineLatestWith ✅
- ExhaustAll
- ForkJoin ✅
- MergeAll 🚧
- MergeWith 🚧
- RaceWith 🚧
- ZipAll ✅
- ZipWith ✅
- SwitchAll
- startWith
- WithLatestFrom

## Transformation Operators

- Buffer ✅
- BufferCount ✅ 📝
- BufferTime ✅
- BufferToggle ✅
- BufferWhen ✅
- ConcatMap ✅
- ExhaustMap ✅
- Expand
- GroupBy 🚧
- Map ✅ 📝
- MergeMap 🚧
- MergeScan
- Pairwise ✅
- Scan ✅
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
- DistinctUntilChanged ✅ 📝
- ElementAt ✅ 📝
- Filter ✅ 📝
- First ✅ 📝
- IgnoreElements ✅ 📝
- Last ✅
- Sample ✅
- SampleTime ✅
- Single ✅
- Skip ✅
- SkipLast ✅
- SkipUntil ✅
- SkipWhile ✅
- Take ✅
- TakeLast ✅
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
- RetryWhen 👎

## Utility Operators

- Do ✅
- Delay ✅
- DelayWhen 🚧
- Dematerialize ✅
- Materialize ✅
- ObserveOn
- SubscribeOn
- Repeat ✅
- RepeatWhen 👎
- TimeInterval ✅
- Timestamp ✅ 📝
- Timeout ✅
- TimeoutWith 👎
- [ToSlice](./to-slice.md) ✅ 📝

## Conditional and Boolean Operators

- DefaultIfEmpty ✅
- Every ✅
- Find ✅
- FindIndex ✅
- IsEmpty ✅
- SequenceEqual ✅
- ThrowIfEmpty ✅

## Mathematical and Aggregate Operators

- Count ✅
- Max ✅
- Min ✅
- Reduce ✅
