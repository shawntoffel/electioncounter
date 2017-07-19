package main

type Aggregate struct {
	ExpectedVersion int
	Changes         []interface{}
}
