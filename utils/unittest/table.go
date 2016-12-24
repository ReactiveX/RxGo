package unittest

// Table is a generic test table for unittesting
type Table struct {
	Actual   interface{}
	Expected interface{}
}

// Tables is a slice of Tables
type Tables []Table
