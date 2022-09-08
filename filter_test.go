package rxgo

import "testing"

func TestDistinct(t *testing.T) {
	t.Run("Distinct with EMPTY", func(t *testing.T) {
		checkObservableResult(t, Pipe1(EMPTY[any](), Distinct(func(value any) int {
			return value.(int)
		})), nil, nil, true)
	})

	t.Run("Distinct with numbers", func(t *testing.T) {
		checkObservableResults(t, Pipe1(Scheduled(1, 1, 2, 2, 2, 1, 2, 3, 4, 3, 2, 1), Distinct(func(value int) int {
			return value
		})), []int{1, 2, 3, 4}, nil, true)
	})

	t.Run("Distinct with struct", func(t *testing.T) {
		type user struct {
			name string
			age  uint
		}

		checkObservableResults(t, Pipe1(Scheduled(
			user{name: "Foo", age: 4},
			user{name: "Bar", age: 7},
			user{name: "Foo", age: 5},
		), Distinct(func(v user) string {
			return v.name
		})), []user{
			{age: 4, name: "Foo"},
			{age: 7, name: "Bar"},
		}, nil, true)
	})
}

func TestDistinctUntilChanged(t *testing.T) {
	t.Run("DistinctUntilChanged with empty", func(t *testing.T) {
		checkObservableResult(t, Pipe1(EMPTY[any](), DistinctUntilChanged[any]()), nil, nil, true)
	})

	t.Run("DistinctUntilChanged with String", func(t *testing.T) {
		checkObservableResults(t,
			Pipe1(Scheduled("a", "a", "b", "a", "c", "c", "d"), DistinctUntilChanged[string]()),
			[]string{"a", "b", "a", "c", "d"}, nil, true)
	})

	t.Run("DistinctUntilChanged with Number", func(t *testing.T) {
		checkObservableResults(t,
			Pipe1(
				Scheduled(30, 31, 20, 34, 33, 29, 35, 20),
				DistinctUntilChanged(func(prev, current int) bool {
					return current <= prev
				}),
			), []int{30, 31, 34, 35}, nil, true)
	})

	t.Run("DistinctUntilChanged with Struct", func(t *testing.T) {
		type build struct {
			engineVersion       string
			transmissionVersion string
		}
		checkObservableResults(t,
			Pipe1(
				Scheduled(
					build{engineVersion: "1.1.0", transmissionVersion: "1.2.0"},
					build{engineVersion: "1.1.0", transmissionVersion: "1.4.0"},
					build{engineVersion: "1.3.0", transmissionVersion: "1.4.0"},
					build{engineVersion: "1.3.0", transmissionVersion: "1.5.0"},
					build{engineVersion: "2.0.0", transmissionVersion: "1.5.0"},
				),
				DistinctUntilChanged(func(prev, curr build) bool {
					return (prev.engineVersion == curr.engineVersion ||
						prev.transmissionVersion == curr.transmissionVersion)
				}),
			),
			[]build{
				{engineVersion: "1.1.0", transmissionVersion: "1.2.0"},
				{engineVersion: "1.3.0", transmissionVersion: "1.4.0"},
				{engineVersion: "2.0.0", transmissionVersion: "1.5.0"},
			}, nil, true)
	})

	t.Run("DistinctUntilChanged with Struct(complex)", func(t *testing.T) {
		type account struct {
			updatedBy string
			data      []string
		}
		checkObservableResults(t,
			Pipe1(
				Scheduled(
					account{updatedBy: "blesh", data: []string{}},
					account{updatedBy: "blesh", data: []string{}},
					account{updatedBy: "jamieson"},
					account{updatedBy: "jamieson"},
					account{updatedBy: "blesh"},
				),
				DistinctUntilChanged[account](),
			),
			[]account{
				{updatedBy: "blesh", data: []string{}},
				{updatedBy: "jamieson"},
				{updatedBy: "blesh"},
			}, nil, true)
	})
}

func TestFilter(t *testing.T) {

}
