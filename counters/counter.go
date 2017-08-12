package counters

import (
	"github.com/shawntoffel/electioncounter/election"
)

type Counter interface {
	Initialize(config election.Config)
	UpdateRound()
	HasEnded() bool
	Result() (*election.Result, error)
}

type CounterState struct {
	Events          election.Events
	ExpectedVersion int
	Error           error
}
