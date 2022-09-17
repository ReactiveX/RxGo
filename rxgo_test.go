package rxgo

import (
	"testing"

	"github.com/stretchr/testify/require"
)

func checkObservableResult[T any](t *testing.T, obs Observable[T], result T, err error, isCompleted bool) {
	var (
		hasCompleted  bool
		collectedErr  error
		collectedData T
	)
	obs.SubscribeSync(func(v T) {
		collectedData = v
	}, func(err error) {
		collectedErr = err
	}, func() {
		hasCompleted = true
	})
	require.Equal(t, collectedData, result)
	require.Equal(t, hasCompleted, isCompleted)
	require.Equal(t, collectedErr, err)
}

func checkObservableResultWithAnyError[T any](t *testing.T, obs Observable[T], result T, err []error, isCompleted bool) {
	var (
		hasCompleted  bool
		collectedErr  error
		collectedData T
	)
	obs.SubscribeSync(func(v T) {
		collectedData = v
	}, func(err error) {
		collectedErr = err
	}, func() {
		hasCompleted = true
	})
	require.Equal(t, collectedData, result)
	require.Equal(t, hasCompleted, isCompleted)
	require.Contains(t, err, collectedErr)
}

func checkObservableHasResults[T any](t *testing.T, obs Observable[T], hasResult bool, err error, isCompleted bool) {
	var (
		hasCompleted  bool
		collectedErr  error
		collectedData = make([]T, 0)
	)
	obs.SubscribeSync(func(v T) {
		collectedData = append(collectedData, v)
	}, func(err error) {
		collectedErr = err
	}, func() {
		hasCompleted = true
	})
	if hasResult {
		require.True(t, len(collectedData) > 0)
	} else {
		require.True(t, len(collectedData) == 0)
	}
	require.Equal(t, hasCompleted, isCompleted)
	require.Equal(t, collectedErr, err)
}

func checkObservableResults[T any](t *testing.T, obs Observable[T], result []T, err error, isCompleted bool) {
	var (
		hasCompleted  bool
		collectedErr  error
		collectedData = make([]T, 0, len(result))
	)
	obs.SubscribeSync(func(v T) {
		collectedData = append(collectedData, v)
	}, func(err error) {
		collectedErr = err
	}, func() {
		hasCompleted = true
	})
	require.ElementsMatch(t, collectedData, result)
	require.Equal(t, hasCompleted, isCompleted)
	require.Equal(t, collectedErr, err)
}
