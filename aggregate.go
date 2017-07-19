package main

type Aggregate interface {
	HandleEvent(event interface{})
}

type aggregate struct {
	ExpectedVersion int
	Changes         []interface{}
}
