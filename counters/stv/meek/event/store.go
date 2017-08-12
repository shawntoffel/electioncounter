package event

type EventStore interface {
	AddEvent(event MeekEvent)
}

type eventStore struct {
	Events   []MeekEvent
	Consumer Consumer
}

func NewEventStore(consumer Consumer) EventStore {
	return &eventStore{[]MeekEvent{}, consumer}
}

func (e *eventStore) AddEvent(event MeekEvent) {
	e.Events = append(e.Events, event)

	e.Consumer.ProcessEvent(event)
}
