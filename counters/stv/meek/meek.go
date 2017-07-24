package meek

import (
	"github.com/shawntoffel/electioncounter/counters"
	"github.com/shawntoffel/electioncounter/counters/stv"
)

type MeekStvCounter interface {
	stv.StvCounter
	HandleEvent(event MeekEvent)
}

type meekStvCounter struct {
	stv.Stv
	Pool      Pool
	Precision int
	Scale     int64
}

func NewMeekStvCounter(events []MeekEvent) MeekStvCounter {
	m := meekStvCounter{}
	m.Pool = NewPool(NewMemoryStorage())
	for _, event := range events {
		m.HandleEvent(event)
		m.ExpectedVersion++
	}
	return &m
}

func (state *meekStvCounter) Create(counterConfig counters.CounterConfig) {
	createCount := CreateCount{}
	createCount.Candidates = counterConfig.Candidates
	createCount.Ballots = counterConfig.Ballots
	createCount.Precision = counterConfig.Precision
	createCount.NumberToElect = counterConfig.NumberToElect

	state.HandleEvent(&createCount)
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

	description := event.Transition(state)

	counterEvent := counters.Event{}
	counterEvent.Description = description

	state.Changes = append(state.Changes, counterEvent)
}
