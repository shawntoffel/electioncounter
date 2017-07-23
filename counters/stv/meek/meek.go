package meek

import (
	"container/list"
	"github.com/shawntoffel/electioncounter/counters"
)

type MeekStvCounter interface {
	StvCounter
	HandleEvent(event MeekEvent)
}

type meekStvCounter struct {
	stvCounter
	Pool      Pool
	Precision int
	Scaler    int64
}

type MeekEvent interface {
	Transition(m *meekStvCounter)
	Describe() string
}

func NewMeekStvCounter(events []MeekEvent) MeekStvCounter {
	m := meekStvCounter{}
	for _, event := range events {
		m.HandleEvent(event)
		m.ExpectedVersion++
	}
	return &m
}

func (state *meekStvCounter) Create(counterConfig counters.CounterConfig) {

}

func (state *meekStvCounter) UpdateRound() {}

func (state *meekStvCounter) HasEnded() bool {
	return true
}

func (state *meekStvCounter) Result() counters.Result {
	result := counters.Result{}
	result.Events = state.Changes

	return result
}

func (state *meekStvCounter) HandleEvent(event MeekEvent) {
	counterEvent := counters.Event{}
	counterEvent.Description = event.Describe()

	state.Changes = append(state.Changes, counterEvent)

	event.Transition(state)
}
