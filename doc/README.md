## Categories of operators

There are operators for different purposes, and they may be categorized as: creation, transformation, filtering, joining, multicasting, error handling, utility, etc. In the following list you will find all the operators organized in categories.

## Creation Operators

<!-- - bindCallback -->
<!-- - bindNodeCallback -->
<!-- - from -->
<!-- - fromEventPattern -->
<!-- - generate -->
<!-- - of -->
- Defer ✅
- EMPTY ✅
- Interval ✅
- NEVER ✅
- Range ✅
- ThrowError ✅
- Timer ✅
- Iif ✅

## Join Creation Operators

> These are Observable creation operators that also have join functionality -- emitting values of multiple source Observables.

<!-- - Concat -->
<!-- - Partition -->
- CombineLatestWith ✅
- ForkJoin ✅
- MergeWith 🚧
- RaceWith 🚧
- Zip 🚧
- combineLatestAll
- concatAll
- exhaustAll
- mergeAll
- switchAll
- startWith
- withLatestFrom

## Transformation Operators

- Buffer 🚧
- BufferCount 🚧
- BufferTime 🚧
- BufferToggle 🚧
- BufferWhen 🚧
- ConcatMap ✅
- ExhaustMap ✅
- Expand
- GroupBy
- Map ✅
- MergeMap
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

- Audit
- AuditTime
- Debounce ✅
- DebounceTime ✅
- Distinct ✅
- DistinctUntilChanged ✅
- ElementAt ✅
- Filter ✅
- First ✅
- IgnoreElements ✅
- Last ✅
- Sample
- SampleTime
- Single
- Skip ✅
- SkipLast ✅
- SkipUntil ✅
- SkipWhile ✅
- Take ✅
- TakeLast ✅
- TakeUntil ✅
- TakeWhile ✅
- Throttle
- ThrottleTime

## Multicasting Operators

- Multicast
- Publish
- PublishBehavior
- PublishLast
- PublishReplay
- Share

## Error Handling Operators

- CatchError
- Retry
- RetryWhen

## Utility Operators

- Tap ✅
- Delay ✅
- DelayWhen
- Dematerialize
- Materialize ✅
- observeOn
- subscribeOn
- TimeInterval ✅
- Timestamp ✅
- Timeout ✅
- TimeoutWith 
- ToArray ✅

## Conditional and Boolean Operators

- DefaultIfEmpty ✅
- Every ✅
- Find ✅
- FindIndex ✅
- IsEmpty ✅

## Mathematical and Aggregate Operators

- Count ✅
- Max ✅
- Min ✅
- Reduce ✅