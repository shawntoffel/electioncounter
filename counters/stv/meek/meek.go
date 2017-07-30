package meek

import (
	"github.com/shawntoffel/electioncounter/counters"
	"github.com/shawntoffel/electioncounter/counters/stv/meek/event"
	"github.com/shawntoffel/electioncounter/election"
)

type MeekStvCounter interface {
	counters.Counter
}

type meekStvCounter struct {
	Meek event.MeekEventProcessor
}

func NewMeekStvCounter(history []event.MeekEvent) MeekStvCounter {
	m := meekStvCounter{}
	m.Meek = event.NewMeekEventProcessor(history)
	return &m
}

func (state *meekStvCounter) Initialize(config election.Config) {
	state.Meek.Create(config)
	state.Meek.ExcludeWithdrawnCandidates(config.WithdrawnCandidates)
	state.Meek.PerformPreliminaryElection()
}

func (state *meekStvCounter) UpdateRound() {
	state.Meek.IncrementRound()

	for {
		state.Meek.DistributeVotes()

		if state.Meek.RoundHasEnded() {
			break
		}
	}

	if !state.Meek.HasEnded() {
		state.Meek.ExcludeLowestCandidate()
	}
}

func (state *meekStvCounter) HasEnded() bool {
	return state.Meek.HasEnded()
}

func (state *meekStvCounter) Result() (*election.Result, error) {
	state.Meek.ExcludeRemainingCandidates()

	changes, err := state.Meek.Changes()

	result := election.Result{}
	result.Events = changes

	return &result, err
}
