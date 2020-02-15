package rxgo

//func Test_NewObserver(t *testing.T) {
//	count := 0
//	var gotErr error
//	done := false
//
//	observer := NewObserver(func(i interface{}) {
//		count += i.(int)
//	}, func(err error) {
//		gotErr = err
//	}, func() {
//		done = true
//	})
//
//	observer.OnNext(3)
//	observer.OnNext(2)
//	observer.OnNext(5)
//	expectedErr := errors.New("foo")
//	observer.OnError(expectedErr)
//	observer.OnDone()
//	observer.Block()
//
//	assert.Equal(t, 10, count)
//	assert.Equal(t, expectedErr, gotErr)
//	assert.True(t, done)
//}
