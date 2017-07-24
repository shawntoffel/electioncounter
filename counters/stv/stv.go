package stv

import (
	"github.com/shawntoffel/electioncounter/counters"
)

type StvCounter interface {
	counters.Counter
}

type Stv struct {
	Quota           int64
	Round           int
	NumberToElect   int
	Ballots         counters.RolledUpBallots
	Changes         counters.Events
	ExpectedVersion int
}
