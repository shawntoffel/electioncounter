package stv

import (
	"github.com/shawntoffel/electioncounter/counters"
	"github.com/shawntoffel/electioncounter/election"
)

type State struct {
	counters.CounterState
	Quota    int64
	Round    int
	NumSeats int
	Ballots  election.RolledUpBallots
}
