package meek

import (
	"github.com/shawntoffel/electioncounter/counters"
	"github.com/shawntoffel/electioncounter/counters/stv"
)

type MeekStvCounter interface {
	stv.StvCounter
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
		HandleEvent(&m, event)
		m.ExpectedVersion++
	}
	return &m
}

func (state *meekStvCounter) Create(counterConfig counters.CounterConfig) {
	createCount(state, counterConfig)

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
