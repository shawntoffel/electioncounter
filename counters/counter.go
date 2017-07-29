package counters

import (
	"github.com/shawntoffel/electioncounter/election"
)

type Counter interface {
	Initialize(config election.Config)
	CountInitialVotes()
	UpdateRound()
	HasEnded() bool
	Result() (*election.Result, error)
}

type CounterState struct {
	Error           error
	Events          election.Events
	ExpectedVersion int
}
