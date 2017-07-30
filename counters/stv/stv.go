package stv

import (
	"github.com/shawntoffel/electioncounter/counters"
	"github.com/shawntoffel/electioncounter/election"
)

type StvState struct {
	counters.CounterState
	Quota    int64
	Round    StvRound
	NumSeats int
	Ballots  election.RolledUpBallots
}

type StvRound struct {
	Number     int
	AnyElected bool
}
