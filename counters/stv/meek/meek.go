package meek

import (
	"github.com/shawntoffel/electioncounter/counters"
	"github.com/shawntoffel/electioncounter/counters/stv/meek/event"
	"github.com/shawntoffel/electioncounter/election"
	"fmt"
)

type MeekStvCounter interface {
	counters.Counter
}

type meekStvCounter struct {
	Meek event.Commander
}

func NewMeekStvCounter(history []event.MeekEvent) MeekStvCounter {
	m := meekStvCounter{}
	m.Meek = event.NewCommander()
	return &m
}

func (state *meekStvCounter) Initialize(config election.Config) {
	state.Meek.Create(config)
	state.Meek.ExcludeWithdrawnCandidates(config.WithdrawnCandidates)
	state.Meek.PerformPreliminaryElection()
}

func (state *meekStvCounter) UpdateRound() {
    fmt.Println("updating round============================================================")
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

	status := state.Meek.Status()

	if status.Error != nil {
		return nil, status.Error
	}

	elected := election.Candidates{}

	for _, candidate := range status.Elected {
		c := election.Candidate{}
		c.Id = candidate.Id
		c.Name = candidate.Name

		elected = append(elected, c)
	}

	result := election.Result{}
	result.Events = status.Events
	result.Candidates = elected

	return &result, nil
}
