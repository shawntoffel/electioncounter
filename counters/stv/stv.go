package stv

import (
	"github.com/shawntoffel/electioncounter/counters"
	"github.com/shawntoffel/electioncounter/election"
)

type Stv struct {
	counters.CounterState
	Quota         int64
	Round         int
	NumberToElect int
	Ballots       election.RolledUpBallots
}
