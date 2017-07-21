package counters

type Aggregate interface {
	HandleEvent(event interface{})
}

type aggregate struct {
	ExpectedVersion int
	Changes         []interface{}
}
